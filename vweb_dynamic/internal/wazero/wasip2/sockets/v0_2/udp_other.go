//go:build !unix && !windows

package v0_2

import (
	"context"
	"errors"
	"net"

	"github.com/456vv/x/vweb_dynamic/internal/wazero/manager/sockets"
	witgo "github.com/456vv/x/vweb_dynamic/internal/wazero/witgo"
)

func (i *udpImpl) StartBind(_ context.Context, this UDPSocket, network Network, localAddress IPSocketAddress) witgo.Result[witgo.Unit, ErrorCode] {
	sock, ok := i.host.UDPSocketManager().Get(this)
	if !ok {
		return witgo.Err[witgo.Unit, ErrorCode](ErrorCodeInvalidArgument)
	}
	if code := i.checkNetwork(network); code != 0 {
		return witgo.Err[witgo.Unit, ErrorCode](code)
	}
	addr, err := fromIPSocketAddressToUDPAddr(localAddress)
	if err != nil {
		return witgo.Err[witgo.Unit, ErrorCode](ErrorCodeInvalidArgument)
	}

	pnetwork := "udp"
	if sock.Family == sockets.IPAddressFamilyIPV6 {
		pnetwork = "udp6"
	}

	// 与 Stream/Drop/析构共用 UDPSocket.mu，避免并发 start-bind 双 ListenUDP。
	err = sock.DoBind(func() error {
		conn, listenErr := net.ListenUDP(pnetwork, addr)
		if listenErr != nil {
			return listenErr
		}
		sock.Conn = conn
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

func (i *udpImpl) FinishBind(_ context.Context, this UDPSocket) witgo.Result[witgo.Unit, ErrorCode] {
	// 我们的 start-bind 是同步的，所以这里直接成功返回
	return witgo.Ok[witgo.Unit, ErrorCode](witgo.Unit{})
}

func (i *udpImpl) UnicastHopLimit(ctx context.Context, this UDPSocket) witgo.Result[uint8, ErrorCode] {
	return witgo.Err[uint8, ErrorCode](ErrorCodeNotSupported)
}

func (i *udpImpl) SetUnicastHopLimit(ctx context.Context, this UDPSocket, value uint8) witgo.Result[witgo.Unit, ErrorCode] {
	return witgo.Err[witgo.Unit, ErrorCode](ErrorCodeNotSupported)
}

func (i *udpImpl) ReceiveBufferSize(ctx context.Context, this UDPSocket) witgo.Result[uint64, ErrorCode] {
	return witgo.Err[uint64, ErrorCode](ErrorCodeNotSupported)
}

func (i *udpImpl) SendBufferSize(ctx context.Context, this UDPSocket) witgo.Result[uint64, ErrorCode] {
	return witgo.Err[uint64, ErrorCode](ErrorCodeNotSupported)
}

func (i *udpImpl) connectUDP(sock *sockets.UDPSocket, remoteAddress witgo.Option[IPSocketAddress]) error {
	if sock == nil || sock.Conn == nil {
		return errors.New("udp socket is not bound")
	}
	if !remoteAddress.IsSome() || remoteAddress.Some == nil {
		return nil
	}
	raddr, err := fromIPSocketAddressToUDPAddr(*remoteAddress.Some)
	if err != nil {
		return err
	}
	laddr, _ := sock.Conn.LocalAddr().(*net.UDPAddr)
	network := "udp"
	if sock.Family == sockets.IPAddressFamilyIPV6 {
		network = "udp6"
	}
	// 同一本地地址必须先关再 Dial；失败时 ListenUDP 恢复，避免 Conn 已关仍挂在 socket 上。
	_ = sock.Conn.Close()
	conn, err := net.DialUDP(network, laddr, raddr)
	if err != nil {
		if restored, rerr := net.ListenUDP(network, laddr); rerr == nil {
			sock.Conn = restored
		} else {
			sock.Conn = nil
		}
		return err
	}
	sock.Conn = conn
	return nil
}
