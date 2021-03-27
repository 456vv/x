// Package config provide Go+ "github.com/456vv/vweb/v2/server/config" package, as "github.com/456vv/vweb/v2/server/config" package in Go.
package config

import (
	io "io"
	reflect "reflect"

	vweb "github.com/456vv/vweb/v2"
	config "github.com/456vv/vweb/v2/server/config"
	gop "github.com/goplus/gop"
	qspec "github.com/goplus/gop/exec.spec"
)

func toType0(v interface{}) io.Reader {
	if v == nil {
		return nil
	}
	return v.(io.Reader)
}

func execConfigDataParse(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := config.ConfigDataParse(args[0].(*config.Config), toType0(args[1]))
	p.Ret(2, ret0)
}

func execConfigFileParse(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := config.ConfigFileParse(args[0].(*config.Config), args[1].(string))
	p.Ret(2, ret0)
}

func execmConfigServerPublicConfigConn(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*config.ConfigServerPublic).ConfigConn(args[1].(*config.ConfigConn), args[2].(func(name string, dsc reflect.Value, src reflect.Value) bool))
	p.Ret(3, ret0)
}

func execmConfigServerPublicConfigServer(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*config.ConfigServerPublic).ConfigServer(args[1].(*config.ConfigServer), args[2].(func(name string, dsc reflect.Value, src reflect.Value) bool))
	p.Ret(3, ret0)
}

func execmConfigServerTLSCipherSuitesAuto(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	args[0].(*config.ConfigServerTLS).CipherSuitesAuto()
	p.PopN(1)
}

func execmConfigSiteDirectoryRootDir(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*config.ConfigSiteDirectory).RootDir(args[1].(string))
	p.Ret(2, ret0)
}

func execmConfigSiteForwardsRewrite(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1, ret2 := args[0].(*config.ConfigSiteForwards).Rewrite(args[1].(string))
	p.Ret(2, ret0, ret1, ret2)
}

func execmConfigSitePluginConfigPluginHTTPClient(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*config.ConfigSitePlugin).ConfigPluginHTTPClient(args[1].(*vweb.PluginHTTPClient))
	p.Ret(2, ret0)
}

func execmConfigSitePluginConfigPluginRPCClient(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*config.ConfigSitePlugin).ConfigPluginRPCClient(args[1].(*vweb.PluginRPCClient))
	p.Ret(2, ret0)
}

func execmConfigSitePluginsConfigSitePluginRPC(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*config.ConfigSitePlugins).ConfigSitePluginRPC(args[1].(*config.ConfigSitePlugin), args[2].(func(name string, dsc reflect.Value, src reflect.Value) bool))
	p.Ret(3, ret0)
}

func execmConfigSitePluginsConfigSitePluginHTTP(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*config.ConfigSitePlugins).ConfigSitePluginHTTP(args[1].(*config.ConfigSitePlugin), args[2].(func(name string, dsc reflect.Value, src reflect.Value) bool))
	p.Ret(3, ret0)
}

func execmConfigSitePublicConfigSiteSession(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*config.ConfigSitePublic).ConfigSiteSession(args[1].(*config.ConfigSiteSession), args[2].(func(name string, dsc reflect.Value, src reflect.Value) bool))
	p.Ret(3, ret0)
}

func execmConfigSitePublicConfigSiteHeader(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*config.ConfigSitePublic).ConfigSiteHeader(args[1].(*config.ConfigSiteHeader), args[2].(func(name string, dsc reflect.Value, src reflect.Value) bool))
	p.Ret(3, ret0)
}

func execmConfigSitePublicConfigSiteForward(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*config.ConfigSitePublic).ConfigSiteForward(args[1].(*config.ConfigSiteForward), args[2].(func(name string, dsc reflect.Value, src reflect.Value) bool))
	p.Ret(3, ret0)
}

func execmConfigSitePublicConfigSiteProperty(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*config.ConfigSitePublic).ConfigSiteProperty(args[1].(*config.ConfigSiteProperty), args[2].(func(name string, dsc reflect.Value, src reflect.Value) bool))
	p.Ret(3, ret0)
}

func execmConfigSitePublicConfigSiteDynamic(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*config.ConfigSitePublic).ConfigSiteDynamic(args[1].(*config.ConfigSiteDynamic), args[2].(func(name string, dsc reflect.Value, src reflect.Value) bool))
	p.Ret(3, ret0)
}

func execiPluginerHTTP(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(config.Pluginer).HTTP(args[1].(string))
	p.Ret(2, ret0, ret1)
}

func execiPluginerRPC(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(config.Pluginer).RPC(args[1].(string))
	p.Ret(2, ret0, ret1)
}

// I is a Go package instance.
var I = gop.NewGoPackage("github.com/456vv/vweb/v2/server/config")

func init() {
	I.RegisterFuncs(
		I.Func("ConfigDataParse", config.ConfigDataParse, execConfigDataParse),
		I.Func("ConfigFileParse", config.ConfigFileParse, execConfigFileParse),
		I.Func("(*ConfigServerPublic).ConfigConn", (*config.ConfigServerPublic).ConfigConn, execmConfigServerPublicConfigConn),
		I.Func("(*ConfigServerPublic).ConfigServer", (*config.ConfigServerPublic).ConfigServer, execmConfigServerPublicConfigServer),
		I.Func("(*ConfigServerTLS).CipherSuitesAuto", (*config.ConfigServerTLS).CipherSuitesAuto, execmConfigServerTLSCipherSuitesAuto),
		I.Func("(*ConfigSiteDirectory).RootDir", (*config.ConfigSiteDirectory).RootDir, execmConfigSiteDirectoryRootDir),
		I.Func("(*ConfigSiteForwards).Rewrite", (*config.ConfigSiteForwards).Rewrite, execmConfigSiteForwardsRewrite),
		I.Func("(*ConfigSitePlugin).ConfigPluginHTTPClient", (*config.ConfigSitePlugin).ConfigPluginHTTPClient, execmConfigSitePluginConfigPluginHTTPClient),
		I.Func("(*ConfigSitePlugin).ConfigPluginRPCClient", (*config.ConfigSitePlugin).ConfigPluginRPCClient, execmConfigSitePluginConfigPluginRPCClient),
		I.Func("(*ConfigSitePlugins).ConfigSitePluginRPC", (*config.ConfigSitePlugins).ConfigSitePluginRPC, execmConfigSitePluginsConfigSitePluginRPC),
		I.Func("(*ConfigSitePlugins).ConfigSitePluginHTTP", (*config.ConfigSitePlugins).ConfigSitePluginHTTP, execmConfigSitePluginsConfigSitePluginHTTP),
		I.Func("(*ConfigSitePublic).ConfigSiteSession", (*config.ConfigSitePublic).ConfigSiteSession, execmConfigSitePublicConfigSiteSession),
		I.Func("(*ConfigSitePublic).ConfigSiteHeader", (*config.ConfigSitePublic).ConfigSiteHeader, execmConfigSitePublicConfigSiteHeader),
		I.Func("(*ConfigSitePublic).ConfigSiteForward", (*config.ConfigSitePublic).ConfigSiteForward, execmConfigSitePublicConfigSiteForward),
		I.Func("(*ConfigSitePublic).ConfigSiteProperty", (*config.ConfigSitePublic).ConfigSiteProperty, execmConfigSitePublicConfigSiteProperty),
		I.Func("(*ConfigSitePublic).ConfigSiteDynamic", (*config.ConfigSitePublic).ConfigSiteDynamic, execmConfigSitePublicConfigSiteDynamic),
		I.Func("(Pluginer).HTTP", (config.Pluginer).HTTP, execiPluginerHTTP),
		I.Func("(Pluginer).RPC", (config.Pluginer).RPC, execiPluginerRPC),
	)
	I.RegisterTypes(
		I.Type("Config", reflect.TypeOf((*config.Config)(nil)).Elem()),
		I.Type("ConfigConn", reflect.TypeOf((*config.ConfigConn)(nil)).Elem()),
		I.Type("ConfigListen", reflect.TypeOf((*config.ConfigListen)(nil)).Elem()),
		I.Type("ConfigServer", reflect.TypeOf((*config.ConfigServer)(nil)).Elem()),
		I.Type("ConfigServerPublic", reflect.TypeOf((*config.ConfigServerPublic)(nil)).Elem()),
		I.Type("ConfigServerTLS", reflect.TypeOf((*config.ConfigServerTLS)(nil)).Elem()),
		I.Type("ConfigServerTLSFile", reflect.TypeOf((*config.ConfigServerTLSFile)(nil)).Elem()),
		I.Type("ConfigServers", reflect.TypeOf((*config.ConfigServers)(nil)).Elem()),
		I.Type("ConfigSite", reflect.TypeOf((*config.ConfigSite)(nil)).Elem()),
		I.Type("ConfigSiteDirectory", reflect.TypeOf((*config.ConfigSiteDirectory)(nil)).Elem()),
		I.Type("ConfigSiteDynamic", reflect.TypeOf((*config.ConfigSiteDynamic)(nil)).Elem()),
		I.Type("ConfigSiteForward", reflect.TypeOf((*config.ConfigSiteForward)(nil)).Elem()),
		I.Type("ConfigSiteForwards", reflect.TypeOf((*config.ConfigSiteForwards)(nil)).Elem()),
		I.Type("ConfigSiteHeader", reflect.TypeOf((*config.ConfigSiteHeader)(nil)).Elem()),
		I.Type("ConfigSiteHeaderType", reflect.TypeOf((*config.ConfigSiteHeaderType)(nil)).Elem()),
		I.Type("ConfigSiteLog", reflect.TypeOf((*config.ConfigSiteLog)(nil)).Elem()),
		I.Type("ConfigSiteLogLevel", reflect.TypeOf((*config.ConfigSiteLogLevel)(nil)).Elem()),
		I.Type("ConfigSitePlugin", reflect.TypeOf((*config.ConfigSitePlugin)(nil)).Elem()),
		I.Type("ConfigSitePluginTLS", reflect.TypeOf((*config.ConfigSitePluginTLS)(nil)).Elem()),
		I.Type("ConfigSitePlugins", reflect.TypeOf((*config.ConfigSitePlugins)(nil)).Elem()),
		I.Type("ConfigSiteProperty", reflect.TypeOf((*config.ConfigSiteProperty)(nil)).Elem()),
		I.Type("ConfigSitePublic", reflect.TypeOf((*config.ConfigSitePublic)(nil)).Elem()),
		I.Type("ConfigSiteSession", reflect.TypeOf((*config.ConfigSiteSession)(nil)).Elem()),
		I.Type("ConfigSites", reflect.TypeOf((*config.ConfigSites)(nil)).Elem()),
		I.Type("Pluginer", reflect.TypeOf((*config.Pluginer)(nil)).Elem()),
	)
	I.RegisterConsts(
		I.Const("ConfigSiteLogLevelDisable", qspec.Int, config.ConfigSiteLogLevelDisable),
	)
}
