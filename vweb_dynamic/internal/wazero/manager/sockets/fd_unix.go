//go:build unix

package sockets

import "golang.org/x/sys/unix"

// AttachFd 绑定系统套接字 fd，并登记平台相关的关闭函数。
// closeFd 未导出，必须在本包内赋值；unix.Close 不能写进无 build tag 的 sockets.go。
// fd==0 可能是合法套接字（stdin 关闭后内核会复用 0），不能当无效哨兵；无效用 fd<0。
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
		err := unix.Close(fd)
		s.Fd = -1 // 关闭后用 -1 标记，避免与合法 fd 0 混淆
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
		err := unix.Close(fd)
		s.Fd = -1
		return err
	}
}
