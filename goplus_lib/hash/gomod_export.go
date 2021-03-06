// Package hash provide Go+ "hash" package, as "hash" package in Go.
package hash

import (
	hash "hash"
	io "io"
	reflect "reflect"

	gop "github.com/goplus/gop"
)

func execiHashBlockSize(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(hash.Hash).BlockSize()
	p.Ret(1, ret0)
}

func execiHashReset(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	args[0].(hash.Hash).Reset()
	p.PopN(1)
}

func execiHashSize(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(hash.Hash).Size()
	p.Ret(1, ret0)
}

func execiHashSum(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(hash.Hash).Sum(args[1].([]byte))
	p.Ret(2, ret0)
}

func execiWriterWrite(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(io.Writer).Write(args[1].([]byte))
	p.Ret(2, ret0, ret1)
}

func execiHash32Sum32(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(hash.Hash32).Sum32()
	p.Ret(1, ret0)
}

func execiHash64Sum64(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(hash.Hash64).Sum64()
	p.Ret(1, ret0)
}

// I is a Go package instance.
var I = gop.NewGoPackage("hash")

func init() {
	I.RegisterFuncs(
		I.Func("(Hash).BlockSize", (hash.Hash).BlockSize, execiHashBlockSize),
		I.Func("(Hash).Reset", (hash.Hash).Reset, execiHashReset),
		I.Func("(Hash).Size", (hash.Hash).Size, execiHashSize),
		I.Func("(Hash).Sum", (hash.Hash).Sum, execiHashSum),
		I.Func("(Writer).Write", (io.Writer).Write, execiWriterWrite),
		I.Func("(Hash32).Sum32", (hash.Hash32).Sum32, execiHash32Sum32),
		I.Func("(Hash64).Sum64", (hash.Hash64).Sum64, execiHash64Sum64),
	)
	I.RegisterTypes(
		I.Type("Hash", reflect.TypeOf((*hash.Hash)(nil)).Elem()),
		I.Type("Hash32", reflect.TypeOf((*hash.Hash32)(nil)).Elem()),
		I.Type("Hash64", reflect.TypeOf((*hash.Hash64)(nil)).Elem()),
	)
}
