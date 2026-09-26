package v0_2

import (
	"context"
	"errors"
	"fmt"
	"net"
	"time"

	manager_io "github.com/456vv/x/vweb_dynamic/internal/wazero/manager/io"
	"github.com/456vv/x/vweb_dynamic/internal/wazero/manager/sockets"
	"github.com/456vv/x/vweb_dynamic/internal/wazero/wasip2"
	wasip2_io "github.com/456vv/x/vweb_dynamic/internal/wazero/wasip2/io/v0_2"
	witgo "github.com/456vv/x/vweb_dynamic/internal/wazero/witgo"
)

type tcpImpl struct {
	host *wasip2.Host
}

func newTCPImpl(h *wasip2.Host) *tcpImpl {
	return &tcpImpl{host: h}
}

func (i *tcpImpl) checkNetwork(network Network) ErrorCode {
	if i.host == nil || i.host.NetworkManager() == nil {
		return ErrorCodeInvalidArgument
	}
	if _, ok := i.host.NetworkManager().Get(network); !ok {
		// network 是 WASI 能力句柄；无效时不应 bind/connect。
		return ErrorCodeInvalidArgument
	}
	return 0
}

func (i *tcpImpl) StartConnect(ctx context.Context, this TCPSocket, network Network, remoteAddress IPSocketAddress) witgo.Result[witgo.Unit, ErrorCode] {
	sock, ok := i.host.TCPSocketManager().Get(this)
	if !ok {
		return witgo.Err[witgo.Unit, ErrorCode](ErrorCodeInvalidArgument)
	}

	if code := i.checkNetwork(network); code != 0 {
		return witgo.Err[witgo.Unit, ErrorCode](code)
	}
	if (remoteAddress.IPV4 != nil && sock.Family != sockets.IPAddressFamilyIPV4) ||
		(remoteAddress.IPV6 != nil && sock.Family != sockets.IPAddressFamilyIPV6) ||
		(remoteAddress.IPV4 == nil && remoteAddress.IPV6 == nil) {
		return witgo.Err[witgo.Unit, ErrorCode](ErrorCodeInvalidArgument)
	}
	okc, conflict := sock.BeginConnect()
	if conflict {
		return witgo.Err[witgo.Unit, ErrorCode](ErrorCodeConcurrencyConflict)
	}
	if !okc {
		return witgo.Err[witgo.Unit, ErrorCode](ErrorCodeInvalidState)
	}
	go func() {
		var conn *net.TCPConn
		var dialErr error
		defer func() {
			if rec := recover(); rec != nil {
				// connect 路径 panic 若不 StoreResult，subscribe 会永久阻塞；已拨通的 Conn 必须关掉
				if conn != nil {
					_ = conn.Close()
				}
				sock.StoreConnectResult(sockets.ConnectResult{Err: fmt.Errorf("tcp connect panic: %v", rec)})
			}
		}()
		conn, dialErr = i.connectTCP(sock, remoteAddress)
		sock.StoreConnectResult(sockets.ConnectResult{Conn: conn, Err: dialErr})
	}()
	return witgo.Ok[witgo.Unit, ErrorCode](witgo.Unit{})
}

func (i *tcpImpl) FinishConnect(ctx context.Context, this TCPSocket) witgo.Result[witgo.Tuple[wasip2_io.InputStream, wasip2_io.OutputStream], ErrorCode] {
	sock, ok := i.host.TCPSocketManager().Get(this)
	if !ok {
		return witgo.Err[witgo.Tuple[wasip2_io.InputStream, wasip2_io.OutputStream], ErrorCode](ErrorCodeInvalidArgument)
	}

	if sock.GetState() != sockets.TCPStateConnecting {
		return witgo.Err[witgo.Tuple[wasip2_io.InputStream, wasip2_io.OutputStream], ErrorCode](ErrorCodeNotInProgress)
	}

	// 不再从 channel 取出结果，避免与 Subscribe 互抢
	result, ready := sock.TryConnectResult()
	if !ready {
		return witgo.Err[witgo.Tuple[wasip2_io.InputStream, wasip2_io.OutputStream], ErrorCode](ErrorCodeWouldBlock)
	}
	if result.Err != nil {
		sock.ClaimConnectFailure()
		return witgo.Err[witgo.Tuple[wasip2_io.InputStream, wasip2_io.OutputStream], ErrorCode](mapOsError(result.Err))
	}
	if result.Conn == nil {
		sock.ClaimConnectFailure()
		return witgo.Err[witgo.Tuple[wasip2_io.InputStream, wasip2_io.OutputStream], ErrorCode](ErrorCodeUnknown)
	}

	conn := result.Conn
	// 从 channel 取出结果时，可能被其他 goroutine 先一步关闭了 socket，导致 Claim 失败
	if !sock.ClaimConnectSuccess(conn) {
		// drop 抢先 Closed 时 Claim 失败，result.Conn 可能尚未被析构拷走
		if conn != nil && sock.GetState() != sockets.TCPStateConnected {
			_ = conn.Close()
		}
		return witgo.Err[witgo.Tuple[wasip2_io.InputStream, wasip2_io.OutputStream], ErrorCode](ErrorCodeNotInProgress)
	}

	// Claim 释放锁后析构可能把 sock.Conn 置 nil；用局部 conn 绑流，避免空 reader 且把真正连接关掉无人读。
	inStream := manager_io.NewAsyncStreamForReader(conn, manager_io.DontCloseReader())
	inStreamHandle := i.host.StreamManager().Add(inStream)
	outStream := manager_io.NewAsyncStreamForWriter(conn, manager_io.DontCloseWriter())
	outStreamHandle := i.host.StreamManager().Add(outStream)

	return witgo.Ok[witgo.Tuple[wasip2_io.InputStream, wasip2_io.OutputStream], ErrorCode](witgo.Tuple[wasip2_io.InputStream, wasip2_io.OutputStream]{
		F0: inStreamHandle,
		F1: outStreamHandle,
	})
}

func (i *tcpImpl) StartListen(ctx context.Context, this TCPSocket) witgo.Result[witgo.Unit, ErrorCode] {
	sock, ok := i.host.TCPSocketManager().Get(this)
	if !ok {
		return witgo.Err[witgo.Unit, ErrorCode](ErrorCodeInvalidArgument)
	}
	err := sock.DoListen(func() error { return i.listenTCP(sock) })
	if err != nil {
		if errors.Is(err, sockets.ErrInvalidSocketState) {
			return witgo.Err[witgo.Unit, ErrorCode](ErrorCodeInvalidState)
		}
		return witgo.Err[witgo.Unit, ErrorCode](mapOsError(err))
	}
	return witgo.Ok[witgo.Unit, ErrorCode](witgo.Unit{})
}

func (i *tcpImpl) FinishListen(ctx context.Context, this TCPSocket) witgo.Result[witgo.Unit, ErrorCode] {
	sock, ok := i.host.TCPSocketManager().Get(this)
	if !ok {
		return witgo.Err[witgo.Unit, ErrorCode](ErrorCodeInvalidArgument)
	}
	if sock.GetState() != sockets.TCPStateListening {
		return witgo.Err[witgo.Unit, ErrorCode](ErrorCodeInvalidState)
	}
	return witgo.Ok[witgo.Unit, ErrorCode](witgo.Unit{})
}

func (i *tcpImpl) Accept(ctx context.Context, this TCPSocket) witgo.Result[witgo.Tuple3[TCPSocket, wasip2_io.InputStream, wasip2_io.OutputStream], ErrorCode] {
	sock, ok := i.host.TCPSocketManager().Get(this)
	if !ok {
		return witgo.Err[witgo.Tuple3[TCPSocket, wasip2_io.InputStream, wasip2_io.OutputStream], ErrorCode](ErrorCodeInvalidArgument)
	}
	if sock.GetState() != sockets.TCPStateListening {
		return witgo.Err[witgo.Tuple3[TCPSocket, wasip2_io.InputStream, wasip2_io.OutputStream], ErrorCode](ErrorCodeInvalidState)
	}
	// unix accept 可走原始 fd；原先强制 Listener!=nil，FileListener 失败或仅 listen(2) 成功时无法 accept。
	if sock.GetListener() == nil {
		if _, hasFd := sock.SnapshotFd(); !hasFd {
			return witgo.Err[witgo.Tuple3[TCPSocket, wasip2_io.InputStream, wasip2_io.OutputStream], ErrorCode](ErrorCodeInvalidState)
		}
	}

	// 走平台 acceptTCP（unix 可从原 fd Accept；deadline=Now 在队列非空时仍会立刻超时）
	conn, err := i.acceptTCP(sock)
	if err != nil {
		var netErr net.Error
		if errors.As(err, &netErr) && netErr.Timeout() {
			return witgo.Err[witgo.Tuple3[TCPSocket, wasip2_io.InputStream, wasip2_io.OutputStream], ErrorCode](ErrorCodeWouldBlock)
		}
		return witgo.Err[witgo.Tuple3[TCPSocket, wasip2_io.InputStream, wasip2_io.OutputStream], ErrorCode](mapOsError(err))
	}

	newSock := &sockets.TCPSocket{
		Conn:   conn,
		Family: sock.Family,
	}
	newSock.SetState(sockets.TCPStateConnected)
	newSockHandle := i.host.TCPSocketManager().Add(newSock)

	inStream := manager_io.NewAsyncStreamForReader(conn, manager_io.DontCloseReader())
	inStreamHandle := i.host.StreamManager().Add(inStream)
	outStream := manager_io.NewAsyncStreamForWriter(conn, manager_io.DontCloseWriter())
	outStreamHandle := i.host.StreamManager().Add(outStream)

	result := witgo.Tuple3[TCPSocket, wasip2_io.InputStream, wasip2_io.OutputStream]{
		F0: newSockHandle,
		F1: inStreamHandle,
		F2: outStreamHandle,
	}
	return witgo.Ok[witgo.Tuple3[TCPSocket, wasip2_io.InputStream, wasip2_io.OutputStream], ErrorCode](result)
}

func (i *tcpImpl) Shutdown(ctx context.Context, this TCPSocket, shutdownType ShutdownType) witgo.Result[witgo.Unit, ErrorCode] {
	sock, ok := i.host.TCPSocketManager().Get(this)
	if !ok {
		return witgo.Err[witgo.Unit, ErrorCode](ErrorCodeInvalidArgument)
	}
	conn := sock.GetConn()
	if sock.GetState() != sockets.TCPStateConnected {
		return witgo.Err[witgo.Unit, ErrorCode](ErrorCodeInvalidState)
	}
	if conn == nil {
		return witgo.Err[witgo.Unit, ErrorCode](ErrorCodeInvalidState)
	}

	var err error
	switch shutdownType {
	case ShutdownTypeReceive:
		err = conn.CloseRead()
	case ShutdownTypeSend:
		err = conn.CloseWrite()
	case ShutdownTypeBoth:
		// WASI shutdown 不应释放 tcp-socket；Conn.Close 会关掉 fd，后续 local-address/drop 未定义
		err = conn.CloseRead()
		if err2 := conn.CloseWrite(); err == nil {
			err = err2
		}
	default:
		return witgo.Err[witgo.Unit, ErrorCode](ErrorCodeInvalidArgument)
	}

	if err != nil {
		return witgo.Err[witgo.Unit, ErrorCode](mapOsError(err))
	}
	return witgo.Ok[witgo.Unit, ErrorCode](witgo.Unit{})
}

func (i *tcpImpl) LocalAddress(ctx context.Context, this TCPSocket) witgo.Result[IPSocketAddress, ErrorCode] {
	sock, ok := i.host.TCPSocketManager().Get(this)
	if !ok {
		return witgo.Err[IPSocketAddress, ErrorCode](ErrorCodeInvalidArgument)
	}

	st := sock.GetState()
	var addr net.Addr
	listen := sock.GetListener()
	conn := sock.GetConn()
	switch {
	case listen != nil && st >= sockets.TCPStateBound:
		addr = listen.Addr()
	case conn != nil && st == sockets.TCPStateConnected:
		addr = conn.LocalAddr()
	default:
		// bind 之后尚未 listen/connect 时只有 Fd，必须 getsockname
		a, err := i.localAddrFromFd(sock)
		if err != nil {
			return witgo.Err[IPSocketAddress, ErrorCode](ErrorCodeInvalidState)
		}
		addr = a
	}

	wasiAddr, err := toIPSocketAddress(addr)
	if err != nil {
		return witgo.Err[IPSocketAddress, ErrorCode](mapOsError(err))
	}
	return witgo.Ok[IPSocketAddress, ErrorCode](wasiAddr)
}

func (i *tcpImpl) RemoteAddress(ctx context.Context, this TCPSocket) witgo.Result[IPSocketAddress, ErrorCode] {
	sock, ok := i.host.TCPSocketManager().Get(this)
	if !ok {
		return witgo.Err[IPSocketAddress, ErrorCode](ErrorCodeInvalidArgument)
	}
	conn := sock.GetConn()
	if sock.GetState() != sockets.TCPStateConnected || conn == nil {
		return witgo.Err[IPSocketAddress, ErrorCode](ErrorCodeInvalidState)
	}

	wasiAddr, err := toIPSocketAddress(conn.RemoteAddr())
	if err != nil {
		return witgo.Err[IPSocketAddress, ErrorCode](mapOsError(err))
	}
	return witgo.Ok[IPSocketAddress, ErrorCode](wasiAddr)
}

func (i *tcpImpl) AddressFamily(ctx context.Context, this TCPSocket) IPAddressFamily {
	sock, ok := i.host.TCPSocketManager().Get(this)
	if !ok {
		return IPAddressFamilyIPV4 // or some other default/error indicator
	}
	family, _ := toIPAddressFamily(sock.Family)
	return family
}

func (i *tcpImpl) IsListening(ctx context.Context, this TCPSocket) bool {
	sock, ok := i.host.TCPSocketManager().Get(this)
	if !ok {
		return false
	}
	return sock.GetState() == sockets.TCPStateListening
}

// SetKeepAliveEnabled 启用或禁用 keep-alive。
func (i *tcpImpl) SetKeepAliveEnabled(ctx context.Context, this TCPSocket, value bool) witgo.Result[witgo.Unit, ErrorCode] {
	sock, ok := i.host.TCPSocketManager().Get(this)
	if !ok {
		return witgo.Err[witgo.Unit, ErrorCode](ErrorCodeInvalidArgument)
	}
	conn := sock.GetConn()
	if conn != nil {
		if err := conn.SetKeepAlive(value); err != nil {
			return witgo.Err[witgo.Unit, ErrorCode](mapOsError(err))
		}
		return witgo.Ok[witgo.Unit, ErrorCode](witgo.Unit{})
	}
	// 连接前也允许设置；走平台 setsockopt（见 setsockoptInt 的 Fd 回退）
	v := 0
	if value {
		v = 1
	}
	return i.setKeepAliveEnabledUnconnected(this, v)
}

// SetReceiveBufferSize 设置接收缓冲区大小。
func (i *tcpImpl) SetReceiveBufferSize(ctx context.Context, this TCPSocket, value uint64) witgo.Result[witgo.Unit, ErrorCode] {
	sock, ok := i.host.TCPSocketManager().Get(this)
	if !ok {
		return witgo.Err[witgo.Unit, ErrorCode](ErrorCodeInvalidArgument)
	}
	conn := sock.GetConn()
	if conn != nil {
		if err := conn.SetReadBuffer(clampToInt(value)); err != nil {
			return witgo.Err[witgo.Unit, ErrorCode](mapOsError(err))
		}
		return witgo.Ok[witgo.Unit, ErrorCode](witgo.Unit{})
	}
	return i.setReceiveBufferUnconnected(this, value)
}

// SetSendBufferSize 设置发送缓冲区大小。
func (i *tcpImpl) SetSendBufferSize(ctx context.Context, this TCPSocket, value uint64) witgo.Result[witgo.Unit, ErrorCode] {
	sock, ok := i.host.TCPSocketManager().Get(this)
	if !ok {
		return witgo.Err[witgo.Unit, ErrorCode](ErrorCodeInvalidArgument)
	}
	conn := sock.GetConn()
	if conn != nil {
		if err := conn.SetWriteBuffer(clampToInt(value)); err != nil {
			return witgo.Err[witgo.Unit, ErrorCode](mapOsError(err))
		}
		return witgo.Ok[witgo.Unit, ErrorCode](witgo.Unit{})
	}
	return i.setSendBufferUnconnected(this, value)
}

// Subscribe 创建一个 pollable 用于异步操作。
func (i *tcpImpl) Subscribe(pctx context.Context, this TCPSocket) wasip2_io.Pollable {
	sock, ok := i.host.TCPSocketManager().Get(this)
	if !ok {
		p := manager_io.NewPollable(nil)
		handle := i.host.PollManager().Add(p)
		p.SetReady()
		return handle
	}

	switch sock.GetState() {
	case sockets.TCPStateConnecting:
		if ch := sock.GetConnectDone(); ch != nil {
			return i.host.PollManager().Add(manager_io.NewPollableByChan(ch, nil))
		}
	case sockets.TCPStateListening:
		// listening 的 subscribe 必须在“有连接可 accept”时才就绪；
		// 原先一律 SetReady 会让 guest 忙等，且与非阻塞 Accept 的 WouldBlock 打架。
		return i.subscribeListen(sock) // 平台实现
	}
	p := manager_io.NewPollable(nil)
	handle := i.host.PollManager().Add(p)
	p.SetReady()
	return handle
}

func nsToSockoptSeconds(ns uint64) int {
	const nsPerSec = uint64(time.Second)
	if ns == 0 {
		return 0
	}
	sec := ns / nsPerSec
	if ns%nsPerSec != 0 {
		sec++
	}
	max := uint64(^uint(0) >> 1)
	if sec > max {
		return int(max)
	}
	if sec == 0 {
		return 1
	}
	return int(sec)
}

func sockoptSecondsToNs(sec int) uint64 {
	if sec <= 0 {
		return 0
	}
	return uint64(sec) * uint64(time.Second)
}

func clampToInt(value uint64) int {
	max := uint64(^uint(0) >> 1)
	if value > max {
		return int(max)
	}
	return int(value)
}
