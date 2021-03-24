// Package html provide Go+ "github.com/88250/lute/html" package, as "github.com/88250/lute/html" package in Go.
package html

import (
	io "io"
	reflect "reflect"

	html "github.com/88250/lute/html"
	gop "github.com/goplus/gop"
	qspec "github.com/goplus/gop/exec.spec"
)

func execEncodeDestination(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := html.EncodeDestination(args[0].([]byte))
	p.Ret(1, ret0)
}

func execEscapeHTML(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := html.EscapeHTML(args[0].([]byte))
	p.Ret(1, ret0)
}

func execEscapeString(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := html.EscapeString(args[0].(string))
	p.Ret(1, ret0)
}

func execHtmlUnescapeString(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := html.HtmlUnescapeString(args[0].(string))
	p.Ret(1, ret0)
}

func toType0(v interface{}) io.Reader {
	if v == nil {
		return nil
	}
	return v.(io.Reader)
}

func execNewTokenizer(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := html.NewTokenizer(toType0(args[0]))
	p.Ret(1, ret0)
}

func execNewTokenizerFragment(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := html.NewTokenizerFragment(toType0(args[0]), args[1].(string))
	p.Ret(2, ret0)
}

func execmNodeUnlink(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	args[0].(*html.Node).Unlink()
	p.PopN(1)
}

func execmNodeInsertBefore(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*html.Node).InsertBefore(args[1].(*html.Node))
	p.PopN(2)
}

func execmNodeInsertAfter(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*html.Node).InsertAfter(args[1].(*html.Node))
	p.PopN(2)
}

func execmNodeInsertChildBefore(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	args[0].(*html.Node).InsertChildBefore(args[1].(*html.Node), args[2].(*html.Node))
	p.PopN(3)
}

func execmNodeAppendChild(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*html.Node).AppendChild(args[1].(*html.Node))
	p.PopN(2)
}

func execmNodeRemoveChild(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*html.Node).RemoveChild(args[1].(*html.Node))
	p.PopN(2)
}

func execParse(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := html.Parse(toType0(args[0]))
	p.Ret(1, ret0, ret1)
}

func execParseFragment(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := html.ParseFragment(toType0(args[0]), args[1].(*html.Node))
	p.Ret(2, ret0, ret1)
}

func toSlice0(args []interface{}) []html.ParseOption {
	ret := make([]html.ParseOption, len(args))
	for i, arg := range args {
		ret[i] = arg.(html.ParseOption)
	}
	return ret
}

func execParseFragmentWithOptions(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	ret0, ret1 := html.ParseFragmentWithOptions(toType0(args[0]), args[1].(*html.Node), toSlice0(args[2:])...)
	p.Ret(arity, ret0, ret1)
}

func execParseOptionEnableScripting(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := html.ParseOptionEnableScripting(args[0].(bool))
	p.Ret(1, ret0)
}

func execParseWithOptions(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	ret0, ret1 := html.ParseWithOptions(toType0(args[0]), toSlice0(args[1:])...)
	p.Ret(arity, ret0, ret1)
}

func toType1(v interface{}) io.Writer {
	if v == nil {
		return nil
	}
	return v.(io.Writer)
}

func execRender(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := html.Render(toType1(args[0]), args[1].(*html.Node))
	p.Ret(2, ret0)
}

func execmTokenString(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(html.Token).String()
	p.Ret(1, ret0)
}

func execmTokenTypeString(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(html.TokenType).String()
	p.Ret(1, ret0)
}

func execmTokenizerAllowCDATA(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*html.Tokenizer).AllowCDATA(args[1].(bool))
	p.PopN(2)
}

func execmTokenizerNextIsNotRawText(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	args[0].(*html.Tokenizer).NextIsNotRawText()
	p.PopN(1)
}

func execmTokenizerErr(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*html.Tokenizer).Err()
	p.Ret(1, ret0)
}

func execmTokenizerBuffered(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*html.Tokenizer).Buffered()
	p.Ret(1, ret0)
}

func execmTokenizerNext(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*html.Tokenizer).Next()
	p.Ret(1, ret0)
}

func execmTokenizerRaw(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*html.Tokenizer).Raw()
	p.Ret(1, ret0)
}

func execmTokenizerText(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*html.Tokenizer).Text()
	p.Ret(1, ret0)
}

func execmTokenizerTagName(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := args[0].(*html.Tokenizer).TagName()
	p.Ret(1, ret0, ret1)
}

func execmTokenizerTagAttr(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1, ret2 := args[0].(*html.Tokenizer).TagAttr()
	p.Ret(1, ret0, ret1, ret2)
}

func execmTokenizerToken(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*html.Tokenizer).Token()
	p.Ret(1, ret0)
}

func execmTokenizerSetMaxBuf(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*html.Tokenizer).SetMaxBuf(args[1].(int))
	p.PopN(2)
}

func execUnescapeBytes(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := html.UnescapeBytes(args[0].([]byte))
	p.Ret(1, ret0)
}

func execUnescapeHTML(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := html.UnescapeHTML(args[0].([]byte))
	p.Ret(1, ret0)
}

func execUnescapeString(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := html.UnescapeString(args[0].(string))
	p.Ret(1, ret0)
}

// I is a Go package instance.
var I = gop.NewGoPackage("github.com/88250/lute/html")

func init() {
	I.RegisterFuncs(
		I.Func("EncodeDestination", html.EncodeDestination, execEncodeDestination),
		I.Func("EscapeHTML", html.EscapeHTML, execEscapeHTML),
		I.Func("EscapeString", html.EscapeString, execEscapeString),
		I.Func("HtmlUnescapeString", html.HtmlUnescapeString, execHtmlUnescapeString),
		I.Func("NewTokenizer", html.NewTokenizer, execNewTokenizer),
		I.Func("NewTokenizerFragment", html.NewTokenizerFragment, execNewTokenizerFragment),
		I.Func("(*Node).Unlink", (*html.Node).Unlink, execmNodeUnlink),
		I.Func("(*Node).InsertBefore", (*html.Node).InsertBefore, execmNodeInsertBefore),
		I.Func("(*Node).InsertAfter", (*html.Node).InsertAfter, execmNodeInsertAfter),
		I.Func("(*Node).InsertChildBefore", (*html.Node).InsertChildBefore, execmNodeInsertChildBefore),
		I.Func("(*Node).AppendChild", (*html.Node).AppendChild, execmNodeAppendChild),
		I.Func("(*Node).RemoveChild", (*html.Node).RemoveChild, execmNodeRemoveChild),
		I.Func("Parse", html.Parse, execParse),
		I.Func("ParseFragment", html.ParseFragment, execParseFragment),
		I.Func("ParseOptionEnableScripting", html.ParseOptionEnableScripting, execParseOptionEnableScripting),
		I.Func("Render", html.Render, execRender),
		I.Func("(Token).String", (html.Token).String, execmTokenString),
		I.Func("(TokenType).String", (html.TokenType).String, execmTokenTypeString),
		I.Func("(*Tokenizer).AllowCDATA", (*html.Tokenizer).AllowCDATA, execmTokenizerAllowCDATA),
		I.Func("(*Tokenizer).NextIsNotRawText", (*html.Tokenizer).NextIsNotRawText, execmTokenizerNextIsNotRawText),
		I.Func("(*Tokenizer).Err", (*html.Tokenizer).Err, execmTokenizerErr),
		I.Func("(*Tokenizer).Buffered", (*html.Tokenizer).Buffered, execmTokenizerBuffered),
		I.Func("(*Tokenizer).Next", (*html.Tokenizer).Next, execmTokenizerNext),
		I.Func("(*Tokenizer).Raw", (*html.Tokenizer).Raw, execmTokenizerRaw),
		I.Func("(*Tokenizer).Text", (*html.Tokenizer).Text, execmTokenizerText),
		I.Func("(*Tokenizer).TagName", (*html.Tokenizer).TagName, execmTokenizerTagName),
		I.Func("(*Tokenizer).TagAttr", (*html.Tokenizer).TagAttr, execmTokenizerTagAttr),
		I.Func("(*Tokenizer).Token", (*html.Tokenizer).Token, execmTokenizerToken),
		I.Func("(*Tokenizer).SetMaxBuf", (*html.Tokenizer).SetMaxBuf, execmTokenizerSetMaxBuf),
		I.Func("UnescapeBytes", html.UnescapeBytes, execUnescapeBytes),
		I.Func("UnescapeHTML", html.UnescapeHTML, execUnescapeHTML),
		I.Func("UnescapeString", html.UnescapeString, execUnescapeString),
	)
	I.RegisterFuncvs(
		I.Funcv("ParseFragmentWithOptions", html.ParseFragmentWithOptions, execParseFragmentWithOptions),
		I.Funcv("ParseWithOptions", html.ParseWithOptions, execParseWithOptions),
	)
	I.RegisterVars(
		I.Var("Entities", &html.Entities),
		I.Var("ErrBufferExceeded", &html.ErrBufferExceeded),
	)
	I.RegisterTypes(
		I.Type("Attribute", reflect.TypeOf((*html.Attribute)(nil)).Elem()),
		I.Type("Node", reflect.TypeOf((*html.Node)(nil)).Elem()),
		I.Type("NodeType", reflect.TypeOf((*html.NodeType)(nil)).Elem()),
		I.Type("ParseOption", reflect.TypeOf((*html.ParseOption)(nil)).Elem()),
		I.Type("Token", reflect.TypeOf((*html.Token)(nil)).Elem()),
		I.Type("TokenType", reflect.TypeOf((*html.TokenType)(nil)).Elem()),
		I.Type("Tokenizer", reflect.TypeOf((*html.Tokenizer)(nil)).Elem()),
	)
	I.RegisterConsts(
		I.Const("BadEntity", qspec.String, html.BadEntity),
		I.Const("CommentNode", qspec.Uint32, html.CommentNode),
		I.Const("CommentToken", qspec.Uint32, html.CommentToken),
		I.Const("DoctypeNode", qspec.Uint32, html.DoctypeNode),
		I.Const("DoctypeToken", qspec.Uint32, html.DoctypeToken),
		I.Const("DocumentNode", qspec.Uint32, html.DocumentNode),
		I.Const("ElementNode", qspec.Uint32, html.ElementNode),
		I.Const("EndTagToken", qspec.Uint32, html.EndTagToken),
		I.Const("ErrorNode", qspec.Uint32, html.ErrorNode),
		I.Const("ErrorToken", qspec.Uint32, html.ErrorToken),
		I.Const("SelfClosingTagToken", qspec.Uint32, html.SelfClosingTagToken),
		I.Const("StartTagToken", qspec.Uint32, html.StartTagToken),
		I.Const("TextNode", qspec.Uint32, html.TextNode),
		I.Const("TextToken", qspec.Uint32, html.TextToken),
	)
}
