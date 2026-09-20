package filesystem

import (
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
}

// Manager is the resource manager for all filesystem descriptors.
type Manager = witgo.ResourceManager[*Descriptor]

// NewManager creates a new filesystem descriptor manager.
func NewManager() *Manager {
	return witgo.NewResourceManager[*Descriptor](func(resource *Descriptor) {
		// 在资源被释放前关闭文件描述符，确保资源正确释放。关闭前判空，避免异常资源导致 panic
		if resource != nil && resource.File != nil {
			_ = resource.File.Close()
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
