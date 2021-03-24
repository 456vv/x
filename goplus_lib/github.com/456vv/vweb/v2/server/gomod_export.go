// Package server provide Go+ "github.com/456vv/vweb/v2/server" package, as "github.com/456vv/vweb/v2/server" package in Go.
package server

import (
	io "io"
	net "net"
	reflect "reflect"

	vweb "github.com/456vv/vweb/v2"
	server "github.com/456vv/vweb/v2/server"
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
	ret0 := server.ConfigDataParse(args[0].(*server.Config), toType0(args[1]))
	p.Ret(2, ret0)
}

func execConfigFileParse(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := server.ConfigFileParse(args[0].(*server.Config), args[1].(string))
	p.Ret(2, ret0)
}

func execmConfigServerPublicConfigConn(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*server.ConfigServerPublic).ConfigConn(args[1].(*server.ConfigConn), args[2].(func(name string, dsc reflect.Value, src reflect.Value) bool))
	p.Ret(3, ret0)
}

func execmConfigServerPublicConfigServer(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*server.ConfigServerPublic).ConfigServer(args[1].(*server.ConfigServer), args[2].(func(name string, dsc reflect.Value, src reflect.Value) bool))
	p.Ret(3, ret0)
}

func execmConfigServerTLSCipherSuitesAuto(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	args[0].(*server.ConfigServerTLS).CipherSuitesAuto()
	p.PopN(1)
}

func execmConfigSiteDirectoryRootDir(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*server.ConfigSiteDirectory).RootDir(args[1].(string))
	p.Ret(2, ret0)
}

func execmConfigSiteForwardsRewrite(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1, ret2 := args[0].(*server.ConfigSiteForwards).Rewrite(args[1].(string))
	p.Ret(2, ret0, ret1, ret2)
}

func execmConfigSitePluginConfigPluginHTTPClient(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*server.ConfigSitePlugin).ConfigPluginHTTPClient(args[1].(*vweb.PluginHTTPClient))
	p.Ret(2, ret0)
}

func execmConfigSitePluginConfigPluginRPCClient(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*server.ConfigSitePlugin).ConfigPluginRPCClient(args[1].(*vweb.PluginRPCClient))
	p.Ret(2, ret0)
}

func execmConfigSitePluginsConfigSitePluginRPC(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*server.ConfigSitePlugins).ConfigSitePluginRPC(args[1].(*server.ConfigSitePlugin), args[2].(func(name string, dsc reflect.Value, src reflect.Value) bool))
	p.Ret(3, ret0)
}

func execmConfigSitePluginsConfigSitePluginHTTP(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*server.ConfigSitePlugins).ConfigSitePluginHTTP(args[1].(*server.ConfigSitePlugin), args[2].(func(name string, dsc reflect.Value, src reflect.Value) bool))
	p.Ret(3, ret0)
}

func execmConfigSitePublicConfigSiteSession(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*server.ConfigSitePublic).ConfigSiteSession(args[1].(*server.ConfigSiteSession), args[2].(func(name string, dsc reflect.Value, src reflect.Value) bool))
	p.Ret(3, ret0)
}

func execmConfigSitePublicConfigSiteHeader(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*server.ConfigSitePublic).ConfigSiteHeader(args[1].(*server.ConfigSiteHeader), args[2].(func(name string, dsc reflect.Value, src reflect.Value) bool))
	p.Ret(3, ret0)
}

func execmConfigSitePublicConfigSiteForward(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*server.ConfigSitePublic).ConfigSiteForward(args[1].(*server.ConfigSiteForward), args[2].(func(name string, dsc reflect.Value, src reflect.Value) bool))
	p.Ret(3, ret0)
}

func execmConfigSitePublicConfigSiteProperty(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*server.ConfigSitePublic).ConfigSiteProperty(args[1].(*server.ConfigSiteProperty), args[2].(func(name string, dsc reflect.Value, src reflect.Value) bool))
	p.Ret(3, ret0)
}

func execmConfigSitePublicConfigSiteDynamic(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*server.ConfigSitePublic).ConfigSiteDynamic(args[1].(*server.ConfigSiteDynamic), args[2].(func(name string, dsc reflect.Value, src reflect.Value) bool))
	p.Ret(3, ret0)
}

func execNewServerGroup(_ int, p *gop.Context) {
	ret0 := server.NewServerGroup()
	p.Ret(0, ret0)
}

func execiPluginerHTTP(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(server.Pluginer).HTTP(args[1].(string))
	p.Ret(2, ret0, ret1)
}

func execiPluginerRPC(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(server.Pluginer).RPC(args[1].(string))
	p.Ret(2, ret0, ret1)
}

func toType1(v interface{}) net.Listener {
	if v == nil {
		return nil
	}
	return v.(net.Listener)
}

func execmServerServe(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*server.Server).Serve(toType1(args[1]))
	p.Ret(2, ret0)
}

func execmServerListenAndServe(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*server.Server).ListenAndServe()
	p.Ret(1, ret0)
}

func execmServerConfigConn(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*server.Server).ConfigConn(args[1].(*server.ConfigConn))
	p.Ret(2, ret0)
}

func execmServerConfigServer(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*server.Server).ConfigServer(args[1].(*server.ConfigServer))
	p.Ret(2, ret0)
}

func execmServerGroupSetServer(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*server.ServerGroup).SetServer(args[1].(string), args[2].(*server.Server))
	p.Ret(3, ret0)
}

func execmServerGroupGetServer(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(*server.ServerGroup).GetServer(args[1].(string))
	p.Ret(2, ret0, ret1)
}

func execmServerGroupSetSitePool(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*server.ServerGroup).SetSitePool(args[1].(*vweb.SitePool))
	p.Ret(2, ret0)
}

func execmServerGroupLoadConfigFile(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1, ret2 := args[0].(*server.ServerGroup).LoadConfigFile(args[1].(string))
	p.Ret(2, ret0, ret1, ret2)
}

func execmServerGroupUpdateConfig(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*server.ServerGroup).UpdateConfig(args[1].(*server.Config))
	p.Ret(2, ret0)
}

func execmServerGroupStart(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*server.ServerGroup).Start()
	p.Ret(1, ret0)
}

func execmServerGroupClose(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*server.ServerGroup).Close()
	p.Ret(1, ret0)
}

// I is a Go package instance.
var I = gop.NewGoPackage("github.com/456vv/vweb/v2/server")

func init() {
	I.RegisterFuncs(
		I.Func("ConfigDataParse", server.ConfigDataParse, execConfigDataParse),
		I.Func("ConfigFileParse", server.ConfigFileParse, execConfigFileParse),
		I.Func("(*ConfigServerPublic).ConfigConn", (*server.ConfigServerPublic).ConfigConn, execmConfigServerPublicConfigConn),
		I.Func("(*ConfigServerPublic).ConfigServer", (*server.ConfigServerPublic).ConfigServer, execmConfigServerPublicConfigServer),
		I.Func("(*ConfigServerTLS).CipherSuitesAuto", (*server.ConfigServerTLS).CipherSuitesAuto, execmConfigServerTLSCipherSuitesAuto),
		I.Func("(*ConfigSiteDirectory).RootDir", (*server.ConfigSiteDirectory).RootDir, execmConfigSiteDirectoryRootDir),
		I.Func("(*ConfigSiteForwards).Rewrite", (*server.ConfigSiteForwards).Rewrite, execmConfigSiteForwardsRewrite),
		I.Func("(*ConfigSitePlugin).ConfigPluginHTTPClient", (*server.ConfigSitePlugin).ConfigPluginHTTPClient, execmConfigSitePluginConfigPluginHTTPClient),
		I.Func("(*ConfigSitePlugin).ConfigPluginRPCClient", (*server.ConfigSitePlugin).ConfigPluginRPCClient, execmConfigSitePluginConfigPluginRPCClient),
		I.Func("(*ConfigSitePlugins).ConfigSitePluginRPC", (*server.ConfigSitePlugins).ConfigSitePluginRPC, execmConfigSitePluginsConfigSitePluginRPC),
		I.Func("(*ConfigSitePlugins).ConfigSitePluginHTTP", (*server.ConfigSitePlugins).ConfigSitePluginHTTP, execmConfigSitePluginsConfigSitePluginHTTP),
		I.Func("(*ConfigSitePublic).ConfigSiteSession", (*server.ConfigSitePublic).ConfigSiteSession, execmConfigSitePublicConfigSiteSession),
		I.Func("(*ConfigSitePublic).ConfigSiteHeader", (*server.ConfigSitePublic).ConfigSiteHeader, execmConfigSitePublicConfigSiteHeader),
		I.Func("(*ConfigSitePublic).ConfigSiteForward", (*server.ConfigSitePublic).ConfigSiteForward, execmConfigSitePublicConfigSiteForward),
		I.Func("(*ConfigSitePublic).ConfigSiteProperty", (*server.ConfigSitePublic).ConfigSiteProperty, execmConfigSitePublicConfigSiteProperty),
		I.Func("(*ConfigSitePublic).ConfigSiteDynamic", (*server.ConfigSitePublic).ConfigSiteDynamic, execmConfigSitePublicConfigSiteDynamic),
		I.Func("NewServerGroup", server.NewServerGroup, execNewServerGroup),
		I.Func("(Pluginer).HTTP", (server.Pluginer).HTTP, execiPluginerHTTP),
		I.Func("(Pluginer).RPC", (server.Pluginer).RPC, execiPluginerRPC),
		I.Func("(*Server).Serve", (*server.Server).Serve, execmServerServe),
		I.Func("(*Server).ListenAndServe", (*server.Server).ListenAndServe, execmServerListenAndServe),
		I.Func("(*Server).ConfigConn", (*server.Server).ConfigConn, execmServerConfigConn),
		I.Func("(*Server).ConfigServer", (*server.Server).ConfigServer, execmServerConfigServer),
		I.Func("(*ServerGroup).SetServer", (*server.ServerGroup).SetServer, execmServerGroupSetServer),
		I.Func("(*ServerGroup).GetServer", (*server.ServerGroup).GetServer, execmServerGroupGetServer),
		I.Func("(*ServerGroup).SetSitePool", (*server.ServerGroup).SetSitePool, execmServerGroupSetSitePool),
		I.Func("(*ServerGroup).LoadConfigFile", (*server.ServerGroup).LoadConfigFile, execmServerGroupLoadConfigFile),
		I.Func("(*ServerGroup).UpdateConfig", (*server.ServerGroup).UpdateConfig, execmServerGroupUpdateConfig),
		I.Func("(*ServerGroup).Start", (*server.ServerGroup).Start, execmServerGroupStart),
		I.Func("(*ServerGroup).Close", (*server.ServerGroup).Close, execmServerGroupClose),
	)
	I.RegisterVars(
		I.Var("Version", &server.Version),
	)
	I.RegisterTypes(
		I.Type("Config", reflect.TypeOf((*server.Config)(nil)).Elem()),
		I.Type("ConfigConn", reflect.TypeOf((*server.ConfigConn)(nil)).Elem()),
		I.Type("ConfigListen", reflect.TypeOf((*server.ConfigListen)(nil)).Elem()),
		I.Type("ConfigServer", reflect.TypeOf((*server.ConfigServer)(nil)).Elem()),
		I.Type("ConfigServerPublic", reflect.TypeOf((*server.ConfigServerPublic)(nil)).Elem()),
		I.Type("ConfigServerTLS", reflect.TypeOf((*server.ConfigServerTLS)(nil)).Elem()),
		I.Type("ConfigServerTLSFile", reflect.TypeOf((*server.ConfigServerTLSFile)(nil)).Elem()),
		I.Type("ConfigServers", reflect.TypeOf((*server.ConfigServers)(nil)).Elem()),
		I.Type("ConfigSite", reflect.TypeOf((*server.ConfigSite)(nil)).Elem()),
		I.Type("ConfigSiteDirectory", reflect.TypeOf((*server.ConfigSiteDirectory)(nil)).Elem()),
		I.Type("ConfigSiteDynamic", reflect.TypeOf((*server.ConfigSiteDynamic)(nil)).Elem()),
		I.Type("ConfigSiteForward", reflect.TypeOf((*server.ConfigSiteForward)(nil)).Elem()),
		I.Type("ConfigSiteForwards", reflect.TypeOf((*server.ConfigSiteForwards)(nil)).Elem()),
		I.Type("ConfigSiteHeader", reflect.TypeOf((*server.ConfigSiteHeader)(nil)).Elem()),
		I.Type("ConfigSiteHeaderType", reflect.TypeOf((*server.ConfigSiteHeaderType)(nil)).Elem()),
		I.Type("ConfigSiteLog", reflect.TypeOf((*server.ConfigSiteLog)(nil)).Elem()),
		I.Type("ConfigSiteLogLevel", reflect.TypeOf((*server.ConfigSiteLogLevel)(nil)).Elem()),
		I.Type("ConfigSitePlugin", reflect.TypeOf((*server.ConfigSitePlugin)(nil)).Elem()),
		I.Type("ConfigSitePluginTLS", reflect.TypeOf((*server.ConfigSitePluginTLS)(nil)).Elem()),
		I.Type("ConfigSitePlugins", reflect.TypeOf((*server.ConfigSitePlugins)(nil)).Elem()),
		I.Type("ConfigSiteProperty", reflect.TypeOf((*server.ConfigSiteProperty)(nil)).Elem()),
		I.Type("ConfigSitePublic", reflect.TypeOf((*server.ConfigSitePublic)(nil)).Elem()),
		I.Type("ConfigSiteSession", reflect.TypeOf((*server.ConfigSiteSession)(nil)).Elem()),
		I.Type("ConfigSites", reflect.TypeOf((*server.ConfigSites)(nil)).Elem()),
		I.Type("Pluginer", reflect.TypeOf((*server.Pluginer)(nil)).Elem()),
		I.Type("Server", reflect.TypeOf((*server.Server)(nil)).Elem()),
		I.Type("ServerGroup", reflect.TypeOf((*server.ServerGroup)(nil)).Elem()),
	)
	I.RegisterConsts(
		I.Const("ConfigSiteLogLevelDisable", qspec.Int, server.ConfigSiteLogLevelDisable),
	)
}
