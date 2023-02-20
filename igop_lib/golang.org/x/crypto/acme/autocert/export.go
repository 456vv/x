// export by github.com/goplus/igop/cmd/qexp

//go:build igop_lib
// +build igop_lib

package autocert

import (
	q "golang.org/x/crypto/acme/autocert"

	"go/constant"
	"reflect"

	"github.com/goplus/igop"
)

func init() {
	igop.RegisterPackage(&igop.Package{
		Name: "autocert",
		Path: "golang.org/x/crypto/acme/autocert",
		Deps: map[string]string{
			"bytes":                    "bytes",
			"context":                  "context",
			"crypto":                   "crypto",
			"crypto/ecdsa":             "ecdsa",
			"crypto/elliptic":          "elliptic",
			"crypto/rand":              "rand",
			"crypto/rsa":               "rsa",
			"crypto/tls":               "tls",
			"crypto/x509":              "x509",
			"crypto/x509/pkix":         "pkix",
			"encoding/pem":             "pem",
			"errors":                   "errors",
			"fmt":                      "fmt",
			"golang.org/x/crypto/acme": "acme",
			"golang.org/x/net/idna":    "idna",
			"io":                       "io",
			"log":                      "log",
			"math/rand":                "rand",
			"net":                      "net",
			"net/http":                 "http",
			"os":                       "os",
			"path":                     "path",
			"path/filepath":            "filepath",
			"runtime":                  "runtime",
			"strings":                  "strings",
			"sync":                     "sync",
			"time":                     "time",
		},
		Interfaces: map[string]reflect.Type{
			"Cache": reflect.TypeOf((*q.Cache)(nil)).Elem(),
		},
		NamedTypes: map[string]reflect.Type{
			"DirCache":   reflect.TypeOf((*q.DirCache)(nil)).Elem(),
			"HostPolicy": reflect.TypeOf((*q.HostPolicy)(nil)).Elem(),
			"Manager":    reflect.TypeOf((*q.Manager)(nil)).Elem(),
		},
		AliasTypes: map[string]reflect.Type{},
		Vars: map[string]reflect.Value{
			"ErrCacheMiss": reflect.ValueOf(&q.ErrCacheMiss),
		},
		Funcs: map[string]reflect.Value{
			"AcceptTOS":     reflect.ValueOf(q.AcceptTOS),
			"HostWhitelist": reflect.ValueOf(q.HostWhitelist),
			"NewListener":   reflect.ValueOf(q.NewListener),
		},
		TypedConsts: map[string]igop.TypedConst{},
		UntypedConsts: map[string]igop.UntypedConst{
			"DefaultACMEDirectory": {"untyped string", constant.MakeString(string(q.DefaultACMEDirectory))},
		},
	})
}
