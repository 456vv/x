//go:build !unix && !windows

package v0_2

import (
	"context"
	"errors"
	"net"
	"sync"
	"syscall"
	"time"

	manager_io "github.com/456vv/x/vweb_dynamic/internal/wazero/manager/io"
	"github.com/456vv/x/vweb_dynamic/internal/wazero/manager/sockets"
	wasip2_io "github.com/456vv/x/vweb_dynamic/internal/wazero/wasip2/io/v0_2"
	witgo "github.com/456vv/x/vweb_dynamic/internal/wazero/witgo"
)

func (i *tcpImpl) DropTCPSocket(_ context.Context, handle TCPSocket) {
	i.host.TCPSocketManager().Remove(handle)
}

func (i *tcpImpl) StartBind(_ context.Context, this TCPSocket, network Network, localAddress IPSocketAddress) witgo.Result[witgo.Unit, ErrorCode] {
	sock, ok := i.host.TCPSocketManager().Get(this)
	if !ok {
		return witgo.Err[witgo.Unit, ErrorCode](ErrorCodeInvalidArgument)
	}
	if code := i.checkNetwork(network); code != 0 {
		return witgo.Err[witgo.Unit, ErrorCode](code)
	}

	addr, err := fromIPSocketAddressToTCPAddr(localAddress)
	if err != nil {
		return witgo.Err[witgo.Unit, ErrorCode](ErrorCodeInvalidArgument)
	}

	pnetwork := "tcp"
	if sock.Family == sockets.IPAddressFamilyIPV6 {
		pnetwork = "tcp6"
	}

	err = sock.DoBind(func() error {
		// 本平台无原始 fd，只能 ListenTCP（会立即进入监听）；状态仍标 Bound，真正 Listening 在 StartListen
		listener, listenErr := net.ListenTCP(pnetwork, addr)
		if listenErr != nil {
			return listenErr
		}
		sock.Listener = listener
		return nil
	})
	if err != nil {
		if errors.Is(err, sockets.ErrInvalidSocketState) {
			return witgo.Err[witgo.Unit, ErrorCode](ErrorCodeInvalidState)
		}
		return witgo.Err[witgo.Unit, ErrorCode](mapOsError(err))
	}
	return witgo.Ok[witgo.Unit, ErrorCode](witgo.Unit{})
}

func (i *tcpImpl) FinishBind(_ context.Context, this TCPSocket) witgo.Result[witgo.Unit, ErrorCode] {
	sock, ok := i.host.TCPSocketManager().Get(this)
	if !ok {
		return witgo.Err[witgo.Unit, ErrorCode](ErrorCodeInvalidArgument)
	}
	if sock.GetState() != sockets.TCPStateBound {
		return witgo.Err[witgo.Unit, ErrorCode](ErrorCodeInvalidState)
	}
	return witgo.Ok[witgo.Unit, ErrorCode](witgo.Unit{})
}

func (i *tcpImpl) listenTCP(sock *sockets.TCPSocket) error {
	if sock.Listener == nil {
		return errors.New("tcp socket is not bound")
	}
	return nil
}

func (i *tcpImpl) connectTCP(sock *sockets.TCPSocket, remoteAddress IPSocketAddress) (*net.TCPConn, error) {
	addr, err := fromIPSocketAddressToTCPAddr(remoteAddress)
	if err != nil {
		return nil, err
	}
	network := "tcp"
	if sock.Family == sockets.IPAddressFamilyIPV6 {
		network = "tcp6"
	}
	var local *net.TCPAddr
	if sock.Listener != nil {
		if la, ok := sock.Listener.Addr().(*net.TCPAddr); ok {
			local = la
		}
	}
	// drop 时 ConnectContext 取消才能打断无原始 fd 平台上的拨号
	d := net.Dialer{LocalAddr: local}
	c, err := d.DialContext(sock.ConnectContext(), network, addr.String())
	if err != nil {
		return nil, err
	}
	tc, ok := c.(*net.TCPConn)
	if !ok {
		_ = c.Close()
		return nil, errors.New("dial did not return TCPConn")
	}
	return tc, nil
}

func (i *tcpImpl) acceptTCP(sock *sockets.TCPSocket) (*net.TCPConn, error) {
	if sock == nil || sock.Listener == nil {
		return nil, syscall.EINVAL
	}
	_ = sock.Listener.SetDeadline(time.Now().Add(time.Millisecond))
	conn, err := sock.Listener.AcceptTCP()
	_ = sock.Listener.SetDeadline(time.Time{})
	return conn, err
}

func (i *tcpImpl) localAddrFromFd(_ *sockets.TCPSocket) (net.Addr, error) {
	// tcp.go LocalAddress 在 Bound 未 listen 时会调用；本平台无 fd
	return nil, errors.New("local address from fd is not supported on this platform")
}

func (i *tcpImpl) setKeepAliveEnabledUnconnected(this TCPSocket, value int) witgo.Result[witgo.Unit, ErrorCode] {
	return witgo.Err[witgo.Unit, ErrorCode](ErrorCodeNotSupported)
}

func (i *tcpImpl) setReceiveBufferUnconnected(this TCPSocket, value uint64) witgo.Result[witgo.Unit, ErrorCode] {
	return witgo.Err[witgo.Unit, ErrorCode](ErrorCodeNotSupported)
}

func (i *tcpImpl) setSendBufferUnconnected(this TCPSocket, value uint64) witgo.Result[witgo.Unit, ErrorCode] {
	return witgo.Err[witgo.Unit, ErrorCode](ErrorCodeNotSupported)
}

func (i *tcpImpl) setBufferSize(this TCPSocket, recv bool, value uint64) witgo.Result[witgo.Unit, ErrorCode] {
	return witgo.Err[witgo.Unit, ErrorCode](ErrorCodeNotSupported)
}

func (i *tcpImpl) SetListenBacklogSize(ctx context.Context, this TCPSocket, value uint64) witgo.Result[witgo.Unit, ErrorCode] {
	return witgo.Err[witgo.Unit, ErrorCode](ErrorCodeNotSupported)
}

func (i *tcpImpl) KeepAliveEnabled(ctx context.Context, this TCPSocket) witgo.Result[bool, ErrorCode] {
	return witgo.Err[bool, ErrorCode](ErrorCodeNotSupported)
}

func (i *tcpImpl) KeepAliveIdleTime(ctx context.Context, this TCPSocket) witgo.Result[uint64, ErrorCode] {
	return witgo.Err[uint64, ErrorCode](ErrorCodeNotSupported)
}

func (i *tcpImpl) SetKeepAliveIdleTime(ctx context.Context, this TCPSocket, value uint64) witgo.Result[witgo.Unit, ErrorCode] {
	return witgo.Err[witgo.Unit, ErrorCode](ErrorCodeNotSupported)
}

func (i *tcpImpl) KeepAliveInterval(ctx context.Context, this TCPSocket) witgo.Result[uint64, ErrorCode] {
	return witgo.Err[uint64, ErrorCode](ErrorCodeNotSupported)
}

func (i *tcpImpl) SetKeepAliveInterval(ctx context.Context, this TCPSocket, value uint64) witgo.Result[witgo.Unit, ErrorCode] {
	return witgo.Err[witgo.Unit, ErrorCode](ErrorCodeNotSupported)
}

func (i *tcpImpl) KeepAliveCount(ctx context.Context, this TCPSocket) witgo.Result[uint32, ErrorCode] {
	return witgo.Err[uint32, ErrorCode](ErrorCodeNotSupported)
}

func (i *tcpImpl) SetKeepAliveCount(ctx context.Context, this TCPSocket, value uint32) witgo.Result[witgo.Unit, ErrorCode] {
	return witgo.Err[witgo.Unit, ErrorCode](ErrorCodeNotSupported)
}

func (i *tcpImpl) HopLimit(ctx context.Context, this TCPSocket) witgo.Result[uint8, ErrorCode] {
	return witgo.Err[uint8, ErrorCode](ErrorCodeNotSupported)
}

func (i *tcpImpl) SetHopLimit(ctx context.Context, this TCPSocket, value uint8) witgo.Result[witgo.Unit, ErrorCode] {
	return witgo.Err[witgo.Unit, ErrorCode](ErrorCodeNotSupported)
}

func (i *tcpImpl) ReceiveBufferSize(ctx context.Context, this TCPSocket) witgo.Result[uint64, ErrorCode] {
	return witgo.Err[uint64, ErrorCode](ErrorCodeNotSupported)
}

func (i *tcpImpl) SendBufferSize(ctx context.Context, this TCPSocket) witgo.Result[uint64, ErrorCode] {
	return witgo.Err[uint64, ErrorCode](ErrorCodeNotSupported)
}

// subscribeListen 在无 poll(2)/WSA 的平台上等待可 accept。
//  1. 不能 AcceptTCP：会从 backlog 取走连接，guest 再 accept 会丢连接。
//  2. js/plan9/wasip1 通常没有可用的原始 fd。
//  3. 若 Listener 实现了 syscall.Conn，RawConn.Read 会等到内核标记可读，
//     回调返回 true 且不读数据，连接仍留在队列里。
//  4. 否则立即就绪：Accept 侧用 deadline=now 映射 WouldBlock，避免永久阻塞 wasm 线程。
func (i *tcpImpl) subscribeListen(sock *sockets.TCPSocket) wasip2_io.Pollable {
	done := make(chan struct{})
	var once sync.Once
	p := manager_io.NewPollable(func() {
		once.Do(func() { close(done) })
	})
	handle := i.host.PollManager().Add(p)

	if sock == nil || sock.Listener == nil {
		p.SetReady()
		return handle
	}

	sc, ok := any(sock.Listener).(syscall.Conn)
	if !ok {
		p.SetReady()
		return handle
	}
	raw, err := sc.SyscallConn()
	if err != nil || raw == nil {
		p.SetReady()
		return handle
	}

	ln := sock.Listener
	go func() {
		waitDone := make(chan struct{})
		go func() {
			defer close(waitDone)
			_ = raw.Read(func(_ uintptr) bool {
				select {
				case <-done:
					return true
				default:
					p.SetReady()
					return true
				}
			})
		}()
		select {
		case <-waitDone:
		case <-done:
			// RawConn.Read 可能一直阻塞；用短 deadline 唤醒，再清掉，避免影响后续 Accept
			_ = ln.SetDeadline(time.Now())
			<-waitDone
			_ = ln.SetDeadline(time.Time{})
		}
	}()
	return handle
}
