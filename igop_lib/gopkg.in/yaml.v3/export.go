// export by github.com/goplus/ixgo/cmd/qexp

//go:build igop_lib
// +build igop_lib

package yaml

import (
	q "gopkg.in/yaml.v3"

	"go/constant"
	"reflect"

	"github.com/goplus/ixgo"
)

func init() {
	ixgo.RegisterPackage(&ixgo.Package{
		Name: "yaml",
		Path: "gopkg.in/yaml.v3",
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
		NamedTypes: map[string]reflect.Type{
			"Decoder":   reflect.TypeOf((*q.Decoder)(nil)).Elem(),
			"Encoder":   reflect.TypeOf((*q.Encoder)(nil)).Elem(),
			"Kind":      reflect.TypeOf((*q.Kind)(nil)).Elem(),
			"Node":      reflect.TypeOf((*q.Node)(nil)).Elem(),
			"Style":     reflect.TypeOf((*q.Style)(nil)).Elem(),
			"TypeError": reflect.TypeOf((*q.TypeError)(nil)).Elem(),
		},
		AliasTypes: map[string]reflect.Type{},
		Vars:       map[string]reflect.Value{},
		Funcs: map[string]reflect.Value{
			"Marshal":    reflect.ValueOf(q.Marshal),
			"NewDecoder": reflect.ValueOf(q.NewDecoder),
			"NewEncoder": reflect.ValueOf(q.NewEncoder),
			"Unmarshal":  reflect.ValueOf(q.Unmarshal),
		},
		TypedConsts: map[string]ixgo.TypedConst{
			"AliasNode":         {Typ: reflect.TypeOf(q.AliasNode), Value: constant.MakeInt64(int64(q.AliasNode))},
			"DocumentNode":      {Typ: reflect.TypeOf(q.DocumentNode), Value: constant.MakeInt64(int64(q.DocumentNode))},
			"DoubleQuotedStyle": {Typ: reflect.TypeOf(q.DoubleQuotedStyle), Value: constant.MakeInt64(int64(q.DoubleQuotedStyle))},
			"FlowStyle":         {Typ: reflect.TypeOf(q.FlowStyle), Value: constant.MakeInt64(int64(q.FlowStyle))},
			"FoldedStyle":       {Typ: reflect.TypeOf(q.FoldedStyle), Value: constant.MakeInt64(int64(q.FoldedStyle))},
			"LiteralStyle":      {Typ: reflect.TypeOf(q.LiteralStyle), Value: constant.MakeInt64(int64(q.LiteralStyle))},
			"MappingNode":       {Typ: reflect.TypeOf(q.MappingNode), Value: constant.MakeInt64(int64(q.MappingNode))},
			"ScalarNode":        {Typ: reflect.TypeOf(q.ScalarNode), Value: constant.MakeInt64(int64(q.ScalarNode))},
			"SequenceNode":      {Typ: reflect.TypeOf(q.SequenceNode), Value: constant.MakeInt64(int64(q.SequenceNode))},
			"SingleQuotedStyle": {Typ: reflect.TypeOf(q.SingleQuotedStyle), Value: constant.MakeInt64(int64(q.SingleQuotedStyle))},
			"TaggedStyle":       {Typ: reflect.TypeOf(q.TaggedStyle), Value: constant.MakeInt64(int64(q.TaggedStyle))},
		},
		UntypedConsts: map[string]ixgo.UntypedConst{},
	})
}
