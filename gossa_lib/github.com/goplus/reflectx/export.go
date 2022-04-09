// export by github.com/goplus/gossa/cmd/qexp

//go:build gossa_lib
// +build gossa_lib

package reflectx

import (
	q "github.com/goplus/reflectx"

	"go/constant"
	"reflect"

	"github.com/goplus/gossa"
)

func init() {
	gossa.RegisterPackage(&gossa.Package{
		Name: "reflectx",
		Path: "github.com/goplus/reflectx",
		Deps: map[string]string{
			"fmt":          "fmt",
			"go/token":     "token",
			"io":           "io",
			"log":          "log",
			"path":         "path",
			"reflect":      "reflect",
			"sort":         "sort",
			"strconv":      "strconv",
			"strings":      "strings",
			"unicode":      "unicode",
			"unicode/utf8": "utf8",
			"unsafe":       "unsafe",
		},
		Interfaces: map[string]reflect.Type{
			"MethodProvider": reflect.TypeOf((*q.MethodProvider)(nil)).Elem(),
		},
		NamedTypes: map[string]gossa.NamedType{
			"ChanDir":    {reflect.TypeOf((*q.ChanDir)(nil)).Elem(), "", ""},
			"Method":     {reflect.TypeOf((*q.Method)(nil)).Elem(), "", ""},
			"MethodInfo": {reflect.TypeOf((*q.MethodInfo)(nil)).Elem(), "", ""},
			"Named":      {reflect.TypeOf((*q.Named)(nil)).Elem(), "", ""},
			"TypeKind":   {reflect.TypeOf((*q.TypeKind)(nil)).Elem(), "", ""},
			"Value":      {reflect.TypeOf((*q.Value)(nil)).Elem(), "", ""},
		},
		AliasTypes: map[string]reflect.Type{},
		Vars: map[string]reflect.Value{
			"EnableExportAllMethod":        reflect.ValueOf(&q.EnableExportAllMethod),
			"EnableStructOfExportAllField": reflect.ValueOf(&q.EnableStructOfExportAllField),
		},
		Funcs: map[string]reflect.Value{
			"AddMethodProvider":  reflect.ValueOf(q.AddMethodProvider),
			"CanSet":             reflect.ValueOf(q.CanSet),
			"DumpType":           reflect.ValueOf(q.DumpType),
			"Field":              reflect.ValueOf(q.Field),
			"FieldByIndex":       reflect.ValueOf(q.FieldByIndex),
			"FieldByIndexX":      reflect.ValueOf(q.FieldByIndexX),
			"FieldByName":        reflect.ValueOf(q.FieldByName),
			"FieldByNameFunc":    reflect.ValueOf(q.FieldByNameFunc),
			"FieldByNameFuncX":   reflect.ValueOf(q.FieldByNameFuncX),
			"FieldByNameX":       reflect.ValueOf(q.FieldByNameX),
			"FieldX":             reflect.ValueOf(q.FieldX),
			"InterfaceOf":        reflect.ValueOf(q.InterfaceOf),
			"IsNamed":            reflect.ValueOf(q.IsNamed),
			"MakeEmptyInterface": reflect.ValueOf(q.MakeEmptyInterface),
			"MakeMethod":         reflect.ValueOf(q.MakeMethod),
			"MethodByIndex":      reflect.ValueOf(q.MethodByIndex),
			"MethodByName":       reflect.ValueOf(q.MethodByName),
			"MethodX":            reflect.ValueOf(q.MethodX),
			"NamedInterfaceOf":   reflect.ValueOf(q.NamedInterfaceOf),
			"NamedStructOf":      reflect.ValueOf(q.NamedStructOf),
			"NamedTypeOf":        reflect.ValueOf(q.NamedTypeOf),
			"NewInterfaceType":   reflect.ValueOf(q.NewInterfaceType),
			"NewMethodSet":       reflect.ValueOf(q.NewMethodSet),
			"NumMethodX":         reflect.ValueOf(q.NumMethodX),
			"ReplaceType":        reflect.ValueOf(q.ReplaceType),
			"Reset":              reflect.ValueOf(q.Reset),
			"SetElem":            reflect.ValueOf(q.SetElem),
			"SetInterfaceType":   reflect.ValueOf(q.SetInterfaceType),
			"SetMethodSet":       reflect.ValueOf(q.SetMethodSet),
			"SetUnderlying":      reflect.ValueOf(q.SetUnderlying),
			"SetValue":           reflect.ValueOf(q.SetValue),
			"StructOf":           reflect.ValueOf(q.StructOf),
			"StructToMethodSet":  reflect.ValueOf(q.StructToMethodSet),
			"ToNamed":            reflect.ValueOf(q.ToNamed),
			"TypeLinks":          reflect.ValueOf(q.TypeLinks),
			"TypesByString":      reflect.ValueOf(q.TypesByString),
			"UpdateField":        reflect.ValueOf(q.UpdateField),
		},
		TypedConsts: map[string]gossa.TypedConst{
			"BothDir":   {reflect.TypeOf(q.BothDir), constant.MakeInt64(int64(q.BothDir))},
			"RecvDir":   {reflect.TypeOf(q.RecvDir), constant.MakeInt64(int64(q.RecvDir))},
			"SendDir":   {reflect.TypeOf(q.SendDir), constant.MakeInt64(int64(q.SendDir))},
			"TkInvalid": {reflect.TypeOf(q.TkInvalid), constant.MakeInt64(int64(q.TkInvalid))},
			"TkMethod":  {reflect.TypeOf(q.TkMethod), constant.MakeInt64(int64(q.TkMethod))},
			"TkType":    {reflect.TypeOf(q.TkType), constant.MakeInt64(int64(q.TkType))},
		},
		UntypedConsts: map[string]gossa.UntypedConst{},
	})
}
