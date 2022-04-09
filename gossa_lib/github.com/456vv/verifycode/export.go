// export by github.com/goplus/gossa/cmd/qexp

//go:build gossa_lib
// +build gossa_lib

package verifycode

import (
	q "github.com/456vv/verifycode"

	"reflect"

	"github.com/goplus/gossa"
)

func init() {
	gossa.RegisterPackage(&gossa.Package{
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
		NamedTypes: map[string]gossa.NamedType{
			"Color":      {reflect.TypeOf((*q.Color)(nil)).Elem(), "", "AddHEX,AddRGBA,Random"},
			"Font":       {reflect.TypeOf((*q.Font)(nil)).Elem(), "", "AddFile,Random"},
			"Glyph":      {reflect.TypeOf((*q.Glyph)(nil)).Elem(), "", "FontGlyph"},
			"Style":      {reflect.TypeOf((*q.Style)(nil)).Elem(), "", ""},
			"VerifyCode": {reflect.TypeOf((*q.VerifyCode)(nil)).Elem(), "", "Draw,GIF,JPEG,PNG,Style"},
		},
		AliasTypes: map[string]reflect.Type{},
		Vars:       map[string]reflect.Value{},
		Funcs: map[string]reflect.Value{
			"Rand":       reflect.ValueOf(q.Rand),
			"RandRange":  reflect.ValueOf(q.RandRange),
			"RandomText": reflect.ValueOf(q.RandomText),
		},
		TypedConsts:   map[string]gossa.TypedConst{},
		UntypedConsts: map[string]gossa.UntypedConst{},
	})
}
