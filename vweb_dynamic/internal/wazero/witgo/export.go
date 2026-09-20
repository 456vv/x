package witgo

import (
	"context"
	"fmt"
	"math"
	"reflect"
	"sync"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
)

// Canonical ABI（组件模型）限制：超过后改为“内存指针”传递。
// 原先结果用 >2 才走 retptr，与 CABI MAX_FLAT_RESULTS=1 不符；
// 参数超过 16 个扁平值时也必须改成单一 i32 指针（MAX_FLAT_PARAMS=16）。
const (
	maxFlatParams  = 16
	maxFlatResults = 1
)

// exportCacheKey 用于缓存导出函数的 wrapper。
// 组合类型和函数指针地址，确保不同实例的闭包不被错误复用。
type exportCacheKey struct {
	typ reflect.Type // 函数类型
	ptr uintptr      // 函数指针地址
}

// Exporter 提供链式 API，用于将 Go 函数导出为 Guest 可调用的函数。
//
// 使用示例:
//
//	exporter := witgo.NewExporter(builder)
//	exporter.Export("greet", func(ctx context.Context, name string) string {
//	    return "Hello, " + name
//	}).Export("add", func(a, b uint32) uint32 { return a + b })
type Exporter struct {
	wazero.HostModuleBuilder
	mu sync.Mutex // 保护 NewFunctionBuilder 的并发调用

	// wrapperCache 保留用于兼容，但不再按代码指针复用闭包。
	wrapperCache map[exportCacheKey]any
}

var exportSigCache sync.Map // reflect.Type -> *exportSig，缓存签名分析结果而非闭包

// exportSig 缓存一个 Go 函数签名的扁平化分析结果。
type exportSig struct {
	flatIn, flatOut   []reflect.Type // 扁平化后的参数/返回值类型（已按 CABI 收缩）
	hasRetptr         bool           // 是否使用间接返回（retptr）
	hasIndirectParams bool           // 是否使用间接参数（单一 i32 指向参数 tuple）
	wrapperType       reflect.Type   // wrapper 函数的 reflect.Type
}

// NewExporter 创建一个 Exporter，基于给定的 HostModuleBuilder。
//
// 参数:
//   - builder: wazero 的 HostModuleBuilder，用于注册导出的函数
//
// 返回:
//   - *Exporter: 初始化的导出器
func NewExporter(builder wazero.HostModuleBuilder) *Exporter {
	return &Exporter{
		HostModuleBuilder: builder,
		wrapperCache:      make(map[exportCacheKey]any),
	}
}

// MustExport 调用 Export，若失败则 panic。
// 适用于启动时确定性的导出，不需要处理错误。
//
// 参数:
//   - funcName: 导出的函数名称（Guest 可见）
//   - goFunc:   要导出的 Go 函数
//
// 返回:
//   - *Exporter: 自身，支持链式调用
//
// 示例:
//
//	exporter.MustExport("my_func", myGoFunction)
func (e *Exporter) MustExport(funcName string, goFunc any) *Exporter {
	if err := e.Export(funcName, goFunc); err != nil {
		panic(err)
	}
	return e
}

// Export 将 Go 函数 goFunc 以 funcName 导出为 Guest 函数。
// goFunc 必须是一个函数，其参数和返回值类型将被自动扁平化，并适配 ABI。
// 若 goFunc 包含 context.Context 作为第一个参数，则会传入调用上下文。
//
// 参数:
//   - funcName: 导出的函数名称
//   - goFunc:   要导出的 Go 函数（必须是 func 类型）
//
// 返回:
//   - error: 函数类型不支持或导出失败时的错误
//
// 示例:
//
//	err := exporter.Export("add", func(a, b uint32) uint32 { return a + b })
func (e *Exporter) Export(funcName string, goFunc interface{}) error {
	if e == nil {
		return fmt.Errorf("Exporter is nil")
	}
	if e.HostModuleBuilder == nil {
		return fmt.Errorf("HostModuleBuilder is nil")
	}
	if goFunc == nil {
		return fmt.Errorf("`goFunc` cannot be nil for export '%s'", funcName)
	}

	funcVal := reflect.ValueOf(goFunc)
	funcType := funcVal.Type()
	if funcType.Kind() != reflect.Func {
		return fmt.Errorf("`goFunc` must be a function, but got %T", goFunc)
	}

	wrapperFunc, err := e.makeWrapperFunc(funcType, funcVal)
	if err != nil {
		return fmt.Errorf("failed to create wrapper for %s: %w", funcName, err)
	}

	e.mu.Lock()
	defer e.mu.Unlock()
	e.NewFunctionBuilder().WithFunc(wrapperFunc).Export(funcName)
	return nil
}

// makeWrapperFunc 构造一个 reflect.MakeFunc 包装器，处理参数扁平化、间接参数、间接返回等。
// 返回的 interface{} 可作为 wazero 的 Go 函数实现。
//
// 参数:
//   - funcType: Go 函数的 reflect.Type
//   - funcVal:  Go 函数的 reflect.Value
//
// 返回:
//   - interface{}: 包装后的函数实现
//   - error: 签名分析失败时的错误
func (e *Exporter) makeWrapperFunc(funcType reflect.Type, funcVal reflect.Value) (interface{}, error) {
	// 尝试从缓存获取已分析的签名
	var sig *exportSig
	if v, ok := exportSigCache.Load(funcType); ok {
		sig, _ = v.(*exportSig)
	}
	if sig == nil {
		flatIn, flatOut, hasRetptr, hasIndirectParams, err := e.flattenSignatureTypes(funcType)
		if err != nil {
			return nil, err
		}
		wrapperIn := append([]reflect.Type{
			typeOf[context.Context](),
			typeOf[api.Module](),
		}, flatIn...)
		sig = &exportSig{
			flatIn:            flatIn,
			flatOut:           flatOut,
			hasRetptr:         hasRetptr,
			hasIndirectParams: hasIndirectParams,
			wrapperType:       reflect.FuncOf(wrapperIn, flatOut, false),
		}
		exportSigCache.Store(funcType, sig)
	}

	flatOut := sig.flatOut
	hasRetptr := sig.hasRetptr
	hasIndirectParams := sig.hasIndirectParams
	wrapperType := sig.wrapperType

	// 构建 wrapper 实现
	wrapperImpl := func(args []reflect.Value) []reflect.Value {
		ctx := args[0].Interface().(context.Context)
		module := args[1].Interface().(api.Module)

		h, err := getOrCreateHost(module)
		if err != nil {
			panic(fmt.Sprintf("failed to get host for calling module: %v", err))
		}

		argIdx := 2
		var paramBase uint32
		if hasIndirectParams {
			if argIdx >= len(args) {
				panic("function expected an indirect-params pointer, but received no parameters")
			}
			paramBase = uint32(args[argIdx].Uint())
			argIdx++
		}

		var retptr uint32
		if hasRetptr {
			if argIdx >= len(args) {
				panic("function expected a return pointer, but received no parameters")
			}
			retptr = uint32(args[argIdx].Uint())
		}

		// 反扁平化参数为 Go 值
		callArgs := make([]reflect.Value, funcType.NumIn())
		funcParamIndex := 0
		if len(callArgs) > 0 && funcType.In(0) == typeOf[context.Context]() {
			callArgs[0] = args[0]
			funcParamIndex = 1
		}

		if hasIndirectParams {
			// 扁平参数 >16 时 guest 只传一个指向“参数 tuple”的 i32。
			// 按记录布局从该指针逐个 Lower，不能再走 unflattenParam 的扁平栈。
			offset := paramBase
			for ; funcParamIndex < len(callArgs); funcParamIndex++ {
				paramType := funcType.In(funcParamIndex)
				layout, lerr := GetOrCalculateLayout(paramType)
				if lerr != nil {
					panic(fmt.Sprintf("failed to layout indirect param %d for %s: %v", funcParamIndex, funcVal.Type().String(), lerr))
				}
				offset = align(offset, layout.Alignment)
				outVal := reflect.New(paramType).Elem()
				if lerr = Lower(ctx, h, offset, outVal); lerr != nil {
					panic(fmt.Sprintf("failed to lower indirect param %d for %s: %v", funcParamIndex, funcVal.Type().String(), lerr))
				}
				callArgs[funcParamIndex] = outVal
				next, nerr := addU32(offset, layout.Size)
				if nerr != nil {
					panic(nerr)
				}
				offset = next
			}
		} else {
			end := len(args)
			if hasRetptr {
				end--
			}
			ps := &paramStream{params: make([]uint64, 0, end-2)}
			for _, arg := range args[2:end] {
				ps.params = append(ps.params, wasmArgToU64(arg))
			}
			for ; funcParamIndex < len(callArgs); funcParamIndex++ {
				paramType := funcType.In(funcParamIndex)
				val, uerr := h.unflattenParam(ctx, module.Memory(), ps, paramType)
				if uerr != nil {
					panic(fmt.Sprintf("failed to unflatten parameter %d for %s: %v", funcParamIndex, funcVal.Type().String(), uerr))
				}
				callArgs[funcParamIndex] = val
			}
		}

		results := funcVal.Call(callArgs)

		if hasRetptr {
			if len(results) == 0 {
				panic("function was expected to return a value but did not")
			}
			offset := retptr
			for i, res := range results {
				layout, lerr := GetOrCalculateLayout(res.Type())
				if lerr != nil {
					panic(fmt.Sprintf("failed to layout result %d: %v", i, lerr))
				}
				offset = align(offset, layout.Alignment)
				if lerr = LiftToPtr(ctx, module.Memory(), h.allocator, res, offset); lerr != nil {
					panic(fmt.Sprintf("failed to lift result %d to retptr: %v", i, lerr))
				}
				next, nerr := addU32(offset, layout.Size)
				if nerr != nil {
					panic(nerr)
				}
				offset = next
			}
			return nil
		}
		// 处理直接返回值（扁平化为 uint64）
		if len(flatOut) > 0 {
			var flats []uint64
			for _, res := range results {
				if ferr := h.flattenParam(ctx, res, &flats); ferr != nil {
					panic(fmt.Sprintf("failed to flatten result: %v", ferr))
				}
			}
			if len(flats) != len(flatOut) {
				panic(fmt.Sprintf("result arity mismatch: got %d flat values, want %d", len(flats), len(flatOut)))
			}
			out := make([]reflect.Value, len(flatOut))
			for i := range flatOut {
				ret := reflect.New(flatOut[i]).Elem()
				setFlatGoValue(ret, flats[i])
				out[i] = ret
			}
			return out
		}
		return nil
	}

	return reflect.MakeFunc(wrapperType, wrapperImpl).Interface(), nil
}

// wasmArgToU64 将任意 reflect.Value 转换为 uint64，用于传递给 Guest 函数。
// 处理各种数值类型、浮点数的位表示转换，以及布尔值。
//
// 参数:
//   - arg: 源 reflect.Value
//
// 返回:
//   - uint64: 转换后的值
func wasmArgToU64(arg reflect.Value) uint64 {
	switch arg.Kind() {
	case reflect.Int8, reflect.Int16, reflect.Int32:
		return uint64(uint32(int32(arg.Int()))) // 有符号扩展到 uint32 再提升到 uint64
	case reflect.Int64:
		return uint64(arg.Int())
	case reflect.Int:
		if arg.Type().Size() == 4 {
			return uint64(uint32(int32(arg.Int())))
		}
		return uint64(arg.Int())
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return arg.Uint()
	case reflect.Uint:
		if arg.Type().Size() == 4 {
			return uint64(uint32(arg.Uint()))
		}
		return arg.Uint()
	case reflect.Float32:
		return uint64(math.Float32bits(float32(arg.Float()))) // 保留原始位模式
	case reflect.Float64:
		return math.Float64bits(arg.Float())
	case reflect.Bool:
		if arg.Bool() {
			return 1
		}
		return 0
	default:
		if arg.CanUint() {
			return arg.Uint()
		}
		if arg.CanInt() {
			return uint64(arg.Int())
		}
		return 0
	}
}

// setFlatGoValue 将 uint64 值设置到 reflect.Value 中，根据其种类转换。
// 与 wasmArgToU64 相反的操作。
//
// 参数:
//   - ret: 目标 reflect.Value（必须是可设置的标量类型）
//   - v:   源 uint64 值
func setFlatGoValue(ret reflect.Value, v uint64) {
	switch ret.Kind() {
	case reflect.Bool:
		ret.SetBool(v != 0)
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
		ret.SetInt(int64(v))
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
		ret.SetUint(v)
	case reflect.Float32:
		ret.SetFloat(float64(math.Float32frombits(uint32(v))))
	case reflect.Float64:
		ret.SetFloat(math.Float64frombits(v))
	}
}

// flattenSignatureTypes 分析函数签名，决定核心 wasm 的扁平参数/返回值，以及是否走间接参数、间接返回。
func (e *Exporter) flattenSignatureTypes(funcType reflect.Type) (inTypes, outTypes []reflect.Type, isIndirectReturn, isIndirectParams bool, err error) {
	for i := 0; i < funcType.NumIn(); i++ {
		typ := funcType.In(i)
		if i == 0 && typ == typeOf[context.Context]() {
			continue
		}
		flatTypes, ferr := flattenType(typ)
		if ferr != nil {
			return nil, nil, false, false, fmt.Errorf("could not get shape for input type %v: %w", typ, ferr)
		}
		inTypes = append(inTypes, flatTypes...)
	}

	var initialOutTypes []reflect.Type
	for i := 0; i < funcType.NumOut(); i++ {
		flatTypes, ferr := flattenType(funcType.Out(i))
		if ferr != nil {
			return nil, nil, false, false, fmt.Errorf("could not get shape for output type %v: %w", funcType.Out(i), ferr)
		}
		initialOutTypes = append(initialOutTypes, flatTypes...)
	}

	// Canonical ABI MAX_FLAT_PARAMS=16；超出则核心函数只收一个 i32（指向参数 tuple）
	if len(inTypes) > maxFlatParams {
		isIndirectParams = true
		inTypes = []reflect.Type{typeOf[uint32]()}
	}

	// Canonical ABI MAX_FLAT_RESULTS=1。原先 >2 才 retptr，
	// 导致 string / result<T,E> 等 2 个扁平值被当成两个核心返回值，与 wit-bindgen guest 对不齐。
	if len(initialOutTypes) > maxFlatResults {
		isIndirectReturn = true
		inTypes = append(inTypes, typeOf[uint32]())
		outTypes = []reflect.Type{}
	} else {
		outTypes = initialOutTypes
	}
	return
}

// flattenType 是 flattenType 的别名，供内部使用。
func (e *Exporter) flattenType(typ reflect.Type) ([]reflect.Type, error) {
	return flattenType(typ)
}
