// Package vmap provide Go+ "github.com/456vv/vmap/v2" package, as "github.com/456vv/vmap/v2" package in Go.
package vmap

import (
	reflect "reflect"
	time "time"

	vmap "github.com/456vv/vmap/v2"
	gop "github.com/goplus/gop"
)

func execmMapNew(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*vmap.Map).New(args[1])
	p.Ret(2, ret0)
}

func execmMapGetNewMap(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*vmap.Map).GetNewMap(args[1])
	p.Ret(2, ret0)
}

func execmMapGetNewMaps(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	ret0 := args[0].(*vmap.Map).GetNewMaps(args[1:]...)
	p.Ret(arity, ret0)
}

func execmMapLen(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*vmap.Map).Len()
	p.Ret(1, ret0)
}

func execmMapSet(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	args[0].(*vmap.Map).Set(args[1], args[2])
	p.PopN(3)
}

func execmMapSetExpired(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	args[0].(*vmap.Map).SetExpired(args[1], args[2].(time.Duration))
	p.PopN(3)
}

func execmMapSetExpiredCall(_ int, p *gop.Context) {
	args := p.GetArgs(4)
	args[0].(*vmap.Map).SetExpiredCall(args[1], args[2].(time.Duration), args[3].(func(interface{})))
	p.PopN(4)
}

func execmMapGet(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*vmap.Map).Get(args[1])
	p.Ret(2, ret0)
}

func execmMapHas(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*vmap.Map).Has(args[1])
	p.Ret(2, ret0)
}

func execmMapGetHas(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(*vmap.Map).GetHas(args[1])
	p.Ret(2, ret0, ret1)
}

func execmMapGetOrDefault(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*vmap.Map).GetOrDefault(args[1], args[2])
	p.Ret(3, ret0)
}

func execmMapIndex(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	ret0 := args[0].(*vmap.Map).Index(args[1:]...)
	p.Ret(arity, ret0)
}

func execmMapIndexHas(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	ret0, ret1 := args[0].(*vmap.Map).IndexHas(args[1:]...)
	p.Ret(arity, ret0, ret1)
}

func execmMapDel(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*vmap.Map).Del(args[1])
	p.PopN(2)
}

func execmMapDels(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*vmap.Map).Dels(args[1].([]interface{}))
	p.PopN(2)
}

func execmMapReadAll(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*vmap.Map).ReadAll()
	p.Ret(1, ret0)
}

func execmMapKeys(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*vmap.Map).Keys()
	p.Ret(1, ret0)
}

func execmMapValues(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*vmap.Map).Values()
	p.Ret(1, ret0)
}

func execmMapReset(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	args[0].(*vmap.Map).Reset()
	p.PopN(1)
}

func execmMapCopy(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	args[0].(*vmap.Map).Copy(args[1].(*vmap.Map), args[2].(bool))
	p.PopN(3)
}

func execmMapMarshalJSON(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := args[0].(*vmap.Map).MarshalJSON()
	p.Ret(1, ret0, ret1)
}

func execmMapString(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*vmap.Map).String()
	p.Ret(1, ret0)
}

func execmMapUnmarshalJSON(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*vmap.Map).UnmarshalJSON(args[1].([]byte))
	p.Ret(2, ret0)
}

func execmMapWriteTo(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*vmap.Map).WriteTo(args[1])
	p.Ret(2, ret0)
}

func execmMapReadFrom(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*vmap.Map).ReadFrom(args[1])
	p.Ret(2, ret0)
}

func execmMapRange(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*vmap.Map).Range(args[1].(func(key interface{}, value interface{}) bool))
	p.PopN(2)
}

func execNewMap(_ int, p *gop.Context) {
	ret0 := vmap.NewMap()
	p.Ret(0, ret0)
}

// I is a Go package instance.
var I = gop.NewGoPackage("github.com/456vv/vmap/v2")

func init() {
	I.RegisterFuncs(
		I.Func("(*Map).New", (*vmap.Map).New, execmMapNew),
		I.Func("(*Map).GetNewMap", (*vmap.Map).GetNewMap, execmMapGetNewMap),
		I.Func("(*Map).Len", (*vmap.Map).Len, execmMapLen),
		I.Func("(*Map).Set", (*vmap.Map).Set, execmMapSet),
		I.Func("(*Map).SetExpired", (*vmap.Map).SetExpired, execmMapSetExpired),
		I.Func("(*Map).SetExpiredCall", (*vmap.Map).SetExpiredCall, execmMapSetExpiredCall),
		I.Func("(*Map).Get", (*vmap.Map).Get, execmMapGet),
		I.Func("(*Map).Has", (*vmap.Map).Has, execmMapHas),
		I.Func("(*Map).GetHas", (*vmap.Map).GetHas, execmMapGetHas),
		I.Func("(*Map).GetOrDefault", (*vmap.Map).GetOrDefault, execmMapGetOrDefault),
		I.Func("(*Map).Del", (*vmap.Map).Del, execmMapDel),
		I.Func("(*Map).Dels", (*vmap.Map).Dels, execmMapDels),
		I.Func("(*Map).ReadAll", (*vmap.Map).ReadAll, execmMapReadAll),
		I.Func("(*Map).Keys", (*vmap.Map).Keys, execmMapKeys),
		I.Func("(*Map).Values", (*vmap.Map).Values, execmMapValues),
		I.Func("(*Map).Reset", (*vmap.Map).Reset, execmMapReset),
		I.Func("(*Map).Copy", (*vmap.Map).Copy, execmMapCopy),
		I.Func("(*Map).MarshalJSON", (*vmap.Map).MarshalJSON, execmMapMarshalJSON),
		I.Func("(*Map).String", (*vmap.Map).String, execmMapString),
		I.Func("(*Map).UnmarshalJSON", (*vmap.Map).UnmarshalJSON, execmMapUnmarshalJSON),
		I.Func("(*Map).WriteTo", (*vmap.Map).WriteTo, execmMapWriteTo),
		I.Func("(*Map).ReadFrom", (*vmap.Map).ReadFrom, execmMapReadFrom),
		I.Func("(*Map).Range", (*vmap.Map).Range, execmMapRange),
		I.Func("NewMap", vmap.NewMap, execNewMap),
	)
	I.RegisterFuncvs(
		I.Funcv("(*Map).GetNewMaps", (*vmap.Map).GetNewMaps, execmMapGetNewMaps),
		I.Funcv("(*Map).Index", (*vmap.Map).Index, execmMapIndex),
		I.Funcv("(*Map).IndexHas", (*vmap.Map).IndexHas, execmMapIndexHas),
	)
	I.RegisterTypes(
		I.Type("Map", reflect.TypeOf((*vmap.Map)(nil)).Elem()),
	)
}
