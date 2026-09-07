package vweb_dynamic

import (
	"bytes"
	"fmt"
	"go/constant"
	"io"
	"maps"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"text/template"

	"github.com/456vv/vweb/v3"
	"github.com/456vv/vweb/v3/builtin"
	"github.com/goplus/ixgo"
)

/*
注：如果你要使用更高级的功能，需要更新标准库源代码。
*/

// 点函数映射缓存（包名 → FuncMap）
// 使用 RWMutex 保证并发安全，读多写少场景下性能更好
var (
	dotPackage   = make(map[string]template.FuncMap)
	dotPackageMu sync.RWMutex
)

// ExtendPackage 扩展模板的包
//
//	name string				包名
//	deputy map[string]any 	函数集
func ExtendPackage(name string, deputy template.FuncMap) {
	dotPackage[name] = addFuncs(deputy)
}

func addFuncs(funcMap template.FuncMap) template.FuncMap {
	m := make(template.FuncMap)
	maps.Copy(m, funcMap)
	return m
}

// templateFuncMapError 从返回值中提取 error（支持单 error 或 []reflect.Value 最后一元素）
// 收紧边界检查，避免对空切片/非 error 类型产生误报或 panic。
func templateFuncMapError(v any) error {
	if errs, ok := v.([]reflect.Value); ok {
		l := len(errs)
		if l == 0 {
			return nil
		}
		err := errs[l-1]
		if !err.IsValid() {
			return nil
		}
		if err.CanInterface() {
			inf := err.Interface()
			if e, ok := inf.(error); ok {
				return e
			}
			if inf == nil {
				return nil
			}
		}
		return fmt.Errorf("the last argument is not of the wrong type: %v", err)
	}
	if err, ok := v.(error); ok {
		return err
	}
	return nil
}

// callMethod 通过反射调用对象方法（支持值接收者与指针接收者）
func callMethod(f any, name string, args ...any) ([]any, error) {
	if f == nil {
		return nil, fmt.Errorf("callMethod: receiver is nil")
	}
	vfn := reflect.ValueOf(f)
	orig := vfn

	// 尝试指针接收者
	if vfn.Kind() == reflect.Ptr {
		if method := vfn.MethodByName(name); method.IsValid() && method.Kind() == reflect.Func {
			return call(method.Interface(), args...)
		}
		// 解引用继续尝试值接收者
		vfn = reflect.Indirect(vfn)
	}

	// 尝试值接收者（Struct）
	if vfn.Kind() == reflect.Struct {
		if method := vfn.MethodByName(name); method.IsValid() && method.Kind() == reflect.Func {
			return call(method.Interface(), args...)
		}
	}

	// 最后尝试原始值上的方法（覆盖 interface 等情况）
	if method := orig.MethodByName(name); method.IsValid() && method.Kind() == reflect.Func {
		return call(method.Interface(), args...)
	}

	return nil, fmt.Errorf("the `%s` method was not found in `%v`", name, f)
}

// call 通用反射调用辅助（依赖 vweb.ExecCall）
func call(f any, args ...any) ([]any, error) {
	if f == nil {
		return nil, fmt.Errorf("call: function is nil")
	}
	var ef vweb.ExecCall
	if err := ef.Func(f, args...); err != nil {
		return nil, err
	}
	return ef.Exec(), nil
}

// pkgFuncs 包装包内函数，用于模板调用
type pkgFuncs struct {
	f reflect.Value
}

// call 调用被包装的函数。返回值规则：
// - 0 个返回值 → nil
// - 1 个返回值 → 该值
// - 多个返回值 → []any
// 出错时直接返回 error 作为 any（模板侧可处理）
func (T *pkgFuncs) call(args ...any) any {
	result, err := call(T.f, args...)
	if err != nil {
		return err
	}
	switch len(result) {
	case 0:
		return nil
	case 1:
		return result[0]
	default:
		return result
	}
}

// pkgTypes 包装命名类型/别名类型，用于实例化
type pkgTypes struct {
	t reflect.Type
}

// new 创建类型实例。
// - 无参数：使用 builtin.Init 初始化零值
// - 有参数：尝试把 args[0] 转换为目标类型；失败则 panic（与原行为一致）
// 多余参数被忽略
func (T *pkgTypes) new(args ...any) any {
	pv := reflect.New(T.t).Elem()
	if len(args) != 0 {
		if args[0] != nil && !builtin.Convert(pv, args[0]) {
			panic("The type are not the same and cannot be converted")
		}
	} else {
		builtin.Init(pv, false)
	}
	return pv.Interface()
}

// pkgInterface 包装命名类型/别名类型，用于实例化
type pkgInterface struct {
	t reflect.Type
}

// to 将任意值转换为目标接口类型（可能 panic，若不可转换）
// 增加 nil 检查，避免对 nil 调用 Convert 导致 panic。
func (T *pkgInterface) to(v any) any {
	if v == nil {
		return nil
	}
	return reflect.ValueOf(v).Convert(T.t).Interface()
}

// pkgTypedConsts 包装有类型常量
type pkgTypedConsts struct {
	tc ixgo.TypedConst
}

// get 返回有类型常量的实际值
func (T *pkgTypedConsts) get() any {
	return constant.Val(T.tc.Value)
}

// importPkg 根据包名导入并构建 template.FuncMap。
// 结果会被缓存；并发安全；包不存在时返回空 FuncMap 并缓存，避免重复查找。
// 保持原有双重检查锁定；对空包名做快速返回，避免无意义查找。
func importPkg(name string) template.FuncMap {
	if name == "" {
		return template.FuncMap{}
	}

	// 快速路径：读锁检查缓存
	dotPackageMu.RLock()
	if fm, ok := dotPackage[name]; ok {
		dotPackageMu.RUnlock()
		return fm
	}
	dotPackageMu.RUnlock()

	// 慢路径：写锁构建
	dotPackageMu.Lock()
	defer dotPackageMu.Unlock()

	// double-check，防止并发重复构建
	if fm, ok := dotPackage[name]; ok {
		return fm
	}

	fm := make(template.FuncMap)
	if pkg, ok := ixgo.LookupPackage(name); ok {
		// 接口，用处不大（主要用于类型断言/转换）
		for n, t := range pkg.Interfaces {
			pt := pkgInterface{t: t}
			fm[n] = pt.to
		}
		// 类型-实质（命名类型）
		for n, t := range pkg.NamedTypes {
			pt := pkgTypes{t: t}
			fm[n] = pt.new
		}
		// 类型-别名
		for n, t := range pkg.AliasTypes {
			pt := pkgTypes{t: t}
			fm[n] = pt.new
		}
		// 变量：注册时快照当前值（与原逻辑一致）
		// 注意：若变量后续被修改，模板中看到的仍是注册时刻的值
		for n, v := range pkg.Vars {
			fm[n] = v.Elem().Interface()
		}
		// 函数
		for n, f := range pkg.Funcs {
			pf := pkgFuncs{f: f}
			fm[n] = pf.call
		}
		// 常量
		for n, t := range pkg.TypedConsts {
			pt := pkgTypedConsts{tc: t}
			fm[n] = pt.get
		}
	}
	// 无论是否找到包，都缓存结果（空 map 也缓存，避免反复 Lookup）
	dotPackage[name] = fm
	return fm
}

type wasmFuncType int

const (
	wasmFuncTypeImport wasmFuncType = iota
	ForMethod
	ForType
	DepthField
	CopyStruct
	CopyStructDeep
	Convert
	Init
	Value
	Call
	CallMethod
)

var wasmFunc = map[wasmFuncType]any{
	Call:       call,
	CallMethod: callMethod,
}

// TemplateFunc 模板函数映射（可导出）。只读使用，初始化后不再修改，并发读安全。
var TemplateFunc = template.FuncMap{
	"Import":         importPkg,
	"ForMethod":      vweb.ForMethod,
	"ForType":        vweb.ForType,
	"Register":       builtin.Register,
	"DepthField":     builtin.DepthField,
	"CopyStruct":     builtin.CopyStruct,
	"CopyStructDeep": builtin.CopyStructDeep,
	"Convert":        builtin.Convert,
	"To":             builtin.To,
	"Init":           builtin.Init,
	"Value":          builtin.Value, // Value(v) reflect.Value
	"Call":           call,
	"CallMethod":     callMethod,
	"Defer":          func(f any, args ...any) func() { return func() { call(f, args...) } },
	"DeferMethod":    func(f any, name string, args ...any) func() { return func() { callMethod(f, name, args...) } },
	"Go":             func(f func()) { go f() },
	"PtrTo":          func(inf any) any { return reflect.ValueOf(inf).Elem().Interface() },
	"ToPtr":          func(inf any) any { return reflect.ValueOf(inf).Addr().Interface() },
	"Nil":            func() any { return nil },
	"NotNil":         func(inf any) bool { return inf != nil },
	"IsNil":          func(inf any) bool { return inf == nil },
	"Type":           builtin.Type, // Type(v) reflect.Type
	"MustType":       builtin.MustType,
	"Panic":          builtin.Panic,  // Panic(v)
	"New":            builtin.New,    //New(T)
	"Make":           builtin.Make,   // Make([]T, length, cap)|Make(map[T]T, length)|Make(Chan, length)
	"Delete":         builtin.Delete, // Delete(map[T]T, "key")
	"SetUnexported":  builtin.SetUnexported,
	"Set":            builtin.Set, // Set([]T, 位置0,值1, ...)|Set(map[T]T, 键名0,值1, ...)|Set(struct{}, 名称0,值1, ...)
	"GetUnexported":  builtin.GetUnexported,
	"Get":            builtin.Get,      // Get(map[T]T/[]T/struct{}/string/number, key)
	"Len":            builtin.Len,      // Len([]T/string/map[T]T)
	"Cap":            builtin.Cap,      // Cap([]T)
	"GetSlice":       builtin.GetSlice, // GetSlice([]T, 1, 5)
	"SetSlice":       builtin.SetSlice, // SetSlice([]T, 1,5, T)
	"Copy":           builtin.Copy,     // Copy([]T, []T)
	"Append":         builtin.Append,   // Append([]T, T...)
	"Compute":        builtin.Compute,  // Compute(1, "+", 2)
	"Or":             builtin.Or,       // Or returns 1 || true
	"And":            builtin.And,      // And returns 1 && true
	"Not":            builtin.Not,      // Not returns !a
	"LT":             builtin.LT,       // LT returns a < b
	"GT":             builtin.GT,       // GT returns a > b
	"LE":             builtin.LE,       // LE returns a <= b
	"GE":             builtin.GE,       // GE returns a >= b
	"EQ":             builtin.EQ,       // EQ returns a == b
	"NE":             builtin.NE,       // NE returns a != b
	"TrySend":        builtin.TrySend,  // TrySend(*Chan, value)	不阻塞
	"TryRecv":        builtin.TryRecv,  // value = TryRecv(*Chan)	不阻塞
	"Send":           builtin.Send,     // Send(*Chan, value)
	"Recv":           builtin.Recv,     // Recv(*Chan)
	"Close":          builtin.Close,    // Close(*Chan)
	"Error": func(v any) bool {
		return templateFuncMapError(v) != nil
	},
	"NotError": func(v any) bool {
		return templateFuncMapError(v) == nil
	},
}

// entryname 根据默认入口名与文件名计算实际入口函数名
func entryname(name1, name2 string) string {
	if name1 != "" {
		return name1
	}

	base := filepath.Base(name2)
	pos := strings.IndexAny(base, ".")
	isDir := true
	if pos != -1 {
		base = base[:pos]
		isDir = false
	}

	if isDir || base == "index" || base == "" {
		return "Main"
	}

	// 仅允许字母数字下划线（保证合法 Go 标识符）
	for _, v := range base {
		if !((v >= '0' && v <= '9') || (v >= 'A' && v <= 'Z') || (v >= 'a' && v <= 'z') || v == '_') {
			return "Main"
		}
	}

	// 首字母大写（ASCII）
	if base[0] >= 'a' && base[0] <= 'z' {
		return string(base[0]-32) + base[1:]
	}
	return base
}

// fileHeaderLine 从缓冲区头部读取连续的 // 注释行（作为模板头）
// 兼容 \r\n，空注释行跳过，EOF 安全处理。
func fileHeaderLine(buf *bytes.Buffer) (l []string) {
	for {
		// 先 peek 前两个字节，避免错误消费后无法完整回退
		data := buf.Bytes()
		if len(data) < 2 || data[0] != '/' || data[1] != '/' {
			return
		}

		// 确认是 //，安全消费整行
		// 先丢掉已经确认的两个 '/'
		buf.Next(2)

		// 读取到行尾
		line, err := buf.ReadBytes('\n')
		if err != nil && err != io.EOF {
			return
		}

		// 去掉换行符与回车
		line = bytes.TrimRight(line, "\r\n")
		content := bytes.TrimSpace(line)
		if len(content) == 0 {
			// 空注释行跳过
			if err == io.EOF {
				return
			}
			continue
		}

		l = append(l, string(content))
		if err == io.EOF {
			return
		}
	}
}
