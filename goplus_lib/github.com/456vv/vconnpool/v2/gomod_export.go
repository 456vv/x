// Package vconnpool provide Go+ "github.com/456vv/vconnpool/v2" package, as "github.com/456vv/vconnpool/v2" package in Go.
package vconnpool

import (
	context "context"
	net "net"
	reflect "reflect"
	time "time"

	vconnpool "github.com/456vv/vconnpool/v2"
	gop "github.com/goplus/gop"
)

func execiConnClose(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(net.Conn).Close()
	p.Ret(1, ret0)
}

func execiConnDiscard(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(vconnpool.Conn).Discard()
	p.Ret(1, ret0)
}

func execiConnIsReuseConn(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(vconnpool.Conn).IsReuseConn()
	p.Ret(1, ret0)
}

func execiConnLocalAddr(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(net.Conn).LocalAddr()
	p.Ret(1, ret0)
}

func execiConnRawConn(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(vconnpool.Conn).RawConn()
	p.Ret(1, ret0)
}

func execiConnRead(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(net.Conn).Read(args[1].([]byte))
	p.Ret(2, ret0, ret1)
}

func execiConnRemoteAddr(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(net.Conn).RemoteAddr()
	p.Ret(1, ret0)
}

func execiConnSetDeadline(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(net.Conn).SetDeadline(args[1].(time.Time))
	p.Ret(2, ret0)
}

func execiConnSetReadDeadline(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(net.Conn).SetReadDeadline(args[1].(time.Time))
	p.Ret(2, ret0)
}

func execiConnSetWriteDeadline(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(net.Conn).SetWriteDeadline(args[1].(time.Time))
	p.Ret(2, ret0)
}

func execiConnWrite(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(net.Conn).Write(args[1].([]byte))
	p.Ret(2, ret0, ret1)
}

func execmConnPoolDial(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0, ret1 := args[0].(*vconnpool.ConnPool).Dial(args[1].(string), args[2].(string))
	p.Ret(3, ret0, ret1)
}

func toType0(v interface{}) context.Context {
	if v == nil {
		return nil
	}
	return v.(context.Context)
}

func execmConnPoolDialContext(_ int, p *gop.Context) {
	args := p.GetArgs(4)
	ret0, ret1 := args[0].(*vconnpool.ConnPool).DialContext(toType0(args[1]), args[2].(string), args[3].(string))
	p.Ret(4, ret0, ret1)
}

func toType1(v interface{}) net.Addr {
	if v == nil {
		return nil
	}
	return v.(net.Addr)
}

func execmConnPoolGet(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(*vconnpool.ConnPool).Get(toType1(args[1]))
	p.Ret(2, ret0, ret1)
}

func toType2(v interface{}) net.Conn {
	if v == nil {
		return nil
	}
	return v.(net.Conn)
}

func execmConnPoolAdd(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*vconnpool.ConnPool).Add(toType2(args[1]))
	p.Ret(2, ret0)
}

func execmConnPoolPut(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*vconnpool.ConnPool).Put(toType2(args[1]), toType1(args[2]))
	p.Ret(3, ret0)
}

func execmConnPoolConnNum(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*vconnpool.ConnPool).ConnNum()
	p.Ret(1, ret0)
}

func execmConnPoolConnNumIde(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*vconnpool.ConnPool).ConnNumIde(args[1].(string), args[2].(string))
	p.Ret(3, ret0)
}

func execmConnPoolCloseIdleConnections(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	args[0].(*vconnpool.ConnPool).CloseIdleConnections()
	p.PopN(1)
}

func execmConnPoolClose(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*vconnpool.ConnPool).Close()
	p.Ret(1, ret0)
}

func execiDialerDial(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0, ret1 := args[0].(vconnpool.Dialer).Dial(args[1].(string), args[2].(string))
	p.Ret(3, ret0, ret1)
}

func execiDialerDialContext(_ int, p *gop.Context) {
	args := p.GetArgs(4)
	ret0, ret1 := args[0].(vconnpool.Dialer).DialContext(toType0(args[1]), args[2].(string), args[3].(string))
	p.Ret(4, ret0, ret1)
}

func execParseAddr(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := vconnpool.ParseAddr(args[0].(string), args[1].(string))
	p.Ret(2, ret0)
}

// I is a Go package instance.
var I = gop.NewGoPackage("github.com/456vv/vconnpool/v2")

func init() {
	I.RegisterFuncs(
		I.Func("(Conn).Close", (net.Conn).Close, execiConnClose),
		I.Func("(Conn).Discard", (vconnpool.Conn).Discard, execiConnDiscard),
		I.Func("(Conn).IsReuseConn", (vconnpool.Conn).IsReuseConn, execiConnIsReuseConn),
		I.Func("(Conn).LocalAddr", (net.Conn).LocalAddr, execiConnLocalAddr),
		I.Func("(Conn).RawConn", (vconnpool.Conn).RawConn, execiConnRawConn),
		I.Func("(Conn).Read", (net.Conn).Read, execiConnRead),
		I.Func("(Conn).RemoteAddr", (net.Conn).RemoteAddr, execiConnRemoteAddr),
		I.Func("(Conn).SetDeadline", (net.Conn).SetDeadline, execiConnSetDeadline),
		I.Func("(Conn).SetReadDeadline", (net.Conn).SetReadDeadline, execiConnSetReadDeadline),
		I.Func("(Conn).SetWriteDeadline", (net.Conn).SetWriteDeadline, execiConnSetWriteDeadline),
		I.Func("(Conn).Write", (net.Conn).Write, execiConnWrite),
		I.Func("(*ConnPool).Dial", (*vconnpool.ConnPool).Dial, execmConnPoolDial),
		I.Func("(*ConnPool).DialContext", (*vconnpool.ConnPool).DialContext, execmConnPoolDialContext),
		I.Func("(*ConnPool).Get", (*vconnpool.ConnPool).Get, execmConnPoolGet),
		I.Func("(*ConnPool).Add", (*vconnpool.ConnPool).Add, execmConnPoolAdd),
		I.Func("(*ConnPool).Put", (*vconnpool.ConnPool).Put, execmConnPoolPut),
		I.Func("(*ConnPool).ConnNum", (*vconnpool.ConnPool).ConnNum, execmConnPoolConnNum),
		I.Func("(*ConnPool).ConnNumIde", (*vconnpool.ConnPool).ConnNumIde, execmConnPoolConnNumIde),
		I.Func("(*ConnPool).CloseIdleConnections", (*vconnpool.ConnPool).CloseIdleConnections, execmConnPoolCloseIdleConnections),
		I.Func("(*ConnPool).Close", (*vconnpool.ConnPool).Close, execmConnPoolClose),
		I.Func("(Dialer).Dial", (vconnpool.Dialer).Dial, execiDialerDial),
		I.Func("(Dialer).DialContext", (vconnpool.Dialer).DialContext, execiDialerDialContext),
		I.Func("ParseAddr", vconnpool.ParseAddr, execParseAddr),
	)
	I.RegisterTypes(
		I.Type("Conn", reflect.TypeOf((*vconnpool.Conn)(nil)).Elem()),
		I.Type("ConnPool", reflect.TypeOf((*vconnpool.ConnPool)(nil)).Elem()),
		I.Type("Dialer", reflect.TypeOf((*vconnpool.Dialer)(nil)).Elem()),
	)
}
