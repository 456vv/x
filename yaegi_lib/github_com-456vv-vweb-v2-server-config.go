// Code generated by 'yaegi extract github.com/456vv/vweb/v2/server/config'. DO NOT EDIT.

//go:build yaegi_lib
// +build yaegi_lib

package yaegi_lib

import (
	"github.com/456vv/vweb/v2/server/config"
	"reflect"
)

func init() {
	Symbols["github.com/456vv/vweb/v2/server/config/config"] = map[string]reflect.Value{
		// function, constant and variable definitions
		"SiteLogLevelDisable": reflect.ValueOf(config.SiteLogLevelDisable),

		// type definitions
		"Config":         reflect.ValueOf((*config.Config)(nil)),
		"Conn":           reflect.ValueOf((*config.Conn)(nil)),
		"Listen":         reflect.ValueOf((*config.Listen)(nil)),
		"Server":         reflect.ValueOf((*config.Server)(nil)),
		"ServerPublic":   reflect.ValueOf((*config.ServerPublic)(nil)),
		"ServerTLS":      reflect.ValueOf((*config.ServerTLS)(nil)),
		"ServerTLSFile":  reflect.ValueOf((*config.ServerTLSFile)(nil)),
		"Servers":        reflect.ValueOf((*config.Servers)(nil)),
		"Site":           reflect.ValueOf((*config.Site)(nil)),
		"SiteDirectory":  reflect.ValueOf((*config.SiteDirectory)(nil)),
		"SiteDynamic":    reflect.ValueOf((*config.SiteDynamic)(nil)),
		"SiteForward":    reflect.ValueOf((*config.SiteForward)(nil)),
		"SiteForwards":   reflect.ValueOf((*config.SiteForwards)(nil)),
		"SiteHeader":     reflect.ValueOf((*config.SiteHeader)(nil)),
		"SiteHeaderType": reflect.ValueOf((*config.SiteHeaderType)(nil)),
		"SiteLog":        reflect.ValueOf((*config.SiteLog)(nil)),
		"SiteLogLevel":   reflect.ValueOf((*config.SiteLogLevel)(nil)),
		"SitePlugin":     reflect.ValueOf((*config.SitePlugin)(nil)),
		"SitePluginTLS":  reflect.ValueOf((*config.SitePluginTLS)(nil)),
		"SitePlugins":    reflect.ValueOf((*config.SitePlugins)(nil)),
		"SiteProperty":   reflect.ValueOf((*config.SiteProperty)(nil)),
		"SitePublic":     reflect.ValueOf((*config.SitePublic)(nil)),
		"SiteSession":    reflect.ValueOf((*config.SiteSession)(nil)),
		"Sites":          reflect.ValueOf((*config.Sites)(nil)),
	}
}
