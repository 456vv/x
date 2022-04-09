// export by github.com/goplus/gossa/cmd/qexp

//go:build gossa_lib
// +build gossa_lib

package sjson

import (
	q "github.com/tidwall/sjson"

	"reflect"

	"github.com/goplus/gossa"
)

func init() {
	gossa.RegisterPackage(&gossa.Package{
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
		NamedTypes: map[string]gossa.NamedType{
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
		TypedConsts:   map[string]gossa.TypedConst{},
		UntypedConsts: map[string]gossa.UntypedConst{},
	})
}
