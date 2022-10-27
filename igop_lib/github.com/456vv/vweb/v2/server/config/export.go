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
			"Config":               reflect.TypeOf((*q.Config)(nil)).Elem(),
			"ConfigConn":           reflect.TypeOf((*q.ConfigConn)(nil)).Elem(),
			"ConfigListen":         reflect.TypeOf((*q.ConfigListen)(nil)).Elem(),
			"ConfigServer":         reflect.TypeOf((*q.ConfigServer)(nil)).Elem(),
			"ConfigServerPublic":   reflect.TypeOf((*q.ConfigServerPublic)(nil)).Elem(),
			"ConfigServerTLS":      reflect.TypeOf((*q.ConfigServerTLS)(nil)).Elem(),
			"ConfigServerTLSFile":  reflect.TypeOf((*q.ConfigServerTLSFile)(nil)).Elem(),
			"ConfigServers":        reflect.TypeOf((*q.ConfigServers)(nil)).Elem(),
			"ConfigSite":           reflect.TypeOf((*q.ConfigSite)(nil)).Elem(),
			"ConfigSiteDirectory":  reflect.TypeOf((*q.ConfigSiteDirectory)(nil)).Elem(),
			"ConfigSiteDynamic":    reflect.TypeOf((*q.ConfigSiteDynamic)(nil)).Elem(),
			"ConfigSiteForward":    reflect.TypeOf((*q.ConfigSiteForward)(nil)).Elem(),
			"ConfigSiteForwards":   reflect.TypeOf((*q.ConfigSiteForwards)(nil)).Elem(),
			"ConfigSiteHeader":     reflect.TypeOf((*q.ConfigSiteHeader)(nil)).Elem(),
			"ConfigSiteHeaderType": reflect.TypeOf((*q.ConfigSiteHeaderType)(nil)).Elem(),
			"ConfigSiteLog":        reflect.TypeOf((*q.ConfigSiteLog)(nil)).Elem(),
			"ConfigSiteLogLevel":   reflect.TypeOf((*q.ConfigSiteLogLevel)(nil)).Elem(),
			"ConfigSitePlugin":     reflect.TypeOf((*q.ConfigSitePlugin)(nil)).Elem(),
			"ConfigSitePluginTLS":  reflect.TypeOf((*q.ConfigSitePluginTLS)(nil)).Elem(),
			"ConfigSitePlugins":    reflect.TypeOf((*q.ConfigSitePlugins)(nil)).Elem(),
			"ConfigSiteProperty":   reflect.TypeOf((*q.ConfigSiteProperty)(nil)).Elem(),
			"ConfigSitePublic":     reflect.TypeOf((*q.ConfigSitePublic)(nil)).Elem(),
			"ConfigSiteSession":    reflect.TypeOf((*q.ConfigSiteSession)(nil)).Elem(),
			"ConfigSites":          reflect.TypeOf((*q.ConfigSites)(nil)).Elem(),
		},
		AliasTypes: map[string]reflect.Type{},
		Vars:       map[string]reflect.Value{},
		Funcs:      map[string]reflect.Value{},
		TypedConsts: map[string]igop.TypedConst{
			"ConfigSiteLogLevelDisable": {reflect.TypeOf(q.ConfigSiteLogLevelDisable), constant.MakeInt64(int64(q.ConfigSiteLogLevelDisable))},
		},
		UntypedConsts: map[string]igop.UntypedConst{},
	})
}
