//go:build windows

package v0_2

import (
	"context"
	"unsafe"

	"github.com/456vv/x/vweb_dynamic/internal/wazero/manager/sockets"
	witgo "github.com/456vv/x/vweb_dynamic/internal/wazero/witgo"

	"golang.org/x/sys/windows"
)

func (i *udpCreateSocketImpl) CreateUDPSocket(_ context.Context, addressFamily IPAddressFamily) witgo.Result[UDPSocket, ErrorCode] {
	family, err := fromIPAddressFamily(addressFamily)
	if err != nil {
		return witgo.Err[UDPSocket, ErrorCode](ErrorCodeNotSupported)
	}

	domain := windows.AF_INET
	if family == sockets.IPAddressFamilyIPV6 {
		domain = windows.AF_INET6
	}

	handle, err := windows.Socket(domain, windows.SOCK_DGRAM, windows.IPPROTO_UDP)
	if err != nil {
		return witgo.Err[UDPSocket, ErrorCode](mapOsError(err))
	}

	// Windows SOCKET 默认可继承，子进程会泄漏句柄。
	windows.SetHandleInformation(handle, windows.HANDLE_FLAG_INHERIT, 0)

	// 使用 WSAIoctl 将套接字设置为非阻塞模式
	var nonBlockingMode uint32 = 1
	var bytesReturned uint32
	err = windows.WSAIoctl(
		handle,
		FIONBIO,
		(*byte)(unsafe.Pointer(&nonBlockingMode)),
		uint32(unsafe.Sizeof(nonBlockingMode)),
		nil,
		0,
		&bytesReturned,
		nil,
		0,
	)
	if err != nil {
		windows.Closesocket(handle)
		return witgo.Err[UDPSocket, ErrorCode](mapOsError(err))
	}

	if family == sockets.IPAddressFamilyIPV6 {
		val := 1
		err = windows.SetsockoptInt(handle, windows.IPPROTO_IPV6, windows.IPV6_V6ONLY, val)
		if err != nil {
			windows.Closesocket(handle)
			return witgo.Err[UDPSocket, ErrorCode](mapOsError(err))
		}
	}

	udpSocket := &sockets.UDPSocket{
		Fd:     int(handle),
		Family: family,
	}

	udpSocket.AttachFd(int(handle))
	h := i.host.UDPSocketManager().Add(udpSocket)
	return witgo.Ok[UDPSocket, ErrorCode](h)
}

func (i *udpImpl) DropUDPSocket(_ context.Context, handle UDPSocket) {
	i.host.UDPSocketManager().Remove(handle)
}
