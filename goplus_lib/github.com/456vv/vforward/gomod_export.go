// Package vforward provide Go+ "github.com/456vv/vforward" package, as "github.com/456vv/vforward" package in Go.
package vforward

import (
	reflect "reflect"

	vforward "github.com/456vv/vforward"
	gop "github.com/goplus/gop"
	qspec "github.com/goplus/gop/exec.spec"
)

func execmD2DTransport(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0, ret1 := args[0].(*vforward.D2D).Transport(args[1].(*vforward.Addr), args[2].(*vforward.Addr))
	p.Ret(3, ret0, ret1)
}

func execmD2DClose(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*vforward.D2D).Close()
	p.Ret(1, ret0)
}

func execmD2DSwapConnNum(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*vforward.D2DSwap).ConnNum()
	p.Ret(1, ret0)
}

func execmD2DSwapSwap(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*vforward.D2DSwap).Swap()
	p.Ret(1, ret0)
}

func execmD2DSwapClose(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*vforward.D2DSwap).Close()
	p.Ret(1, ret0)
}

func execmL2DTransport(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0, ret1 := args[0].(*vforward.L2D).Transport(args[1].(*vforward.Addr), args[2].(*vforward.Addr))
	p.Ret(3, ret0, ret1)
}

func execmL2DClose(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*vforward.L2D).Close()
	p.Ret(1, ret0)
}

func execmL2DSwapConnNum(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*vforward.L2DSwap).ConnNum()
	p.Ret(1, ret0)
}

func execmL2DSwapSwap(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*vforward.L2DSwap).Swap()
	p.Ret(1, ret0)
}

func execmL2DSwapClose(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*vforward.L2DSwap).Close()
	p.Ret(1, ret0)
}

func execmL2LTransport(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0, ret1 := args[0].(*vforward.L2L).Transport(args[1].(*vforward.Addr), args[2].(*vforward.Addr))
	p.Ret(3, ret0, ret1)
}

func execmL2LClose(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*vforward.L2L).Close()
	p.Ret(1, ret0)
}

func execmL2LSwapConnNum(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*vforward.L2LSwap).ConnNum()
	p.Ret(1, ret0)
}

func execmL2LSwapSwap(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*vforward.L2LSwap).Swap()
	p.Ret(1, ret0)
}

func execmL2LSwapClose(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*vforward.L2LSwap).Close()
	p.Ret(1, ret0)
}

// I is a Go package instance.
var I = gop.NewGoPackage("github.com/456vv/vforward")

func init() {
	I.RegisterFuncs(
		I.Func("(*D2D).Transport", (*vforward.D2D).Transport, execmD2DTransport),
		I.Func("(*D2D).Close", (*vforward.D2D).Close, execmD2DClose),
		I.Func("(*D2DSwap).ConnNum", (*vforward.D2DSwap).ConnNum, execmD2DSwapConnNum),
		I.Func("(*D2DSwap).Swap", (*vforward.D2DSwap).Swap, execmD2DSwapSwap),
		I.Func("(*D2DSwap).Close", (*vforward.D2DSwap).Close, execmD2DSwapClose),
		I.Func("(*L2D).Transport", (*vforward.L2D).Transport, execmL2DTransport),
		I.Func("(*L2D).Close", (*vforward.L2D).Close, execmL2DClose),
		I.Func("(*L2DSwap).ConnNum", (*vforward.L2DSwap).ConnNum, execmL2DSwapConnNum),
		I.Func("(*L2DSwap).Swap", (*vforward.L2DSwap).Swap, execmL2DSwapSwap),
		I.Func("(*L2DSwap).Close", (*vforward.L2DSwap).Close, execmL2DSwapClose),
		I.Func("(*L2L).Transport", (*vforward.L2L).Transport, execmL2LTransport),
		I.Func("(*L2L).Close", (*vforward.L2L).Close, execmL2LClose),
		I.Func("(*L2LSwap).ConnNum", (*vforward.L2LSwap).ConnNum, execmL2LSwapConnNum),
		I.Func("(*L2LSwap).Swap", (*vforward.L2LSwap).Swap, execmL2LSwapSwap),
		I.Func("(*L2LSwap).Close", (*vforward.L2LSwap).Close, execmL2LSwapClose),
	)
	I.RegisterTypes(
		I.Type("Addr", reflect.TypeOf((*vforward.Addr)(nil)).Elem()),
		I.Type("D2D", reflect.TypeOf((*vforward.D2D)(nil)).Elem()),
		I.Type("D2DSwap", reflect.TypeOf((*vforward.D2DSwap)(nil)).Elem()),
		I.Type("L2D", reflect.TypeOf((*vforward.L2D)(nil)).Elem()),
		I.Type("L2DSwap", reflect.TypeOf((*vforward.L2DSwap)(nil)).Elem()),
		I.Type("L2L", reflect.TypeOf((*vforward.L2L)(nil)).Elem()),
		I.Type("L2LSwap", reflect.TypeOf((*vforward.L2LSwap)(nil)).Elem()),
	)
	I.RegisterConsts(
		I.Const("DefaultReadBufSize", qspec.Int, vforward.DefaultReadBufSize),
	)
}
