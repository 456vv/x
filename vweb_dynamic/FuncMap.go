package vweb_dynamic
import(
    "text/template"
    "reflect"
    "fmt"
	"github.com/456vv/vweb/v2"
	"github.com/456vv/vweb/v2/builtin"
)

/*
注：如果你要使用更高级的功能，需要更新标准库源代码。
*/

// 点函数映射
var dotPackage = make(map[string]map[string]interface{})

//ExtendTemplatePackage 扩展模板的包
//	pkgName string					包名
//	deputy map[string]interface{} 	函数集
func ExtendPackage(pkgName string, deputy template.FuncMap) {
	if _, ok := dotPackage[pkgName]; !ok {
		dotPackage[pkgName] = make(template.FuncMap)
	}
	for name, fn  := range deputy {
		dotPackage[pkgName][name]=fn
	}
}

func templateFuncMapError(v interface{}) error {
    if errs, ok := v.([]reflect.Value); ok {
        l   := len(errs)
        if l==0 {return nil}
        err := errs[l-1]
        if err.CanInterface() {
        	inf := err.Interface()
            if e, ok := inf.(error); ok {
                return e
            }else if inf == nil {
            	return nil
            }
        }
        return fmt.Errorf("Error: 判断最后一个参数不是错误类型。%s", err)
    }else if err, ok := v.(error); ok {
        return err
    }
    return nil
}

func callMethod(f interface{}, name string, args ...interface{}) ([]interface{}, error) {
	var isMethodFunc bool
	vfn := reflect.ValueOf(f)
	if vfn.Kind() == reflect.Ptr {
		//func (T *A) B(){}
		if vfn.NumMethod() > 0  {
			t := vfn.MethodByName(name)
			if t.Kind() == reflect.Func {
				vfn = t
				isMethodFunc = true
			}
		}else{
			vfn = reflect.Indirect(vfn)
		}
	}
	if vfn.Kind() == reflect.Struct {
		if vfn.NumMethod() > 0  {
			//func (T A) C(){}
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

func call(f interface{}, args ...interface{}) ([]interface{}, error){
	var ef vweb.ExecCall
	if err := ef.Func(f, args...); err != nil {
		return nil, err
	}
	return ef.Exec(), nil
}

// 模板函数映射
var TemplateFunc = template.FuncMap{
	"Import": func(pkgName string) template.FuncMap {return dotPackage[pkgName]},
    "ForMethod": vweb.ForMethod,
    "ForType": vweb.ForType,
    "InDirect": vweb.InDirect,
    "DepthField": vweb.DepthField,
    "CopyStruct": vweb.CopyStruct,
    "CopyStructDeep": vweb.CopyStructDeep,
    "GoTypeTo":builtin.GoTypeTo,
    "GoTypeInit":builtin.GoTypeInit,
    "Value":builtin.Value,						//Value(v) reflect.Value
	"_Value_":func(s []reflect.Value, v ...reflect.Value) []reflect.Value {return append(s, v...)},
	"Call":call,
	"CallMethod":callMethod,
	"Defer":func(f interface{}, args ...interface{}) func() {return func(){call(f, args...)}},
	"DeferMethod":func(f interface{}, name string, args ...interface{}) func() {return func(){callMethod(f, name, args...)}},
	"Go":func(f func()){go f()},
	"PtrTo":func(inf interface{}) interface{} {v := reflect.Indirect(reflect.ValueOf(inf));return typeSelect(v)},
    "ToPtr":func(inf interface{}) interface{} {return &inf},
	"Nil":func() interface{} {return nil},
	"NotNil":func(inf interface{}) bool {return inf != nil},
	"IsNil":func(inf interface{}) bool {return inf == nil},
    "Bytes": builtin.Bytes,
    "Runes": builtin.Runs,
    "Pointer":builtin.Pointer,
    "Append": builtin.Append,			//Append([]T, value...)
    "_Interface_": func(inf interface{}, v ...interface{}) []interface{} {if inf==nil && len(v)==0 {return make([]interface{}, 0, 0)};if s, ok := inf.([]interface{}); ok {return append(s, v...)};return nil},
    "_Interface": func(s interface{}) *interface{} {return &s},
    "Interface_": func(s *interface{}) interface{} {return *s},
    "SetInterface": func(s *interface{}, v interface{}) *interface{} {*s = v;return s},
    "Uintptr":builtin.Uintptr,
    "_Uintptr_": func(inf interface{}, v ...uintptr) []uintptr {if inf==nil && len(v)==0 {return make([]uintptr, 0, 0)};if s, ok := inf.([]uintptr); ok {return append(s, v...)};return nil},
    "_Uintptr": func(s uintptr) *uintptr {return &s},
    "Uintptr_": func(s *uintptr) uintptr {return *s},
    "SetUintptr": func(s *uintptr, v uintptr) *uintptr {*s = v;return s},
    "Byte":builtin.Byte,
    "_Byte_": func(inf interface{}, v ...byte) []byte {if inf==nil && len(v)==0 {return make([]byte, 0, 0)};if s, ok := inf.([]byte); ok {return append(s, v...)};return nil},
    "_Byte": func(s byte) *byte {return &s},
    "Byte_": func(s *byte) byte {return *s},
    "SetByte": func(s *byte, v byte) *byte {*s = v;return s},
    "Rune":builtin.Rune,
    "_Rune_": func(inf interface{}, v ...rune) []rune {if inf==nil && len(v)==0 {return make([]rune, 0, 0)};if s, ok := inf.([]rune); ok {return append(s, v...)};return nil},
    "_Rune": func(s rune) *rune {return &s},
    "Rune_": func(s *rune) rune {return *s},
    "SetRune": func(s *rune, v rune) *rune {*s = v;return s},
    "String":builtin.String,
    "_String_": func(inf interface{}, v ...string) []string {if inf==nil && len(v)==0 {return make([]string, 0, 0)};if s, ok := inf.([]string); ok {return append(s, v...)};return nil},
    "_String": func(s string) *string {return &s},
    "String_": func(s *string) string {return *s},
    "SetString": func(s *string, v string) *string {*s = v;return s},
    "Bool":builtin.Bool,
    "_Bool_": func(inf interface{}, v ...bool) []bool {if inf==nil && len(v)==0 {return make([]bool, 0, 0)};if s, ok := inf.([]bool); ok {return append(s, v...)};return nil},
    "_Bool": func(i bool) *bool {return &i},
    "Bool_": func(i *bool) bool {return *i},
    "SetBool": func(i *bool, v bool) *bool {*i = v;return i},
    "Int":builtin.Int,
    "_Int_": func(inf interface{}, v ...int) []int {if inf==nil && len(v)==0 {return make([]int, 0, 0)};if s, ok := inf.([]int); ok {return append(s, v...)};return nil},
    "_Int": func(i int) *int {return &i},
    "Int_": func(i *int) int {return *i},
    "SetInt": func(i *int, v int) *int {*i = v;return i},
    "Int8":builtin.Int8,
    "_Int8_": func(inf interface{}, v ...int8) []int8 {if inf==nil && len(v)==0 {return make([]int8, 0, 0)};if s, ok := inf.([]int8); ok {return append(s, v...)};return nil},
    "_Int8": func(i int8) *int8 {return &i},
    "Int8_": func(i *int8) int8 {return *i},
    "SetInt8": func(i *int8, v int8) *int8 {*i = v;return i},
    "Int16":builtin.Int16,
     "_Int16_": func(inf interface{}, v ...int16) []int16 {if inf==nil && len(v)==0 {return make([]int16, 0, 0)};if s, ok := inf.([]int16); ok {return append(s, v...)};return nil},
    "_Int16": func(i int16) *int16 {return &i},
    "Int16_": func(i *int16) int16 {return *i},
    "SetInt16": func(i *int16, v int16) *int16 {*i = v;return i},
    "Int32":builtin.Int32,
     "_Int32_": func(inf interface{}, v ...int32) []int32 {if inf==nil && len(v)==0 {return make([]int32, 0, 0)};if s, ok := inf.([]int32); ok {return append(s, v...)};return nil},
    "_Int32": func(i int32) *int32 {return &i},
    "Int32_": func(i *int32) int32 {return *i},
    "SetInt32": func(i *int32, v int32) *int32 {*i = v;return i},
    "Int64":builtin.Int64,
     "_Int64_": func(inf interface{}, v ...int64) []int64 {if inf==nil && len(v)==0 {return make([]int64, 0, 0)};if s, ok := inf.([]int64); ok {return append(s, v...)};return nil},
    "_Int64": func(i int64) *int64 {return &i},
    "Int64_": func(i *int64) int64 {return *i},
    "SetInt64": func(i *int64, v int64) *int64 {*i = v;return i},
    "Uint":builtin.Uint,
     "_Uint_": func(inf interface{}, v ...uint) []uint {if inf==nil && len(v)==0 {return make([]uint, 0, 0)};if s, ok := inf.([]uint); ok {return append(s, v...)};return nil},
    "_Uint": func(i uint) *uint {return &i},
    "Uint_": func(i *uint) uint {return *i},
    "SetUint": func(i *uint, v uint) *uint {*i = v;return i},
    "Uint8":builtin.Uint8,
     "_Uint8_": func(inf interface{}, v ...uint8) []uint8 {if inf==nil && len(v)==0 {return make([]uint8, 0, 0)};if s, ok := inf.([]uint8); ok {return append(s, v...)};return nil},
    "_Uint8": func(i uint8) *uint8 {return &i},
    "Uint8_": func(i *uint8) uint8 {return *i},
    "SetUint8": func(i *uint8, v uint8) *uint8 {*i = v;return i},
    "Uint16":builtin.Uint16,
     "_Uint16_": func(inf interface{}, v ...uint16) []uint16 {if inf==nil && len(v)==0 {return make([]uint16, 0, 0)};if s, ok := inf.([]uint16); ok {return append(s, v...)};return nil},
    "_Uint16": func(i uint16) *uint16 {return &i},
    "Uint16_": func(i *uint16) uint16 {return *i},
    "SetUint16": func(i *uint16, v uint16) *uint16 {*i = v;return i},
    "Uint32":builtin.Uint32,
     "_Uint32_": func(inf interface{}, v ...uint32) []uint32 {if inf==nil && len(v)==0 {return make([]uint32, 0, 0)};if s, ok := inf.([]uint32); ok {return append(s, v...)};return nil},
    "_Uint32": func(i uint32) *uint32 {return &i},
    "Uint32_": func(i *uint32) uint32 {return *i},
    "SetUint32": func(i *uint32, v uint32) *uint32 {*i = v;return i},
    "Uint64":builtin.Uint64,
     "_Uint64_": func(inf interface{}, v ...uint64) []uint64 {if inf==nil && len(v)==0 {return make([]uint64, 0, 0)};if s, ok := inf.([]uint64); ok {return append(s, v...)};return nil},
    "_Uint64": func(i uint64) *uint64 {return &i},
    "Uint64_": func(i *uint64) uint64 {return *i},
    "SetUint64": func(i *uint64, v uint64) *uint64 {*i = v;return i},
    "Float32":builtin.Float32,
     "_Float32_": func(inf interface{}, v ...float32) []float32 {if inf==nil && len(v)==0 {return make([]float32, 0, 0)};if s, ok := inf.([]float32); ok {return append(s, v...)};return nil},
    "_Float32": func(f float32) *float32 {return &f},
    "Float32_": func(f *float32) float32 {return *f},
    "SetFloat32": func(f *float32, v float32) *float32 {*f = v;return f},
    "Float64":builtin.Float64,
     "_Float64_": func(inf interface{}, v ...float64) []float64 {if inf==nil && len(v)==0 {return make([]float64, 0, 0)};if s, ok := inf.([]float64); ok {return append(s, v...)};return nil},
    "_Float64": func(f float64) *float64 {return &f},
    "Float64_": func(f *float64) float64 {return *f},
    "SetFloat64": func(f *float64, v float64) *float64 {*f = v;return f},
    "Complex64":builtin.Complex64,
     "_Complex64_": func(inf interface{}, v ...complex64) []complex64 {if inf==nil && len(v)==0 {return make([]complex64, 0, 0)};if s, ok := inf.([]complex64); ok {return append(s, v...)};return nil},
    "_Complex64": func(c complex64) *complex64 {return &c},
    "Complex64_": func(c *complex64) complex64 {return *c},
    "SetComplex64": func(c *complex64, v complex64) *complex64 {*c = v;return c},
    "Complex128":builtin.Complex128,
     "_Complex128_": func(inf interface{}, v ...complex128) []complex128 {if inf==nil && len(v)==0 {return make([]complex128, 0, 0)};if s, ok := inf.([]complex128); ok {return append(s, v...)};return nil},
    "_Complex128": func(c complex128) *complex128 {return &c},
    "Complex128_": func(c *complex128) complex128 {return *c},
    "SetComplex128": func(c *complex128, v complex128) *complex128 {*c = v;return c},
    "Type":builtin.Type,						//Type(v) reflect.Type
    "Panic":builtin.Panic,						//Panic(v)
    "Make":builtin.Make,						//Make([]T, length, cap)|Make([T]T, length)|Make(Chan, length)
    "MapFrom":builtin.MapFrom,					//MapFrom(M, T1,V1, T2, V2, ...)
    "SliceFrom":builtin.SliceFrom,				//SliceFrom(S, 值0, 值1,...)
    "Delete":builtin.Delete,					//Delete(map[T]T, "key")
    "Set":builtin.Set,							//Set([]T, 位置0,值1, 位置1,值2, 位置2,值3)|Set(map[T]T, 键名0,值1, 键名1,值2, 键名2,值3)|Set(struct{}, 名称0,值1, 名称1,值2, 名称2,值3)
    "Get":builtin.Get,							//Get(map[T]T/[]T/struct{}/string/number, key)
    "Len":builtin.Len,							//Len([]T/string/map[T]T)
    "Cap":builtin.Cap,							//Cap([]T)
    "GetSlice":builtin.GetSlice,				//GetSlice([]T, 1, 5)
    "GetSlice3":builtin.GetSlice3,				//GetSlice3([]T, 1, 5, 7)
    "Copy":builtin.Copy,						//Copy([]T, []T)
    "Compute": builtin.Compute,					//Compute(1, "+", 2)
    "Inc":builtin.Inc,							//Inc returns a+1
    "Dec":builtin.Dec,							//Dec returns a-1
    "Neg":builtin.Neg,							//Neg returns -a
    "Mul":builtin.Mul,							//Mul returns a*b
    "Quo":builtin.Quo,							//Quo returns a/b
    "Mod":builtin.Mod,							//Mod returns a%b
    "Add":builtin.Add,							//Add returns a+b
    "Sub":builtin.Sub,							//Sub returns a-b
    "BitLshr":builtin.BitLshr,					//BitLshr returns a << b
    "BitRshr":builtin.BitRshr,					//BitRshr returns a >> b
    "BitXor":builtin.BitXor,					//BitXor returns a ^ b
    "BitAnd":builtin.BitAnd,					//BitAnd returns a & b
    "BitOr":builtin.BitOr,						//BitOr returns a | b
    "BitNot":builtin.BitNot,					//BitNot returns ^a
    "BitAndNot":builtin.BitAndNot,				//BitAndNot returns a &^ b
    "Or":builtin.Or,							//Or returns 1 || true
    "And":builtin.And,							//And returns 1 && true
    "Not":builtin.Not,							//Not returns !a
    "LT":builtin.LT,							//LT returns a < b
    "GT":builtin.GT,							//GT returns a > b
    "LE":builtin.LE,							//LE returns a <= b
    "GE":builtin.GE,							//GE returns a >= b
    "EQ":builtin.EQ,							//EQ returns a == b
    "NE":builtin.NE,							//NE returns a != b
    "TrySend":builtin.TrySend,					//TrySend(*Chan, value)	不阻塞
    "TryRecv":builtin.TryRecv,					//TryRecv(*Chan, value)	不阻塞
    "Send":builtin.Send,						//Send(*Chan, value)
    "Recv":builtin.Recv,						//Recv(*Chan)
    "Close":builtin.Close,						//Close(*Chan)
    "ChanOf":builtin.ChanOf,					//ChanOf(T)
    "MakeChan":builtin.MakeChan,				//MakeChan(T, size)
    "Max":builtin.Max,							//Max(a1, a2 ...)
    "Min":builtin.Min,							//Min(a1, a2 ...)
    "Error": func(v interface{}) bool {
		return templateFuncMapError(v) != nil
    },
    "NotError": func(v interface{}) bool {
       return templateFuncMapError(v) == nil
    },
}


