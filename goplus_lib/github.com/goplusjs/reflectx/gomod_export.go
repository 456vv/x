// Package reflectx provide Go+ "github.com/goplusjs/reflectx" package, as "github.com/goplusjs/reflectx" package in Go.
package reflectx

import (
	reflect "reflect"

	gop "github.com/goplus/gop"
	qspec "github.com/goplus/gop/exec.spec"
	reflectx "github.com/goplusjs/reflectx"
)

func execCanSet(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := reflectx.CanSet(args[0].(reflect.Value))
	p.Ret(1, ret0)
}

func execField(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := reflectx.Field(args[0].(reflect.Value), args[1].(int))
	p.Ret(2, ret0)
}

func execFieldByIndex(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := reflectx.FieldByIndex(args[0].(reflect.Value), args[1].([]int))
	p.Ret(2, ret0)
}

func execFieldByName(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := reflectx.FieldByName(args[0].(reflect.Value), args[1].(string))
	p.Ret(2, ret0)
}

func execFieldByNameFunc(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := reflectx.FieldByNameFunc(args[0].(reflect.Value), args[1].(func(name string) bool))
	p.Ret(2, ret0)
}

func execInterface(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := reflectx.Interface(args[0].(reflect.Value))
	p.Ret(1, ret0)
}

func toType0(v interface{}) reflect.Type {
	if v == nil {
		return nil
	}
	return v.(reflect.Type)
}

func execInterfaceOf(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := reflectx.InterfaceOf(toType0(args[0]), args[1].([]reflect.Type), args[2].([]reflectx.Method))
	p.Ret(3, ret0)
}

func execIsNamed(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := reflectx.IsNamed(toType0(args[0]))
	p.Ret(1, ret0)
}

func execMakeEmptyInterface(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := reflectx.MakeEmptyInterface(args[0].(string), args[1].(string))
	p.Ret(2, ret0)
}

func execMakeMethod(_ int, p *gop.Context) {
	args := p.GetArgs(4)
	ret0 := reflectx.MakeMethod(args[0].(string), args[1].(bool), toType0(args[2]), args[3].(func(args []reflect.Value) (result []reflect.Value)))
	p.Ret(4, ret0)
}

func execMethodByIndex(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := reflectx.MethodByIndex(toType0(args[0]), args[1].(int))
	p.Ret(2, ret0)
}

func execMethodByName(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := reflectx.MethodByName(toType0(args[0]), args[1].(string))
	p.Ret(2, ret0, ret1)
}

func execMethodOf(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := reflectx.MethodOf(toType0(args[0]), args[1].([]reflectx.Method))
	p.Ret(2, ret0)
}

func execNamedInterfaceOf(_ int, p *gop.Context) {
	args := p.GetArgs(4)
	ret0 := reflectx.NamedInterfaceOf(args[0].(string), args[1].(string), args[2].([]reflect.Type), args[3].([]reflectx.Method))
	p.Ret(4, ret0)
}

func execNamedStructOf(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := reflectx.NamedStructOf(args[0].(string), args[1].(string), args[2].([]reflect.StructField))
	p.Ret(3, ret0)
}

func execNamedTypeOf(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := reflectx.NamedTypeOf(args[0].(string), args[1].(string), toType0(args[2]))
	p.Ret(3, ret0)
}

func execNew(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := reflectx.New(toType0(args[0]))
	p.Ret(1, ret0)
}

func execSetValue(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	reflectx.SetValue(args[0].(reflect.Value), args[1].(reflect.Value))
	p.PopN(2)
}

func execStructOf(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := reflectx.StructOf(args[0].([]reflect.StructField))
	p.Ret(1, ret0)
}

func execToNamed(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := reflectx.ToNamed(toType0(args[0]))
	p.Ret(1, ret0, ret1)
}

// I is a Go package instance.
var I = gop.NewGoPackage("github.com/goplusjs/reflectx")

func init() {
	I.RegisterFuncs(
		I.Func("CanSet", reflectx.CanSet, execCanSet),
		I.Func("Field", reflectx.Field, execField),
		I.Func("FieldByIndex", reflectx.FieldByIndex, execFieldByIndex),
		I.Func("FieldByName", reflectx.FieldByName, execFieldByName),
		I.Func("FieldByNameFunc", reflectx.FieldByNameFunc, execFieldByNameFunc),
		I.Func("Interface", reflectx.Interface, execInterface),
		I.Func("InterfaceOf", reflectx.InterfaceOf, execInterfaceOf),
		I.Func("IsNamed", reflectx.IsNamed, execIsNamed),
		I.Func("MakeEmptyInterface", reflectx.MakeEmptyInterface, execMakeEmptyInterface),
		I.Func("MakeMethod", reflectx.MakeMethod, execMakeMethod),
		I.Func("MethodByIndex", reflectx.MethodByIndex, execMethodByIndex),
		I.Func("MethodByName", reflectx.MethodByName, execMethodByName),
		I.Func("MethodOf", reflectx.MethodOf, execMethodOf),
		I.Func("NamedInterfaceOf", reflectx.NamedInterfaceOf, execNamedInterfaceOf),
		I.Func("NamedStructOf", reflectx.NamedStructOf, execNamedStructOf),
		I.Func("NamedTypeOf", reflectx.NamedTypeOf, execNamedTypeOf),
		I.Func("New", reflectx.New, execNew),
		I.Func("SetValue", reflectx.SetValue, execSetValue),
		I.Func("StructOf", reflectx.StructOf, execStructOf),
		I.Func("ToNamed", reflectx.ToNamed, execToNamed),
	)
	I.RegisterVars(
		I.Var("DisableStructOfExportAllField", &reflectx.DisableStructOfExportAllField),
	)
	I.RegisterTypes(
		I.Type("ChanDir", reflect.TypeOf((*reflectx.ChanDir)(nil)).Elem()),
		I.Type("Method", reflect.TypeOf((*reflectx.Method)(nil)).Elem()),
		I.Type("Named", reflect.TypeOf((*reflectx.Named)(nil)).Elem()),
		I.Type("TypeKind", reflect.TypeOf((*reflectx.TypeKind)(nil)).Elem()),
		I.Type("Value", reflect.TypeOf((*reflectx.Value)(nil)).Elem()),
	)
	I.RegisterConsts(
		I.Const("BothDir", qspec.Int, reflectx.BothDir),
		I.Const("RecvDir", qspec.Int, reflectx.RecvDir),
		I.Const("SendDir", qspec.Int, reflectx.SendDir),
		I.Const("TkInvalid", qspec.Int, reflectx.TkInvalid),
		I.Const("TkMethod", qspec.Int, reflectx.TkMethod),
		I.Const("TkType", qspec.Int, reflectx.TkType),
	)
}
