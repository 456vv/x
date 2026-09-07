package witgo

import (
	"math"
	"sync"
)

// ============================================================
// ResourceManager：Wasm 资源句柄管理
// ============================================================

// DestructorFunc 定义资源清理函数的签名。
// 当资源被删除时调用，用于释放底层资源（文件、连接等）。
type DestructorFunc[T any] func(resource T)

// ResourceManager[T] 是一个泛型的、线程安全的句柄表，用于管理宿主拥有的 T 类型资源。
// 每个资源通过唯一句柄（uint32）访问，句柄由管理器分配。
//
// 线程安全：所有方法均使用 RWMutex 保护。
type ResourceManager[T any] struct {
	mu         sync.RWMutex      // 读写锁
	handles    map[uint32]T      // 句柄 -> 资源的映射
	nextID     uint32            // 下一个候选句柄 ID
	destructor DestructorFunc[T] // 资源删除时的清理函数
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

// Set 将指定 handle 关联的资源设置为 resource，若 handle 已存在则替换旧资源（并调用 destructor）。
//
// 参数:
//   - handle:   资源句柄
//   - resource: 要关联的资源值
func (m *ResourceManager[T]) Set(handle uint32, resource T) {
	if m == nil {
		return
	}
	m.mu.Lock()
	old, exists := m.handles[handle]
	m.handles[handle] = resource
	if m.nextID < handle {
		m.nextID = handle
	}
	dt := m.destructor
	m.mu.Unlock()
	if exists && dt != nil {
		dt(old) // 释放旧资源
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
	m.mu.Lock()
	defer m.mu.Unlock()

	if uint64(len(m.handles)) >= math.MaxUint32-1 {
		panic("ResourceManager handle table is full")
	}

	m.nextID++
	if m.nextID == 0 {
		m.nextID = 1
	}
	start := m.nextID
	for {
		if _, exists := m.handles[m.nextID]; !exists {
			break
		}
		m.nextID++
		if m.nextID == 0 {
			m.nextID = 1
		}
		if m.nextID == start {
			panic("ResourceManager handle table is full")
		}
	}
	handle := m.nextID
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
	dt := m.destructor
	m.mu.Unlock()

	if dt != nil {
		dt(res) // 清理资源
	}
	return true
}

// Pop 移除并返回句柄对应的资源，若成功则调用 destructor。
// 与 Remove 不同，Pop 返回被移除的资源值。
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
	dt := m.destructor
	m.mu.Unlock()

	if dt != nil {
		dt(res)
	}
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
