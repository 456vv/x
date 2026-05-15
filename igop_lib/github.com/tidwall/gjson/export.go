// export by github.com/goplus/ixgo/cmd/qexp

//go:build igop_lib
// +build igop_lib

package gjson

import (
	q "github.com/tidwall/gjson"

	"go/constant"
	"reflect"

	"github.com/goplus/ixgo"
)

func init() {
	ixgo.RegisterPackage(&ixgo.Package{
		Name: "gjson",
		Path: "github.com/tidwall/gjson",
		Deps: map[string]string{
			"github.com/tidwall/match":  "match",
			"github.com/tidwall/pretty": "pretty",
			"iter":                      "iter",
			"strconv":                   "strconv",
			"strings":                   "strings",
			"time":                      "time",
			"unicode/utf16":             "utf16",
			"unicode/utf8":              "utf8",
			"unsafe":                    "unsafe",
		},
		Interfaces: map[string]reflect.Type{},
		NamedTypes: map[string]reflect.Type{
			"Result": reflect.TypeOf((*q.Result)(nil)).Elem(),
			"Type":   reflect.TypeOf((*q.Type)(nil)).Elem(),
		},
		AliasTypes: map[string]reflect.Type{},
		Vars: map[string]reflect.Value{
			"DisableEscapeHTML": reflect.ValueOf(&q.DisableEscapeHTML),
			"DisableModifiers":  reflect.ValueOf(&q.DisableModifiers),
		},
		Funcs: map[string]reflect.Value{
			"AddModifier":      reflect.ValueOf(q.AddModifier),
			"AppendJSONString": reflect.ValueOf(q.AppendJSONString),
			"Escape":           reflect.ValueOf(q.Escape),
			"ForEachLine":      reflect.ValueOf(q.ForEachLine),
			"Get":              reflect.ValueOf(q.Get),
			"GetBytes":         reflect.ValueOf(q.GetBytes),
			"GetMany":          reflect.ValueOf(q.GetMany),
			"GetManyBytes":     reflect.ValueOf(q.GetManyBytes),
			"ModifierExists":   reflect.ValueOf(q.ModifierExists),
			"Parse":            reflect.ValueOf(q.Parse),
			"ParseBytes":       reflect.ValueOf(q.ParseBytes),
			"Valid":            reflect.ValueOf(q.Valid),
			"ValidBytes":       reflect.ValueOf(q.ValidBytes),
		},
		TypedConsts: map[string]ixgo.TypedConst{
			"False":  {Typ: reflect.TypeOf(q.False), Value: constant.MakeInt64(int64(q.False))},
			"JSON":   {Typ: reflect.TypeOf(q.JSON), Value: constant.MakeInt64(int64(q.JSON))},
			"Null":   {Typ: reflect.TypeOf(q.Null), Value: constant.MakeInt64(int64(q.Null))},
			"Number": {Typ: reflect.TypeOf(q.Number), Value: constant.MakeInt64(int64(q.Number))},
			"String": {Typ: reflect.TypeOf(q.String), Value: constant.MakeInt64(int64(q.String))},
			"True":   {Typ: reflect.TypeOf(q.True), Value: constant.MakeInt64(int64(q.True))},
		},
		UntypedConsts: map[string]ixgo.UntypedConst{},
	})
}
