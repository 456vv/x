// export by github.com/goplus/igop/cmd/qexp

//go:build igop_lib
// +build igop_lib

package gjson

import (
	q "github.com/tidwall/gjson"

	"go/constant"
	"reflect"

	"github.com/goplus/igop"
)

func init() {
	igop.RegisterPackage(&igop.Package{
		Name: "gjson",
		Path: "github.com/tidwall/gjson",
		Deps: map[string]string{
			"github.com/tidwall/match":  "match",
			"github.com/tidwall/pretty": "pretty",
			"strconv":                   "strconv",
			"strings":                   "strings",
			"time":                      "time",
			"unicode/utf16":             "utf16",
			"unicode/utf8":              "utf8",
			"unsafe":                    "unsafe",
		},
		Interfaces: map[string]reflect.Type{},
		NamedTypes: map[string]igop.NamedType{
			"Result": {reflect.TypeOf((*q.Result)(nil)).Elem(), "Array,Bool,Exists,Float,ForEach,Get,Int,IsArray,IsBool,IsObject,Less,Map,Path,Paths,String,Time,Uint,Value,arrayOrMap", ""},
			"Type":   {reflect.TypeOf((*q.Type)(nil)).Elem(), "String", ""},
		},
		AliasTypes: map[string]reflect.Type{},
		Vars: map[string]reflect.Value{
			"DisableModifiers": reflect.ValueOf(&q.DisableModifiers),
		},
		Funcs: map[string]reflect.Value{
			"AddModifier":      reflect.ValueOf(q.AddModifier),
			"AppendJSONString": reflect.ValueOf(q.AppendJSONString),
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
		TypedConsts: map[string]igop.TypedConst{
			"False":  {reflect.TypeOf(q.False), constant.MakeInt64(int64(q.False))},
			"JSON":   {reflect.TypeOf(q.JSON), constant.MakeInt64(int64(q.JSON))},
			"Null":   {reflect.TypeOf(q.Null), constant.MakeInt64(int64(q.Null))},
			"Number": {reflect.TypeOf(q.Number), constant.MakeInt64(int64(q.Number))},
			"String": {reflect.TypeOf(q.String), constant.MakeInt64(int64(q.String))},
			"True":   {reflect.TypeOf(q.True), constant.MakeInt64(int64(q.True))},
		},
		UntypedConsts: map[string]igop.UntypedConst{},
	})
}
