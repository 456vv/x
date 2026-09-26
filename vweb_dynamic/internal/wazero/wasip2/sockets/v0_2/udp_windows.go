//go:build windows

package v0_2

import (
	"context"
	"errors"
	"net"
	"os"
	"syscall"

	"github.com/456vv/x/vweb_dynamic/internal/wazero/manager/sockets"
	witgo "github.com/456vv/x/vweb_dynamic/internal/wazero/witgo"

	"golang.org/x/sys/windows"
)

func (i *udpImpl) StartBind(_ context.Context, this UDPSocket, network Network, localAddress IPSocketAddress) witgo.Result[witgo.Unit, ErrorCode] {
	sock, ok := i.host.UDPSocketManager().Get(this)
	if !ok || sock == nil {
		return witgo.Err[witgo.Unit, ErrorCode](ErrorCodeInvalidArgument)
	}
	if code := i.checkNetwork(network); code != 0 {
		return witgo.Err[witgo.Unit, ErrorCode](code)
	}

	sockaddr, err := fromIPSocketAddressToSockaddr(localAddress)
	if err != nil {
		return witgo.Err[witgo.Unit, ErrorCode](ErrorCodeInvalidArgument)
	}

	// 无效 SOCKET 是 INVALID_SOCKET(-1) 不是 0；Conn/fd 交接必须与析构同一把锁。
	err = sock.DoBind(func() error {
		if !sock.HasFd() {
			return sockets.ErrInvalidSocketState
		}
		if bindErr := syscall.Bind(syscall.Handle(sock.Fd), sockaddr); bindErr != nil {
			return bindErr
		}

		file := os.NewFile(uintptr(sock.Fd), "")
		if file == nil {
			return sockets.ErrInvalidSocketState
		}
		conn, connErr := net.FileConn(file)
		// NewFile 接管 SOCKET，FileConn 再 dup；必须 file.Close()，成功后不要再 Closesocket(sock.Fd)。
		_ = file.Close()
		if connErr != nil {
			sock.DetachFd()
			return connErr
		}
		udpConn, ok := conn.(*net.UDPConn)
		if !ok {
			_ = conn.Close()
			sock.DetachFd()
			return sockets.ErrInvalidSocketState
		}
		sock.Conn = udpConn
		sock.DetachFd()
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

func (i *udpImpl) connectUDP(sock *sockets.UDPSocket, remoteAddress witgo.Option[IPSocketAddress]) error {
	// udp.go Stream() 调用 connectUDP，本文件原先未定义导致 Windows 无法编译
	if sock == nil || sock.Conn == nil {
		return syscall.EINVAL
	}
	if !remoteAddress.IsSome() || remoteAddress.Some == nil {
		// WASI stream(none) 必须解除上次 connect 的默认远端；
		// unix/windows 上对「空地址」syscall.Connect 不可移植，统一按原本地地址 Listen 恢复未连接。
		return rebindUnconnectedUDP(sock)
	}

	sa, err := fromIPSocketAddressToSockaddr(*remoteAddress.Some)
	if err != nil {
		return err
	}
	raw, err := sock.Conn.SyscallConn()
	if err != nil {
		return err
	}
	var cerr error
	if err := raw.Control(func(fd uintptr) {
		cerr = syscall.Connect(syscall.Handle(fd), sa)
		if cerr == syscall.EWOULDBLOCK || cerr == windows.WSAEWOULDBLOCK ||
			cerr == syscall.EINPROGRESS || cerr == windows.WSAEINPROGRESS {
			cerr = nil
		}
	}); err != nil {
		return err
	}
	return cerr
}

func (i *udpImpl) UnicastHopLimit(ctx context.Context, this UDPSocket) witgo.Result[uint8, ErrorCode] {
	sock, ok := i.host.UDPSocketManager().Get(this)
	if !ok {
		return witgo.Err[uint8, ErrorCode](ErrorCodeInvalidArgument)
	}
	if sock.Family == sockets.IPAddressFamilyIPV6 {
		return getUDPSockoptInt[uint8](i, this, windows.IPPROTO_IPV6, windows.IPV6_UNICAST_HOPS)
	}
	return getUDPSockoptInt[uint8](i, this, windows.IPPROTO_IP, windows.IP_TTL)
}

func (i *udpImpl) SetUnicastHopLimit(ctx context.Context, this UDPSocket, value uint8) witgo.Result[witgo.Unit, ErrorCode] {
	sock, ok := i.host.UDPSocketManager().Get(this)
	if !ok {
		return witgo.Err[witgo.Unit, ErrorCode](ErrorCodeInvalidArgument)
	}
	if sock.Family == sockets.IPAddressFamilyIPV6 {
		return setUDPSockoptInt(i, this, windows.IPPROTO_IPV6, windows.IPV6_UNICAST_HOPS, int(value))
	}
	return setUDPSockoptInt(i, this, windows.IPPROTO_IP, windows.IP_TTL, int(value))
}

func (i *udpImpl) ReceiveBufferSize(ctx context.Context, this UDPSocket) witgo.Result[uint64, ErrorCode] {
	return getUDPSockoptInt[uint64](i, this, windows.SOL_SOCKET, windows.SO_RCVBUF)
}

func (i *udpImpl) SendBufferSize(ctx context.Context, this UDPSocket) witgo.Result[uint64, ErrorCode] {
	return getUDPSockoptInt[uint64](i, this, windows.SOL_SOCKET, windows.SO_SNDBUF)
}

func getUDPSockoptInt[T ~int | ~uint64 | ~uint32 | ~uint8](i *udpImpl, this UDPSocket, level, opt int) witgo.Result[T, ErrorCode] {
	sock, ok := i.host.UDPSocketManager().Get(this)
	if !ok {
		return witgo.Err[T, ErrorCode](ErrorCodeInvalidArgument)
	}

	var val int
	var getErr error

	conn, _, hasFd := sock.SnapshotConnOrFd()
	switch {
	case conn != nil:
		rawConn, err := conn.SyscallConn()
		if err != nil {
			return witgo.Err[T, ErrorCode](mapOsError(err))
		}
		err = rawConn.Control(func(fd uintptr) {
			val, getErr = windows.GetsockoptInt(windows.Handle(fd), level, opt)
		})
		if err != nil {
			return witgo.Err[T, ErrorCode](mapOsError(err))
		}
	case hasFd:
		getErr = sock.ControlFd(func(fd int) error {
			var err error
			val, err = windows.GetsockoptInt(windows.Handle(fd), level, opt)
			return err
		})
	default:
		return witgo.Err[T, ErrorCode](ErrorCodeInvalidArgument)
	}
	if getErr != nil {
		return witgo.Err[T, ErrorCode](mapOsError(getErr))
	}

	return witgo.Ok[T, ErrorCode](T(val))
}

func setUDPSockoptInt(i *udpImpl, this UDPSocket, level, opt, value int) witgo.Result[witgo.Unit, ErrorCode] {
	sock, ok := i.host.UDPSocketManager().Get(this)
	if !ok {
		return witgo.Err[witgo.Unit, ErrorCode](ErrorCodeInvalidArgument)
	}

	var setErr error
	conn, _, hasFd := sock.SnapshotConnOrFd()
	switch {
	case conn != nil:
		rawConn, err := conn.SyscallConn()
		if err != nil {
			return witgo.Err[witgo.Unit, ErrorCode](mapOsError(err))
		}
		err = rawConn.Control(func(fd uintptr) {
			setErr = windows.SetsockoptInt(windows.Handle(fd), level, opt, value)
		})
		if err != nil {
			return witgo.Err[witgo.Unit, ErrorCode](mapOsError(err))
		}
	case hasFd:
		setErr = sock.ControlFd(func(fd int) error {
			return windows.SetsockoptInt(windows.Handle(fd), level, opt, value)
		})
	default:
		return witgo.Err[witgo.Unit, ErrorCode](ErrorCodeInvalidArgument)
	}
	if setErr != nil {
		return witgo.Err[witgo.Unit, ErrorCode](mapOsError(setErr))
	}

	return witgo.Ok[witgo.Unit, ErrorCode](witgo.Unit{})
}

func (i *udpImpl) setReceiveBufferUnconnected(this UDPSocket, value uint64) witgo.Result[witgo.Unit, ErrorCode] {
	return setUDPSockoptInt(i, this, windows.SOL_SOCKET, windows.SO_RCVBUF, clampToInt(value))
}

func (i *udpImpl) setSendBufferUnconnected(this UDPSocket, value uint64) witgo.Result[witgo.Unit, ErrorCode] {
	return setUDPSockoptInt(i, this, windows.SOL_SOCKET, windows.SO_SNDBUF, clampToInt(value))
}
