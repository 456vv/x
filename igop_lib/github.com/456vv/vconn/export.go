// export by github.com/goplus/igop/cmd/qexp

//go:build igop_lib
// +build igop_lib

package vconn

import (
	q "github.com/456vv/vconn"

	"reflect"

	"github.com/goplus/igop"
)

func init() {
	igop.RegisterPackage(&igop.Package{
		Name: "vconn",
		Path: "github.com/456vv/vconn",
		Deps: map[string]string{
			"errors":      "errors",
			"io":          "io",
			"math":        "math",
			"net":         "net",
			"os":          "os",
			"sync":        "sync",
			"sync/atomic": "atomic",
			"syscall":     "syscall",
			"time":        "time",
		},
		Interfaces: map[string]reflect.Type{
			"CloseNotifier": reflect.TypeOf((*q.CloseNotifier)(nil)).Elem(),
		},
		NamedTypes: map[string]igop.NamedType{
			"Conn": {reflect.TypeOf((*q.Conn)(nil)).Elem(), "", "Close,CloseNotify,DisableBackgroundRead,LocalAddr,RawConn,Read,RemoteAddr,SetBackgroundReadDiscard,SetDeadline,SetReadDeadline,SetReadLimit,SetWriteDeadline,Write,close,closeNotify"},
		},
		AliasTypes: map[string]reflect.Type{},
		Vars:       map[string]reflect.Value{},
		Funcs: map[string]reflect.Value{
			"New":     reflect.ValueOf(q.New),
			"NewConn": reflect.ValueOf(q.NewConn),
		},
		TypedConsts:   map[string]igop.TypedConst{},
		UntypedConsts: map[string]igop.UntypedConst{},
	})
}
