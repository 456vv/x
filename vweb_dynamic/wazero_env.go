package vweb_dynamic

import (
	"context"
	"crypto/rand"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
)

const (
	// maxEnvGuestBytes 限制从 guest 线性内存读入的单次字符串/参数大小。
	maxEnvGuestBytes = 16 << 20 // 16 MiB
	// errWriteFailed 是 tf_last_error 写失败时的特殊返回值（i32 -1）。
	// 与“无错误(0)”区分；guest 应视为不可恢复。
	errWriteFailed = 0xFFFFFFFF
)

// envHandleTable Host 侧对象表。Guest 只持有不透明 uint64 句柄。
//
// 句柄生成：salt（table 级随机）+ 递增 seq 乘法混合（Knuth 常量）。
//
// scopes：按 api.Module 隔离句柄与函数名缓存。
//   - tf_lookup / tf_call / tf_root_* 创建的句柄登记到对应 scope；
//   - ExecuteContext 通过 releaseModule 在入口调用返回后逆序回收；
//   - fnCache 从全局改为 per-module，消除跨实例 drop 失效问题。
type envHandleTable struct {
	mu      sync.RWMutex
	salt    uint64
	seq     uint64
	objects map[uint64]any

	scopes sync.Map // map[api.Module]*moduleScope
	roots  sync.Map // map[api.Module]uint64
	errs   sync.Map // map[api.Module]string
}

type moduleScope struct {
	mu      sync.Mutex
	handles map[uint64]struct{}
	fnCache map[string]uint64
}

func newModuleScope() *moduleScope {
	return &moduleScope{
		handles: make(map[uint64]struct{}),
		fnCache: make(map[string]uint64),
	}
}

func newEnvHandleTable() *envHandleTable {
	var b [8]byte
	var salt uint64
	if _, err := rand.Read(b[:]); err == nil {
		salt = binary.LittleEndian.Uint64(b[:])
	} else {
		salt = uint64(time.Now().UnixNano()) ^ 0x9E3779B97F4A7C15
	}
	if salt == 0 {
		salt = 0xA5A5A5A5A5A5A5A5
	}
	return &envHandleTable{salt: salt, objects: make(map[uint64]any)}
}

// scope 获取/创建 module 对应的句柄集合。mod 为 nil 时返回 nil（不登记）。
func (t *envHandleTable) scope(mod api.Module) *moduleScope {
	if mod == nil {
		return nil
	}
	if v, ok := t.scopes.Load(mod); ok {
		return v.(*moduleScope)
	}
	s := newModuleScope()
	actual, _ := t.scopes.LoadOrStore(mod, s)
	return actual.(*moduleScope)
}

// track 把 handle 登记到 module scope。putFor 与 tf_call 的 pack 句柄都走这里。
func (t *envHandleTable) track(mod api.Module, handle uint64) {
	if handle == 0 || mod == nil {
		return
	}
	if sc := t.scope(mod); sc != nil {
		sc.mu.Lock()
		sc.handles[handle] = struct{}{}
		sc.mu.Unlock()
	}
}

// putFor 保存对象并登记到 module scope，供 releaseModule 兜底回收。
// 所有返回给 guest 的“对象/结果句柄”都通过这里创建。
func (t *envHandleTable) putFor(mod api.Module, v any) uint64 {
	h := t.put(v)
	t.track(mod, h)
	return h
}

func (t *envHandleTable) releaseModule(mod api.Module) {
	if mod == nil {
		return
	}
	if v, ok := t.scopes.LoadAndDelete(mod); ok {
		s := v.(*moduleScope)
		s.mu.Lock()
		handles := make([]uint64, 0, len(s.handles))
		for h := range s.handles {
			handles = append(handles, h)
		}
		s.handles = make(map[uint64]struct{})
		s.fnCache = make(map[string]uint64)
		s.mu.Unlock()

		if len(handles) > 0 {
			t.mu.Lock()
			for _, h := range handles {
				delete(t.objects, h)
			}
			t.mu.Unlock()
		}
	}
	t.roots.Delete(mod)
	t.errs.Delete(mod)
}

func (t *envHandleTable) put(v any) uint64 {
	t.mu.Lock()
	defer t.mu.Unlock()
	for {
		t.seq++
		id := t.salt ^ (t.seq * 0x9E3779B97F4A7C15)
		if id == 0 {
			continue
		}
		if _, exists := t.objects[id]; exists {
			continue
		}
		t.objects[id] = v
		return id
	}
}

func (t *envHandleTable) get(id uint64) (any, bool) {
	if id == 0 {
		return nil, false
	}
	t.mu.RLock()
	defer t.mu.RUnlock()
	v, ok := t.objects[id]
	return v, ok
}

func (t *envHandleTable) drop(id uint64) {
	if id == 0 {
		return
	}
	t.mu.Lock()
	delete(t.objects, id)
	t.mu.Unlock()
}

func (t *envHandleTable) clear() {
	t.mu.Lock()
	t.objects = make(map[uint64]any)
	t.mu.Unlock()

	t.scopes.Range(func(k, _ any) bool {
		t.scopes.Delete(k)
		return true
	})
	t.roots.Range(func(k, _ any) bool {
		t.roots.Delete(k)
		return true
	})
	t.errs.Range(func(k, _ any) bool {
		t.errs.Delete(k)
		return true
	})
}

// setErr 按 module 记录/清除最近一次错误。err == nil 表示清空。
func (t *envHandleTable) setErr(mod api.Module, err error) {
	if mod == nil {
		return
	}
	if err == nil {
		t.errs.Delete(mod)
		return
	}
	t.errs.Store(mod, err.Error())
}

func (t *envHandleTable) getErr(mod api.Module) string {
	if mod == nil {
		return ""
	}
	if v, ok := t.errs.Load(mod); ok {
		return v.(string)
	}
	return ""
}

func (t *envHandleTable) setRoot(mod api.Module, h uint64) {
	if mod == nil {
		return
	}
	t.roots.Store(mod, h)
}

func (t *envHandleTable) rootFor(mod api.Module) uint64 {
	if mod == nil {
		return 0
	}
	if v, ok := t.roots.Load(mod); ok {
		return v.(uint64)
	}
	return 0
}

func (t *envHandleTable) clearRoot(mod api.Module) {
	if mod == nil {
		return
	}
	t.roots.Delete(mod)
}

// ---------------------------------------------------------------------------
// guest 线性内存读写
// ---------------------------------------------------------------------------

func readGuestBytes(mem api.Memory, ptr, length uint32) ([]byte, error) {
	if mem == nil {
		return nil, fmt.Errorf("module memory is nil")
	}
	if length == 0 {
		return nil, nil
	}
	if length > maxEnvGuestBytes {
		return nil, fmt.Errorf("guest payload %d exceeds limit %d", length, maxEnvGuestBytes)
	}
	b, ok := mem.Read(ptr, length)
	if !ok {
		return nil, fmt.Errorf("guest memory read failed ptr=%d len=%d", ptr, length)
	}
	out := make([]byte, length)
	copy(out, b)
	return out, nil
}

func readGuestString(mem api.Memory, ptr, length uint32) (string, error) {
	b, err := readGuestBytes(mem, ptr, length)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func writeGuestBytes(mem api.Memory, ptr uint32, data []byte) error {
	if len(data) == 0 {
		return nil
	}
	if mem == nil {
		return fmt.Errorf("module memory is nil")
	}
	if !mem.Write(ptr, data) {
		return fmt.Errorf("guest memory write failed ptr=%d len=%d", ptr, len(data))
	}
	return nil
}

// writeBackToGuest 把 data 写回 guest。返回值约定：
//   - (need, nil) 且 need <= outCap  → 已成功写入 need 字节
//   - (need, nil) 且 need >  outCap  → 未写入（容量不足），guest 需按 need 扩容后重试
//   - (0,    err)                    → 写失败
func writeBackToGuest(mod api.Module, outPtr, outCap uint32, data []byte) (uint64, error) {
	need := uint32(len(data))
	if need > outCap {
		return uint64(need), nil
	}
	if err := writeGuestBytes(mod.Memory(), outPtr, data); err != nil {
		return 0, err
	}
	return uint64(need), nil
}

// ---------------------------------------------------------------------------
// JSON 参数 -> Go 形参
// ---------------------------------------------------------------------------

// convertJSONArg 把 JSON 反序列化出来的 any 转成目标形参类型。
//   - 数字统一从 float64 按目标 kind 收窄
//   - 结构体/映射/切片/数组走「再序列化 + 反序列化」
//   - 指针形参：分配 *Elem 再 Unmarshal，返回类型与 want 一致（避免 New(want) 得到 **T）
func convertJSONArg(a any, want reflect.Type) (reflect.Value, error) {
	if want.Kind() == reflect.Interface {
		if a == nil {
			return reflect.Zero(want), nil
		}
		return reflect.ValueOf(a), nil
	}
	if a == nil {
		return reflect.Zero(want), nil
	}

	rv := reflect.ValueOf(a)
	if rv.Type().AssignableTo(want) {
		return rv, nil
	}
	if rv.Type().ConvertibleTo(want) {
		return rv.Convert(want), nil
	}

	if f, ok := a.(float64); ok {
		switch want.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return reflect.ValueOf(int64(f)).Convert(want), nil
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return reflect.ValueOf(uint64(f)).Convert(want), nil
		case reflect.Float32, reflect.Float64:
			return reflect.ValueOf(f).Convert(want), nil
		case reflect.Bool:
			return reflect.ValueOf(f != 0).Convert(want), nil
		}
	}

	switch want.Kind() {
	case reflect.Pointer:
		raw, err := json.Marshal(a)
		if err != nil {
			return reflect.Value{}, fmt.Errorf("%T -> %v: %w", a, want, err)
		}
		out := reflect.New(want.Elem())
		if err := json.Unmarshal(raw, out.Interface()); err != nil {
			return reflect.Value{}, fmt.Errorf("%T -> %v: %w", a, want, err)
		}
		return out, nil

	case reflect.Struct, reflect.Map, reflect.Slice, reflect.Array:
		raw, err := json.Marshal(a)
		if err != nil {
			return reflect.Value{}, fmt.Errorf("%T -> %v: %w", a, want, err)
		}
		out := reflect.New(want)
		if err := json.Unmarshal(raw, out.Interface()); err != nil {
			return reflect.Value{}, fmt.Errorf("%T -> %v: %w", a, want, err)
		}
		return out.Elem(), nil
	}

	return reflect.Value{}, fmt.Errorf("cannot convert %T to %v", a, want)
}

// valueToDisplayString 默认走 fmt.Sprint，与 text/template 的渲染语义一致；
// Stringer / error / []byte 单独处理，避免 fmt.Sprint 走的是 reflect 内部表示。
func valueToDisplayString(v any) string {
	switch x := v.(type) {
	case nil:
		return ""
	case string:
		return x
	case []byte:
		return string(x)
	case fmt.Stringer:
		return x.String()
	case error:
		return x.Error()
	default:
		return fmt.Sprint(x)
	}
}

// unwrapResultPack 把结果句柄里的 []any 拆包：
//   - 0 元素 → 空串（tf_stringify）/ nil（tf_stringify_json 走 json.Marshal(nil) = "null"）
//   - 1 元素 → 该元素
//   - 多元素 → 原样
func unwrapResultPack(v any) any {
	if pack, ok := v.([]any); ok {
		switch len(pack) {
		case 0:
			return nil
		case 1:
			return pack[0]
		}
	}
	return v
}

// valueToJSON 把值序列化为 JSON 字节串，供 tf_stringify_json 使用。
// 与 valueToDisplayString 对齐：对 []any 结果句柄做拆包。
func valueToJSON(v any) ([]byte, bool) {
	b, err := json.Marshal(unwrapResultPack(v))
	if err != nil {
		return nil, false
	}
	return b, true
}

// isLikelyFlatABIValue 粗判是否适合直接交给 witgo.Host.Call 扁平化。
//
// 注意：这**不是** witgo.flattenParam 的等价物。witgo 对 struct/array 也可扁平化，
// 但模板点对象通常是复杂 struct，本函数把它们判为不可扁平，以便提前提示用户：
//   - 走 RootAsFirstArg=true + tf_root_handle；
//   - 或只传标量 / string / 标量元素 slice。
//
// nil 返回 true（调用方通常已用 `v != nil` 保护 flatIn[0]）。
func isLikelyFlatABIValue(v any) bool {
	if v == nil {
		return true
	}
	rv := reflect.ValueOf(v)
	for rv.Kind() == reflect.Pointer {
		if rv.IsNil() {
			return false
		}
		rv = rv.Elem()
	}
	switch rv.Kind() {
	case reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64,
		reflect.String:
		return true
	case reflect.Slice:
		switch rv.Type().Elem().Kind() {
		case reflect.Bool,
			reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
			reflect.Float32, reflect.Float64,
			reflect.String:
			return true
		default:
			return false
		}
	}
	return false
}

// ---------------------------------------------------------------------------
// 注册 env API
// ---------------------------------------------------------------------------

// registerEnvTemplateAPI 在 builder（通常为 "env"）上导出句柄式 TemplateFunc 桥。
// 关键约定见每个导出函数的注释。
func registerEnvTemplateAPI(builder wazero.HostModuleBuilder, table *envHandleTable) {
	const (
		i32 = api.ValueTypeI32
		i64 = api.ValueTypeI64
	)

	builder.NewFunctionBuilder().
		WithGoModuleFunction(api.GoModuleFunc(func(_ context.Context, _ api.Module, stack []uint64) {
			stack[0] = 1
		}), nil, []api.ValueType{i32}).
		Export("tf_protocol_version")

	// tf_lookup(name_ptr, name_len) -> handle
	// fnCache 按 module 隔离；命中时复检句柄是否仍存活。
	builder.NewFunctionBuilder().
		WithGoModuleFunction(api.GoModuleFunc(func(_ context.Context, mod api.Module, stack []uint64) {
			table.setErr(mod, nil)
			name, err := readGuestString(mod.Memory(), uint32(stack[0]), uint32(stack[1]))
			if err != nil {
				table.setErr(mod, err)
				stack[0] = 0
				return
			}
			if name == "" {
				table.setErr(mod, fmt.Errorf("empty function name"))
				stack[0] = 0
				return
			}

			sc := table.scope(mod)

			var cached uint64
			if sc != nil {
				sc.mu.Lock()
				cached = sc.fnCache[name]
				sc.mu.Unlock()
			}
			if cached != 0 {
				if _, alive := table.get(cached); alive {
					stack[0] = cached
					return
				}
				if sc != nil {
					sc.mu.Lock()
					if sc.fnCache[name] == cached {
						delete(sc.fnCache, name)
					}
					sc.mu.Unlock()
				}
			}

			fn, ok := TemplateFunc[name]
			if !ok || fn == nil {
				table.setErr(mod, fmt.Errorf("TemplateFunc %q not found", name))
				stack[0] = 0
				return
			}

			h := table.putFor(mod, fn)
			if sc != nil {
				sc.mu.Lock()
				sc.fnCache[name] = h
				sc.mu.Unlock()
			}
			stack[0] = h
		}), []api.ValueType{i32, i32}, []api.ValueType{i64}).
		Export("tf_lookup")

	// tf_call(fn_handle, args_ptr, args_len) -> result_handle
	// 成功时 result 恒为 []any（0 个返回值 → 空切片）。
	builder.NewFunctionBuilder().
		WithGoModuleFunction(api.GoModuleFunc(func(_ context.Context, mod api.Module, stack []uint64) {
			fnHandle := stack[0]
			ptr, length := uint32(stack[1]), uint32(stack[2])

			fail := func(format string, args ...any) {
				table.setErr(mod, fmt.Errorf(format, args...))
				stack[0] = 0
			}

			fnObj, ok := table.get(fnHandle)
			if !ok {
				fail("invalid function handle")
				return
			}
			fv := reflect.ValueOf(fnObj)
			if fv.Kind() != reflect.Func {
				fail("handle is not a func (got %T)", fnObj)
				return
			}

			raw, err := readGuestString(mod.Memory(), ptr, length)
			if err != nil {
				fail("read args: %v", err)
				return
			}
			var jsonArgs []any
			if length > 0 {
				if err := json.Unmarshal([]byte(raw), &jsonArgs); err != nil {
					fail("json args: %v", err)
					return
				}
			}

			ft := fv.Type()
			nFixed := ft.NumIn()
			if ft.IsVariadic() {
				nFixed--
			}
			if len(jsonArgs) < nFixed {
				fail("too few args: got %d, want >= %d", len(jsonArgs), nFixed)
				return
			}
			if !ft.IsVariadic() && len(jsonArgs) != ft.NumIn() {
				fail("arg count: got %d, want %d", len(jsonArgs), ft.NumIn())
				return
			}

			in := make([]reflect.Value, 0, len(jsonArgs))
			for i := 0; i < nFixed; i++ {
				v, cerr := convertJSONArg(jsonArgs[i], ft.In(i))
				if cerr != nil {
					fail("arg %d: %v", i, cerr)
					return
				}
				in = append(in, v)
			}
			if ft.IsVariadic() {
				elemT := ft.In(nFixed).Elem()
				tail := reflect.MakeSlice(ft.In(nFixed), 0, len(jsonArgs)-nFixed)
				for i := nFixed; i < len(jsonArgs); i++ {
					v, cerr := convertJSONArg(jsonArgs[i], elemT)
					if cerr != nil {
						fail("arg %d: %v", i, cerr)
						return
					}
					tail = reflect.Append(tail, v)
				}
				in = append(in, tail)
			}

			var out []reflect.Value
			var callErr error
			func() {
				defer func() {
					if rec := recover(); rec != nil {
						callErr = fmt.Errorf("template func panic: %v", rec)
					}
				}()
				out = fv.Call(in)
			}()
			if callErr != nil {
				fail("%v", callErr)
				return
			}

			pack := make([]any, len(out))
			for i := range out {
				if out[i].IsValid() && out[i].CanInterface() {
					pack[i] = out[i].Interface()
				}
			}
			table.setErr(mod, nil)
			stack[0] = table.putFor(mod, pack)
		}), []api.ValueType{i64, i32, i32}, []api.ValueType{i64}).
		Export("tf_call")

	// tf_drop(handle)
	builder.NewFunctionBuilder().
		WithGoModuleFunction(api.GoModuleFunc(func(_ context.Context, mod api.Module, stack []uint64) {
			h := stack[0]
			table.drop(h)
			if v, ok := table.scopes.Load(mod); ok {
				sc := v.(*moduleScope)
				sc.mu.Lock()
				delete(sc.handles, h)
				sc.mu.Unlock()
			}
		}), []api.ValueType{i64}, nil).
		Export("tf_drop")

	// tf_stringify(handle, out_ptr, out_cap) -> nbytes
	// 返回语义与 writeBackToGuest 一致。
	builder.NewFunctionBuilder().
		WithGoModuleFunction(api.GoModuleFunc(func(_ context.Context, mod api.Module, stack []uint64) {
			h, outPtr, outCap := stack[0], uint32(stack[1]), uint32(stack[2])
			v, ok := table.get(h)
			if !ok {
				table.setErr(mod, fmt.Errorf("invalid handle for stringify"))
				stack[0] = 0
				return
			}
			n, err := writeBackToGuest(mod, outPtr, outCap, []byte(valueToDisplayString(unwrapResultPack(v))))
			if err != nil {
				table.setErr(mod, err)
				stack[0] = 0
				return
			}
			table.setErr(mod, nil)
			stack[0] = n
		}), []api.ValueType{i64, i32, i32}, []api.ValueType{i32}).
		Export("tf_stringify")

	// tf_stringify_json(handle, out_ptr, out_cap) -> nbytes
	// ret == 0 一定是失败（JSON 序列化任何值至少 1 字节），查 tf_last_error。
	builder.NewFunctionBuilder().
		WithGoModuleFunction(api.GoModuleFunc(func(_ context.Context, mod api.Module, stack []uint64) {
			h, outPtr, outCap := stack[0], uint32(stack[1]), uint32(stack[2])
			v, ok := table.get(h)
			if !ok {
				table.setErr(mod, fmt.Errorf("invalid handle for stringify_json"))
				stack[0] = 0
				return
			}
			b, ok := valueToJSON(v)
			if !ok {
				table.setErr(mod, fmt.Errorf("value of type %T is not JSON-serializable", v))
				stack[0] = 0
				return
			}
			n, err := writeBackToGuest(mod, outPtr, outCap, b)
			if err != nil {
				table.setErr(mod, err)
				stack[0] = 0
				return
			}
			table.setErr(mod, nil)
			stack[0] = n
		}), []api.ValueType{i64, i32, i32}, []api.ValueType{i32}).
		Export("tf_stringify_json")

	// tf_last_error(out_ptr, out_cap) -> nbytes
	//   0                    → 无错误
	//   need <= out_cap      → 已写出 need 字节
	//   need >  out_cap      → 未写出，需扩容后重试
	//   0xFFFFFFFF (i32 -1)  → 写出失败（不可恢复）
	builder.NewFunctionBuilder().
		WithGoModuleFunction(api.GoModuleFunc(func(_ context.Context, mod api.Module, stack []uint64) {
			s := table.getErr(mod)
			outPtr, outCap := uint32(stack[0]), uint32(stack[1])
			n, err := writeBackToGuest(mod, outPtr, outCap, []byte(s))
			if err != nil {
				table.setErr(mod, err)
				stack[0] = errWriteFailed
				return
			}
			stack[0] = n
		}), []api.ValueType{i32, i32}, []api.ValueType{i32}).
		Export("tf_last_error")

	// tf_root_handle() -> handle
	builder.NewFunctionBuilder().
		WithGoModuleFunction(api.GoModuleFunc(func(_ context.Context, mod api.Module, stack []uint64) {
			stack[0] = table.rootFor(mod)
		}), nil, []api.ValueType{i64}).
		Export("tf_root_handle")

	// tf_root_call(root_handle, method_ptr, method_len, args_ptr, args_len) -> result_handle
	builder.NewFunctionBuilder().
		WithGoModuleFunction(api.GoModuleFunc(func(_ context.Context, mod api.Module, stack []uint64) {
			rootHandle := stack[0]
			method, err := readGuestString(mod.Memory(), uint32(stack[1]), uint32(stack[2]))
			if err != nil {
				table.setErr(mod, err)
				stack[0] = 0
				return
			}
			raw, err := readGuestBytes(mod.Memory(), uint32(stack[3]), uint32(stack[4]))
			if err != nil {
				table.setErr(mod, err)
				stack[0] = 0
				return
			}
			obj, ok := table.get(rootHandle)
			if !ok {
				table.setErr(mod, fmt.Errorf("invalid root handle"))
				stack[0] = 0
				return
			}
			var jsonArgs []any
			if len(raw) > 0 {
				if err := json.Unmarshal(raw, &jsonArgs); err != nil {
					table.setErr(mod, fmt.Errorf("json args: %w", err))
					stack[0] = 0
					return
				}
			}
			results, err := callMethodByReflect(obj, method, jsonArgs)
			if err != nil {
				table.setErr(mod, err)
				stack[0] = 0
				return
			}
			pack := make([]any, len(results))
			for i, r := range results {
				pack[i] = r
			}
			table.setErr(mod, nil)
			stack[0] = table.putFor(mod, pack)
		}), []api.ValueType{i64, i32, i32, i32, i32}, []api.ValueType{i64}).
		Export("tf_root_call")

	// tf_root_get(root_handle, key_ptr, key_len) -> value_handle
	builder.NewFunctionBuilder().
		WithGoModuleFunction(api.GoModuleFunc(func(_ context.Context, mod api.Module, stack []uint64) {
			rootHandle := stack[0]
			key, err := readGuestString(mod.Memory(), uint32(stack[1]), uint32(stack[2]))
			if err != nil {
				table.setErr(mod, err)
				stack[0] = 0
				return
			}
			obj, ok := table.get(rootHandle)
			if !ok {
				table.setErr(mod, fmt.Errorf("invalid root handle"))
				stack[0] = 0
				return
			}
			val, ok := getFieldByKey(obj, key)
			if !ok {
				table.setErr(mod, fmt.Errorf("key %q not found", key))
				stack[0] = 0
				return
			}
			table.setErr(mod, nil)
			stack[0] = table.putFor(mod, val)
		}), []api.ValueType{i64, i32, i32}, []api.ValueType{i64}).
		Export("tf_root_get")

	// tf_root_set(root_handle, key_ptr, key_len, val_ptr, val_len) -> i32
	builder.NewFunctionBuilder().
		WithGoModuleFunction(api.GoModuleFunc(func(_ context.Context, mod api.Module, stack []uint64) {
			rootHandle := stack[0]
			key, err := readGuestString(mod.Memory(), uint32(stack[1]), uint32(stack[2]))
			if err != nil {
				table.setErr(mod, err)
				stack[0] = errWriteFailed
				return
			}
			raw, err := readGuestBytes(mod.Memory(), uint32(stack[3]), uint32(stack[4]))
			if err != nil {
				table.setErr(mod, err)
				stack[0] = errWriteFailed
				return
			}
			obj, ok := table.get(rootHandle)
			if !ok {
				table.setErr(mod, fmt.Errorf("invalid root handle"))
				stack[0] = errWriteFailed
				return
			}
			var val any
			if err := json.Unmarshal(raw, &val); err != nil {
				table.setErr(mod, fmt.Errorf("json val: %w", err))
				stack[0] = errWriteFailed
				return
			}
			if err := setFieldByKey(obj, key, val); err != nil {
				table.setErr(mod, err)
				stack[0] = errWriteFailed
				return
			}
			table.setErr(mod, nil)
			stack[0] = 0
		}), []api.ValueType{i64, i32, i32, i32, i32}, []api.ValueType{i32}).
		Export("tf_root_set")

	// tf_root_has(root_handle, key_ptr, key_len) -> i32
	builder.NewFunctionBuilder().
		WithGoModuleFunction(api.GoModuleFunc(func(_ context.Context, mod api.Module, stack []uint64) {
			rootHandle := stack[0]
			key, err := readGuestString(mod.Memory(), uint32(stack[1]), uint32(stack[2]))
			if err != nil {
				table.setErr(mod, err)
				stack[0] = errWriteFailed
				return
			}
			obj, ok := table.get(rootHandle)
			if !ok {
				table.setErr(mod, fmt.Errorf("invalid root handle"))
				stack[0] = errWriteFailed
				return
			}
			if _, ok := getFieldByKey(obj, key); ok {
				stack[0] = 1
			} else {
				stack[0] = 0
			}
		}), []api.ValueType{i64, i32, i32}, []api.ValueType{i32}).
		Export("tf_root_has")

	// tf_root_len(root_handle, key_ptr, key_len) -> i32
	builder.NewFunctionBuilder().
		WithGoModuleFunction(api.GoModuleFunc(func(_ context.Context, mod api.Module, stack []uint64) {
			rootHandle := stack[0]
			key, err := readGuestString(mod.Memory(), uint32(stack[1]), uint32(stack[2]))
			if err != nil {
				table.setErr(mod, err)
				stack[0] = errWriteFailed
				return
			}
			obj, ok := table.get(rootHandle)
			if !ok {
				table.setErr(mod, fmt.Errorf("invalid root handle"))
				stack[0] = errWriteFailed
				return
			}
			target := obj
			if key != "" {
				v, ok := getFieldByKey(obj, key)
				if !ok {
					table.setErr(mod, fmt.Errorf("key %q not found", key))
					stack[0] = errWriteFailed
					return
				}
				target = v
			}
			v := reflect.ValueOf(target)
			for v.Kind() == reflect.Pointer || v.Kind() == reflect.Interface {
				if v.IsNil() {
					table.setErr(mod, fmt.Errorf("nil value"))
					stack[0] = errWriteFailed
					return
				}
				v = v.Elem()
			}
			switch v.Kind() {
			case reflect.Map, reflect.Slice, reflect.Array, reflect.String:
				table.setErr(mod, nil)
				stack[0] = uint64(v.Len())
			default:
				table.setErr(mod, fmt.Errorf("type %T has no length", target))
				stack[0] = errWriteFailed
			}
		}), []api.ValueType{i64, i32, i32}, []api.ValueType{i32}).
		Export("tf_root_len")

	// tf_root_get_path(root_handle, path_ptr, path_len) -> value_handle
	builder.NewFunctionBuilder().
		WithGoModuleFunction(api.GoModuleFunc(func(_ context.Context, mod api.Module, stack []uint64) {
			rootHandle := stack[0]
			path, err := readGuestString(mod.Memory(), uint32(stack[1]), uint32(stack[2]))
			if err != nil {
				table.setErr(mod, err)
				stack[0] = 0
				return
			}
			obj, ok := table.get(rootHandle)
			if !ok {
				table.setErr(mod, fmt.Errorf("invalid root handle"))
				stack[0] = 0
				return
			}
			val, err := getByPath(obj, path)
			if err != nil {
				table.setErr(mod, err)
				stack[0] = 0
				return
			}
			table.setErr(mod, nil)
			stack[0] = table.putFor(mod, val)
		}), []api.ValueType{i64, i32, i32}, []api.ValueType{i64}).
		Export("tf_root_get_path")

	// tf_root_set_path(root_handle, path_ptr, path_len, val_ptr, val_len) -> i32
	builder.NewFunctionBuilder().
		WithGoModuleFunction(api.GoModuleFunc(func(_ context.Context, mod api.Module, stack []uint64) {
			rootHandle := stack[0]
			path, err := readGuestString(mod.Memory(), uint32(stack[1]), uint32(stack[2]))
			if err != nil {
				table.setErr(mod, err)
				stack[0] = errWriteFailed
				return
			}
			raw, err := readGuestBytes(mod.Memory(), uint32(stack[3]), uint32(stack[4]))
			if err != nil {
				table.setErr(mod, err)
				stack[0] = errWriteFailed
				return
			}
			obj, ok := table.get(rootHandle)
			if !ok {
				table.setErr(mod, fmt.Errorf("invalid root handle"))
				stack[0] = errWriteFailed
				return
			}
			var val any
			if err := json.Unmarshal(raw, &val); err != nil {
				table.setErr(mod, fmt.Errorf("json val: %w", err))
				stack[0] = errWriteFailed
				return
			}
			if err := setByPath(obj, path, val); err != nil {
				table.setErr(mod, err)
				stack[0] = errWriteFailed
				return
			}
			table.setErr(mod, nil)
			stack[0] = 0
		}), []api.ValueType{i64, i32, i32, i32, i32}, []api.ValueType{i32}).
		Export("tf_root_set_path")
}

// getFieldByKey 从 map / slice / array / struct 取字段值。
// 指针/接口会自动解引用。返回 (值, true) 表示命中；(nil, false) 表示不存在。
func getFieldByKey(obj any, key string) (any, bool) {
	v := reflect.ValueOf(obj)
	for v.Kind() == reflect.Pointer || v.Kind() == reflect.Interface {
		if v.IsNil() {
			return nil, false
		}
		v = v.Elem()
	}
	switch v.Kind() {
	case reflect.Map:
		kv, err := convertKey(key, v.Type().Key())
		if err != nil {
			return nil, false
		}
		mv := v.MapIndex(kv)
		if !mv.IsValid() {
			return nil, false
		}
		return mv.Interface(), true
	case reflect.Slice, reflect.Array:
		i, err := strconv.Atoi(key)
		if err != nil || i < 0 || i >= v.Len() {
			return nil, false
		}
		return v.Index(i).Interface(), true
	case reflect.Struct:
		f := v.FieldByName(key)
		if !f.IsValid() {
			return nil, false
		}
		return f.Interface(), true
	}
	return nil, false
}

// setFieldByKey 修改 map / slice 元素或 struct 可导出字段。
func setFieldByKey(obj any, key string, val any) error {
	v := reflect.ValueOf(obj)
	for v.Kind() == reflect.Pointer || v.Kind() == reflect.Interface {
		if v.IsNil() {
			return fmt.Errorf("nil container")
		}
		v = v.Elem()
	}
	coerce := func(target reflect.Type) (reflect.Value, error) {
		vv := reflect.ValueOf(val)
		if !vv.IsValid() {
			return reflect.Zero(target), nil
		}
		if vv.Type().AssignableTo(target) {
			return vv, nil
		}
		if vv.Type().ConvertibleTo(target) {
			return vv.Convert(target), nil
		}
		return reflect.Value{}, fmt.Errorf("val type %T not assignable to %v", val, target)
	}
	switch v.Kind() {
	case reflect.Map:
		kv, err := convertKey(key, v.Type().Key())
		if err != nil {
			return err
		}
		vv, err := coerce(v.Type().Elem())
		if err != nil {
			return err
		}
		if v.IsNil() {
			return fmt.Errorf("map is nil")
		}
		v.SetMapIndex(kv, vv)
		return nil
	case reflect.Slice:
		i, err := strconv.Atoi(key)
		if err != nil || i < 0 || i >= v.Len() {
			return fmt.Errorf("index %q out of range", key)
		}
		vv, err := coerce(v.Type().Elem())
		if err != nil {
			return err
		}
		v.Index(i).Set(vv)
		return nil
	case reflect.Struct:
		f := v.FieldByName(key)
		if !f.IsValid() {
			return fmt.Errorf("field %q not found", key)
		}
		if !f.CanSet() {
			return fmt.Errorf("field %q not settable (unexported?)", key)
		}
		vv, err := coerce(f.Type())
		if err != nil {
			return err
		}
		f.Set(vv)
		return nil
	}
	return fmt.Errorf("container type %T not supported for set", obj)
}

// derefValue 解引用指针/接口，nil 返回错误。
func derefValue(v reflect.Value) (reflect.Value, error) {
	for v.Kind() == reflect.Pointer || v.Kind() == reflect.Interface {
		if v.IsNil() {
			return reflect.Value{}, fmt.Errorf("nil value")
		}
		v = v.Elem()
	}
	return v, nil
}

// getChildValue 在 v 上按 key 取子节点，返回 reflect.Value（保留可寻址性用于 set）。
func getChildValue(v reflect.Value, key string) (reflect.Value, error) {
	switch v.Kind() {
	case reflect.Map:
		kv, err := convertKey(key, v.Type().Key())
		if err != nil {
			return reflect.Value{}, err
		}
		mv := v.MapIndex(kv)
		if !mv.IsValid() {
			return reflect.Value{}, fmt.Errorf("key %q not found", key)
		}
		return mv, nil
	case reflect.Slice, reflect.Array:
		i, err := strconv.Atoi(key)
		if err != nil || i < 0 || i >= v.Len() {
			return reflect.Value{}, fmt.Errorf("index %q out of range", key)
		}
		return v.Index(i), nil
	case reflect.Struct:
		f := v.FieldByName(key)
		if !f.IsValid() {
			return reflect.Value{}, fmt.Errorf("field %q not found", key)
		}
		return f, nil
	}
	return reflect.Value{}, fmt.Errorf("cannot get %q from %v", key, v.Kind())
}

// getByPath 从 root 按 "." 路径取叶子值（只读，返回 Interface()）。
func getByPath(root any, path string) (any, error) {
	if path == "" {
		return nil, fmt.Errorf("empty path")
	}
	parts := strings.Split(path, ".")
	v := reflect.ValueOf(root)
	for _, part := range parts {
		var err error
		v, err = derefValue(v)
		if err != nil {
			return nil, fmt.Errorf("path %q: %w", path, err)
		}
		v, err = getChildValue(v, part)
		if err != nil {
			return nil, fmt.Errorf("path %q: %w", path, err)
		}
	}
	v, err := derefValue(v)
	if err != nil {
		return nil, fmt.Errorf("path %q: %w", path, err)
	}
	return v.Interface(), nil
}

// setByPath 从 root 按 "." 路径走到父容器，在父容器上 set 叶子。
// 这样即使中间是 struct 值字段，也能在父容器上原地改叶子。
//   - 父容器 map   → SetMapIndex（值类型叶子也可替换整个 entry）
//   - 父容器 slice → Index(i).Set（元素可寻址）
//   - 父容器 *struct → FieldByName().Set（指针出发，可寻址）
//   - 父容器 struct 值 → 报错（不可寻址），需把 root 改为 *T
func setByPath(root any, path string, val any) error {
	if path == "" {
		return fmt.Errorf("empty path")
	}
	parts := strings.Split(path, ".")
	v := reflect.ValueOf(root)

	// 走到父容器
	for _, part := range parts[:len(parts)-1] {
		var err error
		v, err = derefValue(v)
		if err != nil {
			return fmt.Errorf("path %q: %w", path, err)
		}
		v, err = getChildValue(v, part)
		if err != nil {
			return fmt.Errorf("path %q: %w", path, err)
		}
	}
	leaf := parts[len(parts)-1]

	// 父容器解引用（*struct → struct 可寻址）
	var err error
	v, err = derefValue(v)
	if err != nil {
		return fmt.Errorf("path %q: %w", path, err)
	}

	coerce := func(target reflect.Type) (reflect.Value, error) {
		vv := reflect.ValueOf(val)
		if !vv.IsValid() {
			return reflect.Zero(target), nil
		}
		if vv.Type().AssignableTo(target) {
			return vv, nil
		}
		if vv.Type().ConvertibleTo(target) {
			return vv.Convert(target), nil
		}
		// 尝试 JSON 往返（结构化数据）
		raw, jerr := json.Marshal(val)
		if jerr != nil {
			return reflect.Value{}, fmt.Errorf("val %T -> %v: %w", val, target, jerr)
		}
		out := reflect.New(target)
		if jerr := json.Unmarshal(raw, out.Interface()); jerr != nil {
			return reflect.Value{}, fmt.Errorf("val %T -> %v: %w", val, target, jerr)
		}
		return out.Elem(), nil
	}

	switch v.Kind() {
	case reflect.Map:
		kv, err := convertKey(leaf, v.Type().Key())
		if err != nil {
			return err
		}
		vv, err := coerce(v.Type().Elem())
		if err != nil {
			return err
		}
		if v.IsNil() {
			return fmt.Errorf("map is nil")
		}
		v.SetMapIndex(kv, vv)
		return nil
	case reflect.Slice, reflect.Array:
		i, err := strconv.Atoi(leaf)
		if err != nil || i < 0 || i >= v.Len() {
			return fmt.Errorf("index %q out of range", leaf)
		}
		vv, err := coerce(v.Type().Elem())
		if err != nil {
			return err
		}
		if !v.Index(i).CanSet() {
			return fmt.Errorf("slice/array element not settable (root 是值类型？改用 *T)")
		}
		v.Index(i).Set(vv)
		return nil
	case reflect.Struct:
		f := v.FieldByName(leaf)
		if !f.IsValid() {
			return fmt.Errorf("field %q not found", leaf)
		}
		if !f.CanSet() {
			return fmt.Errorf("field %q not settable（root 不可寻址？改用 *T 作 root）", leaf)
		}
		vv, err := coerce(f.Type())
		if err != nil {
			return err
		}
		f.Set(vv)
		return nil
	}
	return fmt.Errorf("parent type %v not supported for set", v.Kind())
}

// convertKey 把字符串 key 转成 map 的 key 类型。
func convertKey(key string, kt reflect.Type) (reflect.Value, error) {
	switch kt.Kind() {
	case reflect.String:
		return reflect.ValueOf(key).Convert(kt), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i, err := strconv.ParseInt(key, 10, 64)
		if err != nil {
			return reflect.Value{}, err
		}
		return reflect.ValueOf(i).Convert(kt), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		u, err := strconv.ParseUint(key, 10, 64)
		if err != nil {
			return reflect.Value{}, err
		}
		return reflect.ValueOf(u).Convert(kt), nil
	}
	return reflect.Value{}, fmt.Errorf("unsupported map key type %v", kt)
}

// callMethodByReflect 反射调用 obj 的方法，args 用 JSON 值（会按形参类型转换）。
// 返回所有返回值的 Interface() 列表。
func callMethodByReflect(obj any, method string, jsonArgs []any) ([]any, error) {
	v := reflect.ValueOf(obj)
	m := v.MethodByName(method)
	if !m.IsValid() {
		return nil, fmt.Errorf("method %q not found on %T", method, obj)
	}
	mt := m.Type()
	nFixed := mt.NumIn()
	if mt.IsVariadic() {
		nFixed--
	}
	if len(jsonArgs) < nFixed {
		return nil, fmt.Errorf("too few args: got %d, want >= %d", len(jsonArgs), nFixed)
	}
	if !mt.IsVariadic() && len(jsonArgs) != mt.NumIn() {
		return nil, fmt.Errorf("arg count: got %d, want %d", len(jsonArgs), mt.NumIn())
	}
	in := make([]reflect.Value, 0, len(jsonArgs))
	for i := 0; i < nFixed; i++ {
		// 复用 wazero_env.go 已有的 convertJSONArg
		cv, err := convertJSONArg(jsonArgs[i], mt.In(i))
		if err != nil {
			return nil, fmt.Errorf("arg %d: %w", i, err)
		}
		in = append(in, cv)
	}
	if mt.IsVariadic() {
		elemT := mt.In(nFixed).Elem()
		tail := reflect.MakeSlice(mt.In(nFixed), 0, len(jsonArgs)-nFixed)
		for i := nFixed; i < len(jsonArgs); i++ {
			cv, err := convertJSONArg(jsonArgs[i], elemT)
			if err != nil {
				return nil, fmt.Errorf("arg %d: %w", i, err)
			}
			tail = reflect.Append(tail, cv)
		}
		in = append(in, tail)
	}
	out := m.Call(in)
	res := make([]any, len(out))
	for i, r := range out {
		if r.IsValid() && r.CanInterface() {
			res[i] = r.Interface()
		}
	}
	return res, nil
}
