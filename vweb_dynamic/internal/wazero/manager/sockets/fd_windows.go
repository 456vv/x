//go:build windows

package sockets

import "golang.org/x/sys/windows"

// AttachFd 绑定 SOCKET 句柄；关闭必须用 Closesocket 而非 Close。
// windows.Close 对套接字不正确；closeFd 只能在本包赋值。
// 存成 int 时无效句柄是 INVALID_SOCKET（-1），不是 0。
func (s *TCPSocket) AttachFd(fd int) {
	if s == nil {
		return
	}
	s.Fd = fd
	if fd < 0 {
		s.closeFd = nil
		return
	}
	s.closeFd = func() error {
		return windows.Closesocket(windows.Handle(fd))
	}
}

func (s *UDPSocket) AttachFd(fd int) {
	if s == nil {
		return
	}
	s.Fd = fd
	if fd < 0 {
		s.closeFd = nil
		return
	}
	s.closeFd = func() error {
		return windows.Closesocket(windows.Handle(fd))
	}
}

// DupOwnedFd 在锁内 DuplicateHandle。调用方必须 Closesocket 返回值。
// 同 unix，避免锁外使用已被关闭并复用的 SOCKET。
func (s *TCPSocket) DupOwnedFd() (int, error) {
	if s == nil {
		return -1, ErrInvalidSocketState
	}
	s.connectMu.Lock()
	defer s.connectMu.Unlock()
	if s.closeFd == nil || s.Fd < 0 {
		return -1, ErrInvalidSocketState
	}
	var nfd windows.Handle
	p := windows.CurrentProcess()
	if err := windows.DuplicateHandle(p, windows.Handle(s.Fd), p, &nfd, 0, false, windows.DUPLICATE_SAME_ACCESS); err != nil {
		return -1, err
	}
	return int(nfd), nil
}
