package vweb_dynamic

import (
	"fmt"
	"go/constant"
	"path/filepath"
	"reflect"
	"strings"
	"text/template"

	"github.com/456vv/vweb/v2"
	"github.com/456vv/vweb/v2/builtin"
	"github.com/goplus/igop"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

/*
注：如果你要使用更高级的功能，需要更新标准库源代码。
*/

// 点函数映射
var dotPackage = make(map[string]template.FuncMap)

// ExtendPackage 扩展模板的包
//
//	name string				包名
//	deputy map[string]any 	函数集
func ExtendPackage(name string, deputy template.FuncMap) {
	dotPackage[name] = addFuncs(deputy)
}

func addFuncs(funcMap template.FuncMap) template.FuncMap {
	m := make(template.FuncMap)
	for n, v := range funcMap {
		m[n] = v
	}
	return m
}

func templateFuncMapError(v any) error {
	if errs, ok := v.([]reflect.Value); ok {
		l := len(errs)
		if l == 0 {
			return nil
		}
		err := errs[l-1]
		if err.CanInterface() {
			inf := err.Interface()
			if e, ok := inf.(error); ok {
				return e
			} else if inf == nil {
				return nil
			}
		}
		return fmt.Errorf("the last argument is not of the wrong type: %v", err)
	} else if err, ok := v.(error); ok {
		return err
	}
	return nil
}

func callMethod(f any, name string, args ...any) ([]any, error) {
	var isMethodFunc bool
	vfn := reflect.ValueOf(f)
	if vfn.Kind() == reflect.Ptr {
		// func (T *A) B(){}
		if vfn.NumMethod() > 0 {
			t := vfn.MethodByName(name)
			if t.Kind() == reflect.Func {
				vfn = t
				isMethodFunc = true
			}
		} else {
			vfn = reflect.Indirect(vfn)
		}
	}
	if vfn.Kind() == reflect.Struct {
		if vfn.NumMethod() > 0 {
			// func (T A) C(){}
			t := vfn.MethodByName(name)
			if t.Kind() == reflect.Func {
				vfn = t
				isMethodFunc = true
			}
		}
	}
	if !isMethodFunc {
		return nil, fmt.Errorf("the `%s` method was not found in `%v`", name, f)
	}
	return call(vfn, args...)
}

func call(f any, args ...any) ([]any, error) {
	var ef vweb.ExecCall
	if err := ef.Func(f, args...); err != nil {
		return nil, err
	}
	return ef.Exec(), nil
}

type pkgFuncs struct {
	f reflect.Value
}

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

type pkgTypes struct {
	t reflect.Type
}

func (T *pkgTypes) new(args ...any) any {
	var pv reflect.Value
	if len(args) != 0 {
		ptyp := reflect.PointerTo(T.t)
		pv = reflect.New(ptyp).Elem()
		if args[0] != nil && !builtin.Convert(pv, args[0]) {
			panic("The type are not the same and cannot be converted")
		}
	} else {
		pv = reflect.New(T.t)
		builtin.Init(pv)
	}
	return pv.Interface()
}

type pkgInterface struct {
	t reflect.Type
}

func (T *pkgInterface) to(v any) any {
	return reflect.ValueOf(v).Convert(T.t).Interface()
}

type pkgTypedConsts struct {
	tc igop.TypedConst
}

func (T *pkgTypedConsts) get() any {
	return constant.Val(T.tc.Value)
}

func importPkg(name string) template.FuncMap {
	if fm, ok := dotPackage[name]; ok {
		return fm
	}

	fm := make(template.FuncMap)
	if pkg, ok := igop.LookupPackage(name); ok {
		// 接口，用处不大
		for n, t := range pkg.Interfaces {
			pt := pkgInterface{t: t}
			fm[n] = pt.to
		}
		// 类型-实质
		for n, t := range pkg.NamedTypes {
			pt := pkgTypes{t: t}
			fm[n] = pt.new
		}
		// 类型-别名
		for n, t := range pkg.AliasTypes {
			pt := pkgTypes{t: t}
			fm[n] = pt.new
		}
		// 变量
		for n, v := range pkg.Vars {
			fm[n] = v.Elem().Interface
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
	dotPackage[name] = fm
	return fm
}

// 模板函数映射
var TemplateFunc = template.FuncMap{
	"Import":         importPkg,
	"ForMethod":      vweb.ForMethod,
	"ForType":        vweb.ForType,
	"InDirect":       vweb.InDirect,
	"DepthField":     vweb.DepthField,
	"CopyStruct":     vweb.CopyStruct,
	"CopyStructDeep": vweb.CopyStructDeep,
	"Convert":        builtin.Convert,
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
	"Bytes":          builtin.Bytes,
	"Runes":          builtin.Runes,
	"Pointer":        builtin.Pointer,
	"Append":         builtin.Append, // Append([]T, value...)
	"Uintptr":        builtin.Uintptr,
	"SetUintptr":     func(s *uintptr, v uintptr) *uintptr { *s = v; return s },
	"Byte":           builtin.Byte,
	"SetByte":        func(s *byte, v byte) *byte { *s = v; return s },
	"Rune":           builtin.Rune,
	"SetRune":        func(s *rune, v rune) *rune { *s = v; return s },
	"String":         builtin.String,
	"SetString":      func(s *string, v string) *string { *s = v; return s },
	"Bool":           builtin.Bool,
	"SetBool":        func(i *bool, v bool) *bool { *i = v; return i },
	"Int":            builtin.Int,
	"SetInt":         func(i *int, v int) *int { *i = v; return i },
	"Int8":           builtin.Int8,
	"SetInt8":        func(i *int8, v int8) *int8 { *i = v; return i },
	"Int16":          builtin.Int16,
	"SetInt16":       func(i *int16, v int16) *int16 { *i = v; return i },
	"Int32":          builtin.Int32,
	"SetInt32":       func(i *int32, v int32) *int32 { *i = v; return i },
	"Int64":          builtin.Int64,
	"SetInt64":       func(i *int64, v int64) *int64 { *i = v; return i },
	"Uint":           builtin.Uint,
	"SetUint":        func(i *uint, v uint) *uint { *i = v; return i },
	"Uint8":          builtin.Uint8,
	"SetUint8":       func(i *uint8, v uint8) *uint8 { *i = v; return i },
	"Uint16":         builtin.Uint16,
	"SetUint16":      func(i *uint16, v uint16) *uint16 { *i = v; return i },
	"Uint32":         builtin.Uint32,
	"SetUint32":      func(i *uint32, v uint32) *uint32 { *i = v; return i },
	"Uint64":         builtin.Uint64,
	"SetUint64":      func(i *uint64, v uint64) *uint64 { *i = v; return i },
	"Float32":        builtin.Float32,
	"SetFloat32":     func(f *float32, v float32) *float32 { *f = v; return f },
	"Float64":        builtin.Float64,
	"SetFloat64":     func(f *float64, v float64) *float64 { *f = v; return f },
	"Complex64":      builtin.Complex64,
	"SetComplex64":   func(c *complex64, v complex64) *complex64 { *c = v; return c },
	"Complex128":     builtin.Complex128,
	"SetComplex128":  func(c *complex128, v complex128) *complex128 { *c = v; return c },
	"Type":           builtin.Type,      // Type(v) reflect.Type
	"Panic":          builtin.Panic,     // Panic(v)
	"Make":           builtin.Make,      // Make([]T, length, cap)|Make([T]T, length)|Make(Chan, length)
	"MapFrom":        builtin.MapFrom,   // MapFrom(M, T1,V1, T2, V2, ...)
	"SliceFrom":      builtin.SliceFrom, // SliceFrom(S, 值0, 值1,...)
	"Delete":         builtin.Delete,    // Delete(map[T]T, "key")
	"Set":            builtin.Set,       // Set([]T, 位置0,值1, 位置1,值2, 位置2,值3)|Set(map[T]T, 键名0,值1, 键名1,值2, 键名2,值3)|Set(struct{}, 名称0,值1, 名称1,值2, 名称2,值3)
	"Get":            builtin.Get,       // Get(map[T]T/[]T/struct{}/string/number, key)
	"Len":            builtin.Len,       // Len([]T/string/map[T]T)
	"Cap":            builtin.Cap,       // Cap([]T)
	"GetSlice":       builtin.GetSlice,  // GetSlice([]T, 1, 5)
	"GetSlice3":      builtin.GetSlice3, // GetSlice3([]T, 1, 5, 7)
	"Copy":           builtin.Copy,      // Copy([]T, []T)
	"Compute":        builtin.Compute,   // Compute(1, "+", 2)
	"Inc":            builtin.Inc,       // Inc returns a+1
	"Dec":            builtin.Dec,       // Dec returns a-1
	"Neg":            builtin.Neg,       // Neg returns -a
	"Mul":            builtin.Mul,       // Mul returns a*b
	"Quo":            builtin.Quo,       // Quo returns a/b
	"Mod":            builtin.Mod,       // Mod returns a%b
	"Add":            builtin.Add,       // Add returns a+b
	"Sub":            builtin.Sub,       // Sub returns a-b
	"BitLshr":        builtin.BitLshr,   // BitLshr returns a << b
	"BitRshr":        builtin.BitRshr,   // BitRshr returns a >> b
	"BitXor":         builtin.BitXor,    // BitXor returns a ^ b
	"BitAnd":         builtin.BitAnd,    // BitAnd returns a & b
	"BitOr":          builtin.BitOr,     // BitOr returns a | b
	"BitNot":         builtin.BitNot,    // BitNot returns ^a
	"BitAndNot":      builtin.BitAndNot, // BitAndNot returns a &^ b
	"Or":             builtin.Or,        // Or returns 1 || true
	"And":            builtin.And,       // And returns 1 && true
	"Not":            builtin.Not,       // Not returns !a
	"LT":             builtin.LT,        // LT returns a < b
	"GT":             builtin.GT,        // GT returns a > b
	"LE":             builtin.LE,        // LE returns a <= b
	"GE":             builtin.GE,        // GE returns a >= b
	"EQ":             builtin.EQ,        // EQ returns a == b
	"NE":             builtin.NE,        // NE returns a != b
	"TrySend":        builtin.TrySend,   // TrySend(*Chan, value)	不阻塞
	"TryRecv":        builtin.TryRecv,   // value = TryRecv(*Chan)	不阻塞
	"Send":           builtin.Send,      // Send(*Chan, value)
	"Recv":           builtin.Recv,      // Recv(*Chan)
	"Close":          builtin.Close,     // Close(*Chan)
	"ChanOf":         builtin.ChanOf,    // ChanOf(T)
	"MakeChan":       builtin.MakeChan,  // MakeChan(T, size)
	"Max":            builtin.Max,       // Max(a1, a2 ...)
	"Min":            builtin.Min,       // Min(a1, a2 ...)
	"Error": func(v any) bool {
		return templateFuncMapError(v) != nil
	},
	"NotError": func(v any) bool {
		return templateFuncMapError(v) == nil
	},
}

func entryname(name string) string {
	base := filepath.Base(name)
	pos := strings.IndexAny(base, ".")
	if pos != -1 {
		base = base[:pos]
	}
	base = strings.TrimSpace(base)
	base = cases.Title(language.English).String(base) // strings.Title(base)
	if base == "\\" || base == "Index" || base == "" || base[0] < 'A' && base[0] > 'Z' {
		base = "Main"
	}
	return base
}
