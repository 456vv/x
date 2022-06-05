// export by github.com/goplus/igop/cmd/qexp

//go:build igop_lib
// +build igop_lib

package yaml

import (
	q "gopkg.in/yaml.v2"

	"reflect"

	"github.com/goplus/igop"
)

func init() {
	igop.RegisterPackage(&igop.Package{
		Name: "yaml",
		Path: "gopkg.in/yaml.v2",
		Deps: map[string]string{
			"bytes":           "bytes",
			"encoding":        "encoding",
			"encoding/base64": "base64",
			"errors":          "errors",
			"fmt":             "fmt",
			"io":              "io",
			"math":            "math",
			"reflect":         "reflect",
			"regexp":          "regexp",
			"sort":            "sort",
			"strconv":         "strconv",
			"strings":         "strings",
			"sync":            "sync",
			"time":            "time",
			"unicode":         "unicode",
			"unicode/utf8":    "utf8",
		},
		Interfaces: map[string]reflect.Type{
			"IsZeroer":    reflect.TypeOf((*q.IsZeroer)(nil)).Elem(),
			"Marshaler":   reflect.TypeOf((*q.Marshaler)(nil)).Elem(),
			"Unmarshaler": reflect.TypeOf((*q.Unmarshaler)(nil)).Elem(),
		},
		NamedTypes: map[string]igop.NamedType{
			"Decoder":   {reflect.TypeOf((*q.Decoder)(nil)).Elem(), "", "Decode,SetStrict"},
			"Encoder":   {reflect.TypeOf((*q.Encoder)(nil)).Elem(), "", "Close,Encode"},
			"MapItem":   {reflect.TypeOf((*q.MapItem)(nil)).Elem(), "", ""},
			"MapSlice":  {reflect.TypeOf((*q.MapSlice)(nil)).Elem(), "", ""},
			"TypeError": {reflect.TypeOf((*q.TypeError)(nil)).Elem(), "", "Error"},
		},
		AliasTypes: map[string]reflect.Type{},
		Vars:       map[string]reflect.Value{},
		Funcs: map[string]reflect.Value{
			"FutureLineWrap":  reflect.ValueOf(q.FutureLineWrap),
			"Marshal":         reflect.ValueOf(q.Marshal),
			"NewDecoder":      reflect.ValueOf(q.NewDecoder),
			"NewEncoder":      reflect.ValueOf(q.NewEncoder),
			"Unmarshal":       reflect.ValueOf(q.Unmarshal),
			"UnmarshalStrict": reflect.ValueOf(q.UnmarshalStrict),
		},
		TypedConsts:   map[string]igop.TypedConst{},
		UntypedConsts: map[string]igop.UntypedConst{},
	})
}
