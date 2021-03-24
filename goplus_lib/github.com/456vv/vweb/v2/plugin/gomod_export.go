// Package plugin provide Go+ "github.com/456vv/vweb/v2/plugin" package, as "github.com/456vv/vweb/v2/plugin" package in Go.
package plugin

import (
	tls "crypto/tls"
	net "net"
	reflect "reflect"

	plugin "github.com/456vv/vweb/v2/plugin"
	gop "github.com/goplus/gop"
)

func execNewServerHTTP(_ int, p *gop.Context) {
	ret0 := plugin.NewServerHTTP()
	p.Ret(0, ret0)
}

func execNewServerRPC(_ int, p *gop.Context) {
	ret0 := plugin.NewServerRPC()
	p.Ret(0, ret0)
}

func execmServerHTTPLoadTLS(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*plugin.ServerHTTP).LoadTLS(args[1].(*tls.Config), args[2].([]plugin.ServerTLSFile))
	p.Ret(3, ret0)
}

func execmServerHTTPListenAndServe(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*plugin.ServerHTTP).ListenAndServe()
	p.Ret(1, ret0)
}

func toType0(v interface{}) net.Listener {
	if v == nil {
		return nil
	}
	return v.(net.Listener)
}

func execmServerHTTPServe(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*plugin.ServerHTTP).Serve(toType0(args[1]))
	p.Ret(2, ret0)
}

func execmServerHTTPClose(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*plugin.ServerHTTP).Close()
	p.Ret(1, ret0)
}

func execmServerRPCRegister(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*plugin.ServerRPC).Register(args[1])
	p.PopN(2)
}

func execmServerRPCRegisterName(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*plugin.ServerRPC).RegisterName(args[1].(string), args[2])
	p.Ret(3, ret0)
}

func execmServerRPCHandleHTTP(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	args[0].(*plugin.ServerRPC).HandleHTTP(args[1].(string), args[2].(string))
	p.PopN(3)
}

func execmServerRPCListenAndServe(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*plugin.ServerRPC).ListenAndServe()
	p.Ret(1, ret0)
}

func execmServerRPCServe(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*plugin.ServerRPC).Serve(toType0(args[1]))
	p.Ret(2, ret0)
}

func execmServerRPCClose(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*plugin.ServerRPC).Close()
	p.Ret(1, ret0)
}

// I is a Go package instance.
var I = gop.NewGoPackage("github.com/456vv/vweb/v2/plugin")

func init() {
	I.RegisterFuncs(
		I.Func("NewServerHTTP", plugin.NewServerHTTP, execNewServerHTTP),
		I.Func("NewServerRPC", plugin.NewServerRPC, execNewServerRPC),
		I.Func("(*ServerHTTP).LoadTLS", (*plugin.ServerHTTP).LoadTLS, execmServerHTTPLoadTLS),
		I.Func("(*ServerHTTP).ListenAndServe", (*plugin.ServerHTTP).ListenAndServe, execmServerHTTPListenAndServe),
		I.Func("(*ServerHTTP).Serve", (*plugin.ServerHTTP).Serve, execmServerHTTPServe),
		I.Func("(*ServerHTTP).Close", (*plugin.ServerHTTP).Close, execmServerHTTPClose),
		I.Func("(*ServerRPC).Register", (*plugin.ServerRPC).Register, execmServerRPCRegister),
		I.Func("(*ServerRPC).RegisterName", (*plugin.ServerRPC).RegisterName, execmServerRPCRegisterName),
		I.Func("(*ServerRPC).HandleHTTP", (*plugin.ServerRPC).HandleHTTP, execmServerRPCHandleHTTP),
		I.Func("(*ServerRPC).ListenAndServe", (*plugin.ServerRPC).ListenAndServe, execmServerRPCListenAndServe),
		I.Func("(*ServerRPC).Serve", (*plugin.ServerRPC).Serve, execmServerRPCServe),
		I.Func("(*ServerRPC).Close", (*plugin.ServerRPC).Close, execmServerRPCClose),
	)
	I.RegisterTypes(
		I.Type("ServerHTTP", reflect.TypeOf((*plugin.ServerHTTP)(nil)).Elem()),
		I.Type("ServerRPC", reflect.TypeOf((*plugin.ServerRPC)(nil)).Elem()),
		I.Type("ServerTLSFile", reflect.TypeOf((*plugin.ServerTLSFile)(nil)).Elem()),
	)
}
