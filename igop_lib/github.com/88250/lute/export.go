// export by github.com/goplus/igop/cmd/qexp

//go:build igop_lib
// +build igop_lib

package lute

import (
	q "github.com/88250/lute"

	"go/constant"
	"reflect"

	"github.com/goplus/igop"
)

func init() {
	igop.RegisterPackage(&igop.Package{
		Name: "lute",
		Path: "github.com/88250/lute",
		Deps: map[string]string{
			"bytes":                           "bytes",
			"errors":                          "errors",
			"github.com/88250/lute/ast":       "ast",
			"github.com/88250/lute/editor":    "editor",
			"github.com/88250/lute/html":      "html",
			"github.com/88250/lute/html/atom": "atom",
			"github.com/88250/lute/lex":       "lex",
			"github.com/88250/lute/parse":     "parse",
			"github.com/88250/lute/render":    "render",
			"github.com/88250/lute/util":      "util",
			"github.com/gopherjs/gopherjs/js": "js",
			"strconv":                         "strconv",
			"strings":                         "strings",
			"sync":                            "sync",
			"unicode":                         "unicode",
		},
		Interfaces: map[string]reflect.Type{},
		NamedTypes: map[string]reflect.Type{
			"Lute":        reflect.TypeOf((*q.Lute)(nil)).Elem(),
			"ParseOption": reflect.TypeOf((*q.ParseOption)(nil)).Elem(),
		},
		AliasTypes: map[string]reflect.Type{},
		Vars:       map[string]reflect.Value{},
		Funcs: map[string]reflect.Value{
			"FormatNodeSync":          reflect.ValueOf(q.FormatNodeSync),
			"New":                     reflect.ValueOf(q.New),
			"ProtyleExportMdNodeSync": reflect.ValueOf(q.ProtyleExportMdNodeSync),
			"RenderNodeBlockDOM":      reflect.ValueOf(q.RenderNodeBlockDOM),
		},
		TypedConsts: map[string]igop.TypedConst{},
		UntypedConsts: map[string]igop.UntypedConst{
			"Version": {"untyped string", constant.MakeString(string(q.Version))},
		},
	})
}
