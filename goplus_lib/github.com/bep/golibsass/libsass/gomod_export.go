// Package libsass provide Go+ "github.com/bep/golibsass/libsass" package, as "github.com/bep/golibsass/libsass" package in Go.
package libsass

import (
	reflect "reflect"

	libsass "github.com/bep/golibsass/libsass"
	gop "github.com/goplus/gop"
	qspec "github.com/goplus/gop/exec.spec"
)

func execmErrorError(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(libsass.Error).Error()
	p.Ret(1, ret0)
}

func execNew(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := libsass.New(args[0].(libsass.Options))
	p.Ret(1, ret0, ret1)
}

func execParseOutputStyle(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := libsass.ParseOutputStyle(args[0].(string))
	p.Ret(1, ret0)
}

func execiTranspilerExecute(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(libsass.Transpiler).Execute(args[1].(string))
	p.Ret(2, ret0, ret1)
}

// I is a Go package instance.
var I = gop.NewGoPackage("github.com/bep/golibsass/libsass")

func init() {
	I.RegisterFuncs(
		I.Func("(Error).Error", (libsass.Error).Error, execmErrorError),
		I.Func("New", libsass.New, execNew),
		I.Func("ParseOutputStyle", libsass.ParseOutputStyle, execParseOutputStyle),
		I.Func("(Transpiler).Execute", (libsass.Transpiler).Execute, execiTranspilerExecute),
	)
	I.RegisterTypes(
		I.Type("Error", reflect.TypeOf((*libsass.Error)(nil)).Elem()),
		I.Type("Options", reflect.TypeOf((*libsass.Options)(nil)).Elem()),
		I.Type("OutputStyle", reflect.TypeOf((*libsass.OutputStyle)(nil)).Elem()),
		I.Type("Result", reflect.TypeOf((*libsass.Result)(nil)).Elem()),
		I.Type("SourceMapOptions", reflect.TypeOf((*libsass.SourceMapOptions)(nil)).Elem()),
		I.Type("Transpiler", reflect.TypeOf((*libsass.Transpiler)(nil)).Elem()),
	)
	I.RegisterConsts(
		I.Const("CompactStyle", qspec.Int, libsass.CompactStyle),
		I.Const("CompressedStyle", qspec.Int, libsass.CompressedStyle),
		I.Const("ExpandedStyle", qspec.Int, libsass.ExpandedStyle),
		I.Const("NestedStyle", qspec.Int, libsass.NestedStyle),
	)
}
