//go:build !unix && !windows

package filesystem

import "os"

// DupFile 本平台没有可移植 dup。文件仍打开时返回 os.ErrInvalid，调用方改走加锁读。
// 不能把“无 dup”报成 os.ErrClosed，否则未关闭文件会被当成 bad-descriptor。
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
	return nil, os.ErrInvalid
}

// openDirectorySnapshot 本平台无 openat，按路径重开。
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
