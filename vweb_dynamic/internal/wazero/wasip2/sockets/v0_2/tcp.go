package v0_2

import (
	"context"
	"errors"
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
		conn, dialErr := i.connectTCP(sock, remoteAddress)
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

	if !sock.ClaimConnectSuccess(result.Conn) {
		// drop 抢先 Closed 时 Claim 失败，result.Conn 可能尚未被析构拷走
		if result.Conn != nil && sock.GetState() != sockets.TCPStateConnected {
			_ = result.Conn.Close()
		}
		return witgo.Err[witgo.Tuple[wasip2_io.InputStream, wasip2_io.OutputStream], ErrorCode](ErrorCodeNotInProgress)
	}

	inStream := manager_io.NewAsyncStreamForReader(sock.Conn, manager_io.DontCloseReader())
	inStreamHandle := i.host.StreamManager().Add(inStream)
	outStream := manager_io.NewAsyncStreamForWriter(sock.Conn, manager_io.DontCloseWriter())
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
	if sock.GetState() != sockets.TCPStateListening || sock.Listener == nil {
		return witgo.Err[witgo.Tuple3[TCPSocket, wasip2_io.InputStream, wasip2_io.OutputStream], ErrorCode](ErrorCodeInvalidState)
	}

	// 走平台 acceptTCP（unix 可从原 fd Accept；deadline=Now 在队列非空时仍会立刻超时）
	conn, err := i.acceptTCP(sock)
	if err != nil {
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
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
	if sock.GetState() != sockets.TCPStateConnected {
		return witgo.Err[witgo.Unit, ErrorCode](ErrorCodeInvalidState)
	}
	if sock.Conn == nil {
		return witgo.Err[witgo.Unit, ErrorCode](ErrorCodeInvalidState)
	}

	var err error
	switch shutdownType {
	case ShutdownTypeReceive:
		err = sock.Conn.CloseRead()
	case ShutdownTypeSend:
		err = sock.Conn.CloseWrite()
	case ShutdownTypeBoth:
		err = sock.Conn.Close() // Close会同时关闭读和写
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
	switch {
	case sock.Listener != nil && st >= sockets.TCPStateBound:
		addr = sock.Listener.Addr()
	case sock.Conn != nil && st == sockets.TCPStateConnected:
		addr = sock.Conn.LocalAddr()
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
	if sock.GetState() != sockets.TCPStateConnected || sock.Conn == nil {
		return witgo.Err[IPSocketAddress, ErrorCode](ErrorCodeInvalidState)
	}

	wasiAddr, err := toIPSocketAddress(sock.Conn.RemoteAddr())
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
	if sock.Conn != nil {
		if err := sock.Conn.SetKeepAlive(value); err != nil {
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
	if sock.Conn != nil {
		if err := sock.Conn.SetReadBuffer(clampToInt(value)); err != nil {
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
	if sock.Conn != nil {
		if err := sock.Conn.SetWriteBuffer(clampToInt(value)); err != nil {
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
		if sock.ConnectDone != nil {
			return i.host.PollManager().Add(manager_io.NewPollableByChan(sock.ConnectDone, nil))
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
