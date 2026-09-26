//go:build windows

package v0_2

import (
	"context"
	"errors"
	"net"
	"os"
	"sync"
	"syscall"
	"time"
	"unsafe"

	manager_io "github.com/456vv/x/vweb_dynamic/internal/wazero/manager/io"
	"github.com/456vv/x/vweb_dynamic/internal/wazero/manager/sockets"
	wasip2_io "github.com/456vv/x/vweb_dynamic/internal/wazero/wasip2/io/v0_2"
	witgo "github.com/456vv/x/vweb_dynamic/internal/wazero/witgo"

	"golang.org/x/sys/windows"
)

func (i *tcpImpl) DropTCPSocket(_ context.Context, handle TCPSocket) {
	i.host.TCPSocketManager().Remove(handle)
}

func (i *tcpImpl) StartBind(_ context.Context, this TCPSocket, network Network, localAddress IPSocketAddress) witgo.Result[witgo.Unit, ErrorCode] {
	sock, ok := i.host.TCPSocketManager().Get(this)
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

	err = sock.DoBind(func() error {
		if !sock.HasFd() {
			// 无效句柄是 INVALID_SOCKET(-1)，不是 0
			return sockets.ErrInvalidSocketState
		}
		bindErr := syscall.Bind(syscall.Handle(sock.Fd), sockaddr)
		return bindErr
	})
	if err != nil {
		if errors.Is(err, sockets.ErrInvalidSocketState) {
			return witgo.Err[witgo.Unit, ErrorCode](ErrorCodeInvalidState)
		}
		return witgo.Err[witgo.Unit, ErrorCode](mapOsError(err))
	}
	return witgo.Ok[witgo.Unit, ErrorCode](witgo.Unit{})
}

func (i *tcpImpl) listenTCP(sock *sockets.TCPSocket) error {
	if sock.Listener != nil {
		return nil
	}
	if !sock.HasFd() {
		return syscall.EBADF
	}

	if err := windows.Listen(windows.Handle(sock.Fd), syscall.SOMAXCONN); err != nil {
		return err
	}
	var nfd windows.Handle
	p := windows.CurrentProcess()
	if err := windows.DuplicateHandle(p, windows.Handle(sock.Fd), p, &nfd, 0, false, windows.DUPLICATE_SAME_ACCESS); err != nil {
		return nil
	}
	file := os.NewFile(uintptr(nfd), "")
	if file == nil {
		// NewFile 失败必须 Closesocket(nfd)，否则泄漏 SOCKET。
		windows.Closesocket(nfd)
		return syscall.EINVAL
	}
	ln, err := net.FileListener(file)
	_ = file.Close()
	if err != nil {
		return err
	}
	tl, ok := ln.(*net.TCPListener)
	if !ok {
		_ = ln.Close()
		return syscall.EINVAL
	}
	sock.Listener = tl
	return nil
}

// soError 是 winsock2.h 的 SO_ERROR（0x1007）。
// Go 的 syscall / 部分 x/sys/windows 不导出该常量。
const soError = 0x1007

// winFdSet 对应 Winsock fd_set：{fd_count, SOCKET fd_array[]}，不是 Unix 位图。
// syscall.FdSet/Select 在 Windows 上不存在；windows.WSAPoll 在旧 x/sys 上也不存在。
const winFdSetSize = 64

type winFdSet struct {
	fdCount uint32
	fdArray [winFdSetSize]windows.Handle
}

var (
	modws2_32  = windows.NewLazySystemDLL("ws2_32.dll")
	procSelect = modws2_32.NewProc("select")
	procAccept = modws2_32.NewProc("accept") // syscall/windows.Accept 在 Windows 上都是 EWINDOWS 桩。
)

// acceptOnHandle 调用 ws2_32.accept。rsa 仅占位，对端地址由后续 FileConn 再取。
func acceptOnHandle(fd windows.Handle) (windows.Handle, error) {
	var rsa windows.RawSockaddrAny
	addrlen := int32(unsafe.Sizeof(rsa))
	r1, _, callErr := procAccept.Call(
		uintptr(fd),
		uintptr(unsafe.Pointer(&rsa)),
		uintptr(unsafe.Pointer(&addrlen)),
	)
	nfd := windows.Handle(r1)
	// Winsock：失败返回 INVALID_SOCKET（^Handle(0)），不是 Unix 的 -1 语义下的 0。
	if nfd == windows.InvalidHandle {
		if callErr != nil {
			return windows.InvalidHandle, callErr
		}
		return windows.InvalidHandle, syscall.EINVAL
	}
	return nfd, nil
}

// waitTCPConnect 等待非阻塞 connect 完成（可写或 except）。
// select 失败时用 Call 的 lastErr；SOCKET_ERROR 恒为 -1。
func waitTCPConnect(fd windows.Handle, ctx context.Context) error {
	for {
		if ctx != nil {
			if err := ctx.Err(); err != nil {
				return err
			}
		}
		var wset, eset winFdSet
		wset.fdCount = 1
		wset.fdArray[0] = fd
		eset.fdCount = 1
		eset.fdArray[0] = fd
		// timeout=NULL 时 drop 无法打断；200ms 轮询观察 ConnectContext
		tv := syscall.Timeval{Sec: 0, Usec: 200000}
		r1, _, callErr := procSelect.Call(
			0,
			0,
			uintptr(unsafe.Pointer(&wset)),
			uintptr(unsafe.Pointer(&eset)),
			uintptr(unsafe.Pointer(&tv)),
		)
		if int32(r1) == -1 {
			if callErr == windows.WSAEINTR {
				continue
			}
			if callErr != nil {
				return callErr
			}
			return syscall.EINVAL
		}
		if r1 > 0 {
			return nil
		}
	}
}

func (i *tcpImpl) connectTCP(sock *sockets.TCPSocket, remoteAddress IPSocketAddress) (*net.TCPConn, error) {
	_, hasFd := sock.SnapshotFd()
	if !hasFd {
		addr, err := fromIPSocketAddressToTCPAddr(remoteAddress)
		if err != nil {
			return nil, err
		}
		network := "tcp"
		if sock.Family == sockets.IPAddressFamilyIPV6 {
			network = "tcp6"
		}
		// drop 必须能取消拨号
		c, err := (&net.Dialer{}).DialContext(sock.ConnectContext(), network, addr.String())
		if err != nil {
			return nil, err
		}
		tc, ok := c.(*net.TCPConn)
		if !ok {
			_ = c.Close()
			return nil, syscall.EINVAL
		}
		return tc, nil
	}

	fd, err := sock.DupOwnedFd()
	if err != nil {
		return nil, err
	}
	defer windows.Closesocket(windows.Handle(fd))

	sa, err := fromIPSocketAddressToSockaddr(remoteAddress)
	if err != nil {
		return nil, err
	}
	cerr := syscall.Connect(syscall.Handle(fd), sa)
	if cerr != nil &&
		cerr != syscall.EWOULDBLOCK &&
		cerr != windows.WSAEWOULDBLOCK &&
		cerr != syscall.EINPROGRESS &&
		cerr != windows.WSAEINPROGRESS {
		return nil, cerr
	}
	if cerr != nil {
		if err := waitTCPConnect(windows.Handle(fd), sock.ConnectContext()); err != nil {
			return nil, err
		}
		soerr, gerr := windows.GetsockoptInt(windows.Handle(fd), windows.SOL_SOCKET, soError)
		if gerr != nil {
			return nil, gerr
		}
		if soerr != 0 {
			return nil, syscall.Errno(soerr)
		}
	}
	var nfd windows.Handle
	p := windows.CurrentProcess()
	if err := windows.DuplicateHandle(p, windows.Handle(fd), p, &nfd, 0, false, windows.DUPLICATE_SAME_ACCESS); err != nil {
		return nil, err
	}
	file := os.NewFile(uintptr(nfd), "")
	if file == nil {
		// NewFile 失败必须关 duplicated SOCKET
		windows.Closesocket(nfd)
		return nil, syscall.EINVAL
	}
	c, err := net.FileConn(file)
	_ = file.Close()
	if err != nil {
		return nil, err
	}
	tc, ok := c.(*net.TCPConn)
	if !ok {
		_ = c.Close()
		return nil, syscall.EINVAL
	}
	return tc, nil
}

func (i *tcpImpl) acceptTCP(sock *sockets.TCPSocket) (*net.TCPConn, error) {
	if sock == nil {
		return nil, syscall.EINVAL
	}

	ln := sock.GetListener()
	if ln != nil {
		ln.SetDeadline(time.Now().Add(time.Millisecond))
		conn, err := ln.AcceptTCP()
		ln.SetDeadline(time.Time{})
		return conn, err
	}

	_, hasFd := sock.SnapshotFd()
	if !hasFd {
		return nil, syscall.EINVAL
	}
	fd, err := sock.DupOwnedFd()
	if err != nil {
		return nil, err
	}
	defer windows.Closesocket(windows.Handle(fd))

	// 不能用 windows.Accept / syscall.Accept，两者在 Windows 上都是 EWINDOWS 桩。
	accepted, err := acceptOnHandle(windows.Handle(fd))
	if err != nil {
		return nil, err
	}

	windows.SetHandleInformation(accepted, windows.HANDLE_FLAG_INHERIT, 0)
	var nonBlockingMode uint32 = 1
	var bytesReturned uint32
	if ioctlErr := windows.WSAIoctl(
		accepted,
		FIONBIO,
		(*byte)(unsafe.Pointer(&nonBlockingMode)),
		uint32(unsafe.Sizeof(nonBlockingMode)),
		nil,
		0,
		&bytesReturned,
		nil,
		0,
	); ioctlErr != nil {
		_ = windows.Closesocket(accepted)
		return nil, ioctlErr
	}

	file := os.NewFile(uintptr(accepted), "")
	if file == nil {
		windows.Closesocket(accepted)
		return nil, syscall.EINVAL
	}
	c, err := net.FileConn(file)
	file.Close() // FileConn 再 dup；关 file 不关连接
	if err != nil {
		return nil, err
	}
	tc, ok := c.(*net.TCPConn)
	if !ok {
		c.Close()
		return nil, syscall.EINVAL
	}
	return tc, nil
}

func (i *tcpImpl) localAddrFromFd(sock *sockets.TCPSocket) (net.Addr, error) {
	var sa windows.Sockaddr
	err := sock.ControlFd(func(fd int) error {
		var gerr error
		sa, gerr = windows.Getsockname(windows.Handle(fd))
		return gerr
	})
	if err != nil {
		return nil, err
	}
	switch a := sa.(type) {
	case *windows.SockaddrInet4:
		ip := make(net.IP, net.IPv4len)
		copy(ip, a.Addr[:])
		return &net.TCPAddr{IP: ip, Port: a.Port}, nil
	case *windows.SockaddrInet6:
		ip := make(net.IP, net.IPv6len)
		copy(ip, a.Addr[:])
		return &net.TCPAddr{IP: ip, Port: a.Port, Zone: zoneFromScopeID(a.ZoneId)}, nil
	default:
		return nil, syscall.EAFNOSUPPORT
	}
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

func (i *tcpImpl) SetListenBacklogSize(ctx context.Context, this TCPSocket, value uint64) witgo.Result[witgo.Unit, ErrorCode] {
	sock, ok := i.host.TCPSocketManager().Get(this)
	if !ok {
		return witgo.Err[witgo.Unit, ErrorCode](ErrorCodeInvalidArgument)
	}
	n := clampToInt(value)
	ctrlErr := sock.ControlFd(func(fd int) error {
		return windows.Listen(windows.Handle(fd), n)
	})
	if ctrlErr == nil {
		return witgo.Ok[witgo.Unit, ErrorCode](witgo.Unit{})
	}
	if !errors.Is(ctrlErr, sockets.ErrInvalidSocketState) {
		return witgo.Err[witgo.Unit, ErrorCode](mapOsError(ctrlErr))
	}
	ln := sock.GetListener()
	if ln == nil {
		return witgo.Err[witgo.Unit, ErrorCode](ErrorCodeInvalidState)
	}
	raw, err := ln.SyscallConn()
	if err != nil {
		return witgo.Err[witgo.Unit, ErrorCode](mapOsError(err))
	}
	var listenErr error
	if ctrlErr = raw.Control(func(fd uintptr) {
		listenErr = windows.Listen(windows.Handle(fd), n)
	}); ctrlErr != nil {
		return witgo.Err[witgo.Unit, ErrorCode](mapOsError(ctrlErr))
	}
	if listenErr != nil {
		return witgo.Err[witgo.Unit, ErrorCode](mapOsError(listenErr))
	}
	return witgo.Ok[witgo.Unit, ErrorCode](witgo.Unit{})
}

func (i *tcpImpl) KeepAliveEnabled(ctx context.Context, this TCPSocket) witgo.Result[bool, ErrorCode] {
	result := getTCPSockopt[int](i, this, windows.SOL_SOCKET, windows.SO_KEEPALIVE)
	if result.Err != nil {
		return witgo.Err[bool, ErrorCode](*result.Err)
	}
	return witgo.Ok[bool, ErrorCode](*result.Ok == 1)
}

func (i *tcpImpl) setKeepAliveEnabledUnconnected(this TCPSocket, value int) witgo.Result[witgo.Unit, ErrorCode] {
	return setTCPSockopt(i, this, windows.SOL_SOCKET, windows.SO_KEEPALIVE, value)
}

func (i *tcpImpl) KeepAliveIdleTime(ctx context.Context, this TCPSocket) witgo.Result[uint64, ErrorCode] {
	// WASI 是纳秒，TCP_KEEPIDLE 是秒
	result := getTCPSockopt[int](i, this, windows.IPPROTO_TCP, windows.TCP_KEEPIDLE)
	if result.Err != nil {
		return witgo.Err[uint64, ErrorCode](*result.Err)
	}
	return witgo.Ok[uint64, ErrorCode](sockoptSecondsToNs(*result.Ok))
}

func (i *tcpImpl) SetKeepAliveIdleTime(ctx context.Context, this TCPSocket, value uint64) witgo.Result[witgo.Unit, ErrorCode] {
	return setTCPSockopt(i, this, windows.IPPROTO_TCP, windows.TCP_KEEPIDLE, nsToSockoptSeconds(value))
}

func (i *tcpImpl) KeepAliveInterval(ctx context.Context, this TCPSocket) witgo.Result[uint64, ErrorCode] {
	result := getTCPSockopt[int](i, this, windows.IPPROTO_TCP, windows.TCP_KEEPINTVL)
	if result.Err != nil {
		return witgo.Err[uint64, ErrorCode](*result.Err)
	}
	return witgo.Ok[uint64, ErrorCode](sockoptSecondsToNs(*result.Ok))
}

func (i *tcpImpl) SetKeepAliveInterval(ctx context.Context, this TCPSocket, value uint64) witgo.Result[witgo.Unit, ErrorCode] {
	return setTCPSockopt(i, this, windows.IPPROTO_TCP, windows.TCP_KEEPINTVL, nsToSockoptSeconds(value))
}

func (i *tcpImpl) KeepAliveCount(ctx context.Context, this TCPSocket) witgo.Result[uint32, ErrorCode] {
	return getTCPSockopt[uint32](i, this, windows.IPPROTO_TCP, windows.TCP_KEEPCNT)
}

func (i *tcpImpl) SetKeepAliveCount(ctx context.Context, this TCPSocket, value uint32) witgo.Result[witgo.Unit, ErrorCode] {
	return setTCPSockopt(i, this, windows.IPPROTO_TCP, windows.TCP_KEEPCNT, clampToInt(uint64(value)))
}

func (i *tcpImpl) HopLimit(ctx context.Context, this TCPSocket) witgo.Result[uint8, ErrorCode] {
	sock, ok := i.host.TCPSocketManager().Get(this)
	if !ok {
		return witgo.Err[uint8, ErrorCode](ErrorCodeInvalidArgument)
	}
	if sock.Family == sockets.IPAddressFamilyIPV6 {
		return getTCPSockopt[uint8](i, this, windows.IPPROTO_IPV6, windows.IPV6_UNICAST_HOPS)
	}
	return getTCPSockopt[uint8](i, this, windows.IPPROTO_IP, windows.IP_TTL)
}

func (i *tcpImpl) SetHopLimit(ctx context.Context, this TCPSocket, value uint8) witgo.Result[witgo.Unit, ErrorCode] {
	sock, ok := i.host.TCPSocketManager().Get(this)
	if !ok {
		return witgo.Err[witgo.Unit, ErrorCode](ErrorCodeInvalidArgument)
	}
	if sock.Family == sockets.IPAddressFamilyIPV6 {
		return setTCPSockopt(i, this, windows.IPPROTO_IPV6, windows.IPV6_UNICAST_HOPS, int(value))
	}
	return setTCPSockopt(i, this, windows.IPPROTO_IP, windows.IP_TTL, int(value))
}

func (i *tcpImpl) ReceiveBufferSize(ctx context.Context, this TCPSocket) witgo.Result[uint64, ErrorCode] {
	return getTCPSockopt[uint64](i, this, windows.SOL_SOCKET, windows.SO_RCVBUF)
}

func (i *tcpImpl) SendBufferSize(ctx context.Context, this TCPSocket) witgo.Result[uint64, ErrorCode] {
	return getTCPSockopt[uint64](i, this, windows.SOL_SOCKET, windows.SO_SNDBUF)
}

func (i *tcpImpl) setReceiveBufferUnconnected(this TCPSocket, value uint64) witgo.Result[witgo.Unit, ErrorCode] {
	return setTCPSockopt(i, this, windows.SOL_SOCKET, windows.SO_RCVBUF, clampToInt(value))
}

func (i *tcpImpl) setSendBufferUnconnected(this TCPSocket, value uint64) witgo.Result[witgo.Unit, ErrorCode] {
	return setTCPSockopt(i, this, windows.SOL_SOCKET, windows.SO_SNDBUF, clampToInt(value))
}

func getTCPSockopt[T ~int | ~uint64 | ~uint32 | ~uint8](i *tcpImpl, this TCPSocket, level, opt int) witgo.Result[T, ErrorCode] {
	sock, ok := i.host.TCPSocketManager().Get(this)
	if !ok || sock == nil {
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

func setTCPSockopt(i *tcpImpl, this TCPSocket, level, opt, value int) witgo.Result[witgo.Unit, ErrorCode] {
	sock, ok := i.host.TCPSocketManager().Get(this)
	if !ok || sock == nil {
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

// waitTCPAccept 等待非阻塞 listen 套接字可读（有连接或错误）。
// windows.WSAPoll 在旧 x/sys 上不存在；syscall.Select/FdSet 在 Windows 上不是位图。
// timeout 200ms，以便 pollable drop 时退出，避免 goroutine 泄漏。
func waitTCPAccept(fd windows.Handle, done <-chan struct{}) {
	for {
		select {
		case <-done:
			return
		default:
		}
		var rset, eset winFdSet
		rset.fdCount = 1
		rset.fdArray[0] = fd
		eset.fdCount = 1
		eset.fdArray[0] = fd
		tv := syscall.Timeval{Sec: 0, Usec: 200000}
		r1, _, callErr := procSelect.Call(
			0, // Windows 忽略 nfds
			uintptr(unsafe.Pointer(&rset)),
			0,
			uintptr(unsafe.Pointer(&eset)),
			uintptr(unsafe.Pointer(&tv)),
		)
		if int32(r1) == -1 {
			if callErr == windows.WSAEINTR {
				continue
			}
			return
		}
		if r1 > 0 {
			return
		}
	}
}

func (i *tcpImpl) subscribeListen(sock *sockets.TCPSocket) wasip2_io.Pollable {
	fdInt, err := sock.DupOwnedFd()
	if err != nil {
		p := manager_io.NewPollable(nil)
		handle := i.host.PollManager().Add(p)
		p.SetReady()
		return handle
	}
	fd := windows.Handle(fdInt)

	done := make(chan struct{})
	var once sync.Once
	p := manager_io.NewPollable(func() {
		once.Do(func() { close(done) })
	})
	handle := i.host.PollManager().Add(p)

	go func() {
		defer windows.Closesocket(fd)
		defer p.SetReady()
		waitTCPAccept(fd, done)
	}()
	return handle
}
