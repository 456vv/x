// export by github.com/goplus/igop/cmd/qexp

//go:build igop_lib
// +build igop_lib

package vforward

import (
	q "github.com/456vv/vforward"

	"go/constant"
	"reflect"

	"github.com/goplus/igop"
)

func init() {
	igop.RegisterPackage(&igop.Package{
		Name: "vforward",
		Path: "github.com/456vv/vforward",
		Deps: map[string]string{
			"context":                       "context",
			"errors":                        "errors",
			"fmt":                           "fmt",
			"github.com/456vv/vconnpool/v2": "vconnpool",
			"github.com/456vv/vmap/v2":      "vmap",
			"io":                            "io",
			"log":                           "log",
			"net":                           "net",
			"strings":                       "strings",
			"sync/atomic":                   "atomic",
			"time":                          "time",
		},
		Interfaces: map[string]reflect.Type{},
		NamedTypes: map[string]reflect.Type{
			"Addr":    reflect.TypeOf((*q.Addr)(nil)).Elem(),
			"D2D":     reflect.TypeOf((*q.D2D)(nil)).Elem(),
			"D2DSwap": reflect.TypeOf((*q.D2DSwap)(nil)).Elem(),
			"L2D":     reflect.TypeOf((*q.L2D)(nil)).Elem(),
			"L2DSwap": reflect.TypeOf((*q.L2DSwap)(nil)).Elem(),
			"L2L":     reflect.TypeOf((*q.L2L)(nil)).Elem(),
			"L2LSwap": reflect.TypeOf((*q.L2LSwap)(nil)).Elem(),
		},
		AliasTypes: map[string]reflect.Type{},
		Vars:       map[string]reflect.Value{},
		Funcs:      map[string]reflect.Value{},
		TypedConsts: map[string]igop.TypedConst{
			"DefaultReadBufSize": {reflect.TypeOf(q.DefaultReadBufSize), constant.MakeInt64(int64(q.DefaultReadBufSize))},
		},
		UntypedConsts: map[string]igop.UntypedConst{},
	})
}
