//go:build !unix && !windows

package v0_2

import (
	"context"

	"github.com/456vv/x/vweb_dynamic/internal/wazero/manager/sockets"
	witgo "github.com/456vv/x/vweb_dynamic/internal/wazero/witgo"
)

func (i *tcpCreateSocketImpl) CreateTCPSocket(_ context.Context, addressFamily IPAddressFamily) witgo.Result[TCPSocket, ErrorCode] {
	family, err := fromIPAddressFamily(addressFamily)
	if err != nil {
		return witgo.Err[TCPSocket, ErrorCode](ErrorCodeNotSupported)
	}

	tcpSocket := &sockets.TCPSocket{
		Family: family,
		State:  sockets.TCPStateUnbound,
	}
	tcpSocket.SetState(sockets.TCPStateUnbound)
	handle := i.host.TCPSocketManager().Add(tcpSocket)
	return witgo.Ok[TCPSocket, ErrorCode](handle)
}
