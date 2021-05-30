// Package ticker provide Go+ "github.com/456vv/x/ticker" package, as "github.com/456vv/x/ticker" package in Go.
package ticker

import (
	reflect "reflect"
	time "time"

	ticker "github.com/456vv/x/ticker"
	gop "github.com/goplus/gop"
)

func execNewTicker(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := ticker.NewTicker(args[0].(time.Duration))
	p.Ret(1, ret0)
}

func execmTickerFunc(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*ticker.Ticker).Func(args[1].(func()))
	p.PopN(2)
}

func execmTickerReset(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*ticker.Ticker).Reset(args[1].(time.Duration))
	p.PopN(2)
}

func execmTickerStop(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	args[0].(*ticker.Ticker).Stop()
	p.PopN(1)
}

// I is a Go package instance.
var I = gop.NewGoPackage("github.com/456vv/x/ticker")

func init() {
	I.RegisterFuncs(
		I.Func("NewTicker", ticker.NewTicker, execNewTicker),
		I.Func("(*Ticker).Func", (*ticker.Ticker).Func, execmTickerFunc),
		I.Func("(*Ticker).Reset", (*ticker.Ticker).Reset, execmTickerReset),
		I.Func("(*Ticker).Stop", (*ticker.Ticker).Stop, execmTickerStop),
	)
	I.RegisterTypes(
		I.Type("Ticker", reflect.TypeOf((*ticker.Ticker)(nil)).Elem()),
	)
}
