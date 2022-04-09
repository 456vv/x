// export by github.com/goplus/gossa/cmd/qexp

//go:build gossa_lib
// +build gossa_lib

package vconn

import (
	q "github.com/456vv/vconn"

	"reflect"

	"github.com/goplus/gossa"
)

func init() {
	gossa.RegisterPackage(&gossa.Package{
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
		NamedTypes: map[string]gossa.NamedType{
			"Conn": {reflect.TypeOf((*q.Conn)(nil)).Elem(), "", "Close,CloseNotify,DisableBackgroundRead,LocalAddr,RawConn,Read,RemoteAddr,SetBackgroundReadDiscard,SetDeadline,SetReadDeadline,SetReadLimit,SetWriteDeadline,Write,close,closeNotify"},
		},
		AliasTypes: map[string]reflect.Type{},
		Vars:       map[string]reflect.Value{},
		Funcs: map[string]reflect.Value{
			"New":     reflect.ValueOf(q.New),
			"NewConn": reflect.ValueOf(q.NewConn),
		},
		TypedConsts:   map[string]gossa.TypedConst{},
		UntypedConsts: map[string]gossa.UntypedConst{},
	})
}
