// export by github.com/goplus/gossa/cmd/qexp

//go:build gossa_lib
// +build gossa_lib

package vforward

import (
	q "github.com/456vv/vforward"

	"go/constant"
	"reflect"

	"github.com/goplus/gossa"
)

func init() {
	gossa.RegisterPackage(&gossa.Package{
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
		NamedTypes: map[string]gossa.NamedType{
			"Addr":    {reflect.TypeOf((*q.Addr)(nil)).Elem(), "", ""},
			"D2D":     {reflect.TypeOf((*q.D2D)(nil)).Elem(), "", "Close,Transport,aGetConn,bGetConn,bufConn,currUseConns,init,logf"},
			"D2DSwap": {reflect.TypeOf((*q.D2DSwap)(nil)).Elem(), "", "Close,ConnNum,Swap"},
			"L2D":     {reflect.TypeOf((*q.L2D)(nil)).Elem(), "", "Close,Transport,logf"},
			"L2DSwap": {reflect.TypeOf((*q.L2DSwap)(nil)).Elem(), "", "Close,ConnNum,Swap,connReadReply,connRemoteTCP,currUseConns,keepAvailable"},
			"L2L":     {reflect.TypeOf((*q.L2L)(nil)).Elem(), "", "Close,Transport,aGetConn,bGetConn,bufConn,currUseConns,init,logf"},
			"L2LSwap": {reflect.TypeOf((*q.L2LSwap)(nil)).Elem(), "", "Close,ConnNum,Swap"},
		},
		AliasTypes: map[string]reflect.Type{},
		Vars:       map[string]reflect.Value{},
		Funcs:      map[string]reflect.Value{},
		TypedConsts: map[string]gossa.TypedConst{
			"DefaultReadBufSize": {reflect.TypeOf(q.DefaultReadBufSize), constant.MakeInt64(int64(q.DefaultReadBufSize))},
		},
		UntypedConsts: map[string]gossa.UntypedConst{},
	})
}
