// export by github.com/goplus/gossa/cmd/qexp

//go:build gossa_lib
// +build gossa_lib

package gjson

import (
	q "github.com/tidwall/gjson"

	"go/constant"
	"reflect"

	"github.com/goplus/gossa"
)

func init() {
	gossa.RegisterPackage(&gossa.Package{
		Name: "gjson",
		Path: "github.com/tidwall/gjson",
		Deps: map[string]string{
			"encoding/json":             "json",
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
		NamedTypes: map[string]gossa.NamedType{
			"Result": {reflect.TypeOf((*q.Result)(nil)).Elem(), "Array,Bool,Exists,Float,ForEach,Get,Int,IsArray,IsBool,IsObject,Less,Map,Path,Paths,String,Time,Uint,Value,arrayOrMap", ""},
			"Type":   {reflect.TypeOf((*q.Type)(nil)).Elem(), "String", ""},
		},
		AliasTypes: map[string]reflect.Type{},
		Vars: map[string]reflect.Value{
			"DisableModifiers": reflect.ValueOf(&q.DisableModifiers),
		},
		Funcs: map[string]reflect.Value{
			"AddModifier":    reflect.ValueOf(q.AddModifier),
			"ForEachLine":    reflect.ValueOf(q.ForEachLine),
			"Get":            reflect.ValueOf(q.Get),
			"GetBytes":       reflect.ValueOf(q.GetBytes),
			"GetMany":        reflect.ValueOf(q.GetMany),
			"GetManyBytes":   reflect.ValueOf(q.GetManyBytes),
			"ModifierExists": reflect.ValueOf(q.ModifierExists),
			"Parse":          reflect.ValueOf(q.Parse),
			"ParseBytes":     reflect.ValueOf(q.ParseBytes),
			"Valid":          reflect.ValueOf(q.Valid),
			"ValidBytes":     reflect.ValueOf(q.ValidBytes),
		},
		TypedConsts: map[string]gossa.TypedConst{
			"False":  {reflect.TypeOf(q.False), constant.MakeInt64(int64(q.False))},
			"JSON":   {reflect.TypeOf(q.JSON), constant.MakeInt64(int64(q.JSON))},
			"Null":   {reflect.TypeOf(q.Null), constant.MakeInt64(int64(q.Null))},
			"Number": {reflect.TypeOf(q.Number), constant.MakeInt64(int64(q.Number))},
			"String": {reflect.TypeOf(q.String), constant.MakeInt64(int64(q.String))},
			"True":   {reflect.TypeOf(q.True), constant.MakeInt64(int64(q.True))},
		},
		UntypedConsts: map[string]gossa.UntypedConst{},
	})
}
