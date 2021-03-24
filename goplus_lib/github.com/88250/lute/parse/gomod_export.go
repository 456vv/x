// Package parse provide Go+ "github.com/88250/lute/parse" package, as "github.com/88250/lute/parse" package in Go.
package parse

import (
	reflect "reflect"

	ast "github.com/88250/lute/ast"
	parse "github.com/88250/lute/parse"
	gop "github.com/goplus/gop"
	qspec "github.com/goplus/gop/exec.spec"
)

func execAddAutoLinkDomainSuffix(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	parse.AddAutoLinkDomainSuffix(args[0].(string))
	p.PopN(1)
}

func execBlockquoteContinue(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := parse.BlockquoteContinue(args[0].(*ast.Node), args[1].(*parse.Context))
	p.Ret(2, ret0)
}

func execCodeBlockContinue(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := parse.CodeBlockContinue(args[0].(*ast.Node), args[1].(*parse.Context))
	p.Ret(2, ret0)
}

func execmContextParentTip(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	args[0].(*parse.Context).ParentTip()
	p.PopN(1)
}

func execmContextTipAppendChild(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*parse.Context).TipAppendChild(args[1].(*ast.Node))
	p.PopN(2)
}

func execFootnotesContinue(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := parse.FootnotesContinue(args[0].(*ast.Node), args[1].(*parse.Context))
	p.Ret(2, ret0)
}

func execGitConflictContinue(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := parse.GitConflictContinue(args[0].(*ast.Node), args[1].(*parse.Context))
	p.Ret(2, ret0)
}

func execHtmlBlockContinue(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := parse.HtmlBlockContinue(args[0].(*ast.Node), args[1].(*parse.Context))
	p.Ret(2, ret0)
}

func execIAL2Tokens(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := parse.IAL2Tokens(args[0].([][]string))
	p.Ret(1, ret0)
}

func execInline(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := parse.Inline(args[0].(string), args[1].([]byte), args[2].(*parse.Options))
	p.Ret(3, ret0)
}

func execListItemContinue(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := parse.ListItemContinue(args[0].(*ast.Node), args[1].(*parse.Context))
	p.Ret(2, ret0)
}

func execMathBlockContinue(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := parse.MathBlockContinue(args[0].(*ast.Node), args[1].(*parse.Context))
	p.Ret(2, ret0)
}

func execNewOptions(_ int, p *gop.Context) {
	ret0 := parse.NewOptions()
	p.Ret(0, ret0)
}

func execParagraphContinue(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := parse.ParagraphContinue(args[0].(*ast.Node), args[1].(*parse.Context))
	p.Ret(2, ret0)
}

func execParse(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := parse.Parse(args[0].(string), args[1].([]byte), args[2].(*parse.Options))
	p.Ret(3, ret0)
}

func execSuperBlockContinue(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := parse.SuperBlockContinue(args[0].(*ast.Node), args[1].(*parse.Context))
	p.Ret(2, ret0)
}

func execmTreeEmojiImgTokens(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*parse.Tree).EmojiImgTokens(args[1].(string), args[2].(string))
	p.Ret(3, ret0)
}

func execmTreeFindFootnotesDef(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(*parse.Tree).FindFootnotesDef(args[1].([]byte))
	p.Ret(2, ret0, ret1)
}

func execmTreeFindLinkRefDefLink(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*parse.Tree).FindLinkRefDefLink(args[1].([]byte))
	p.Ret(2, ret0)
}

func execYamlFrontMatterContinue(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := parse.YamlFrontMatterContinue(args[0].(*ast.Node), args[1].(*parse.Context))
	p.Ret(2, ret0)
}

// I is a Go package instance.
var I = gop.NewGoPackage("github.com/88250/lute/parse")

func init() {
	I.RegisterFuncs(
		I.Func("AddAutoLinkDomainSuffix", parse.AddAutoLinkDomainSuffix, execAddAutoLinkDomainSuffix),
		I.Func("BlockquoteContinue", parse.BlockquoteContinue, execBlockquoteContinue),
		I.Func("CodeBlockContinue", parse.CodeBlockContinue, execCodeBlockContinue),
		I.Func("(*Context).ParentTip", (*parse.Context).ParentTip, execmContextParentTip),
		I.Func("(*Context).TipAppendChild", (*parse.Context).TipAppendChild, execmContextTipAppendChild),
		I.Func("FootnotesContinue", parse.FootnotesContinue, execFootnotesContinue),
		I.Func("GitConflictContinue", parse.GitConflictContinue, execGitConflictContinue),
		I.Func("HtmlBlockContinue", parse.HtmlBlockContinue, execHtmlBlockContinue),
		I.Func("IAL2Tokens", parse.IAL2Tokens, execIAL2Tokens),
		I.Func("Inline", parse.Inline, execInline),
		I.Func("ListItemContinue", parse.ListItemContinue, execListItemContinue),
		I.Func("MathBlockContinue", parse.MathBlockContinue, execMathBlockContinue),
		I.Func("NewOptions", parse.NewOptions, execNewOptions),
		I.Func("ParagraphContinue", parse.ParagraphContinue, execParagraphContinue),
		I.Func("Parse", parse.Parse, execParse),
		I.Func("SuperBlockContinue", parse.SuperBlockContinue, execSuperBlockContinue),
		I.Func("(*Tree).EmojiImgTokens", (*parse.Tree).EmojiImgTokens, execmTreeEmojiImgTokens),
		I.Func("(*Tree).FindFootnotesDef", (*parse.Tree).FindFootnotesDef, execmTreeFindFootnotesDef),
		I.Func("(*Tree).FindLinkRefDefLink", (*parse.Tree).FindLinkRefDefLink, execmTreeFindLinkRefDefLink),
		I.Func("YamlFrontMatterContinue", parse.YamlFrontMatterContinue, execYamlFrontMatterContinue),
	)
	I.RegisterVars(
		I.Var("EmojiAliasUnicode", &parse.EmojiAliasUnicode),
		I.Var("EmojiSitePlaceholder", &parse.EmojiSitePlaceholder),
		I.Var("EmojiUnicodeAlias", &parse.EmojiUnicodeAlias),
		I.Var("MathBlockMarker", &parse.MathBlockMarker),
		I.Var("MathBlockMarkerCaret", &parse.MathBlockMarkerCaret),
		I.Var("MathBlockMarkerCaretNewline", &parse.MathBlockMarkerCaretNewline),
		I.Var("MathBlockMarkerNewline", &parse.MathBlockMarkerNewline),
		I.Var("YamlFrontMatterMarker", &parse.YamlFrontMatterMarker),
		I.Var("YamlFrontMatterMarkerCaret", &parse.YamlFrontMatterMarkerCaret),
		I.Var("YamlFrontMatterMarkerCaretNewline", &parse.YamlFrontMatterMarkerCaretNewline),
		I.Var("YamlFrontMatterMarkerNewline", &parse.YamlFrontMatterMarkerNewline),
	)
	I.RegisterTypes(
		I.Type("Context", reflect.TypeOf((*parse.Context)(nil)).Elem()),
		I.Type("InlineContext", reflect.TypeOf((*parse.InlineContext)(nil)).Elem()),
		I.Type("Options", reflect.TypeOf((*parse.Options)(nil)).Elem()),
		I.Type("Tree", reflect.TypeOf((*parse.Tree)(nil)).Elem()),
	)
	I.RegisterConsts(
		I.Const("Zwsp", qspec.ConstBoundString, parse.Zwsp),
	)
}
