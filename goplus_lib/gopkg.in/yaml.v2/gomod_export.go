// Package yaml provide Go+ "gopkg.in/yaml.v2" package, as "gopkg.in/yaml.v2" package in Go.
package yaml

import (
	io "io"
	reflect "reflect"

	gop "github.com/goplus/gop"
	yaml "gopkg.in/yaml.v2"
)

func execmDecoderSetStrict(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*yaml.Decoder).SetStrict(args[1].(bool))
	p.PopN(2)
}

func execmDecoderDecode(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*yaml.Decoder).Decode(args[1])
	p.Ret(2, ret0)
}

func execmEncoderEncode(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*yaml.Encoder).Encode(args[1])
	p.Ret(2, ret0)
}

func execmEncoderClose(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*yaml.Encoder).Close()
	p.Ret(1, ret0)
}

func execFutureLineWrap(_ int, p *gop.Context) {
	yaml.FutureLineWrap()
	p.PopN(0)
}

func execiIsZeroerIsZero(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(yaml.IsZeroer).IsZero()
	p.Ret(1, ret0)
}

func execMarshal(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := yaml.Marshal(args[0])
	p.Ret(1, ret0, ret1)
}

func execiMarshalerMarshalYAML(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := args[0].(yaml.Marshaler).MarshalYAML()
	p.Ret(1, ret0, ret1)
}

func toType0(v interface{}) io.Reader {
	if v == nil {
		return nil
	}
	return v.(io.Reader)
}

func execNewDecoder(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := yaml.NewDecoder(toType0(args[0]))
	p.Ret(1, ret0)
}

func toType1(v interface{}) io.Writer {
	if v == nil {
		return nil
	}
	return v.(io.Writer)
}

func execNewEncoder(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := yaml.NewEncoder(toType1(args[0]))
	p.Ret(1, ret0)
}

func execmTypeErrorError(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*yaml.TypeError).Error()
	p.Ret(1, ret0)
}

func execUnmarshal(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := yaml.Unmarshal(args[0].([]byte), args[1])
	p.Ret(2, ret0)
}

func execUnmarshalStrict(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := yaml.UnmarshalStrict(args[0].([]byte), args[1])
	p.Ret(2, ret0)
}

func execiUnmarshalerUnmarshalYAML(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(yaml.Unmarshaler).UnmarshalYAML(args[1].(func(interface{}) error))
	p.Ret(2, ret0)
}

// I is a Go package instance.
var I = gop.NewGoPackage("gopkg.in/yaml.v2")

func init() {
	I.RegisterFuncs(
		I.Func("(*Decoder).SetStrict", (*yaml.Decoder).SetStrict, execmDecoderSetStrict),
		I.Func("(*Decoder).Decode", (*yaml.Decoder).Decode, execmDecoderDecode),
		I.Func("(*Encoder).Encode", (*yaml.Encoder).Encode, execmEncoderEncode),
		I.Func("(*Encoder).Close", (*yaml.Encoder).Close, execmEncoderClose),
		I.Func("FutureLineWrap", yaml.FutureLineWrap, execFutureLineWrap),
		I.Func("(IsZeroer).IsZero", (yaml.IsZeroer).IsZero, execiIsZeroerIsZero),
		I.Func("Marshal", yaml.Marshal, execMarshal),
		I.Func("(Marshaler).MarshalYAML", (yaml.Marshaler).MarshalYAML, execiMarshalerMarshalYAML),
		I.Func("NewDecoder", yaml.NewDecoder, execNewDecoder),
		I.Func("NewEncoder", yaml.NewEncoder, execNewEncoder),
		I.Func("(*TypeError).Error", (*yaml.TypeError).Error, execmTypeErrorError),
		I.Func("Unmarshal", yaml.Unmarshal, execUnmarshal),
		I.Func("UnmarshalStrict", yaml.UnmarshalStrict, execUnmarshalStrict),
		I.Func("(Unmarshaler).UnmarshalYAML", (yaml.Unmarshaler).UnmarshalYAML, execiUnmarshalerUnmarshalYAML),
	)
	I.RegisterTypes(
		I.Type("Decoder", reflect.TypeOf((*yaml.Decoder)(nil)).Elem()),
		I.Type("Encoder", reflect.TypeOf((*yaml.Encoder)(nil)).Elem()),
		I.Type("IsZeroer", reflect.TypeOf((*yaml.IsZeroer)(nil)).Elem()),
		I.Type("MapItem", reflect.TypeOf((*yaml.MapItem)(nil)).Elem()),
		I.Type("MapSlice", reflect.TypeOf((*yaml.MapSlice)(nil)).Elem()),
		I.Type("Marshaler", reflect.TypeOf((*yaml.Marshaler)(nil)).Elem()),
		I.Type("TypeError", reflect.TypeOf((*yaml.TypeError)(nil)).Elem()),
		I.Type("Unmarshaler", reflect.TypeOf((*yaml.Unmarshaler)(nil)).Elem()),
	)
}
