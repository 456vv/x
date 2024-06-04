// export by github.com/goplus/igop/cmd/qexp

//go:build igop_lib
// +build igop_lib

package server

import (
	q "github.com/456vv/vweb/v2/server"

	"reflect"

	"github.com/goplus/igop"
)

func init() {
	igop.RegisterPackage(&igop.Package{
		Name: "server",
		Path: "github.com/456vv/vweb/v2/server",
		Deps: map[string]string{
			"bytes":                                  "bytes",
			"context":                                "context",
			"crypto/tls":                             "tls",
			"crypto/x509":                            "x509",
			"errors":                                 "errors",
			"fmt":                                    "fmt",
			"github.com/456vv/verror":                "verror",
			"github.com/456vv/vmap/v2":               "vmap",
			"github.com/456vv/vweb/v2":               "vweb",
			"github.com/456vv/vweb/v2/server/config": "config",
			"golang.org/x/crypto/acme/autocert":      "autocert",
			"io":                                     "io",
			"log":                                    "log",
			"mime":                                   "mime",
			"net":                                    "net",
			"net/http":                               "http",
			"os":                                     "os",
			"path":                                   "path",
			"path/filepath":                          "filepath",
			"reflect":                                "reflect",
			"strconv":                                "strconv",
			"strings":                                "strings",
			"sync":                                   "sync",
			"sync/atomic":                            "atomic",
			"time":                                   "time",
		},
		Interfaces: map[string]reflect.Type{},
		NamedTypes: map[string]reflect.Type{
			"Group":  reflect.TypeOf((*q.Group)(nil)).Elem(),
			"Server": reflect.TypeOf((*q.Server)(nil)).Elem(),
		},
		AliasTypes: map[string]reflect.Type{},
		Vars: map[string]reflect.Value{
			"ServerContextKey": reflect.ValueOf(&q.ServerContextKey),
			"Version":          reflect.ValueOf(&q.Version),
		},
		Funcs: map[string]reflect.Value{
			"NewGroup": reflect.ValueOf(q.NewGroup),
		},
		TypedConsts:   map[string]igop.TypedConst{},
		UntypedConsts: map[string]igop.UntypedConst{},
	})
}
