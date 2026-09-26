package witgo

import (
	"context"
	"fmt"
	"sync"

	"github.com/tetratelabs/wazero/api"
)

// ============================================================
// 内存管理模式检测与基础类型定义
// ============================================================

// allocatorMode 表示检测到的 Guest 内存管理模式。
// Wasm Guest（如 Rust/TinyGo/C）可能使用不同的分配接口，
// 此枚举用于运行时自动适配。
type allocatorMode int

const (
	modeUnsupported allocatorMode = iota // 未检测到支持的分配器
	modeCabiRealloc                      // 标准组件模型: cabi_realloc(ptr, old_size, align, new_size)
	modeRealloc                          // TinyGo 或 C-style: realloc(ptr, new_size)
	modeMallocFree                       // C-style: malloc(size), free(ptr)
)

// allocRec 记录一次由 Host 发起的 Guest 分配，便于调用结束后自动回收，避免泄漏。
// ptr:    分配的内存指针（Guest 地址空间内）
// size:   请求的有效字节数
// align:  对齐要求（2 的幂）
type allocRec struct {
	ptr, size, align uint32
}

// rawAlloc 记录非 cabi 模式下为满足对齐而"过量分配"的原始指针。
// 当 Guest 使用 realloc/malloc 时，需要对齐，因此实际分配大小 > 请求大小。
// raw:    malloc/realloc 返回的原始指针（Free 时必须还回此值）
// size:   实际分配的总字节数（含对齐填充）
// align:  对齐要求
type rawAlloc struct {
	raw   uint32
	size  uint32
	align uint32
}

// allocScopeKey 是私有类型，用作 context key 避免与用户 key 冲突。
// Go 的 context.WithValue 要求 key 类型为 comparable，私有 struct{} 满足此要求。
type allocScopeKey struct{}

// allocScope 挂载在 context 上，保证并发 Host.Call 不会共用同一条分配栈。
// 每次 Call 创建独立的 scope，调用结束后统一释放，防止内存泄漏。
type allocScope struct {
	mu   sync.Mutex // 保护 recs 切片，支持并发写入
	recs []allocRec // 本次调用中所有分配的列表，freeAll 时逆序释放
}

// contextWithAllocScope 创建带分配作用域的 context。
// 参数:
//   - ctx: 原始上下文，可为 nil（自动使用 context.Background()）
//
// 返回:
//   - context.Context: 携带 allocScope 的新 context
//   - *allocScope: 可直接用于 add/remove/freeAll
//
// 示例:
//
//	ctx, scope := contextWithAllocScope(context.Background())
//	defer scope.freeAll(ctx, allocator) // 确保释放
func contextWithAllocScope(ctx context.Context) (context.Context, *allocScope) {
	if ctx == nil {
		ctx = context.Background()
	}
	s := &allocScope{}
	return context.WithValue(ctx, allocScopeKey{}, s), s
}

// scopeFromCtx 从 context 中提取分配作用域，若不存在则返回 nil。
// 参数:
//   - ctx: 包含 allocScope 的上下文
//
// 返回:
//   - *allocScope: 作用域对象，可能为 nil
func scopeFromCtx(ctx context.Context) *allocScope {
	if ctx == nil {
		return nil
	}
	s, _ := ctx.Value(allocScopeKey{}).(*allocScope)
	return s
}

// add 将一条分配记录添加到作用域，供后续释放。
// 参数:
//   - rec: 分配记录，ptr 为 0 时跳过（避免记录无效分配）
func (s *allocScope) add(rec allocRec) {
	if s == nil || rec.ptr == 0 {
		return
	}
	s.mu.Lock()
	s.recs = append(s.recs, rec)
	s.mu.Unlock()
}

// remove 从作用域中移除指定指针的记录（当显式 Free 时调用，避免 double-free）。
// 参数:
//   - ptr: 要移除的分配指针
func (s *allocScope) remove(ptr uint32) {
	if s == nil || ptr == 0 {
		return
	}
	s.mu.Lock()
	for i := len(s.recs) - 1; i >= 0; i-- {
		if s.recs[i].ptr == ptr {
			s.recs = append(s.recs[:i], s.recs[i+1:]...)
			break
		}
	}
	s.mu.Unlock()
}

// freeAll 逆序释放作用域中所有未显式释放的分配，返回第一个遇到的错误。
// 逆序释放遵循"后进先出"原则，避免释放已被覆盖的指针。
// 参数:
//   - ctx: 上下文（用于传递取消信号）
//   - a:   GuestAllocator 实例，用于执行实际释放
//
// 返回:
//   - error: 第一个发生的错误，后续错误被忽略
func (s *allocScope) freeAll(ctx context.Context, a *GuestAllocator) error {
	if s == nil || a == nil {
		return nil
	}
	s.mu.Lock()
	recs := s.recs
	s.recs = nil // 清空，避免重复释放
	s.mu.Unlock()

	var first error
	for i := len(recs) - 1; i >= 0; i-- {
		if err := a.Free(ctx, recs[i].ptr, recs[i].size, recs[i].align); err != nil && first == nil {
			first = err
		}
	}
	return first
}

// ============================================================
// GuestAllocator：自适应内存分配器
// ============================================================

// GuestAllocator 能够智能管理不同 Wasm Guest 的内存。
// 它会自动检测并适配多种内存管理接口，并提供并发安全的分配与释放。
//
// 工作流程:
// 1. NewGuestAllocator 检测模块导出的函数，确定分配模式
// 2. Allocate 按检测到的模式调用 Guest 分配函数
// 3. Free 按模式正确释放（cabi_realloc 直接释放，realloc/malloc 需还原原始指针）
type GuestAllocator struct {
	mu      sync.Mutex // 保证多线程并发分配与释放 Guest 内存时的线程安全
	mode    allocatorMode
	realloc api.Function // cabi_realloc 或 realloc 函数引用
	malloc  api.Function // malloc 函数引用（modeMallocFree 模式）
	free    api.Function // free 函数引用

	// aligned 用于 realloc/malloc 模式：对齐后的指针 -> 原始 malloc 指针，保证 Free 能还回真正基址。
	// 例如: aligned[0x1004] = {raw: 0x1000, size: 12, align: 8}
	aligned map[uint32]rawAlloc
}

// NewGuestAllocator 通过检查模块导出的函数，创建一个自适应的分配器。
// 按 `cabi_realloc` -> `realloc` -> `malloc`/`free` 的优先级回退。
//
// 参数:
//   - module: Wasm 模块实例，必须非空且至少导出一个可识别的内存管理函数
//
// 返回:
//   - *GuestAllocator: 创建的分配器实例
//   - error: 模块无效或无可识别函数时的错误
//
// 示例:
//
//	alloc, err := NewGuestAllocator(module)
//	if err != nil { panic(err) }
//	ptr, err := alloc.Allocate(ctx, 1024, 8)
func NewGuestAllocator(module api.Module) (*GuestAllocator, error) {
	if module == nil {
		return nil, fmt.Errorf("module 不能为空")
	}

	// 优先级 1: 标准组件模型 `cabi_realloc`
	// 签名: cabi_realloc(ptr, old_size, align, size) -> new_ptr
	if reallocFunc := module.ExportedFunction("cabi_realloc"); reallocFunc != nil {
		if err := validateExportedSig(reallocFunc, "cabi_realloc", 4, 1); err != nil {
			return nil, err
		}
		return &GuestAllocator{
			mode:    modeCabiRealloc,
			realloc: reallocFunc,
		}, nil
	}

	// 优先级 2: TinyGo 或 C-style `realloc`
	// 签名: realloc(ptr, size) -> new_ptr
	if reallocFunc := module.ExportedFunction("realloc"); reallocFunc != nil {
		if err := validateExportedSig(reallocFunc, "realloc", 2, 1); err != nil {
			return nil, err
		}
		freeFunc := module.ExportedFunction("free")
		if freeFunc != nil {
			if err := validateExportedSig(freeFunc, "free", 1, 0); err != nil {
				return nil, err
			}
		}
		return &GuestAllocator{
			mode:    modeRealloc,
			realloc: reallocFunc,
			free:    freeFunc, // freeFunc 为 nil 时 Free() 对 Guest 为空操作（TinyGo GC）
			aligned: make(map[uint32]rawAlloc),
		}, nil
	}

	// 优先级 3: C-style `malloc`/`free`
	// 签名: malloc(size) -> ptr, free(ptr)
	if mallocFunc := module.ExportedFunction("malloc"); mallocFunc != nil {
		if err := validateExportedSig(mallocFunc, "malloc", 1, 1); err != nil {
			return nil, err
		}
		if freeFunc := module.ExportedFunction("free"); freeFunc != nil {
			if err := validateExportedSig(freeFunc, "free", 1, 0); err != nil {
				return nil, err
			}
			return &GuestAllocator{
				mode:    modeMallocFree,
				malloc:  mallocFunc,
				free:    freeFunc,
				aligned: make(map[uint32]rawAlloc),
			}, nil
		}
	}

	return nil, fmt.Errorf("模块未导出任何可识别的内存管理函数 (cabi_realloc, realloc, or malloc/free)")
}

// validateExportedSig 校验导出函数的参数/返回个数；签名不可用时放行，避免误杀。
// 某些 Wasm 模块不暴露完整签名信息，此时跳过严格校验。
//
// 参数:
//   - fn:          要校验的导出函数
//   - name:        函数名称（用于错误信息）
//   - wantParams:  期望的参数个数
//   - minResults:  最小的返回值个数
//
// 返回:
//   - error: 参数/返回值不匹配时的错误，或 nil
func validateExportedSig(fn api.Function, name string, wantParams, minResults int) error {
	if fn == nil {
		return fmt.Errorf("exported function %s is nil", name)
	}
	def := fn.Definition()
	if def == nil {
		return nil // 无签名信息，放行
	}
	pt, rt := def.ParamTypes(), def.ResultTypes()
	if len(pt) == 0 && len(rt) == 0 {
		return nil // 空签名，放行
	}
	if len(pt) != wantParams {
		return fmt.Errorf("%s 参数个数不匹配: got %d want %d", name, len(pt), wantParams)
	}
	if len(rt) < minResults {
		return fmt.Errorf("%s 返回值个数不足: got %d want >= %d", name, len(rt), minResults)
	}
	return nil
}

// normalizeAlign 将 alignment 规范化为 2 的幂，非 2 幂则上取整到最近的 2 的幂，符合 Canonical ABI。
// 例如: 3 -> 4, 5 -> 8, 8 -> 8
//
// 参数:
//   - alignment: 原始对齐值
//
// 返回:
//   - uint32: 规范化后的 2 的幂对齐值
func normalizeAlign(alignment uint32) uint32 {
	if alignment <= 1 {
		return 1
	}
	if alignment&(alignment-1) == 0 {
		return alignment // 已是 2 的幂
	}
	// 上取整到下一个 2 的幂（位运算技巧）
	a := alignment - 1
	a |= a >> 1
	a |= a >> 2
	a |= a >> 4
	a |= a >> 8
	a |= a >> 16
	a++
	if a == 0 {
		// 当 alignment 接近 2^32 时 +1 溢出成 0，cabi_realloc 对齐 0 未定义
		return 1 << 31
	}
	return a
}

// callU32 调用给定的 api.Function，将参数转换为 uint64 切片，并返回第一个结果（uint32）。
// 封装了统一的函数调用、错误处理和类型转换逻辑。
//
// 参数:
//   - ctx:   上下文
//   - fn:    要调用的 Wasm 函数
//   - name:  函数名称（用于错误信息）
//   - args:  传递给 Wasm 函数的 uint64 参数列表
//
// 返回:
//   - uint32: 函数的第一个返回值的低 32 位
//   - error:  调用失败时的错误
//
// 示例:
//
//	ptr, err := alloc.callU32(ctx, reallocFn, "realloc", 0, 1024)
func (a *GuestAllocator) callU32(ctx context.Context, fn api.Function, name string, args ...uint64) (uint32, error) {
	if fn == nil {
		return 0, fmt.Errorf("%s function is nil", name)
	}
	results, err := fn.Call(ctx, args...)
	if err != nil {
		return 0, fmt.Errorf("%s failed: %w", name, err)
	}
	if len(results) == 0 {
		return 0, fmt.Errorf("%s returned no results", name)
	}
	return uint32(results[0]), nil
}

// Allocate 在 Guest 中分配一块内存，自动选择正确的分配方式。
// 分配记录会自动添加到当前 context 的作用域，以便调用结束后自动释放。
func (a *GuestAllocator) Allocate(ctx context.Context, size, alignment uint32) (uint32, error) {
	if a == nil {
		return 0, fmt.Errorf("GuestAllocator 实例不能为 nil")
	}
	if size == 0 {
		return 0, nil // 零尺寸分配返回空指针
	}
	if ctx == nil {
		ctx = context.Background()
	}
	if err := ctx.Err(); err != nil {
		return 0, err
	}

	alignment = normalizeAlign(alignment)

	var (
		ptr uint32
		err error
	)

	// 原先整段 Allocate（含 cabi_realloc/malloc 的 fn.Call）持 a.mu。
	// Host.Call 进入 guest 后，guest 再调 host export 并 Allocate，会在同一把锁上自死锁。
	// mode 与函数指针在 NewGuestAllocator 之后只读；只有 aligned 映射需要锁。
	switch a.mode {
	case modeCabiRealloc:
		ptr, err = a.callU32(ctx, a.realloc, "cabi_realloc for allocate", 0, 0, uint64(alignment), uint64(size))
	case modeRealloc:
		ptr, err = a.allocAligned(ctx, size, alignment, true)
	case modeMallocFree:
		ptr, err = a.allocAligned(ctx, size, alignment, false)
	default:
		return 0, fmt.Errorf("不支持的分配器模式")
	}
	if err != nil {
		return 0, err
	}
	if ptr == 0 {
		return 0, fmt.Errorf("guest allocator returned NULL for size=%d align=%d", size, alignment)
	}

	scopeFromCtx(ctx).add(allocRec{ptr: ptr, size: size, align: alignment})
	return ptr, nil
}

// allocAligned 在 realloc/malloc 模式下过量分配并返回对齐指针。
func (a *GuestAllocator) allocAligned(ctx context.Context, size, alignment uint32, useRealloc bool) (uint32, error) {
	need := uint64(size) + uint64(alignment) - 1
	if need > uint64(^uint32(0)) {
		return 0, fmt.Errorf("aligned allocation size overflow: size=%d align=%d", size, alignment)
	}
	rawSize := uint32(need)

	var (
		raw uint32
		err error
	)
	if useRealloc {
		raw, err = a.callU32(ctx, a.realloc, "realloc for allocate", 0, uint64(rawSize))
	} else {
		raw, err = a.callU32(ctx, a.malloc, "malloc", uint64(rawSize))
	}
	if err != nil {
		return 0, err
	}
	if raw == 0 {
		return 0, fmt.Errorf("guest allocator returned NULL")
	}
	aligned := align(raw, alignment)
	a.mu.Lock()
	if a.aligned == nil {
		a.aligned = make(map[uint32]rawAlloc)
	}
	a.aligned[aligned] = rawAlloc{raw: raw, size: rawSize, align: alignment}
	a.mu.Unlock()
	return aligned, nil
}

// Free 在 Guest 中释放一块内存，自动选择正确的释放方式。
func (a *GuestAllocator) Free(ctx context.Context, ptr, size, alignment uint32) error {
	if a == nil {
		return fmt.Errorf("GuestAllocator 实例不能为 nil")
	}
	if ptr == 0 {
		return nil
	}
	if ctx == nil {
		ctx = context.Background()
	}

	alignment = normalizeAlign(alignment)
	scopeFromCtx(ctx).remove(ptr)

	a.mu.Lock()
	mode := a.mode
	realloc := a.realloc
	freeFn := a.free
	freePtr := ptr
	if rec, ok := a.aligned[ptr]; ok {
		freePtr = rec.raw
		delete(a.aligned, ptr)
	}
	a.mu.Unlock()

	// 与 Allocate 对称，fn.Call 必须在锁外，避免嵌套 host 分配死锁。
	switch mode {
	case modeCabiRealloc:
		_, err := a.callU32(ctx, realloc, "cabi_realloc for free", uint64(ptr), uint64(size), uint64(alignment), 0)
		return err

	case modeRealloc:
		if freeFn != nil {
			if _, err := freeFn.Call(ctx, uint64(freePtr)); err != nil {
				return fmt.Errorf("free (paired with realloc) failed: %w", err)
			}
		}
		return nil

	case modeMallocFree:
		if freeFn == nil {
			return fmt.Errorf("free function is nil")
		}
		if _, err := freeFn.Call(ctx, uint64(freePtr)); err != nil {
			return fmt.Errorf("free failed: %w", err)
		}
		return nil

	default:
		return fmt.Errorf("不支持的分配器模式")
	}
}
