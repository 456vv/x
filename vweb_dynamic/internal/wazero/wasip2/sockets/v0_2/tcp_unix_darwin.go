//go:build unix && darwin

package v0_2

import (
	"context"

	witgo "github.com/456vv/x/vweb_dynamic/internal/wazero/witgo"

	"golang.org/x/sys/unix"
)

func (i *tcpImpl) KeepAliveIdleTime(ctx context.Context, this TCPSocket) witgo.Result[uint64, ErrorCode] {
	result := getsockoptInt[int](i, this, unix.IPPROTO_TCP, unix.TCP_KEEPALIVE)
	if result.Err != nil {
		return witgo.Err[uint64, ErrorCode](*result.Err)
	}
	return witgo.Ok[uint64, ErrorCode](sockoptSecondsToNs(*result.Ok))
}

func (i *tcpImpl) SetKeepAliveIdleTime(ctx context.Context, this TCPSocket, value uint64) witgo.Result[witgo.Unit, ErrorCode] {
	return setsockoptInt(i, this, unix.IPPROTO_TCP, unix.TCP_KEEPALIVE, nsToSockoptSeconds(value))
}
