// export by github.com/goplus/igop/cmd/qexp

//go:build igop_lib
// +build igop_lib

package vconnpool

import (
	q "github.com/456vv/vconnpool/v2"

	"reflect"

	"github.com/goplus/igop"
)

func init() {
	igop.RegisterPackage(&igop.Package{
		Name: "vconnpool",
		Path: "github.com/456vv/vconnpool/v2",
		Deps: map[string]string{
			"context":                "context",
			"errors":                 "errors",
			"fmt":                    "fmt",
			"github.com/456vv/vconn": "vconn",
			"io":                     "io",
			"net":                    "net",
			"sync":                   "sync",
			"sync/atomic":            "atomic",
			"time":                   "time",
		},
		Interfaces: map[string]reflect.Type{
			"Conn":   reflect.TypeOf((*q.Conn)(nil)).Elem(),
			"Dialer": reflect.TypeOf((*q.Dialer)(nil)).Elem(),
		},
		NamedTypes: map[string]reflect.Type{
			"ConnPool": reflect.TypeOf((*q.ConnPool)(nil)).Elem(),
		},
		AliasTypes: map[string]reflect.Type{},
		Vars: map[string]reflect.Value{
			"ErrConnNotAvailable": reflect.ValueOf(&q.ErrConnNotAvailable),
			"ErrConnPoolMax":      reflect.ValueOf(&q.ErrConnPoolMax),
			"ErrPoolFull":         reflect.ValueOf(&q.ErrPoolFull),
			"PriorityContextKey":  reflect.ValueOf(&q.PriorityContextKey),
		},
		Funcs: map[string]reflect.Value{
			"ResolveAddr": reflect.ValueOf(q.ResolveAddr),
		},
		TypedConsts:   map[string]igop.TypedConst{},
		UntypedConsts: map[string]igop.UntypedConst{},
	})
}
