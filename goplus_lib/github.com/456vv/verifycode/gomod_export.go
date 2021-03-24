// Package verifycode provide Go+ "github.com/456vv/verifycode" package, as "github.com/456vv/verifycode" package in Go.
package verifycode

import (
	color "image/color"
	gif "image/gif"
	jpeg "image/jpeg"
	io "io"
	reflect "reflect"

	verifycode "github.com/456vv/verifycode"
	truetype "github.com/golang/freetype/truetype"
	gop "github.com/goplus/gop"
)

func execmColorAddHEX(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*verifycode.Color).AddHEX(args[1].(string))
	p.Ret(2, ret0)
}

func execmColorAddRGBA(_ int, p *gop.Context) {
	args := p.GetArgs(5)
	ret0 := args[0].(*verifycode.Color).AddRGBA(args[1].(uint8), args[2].(uint8), args[3].(uint8), args[4].(uint8))
	p.Ret(5, ret0)
}

func execmColorRandom(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*verifycode.Color).Random()
	p.Ret(1, ret0)
}

func execmFontAddFile(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*verifycode.Font).AddFile(args[1].(string))
	p.Ret(2, ret0)
}

func execmFontRandom(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := args[0].(*verifycode.Font).Random()
	p.Ret(1, ret0, ret1)
}

func toType0(v interface{}) color.Color {
	if v == nil {
		return nil
	}
	return v.(color.Color)
}

func execmGlyphFontGlyph(_ int, p *gop.Context) {
	args := p.GetArgs(4)
	ret0, ret1 := args[0].(*verifycode.Glyph).FontGlyph(args[1].(*truetype.Font), args[2].(rune), toType0(args[3]))
	p.Ret(4, ret0, ret1)
}

func execRand(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := verifycode.Rand(args[0].(int))
	p.Ret(1, ret0)
}

func execRandRange(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := verifycode.RandRange(args[0].(int), args[1].(int))
	p.Ret(2, ret0)
}

func execRandomText(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := verifycode.RandomText(args[0].(string), args[1].(int))
	p.Ret(2, ret0)
}

func execmVerifyCodeStyle(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*verifycode.VerifyCode).Style(args[1].(*verifycode.Style))
	p.Ret(2, ret0)
}

func execmVerifyCodeDraw(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(*verifycode.VerifyCode).Draw(args[1].(string))
	p.Ret(2, ret0, ret1)
}

func toType1(v interface{}) io.Writer {
	if v == nil {
		return nil
	}
	return v.(io.Writer)
}

func execmVerifyCodePNG(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*verifycode.VerifyCode).PNG(args[1].(string), toType1(args[2]))
	p.Ret(3, ret0)
}

func execmVerifyCodeGIF(_ int, p *gop.Context) {
	args := p.GetArgs(4)
	ret0 := args[0].(*verifycode.VerifyCode).GIF(args[1].(string), toType1(args[2]), args[3].(*gif.Options))
	p.Ret(4, ret0)
}

func execmVerifyCodeJPEG(_ int, p *gop.Context) {
	args := p.GetArgs(4)
	ret0 := args[0].(*verifycode.VerifyCode).JPEG(args[1].(string), toType1(args[2]), args[3].(*jpeg.Options))
	p.Ret(4, ret0)
}

// I is a Go package instance.
var I = gop.NewGoPackage("github.com/456vv/verifycode")

func init() {
	I.RegisterFuncs(
		I.Func("(*Color).AddHEX", (*verifycode.Color).AddHEX, execmColorAddHEX),
		I.Func("(*Color).AddRGBA", (*verifycode.Color).AddRGBA, execmColorAddRGBA),
		I.Func("(*Color).Random", (*verifycode.Color).Random, execmColorRandom),
		I.Func("(*Font).AddFile", (*verifycode.Font).AddFile, execmFontAddFile),
		I.Func("(*Font).Random", (*verifycode.Font).Random, execmFontRandom),
		I.Func("(*Glyph).FontGlyph", (*verifycode.Glyph).FontGlyph, execmGlyphFontGlyph),
		I.Func("Rand", verifycode.Rand, execRand),
		I.Func("RandRange", verifycode.RandRange, execRandRange),
		I.Func("RandomText", verifycode.RandomText, execRandomText),
		I.Func("(*VerifyCode).Style", (*verifycode.VerifyCode).Style, execmVerifyCodeStyle),
		I.Func("(*VerifyCode).Draw", (*verifycode.VerifyCode).Draw, execmVerifyCodeDraw),
		I.Func("(*VerifyCode).PNG", (*verifycode.VerifyCode).PNG, execmVerifyCodePNG),
		I.Func("(*VerifyCode).GIF", (*verifycode.VerifyCode).GIF, execmVerifyCodeGIF),
		I.Func("(*VerifyCode).JPEG", (*verifycode.VerifyCode).JPEG, execmVerifyCodeJPEG),
	)
	I.RegisterTypes(
		I.Type("Color", reflect.TypeOf((*verifycode.Color)(nil)).Elem()),
		I.Type("Font", reflect.TypeOf((*verifycode.Font)(nil)).Elem()),
		I.Type("Glyph", reflect.TypeOf((*verifycode.Glyph)(nil)).Elem()),
		I.Type("Style", reflect.TypeOf((*verifycode.Style)(nil)).Elem()),
		I.Type("VerifyCode", reflect.TypeOf((*verifycode.VerifyCode)(nil)).Elem()),
	)
}
