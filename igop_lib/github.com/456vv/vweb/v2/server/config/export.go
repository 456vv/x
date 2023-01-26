// export by github.com/goplus/igop/cmd/qexp

//go:build igop_lib
// +build igop_lib

package config

import (
	q "github.com/456vv/vweb/v2/server/config"

	"go/constant"
	"reflect"

	"github.com/goplus/igop"
)

func init() {
	igop.RegisterPackage(&igop.Package{
		Name: "config",
		Path: "github.com/456vv/vweb/v2/server/config",
		Deps: map[string]string{
			"crypto/tls":                    "tls",
			"crypto/x509":                   "x509",
			"encoding/json":                 "json",
			"github.com/456vv/vconnpool/v2": "vconnpool",
			"github.com/456vv/verror":       "verror",
			"github.com/456vv/vweb/v2":      "vweb",
			"io":                            "io",
			"io/ioutil":                     "ioutil",
			"net":                           "net",
			"net/http":                      "http",
			"net/url":                       "url",
			"os":                            "os",
			"path/filepath":                 "filepath",
			"reflect":                       "reflect",
			"strings":                       "strings",
			"time":                          "time",
		},
		Interfaces: map[string]reflect.Type{},
		NamedTypes: map[string]reflect.Type{
			"Config":         reflect.TypeOf((*q.Config)(nil)).Elem(),
			"Conn":           reflect.TypeOf((*q.Conn)(nil)).Elem(),
			"Listen":         reflect.TypeOf((*q.Listen)(nil)).Elem(),
			"Server":         reflect.TypeOf((*q.Server)(nil)).Elem(),
			"ServerPublic":   reflect.TypeOf((*q.ServerPublic)(nil)).Elem(),
			"ServerTLS":      reflect.TypeOf((*q.ServerTLS)(nil)).Elem(),
			"ServerTLSFile":  reflect.TypeOf((*q.ServerTLSFile)(nil)).Elem(),
			"Servers":        reflect.TypeOf((*q.Servers)(nil)).Elem(),
			"Site":           reflect.TypeOf((*q.Site)(nil)).Elem(),
			"SiteDirectory":  reflect.TypeOf((*q.SiteDirectory)(nil)).Elem(),
			"SiteDynamic":    reflect.TypeOf((*q.SiteDynamic)(nil)).Elem(),
			"SiteForward":    reflect.TypeOf((*q.SiteForward)(nil)).Elem(),
			"SiteForwards":   reflect.TypeOf((*q.SiteForwards)(nil)).Elem(),
			"SiteHeader":     reflect.TypeOf((*q.SiteHeader)(nil)).Elem(),
			"SiteHeaderType": reflect.TypeOf((*q.SiteHeaderType)(nil)).Elem(),
			"SiteLog":        reflect.TypeOf((*q.SiteLog)(nil)).Elem(),
			"SiteLogLevel":   reflect.TypeOf((*q.SiteLogLevel)(nil)).Elem(),
			"SitePlugin":     reflect.TypeOf((*q.SitePlugin)(nil)).Elem(),
			"SitePluginTLS":  reflect.TypeOf((*q.SitePluginTLS)(nil)).Elem(),
			"SitePlugins":    reflect.TypeOf((*q.SitePlugins)(nil)).Elem(),
			"SiteProperty":   reflect.TypeOf((*q.SiteProperty)(nil)).Elem(),
			"SitePublic":     reflect.TypeOf((*q.SitePublic)(nil)).Elem(),
			"SiteSession":    reflect.TypeOf((*q.SiteSession)(nil)).Elem(),
			"Sites":          reflect.TypeOf((*q.Sites)(nil)).Elem(),
		},
		AliasTypes: map[string]reflect.Type{},
		Vars:       map[string]reflect.Value{},
		Funcs:      map[string]reflect.Value{},
		TypedConsts: map[string]igop.TypedConst{
			"SiteLogLevelDisable": {reflect.TypeOf(q.SiteLogLevelDisable), constant.MakeInt64(int64(q.SiteLogLevelDisable))},
		},
		UntypedConsts: map[string]igop.UntypedConst{},
	})
}
