// Package errors provide Go+ "errors" package, as "errors" package in Go.
package errors

import (
	errors "errors"

	gop "github.com/goplus/gop"
)

func toType0(v interface{}) error {
	if v == nil {
		return nil
	}
	return v.(error)
}

func execAs(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := errors.As(toType0(args[0]), args[1])
	p.Ret(2, ret0)
}

func execIs(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := errors.Is(toType0(args[0]), toType0(args[1]))
	p.Ret(2, ret0)
}

func execNew(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := errors.New(args[0].(string))
	p.Ret(1, ret0)
}

func execUnwrap(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := errors.Unwrap(toType0(args[0]))
	p.Ret(1, ret0)
}

// I is a Go package instance.
var I = gop.NewGoPackage("errors")

func init() {
	I.RegisterFuncs(
		I.Func("As", errors.As, execAs),
		I.Func("Is", errors.Is, execIs),
		I.Func("New", errors.New, execNew),
		I.Func("Unwrap", errors.Unwrap, execUnwrap),
	)
}
