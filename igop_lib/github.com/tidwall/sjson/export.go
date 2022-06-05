// export by github.com/goplus/igop/cmd/qexp

//go:build igop_lib
// +build igop_lib

package sjson

import (
	q "github.com/tidwall/sjson"

	"reflect"

	"github.com/goplus/igop"
)

func init() {
	igop.RegisterPackage(&igop.Package{
		Name: "sjson",
		Path: "github.com/tidwall/sjson",
		Deps: map[string]string{
			"encoding/json":            "json",
			"github.com/tidwall/gjson": "gjson",
			"sort":                     "sort",
			"strconv":                  "strconv",
			"unsafe":                   "unsafe",
		},
		Interfaces: map[string]reflect.Type{},
		NamedTypes: map[string]igop.NamedType{
			"Options": {reflect.TypeOf((*q.Options)(nil)).Elem(), "", ""},
		},
		AliasTypes: map[string]reflect.Type{},
		Vars:       map[string]reflect.Value{},
		Funcs: map[string]reflect.Value{
			"Delete":             reflect.ValueOf(q.Delete),
			"DeleteBytes":        reflect.ValueOf(q.DeleteBytes),
			"Set":                reflect.ValueOf(q.Set),
			"SetBytes":           reflect.ValueOf(q.SetBytes),
			"SetBytesOptions":    reflect.ValueOf(q.SetBytesOptions),
			"SetOptions":         reflect.ValueOf(q.SetOptions),
			"SetRaw":             reflect.ValueOf(q.SetRaw),
			"SetRawBytes":        reflect.ValueOf(q.SetRawBytes),
			"SetRawBytesOptions": reflect.ValueOf(q.SetRawBytesOptions),
			"SetRawOptions":      reflect.ValueOf(q.SetRawOptions),
		},
		TypedConsts:   map[string]igop.TypedConst{},
		UntypedConsts: map[string]igop.UntypedConst{},
	})
}
