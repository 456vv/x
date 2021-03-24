// Package vbody provide Go+ "github.com/456vv/vbody" package, as "github.com/456vv/vbody" package in Go.
package vbody

import (
	io "io"
	reflect "reflect"

	vbody "github.com/456vv/vbody"
	gop "github.com/goplus/gop"
)

func execNewReader(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := vbody.NewReader(args[0])
	p.Ret(1, ret0)
}

func execNewWriter(_ int, p *gop.Context) {
	ret0 := vbody.NewWriter()
	p.Ret(0, ret0)
}

func execmReaderErr(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*vbody.Reader).Err()
	p.Ret(1, ret0)
}

func execmReaderIsNil(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	ret0 := args[0].(*vbody.Reader).IsNil(args[1:]...)
	p.Ret(arity, ret0)
}

func execmReaderHas(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	ret0 := args[0].(*vbody.Reader).Has(args[1:]...)
	p.Ret(arity, ret0)
}

func execmReaderString(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*vbody.Reader).String(args[1].(string), args[2].(string))
	p.Ret(3, ret0)
}

func execmReaderStringAnyEqual(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	ret0 := args[0].(*vbody.Reader).StringAnyEqual(args[1].(string), gop.ToStrings(args[2:])...)
	p.Ret(arity, ret0)
}

func execmReaderBool(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*vbody.Reader).Bool(args[1].(string), args[2].(bool))
	p.Ret(3, ret0)
}

func execmReaderBoolAnyEqual(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	ret0 := args[0].(*vbody.Reader).BoolAnyEqual(args[1].(bool), gop.ToStrings(args[2:])...)
	p.Ret(arity, ret0)
}

func execmReaderFloat64(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*vbody.Reader).Float64(args[1].(string), args[2].(float64))
	p.Ret(3, ret0)
}

func execmReaderFloat64AnyEqual(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	ret0 := args[0].(*vbody.Reader).Float64AnyEqual(args[1].(float64), gop.ToStrings(args[2:])...)
	p.Ret(arity, ret0)
}

func execmReaderInt64(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*vbody.Reader).Int64(args[1].(string), args[2].(int64))
	p.Ret(3, ret0)
}

func execmReaderInt64AnyEqual(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	ret0 := args[0].(*vbody.Reader).Int64AnyEqual(args[1].(int64), gop.ToStrings(args[2:])...)
	p.Ret(arity, ret0)
}

func execmReaderInterface(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*vbody.Reader).Interface(args[1].(string), args[2])
	p.Ret(3, ret0)
}

func execmReaderNewInterface(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*vbody.Reader).NewInterface(args[1].(string), args[2])
	p.Ret(3, ret0)
}

func execmReaderArray(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*vbody.Reader).Array(args[1].(string), args[2].([]interface{}))
	p.Ret(3, ret0)
}

func execmReaderNewArray(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*vbody.Reader).NewArray(args[1].(string), args[2].([]interface{}))
	p.Ret(3, ret0)
}

func execmReaderSlice(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*vbody.Reader).Slice(args[1].(int), args[2].(int))
	p.Ret(3, ret0)
}

func execmReaderNewSlice(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*vbody.Reader).NewSlice(args[1].(int), args[2].(int))
	p.Ret(3, ret0)
}

func execmReaderIndex(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*vbody.Reader).Index(args[1].(int), args[2])
	p.Ret(3, ret0)
}

func execmReaderNewIndex(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*vbody.Reader).NewIndex(args[1].(int), args[2])
	p.Ret(3, ret0)
}

func execmReaderIndexString(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*vbody.Reader).IndexString(args[1].(int), args[2].(string))
	p.Ret(3, ret0)
}

func execmReaderIndexFloat64(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*vbody.Reader).IndexFloat64(args[1].(int), args[2].(float64))
	p.Ret(3, ret0)
}

func execmReaderIndexInt64(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*vbody.Reader).IndexInt64(args[1].(int), args[2].(int64))
	p.Ret(3, ret0)
}

func execmReaderIndexArray(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*vbody.Reader).IndexArray(args[1].(int), args[2].([]interface{}))
	p.Ret(3, ret0)
}

func execmReaderNewIndexArray(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*vbody.Reader).NewIndexArray(args[1].(int), args[2].([]interface{}))
	p.Ret(3, ret0)
}

func execmReaderReset(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*vbody.Reader).Reset(args[1])
	p.Ret(2, ret0)
}

func toType0(v interface{}) io.Reader {
	if v == nil {
		return nil
	}
	return v.(io.Reader)
}

func execmReaderReadFrom(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(*vbody.Reader).ReadFrom(toType0(args[1]))
	p.Ret(2, ret0, ret1)
}

func execmReaderMarshalJSON(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := args[0].(*vbody.Reader).MarshalJSON()
	p.Ret(1, ret0, ret1)
}

func execmReaderUnmarshalJSON(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*vbody.Reader).UnmarshalJSON(args[1].([]byte))
	p.Ret(2, ret0)
}

func execmWriterStatus(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*vbody.Writer).Status(args[1].(int))
	p.PopN(2)
}

func execmWriterMessage(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*vbody.Writer).Message(args[1])
	p.PopN(2)
}

func execmWriterMessagef(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	args[0].(*vbody.Writer).Messagef(args[1].(string), args[2:]...)
	p.PopN(arity)
}

func execmWriterSetResult(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*vbody.Writer).SetResult(args[1])
	p.PopN(2)
}

func execmWriterResult(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	args[0].(*vbody.Writer).Result(args[1].(string), args[2])
	p.PopN(3)
}

func toType1(v interface{}) io.Writer {
	if v == nil {
		return nil
	}
	return v.(io.Writer)
}

func execmWriterWriteTo(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(*vbody.Writer).WriteTo(toType1(args[1]))
	p.Ret(2, ret0, ret1)
}

func execmWriterString(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*vbody.Writer).String()
	p.Ret(1, ret0)
}

// I is a Go package instance.
var I = gop.NewGoPackage("github.com/456vv/vbody")

func init() {
	I.RegisterFuncs(
		I.Func("NewReader", vbody.NewReader, execNewReader),
		I.Func("NewWriter", vbody.NewWriter, execNewWriter),
		I.Func("(*Reader).Err", (*vbody.Reader).Err, execmReaderErr),
		I.Func("(*Reader).String", (*vbody.Reader).String, execmReaderString),
		I.Func("(*Reader).Bool", (*vbody.Reader).Bool, execmReaderBool),
		I.Func("(*Reader).Float64", (*vbody.Reader).Float64, execmReaderFloat64),
		I.Func("(*Reader).Int64", (*vbody.Reader).Int64, execmReaderInt64),
		I.Func("(*Reader).Interface", (*vbody.Reader).Interface, execmReaderInterface),
		I.Func("(*Reader).NewInterface", (*vbody.Reader).NewInterface, execmReaderNewInterface),
		I.Func("(*Reader).Array", (*vbody.Reader).Array, execmReaderArray),
		I.Func("(*Reader).NewArray", (*vbody.Reader).NewArray, execmReaderNewArray),
		I.Func("(*Reader).Slice", (*vbody.Reader).Slice, execmReaderSlice),
		I.Func("(*Reader).NewSlice", (*vbody.Reader).NewSlice, execmReaderNewSlice),
		I.Func("(*Reader).Index", (*vbody.Reader).Index, execmReaderIndex),
		I.Func("(*Reader).NewIndex", (*vbody.Reader).NewIndex, execmReaderNewIndex),
		I.Func("(*Reader).IndexString", (*vbody.Reader).IndexString, execmReaderIndexString),
		I.Func("(*Reader).IndexFloat64", (*vbody.Reader).IndexFloat64, execmReaderIndexFloat64),
		I.Func("(*Reader).IndexInt64", (*vbody.Reader).IndexInt64, execmReaderIndexInt64),
		I.Func("(*Reader).IndexArray", (*vbody.Reader).IndexArray, execmReaderIndexArray),
		I.Func("(*Reader).NewIndexArray", (*vbody.Reader).NewIndexArray, execmReaderNewIndexArray),
		I.Func("(*Reader).Reset", (*vbody.Reader).Reset, execmReaderReset),
		I.Func("(*Reader).ReadFrom", (*vbody.Reader).ReadFrom, execmReaderReadFrom),
		I.Func("(*Reader).MarshalJSON", (*vbody.Reader).MarshalJSON, execmReaderMarshalJSON),
		I.Func("(*Reader).UnmarshalJSON", (*vbody.Reader).UnmarshalJSON, execmReaderUnmarshalJSON),
		I.Func("(*Writer).Status", (*vbody.Writer).Status, execmWriterStatus),
		I.Func("(*Writer).Message", (*vbody.Writer).Message, execmWriterMessage),
		I.Func("(*Writer).SetResult", (*vbody.Writer).SetResult, execmWriterSetResult),
		I.Func("(*Writer).Result", (*vbody.Writer).Result, execmWriterResult),
		I.Func("(*Writer).WriteTo", (*vbody.Writer).WriteTo, execmWriterWriteTo),
		I.Func("(*Writer).String", (*vbody.Writer).String, execmWriterString),
	)
	I.RegisterFuncvs(
		I.Funcv("(*Reader).IsNil", (*vbody.Reader).IsNil, execmReaderIsNil),
		I.Funcv("(*Reader).Has", (*vbody.Reader).Has, execmReaderHas),
		I.Funcv("(*Reader).StringAnyEqual", (*vbody.Reader).StringAnyEqual, execmReaderStringAnyEqual),
		I.Funcv("(*Reader).BoolAnyEqual", (*vbody.Reader).BoolAnyEqual, execmReaderBoolAnyEqual),
		I.Funcv("(*Reader).Float64AnyEqual", (*vbody.Reader).Float64AnyEqual, execmReaderFloat64AnyEqual),
		I.Funcv("(*Reader).Int64AnyEqual", (*vbody.Reader).Int64AnyEqual, execmReaderInt64AnyEqual),
		I.Funcv("(*Writer).Messagef", (*vbody.Writer).Messagef, execmWriterMessagef),
	)
	I.RegisterTypes(
		I.Type("Reader", reflect.TypeOf((*vbody.Reader)(nil)).Elem()),
		I.Type("Writer", reflect.TypeOf((*vbody.Writer)(nil)).Elem()),
	)
}
