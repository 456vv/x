// export by github.com/goplus/ixgo/cmd/qexp

//go:build igop_lib
// +build igop_lib

package vconnpool

import (
	q "github.com/456vv/vconnpool/v3"

	"reflect"

	"github.com/goplus/ixgo"
)

func init() {
	ixgo.RegisterPackage(&ixgo.Package{
		Name: "vconnpool",
		Path: "github.com/456vv/vconnpool/v3",
		Deps: map[string]string{
			"context":                "context",
			"errors":                 "errors",
			"fmt":                    "fmt",
			"github.com/456vv/vconn": "vconn",
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
			"Addr": reflect.TypeOf((*q.Addr)(nil)).Elem(),
			"Pool": reflect.TypeOf((*q.Pool)(nil)).Elem(),
		},
		AliasTypes: map[string]reflect.Type{},
		Vars: map[string]reflect.Value{
			"ErrConnAlreadyExists": reflect.ValueOf(&q.ErrConnAlreadyExists),
			"ErrConnClose":         reflect.ValueOf(&q.ErrConnClose),
			"ErrConnIdleMax":       reflect.ValueOf(&q.ErrConnIdleMax),
			"ErrConnNotAvailable":  reflect.ValueOf(&q.ErrConnNotAvailable),
			"ErrConnPoolClosed":    reflect.ValueOf(&q.ErrConnPoolClosed),
			"ErrConnPoolMax":       reflect.ValueOf(&q.ErrConnPoolMax),
			"ErrConnRAWRead":       reflect.ValueOf(&q.ErrConnRAWRead),
			"PriorityContextKey":   reflect.ValueOf(&q.PriorityContextKey),
		},
		Funcs: map[string]reflect.Value{
			"ResolveAddr": reflect.ValueOf(q.ResolveAddr),
		},
		TypedConsts:   map[string]ixgo.TypedConst{},
		UntypedConsts: map[string]ixgo.UntypedConst{},
	})
}
