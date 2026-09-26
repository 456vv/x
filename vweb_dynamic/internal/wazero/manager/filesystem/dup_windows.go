//go:build windows

package filesystem

import (
	"os"

	"golang.org/x/sys/windows"
)

// DupFile 复制文件 HANDLE。返回的 *os.File 由 CloseHandle 释放。
// 修改原因：与 unix dup 相同，只给 ReadAt 流使用；目录列举不能依赖 DuplicateHandle。
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
	var nfd windows.Handle
	var sysErr error
	p := windows.CurrentProcess()
	if err = raw.Control(func(fd uintptr) {
		sysErr = windows.DuplicateHandle(p, windows.Handle(fd), p, &nfd, 0, false, windows.DUPLICATE_SAME_ACCESS)
	}); err != nil {
		return nil, err
	}
	if sysErr != nil {
		return nil, sysErr
	}
	f := os.NewFile(uintptr(nfd), d.File.Name())
	if f == nil {
		_ = windows.CloseHandle(nfd)
		return nil, os.ErrInvalid
	}
	return f, nil
}

// openDirectorySnapshot 按路径重新打开目录。
// 按名打开存在 TOCTOU，是无 openat 时的尽力而为。
func openDirectorySnapshot(f *os.File) (*os.File, error) {
	if f == nil {
		return nil, os.ErrInvalid
	}
	name := f.Name()
	if name == "" {
		return nil, os.ErrInvalid
	}
	return os.Open(name)
}
