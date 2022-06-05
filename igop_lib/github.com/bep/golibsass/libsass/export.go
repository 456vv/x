// export by github.com/goplus/igop/cmd/qexp

//go:build igop_lib
// +build igop_lib

package libsass

import (
	q "github.com/bep/golibsass/libsass"

	"go/constant"
	"reflect"

	"github.com/goplus/igop"
)

func init() {
	igop.RegisterPackage(&igop.Package{
		Name: "libsass",
		Path: "github.com/bep/golibsass/libsass",
		Deps: map[string]string{
			"github.com/bep/golibsass/internal/libsass":      "libsass",
			"github.com/bep/golibsass/libsass/libsasserrors": "libsasserrors",
			"os":      "os",
			"strings": "strings",
		},
		Interfaces: map[string]reflect.Type{
			"Transpiler": reflect.TypeOf((*q.Transpiler)(nil)).Elem(),
		},
		NamedTypes: map[string]igop.NamedType{
			"Options":          {reflect.TypeOf((*q.Options)(nil)).Elem(), "", ""},
			"OutputStyle":      {reflect.TypeOf((*q.OutputStyle)(nil)).Elem(), "", ""},
			"Result":           {reflect.TypeOf((*q.Result)(nil)).Elem(), "", ""},
			"SourceMapOptions": {reflect.TypeOf((*q.SourceMapOptions)(nil)).Elem(), "", ""},
		},
		AliasTypes: map[string]reflect.Type{},
		Vars:       map[string]reflect.Value{},
		Funcs: map[string]reflect.Value{
			"New":              reflect.ValueOf(q.New),
			"ParseOutputStyle": reflect.ValueOf(q.ParseOutputStyle),
		},
		TypedConsts: map[string]igop.TypedConst{
			"CompactStyle":    {reflect.TypeOf(q.CompactStyle), constant.MakeInt64(int64(q.CompactStyle))},
			"CompressedStyle": {reflect.TypeOf(q.CompressedStyle), constant.MakeInt64(int64(q.CompressedStyle))},
			"ExpandedStyle":   {reflect.TypeOf(q.ExpandedStyle), constant.MakeInt64(int64(q.ExpandedStyle))},
			"NestedStyle":     {reflect.TypeOf(q.NestedStyle), constant.MakeInt64(int64(q.NestedStyle))},
		},
		UntypedConsts: map[string]igop.UntypedConst{},
	})
}
