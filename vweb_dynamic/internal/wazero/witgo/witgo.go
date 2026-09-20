package witgo

import (
	"context"
	"errors"
	"fmt"
	"math"
	"reflect"
	"sync"

	"github.com/tetratelabs/wazero/api"
)

// ============================================================
// Host：Wasm 组件交互的高层接口
// ============================================================

// Host 提供与 Wasm 组件交互的高层接口。
// 封装了内存分配器和模块引用，统一管理 Go <-> Wasm 的数据转换。
//
// 每个 Host 实例绑定到一个 Wasm 模块，持有:
//   - module: Wasm 模块引用（用于调用导出函数）
//   - allocator: Guest 内存分配器（用于分配 Guest 内存）
type Host struct {
	module    api.Module      // Wasm 模块
	allocator *GuestAllocator // Guest 内存分配器
}

var hostCache sync.Map // api.Module -> *Host，复用 Host 实例避免重复初始化

// getOrCreateHost 从缓存获取或新建 Host。
// 同一模块的 Host 实例会被缓存复用，避免重复创建分配器。
//
// 参数:
//   - module: Wasm 模块
//
// 返回:
//   - *Host: 对应的 Host 实例
//   - error: 创建失败时的错误
func getOrCreateHost(module api.Module) (*Host, error) {
	if module == nil {
		return nil, fmt.Errorf("module cannot be nil")
	}
	if v, ok := hostCache.Load(module); ok {
		if h, ok := v.(*Host); ok && h != nil {
			return h, nil
		}
	}
	h, err := NewHost(module)
	if err != nil {
		return nil, err
	}
	actual, _ := hostCache.LoadOrStore(module, h)
	if cached, ok := actual.(*Host); ok && cached != nil {
		return cached, nil
	}
	return h, nil
}

// NewHost 基于模块创建 Host，并初始化 GuestAllocator。
//
// 参数:
//   - module: Wasm 模块
//
// 返回:
//   - *Host: 创建的 Host 实例
//   - error: 分配器初始化失败时的错误
func NewHost(module api.Module) (*Host, error) {
	if module == nil {
		return nil, fmt.Errorf("module cannot be nil")
	}
	alloc, err := NewGuestAllocator(module)
	if err != nil {
		return nil, err
	}
	return &Host{module: module, allocator: alloc}, nil
}

// ErrNotExportFunc 表示 Guest 未导出指定函数。
var ErrNotExportFunc = errors.New("guest function not exports")

// setScalarFromU64 根据 outVal 的类型将 uint64 值转换并设置。
// 用于 decodeResult 中将 Wasm 返回值直接写入标量类型。
func setScalarFromU64(outVal reflect.Value, resultValue uint64) {
	switch outVal.Kind() {
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		outVal.SetUint(resultValue)
	case reflect.Int8:
		outVal.SetInt(int64(int8(resultValue)))
	case reflect.Int16:
		outVal.SetInt(int64(int16(resultValue)))
	case reflect.Int32:
		outVal.SetInt(int64(int32(resultValue)))
	case reflect.Int64:
		outVal.SetInt(int64(resultValue))
	case reflect.Int:
		if outVal.Type().Size() == 4 {
			outVal.SetInt(int64(int32(resultValue)))
		} else {
			outVal.SetInt(int64(resultValue))
		}
	case reflect.Uint:
		if outVal.Type().Size() == 4 {
			outVal.SetUint(uint64(uint32(resultValue)))
		} else {
			outVal.SetUint(resultValue)
		}
	case reflect.Float32:
		outVal.SetFloat(float64(math.Float32frombits(uint32(resultValue))))
	case reflect.Float64:
		outVal.SetFloat(math.Float64frombits(resultValue))
	case reflect.Bool:
		outVal.SetBool(resultValue != 0)
	}
}

// liftParamsAsTuple 把一组 Go 参数按记录（tuple）布局写入 guest 内存，返回起始指针。
// Canonical ABI 在扁平参数超过 MAX_FLAT_PARAMS 时，核心函数只接收这个指针。
func (h *Host) liftParamsAsTuple(ctx context.Context, params []interface{}) (uint32, error) {
	type item struct {
		val    reflect.Value
		layout *TypeLayout
	}
	items := make([]item, 0, len(params))
	var totalSize, maxAlign uint32 = 0, 1
	for _, p := range params {
		val := reflect.ValueOf(p)
		for val.Kind() == reflect.Pointer {
			if !val.IsValid() || val.IsNil() {
				return 0, fmt.Errorf("cannot lift nil pointer of type %v", val.Type())
			}
			val = val.Elem()
		}
		if !val.IsValid() {
			return 0, fmt.Errorf("invalid param value")
		}
		layout, err := GetOrCalculateLayout(val.Type())
		if err != nil {
			return 0, err
		}
		if layout.Alignment > maxAlign {
			maxAlign = layout.Alignment
		}
		totalSize = align(totalSize, layout.Alignment)
		next, err := addU32(totalSize, layout.Size)
		if err != nil {
			return 0, err
		}
		totalSize = next
		items = append(items, item{val: val, layout: layout})
	}
	if maxAlign == 0 {
		maxAlign = 1
	}
	totalSize = align(totalSize, maxAlign)
	ptr, err := h.allocator.Allocate(ctx, totalSize, maxAlign)
	if err != nil {
		return 0, err
	}
	var offset uint32
	for _, it := range items {
		offset = align(offset, it.layout.Alignment)
		fieldPtr, err := addU32(ptr, offset)
		if err != nil {
			return 0, err
		}
		if err := LiftToPtr(ctx, h.module.Memory(), h.allocator, it.val, fieldPtr); err != nil {
			return 0, err
		}
		offset, err = addU32(offset, it.layout.Size)
		if err != nil {
			return 0, err
		}
	}
	return ptr, nil
}

// Call 调用 Guest 导出的函数 funcName，传入 params，并将结果写入 resultPtr（可选）。
// params 可以是任意 Go 值，会自动扁平化为 Guest 参数。
// 分配内存在返回前自动释放（通过 context scope）。
//
// 参数:
//   - ctx:       上下文（用于取消和分配作用域）
//   - funcName:  Guest 导出函数名
//   - resultPtr: 接收返回值的指针（可为 nil 表示无返回值）
//   - params:    可变参数列表（Go 值）
//
// 返回:
//   - error: 调用失败时的错误
//
// 示例:
//
//	var result uint32
//	err := host.Call(ctx, "add", &result, 1, 2)
//
//	var greeting string
//	err := host.Call(ctx, "greet", &greeting, "World")
func (h *Host) Call(ctx context.Context, funcName string, resultPtr interface{}, params ...interface{}) error {
	if h == nil || h.module == nil || h.allocator == nil {
		return fmt.Errorf("host instance or module is not initialized")
	}
	if ctx == nil {
		ctx = context.Background()
	}

	// 获取 Guest 函数
	fn := h.module.ExportedFunction(funcName)
	if fn == nil {
		return fmt.Errorf("函数 '%s' 在 Guest 导出中未找到 %w", funcName, ErrNotExportFunc)
	}

	// 创建分配作用域，确保内存自动释放
	ctx, scope := contextWithAllocScope(ctx)
	defer func() { _ = scope.freeAll(ctx, h.allocator) }()

	def := fn.Definition()
	if def == nil {
		// 部分 runtime 可能不提供签名；解引用会 panic
		return fmt.Errorf("函数 '%s' 缺少 signature definition", funcName)
	}
	paramDefs := def.ParamTypes()
	resultDefs := def.ResultTypes()

	useRetArea := resultPtr != nil && len(resultDefs) == 0 && len(paramDefs) > 0 &&
		paramDefs[len(paramDefs)-1] == api.ValueTypeI32
	argSlots := len(paramDefs)
	if useRetArea {
		argSlots--
	}

	flatCount, err := flattenParamsCount(params)
	if err != nil {
		return fmt.Errorf("扁平化参数失败: %w", err)
	}
	// 间接 tuple 路径若先 flattenParam 再 LiftToPtr，string/list 会在 guest 里分配两份。
	needPacked := argSlots == 1 && len(paramDefs) > 0 && paramDefs[0] == api.ValueTypeI32 && flatCount != 1

	var flatParams []uint64
	if needPacked {
		packed, perr := h.liftParamsAsTuple(ctx, params)
		if perr != nil {
			return fmt.Errorf("为函数 '%s' 打包间接参数失败: %w", funcName, perr)
		}
		flatParams = []uint64{uint64(packed)}
	} else {
		flatParams = make([]uint64, 0, 16)
		for _, p := range params {
			if ferr := h.flattenParam(ctx, reflect.ValueOf(p), &flatParams); ferr != nil {
				return fmt.Errorf("扁平化参数 %#v 失败: %w", p, ferr)
			}
		}
	}

	var retArea uint32
	if useRetArea {
		rv := reflect.ValueOf(resultPtr)
		if rv.Kind() != reflect.Pointer || rv.IsNil() {
			return fmt.Errorf("resultPtr 必须是非空的指针类型，得到: %T", resultPtr)
		}
		layout, lerr := GetOrCalculateLayout(rv.Elem().Type())
		if lerr != nil {
			return fmt.Errorf("计算返回值布局失败: %w", lerr)
		}
		retArea, err = h.allocator.Allocate(ctx, layout.Size, layout.Alignment)
		if err != nil {
			return fmt.Errorf("分配 retptr 失败: %w", err)
		}
		flatParams = append(flatParams, uint64(retArea))
	}

	// 校验参数数量
	if want := len(paramDefs); want != 0 && len(flatParams) != want {
		return fmt.Errorf("函数 '%s' 参数个数不匹配: flattened %d, guest wants %d", funcName, len(flatParams), want)
	}

	// 调用 Guest 函数
	results, err := fn.Call(ctx, flatParams...)
	if err != nil {
		return fmt.Errorf("Guest 函数 '%s' 调用失败: %w", funcName, err)
	}

	// 解码返回值
	if resultPtr != nil {
		rv := reflect.ValueOf(resultPtr)
		if rv.Kind() != reflect.Pointer || rv.IsNil() {
			return fmt.Errorf("resultPtr 必须是非空的指针类型，得到: %T", resultPtr)
		}
		if useRetArea {
			// 无核心返回值，从 retptr 指向的内存 Lower 到 Go
			if err := Lower(ctx, h, retArea, rv.Elem()); err != nil {
				return fmt.Errorf("从 retptr 降低返回值失败: %w", err)
			}
			return nil
		}
		if len(results) == 0 {
			if len(resultDefs) > 0 {
				return fmt.Errorf("函数期望有返回值，但实际没有返回")
			}
			return nil
		}
		if err := h.decodeResult(ctx, results, rv.Elem()); err != nil {
			return err
		}
	}
	return nil
}

// decodeResult 将 Guest 返回的 uint64 结果解码到 outVal。
// 根据 outVal 的类型选择合适的解码策略。
func (h *Host) decodeResult(ctx context.Context, results []uint64, outVal reflect.Value) error {
	if !outVal.CanSet() {
		return fmt.Errorf("result value is not settable")
	}
	// 标量类型直接赋值
	if isScalarKind(outVal.Kind()) && len(results) == 1 {
		setScalarFromU64(outVal, results[0])
		return nil
	}

	// 检查是否可以用扁平化类型匹配
	flatWant, ferr := flattenType(outVal.Type())
	if ferr == nil && len(flatWant) == len(results) {
		kind := outVal.Kind()
		// 字符串/切片特殊情况：单个指针结果
		if (kind == reflect.String || kind == reflect.Slice) && len(results) == 1 {
			return Lower(ctx, h, uint32(results[0]), outVal)
		}
		// 通过 paramStream 反扁平化
		ps := &paramStream{params: results}
		val, err := h.unflattenParam(ctx, h.module.Memory(), ps, outVal.Type())
		if err != nil {
			return fmt.Errorf("failed to unflatten result: %w", err)
		}
		outVal.Set(val)
		return nil
	}

	// 回退到指针模式（单个指针结果）
	return Lower(ctx, h, uint32(results[0]), outVal)
}

func flattenParamsCount(params []any) (int, error) {
	n := 0
	for _, p := range params {
		val := reflect.ValueOf(p)
		for val.Kind() == reflect.Pointer {
			if !val.IsValid() || val.IsNil() {
				return 0, fmt.Errorf("cannot flatten nil pointer of type %v", val.Type())
			}
			val = val.Elem()
		}
		if !val.IsValid() {
			return 0, fmt.Errorf("invalid param value")
		}
		ft, err := flattenType(val.Type())
		if err != nil {
			return 0, err
		}
		n += len(ft)
	}
	return n, nil
}

// ReleaseHost 从模块缓存去掉 Host，避免 wasm 实例销毁后 allocator.aligned 永久泄漏。
func ReleaseHost(module api.Module) {
	if module != nil {
		hostCache.Delete(module)
	}
}
