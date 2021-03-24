// Package builtin provide Go+ "github.com/456vv/vweb/v2/builtin" package, as "github.com/456vv/vweb/v2/builtin" package in Go.
package builtin

import (
	reflect "reflect"

	builtin "github.com/456vv/vweb/v2/builtin"
	gop "github.com/goplus/gop"
)

func execAdd(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := builtin.Add(args[0], args[1])
	p.Ret(2, ret0)
}

func execAnd(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	ret0 := builtin.And(args[0], args[1:]...)
	p.Ret(arity, ret0)
}

func execAppend(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	ret0 := builtin.Append(args[0], args[1:]...)
	p.Ret(arity, ret0)
}

func execBitAnd(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := builtin.BitAnd(args[0], args[1])
	p.Ret(2, ret0)
}

func execBitAndNot(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := builtin.BitAndNot(args[0], args[1])
	p.Ret(2, ret0)
}

func execBitLshr(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := builtin.BitLshr(args[0], args[1])
	p.Ret(2, ret0)
}

func execBitNot(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := builtin.BitNot(args[0])
	p.Ret(1, ret0)
}

func execBitOr(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := builtin.BitOr(args[0], args[1])
	p.Ret(2, ret0)
}

func execBitRshr(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := builtin.BitRshr(args[0], args[1])
	p.Ret(2, ret0)
}

func execBitXor(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := builtin.BitXor(args[0], args[1])
	p.Ret(2, ret0)
}

func execBool(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := builtin.Bool(args[0])
	p.Ret(1, ret0)
}

func execByte(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := builtin.Byte(args[0])
	p.Ret(1, ret0)
}

func execBytes(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := builtin.Bytes(args[0])
	p.Ret(1, ret0)
}

func execCap(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := builtin.Cap(args[0])
	p.Ret(1, ret0)
}

func execChanOf(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := builtin.ChanOf(args[0])
	p.Ret(1, ret0)
}

func execClose(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	builtin.Close(args[0])
	p.PopN(1)
}

func execComplex128(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := builtin.Complex128(args[0])
	p.Ret(1, ret0)
}

func execComplex64(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := builtin.Complex64(args[0])
	p.Ret(1, ret0)
}

func execCompute(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0, ret1 := builtin.Compute(args[0], args[1].(string), args[2])
	p.Ret(3, ret0, ret1)
}

func execCopy(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := builtin.Copy(args[0], args[1])
	p.Ret(2, ret0)
}

func execDec(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := builtin.Dec(args[0])
	p.Ret(1, ret0)
}

func execDelete(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	builtin.Delete(args[0], args[1])
	p.PopN(2)
}

func execEQ(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := builtin.EQ(args[0], args[1])
	p.Ret(2, ret0)
}

func execFloat32(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := builtin.Float32(args[0])
	p.Ret(1, ret0)
}

func execFloat64(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := builtin.Float64(args[0])
	p.Ret(1, ret0)
}

func execGE(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := builtin.GE(args[0], args[1])
	p.Ret(2, ret0)
}

func execGT(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := builtin.GT(args[0], args[1])
	p.Ret(2, ret0)
}

func execGet(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := builtin.Get(args[0], args[1])
	p.Ret(2, ret0)
}

func execGetSlice(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := builtin.GetSlice(args[0], args[1], args[2])
	p.Ret(3, ret0)
}

func execGetSlice3(_ int, p *gop.Context) {
	args := p.GetArgs(4)
	ret0 := builtin.GetSlice3(args[0], args[1], args[2], args[3])
	p.Ret(4, ret0)
}

func execGoTypeInit(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	builtin.GoTypeInit(args[0], args[1:]...)
	p.PopN(arity)
}

func execGoTypeTo(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	ret0 := builtin.GoTypeTo(args[0], args[1:]...)
	p.Ret(arity, ret0)
}

func execInc(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := builtin.Inc(args[0])
	p.Ret(1, ret0)
}

func execInt(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := builtin.Int(args[0])
	p.Ret(1, ret0)
}

func execInt16(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := builtin.Int16(args[0])
	p.Ret(1, ret0)
}

func execInt32(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := builtin.Int32(args[0])
	p.Ret(1, ret0)
}

func execInt64(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := builtin.Int64(args[0])
	p.Ret(1, ret0)
}

func execInt8(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := builtin.Int8(args[0])
	p.Ret(1, ret0)
}

func execLE(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := builtin.LE(args[0], args[1])
	p.Ret(2, ret0)
}

func execLT(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := builtin.LT(args[0], args[1])
	p.Ret(2, ret0)
}

func execLen(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := builtin.Len(args[0])
	p.Ret(1, ret0)
}

func execMake(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	ret0 := builtin.Make(args[0], args[1:]...)
	p.Ret(arity, ret0)
}

func execMakeChan(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	ret0 := builtin.MakeChan(args[0], args[1:]...)
	p.Ret(arity, ret0)
}

func execMapFrom(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	ret0 := builtin.MapFrom(args[0], args[1:]...)
	p.Ret(arity, ret0)
}

func execMax(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	ret0 := builtin.Max(args...)
	p.Ret(arity, ret0)
}

func execMin(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	ret0 := builtin.Min(args...)
	p.Ret(arity, ret0)
}

func execMod(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := builtin.Mod(args[0], args[1])
	p.Ret(2, ret0)
}

func execMul(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := builtin.Mul(args[0], args[1])
	p.Ret(2, ret0)
}

func execNE(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := builtin.NE(args[0], args[1])
	p.Ret(2, ret0)
}

func execNeg(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := builtin.Neg(args[0])
	p.Ret(1, ret0)
}

func execNot(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := builtin.Not(args[0])
	p.Ret(1, ret0)
}

func execOr(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	ret0 := builtin.Or(args[0], args[1:]...)
	p.Ret(arity, ret0)
}

func execPanic(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	builtin.Panic(args[0])
	p.PopN(1)
}

func execPointer(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := builtin.Pointer(args[0])
	p.Ret(1, ret0)
}

func execQuo(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := builtin.Quo(args[0], args[1])
	p.Ret(2, ret0)
}

func execRecv(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := builtin.Recv(args[0])
	p.Ret(1, ret0)
}

func execRune(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := builtin.Rune(args[0])
	p.Ret(1, ret0)
}

func execRuns(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := builtin.Runs(args[0])
	p.Ret(1, ret0)
}

func execSend(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	builtin.Send(args[0], args[1])
	p.PopN(2)
}

func execSet(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	builtin.Set(args[0], args[1:]...)
	p.PopN(arity)
}

func execSliceFrom(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	ret0 := builtin.SliceFrom(args[0], args[1:]...)
	p.Ret(arity, ret0)
}

func execString(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := builtin.String(args[0])
	p.Ret(1, ret0)
}

func execSub(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := builtin.Sub(args[0], args[1])
	p.Ret(2, ret0)
}

func execTryRecv(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := builtin.TryRecv(args[0])
	p.Ret(1, ret0)
}

func execTrySend(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := builtin.TrySend(args[0], args[1])
	p.Ret(2, ret0)
}

func execType(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := builtin.Type(args[0])
	p.Ret(1, ret0)
}

func execUint(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := builtin.Uint(args[0])
	p.Ret(1, ret0)
}

func execUint16(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := builtin.Uint16(args[0])
	p.Ret(1, ret0)
}

func execUint32(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := builtin.Uint32(args[0])
	p.Ret(1, ret0)
}

func execUint64(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := builtin.Uint64(args[0])
	p.Ret(1, ret0)
}

func execUint8(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := builtin.Uint8(args[0])
	p.Ret(1, ret0)
}

func execUintptr(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := builtin.Uintptr(args[0])
	p.Ret(1, ret0)
}

func execValue(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := builtin.Value(args[0])
	p.Ret(1, ret0)
}

// I is a Go package instance.
var I = gop.NewGoPackage("github.com/456vv/vweb/v2/builtin")

func init() {
	I.RegisterFuncs(
		I.Func("Add", builtin.Add, execAdd),
		I.Func("BitAnd", builtin.BitAnd, execBitAnd),
		I.Func("BitAndNot", builtin.BitAndNot, execBitAndNot),
		I.Func("BitLshr", builtin.BitLshr, execBitLshr),
		I.Func("BitNot", builtin.BitNot, execBitNot),
		I.Func("BitOr", builtin.BitOr, execBitOr),
		I.Func("BitRshr", builtin.BitRshr, execBitRshr),
		I.Func("BitXor", builtin.BitXor, execBitXor),
		I.Func("Bool", builtin.Bool, execBool),
		I.Func("Byte", builtin.Byte, execByte),
		I.Func("Bytes", builtin.Bytes, execBytes),
		I.Func("Cap", builtin.Cap, execCap),
		I.Func("ChanOf", builtin.ChanOf, execChanOf),
		I.Func("Close", builtin.Close, execClose),
		I.Func("Complex128", builtin.Complex128, execComplex128),
		I.Func("Complex64", builtin.Complex64, execComplex64),
		I.Func("Compute", builtin.Compute, execCompute),
		I.Func("Copy", builtin.Copy, execCopy),
		I.Func("Dec", builtin.Dec, execDec),
		I.Func("Delete", builtin.Delete, execDelete),
		I.Func("EQ", builtin.EQ, execEQ),
		I.Func("Float32", builtin.Float32, execFloat32),
		I.Func("Float64", builtin.Float64, execFloat64),
		I.Func("GE", builtin.GE, execGE),
		I.Func("GT", builtin.GT, execGT),
		I.Func("Get", builtin.Get, execGet),
		I.Func("GetSlice", builtin.GetSlice, execGetSlice),
		I.Func("GetSlice3", builtin.GetSlice3, execGetSlice3),
		I.Func("Inc", builtin.Inc, execInc),
		I.Func("Int", builtin.Int, execInt),
		I.Func("Int16", builtin.Int16, execInt16),
		I.Func("Int32", builtin.Int32, execInt32),
		I.Func("Int64", builtin.Int64, execInt64),
		I.Func("Int8", builtin.Int8, execInt8),
		I.Func("LE", builtin.LE, execLE),
		I.Func("LT", builtin.LT, execLT),
		I.Func("Len", builtin.Len, execLen),
		I.Func("Mod", builtin.Mod, execMod),
		I.Func("Mul", builtin.Mul, execMul),
		I.Func("NE", builtin.NE, execNE),
		I.Func("Neg", builtin.Neg, execNeg),
		I.Func("Not", builtin.Not, execNot),
		I.Func("Panic", builtin.Panic, execPanic),
		I.Func("Pointer", builtin.Pointer, execPointer),
		I.Func("Quo", builtin.Quo, execQuo),
		I.Func("Recv", builtin.Recv, execRecv),
		I.Func("Rune", builtin.Rune, execRune),
		I.Func("Runs", builtin.Runs, execRuns),
		I.Func("Send", builtin.Send, execSend),
		I.Func("String", builtin.String, execString),
		I.Func("Sub", builtin.Sub, execSub),
		I.Func("TryRecv", builtin.TryRecv, execTryRecv),
		I.Func("TrySend", builtin.TrySend, execTrySend),
		I.Func("Type", builtin.Type, execType),
		I.Func("Uint", builtin.Uint, execUint),
		I.Func("Uint16", builtin.Uint16, execUint16),
		I.Func("Uint32", builtin.Uint32, execUint32),
		I.Func("Uint64", builtin.Uint64, execUint64),
		I.Func("Uint8", builtin.Uint8, execUint8),
		I.Func("Uintptr", builtin.Uintptr, execUintptr),
		I.Func("Value", builtin.Value, execValue),
	)
	I.RegisterFuncvs(
		I.Funcv("And", builtin.And, execAnd),
		I.Funcv("Append", builtin.Append, execAppend),
		I.Funcv("GoTypeInit", builtin.GoTypeInit, execGoTypeInit),
		I.Funcv("GoTypeTo", builtin.GoTypeTo, execGoTypeTo),
		I.Funcv("Make", builtin.Make, execMake),
		I.Funcv("MakeChan", builtin.MakeChan, execMakeChan),
		I.Funcv("MapFrom", builtin.MapFrom, execMapFrom),
		I.Funcv("Max", builtin.Max, execMax),
		I.Funcv("Min", builtin.Min, execMin),
		I.Funcv("Or", builtin.Or, execOr),
		I.Funcv("Set", builtin.Set, execSet),
		I.Funcv("SliceFrom", builtin.SliceFrom, execSliceFrom),
	)
	I.RegisterTypes(
		I.Type("Chan", reflect.TypeOf((*builtin.Chan)(nil)).Elem()),
	)
}
