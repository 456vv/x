// export by github.com/goplus/ixgo/cmd/qexp

//go:build igop_lib
// +build igop_lib

package vforward

import (
	q "github.com/456vv/vforward/v2"

	"go/constant"
	"reflect"

	"github.com/goplus/ixgo"
)

func init() {
	ixgo.RegisterPackage(&ixgo.Package{
		Name: "vforward",
		Path: "github.com/456vv/vforward/v2",
		Deps: map[string]string{
			"context":                        "context",
			"errors":                         "errors",
			"github.com/456vv/vconn":         "vconn",
			"github.com/456vv/vconnpool/v3":  "vconnpool",
			"github.com/libp2p/go-reuseport": "reuseport",
			"io":                             "io",
			"net":                            "net",
			"strings":                        "strings",
			"sync":                           "sync",
			"sync/atomic":                    "atomic",
			"time":                           "time",
		},
		Interfaces: map[string]reflect.Type{},
		NamedTypes: map[string]reflect.Type{
			"D2D":     reflect.TypeOf((*q.D2D)(nil)).Elem(),
			"D2DSwap": reflect.TypeOf((*q.D2DSwap)(nil)).Elem(),
			"L2D":     reflect.TypeOf((*q.L2D)(nil)).Elem(),
			"L2DSwap": reflect.TypeOf((*q.L2DSwap)(nil)).Elem(),
			"L2L":     reflect.TypeOf((*q.L2L)(nil)).Elem(),
			"L2LSwap": reflect.TypeOf((*q.L2LSwap)(nil)).Elem(),
			"Option":  reflect.TypeOf((*q.Option)(nil)).Elem(),
		},
		AliasTypes:  map[string]reflect.Type{},
		Vars:        map[string]reflect.Value{},
		Funcs:       map[string]reflect.Value{},
		TypedConsts: map[string]ixgo.TypedConst{},
		UntypedConsts: map[string]ixgo.UntypedConst{
			"DefaultReadBufSize": {Typ: "untyped int", Value: constant.MakeInt64(int64(q.DefaultReadBufSize))},
		},
	})
}
