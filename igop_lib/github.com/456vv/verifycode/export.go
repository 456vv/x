// export by github.com/goplus/igop/cmd/qexp

//go:build igop_lib
// +build igop_lib

package verifycode

import (
	q "github.com/456vv/verifycode"

	"reflect"

	"github.com/goplus/igop"
)

func init() {
	igop.RegisterPackage(&igop.Package{
		Name: "verifycode",
		Path: "github.com/456vv/verifycode",
		Deps: map[string]string{
			"crypto/rand":                         "rand",
			"fmt":                                 "fmt",
			"github.com/golang/freetype":          "freetype",
			"github.com/golang/freetype/truetype": "truetype",
			"golang.org/x/image/font":             "font",
			"golang.org/x/image/math/fixed":       "fixed",
			"image":                               "image",
			"image/color":                         "color",
			"image/draw":                          "draw",
			"image/gif":                           "gif",
			"image/jpeg":                          "jpeg",
			"image/png":                           "png",
			"io":                                  "io",
			"io/ioutil":                           "ioutil",
			"math/big":                            "big",
			"strconv":                             "strconv",
			"strings":                             "strings",
		},
		Interfaces: map[string]reflect.Type{},
		NamedTypes: map[string]reflect.Type{
			"Color":      reflect.TypeOf((*q.Color)(nil)).Elem(),
			"Font":       reflect.TypeOf((*q.Font)(nil)).Elem(),
			"Glyph":      reflect.TypeOf((*q.Glyph)(nil)).Elem(),
			"Style":      reflect.TypeOf((*q.Style)(nil)).Elem(),
			"VerifyCode": reflect.TypeOf((*q.VerifyCode)(nil)).Elem(),
		},
		AliasTypes: map[string]reflect.Type{},
		Vars:       map[string]reflect.Value{},
		Funcs: map[string]reflect.Value{
			"Rand":       reflect.ValueOf(q.Rand),
			"RandRange":  reflect.ValueOf(q.RandRange),
			"RandomText": reflect.ValueOf(q.RandomText),
		},
		TypedConsts:   map[string]igop.TypedConst{},
		UntypedConsts: map[string]igop.UntypedConst{},
	})
}
