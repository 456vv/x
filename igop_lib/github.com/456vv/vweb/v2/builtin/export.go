// export by github.com/goplus/igop/cmd/qexp

//go:build igop_lib
// +build igop_lib

package builtin

import (
	q "github.com/456vv/vweb/v2/builtin"

	"reflect"

	"github.com/goplus/igop"
)

func init() {
	igop.RegisterPackage(&igop.Package{
		Name: "builtin",
		Path: "github.com/456vv/vweb/v2/builtin",
		Deps: map[string]string{
			"fmt":     "fmt",
			"reflect": "reflect",
			"strconv": "strconv",
			"strings": "strings",
			"unsafe":  "unsafe",
		},
		Interfaces: map[string]reflect.Type{},
		NamedTypes: map[string]reflect.Type{
			"Chan": reflect.TypeOf((*q.Chan)(nil)).Elem(),
		},
		AliasTypes: map[string]reflect.Type{},
		Vars:       map[string]reflect.Value{},
		Funcs: map[string]reflect.Value{
			"Add":        reflect.ValueOf(q.Add),
			"And":        reflect.ValueOf(q.And),
			"Append":     reflect.ValueOf(q.Append),
			"BitAnd":     reflect.ValueOf(q.BitAnd),
			"BitAndNot":  reflect.ValueOf(q.BitAndNot),
			"BitLshr":    reflect.ValueOf(q.BitLshr),
			"BitNot":     reflect.ValueOf(q.BitNot),
			"BitOr":      reflect.ValueOf(q.BitOr),
			"BitRshr":    reflect.ValueOf(q.BitRshr),
			"BitXor":     reflect.ValueOf(q.BitXor),
			"Bool":       reflect.ValueOf(q.Bool),
			"Byte":       reflect.ValueOf(q.Byte),
			"Bytes":      reflect.ValueOf(q.Bytes),
			"Cap":        reflect.ValueOf(q.Cap),
			"ChanOf":     reflect.ValueOf(q.ChanOf),
			"Close":      reflect.ValueOf(q.Close),
			"Complex128": reflect.ValueOf(q.Complex128),
			"Complex64":  reflect.ValueOf(q.Complex64),
			"Compute":    reflect.ValueOf(q.Compute),
			"Convert":    reflect.ValueOf(q.Convert),
			"Copy":       reflect.ValueOf(q.Copy),
			"Dec":        reflect.ValueOf(q.Dec),
			"Delete":     reflect.ValueOf(q.Delete),
			"EQ":         reflect.ValueOf(q.EQ),
			"Float32":    reflect.ValueOf(q.Float32),
			"Float64":    reflect.ValueOf(q.Float64),
			"GE":         reflect.ValueOf(q.GE),
			"GT":         reflect.ValueOf(q.GT),
			"Get":        reflect.ValueOf(q.Get),
			"GetSlice":   reflect.ValueOf(q.GetSlice),
			"GetSlice3":  reflect.ValueOf(q.GetSlice3),
			"Inc":        reflect.ValueOf(q.Inc),
			"Init":       reflect.ValueOf(q.Init),
			"Int":        reflect.ValueOf(q.Int),
			"Int16":      reflect.ValueOf(q.Int16),
			"Int32":      reflect.ValueOf(q.Int32),
			"Int64":      reflect.ValueOf(q.Int64),
			"Int8":       reflect.ValueOf(q.Int8),
			"LE":         reflect.ValueOf(q.LE),
			"LT":         reflect.ValueOf(q.LT),
			"Len":        reflect.ValueOf(q.Len),
			"Make":       reflect.ValueOf(q.Make),
			"MakeChan":   reflect.ValueOf(q.MakeChan),
			"MapFrom":    reflect.ValueOf(q.MapFrom),
			"Max":        reflect.ValueOf(q.Max),
			"Min":        reflect.ValueOf(q.Min),
			"Mod":        reflect.ValueOf(q.Mod),
			"Mul":        reflect.ValueOf(q.Mul),
			"NE":         reflect.ValueOf(q.NE),
			"Neg":        reflect.ValueOf(q.Neg),
			"Not":        reflect.ValueOf(q.Not),
			"Or":         reflect.ValueOf(q.Or),
			"Panic":      reflect.ValueOf(q.Panic),
			"Pointer":    reflect.ValueOf(q.Pointer),
			"Quo":        reflect.ValueOf(q.Quo),
			"Recv":       reflect.ValueOf(q.Recv),
			"Rune":       reflect.ValueOf(q.Rune),
			"Runs":       reflect.ValueOf(q.Runs),
			"Send":       reflect.ValueOf(q.Send),
			"Set":        reflect.ValueOf(q.Set),
			"SliceFrom":  reflect.ValueOf(q.SliceFrom),
			"String":     reflect.ValueOf(q.String),
			"Sub":        reflect.ValueOf(q.Sub),
			"TryRecv":    reflect.ValueOf(q.TryRecv),
			"TrySend":    reflect.ValueOf(q.TrySend),
			"Type":       reflect.ValueOf(q.Type),
			"Uint":       reflect.ValueOf(q.Uint),
			"Uint16":     reflect.ValueOf(q.Uint16),
			"Uint32":     reflect.ValueOf(q.Uint32),
			"Uint64":     reflect.ValueOf(q.Uint64),
			"Uint8":      reflect.ValueOf(q.Uint8),
			"Uintptr":    reflect.ValueOf(q.Uintptr),
			"Value":      reflect.ValueOf(q.Value),
		},
		TypedConsts:   map[string]igop.TypedConst{},
		UntypedConsts: map[string]igop.UntypedConst{},
	})
}
