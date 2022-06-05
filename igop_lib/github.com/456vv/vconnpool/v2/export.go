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
		NamedTypes: map[string]igop.NamedType{
			"ConnPool": {reflect.TypeOf((*q.ConnPool)(nil)).Elem(), "", "Add,Close,CloseIdleConnections,ConnNum,ConnNumIde,Dial,DialContext,Get,Put,clearPoolConn,dialCtx,getConn,getPoolConn,getPoolConnCount,init,putPoolConn"},
		},
		AliasTypes: map[string]reflect.Type{},
		Vars:       map[string]reflect.Value{},
		Funcs: map[string]reflect.Value{
			"ParseAddr": reflect.ValueOf(q.ParseAddr),
		},
		TypedConsts:   map[string]igop.TypedConst{},
		UntypedConsts: map[string]igop.UntypedConst{},
	})
}
