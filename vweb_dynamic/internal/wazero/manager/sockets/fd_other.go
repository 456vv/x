//go:build !unix && !windows

package sockets

// AttachFd 在无原始 fd 的平台上仅为兼容保留；无需 closeFd。
func (s *TCPSocket) AttachFd(fd int) {
	if s != nil {
		s.Fd = fd
		s.closeFd = nil
	}
}

func (s *UDPSocket) AttachFd(fd int) {
	if s != nil {
		s.Fd = fd
		s.closeFd = nil
	}
}
