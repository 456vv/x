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
		err := windows.Closesocket(windows.Handle(fd))
		s.Fd = -1
		return err
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
		err := windows.Closesocket(windows.Handle(fd))
		s.Fd = -1
		return err
	}
}
