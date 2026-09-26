package sockets

import (
	"encoding/binary"
	"errors"
	"io"
	"net"
	"sync"
	"time"

	manager_io "github.com/456vv/x/vweb_dynamic/internal/wazero/manager/io"
)

const defaultUDPBufferSize = 256 // 默认缓冲 256 个数据报

var udpReadBufPool = sync.Pool{
	New: func() any {
		b := make([]byte, 65535)
		return &b
	},
}

// cloneIPSocketAddress 深拷贝地址，避免与 guest 复用结构体指针形成 data race。
func cloneIPSocketAddress(a IPSocketAddress) IPSocketAddress {
	var out IPSocketAddress
	if a.IPV4 != nil {
		v := *a.IPV4
		out.IPV4 = &v
	}
	if a.IPV6 != nil {
		v := *a.IPV6
		out.IPV6 = &v
	}
	return out
}

// --- Asynchronous UDP Reader ---

type AsyncUDPReader struct {
	conn          *net.UDPConn
	buffer        []IncomingDatagram
	mutex         sync.Mutex
	cond          *sync.Cond // 队列满时让 run 等待，Receive/Close 再唤醒，避免无界内存
	ready         *manager_io.ChannelPollable
	done          chan struct{}
	err           error
	once          sync.Once
	maxBufferSize int
	exited        chan struct{}
}

func NewAsyncUDPReader(conn *net.UDPConn) *AsyncUDPReader {
	wrapper := &AsyncUDPReader{
		conn:          conn,
		buffer:        make([]IncomingDatagram, 0, defaultUDPBufferSize),
		ready:         manager_io.NewPollable(nil),
		done:          make(chan struct{}),
		maxBufferSize: defaultUDPBufferSize,
		exited:        make(chan struct{}),
	}
	wrapper.cond = sync.NewCond(&wrapper.mutex)
	// nil *UDPConn 会在后台 ReadFromUDP 解引用 panic
	if conn == nil {
		wrapper.err = net.ErrClosed
		close(wrapper.done)
		close(wrapper.exited)
		wrapper.ready.SetReady()
		return wrapper
	}
	go wrapper.run()
	return wrapper
}

func (ar *AsyncUDPReader) waitExit() {
	if ar == nil || ar.exited == nil {
		return
	}
	<-ar.exited
}

// WaitExit 供 wasip2 在替换/drop stream 时等待后台 goroutine。
func (ar *AsyncUDPReader) WaitExit() { ar.waitExit() }

func (ar *AsyncUDPReader) closed() bool {
	select {
	case <-ar.done:
		return true
	default:
		return false
	}
}

func (ar *AsyncUDPReader) readableLocked() bool {
	return len(ar.buffer) > 0 || ar.err != nil || ar.closed()
}

func (ar *AsyncUDPReader) readable() bool {
	ar.mutex.Lock()
	defer ar.mutex.Unlock()
	return ar.readableLocked()
}

func (ar *AsyncUDPReader) run() {
	defer func() {
		close(ar.exited)
		// 退出时确保 pollable 就绪，避免订阅方永久等待
		ar.mutex.Lock()
		ar.ready.SetReady()
		ar.cond.Broadcast() // 退出路径也唤醒可能卡在 Wait 上的逻辑
		ar.mutex.Unlock()
	}()

	bufp := udpReadBufPool.Get().(*[]byte)
	buf := *bufp
	defer udpReadBufPool.Put(bufp) // 每条 UDP 流常驻 64KiB；退出时归还，降低 GC 压力

	for {
		ar.mutex.Lock()
		// 队列已满时不再收包，避免 buffer 无限增长导致 OOM
		for len(ar.buffer) >= ar.maxBufferSize && ar.err == nil && !ar.closed() {
			ar.cond.Wait() // Wait 会释放 mutex，被 Broadcast 后再抢锁
		}
		if ar.closed() {
			ar.mutex.Unlock()
			return
		}
		ar.mutex.Unlock()

		n, remoteAddr, err := ar.conn.ReadFromUDP(buf)

		ar.mutex.Lock()
		if err != nil {
			if ar.closed() {
				ar.mutex.Unlock()
				return
			}
			ar.err = err
			ar.ready.SetReady()
			ar.cond.Broadcast()
			ar.mutex.Unlock()
			return
		}

		// Close 之后迟到的成功收包不应再入队，否则 drop 后仍占内存
		if ar.closed() {
			ar.mutex.Unlock()
			return
		}

		// UDP 允许 0 长度数据报；原先 if n>0 会静默丢包。
		wasEmpty := len(ar.buffer) == 0
		data := make([]byte, n)
		if n > 0 {
			copy(data, buf[:n]) // buf 会被下一轮 Read 覆盖，必须拷贝
		}

		wasiAddr, factoryErr := ToIPSocketAddress(remoteAddr)
		if factoryErr != nil {
			// 地址转换失败不再静默丢包后继续；记录错误并唤醒消费者
			ar.err = factoryErr
			ar.ready.SetReady()
			ar.cond.Broadcast()
			ar.mutex.Unlock()
			return
		}
		ar.buffer = append(ar.buffer, IncomingDatagram{
			Data:          data,
			RemoteAddress: wasiAddr,
		})
		if wasEmpty {
			ar.ready.SetReady()
		}
		ar.mutex.Unlock()
	}
}

func (ar *AsyncUDPReader) Receive(maxResults uint64) ([]IncomingDatagram, error) {
	ar.mutex.Lock()
	defer ar.mutex.Unlock()

	if len(ar.buffer) > 0 {
		count := uint64(len(ar.buffer))
		if count > maxResults {
			count = maxResults
		}

		datagrams := make([]IncomingDatagram, count)
		copy(datagrams, ar.buffer[:count])
		ar.buffer = ar.buffer[count:]

		// 消费后队列出现空位，必须 Broadcast 叫醒因“满”而 Wait 的 run
		ar.cond.Broadcast()

		// 长期运行后 buffer[count:] 仍占着原 backing array，收缩以免内存只增不减
		if cap(ar.buffer) > ar.maxBufferSize*2 && len(ar.buffer) <= ar.maxBufferSize/4 {
			nbuf := make([]IncomingDatagram, len(ar.buffer), ar.maxBufferSize)
			copy(nbuf, ar.buffer)
			ar.buffer = nbuf
		}

		if len(ar.buffer) == 0 && ar.err == nil && !ar.closed() {
			ar.ready.Reset()
		}
		return datagrams, nil
	}
	if ar.err != nil {
		// 关闭后应让上层映射为 InvalidState，而不是空列表被 guest 当成“暂无数据”空转。
		if ar.err == io.EOF {
			return nil, net.ErrClosed
		}
		return nil, ar.err
	}
	return nil, nil
}

func (ar *AsyncUDPReader) Subscribe() manager_io.IPollable {
	ar.mutex.Lock()
	if ar.readableLocked() {
		ar.ready.SetReady()
	} else {
		ar.ready.Reset()
		if ar.readableLocked() {
			ar.ready.SetReady()
		}
	}
	ar.mutex.Unlock()
	return manager_io.NewLevelPollable(ar.readable, ar.ready)
}

func (ar *AsyncUDPReader) Close() {
	ar.once.Do(func() {
		close(ar.done)
		ar.mutex.Lock()
		// 关闭视为流结束，Subscribe 看到 err 就不会 Reset
		if ar.err == nil {
			ar.err = io.EOF
		}
		ar.ready.SetReady()
		ar.cond.Broadcast() // run 可能卡在“队列满”的 Wait 上，不唤醒就会泄漏 goroutine
		ar.mutex.Unlock()
		if ar.conn != nil {
			// run 可能堵在 ReadFromUDP；仅 close(done) 不能打断
			ar.conn.SetReadDeadline(time.Now())
		}
	})
}

// --- Asynchronous UDP Writer with Backpressure ---

type AsyncUDPWriter struct {
	conn          *net.UDPConn
	buffer        []OutgoingDatagram
	mutex         sync.Mutex
	cond          *sync.Cond
	ready         *manager_io.ChannelPollable
	done          chan struct{}
	err           error
	maxBufferSize int
	once          sync.Once
	closed        bool // Close 后禁止再 Send，避免往已唤醒退出的 run 里塞包
	exited        chan struct{}
}

func NewAsyncUDPWriter(conn *net.UDPConn) *AsyncUDPWriter {
	wrapper := &AsyncUDPWriter{
		conn:          conn,
		buffer:        make([]OutgoingDatagram, 0, defaultUDPBufferSize),
		ready:         manager_io.NewPollable(nil),
		done:          make(chan struct{}),
		maxBufferSize: defaultUDPBufferSize,
		exited:        make(chan struct{}),
	}
	wrapper.cond = sync.NewCond(&wrapper.mutex)
	// nil *UDPConn 会在后台 Write/WriteToUDP 解引用 panic；cond 必须非 nil，Close 会 Broadcast
	if conn == nil {
		wrapper.err = net.ErrClosed
		wrapper.closed = true
		close(wrapper.done)
		close(wrapper.exited)
		return wrapper
	}
	wrapper.ready.SetReady()
	go wrapper.run()
	return wrapper
}

func (aw *AsyncUDPWriter) waitExit() {
	if aw == nil || aw.exited == nil {
		return
	}
	<-aw.exited
}

// WaitExit 供 wasip2 在替换/drop stream 时等待后台 goroutine。
func (aw *AsyncUDPWriter) WaitExit() { aw.waitExit() }

func (aw *AsyncUDPWriter) writableLocked() bool {
	return aw.closed || aw.err != nil || len(aw.buffer) < aw.maxBufferSize
}

func (aw *AsyncUDPWriter) writable() bool {
	aw.mutex.Lock()
	defer aw.mutex.Unlock()
	return aw.writableLocked()
}

func (aw *AsyncUDPWriter) run() {
	defer close(aw.exited)
	// 「持锁等待 → 解锁发送 → 再加锁」，避免持锁做网络 IO。
	for {
		aw.mutex.Lock()
		for len(aw.buffer) == 0 {
			select {
			case <-aw.done:
				aw.mutex.Unlock()
				return
			default:
				aw.cond.Wait()
				select {
				case <-aw.done:
					if len(aw.buffer) == 0 {
						aw.mutex.Unlock()
						return
					}
				default:
				}
			}
		}

		datagramsToSend := make([]OutgoingDatagram, len(aw.buffer))
		copy(datagramsToSend, aw.buffer)
		aw.buffer = aw.buffer[:0]
		aw.mutex.Unlock()

		var writeErr error
		failAt := -1
		for i, dg := range datagramsToSend {
			var remoteAddr *net.UDPAddr
			if dg.RemoteAddress.Some != nil {
				remoteAddr, writeErr = FromIPSocketAddressToUDPAddr(*dg.RemoteAddress.Some)
				if writeErr != nil {
					failAt = i
					break
				}
			}
			var n int
			if remoteAddr != nil {
				n, writeErr = aw.conn.WriteToUDP(dg.Data, remoteAddr)
			} else {
				n, writeErr = aw.conn.Write(dg.Data)
			}
			// 空数据报 n==0 且 len==0 仍然是成功。
			if writeErr == nil && n < len(dg.Data) {
				writeErr = io.ErrShortWrite
			}
			if writeErr != nil {
				failAt = i
				break
			}
		}

		aw.mutex.Lock()
		if writeErr != nil {
			aw.err = writeErr
			// 批次中失败点及之后的数据报原先被丢弃
			if failAt >= 0 && failAt < len(datagramsToSend) && !aw.closed {
				aw.buffer = append(datagramsToSend[failAt:], aw.buffer...)
			}
		}
		if aw.writableLocked() {
			aw.ready.SetReady()
		}
		aw.cond.Broadcast()
		shouldExit := aw.err != nil
		aw.mutex.Unlock()

		if shouldExit {
			return
		}
	}
}

func (aw *AsyncUDPWriter) Send(datagrams []OutgoingDatagram) (uint64, error) {
	aw.mutex.Lock()
	defer aw.mutex.Unlock()

	if aw.closed {
		return 0, net.ErrClosed
	}
	if aw.err != nil {
		return 0, aw.err
	}
	available := aw.maxBufferSize - len(aw.buffer)
	if available <= 0 {
		return 0, nil
	}
	count := uint64(len(datagrams))
	if count > uint64(available) {
		count = uint64(available)
	}
	// guest/调用方可能复用 Data 缓冲区，必须拷贝后再入队
	for i := uint64(0); i < count; i++ {
		dg := datagrams[i]
		if len(dg.Data) > 0 {
			d := make([]byte, len(dg.Data))
			copy(d, dg.Data)
			dg.Data = d
		}
		if dg.RemoteAddress.Some != nil {
			cp := cloneIPSocketAddress(*dg.RemoteAddress.Some)
			dg.RemoteAddress.Some = &cp
		}
		aw.buffer = append(aw.buffer, dg)
	}
	if count > 0 {
		aw.cond.Signal()
	}
	if len(aw.buffer) >= aw.maxBufferSize {
		aw.ready.Reset()
	}
	return count, nil
}

// AvailableSpace 返回缓冲区中可用于写入的数据报数量。
func (aw *AsyncUDPWriter) AvailableSpace() uint64 {
	aw.mutex.Lock()
	defer aw.mutex.Unlock()
	if aw.closed || aw.err != nil {
		return 0
	}
	avail := aw.maxBufferSize - len(aw.buffer)
	if avail < 0 {
		return 0
	}
	return uint64(avail)
}

// ClosedOrErr 在 writer 已关闭或失败时返回错误，供 check-send 映射 InvalidState。
// 关闭后 AvailableSpace()==0 且 Subscribe 仍就绪，guest 会当成“暂无空间”空转。
func (aw *AsyncUDPWriter) ClosedOrErr() error {
	if aw == nil {
		return net.ErrClosed
	}
	aw.mutex.Lock()
	defer aw.mutex.Unlock()
	if aw.err != nil {
		return aw.err
	}
	if aw.closed {
		return net.ErrClosed
	}
	return nil
}

func (aw *AsyncUDPWriter) Subscribe() manager_io.IPollable {
	aw.mutex.Lock()
	if aw.writableLocked() {
		aw.ready.SetReady()
	} else {
		aw.ready.Reset()
		if aw.writableLocked() {
			aw.ready.SetReady()
		}
	}
	aw.mutex.Unlock()
	return manager_io.NewLevelPollable(aw.writable, aw.ready)
}

func (aw *AsyncUDPWriter) Close() {
	aw.once.Do(func() {
		aw.mutex.Lock()
		aw.closed = true // 先禁止再入队
		// 原先 cond.Wait 排空时 run 若堵在 WriteToUDP，destructor 尚未 Conn.Close 会死锁。
		// 只唤醒并 close(done)；用写 deadline 打断阻塞 Write，剩余包由 run 尽量发完或随 err 退出。
		select {
		case <-aw.done:
		default:
			close(aw.done)
		}
		aw.cond.Broadcast()
		aw.ready.SetReady()
		aw.mutex.Unlock()
		if aw.conn != nil {
			_ = aw.conn.SetWriteDeadline(time.Now())
		}
	})
}

func zoneFromScopeID(id uint32) string {
	if id == 0 {
		return ""
	}
	ifi, err := net.InterfaceByIndex(int(id))
	if err != nil {
		return ""
	}
	return ifi.Name
}

func scopeIDFromZone(zone string) uint32 {
	if zone == "" {
		return 0
	}
	ifi, err := net.InterfaceByName(zone)
	if err != nil {
		return 0
	}
	if ifi.Index < 0 {
		return 0
	}
	return uint32(ifi.Index)
}

// FromIPSocketAddressToUDPAddr 将 WIT 的 IPSocketAddress 转换为 Go 的 *net.UDPAddr。
func FromIPSocketAddressToUDPAddr(addr IPSocketAddress) (*net.UDPAddr, error) {
	if addr.IPV4 != nil {
		ip := make(net.IP, 4)
		copy(ip, addr.IPV4.Address[:]) // net.IP(addr[:]) 别名数组，后续改 Address 会污染 UDPAddr
		return &net.UDPAddr{IP: ip, Port: int(addr.IPV4.Port)}, nil
	}
	if addr.IPV6 != nil {
		ip := make(net.IP, 16)
		for i, part := range addr.IPV6.Address {
			binary.BigEndian.PutUint16(ip[i*2:], part)
		}
		return &net.UDPAddr{IP: ip, Port: int(addr.IPV6.Port), Zone: zoneFromScopeID(addr.IPV6.ScopeID)}, nil
	}
	return nil, errors.New("invalid ip-socket-address")
}

func ToIPSocketAddress(addr net.Addr) (IPSocketAddress, error) {
	switch tcpAddr := addr.(type) {
	case *net.TCPAddr:
		if ipv4 := tcpAddr.IP.To4(); ipv4 != nil {
			var wasiAddr IPv4Address
			copy(wasiAddr[:], ipv4)
			return IPSocketAddress{
				IPV4: &IPv4SocketAddress{
					Port:    uint16(tcpAddr.Port),
					Address: wasiAddr,
				},
			}, nil
		}

		if ipv6 := tcpAddr.IP.To16(); ipv6 != nil {
			var wasiAddr IPv6Address
			for i := 0; i < 8; i++ { // for i := range 8 需要 Go 1.22
				wasiAddr[i] = binary.BigEndian.Uint16(ipv6[i*2:])
			}
			return IPSocketAddress{
				IPV6: &IPv6SocketAddress{
					Port:    uint16(tcpAddr.Port),
					Address: wasiAddr,
					ScopeID: scopeIDFromZone(tcpAddr.Zone),
				},
			}, nil
		}
	case *net.UDPAddr:
		if ipv4 := tcpAddr.IP.To4(); ipv4 != nil {
			var wasiAddr IPv4Address
			copy(wasiAddr[:], ipv4)
			return IPSocketAddress{
				IPV4: &IPv4SocketAddress{
					Port:    uint16(tcpAddr.Port),
					Address: wasiAddr,
				},
			}, nil
		}

		if ipv6 := tcpAddr.IP.To16(); ipv6 != nil {
			var wasiAddr IPv6Address
			for i := 0; i < 8; i++ {
				wasiAddr[i] = binary.BigEndian.Uint16(ipv6[i*2:])
			}
			return IPSocketAddress{
				IPV6: &IPv6SocketAddress{
					Port:    uint16(tcpAddr.Port),
					Address: wasiAddr,
					ScopeID: scopeIDFromZone(tcpAddr.Zone),
				},
			}, nil
		}
	}
	return IPSocketAddress{}, errors.New("unsupported address type")
}
