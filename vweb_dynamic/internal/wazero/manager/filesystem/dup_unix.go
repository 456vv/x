//go:build unix

package filesystem

import (
	"os"

	"golang.org/x/sys/unix"
)

// DupFile 复制独立 fd。调用方必须 Close 返回值。
// read-via-stream 使用 ReadAt，不依赖偏移；descriptor 先关闭时不能关掉流仍在读的 fd。
// 不能用于 ReadDir：dup 与原 fd 共享目录偏移。
func (d *Descriptor) DupFile() (*os.File, error) {
	if d == nil {
		return nil, os.ErrInvalid
	}
	d.mu.RLock()
	defer d.mu.RUnlock()
	if d.File == nil {
		if d.closed {
			return nil, os.ErrClosed
		}
		return nil, os.ErrInvalid
	}
	raw, err := d.File.SyscallConn()
	if err != nil {
		return nil, err
	}
	var nfd int
	var sysErr error
	if err = raw.Control(func(fd uintptr) {
		nfd, sysErr = unix.Dup(int(fd))
	}); err != nil {
		return nil, err
	}
	if sysErr != nil {
		return nil, sysErr
	}
	unix.CloseOnExec(nfd)
	f := os.NewFile(uintptr(nfd), d.File.Name())
	if f == nil {
		unix.Close(nfd)
		return nil, os.ErrInvalid
	}
	return f, nil
}

// openDirectorySnapshot 用 openat(dirfd, ".") 打开新的目录 fd，偏移从 0 起。
func openDirectorySnapshot(f *os.File) (*os.File, error) {
	if f == nil {
		return nil, os.ErrInvalid
	}
	raw, err := f.SyscallConn()
	if err != nil {
		return nil, err
	}
	var nfd int
	var sysErr error
	if err = raw.Control(func(fd uintptr) {
		nfd, sysErr = unix.Openat(int(fd), ".", unix.O_RDONLY|unix.O_DIRECTORY|unix.O_CLOEXEC, 0)
	}); err != nil {
		return nil, err
	}
	if sysErr != nil {
		return nil, sysErr
	}
	unix.CloseOnExec(nfd) // 部分 unix 上 Openat 可能不认 O_CLOEXEC。
	nf := os.NewFile(uintptr(nfd), f.Name())
	if nf == nil {
		unix.Close(nfd)
		return nil, os.ErrInvalid
	}
	return nf, nil
}
