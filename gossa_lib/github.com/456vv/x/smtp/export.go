// export by github.com/goplus/gossa/cmd/qexp

//go:build gossa_lib
// +build gossa_lib

package smtp

import (
	q "github.com/456vv/x/smtp"

	"reflect"

	"github.com/goplus/gossa"
)

func init() {
	gossa.RegisterPackage(&gossa.Package{
		Name: "smtp",
		Path: "github.com/456vv/x/smtp",
		Deps: map[string]string{
			"crypto/tls":             "tls",
			"encoding/json":          "json",
			"errors":                 "errors",
			"fmt":                    "fmt",
			"github.com/456vv/vconn": "vconn",
			"net":                    "net",
			"net/smtp":               "smtp",
			"os":                     "os",
			"strings":                "strings",
			"sync":                   "sync",
			"sync/atomic":            "atomic",
		},
		Interfaces: map[string]reflect.Type{},
		NamedTypes: map[string]gossa.NamedType{
			"Info": {reflect.TypeOf((*q.Info)(nil)).Elem(), "", ""},
			"Smtp": {reflect.TypeOf((*q.Smtp)(nil)).Elem(), "", "Close,OpenConfig,Send,connection,notifyConn,reconnection"},
		},
		AliasTypes: map[string]reflect.Type{},
		Vars:       map[string]reflect.Value{},
		Funcs: map[string]reflect.Value{
			"PlainAuth": reflect.ValueOf(q.PlainAuth),
		},
		TypedConsts:   map[string]gossa.TypedConst{},
		UntypedConsts: map[string]gossa.UntypedConst{},
	})
}
