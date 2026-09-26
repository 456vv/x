package v0_2

import (
	"context"
	"errors"
	"io"
	"net"
	"time"

	manager_io "github.com/456vv/x/vweb_dynamic/internal/wazero/manager/io"
	manager_sockets "github.com/456vv/x/vweb_dynamic/internal/wazero/manager/sockets"

	"github.com/456vv/x/vweb_dynamic/internal/wazero/wasip2"
	wasip2_io "github.com/456vv/x/vweb_dynamic/internal/wazero/wasip2/io/v0_2"
	witgo "github.com/456vv/x/vweb_dynamic/internal/wazero/witgo"
)

// rebindUnconnectedUDP 关闭已 connect 的 UDP 再按原本地地址 Listen，恢复「未连接」语义。
// 调用方必须已持有 UDPSocket.mu（ReplaceDatagramStreams）。
func rebindUnconnectedUDP(sock *manager_sockets.UDPSocket) error {
	if sock == nil || sock.Conn == nil {
		return errors.New("udp socket is not bound")
	}
	network := "udp"
	if sock.Family == manager_sockets.IPAddressFamilyIPV6 {
		network = "udp6"
	}
	var local *net.UDPAddr
	if la, ok := sock.Conn.LocalAddr().(*net.UDPAddr); ok && la != nil {
		ip := append(net.IP(nil), la.IP...)
		local = &net.UDPAddr{IP: ip, Port: la.Port, Zone: la.Zone}
	}
	_ = sock.Conn.Close()
	conn, err := net.ListenUDP(network, local)
	if err != nil {
		sock.Conn = nil
		return err
	}
	sock.Conn = conn
	return nil
}

type udpImpl struct {
	host *wasip2.Host
}

func newUDPImpl(h *wasip2.Host) *udpImpl {
	return &udpImpl{host: h}
}

func (i *udpImpl) checkNetwork(network Network) ErrorCode {
	if i.host == nil || i.host.NetworkManager() == nil {
		return ErrorCodeInvalidArgument
	}
	if _, ok := i.host.NetworkManager().Get(network); !ok {
		return ErrorCodeInvalidArgument
	}
	return 0
}

func (i *udpImpl) FinishBind(_ context.Context, this UDPSocket) witgo.Result[witgo.Unit, ErrorCode] {
	sock, ok := i.host.UDPSocketManager().Get(this)
	if !ok || sock == nil {
		return witgo.Err[witgo.Unit, ErrorCode](ErrorCodeInvalidArgument)
	}
	// start-bind 是同步的，成功后才有 Conn。
	// 不能用 HasFd：unix/windows 创建套接字时已有 fd，未 bind 也会被当成成功。
	if sock.GetConn() == nil {
		return witgo.Err[witgo.Unit, ErrorCode](ErrorCodeInvalidState)
	}
	return witgo.Ok[witgo.Unit, ErrorCode](witgo.Unit{})
}

func (i *udpImpl) Stream(_ context.Context, this UDPSocket, remoteAddress witgo.Option[IPSocketAddress]) witgo.Result[witgo.Tuple[IncomingDatagramStream, OutgoingDatagramStream], ErrorCode] {
	sock, ok := i.host.UDPSocketManager().Get(this)
	if !ok || sock == nil || sock.GetConn() == nil {
		return witgo.Err[witgo.Tuple[IncomingDatagramStream, OutgoingDatagramStream], ErrorCode](ErrorCodeInvalidState)
	}
	if err := sock.ReplaceDatagramStreams(func() error {
		return i.connectUDP(sock, remoteAddress)
	}); err != nil {
		if errors.Is(err, manager_sockets.ErrInvalidSocketState) {
			return witgo.Err[witgo.Tuple[IncomingDatagramStream, OutgoingDatagramStream], ErrorCode](ErrorCodeInvalidState)
		}
		return witgo.Err[witgo.Tuple[IncomingDatagramStream, OutgoingDatagramStream], ErrorCode](mapOsError(err))
	}
	return witgo.Ok[witgo.Tuple[IncomingDatagramStream, OutgoingDatagramStream], ErrorCode](
		witgo.Tuple[IncomingDatagramStream, OutgoingDatagramStream]{F0: this, F1: this},
	)
}

func (i *udpImpl) LocalAddress(ctx context.Context, this UDPSocket) witgo.Result[IPSocketAddress, ErrorCode] {
	sock, ok := i.host.UDPSocketManager().Get(this)
	if !ok || sock == nil {
		return witgo.Err[IPSocketAddress, ErrorCode](ErrorCodeInvalidState)
	}
	conn := sock.GetConn()
	if conn == nil {
		return witgo.Err[IPSocketAddress, ErrorCode](ErrorCodeInvalidState)
	}
	addr, err := toIPSocketAddress(conn.LocalAddr())
	if err != nil {
		return witgo.Err[IPSocketAddress, ErrorCode](mapOsError(err))
	}
	return witgo.Ok[IPSocketAddress, ErrorCode](addr)
}

func (i *udpImpl) RemoteAddress(ctx context.Context, this UDPSocket) witgo.Result[IPSocketAddress, ErrorCode] {
	sock, ok := i.host.UDPSocketManager().Get(this)
	if !ok || sock == nil {
		return witgo.Err[IPSocketAddress, ErrorCode](ErrorCodeInvalidState)
	}
	conn := sock.GetConn()
	if conn == nil {
		return witgo.Err[IPSocketAddress, ErrorCode](ErrorCodeInvalidState)
	}
	ra := conn.RemoteAddr()
	// 未连接的 UDP RemoteAddr() 为 nil，不能当 TCPAddr 去转，应 InvalidState。
	if ra == nil {
		return witgo.Err[IPSocketAddress, ErrorCode](ErrorCodeInvalidState)
	}
	addr, err := toIPSocketAddress(ra)
	if err != nil {
		return witgo.Err[IPSocketAddress, ErrorCode](mapOsError(err))
	}
	return witgo.Ok[IPSocketAddress, ErrorCode](addr)
}

func (i *udpImpl) AddressFamily(ctx context.Context, this UDPSocket) IPAddressFamily {
	sock, ok := i.host.UDPSocketManager().Get(this)
	if !ok || sock == nil {
		return IPAddressFamilyIPV4
	}
	family, _ := toIPAddressFamily(sock.Family)
	return family
}

// SetReceiveBufferSize 设置接收缓冲区大小。
func (i *udpImpl) SetReceiveBufferSize(ctx context.Context, this UDPSocket, value uint64) witgo.Result[witgo.Unit, ErrorCode] {
	sock, ok := i.host.UDPSocketManager().Get(this)
	if !ok || sock == nil {
		return witgo.Err[witgo.Unit, ErrorCode](ErrorCodeInvalidArgument)
	}
	conn := sock.GetConn()
	if conn != nil {
		if err := conn.SetReadBuffer(clampToInt(value)); err != nil {
			return witgo.Err[witgo.Unit, ErrorCode](mapOsError(err))
		}
		return witgo.Ok[witgo.Unit, ErrorCode](witgo.Unit{})
	}
	// WASI 允许 bind 前设缓冲区；原先 Conn==nil 直接 InvalidState
	return i.setReceiveBufferUnconnected(this, value)
}

// SetSendBufferSize 设置发送缓冲区大小。
func (i *udpImpl) SetSendBufferSize(ctx context.Context, this UDPSocket, value uint64) witgo.Result[witgo.Unit, ErrorCode] {
	sock, ok := i.host.UDPSocketManager().Get(this)
	if !ok || sock == nil {
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

func (i *udpImpl) Subscribe(ctx context.Context, this UDPSocket) wasip2_io.Pollable {
	return i.host.PollManager().Add(manager_io.NewReadyPollable())
}

func (i *udpImpl) DropIncomingDatagramStream(_ context.Context, handle IncomingDatagramStream) {
	sock, ok := i.host.UDPSocketManager().Get(handle)
	if !ok || sock == nil {
		return
	}
	if r := sock.TakeReader(); r != nil {
		r.Close()
		r.WaitExit()
		conn := sock.GetConn()
		if conn != nil {
			_ = conn.SetReadDeadline(time.Time{})
		}
	}
}

func (i *udpImpl) DropOutgoingDatagramStream(_ context.Context, handle OutgoingDatagramStream) {
	sock, ok := i.host.UDPSocketManager().Get(handle)
	if !ok || sock == nil {
		return
	}
	if w := sock.TakeWriter(); w != nil {
		w.Close()
		w.WaitExit()
		conn := sock.GetConn()
		if conn != nil {
			_ = conn.SetWriteDeadline(time.Time{})
		}
	}
}

func (i *udpImpl) Receive(_ context.Context, this IncomingDatagramStream, maxResults uint64) witgo.Result[[]IncomingDatagram, ErrorCode] {
	sock, ok := i.host.UDPSocketManager().Get(this)
	if !ok || sock == nil {
		return witgo.Err[[]IncomingDatagram, ErrorCode](ErrorCodeInvalidArgument)
	}
	reader := sock.GetReader()
	if reader == nil {
		return witgo.Err[[]IncomingDatagram, ErrorCode](ErrorCodeInvalidArgument)
	}
	if maxResults == 0 {
		return witgo.Ok[[]IncomingDatagram, ErrorCode]([]IncomingDatagram{})
	}
	datagrams, err := reader.Receive(maxResults)
	if err != nil {
		if errors.Is(err, net.ErrClosed) || errors.Is(err, io.EOF) {
			return witgo.Err[[]IncomingDatagram, ErrorCode](ErrorCodeInvalidState)
		}
		return witgo.Err[[]IncomingDatagram, ErrorCode](mapOsError(err))
	}
	return witgo.Ok[[]IncomingDatagram, ErrorCode](datagrams)
}

func (i *udpImpl) Send(_ context.Context, this OutgoingDatagramStream, datagrams []OutgoingDatagram) witgo.Result[uint64, ErrorCode] {
	sock, ok := i.host.UDPSocketManager().Get(this)
	if !ok || sock == nil {
		return witgo.Err[uint64, ErrorCode](ErrorCodeInvalidArgument)
	}
	writer := sock.GetWriter()
	if writer == nil {
		return witgo.Err[uint64, ErrorCode](ErrorCodeInvalidArgument)
	}
	sentCount, err := writer.Send(datagrams)
	if err != nil {
		// writer 已返回 net.ErrClosed。这里若不识别，关闭后的 send 仍是 unknown。
		// GetWriter 对 nil socket 已返回 nil，这里不是 panic 修复。
		if errors.Is(err, net.ErrClosed) || errors.Is(err, io.EOF) {
			return witgo.Err[uint64, ErrorCode](ErrorCodeInvalidState)
		}
		return witgo.Err[uint64, ErrorCode](mapOsError(err))
	}
	return witgo.Ok[uint64, ErrorCode](sentCount)
}

func (i *udpImpl) CheckSend(_ context.Context, this OutgoingDatagramStream) witgo.Result[uint64, ErrorCode] {
	sock, ok := i.host.UDPSocketManager().Get(this)
	// 流已 drop 时应 InvalidState，而不是 Ok(0) 被 guest 当成“暂无空间”空转。
	if !ok || sock == nil {
		return witgo.Err[uint64, ErrorCode](ErrorCodeInvalidState)
	}
	writer := sock.GetWriter()
	if writer == nil {
		return witgo.Err[uint64, ErrorCode](ErrorCodeInvalidState)
	}
	// writer 已关闭时 AvailableSpace()==0 且 Subscribe 仍就绪，guest 会空转。
	if err := writer.ClosedOrErr(); err != nil {
		if errors.Is(err, net.ErrClosed) || errors.Is(err, io.EOF) {
			return witgo.Err[uint64, ErrorCode](ErrorCodeInvalidState)
		}
		return witgo.Err[uint64, ErrorCode](mapOsError(err))
	}
	return witgo.Ok[uint64, ErrorCode](writer.AvailableSpace())
}

func (i *udpImpl) SubscribeIncoming(ctx context.Context, this IncomingDatagramStream) wasip2_io.Pollable {
	sock, ok := i.host.UDPSocketManager().Get(this)
	if !ok || sock == nil {
		return i.host.PollManager().Add(manager_io.NewReadyPollable())
	}
	reader := sock.GetReader()
	if reader == nil {
		return i.host.PollManager().Add(manager_io.NewReadyPollable())
	}
	return i.host.PollManager().Add(reader.Subscribe())
}

func (i *udpImpl) SubscribeOutgoing(ctx context.Context, this OutgoingDatagramStream) wasip2_io.Pollable {
	sock, ok := i.host.UDPSocketManager().Get(this)
	if !ok || sock == nil {
		return i.host.PollManager().Add(manager_io.NewReadyPollable())
	}
	writer := sock.GetWriter()
	if writer == nil {
		return i.host.PollManager().Add(manager_io.NewReadyPollable())
	}
	return i.host.PollManager().Add(writer.Subscribe())
}
