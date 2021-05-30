// Package toml provide Go+ "github.com/pelletier/go-toml" package, as "github.com/pelletier/go-toml" package in Go.
package toml

import (
	io "io"
	reflect "reflect"
	time "time"

	gop "github.com/goplus/gop"
	qspec "github.com/goplus/gop/exec.spec"
	toml "github.com/pelletier/go-toml"
)

func execmDecoderDecode(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*toml.Decoder).Decode(args[1])
	p.Ret(2, ret0)
}

func execmDecoderSetTagName(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*toml.Decoder).SetTagName(args[1].(string))
	p.Ret(2, ret0)
}

func execmDecoderStrict(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*toml.Decoder).Strict(args[1].(bool))
	p.Ret(2, ret0)
}

func execmEncoderEncode(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*toml.Encoder).Encode(args[1])
	p.Ret(2, ret0)
}

func execmEncoderQuoteMapKeys(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*toml.Encoder).QuoteMapKeys(args[1].(bool))
	p.Ret(2, ret0)
}

func execmEncoderArraysWithOneElementPerLine(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*toml.Encoder).ArraysWithOneElementPerLine(args[1].(bool))
	p.Ret(2, ret0)
}

func execmEncoderOrder(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*toml.Encoder).Order(args[1].(toml.MarshalOrder))
	p.Ret(2, ret0)
}

func execmEncoderIndentation(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*toml.Encoder).Indentation(args[1].(string))
	p.Ret(2, ret0)
}

func execmEncoderSetTagName(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*toml.Encoder).SetTagName(args[1].(string))
	p.Ret(2, ret0)
}

func execmEncoderSetTagComment(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*toml.Encoder).SetTagComment(args[1].(string))
	p.Ret(2, ret0)
}

func execmEncoderSetTagCommented(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*toml.Encoder).SetTagCommented(args[1].(string))
	p.Ret(2, ret0)
}

func execmEncoderSetTagMultiline(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*toml.Encoder).SetTagMultiline(args[1].(string))
	p.Ret(2, ret0)
}

func execmEncoderPromoteAnonymous(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*toml.Encoder).PromoteAnonymous(args[1].(bool))
	p.Ret(2, ret0)
}

func execLoad(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := toml.Load(args[0].(string))
	p.Ret(1, ret0, ret1)
}

func execLoadBytes(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := toml.LoadBytes(args[0].([]byte))
	p.Ret(1, ret0, ret1)
}

func execLoadFile(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := toml.LoadFile(args[0].(string))
	p.Ret(1, ret0, ret1)
}

func toType0(v interface{}) io.Reader {
	if v == nil {
		return nil
	}
	return v.(io.Reader)
}

func execLoadReader(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := toml.LoadReader(toType0(args[0]))
	p.Ret(1, ret0, ret1)
}

func execmLocalDateString(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(toml.LocalDate).String()
	p.Ret(1, ret0)
}

func execmLocalDateIsValid(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(toml.LocalDate).IsValid()
	p.Ret(1, ret0)
}

func execmLocalDateIn(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(toml.LocalDate).In(args[1].(*time.Location))
	p.Ret(2, ret0)
}

func execmLocalDateAddDays(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(toml.LocalDate).AddDays(args[1].(int))
	p.Ret(2, ret0)
}

func execmLocalDateDaysSince(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(toml.LocalDate).DaysSince(args[1].(toml.LocalDate))
	p.Ret(2, ret0)
}

func execmLocalDateBefore(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(toml.LocalDate).Before(args[1].(toml.LocalDate))
	p.Ret(2, ret0)
}

func execmLocalDateAfter(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(toml.LocalDate).After(args[1].(toml.LocalDate))
	p.Ret(2, ret0)
}

func execmLocalDateMarshalText(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := args[0].(toml.LocalDate).MarshalText()
	p.Ret(1, ret0, ret1)
}

func execmLocalDateUnmarshalText(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*toml.LocalDate).UnmarshalText(args[1].([]byte))
	p.Ret(2, ret0)
}

func execLocalDateOf(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := toml.LocalDateOf(args[0].(time.Time))
	p.Ret(1, ret0)
}

func execmLocalDateTimeString(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(toml.LocalDateTime).String()
	p.Ret(1, ret0)
}

func execmLocalDateTimeIsValid(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(toml.LocalDateTime).IsValid()
	p.Ret(1, ret0)
}

func execmLocalDateTimeIn(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(toml.LocalDateTime).In(args[1].(*time.Location))
	p.Ret(2, ret0)
}

func execmLocalDateTimeBefore(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(toml.LocalDateTime).Before(args[1].(toml.LocalDateTime))
	p.Ret(2, ret0)
}

func execmLocalDateTimeAfter(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(toml.LocalDateTime).After(args[1].(toml.LocalDateTime))
	p.Ret(2, ret0)
}

func execmLocalDateTimeMarshalText(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := args[0].(toml.LocalDateTime).MarshalText()
	p.Ret(1, ret0, ret1)
}

func execmLocalDateTimeUnmarshalText(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*toml.LocalDateTime).UnmarshalText(args[1].([]byte))
	p.Ret(2, ret0)
}

func execLocalDateTimeOf(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := toml.LocalDateTimeOf(args[0].(time.Time))
	p.Ret(1, ret0)
}

func execmLocalTimeString(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(toml.LocalTime).String()
	p.Ret(1, ret0)
}

func execmLocalTimeIsValid(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(toml.LocalTime).IsValid()
	p.Ret(1, ret0)
}

func execmLocalTimeMarshalText(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := args[0].(toml.LocalTime).MarshalText()
	p.Ret(1, ret0, ret1)
}

func execmLocalTimeUnmarshalText(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*toml.LocalTime).UnmarshalText(args[1].([]byte))
	p.Ret(2, ret0)
}

func execLocalTimeOf(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := toml.LocalTimeOf(args[0].(time.Time))
	p.Ret(1, ret0)
}

func execMarshal(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := toml.Marshal(args[0])
	p.Ret(1, ret0, ret1)
}

func execiMarshalerMarshalTOML(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := args[0].(toml.Marshaler).MarshalTOML()
	p.Ret(1, ret0, ret1)
}

func execNewDecoder(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := toml.NewDecoder(toType0(args[0]))
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
	ret0 := toml.NewEncoder(toType1(args[0]))
	p.Ret(1, ret0)
}

func execParseLocalDate(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := toml.ParseLocalDate(args[0].(string))
	p.Ret(1, ret0, ret1)
}

func execParseLocalDateTime(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := toml.ParseLocalDateTime(args[0].(string))
	p.Ret(1, ret0, ret1)
}

func execParseLocalTime(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := toml.ParseLocalTime(args[0].(string))
	p.Ret(1, ret0, ret1)
}

func execmPositionString(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(toml.Position).String()
	p.Ret(1, ret0)
}

func execmPositionInvalid(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(toml.Position).Invalid()
	p.Ret(1, ret0)
}

func execmtomlValueValue(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*toml.PubTOMLValue).Value()
	p.Ret(1, ret0)
}

func execmtomlValueComment(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*toml.PubTOMLValue).Comment()
	p.Ret(1, ret0)
}

func execmtomlValueCommented(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*toml.PubTOMLValue).Commented()
	p.Ret(1, ret0)
}

func execmtomlValueMultiline(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*toml.PubTOMLValue).Multiline()
	p.Ret(1, ret0)
}

func execmtomlValuePosition(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*toml.PubTOMLValue).Position()
	p.Ret(1, ret0)
}

func execmtomlValueSetValue(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*toml.PubTOMLValue).SetValue(args[1])
	p.PopN(2)
}

func execmtomlValueSetComment(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*toml.PubTOMLValue).SetComment(args[1].(string))
	p.PopN(2)
}

func execmtomlValueSetCommented(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*toml.PubTOMLValue).SetCommented(args[1].(bool))
	p.PopN(2)
}

func execmtomlValueSetMultiline(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*toml.PubTOMLValue).SetMultiline(args[1].(bool))
	p.PopN(2)
}

func execmtomlValueSetPosition(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*toml.PubTOMLValue).SetPosition(args[1].(toml.Position))
	p.PopN(2)
}

func execmTreeUnmarshal(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*toml.Tree).Unmarshal(args[1])
	p.Ret(2, ret0)
}

func execmTreeMarshal(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := args[0].(*toml.Tree).Marshal()
	p.Ret(1, ret0, ret1)
}

func execmTreePosition(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*toml.Tree).Position()
	p.Ret(1, ret0)
}

func execmTreeHas(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*toml.Tree).Has(args[1].(string))
	p.Ret(2, ret0)
}

func execmTreeHasPath(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*toml.Tree).HasPath(args[1].([]string))
	p.Ret(2, ret0)
}

func execmTreeKeys(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*toml.Tree).Keys()
	p.Ret(1, ret0)
}

func execmTreeGet(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*toml.Tree).Get(args[1].(string))
	p.Ret(2, ret0)
}

func execmTreeGetPath(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*toml.Tree).GetPath(args[1].([]string))
	p.Ret(2, ret0)
}

func execmTreeGetArray(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*toml.Tree).GetArray(args[1].(string))
	p.Ret(2, ret0)
}

func execmTreeGetArrayPath(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*toml.Tree).GetArrayPath(args[1].([]string))
	p.Ret(2, ret0)
}

func execmTreeGetPosition(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*toml.Tree).GetPosition(args[1].(string))
	p.Ret(2, ret0)
}

func execmTreeSetPositionPath(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	args[0].(*toml.Tree).SetPositionPath(args[1].([]string), args[2].(toml.Position))
	p.PopN(3)
}

func execmTreeGetPositionPath(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*toml.Tree).GetPositionPath(args[1].([]string))
	p.Ret(2, ret0)
}

func execmTreeGetDefault(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*toml.Tree).GetDefault(args[1].(string), args[2])
	p.Ret(3, ret0)
}

func execmTreeSetWithOptions(_ int, p *gop.Context) {
	args := p.GetArgs(4)
	args[0].(*toml.Tree).SetWithOptions(args[1].(string), args[2].(toml.SetOptions), args[3])
	p.PopN(4)
}

func execmTreeSetPathWithOptions(_ int, p *gop.Context) {
	args := p.GetArgs(4)
	args[0].(*toml.Tree).SetPathWithOptions(args[1].([]string), args[2].(toml.SetOptions), args[3])
	p.PopN(4)
}

func execmTreeSet(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	args[0].(*toml.Tree).Set(args[1].(string), args[2])
	p.PopN(3)
}

func execmTreeSetWithComment(_ int, p *gop.Context) {
	args := p.GetArgs(5)
	args[0].(*toml.Tree).SetWithComment(args[1].(string), args[2].(string), args[3].(bool), args[4])
	p.PopN(5)
}

func execmTreeSetPath(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	args[0].(*toml.Tree).SetPath(args[1].([]string), args[2])
	p.PopN(3)
}

func execmTreeSetPathWithComment(_ int, p *gop.Context) {
	args := p.GetArgs(5)
	args[0].(*toml.Tree).SetPathWithComment(args[1].([]string), args[2].(string), args[3].(bool), args[4])
	p.PopN(5)
}

func execmTreeDelete(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*toml.Tree).Delete(args[1].(string))
	p.Ret(2, ret0)
}

func execmTreeDeletePath(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*toml.Tree).DeletePath(args[1].([]string))
	p.Ret(2, ret0)
}

func execmTreeValues(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*toml.Tree).Values()
	p.Ret(1, ret0)
}

func execmTreeComment(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*toml.Tree).Comment()
	p.Ret(1, ret0)
}

func execmTreeCommented(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*toml.Tree).Commented()
	p.Ret(1, ret0)
}

func execmTreeInline(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*toml.Tree).Inline()
	p.Ret(1, ret0)
}

func execmTreeSetValues(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*toml.Tree).SetValues(args[1].(map[string]interface{}))
	p.PopN(2)
}

func execmTreeSetComment(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*toml.Tree).SetComment(args[1].(string))
	p.PopN(2)
}

func execmTreeSetCommented(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*toml.Tree).SetCommented(args[1].(bool))
	p.PopN(2)
}

func execmTreeSetInline(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*toml.Tree).SetInline(args[1].(bool))
	p.PopN(2)
}

func execmTreeWriteTo(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(*toml.Tree).WriteTo(toType1(args[1]))
	p.Ret(2, ret0, ret1)
}

func execmTreeToTomlString(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := args[0].(*toml.Tree).ToTomlString()
	p.Ret(1, ret0, ret1)
}

func execmTreeString(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*toml.Tree).String()
	p.Ret(1, ret0)
}

func execmTreeToMap(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*toml.Tree).ToMap()
	p.Ret(1, ret0)
}

func execTreeFromMap(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := toml.TreeFromMap(args[0].(map[string]interface{}))
	p.Ret(1, ret0, ret1)
}

func execUnmarshal(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := toml.Unmarshal(args[0].([]byte), args[1])
	p.Ret(2, ret0)
}

func execiUnmarshalerUnmarshalTOML(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(toml.Unmarshaler).UnmarshalTOML(args[1])
	p.Ret(2, ret0)
}

func execValueStringRepresentation(_ int, p *gop.Context) {
	args := p.GetArgs(5)
	ret0, ret1 := toml.ValueStringRepresentation(args[0], args[1].(string), args[2].(string), args[3].(toml.MarshalOrder), args[4].(bool))
	p.Ret(5, ret0, ret1)
}

// I is a Go package instance.
var I = gop.NewGoPackage("github.com/pelletier/go-toml")

func init() {
	I.RegisterFuncs(
		I.Func("(*Decoder).Decode", (*toml.Decoder).Decode, execmDecoderDecode),
		I.Func("(*Decoder).SetTagName", (*toml.Decoder).SetTagName, execmDecoderSetTagName),
		I.Func("(*Decoder).Strict", (*toml.Decoder).Strict, execmDecoderStrict),
		I.Func("(*Encoder).Encode", (*toml.Encoder).Encode, execmEncoderEncode),
		I.Func("(*Encoder).QuoteMapKeys", (*toml.Encoder).QuoteMapKeys, execmEncoderQuoteMapKeys),
		I.Func("(*Encoder).ArraysWithOneElementPerLine", (*toml.Encoder).ArraysWithOneElementPerLine, execmEncoderArraysWithOneElementPerLine),
		I.Func("(*Encoder).Order", (*toml.Encoder).Order, execmEncoderOrder),
		I.Func("(*Encoder).Indentation", (*toml.Encoder).Indentation, execmEncoderIndentation),
		I.Func("(*Encoder).SetTagName", (*toml.Encoder).SetTagName, execmEncoderSetTagName),
		I.Func("(*Encoder).SetTagComment", (*toml.Encoder).SetTagComment, execmEncoderSetTagComment),
		I.Func("(*Encoder).SetTagCommented", (*toml.Encoder).SetTagCommented, execmEncoderSetTagCommented),
		I.Func("(*Encoder).SetTagMultiline", (*toml.Encoder).SetTagMultiline, execmEncoderSetTagMultiline),
		I.Func("(*Encoder).PromoteAnonymous", (*toml.Encoder).PromoteAnonymous, execmEncoderPromoteAnonymous),
		I.Func("Load", toml.Load, execLoad),
		I.Func("LoadBytes", toml.LoadBytes, execLoadBytes),
		I.Func("LoadFile", toml.LoadFile, execLoadFile),
		I.Func("LoadReader", toml.LoadReader, execLoadReader),
		I.Func("(LocalDate).String", (toml.LocalDate).String, execmLocalDateString),
		I.Func("(LocalDate).IsValid", (toml.LocalDate).IsValid, execmLocalDateIsValid),
		I.Func("(LocalDate).In", (toml.LocalDate).In, execmLocalDateIn),
		I.Func("(LocalDate).AddDays", (toml.LocalDate).AddDays, execmLocalDateAddDays),
		I.Func("(LocalDate).DaysSince", (toml.LocalDate).DaysSince, execmLocalDateDaysSince),
		I.Func("(LocalDate).Before", (toml.LocalDate).Before, execmLocalDateBefore),
		I.Func("(LocalDate).After", (toml.LocalDate).After, execmLocalDateAfter),
		I.Func("(LocalDate).MarshalText", (toml.LocalDate).MarshalText, execmLocalDateMarshalText),
		I.Func("(*LocalDate).UnmarshalText", (*toml.LocalDate).UnmarshalText, execmLocalDateUnmarshalText),
		I.Func("LocalDateOf", toml.LocalDateOf, execLocalDateOf),
		I.Func("(LocalDateTime).String", (toml.LocalDateTime).String, execmLocalDateTimeString),
		I.Func("(LocalDateTime).IsValid", (toml.LocalDateTime).IsValid, execmLocalDateTimeIsValid),
		I.Func("(LocalDateTime).In", (toml.LocalDateTime).In, execmLocalDateTimeIn),
		I.Func("(LocalDateTime).Before", (toml.LocalDateTime).Before, execmLocalDateTimeBefore),
		I.Func("(LocalDateTime).After", (toml.LocalDateTime).After, execmLocalDateTimeAfter),
		I.Func("(LocalDateTime).MarshalText", (toml.LocalDateTime).MarshalText, execmLocalDateTimeMarshalText),
		I.Func("(*LocalDateTime).UnmarshalText", (*toml.LocalDateTime).UnmarshalText, execmLocalDateTimeUnmarshalText),
		I.Func("LocalDateTimeOf", toml.LocalDateTimeOf, execLocalDateTimeOf),
		I.Func("(LocalTime).String", (toml.LocalTime).String, execmLocalTimeString),
		I.Func("(LocalTime).IsValid", (toml.LocalTime).IsValid, execmLocalTimeIsValid),
		I.Func("(LocalTime).MarshalText", (toml.LocalTime).MarshalText, execmLocalTimeMarshalText),
		I.Func("(*LocalTime).UnmarshalText", (*toml.LocalTime).UnmarshalText, execmLocalTimeUnmarshalText),
		I.Func("LocalTimeOf", toml.LocalTimeOf, execLocalTimeOf),
		I.Func("Marshal", toml.Marshal, execMarshal),
		I.Func("(Marshaler).MarshalTOML", (toml.Marshaler).MarshalTOML, execiMarshalerMarshalTOML),
		I.Func("NewDecoder", toml.NewDecoder, execNewDecoder),
		I.Func("NewEncoder", toml.NewEncoder, execNewEncoder),
		I.Func("ParseLocalDate", toml.ParseLocalDate, execParseLocalDate),
		I.Func("ParseLocalDateTime", toml.ParseLocalDateTime, execParseLocalDateTime),
		I.Func("ParseLocalTime", toml.ParseLocalTime, execParseLocalTime),
		I.Func("(Position).String", (toml.Position).String, execmPositionString),
		I.Func("(Position).Invalid", (toml.Position).Invalid, execmPositionInvalid),
		I.Func("(*PubTOMLValue).Value", (*toml.PubTOMLValue).Value, execmtomlValueValue),
		I.Func("(*PubTOMLValue).Comment", (*toml.PubTOMLValue).Comment, execmtomlValueComment),
		I.Func("(*PubTOMLValue).Commented", (*toml.PubTOMLValue).Commented, execmtomlValueCommented),
		I.Func("(*PubTOMLValue).Multiline", (*toml.PubTOMLValue).Multiline, execmtomlValueMultiline),
		I.Func("(*PubTOMLValue).Position", (*toml.PubTOMLValue).Position, execmtomlValuePosition),
		I.Func("(*PubTOMLValue).SetValue", (*toml.PubTOMLValue).SetValue, execmtomlValueSetValue),
		I.Func("(*PubTOMLValue).SetComment", (*toml.PubTOMLValue).SetComment, execmtomlValueSetComment),
		I.Func("(*PubTOMLValue).SetCommented", (*toml.PubTOMLValue).SetCommented, execmtomlValueSetCommented),
		I.Func("(*PubTOMLValue).SetMultiline", (*toml.PubTOMLValue).SetMultiline, execmtomlValueSetMultiline),
		I.Func("(*PubTOMLValue).SetPosition", (*toml.PubTOMLValue).SetPosition, execmtomlValueSetPosition),
		I.Func("(*Tree).Unmarshal", (*toml.Tree).Unmarshal, execmTreeUnmarshal),
		I.Func("(*Tree).Marshal", (*toml.Tree).Marshal, execmTreeMarshal),
		I.Func("(*Tree).Position", (*toml.Tree).Position, execmTreePosition),
		I.Func("(*Tree).Has", (*toml.Tree).Has, execmTreeHas),
		I.Func("(*Tree).HasPath", (*toml.Tree).HasPath, execmTreeHasPath),
		I.Func("(*Tree).Keys", (*toml.Tree).Keys, execmTreeKeys),
		I.Func("(*Tree).Get", (*toml.Tree).Get, execmTreeGet),
		I.Func("(*Tree).GetPath", (*toml.Tree).GetPath, execmTreeGetPath),
		I.Func("(*Tree).GetArray", (*toml.Tree).GetArray, execmTreeGetArray),
		I.Func("(*Tree).GetArrayPath", (*toml.Tree).GetArrayPath, execmTreeGetArrayPath),
		I.Func("(*Tree).GetPosition", (*toml.Tree).GetPosition, execmTreeGetPosition),
		I.Func("(*Tree).SetPositionPath", (*toml.Tree).SetPositionPath, execmTreeSetPositionPath),
		I.Func("(*Tree).GetPositionPath", (*toml.Tree).GetPositionPath, execmTreeGetPositionPath),
		I.Func("(*Tree).GetDefault", (*toml.Tree).GetDefault, execmTreeGetDefault),
		I.Func("(*Tree).SetWithOptions", (*toml.Tree).SetWithOptions, execmTreeSetWithOptions),
		I.Func("(*Tree).SetPathWithOptions", (*toml.Tree).SetPathWithOptions, execmTreeSetPathWithOptions),
		I.Func("(*Tree).Set", (*toml.Tree).Set, execmTreeSet),
		I.Func("(*Tree).SetWithComment", (*toml.Tree).SetWithComment, execmTreeSetWithComment),
		I.Func("(*Tree).SetPath", (*toml.Tree).SetPath, execmTreeSetPath),
		I.Func("(*Tree).SetPathWithComment", (*toml.Tree).SetPathWithComment, execmTreeSetPathWithComment),
		I.Func("(*Tree).Delete", (*toml.Tree).Delete, execmTreeDelete),
		I.Func("(*Tree).DeletePath", (*toml.Tree).DeletePath, execmTreeDeletePath),
		I.Func("(*Tree).Values", (*toml.Tree).Values, execmTreeValues),
		I.Func("(*Tree).Comment", (*toml.Tree).Comment, execmTreeComment),
		I.Func("(*Tree).Commented", (*toml.Tree).Commented, execmTreeCommented),
		I.Func("(*Tree).Inline", (*toml.Tree).Inline, execmTreeInline),
		I.Func("(*Tree).SetValues", (*toml.Tree).SetValues, execmTreeSetValues),
		I.Func("(*Tree).SetComment", (*toml.Tree).SetComment, execmTreeSetComment),
		I.Func("(*Tree).SetCommented", (*toml.Tree).SetCommented, execmTreeSetCommented),
		I.Func("(*Tree).SetInline", (*toml.Tree).SetInline, execmTreeSetInline),
		I.Func("(*Tree).WriteTo", (*toml.Tree).WriteTo, execmTreeWriteTo),
		I.Func("(*Tree).ToTomlString", (*toml.Tree).ToTomlString, execmTreeToTomlString),
		I.Func("(*Tree).String", (*toml.Tree).String, execmTreeString),
		I.Func("(*Tree).ToMap", (*toml.Tree).ToMap, execmTreeToMap),
		I.Func("TreeFromMap", toml.TreeFromMap, execTreeFromMap),
		I.Func("Unmarshal", toml.Unmarshal, execUnmarshal),
		I.Func("(Unmarshaler).UnmarshalTOML", (toml.Unmarshaler).UnmarshalTOML, execiUnmarshalerUnmarshalTOML),
		I.Func("ValueStringRepresentation", toml.ValueStringRepresentation, execValueStringRepresentation),
	)
	I.RegisterTypes(
		I.Type("Decoder", reflect.TypeOf((*toml.Decoder)(nil)).Elem()),
		I.Type("Encoder", reflect.TypeOf((*toml.Encoder)(nil)).Elem()),
		I.Type("LocalDate", reflect.TypeOf((*toml.LocalDate)(nil)).Elem()),
		I.Type("LocalDateTime", reflect.TypeOf((*toml.LocalDateTime)(nil)).Elem()),
		I.Type("LocalTime", reflect.TypeOf((*toml.LocalTime)(nil)).Elem()),
		I.Type("MarshalOrder", reflect.TypeOf((*toml.MarshalOrder)(nil)).Elem()),
		I.Type("Marshaler", reflect.TypeOf((*toml.Marshaler)(nil)).Elem()),
		I.Type("Position", reflect.TypeOf((*toml.Position)(nil)).Elem()),
		I.Type("PubTOMLValue", reflect.TypeOf((*toml.PubTOMLValue)(nil)).Elem()),
		I.Type("PubTree", reflect.TypeOf((*toml.PubTree)(nil)).Elem()),
		I.Type("SetOptions", reflect.TypeOf((*toml.SetOptions)(nil)).Elem()),
		I.Type("Tree", reflect.TypeOf((*toml.Tree)(nil)).Elem()),
		I.Type("Unmarshaler", reflect.TypeOf((*toml.Unmarshaler)(nil)).Elem()),
	)
	I.RegisterConsts(
		I.Const("OrderAlphabetical", qspec.Int, toml.OrderAlphabetical),
		I.Const("OrderPreserve", qspec.Int, toml.OrderPreserve),
	)
}
