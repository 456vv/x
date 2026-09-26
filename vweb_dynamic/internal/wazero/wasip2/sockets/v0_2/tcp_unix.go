//go:build unix

package v0_2

import (
	"context"
	"errors"
	"net"
	"os"
	"sync"
	"syscall"
	"time"

	manager_io "github.com/456vv/x/vweb_dynamic/internal/wazero/manager/io"
	"github.com/456vv/x/vweb_dynamic/internal/wazero/manager/sockets"
	wasip2_io "github.com/456vv/x/vweb_dynamic/internal/wazero/wasip2/io/v0_2"
	witgo "github.com/456vv/x/vweb_dynamic/internal/wazero/witgo"

	"golang.org/x/sys/unix"
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
		// 只 syscall.Bind。不可 FileListener——
		// 1) 客户端 bind 后还要 connect，Listener 会占住 fd
		// 2) os.NewFile(fd) 会接管所有权，后续 Close/FileConn 会双关或丢掉 Fd
		if !sock.HasFd() {
			return sockets.ErrInvalidSocketState
		}

		bindErr := syscall.Bind(int(sock.Fd), sockaddr)
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

// listenTCP 在已 bind 的 fd 上 listen，再 Dup 后交给 net.TCPListener。
// 必须 Dup，否则 net.FileListener 与 sock.Fd 抢同一个 fd，Drop 会 double-close。
func (i *tcpImpl) listenTCP(sock *sockets.TCPSocket) error {
	if sock.Listener != nil {
		return nil
	}
	if !sock.HasFd() {
		return syscall.EBADF
	}
	if err := unix.Listen(sock.Fd, syscall.SOMAXCONN); err != nil {
		return err
	}
	nfd, err := unix.Dup(sock.Fd)
	if err != nil {
		// listen 已经生效。这里再返回错误的话，DoListen 不会进入 Listening，但内核队列已经存在。
		// acceptTCP 在 Listener 为空时会改走原始 fd。
		return nil
	}
	unix.CloseOnExec(nfd) // Dup 默认可继承，fork/exec 会把 listener 漏给子进程。
	file := os.NewFile(uintptr(nfd), "")
	if file == nil {
		unix.Close(nfd)
		return nil
	}
	ln, err := net.FileListener(file)
	file.Close() // FileListener 已再 dup；关 file 不关 listener
	if err != nil {
		return nil
	}
	tl, ok := ln.(*net.TCPListener)
	if !ok {
		ln.Close()
		return nil
	}
	sock.Listener = tl
	return nil
}

// connectTCP 在已有 fd（可能已 bind）上 Connect，再 Dup+FileConn。
// 禁止 net.DialTCP 另开套接字，否则 start-bind 的本地地址丢失。
func (i *tcpImpl) connectTCP(sock *sockets.TCPSocket, remoteAddress IPSocketAddress) (*net.TCPConn, error) {
	fd, hasFd := sock.SnapshotFd()
	if !hasFd {
		addr, err := fromIPSocketAddressToTCPAddr(remoteAddress)
		if err != nil {
			return nil, err
		}
		network := "tcp"
		if sock.Family == sockets.IPAddressFamilyIPV6 {
			network = "tcp6"
		}
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
	owned := true
	defer func() {
		if owned {
			_ = unix.Close(fd)
		}
	}()

	sa, err := fromIPSocketAddressToSockaddr(remoteAddress)
	if err != nil {
		return nil, err
	}

	for {
		cerr := syscall.Connect(fd, sa)
		if cerr == nil {
			break
		}
		if cerr == unix.EINTR {
			continue
		}
		if cerr == unix.EINPROGRESS || cerr == unix.EAGAIN || cerr == unix.EWOULDBLOCK {
			fds := []unix.PollFd{{Fd: int32(fd), Events: unix.POLLOUT}}
			ctx := sock.ConnectContext()
			for {
				if ctx != nil {
					if err := ctx.Err(); err != nil {
						return nil, err
					}
				}
				n, perr := unix.Poll(fds, 200)
				if perr == unix.EINTR {
					continue
				}
				if perr != nil {
					return nil, perr
				}
				if n > 0 {
					break
				}
			}
			soerr, gerr := unix.GetsockoptInt(fd, unix.SOL_SOCKET, unix.SO_ERROR)
			if gerr != nil {
				return nil, gerr
			}
			if soerr != 0 {
				return nil, syscall.Errno(soerr)
			}
			break
		}
		return nil, cerr
	}

	// 这是锁内 Dup 出来的私有 fd，可以交给 NewFile。不要再对原 fd 做第二次 Dup。
	file := os.NewFile(uintptr(fd), "")
	if file == nil {
		return nil, syscall.EINVAL
	}
	owned = false
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
	_, hasFd := sock.SnapshotFd()
	if hasFd {
		fd, err := sock.DupOwnedFd()
		if err != nil {
			return nil, err
		}
		defer unix.Close(fd) // 修改原因：Accept 用私有 dup，原 fd 的关闭与号码复用不再影响这次调用
		nfd, _, err := unix.Accept(fd)
		if err != nil {
			return nil, err
		}
		unix.CloseOnExec(nfd)
		if err := unix.SetNonblock(nfd, true); err != nil {
			_ = unix.Close(nfd)
			return nil, err
		}
		file := os.NewFile(uintptr(nfd), "")
		if file == nil {
			_ = unix.Close(nfd)
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

	ln := sock.GetListener()
	if ln == nil {
		return nil, syscall.EINVAL
	}
	// SetDeadline(time.Now()) 在队列非空时仍会立刻 timeout
	_ = ln.SetDeadline(time.Now().Add(time.Millisecond))
	conn, err := ln.AcceptTCP()
	_ = ln.SetDeadline(time.Time{})
	return conn, err
}

func (i *tcpImpl) localAddrFromFd(sock *sockets.TCPSocket) (net.Addr, error) {
	var sa unix.Sockaddr
	err := sock.ControlFd(func(fd int) error {
		var gerr error
		sa, gerr = unix.Getsockname(fd)
		return gerr
	})
	if err != nil {
		return nil, err
	}
	return sockaddrToNetTCPAddr(sa)
}

func sockaddrToNetTCPAddr(sa unix.Sockaddr) (net.Addr, error) {
	switch a := sa.(type) {
	case *unix.SockaddrInet4:
		return &net.TCPAddr{IP: net.IP(a.Addr[:]), Port: a.Port}, nil
	case *unix.SockaddrInet6:
		return &net.TCPAddr{IP: net.IP(a.Addr[:]), Port: a.Port, Zone: zoneFromID(a.ZoneId)}, nil
	default:
		return nil, syscall.EAFNOSUPPORT
	}
}

func zoneFromID(id uint32) string {
	if id == 0 {
		return ""
	}
	ifi, err := net.InterfaceByIndex(int(id))
	if err != nil {
		return ""
	}
	return ifi.Name
}

func (i *tcpImpl) SetListenBacklogSize(ctx context.Context, this TCPSocket, value uint64) witgo.Result[witgo.Unit, ErrorCode] {
	sock, ok := i.host.TCPSocketManager().Get(this)
	if !ok {
		return witgo.Err[witgo.Unit, ErrorCode](ErrorCodeInvalidArgument)
	}
	n := clampToInt(value)
	ctrlErr := sock.ControlFd(func(fd int) error {
		return unix.Listen(fd, n)
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
		listenErr = unix.Listen(int(fd), n)
	}); ctrlErr != nil {
		return witgo.Err[witgo.Unit, ErrorCode](mapOsError(ctrlErr))
	}
	if listenErr != nil {
		return witgo.Err[witgo.Unit, ErrorCode](mapOsError(listenErr))
	}
	return witgo.Ok[witgo.Unit, ErrorCode](witgo.Unit{})
}

func (i *tcpImpl) KeepAliveEnabled(ctx context.Context, this TCPSocket) witgo.Result[bool, ErrorCode] {
	result := getsockoptInt[int](i, this, unix.SOL_SOCKET, unix.SO_KEEPALIVE)
	if result.Err != nil {
		return witgo.Err[bool, ErrorCode](*result.Err)
	}
	return witgo.Ok[bool, ErrorCode](*result.Ok == 1)
}

func (i *tcpImpl) setKeepAliveEnabledUnconnected(this TCPSocket, value int) witgo.Result[witgo.Unit, ErrorCode] {
	return setsockoptInt(i, this, unix.SOL_SOCKET, unix.SO_KEEPALIVE, value)
}

func (i *tcpImpl) KeepAliveInterval(ctx context.Context, this TCPSocket) witgo.Result[uint64, ErrorCode] {
	// WASI 该选项是纳秒，TCP_KEEPINTVL 是秒
	result := getsockoptInt[int](i, this, unix.IPPROTO_TCP, unix.TCP_KEEPINTVL)
	if result.Err != nil {
		return witgo.Err[uint64, ErrorCode](*result.Err)
	}
	return witgo.Ok[uint64, ErrorCode](sockoptSecondsToNs(*result.Ok))
}

func (i *tcpImpl) SetKeepAliveInterval(ctx context.Context, this TCPSocket, value uint64) witgo.Result[witgo.Unit, ErrorCode] {
	return setsockoptInt(i, this, unix.IPPROTO_TCP, unix.TCP_KEEPINTVL, nsToSockoptSeconds(value))
}

func (i *tcpImpl) KeepAliveCount(ctx context.Context, this TCPSocket) witgo.Result[uint32, ErrorCode] {
	return getsockoptInt[uint32](i, this, unix.IPPROTO_TCP, unix.TCP_KEEPCNT)
}

func (i *tcpImpl) SetKeepAliveCount(ctx context.Context, this TCPSocket, value uint32) witgo.Result[witgo.Unit, ErrorCode] {
	return setsockoptInt(i, this, unix.IPPROTO_TCP, unix.TCP_KEEPCNT, clampToInt(uint64(value)))
}

func (i *tcpImpl) HopLimit(ctx context.Context, this TCPSocket) witgo.Result[uint8, ErrorCode] {
	sock, ok := i.host.TCPSocketManager().Get(this)
	if !ok {
		return witgo.Err[uint8, ErrorCode](ErrorCodeInvalidArgument)
	}
	if sock.Family == sockets.IPAddressFamilyIPV6 {
		return getsockoptInt[uint8](i, this, unix.IPPROTO_IPV6, unix.IPV6_UNICAST_HOPS)
	}
	return getsockoptInt[uint8](i, this, unix.IPPROTO_IP, unix.IP_TTL)
}

func (i *tcpImpl) SetHopLimit(ctx context.Context, this TCPSocket, value uint8) witgo.Result[witgo.Unit, ErrorCode] {
	sock, ok := i.host.TCPSocketManager().Get(this)
	if !ok {
		return witgo.Err[witgo.Unit, ErrorCode](ErrorCodeInvalidArgument)
	}
	if sock.Family == sockets.IPAddressFamilyIPV6 {
		return setsockoptInt(i, this, unix.IPPROTO_IPV6, unix.IPV6_UNICAST_HOPS, int(value))
	}
	return setsockoptInt(i, this, unix.IPPROTO_IP, unix.IP_TTL, int(value))
}

func (i *tcpImpl) ReceiveBufferSize(ctx context.Context, this TCPSocket) witgo.Result[uint64, ErrorCode] {
	return getsockoptInt[uint64](i, this, unix.SOL_SOCKET, unix.SO_RCVBUF)
}

func (i *tcpImpl) SendBufferSize(ctx context.Context, this TCPSocket) witgo.Result[uint64, ErrorCode] {
	return getsockoptInt[uint64](i, this, unix.SOL_SOCKET, unix.SO_SNDBUF)
}

func (i *tcpImpl) setReceiveBufferUnconnected(this TCPSocket, value uint64) witgo.Result[witgo.Unit, ErrorCode] {
	return setsockoptInt(i, this, unix.SOL_SOCKET, unix.SO_RCVBUF, clampToInt(value))
}

func (i *tcpImpl) setSendBufferUnconnected(this TCPSocket, value uint64) witgo.Result[witgo.Unit, ErrorCode] {
	return setsockoptInt(i, this, unix.SOL_SOCKET, unix.SO_SNDBUF, clampToInt(value))
}

// subscribeListen 在已 listen 的 fd 上等待 POLLIN（accept 队列非空或出错）。
//
//	WASI tcp-socket.subscribe 在 listening 状态应对应“可 accept”；
//
// unix.Poll 在 linux/darwin/BSD/solaris 行为一致，Darwin 无 SOCK_NONBLOCK 也不影响 poll。
func (i *tcpImpl) subscribeListen(sock *sockets.TCPSocket) wasip2_io.Pollable {
	fd, err := sock.DupOwnedFd()
	if err != nil {
		p := manager_io.NewPollable(nil)
		handle := i.host.PollManager().Add(p)
		p.SetReady()
		return handle
	}

	done := make(chan struct{})
	var once sync.Once
	p := manager_io.NewPollable(func() {
		once.Do(func() { close(done) })
	})
	handle := i.host.PollManager().Add(p)

	go func() {
		// 不要在 Poll 进行中 close(fd)，否则号码复用会让 Poll 观察别的 socket。
		// 200ms 超时用来观察 done，退出后再关闭私有 dup。
		defer unix.Close(fd)
		defer p.SetReady()
		fds := []unix.PollFd{{
			Fd:     int32(fd),
			Events: unix.POLLIN | unix.POLLERR | unix.POLLHUP | unix.POLLNVAL,
		}}
		for {
			select {
			case <-done:
				return
			default:
			}
			n, err := unix.Poll(fds, 200)
			if err == unix.EINTR {
				continue
			}
			if err != nil || n > 0 {
				return
			}
		}
	}()
	return handle
}

func getsockoptInt[T ~int | ~uint64 | ~uint32 | ~uint8](i *tcpImpl, this TCPSocket, level, opt int) witgo.Result[T, ErrorCode] {
	sock, ok := i.host.TCPSocketManager().Get(this)
	if !ok || sock == nil {
		return witgo.Err[T, ErrorCode](ErrorCodeInvalidArgument)
	}

	var val int
	var getErr error
	// GetConn 与 SnapshotFd 分两次加锁，中间 bind/connect/drop 会造成两端都 miss 或对已关 fd 调 getsockopt。
	conn, _, hasFd := sock.SnapshotConnOrFd()
	switch {
	case conn != nil:
		rawConn, err := conn.SyscallConn()
		if err != nil {
			return witgo.Err[T, ErrorCode](mapOsError(err))
		}
		err = rawConn.Control(func(fd uintptr) {
			val, getErr = unix.GetsockoptInt(int(fd), level, opt)
		})
		if err != nil {
			return witgo.Err[T, ErrorCode](mapOsError(err))
		}
	case hasFd:
		// WASI 允许在 bind/connect 前读 hop-limit 等；此时还没有 net.TCPConn
		getErr = sock.ControlFd(func(fd int) error {
			var err error
			val, err = unix.GetsockoptInt(fd, level, opt)
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

func setsockoptInt(i *tcpImpl, this TCPSocket, level, opt, value int) witgo.Result[witgo.Unit, ErrorCode] {
	sock, ok := i.host.TCPSocketManager().Get(this)
	if !ok || sock == nil {
		return witgo.Err[witgo.Unit, ErrorCode](ErrorCodeInvalidArgument)
	}

	var setErr error
	// 与 getsockoptInt 相同，必须在同一快照上选 Conn 或 fd。
	conn, _, hasFd := sock.SnapshotConnOrFd()
	switch {
	case conn != nil:
		rawConn, err := conn.SyscallConn()
		if err != nil {
			return witgo.Err[witgo.Unit, ErrorCode](mapOsError(err))
		}
		err = rawConn.Control(func(fd uintptr) {
			setErr = unix.SetsockoptInt(int(fd), level, opt, value)
		})
		if err != nil {
			return witgo.Err[witgo.Unit, ErrorCode](mapOsError(err))
		}
	case hasFd:
		// 与 getsockoptInt 对称，connect 前设置 TTL/缓冲区/keepalive
		setErr = sock.ControlFd(func(fd int) error {
			return unix.SetsockoptInt(fd, level, opt, value)
		})
	default:
		return witgo.Err[witgo.Unit, ErrorCode](ErrorCodeInvalidArgument)
	}
	if setErr != nil {
		return witgo.Err[witgo.Unit, ErrorCode](mapOsError(setErr))
	}

	return witgo.Ok[witgo.Unit, ErrorCode](witgo.Unit{})
}
