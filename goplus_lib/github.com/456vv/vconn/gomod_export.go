// Package vconn provide Go+ "github.com/456vv/vconn" package, as "github.com/456vv/vconn" package in Go.
package vconn

import (
	net "net"
	reflect "reflect"
	time "time"

	vconn "github.com/456vv/vconn"
	gop "github.com/goplus/gop"
)

func execiCloseNotifierCloseNotify(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(vconn.CloseNotifier).CloseNotify()
	p.Ret(1, ret0)
}

func execmConnSetBackgroundReadDiscard(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*vconn.Conn).SetBackgroundReadDiscard(args[1].(bool))
	p.PopN(2)
}

func execmConnSetReadLimit(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*vconn.Conn).SetReadLimit(args[1].(int))
	p.PopN(2)
}

func execmConnDisableBackgroundRead(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*vconn.Conn).DisableBackgroundRead(args[1].(bool))
	p.PopN(2)
}

func execmConnCloseNotify(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*vconn.Conn).CloseNotify()
	p.Ret(1, ret0)
}

func execmConnRead(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(*vconn.Conn).Read(args[1].([]byte))
	p.Ret(2, ret0, ret1)
}

func execmConnWrite(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(*vconn.Conn).Write(args[1].([]byte))
	p.Ret(2, ret0, ret1)
}

func execmConnClose(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*vconn.Conn).Close()
	p.Ret(1, ret0)
}

func execmConnLocalAddr(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*vconn.Conn).LocalAddr()
	p.Ret(1, ret0)
}

func execmConnRemoteAddr(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*vconn.Conn).RemoteAddr()
	p.Ret(1, ret0)
}

func execmConnSetDeadline(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*vconn.Conn).SetDeadline(args[1].(time.Time))
	p.Ret(2, ret0)
}

func execmConnSetReadDeadline(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*vconn.Conn).SetReadDeadline(args[1].(time.Time))
	p.Ret(2, ret0)
}

func execmConnSetWriteDeadline(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*vconn.Conn).SetWriteDeadline(args[1].(time.Time))
	p.Ret(2, ret0)
}

func toType0(v interface{}) net.Conn {
	if v == nil {
		return nil
	}
	return v.(net.Conn)
}

func execNewConn(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := vconn.NewConn(toType0(args[0]))
	p.Ret(1, ret0)
}

// I is a Go package instance.
var I = gop.NewGoPackage("github.com/456vv/vconn")

func init() {
	I.RegisterFuncs(
		I.Func("(CloseNotifier).CloseNotify", (vconn.CloseNotifier).CloseNotify, execiCloseNotifierCloseNotify),
		I.Func("(*Conn).SetBackgroundReadDiscard", (*vconn.Conn).SetBackgroundReadDiscard, execmConnSetBackgroundReadDiscard),
		I.Func("(*Conn).SetReadLimit", (*vconn.Conn).SetReadLimit, execmConnSetReadLimit),
		I.Func("(*Conn).DisableBackgroundRead", (*vconn.Conn).DisableBackgroundRead, execmConnDisableBackgroundRead),
		I.Func("(*Conn).CloseNotify", (*vconn.Conn).CloseNotify, execmConnCloseNotify),
		I.Func("(*Conn).Read", (*vconn.Conn).Read, execmConnRead),
		I.Func("(*Conn).Write", (*vconn.Conn).Write, execmConnWrite),
		I.Func("(*Conn).Close", (*vconn.Conn).Close, execmConnClose),
		I.Func("(*Conn).LocalAddr", (*vconn.Conn).LocalAddr, execmConnLocalAddr),
		I.Func("(*Conn).RemoteAddr", (*vconn.Conn).RemoteAddr, execmConnRemoteAddr),
		I.Func("(*Conn).SetDeadline", (*vconn.Conn).SetDeadline, execmConnSetDeadline),
		I.Func("(*Conn).SetReadDeadline", (*vconn.Conn).SetReadDeadline, execmConnSetReadDeadline),
		I.Func("(*Conn).SetWriteDeadline", (*vconn.Conn).SetWriteDeadline, execmConnSetWriteDeadline),
		I.Func("NewConn", vconn.NewConn, execNewConn),
	)
	I.RegisterTypes(
		I.Type("CloseNotifier", reflect.TypeOf((*vconn.CloseNotifier)(nil)).Elem()),
		I.Type("Conn", reflect.TypeOf((*vconn.Conn)(nil)).Elem()),
	)
}
