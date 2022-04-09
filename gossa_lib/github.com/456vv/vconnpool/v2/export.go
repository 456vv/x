// export by github.com/goplus/gossa/cmd/qexp

//go:build gossa_lib
// +build gossa_lib

package vconnpool

import (
	q "github.com/456vv/vconnpool/v2"

	"reflect"

	"github.com/goplus/gossa"
)

func init() {
	gossa.RegisterPackage(&gossa.Package{
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
		NamedTypes: map[string]gossa.NamedType{
			"ConnPool": {reflect.TypeOf((*q.ConnPool)(nil)).Elem(), "", "Add,Close,CloseIdleConnections,ConnNum,ConnNumIde,Dial,DialContext,Get,Put,clearPoolConn,dialCtx,getConn,getPoolConn,getPoolConnCount,init,putPoolConn"},
		},
		AliasTypes: map[string]reflect.Type{},
		Vars:       map[string]reflect.Value{},
		Funcs: map[string]reflect.Value{
			"ParseAddr": reflect.ValueOf(q.ParseAddr),
		},
		TypedConsts:   map[string]gossa.TypedConst{},
		UntypedConsts: map[string]gossa.UntypedConst{},
	})
}
