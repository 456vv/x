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
		// 闭包再写 s.Fd 会与 SnapshotFd 发生数据竞争。
		return unix.Close(fd)
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
		// 闭包再写 s.Fd 会与 SnapshotFd 发生数据竞争。
		return unix.Close(fd)
	}
}

// DupOwnedFd 在锁内复制 fd。调用方必须最终 unix.Close 返回值。
// 锁外继续用 SnapshotFd 的号码，和析构 unix.Close 竞态时会操作到别的 socket。
func (s *TCPSocket) DupOwnedFd() (int, error) {
	if s == nil {
		return -1, ErrInvalidSocketState
	}
	s.connectMu.Lock()
	defer s.connectMu.Unlock()
	if s.closeFd == nil || s.Fd < 0 {
		return -1, ErrInvalidSocketState
	}
	nfd, err := unix.Dup(s.Fd)
	if err != nil {
		return -1, err
	}
	unix.CloseOnExec(nfd)
	return nfd, nil
}
