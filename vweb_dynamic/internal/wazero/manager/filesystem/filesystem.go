package filesystem

import (
	"io"
	"io/fs"
	"os"
	"sync"

	witgo "github.com/456vv/x/vweb_dynamic/internal/wazero/witgo"
)

// Descriptor represents a file or directory descriptor.
// It holds the underlying os.File and the pre-opened path.
type Descriptor struct {
	// The underlying file descriptor from Go's os package.
	File *os.File
	// The path this descriptor was pre-opened with, for identification.
	Path string
	// Permissions associated with this descriptor.
	Flags fs.FileMode

	// IsPreopen 为 true 时才出现在 get-directories。
	// 仅凭 IsDir 会把 open-at 得到的普通目录当成 preopen。
	IsPreopen bool

	// 回调内不能再调用本类型会加锁的方法，RWMutex 不可重入。
	mu     sync.RWMutex
	closed bool
}

// NewPreopenDescriptor 构造应出现在 get-directories 中的预打开目录。
func NewPreopenDescriptor(file *os.File, guestPath string) *Descriptor {
	return &Descriptor{File: file, Path: guestPath, IsPreopen: true}
}

// ForPreopen 报告 get-directories 是否应返回该项。
// Close 后 File 为 nil。把所有 nil File 都当目录，会让已关闭描述符重新出现。
// 必须同时要求 IsPreopen，否则普通目录描述符会变成 preopen。
func (d *Descriptor) ForPreopen() bool {
	if d == nil {
		return false
	}
	d.mu.RLock()
	defer d.mu.RUnlock()
	if d.closed || !d.IsPreopen || d.File == nil {
		return false
	}
	if d.File == nil {
		return true
	}
	st, err := d.File.Stat()
	return err == nil && st.IsDir()
}

// ReadDir 列举目录。
// ReadDir 本身也会移动偏移，所以用写锁与 ReadAt/Close 互斥。
// 用 openDirectorySnapshot（unix: openat(fd,".")）拿到独立偏移，失败再 Seek 回退。
func (d *Descriptor) ReadDir(n int) ([]fs.DirEntry, error) {
	var entries []fs.DirEntry
	err := d.Do(func(f *os.File) error {
		df, snapErr := openDirectorySnapshot(f)
		if snapErr != nil {
			_, _ = f.Seek(0, io.SeekStart)
			var readErr error
			entries, readErr = f.ReadDir(n)
			return readErr
		}
		defer df.Close()
		var readErr error
		entries, readErr = df.ReadDir(n)
		return readErr
	})
	return entries, err
}

// DoR 在读锁内使用文件，只用于不移动 fd 偏移的操作。
func (d *Descriptor) DoR(fn func(*os.File) error) error {
	if d == nil || fn == nil {
		return os.ErrInvalid
	}
	d.mu.RLock()
	defer d.mu.RUnlock()
	if d.File == nil {
		if d.closed {
			return os.ErrClosed
		}
		return os.ErrInvalid
	}
	return fn(d.File)
}

// Do 在写锁内使用文件。追加、截断和按偏移写必须独占，避免写到旧尾部。
func (d *Descriptor) Do(fn func(*os.File) error) error {
	if d == nil || fn == nil {
		return os.ErrInvalid
	}
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.File == nil {
		if d.closed {
			return os.ErrClosed
		}
		return os.ErrInvalid
	}
	return fn(d.File)
}

// Name 复制路径。修改原因：关闭后不能再无锁读 File。
func (d *Descriptor) Name() (string, error) {
	var name string
	err := d.DoR(func(f *os.File) error {
		name = f.Name()
		return nil
	})
	return name, err
}

// Stat 与 Close 互斥。
func (d *Descriptor) Stat() (fs.FileInfo, error) {
	var info fs.FileInfo
	err := d.DoR(func(f *os.File) error {
		var statErr error
		info, statErr = f.Stat()
		return statErr
	})
	return info, err
}

// ReadAt 使用 pread，不移动共享 fd 的 seek。
func (d *Descriptor) ReadAt(p []byte, off int64) (int, error) {
	var n int
	err := d.DoR(func(f *os.File) error {
		var readErr error
		n, readErr = f.ReadAt(p, off)
		return readErr
	})
	return n, err
}

// Sync 对应 WASI sync。
func (d *Descriptor) Sync() error {
	return d.DoR(func(f *os.File) error { return f.Sync() })
}

// AppendWrite 在文件当前大小处 WriteAt，不改 fd 的 seek 位置。
func (d *Descriptor) AppendWrite(p []byte) (int, error) {
	var n int
	err := d.Do(func(f *os.File) error {
		st, statErr := f.Stat()
		if statErr != nil {
			return statErr
		}
		off := st.Size()
		if off < 0 {
			off = 0
		}
		var writeErr error
		n, writeErr = f.WriteAt(p, off)
		return writeErr
	})
	return n, err
}

// WriteAt 与 AppendWrite、Truncate、Close 共用写锁。
func (d *Descriptor) WriteAt(p []byte, off int64) (int, error) {
	var n int
	err := d.Do(func(f *os.File) error {
		var writeErr error
		n, writeErr = f.WriteAt(p, off)
		return writeErr
	})
	return n, err
}

// Truncate 与 AppendWrite 共用写锁，避免缩文件和追加交错。
func (d *Descriptor) Truncate(size int64) error {
	return d.Do(func(f *os.File) error { return f.Truncate(size) })
}

// Close 先摘下 *os.File 再关闭，可重复调用。
func (d *Descriptor) Close() error {
	if d == nil {
		return nil
	}
	d.mu.Lock()
	f := d.File
	d.File = nil
	d.closed = true
	d.mu.Unlock()
	if f == nil {
		return nil
	}
	return f.Close()
}

// Manager is the resource manager for all filesystem descriptors.
type Manager = witgo.ResourceManager[*Descriptor]

// NewManager creates a new filesystem descriptor manager.
func NewManager() *Manager {
	return witgo.NewResourceManager[*Descriptor](func(resource *Descriptor) {
		if resource != nil {
			resource.Close()
		}
	})
}

// DirectoryEntryStreamState 用于管理读取目录的状态。
type DirectoryEntryStreamState struct {
	mu sync.Mutex // read-directory-entry 可能被 guest 并发调用
	// 预读的目录条目
	Entries []fs.DirEntry
	// 当前读取到的索引位置
	Index int
}

// Next 取出下一条目录项；没有更多条目时返回 false。
func (s *DirectoryEntryStreamState) Next() (fs.DirEntry, bool) {
	if s == nil {
		return nil, false
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.Index >= len(s.Entries) {
		return nil, false
	}
	e := s.Entries[s.Index]
	s.Index++
	return e, true
}

// DirectoryEntryStreamManager 是用于管理目录条目流的资源管理器。
type DirectoryEntryStreamManager = witgo.ResourceManager[*DirectoryEntryStreamState]

// NewDirectoryEntryStreamManager 创建一个新的目录流管理器。
func NewDirectoryEntryStreamManager() *DirectoryEntryStreamManager {
	return witgo.NewResourceManager[*DirectoryEntryStreamState](nil)
}
