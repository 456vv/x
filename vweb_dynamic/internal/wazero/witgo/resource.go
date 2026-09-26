package witgo

import (
	"math"
	"sync"
)

// DestructorFunc 定义资源清理函数的签名。
// 当资源被删除时调用，用于释放底层资源（文件、连接等）。
type DestructorFunc[T any] func(resource T)

// ResourceManager[T] 是一个泛型的、线程安全的句柄表，用于管理宿主拥有的 T 类型资源。
type ResourceManager[T any] struct {
	mu         sync.RWMutex      // 读写锁
	handles    map[uint32]T      // 句柄 -> 资源的映射
	nextID     uint32            // 下一个候选句柄 ID
	free       []uint32          // Remove/Pop 后回收句柄，避免 nextID 只增并在空洞上线性探测
	destructor DestructorFunc[T] // 资源删除时的清理函数
	closed     bool              // Clear 后禁止再 Add/Set，避免 Close 与 guest 并发路径泄漏新资源
}

// NewResourceManager 创建一个新的资源管理器，可选的 destructor 在删除资源时调用。
//
// 参数:
//   - destructor: 资源清理函数，可为 nil
//
// 返回:
//   - *ResourceManager[T]: 初始化的管理器
//
// 示例:
//
//	rm := witgo.NewResourceManager[myResource](func(r myResource) { r.Close() })
func NewResourceManager[T any](destructor DestructorFunc[T]) *ResourceManager[T] {
	return &ResourceManager[T]{
		handles:    make(map[uint32]T),
		nextID:     0,
		destructor: destructor,
	}
}

func (m *ResourceManager[T]) recycleLocked(handle uint32) {
	if handle == 0 {
		return
	}
	// 极端 churn 时 free 切片 cap 只增不减，收缩常驻内存
	if cap(m.free) > 1024 && len(m.free) < cap(m.free)/4 {
		nfree := make([]uint32, len(m.free), len(m.free)+1)
		copy(nfree, m.free)
		m.free = nfree
	}
	m.free = append(m.free, handle)
}

func (m *ResourceManager[T]) dropFromFreeLocked(handle uint32) {
	if handle == 0 || len(m.free) == 0 {
		return
	}
	for i := len(m.free) - 1; i >= 0; i-- {
		if m.free[i] == handle {
			m.free = append(m.free[:i], m.free[i+1:]...)
			return
		}
	}
}

func (m *ResourceManager[T]) allocHandleLocked() uint32 {
	// 从尾部弹出空闲句柄；Set 可能占用仍留在 free 里的 id，脏项跳过继续弹
	for len(m.free) > 0 {
		h := m.free[len(m.free)-1]
		m.free = m.free[:len(m.free)-1]
		if h == 0 {
			continue
		}
		if _, exists := m.handles[h]; !exists {
			return h
		}
	}

	m.nextID++
	if m.nextID == 0 {
		m.nextID = 1 // 句柄 0 在 WASI 中表示无效资源，跳过
	}
	start := m.nextID
	for {
		if _, exists := m.handles[m.nextID]; !exists {
			return m.nextID
		}
		m.nextID++
		if m.nextID == 0 {
			m.nextID = 1
		}
		if m.nextID == start {
			panic("ResourceManager handle table is full")
		}
	}
}

// Set 将指定 handle 关联的资源设置为 resource，若 handle 已存在则替换旧资源（并调用 destructor）。
//
// 参数:
//   - handle:   资源句柄
//   - resource: 要关联的资源值
func (m *ResourceManager[T]) Set(handle uint32, resource T) {
	if m == nil {
		return
	}
	if handle == 0 {
		// WASI 句柄 0 表示无效资源；写入 0 会让 Get(0) 误成功。
		return
	}
	m.mu.Lock()
	old, exists := m.handles[handle]
	m.handles[handle] = resource
	// 记录已使用的最大句柄，避免随后 Add 从较小 nextID 扫描过久或产生混淆
	if handle >= m.nextID {
		m.nextID = handle
	}
	if !exists {
		m.dropFromFreeLocked(handle)
	}
	dt := m.destructor
	m.mu.Unlock()
	if exists && dt != nil {
		dt(old) // 析构必须在锁外，避免 Close 再入管理器死锁
	}
}

// Add 添加新资源并分配一个唯一句柄（从 1 开始），若表满则 panic。
// 句柄分配策略：优先使用 nextID，若已存在则递增查找空位。
//
// 参数:
//   - resource: 要添加的资源值
//
// 返回:
//   - uint32: 分配的句柄（ caller 需保存以便后续访问）
//
// 示例:
//
//	handle := rm.Add(myResource{})
func (m *ResourceManager[T]) Add(resource T) uint32 {
	if m == nil {
		panic("ResourceManager is nil")
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.closed {
		panic("ResourceManager is closed")
	}
	if uint64(len(m.handles)) >= math.MaxUint32-1 {
		panic("ResourceManager handle table is full")
	}
	handle := m.allocHandleLocked()
	m.handles[handle] = resource
	return handle
}

// Get 根据句柄获取资源，若不存在则返回零值和 false。
//
// 参数:
//   - handle: 资源句柄
//
// 返回:
//   - T:     资源值（不存在时为零值）
//   - bool:  是否存在
func (m *ResourceManager[T]) Get(handle uint32) (T, bool) {
	var zero T
	if m == nil {
		return zero, false
	}
	m.mu.RLock()
	defer m.mu.RUnlock()
	res, ok := m.handles[handle]
	return res, ok
}

// Remove 删除句柄对应的资源，返回是否成功，若成功则调用 destructor。
//
// 参数:
//   - handle: 要删除的资源句柄
//
// 返回:
//   - bool: 是否成功删除
func (m *ResourceManager[T]) Remove(handle uint32) bool {
	if m == nil {
		return false
	}
	m.mu.Lock()
	res, ok := m.handles[handle]
	if !ok {
		m.mu.Unlock()
		return false
	}
	delete(m.handles, handle)
	m.recycleLocked(handle)
	dt := m.destructor
	m.mu.Unlock()

	if dt != nil {
		dt(res) // 清理资源
	}
	return true
}

// RemoveIf 仅当 handle 仍指向 match 认可的那个资源时才删除。
// 句柄回收后会被复用；先 Get 再 Remove 的窗口里可能已经换成新资源。
// match 只能比较指针或值，不能再进 ResourceManager，否则自锁。
func (m *ResourceManager[T]) RemoveIf(handle uint32, match func(T) bool) bool {
	return m.RemoveIfWith(handle, match, nil)
}

// RemoveIfWith 与 RemoveIf 相同，但在仍持有表锁、句柄尚未放回 free 列表时调用 beforeRecycle。
// 先 Unlock 再 UnmarkFieldsImmutable 时，handle 可能已被 Add 复用，
// 会清掉新 incoming headers 的不可变标记。
func (m *ResourceManager[T]) RemoveIfWith(handle uint32, match func(T) bool, beforeRecycle func(T)) bool {
	if m == nil || match == nil {
		return false
	}
	m.mu.Lock()
	res, ok := m.handles[handle]
	if !ok || !match(res) {
		m.mu.Unlock()
		return false
	}
	delete(m.handles, handle)
	if beforeRecycle != nil {
		beforeRecycle(res)
	}
	m.recycleLocked(handle)
	dt := m.destructor
	m.mu.Unlock()
	if dt != nil {
		dt(res)
	}
	return true
}

// Pop 移除并返回句柄对应的资源（不调用 destructor）。
// 与 Remove 不同，Pop 返回被移除的资源值。
// 1) response-outparam.set 向已关闭 channel 发送而 panic
// 2) outgoing-handler.handle 把尚未写完的 PipeReader 提前 Close
// 3) TLS handshake.finish 在握手前关掉底层流
// 资源释放仍由 Remove/Drop 或调用方显式 Close 完成。
//
// 参数:
//   - handle: 要移除的资源句柄
//
// 返回:
//   - T:     被移除的资源值
//   - bool:  是否成功移除
func (m *ResourceManager[T]) Pop(handle uint32) (T, bool) {
	var zero T
	if m == nil {
		return zero, false
	}
	m.mu.Lock()
	res, ok := m.handles[handle]
	if !ok {
		m.mu.Unlock()
		return zero, false
	}
	delete(m.handles, handle)
	m.recycleLocked(handle)
	m.mu.Unlock()
	return res, true
}

// Range 遍历所有句柄-资源对，回调函数返回 false 时停止遍历。
// 遍历在锁外进行（先复制所有条目），因此遍历时新增/删除的元素不会被看到。
//
// 参数:
//   - f: 回调函数，参数为 (handle, resource)，返回 false 时停止
func (m *ResourceManager[T]) Range(f func(handle uint32, resource T) bool) {
	if m == nil || f == nil {
		return
	}
	type entry struct {
		h uint32
		r T
	}

	m.mu.RLock()
	entries := make([]entry, 0, len(m.handles))
	for handle, resource := range m.handles {
		entries = append(entries, entry{h: handle, r: resource})
	}
	m.mu.RUnlock()

	for _, e := range entries {
		if !f(e.h, e.r) {
			break
		}
	}
}

// Clear 释放全部资源。
// Host 关闭时必须并发安全地析构，避免泄漏 fd/conn/goroutine。
func (m *ResourceManager[T]) Clear() {
	if m == nil {
		return
	}
	m.mu.Lock()
	m.closed = true
	handles := m.handles
	m.handles = make(map[uint32]T)
	m.nextID = 0
	m.free = nil
	dt := m.destructor
	m.mu.Unlock()
	if dt == nil {
		return
	}
	for _, res := range handles {
		dt(res)
	}
}
