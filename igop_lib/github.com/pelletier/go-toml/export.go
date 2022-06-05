// export by github.com/goplus/igop/cmd/qexp

//go:build igop_lib
// +build igop_lib

package toml

import (
	q "github.com/pelletier/go-toml"

	"go/constant"
	"reflect"

	"github.com/goplus/igop"
)

func init() {
	igop.RegisterPackage(&igop.Package{
		Name: "toml",
		Path: "github.com/pelletier/go-toml",
		Deps: map[string]string{
			"bytes":     "bytes",
			"encoding":  "encoding",
			"errors":    "errors",
			"fmt":       "fmt",
			"io":        "io",
			"io/ioutil": "ioutil",
			"math":      "math",
			"math/big":  "big",
			"os":        "os",
			"reflect":   "reflect",
			"runtime":   "runtime",
			"sort":      "sort",
			"strconv":   "strconv",
			"strings":   "strings",
			"time":      "time",
		},
		Interfaces: map[string]reflect.Type{
			"Marshaler":   reflect.TypeOf((*q.Marshaler)(nil)).Elem(),
			"Unmarshaler": reflect.TypeOf((*q.Unmarshaler)(nil)).Elem(),
		},
		NamedTypes: map[string]igop.NamedType{
			"Decoder":       {reflect.TypeOf((*q.Decoder)(nil)).Elem(), "", "Decode,SetTagName,Strict,unmarshal,unmarshalText,unwrapPointer,valueFromOtherSlice,valueFromOtherSliceI,valueFromToml,valueFromTree,valueFromTreeSlice"},
			"Encoder":       {reflect.TypeOf((*q.Encoder)(nil)).Elem(), "", "ArraysWithOneElementPerLine,CompactComments,Encode,Indentation,Order,PromoteAnonymous,QuoteMapKeys,SetTagComment,SetTagCommented,SetTagMultiline,SetTagName,appendTree,marshal,nextTree,valueToOtherSlice,valueToToml,valueToTree,valueToTreeSlice,wrapTomlValue"},
			"LocalDate":     {reflect.TypeOf((*q.LocalDate)(nil)).Elem(), "AddDays,After,Before,DaysSince,In,IsValid,MarshalText,String", "UnmarshalText"},
			"LocalDateTime": {reflect.TypeOf((*q.LocalDateTime)(nil)).Elem(), "After,Before,In,IsValid,MarshalText,String", "UnmarshalText"},
			"LocalTime":     {reflect.TypeOf((*q.LocalTime)(nil)).Elem(), "IsValid,MarshalText,String", "UnmarshalText"},
			"MarshalOrder":  {reflect.TypeOf((*q.MarshalOrder)(nil)).Elem(), "", ""},
			"Position":      {reflect.TypeOf((*q.Position)(nil)).Elem(), "Invalid,String", ""},
			"SetOptions":    {reflect.TypeOf((*q.SetOptions)(nil)).Elem(), "", ""},
			"Tree":          {reflect.TypeOf((*q.Tree)(nil)).Elem(), "", "Comment,Commented,Delete,DeletePath,Get,GetArray,GetArrayPath,GetDefault,GetPath,GetPosition,GetPositionPath,Has,HasPath,Inline,Keys,Marshal,Position,Set,SetComment,SetCommented,SetInline,SetPath,SetPathWithComment,SetPathWithOptions,SetPositionPath,SetValues,SetWithComment,SetWithOptions,String,ToMap,ToTomlString,Unmarshal,Values,WriteTo,createSubTree,writeTo,writeToOrdered"},
		},
		AliasTypes: map[string]reflect.Type{
			"PubTOMLValue": reflect.TypeOf((*q.PubTOMLValue)(nil)).Elem(),
			"PubTree":      reflect.TypeOf((*q.PubTree)(nil)).Elem(),
		},
		Vars: map[string]reflect.Value{},
		Funcs: map[string]reflect.Value{
			"Load":                      reflect.ValueOf(q.Load),
			"LoadBytes":                 reflect.ValueOf(q.LoadBytes),
			"LoadFile":                  reflect.ValueOf(q.LoadFile),
			"LoadReader":                reflect.ValueOf(q.LoadReader),
			"LocalDateOf":               reflect.ValueOf(q.LocalDateOf),
			"LocalDateTimeOf":           reflect.ValueOf(q.LocalDateTimeOf),
			"LocalTimeOf":               reflect.ValueOf(q.LocalTimeOf),
			"Marshal":                   reflect.ValueOf(q.Marshal),
			"NewDecoder":                reflect.ValueOf(q.NewDecoder),
			"NewEncoder":                reflect.ValueOf(q.NewEncoder),
			"ParseLocalDate":            reflect.ValueOf(q.ParseLocalDate),
			"ParseLocalDateTime":        reflect.ValueOf(q.ParseLocalDateTime),
			"ParseLocalTime":            reflect.ValueOf(q.ParseLocalTime),
			"TreeFromMap":               reflect.ValueOf(q.TreeFromMap),
			"Unmarshal":                 reflect.ValueOf(q.Unmarshal),
			"ValueStringRepresentation": reflect.ValueOf(q.ValueStringRepresentation),
		},
		TypedConsts: map[string]igop.TypedConst{
			"OrderAlphabetical": {reflect.TypeOf(q.OrderAlphabetical), constant.MakeInt64(int64(q.OrderAlphabetical))},
			"OrderPreserve":     {reflect.TypeOf(q.OrderPreserve), constant.MakeInt64(int64(q.OrderPreserve))},
		},
		UntypedConsts: map[string]igop.UntypedConst{},
	})
}
