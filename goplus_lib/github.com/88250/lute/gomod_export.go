// Package lute provide Go+ "github.com/88250/lute" package, as "github.com/88250/lute" package in Go.
package lute

import (
	reflect "reflect"

	lute "github.com/88250/lute"
	ast "github.com/88250/lute/ast"
	parse "github.com/88250/lute/parse"
	render "github.com/88250/lute/render"
	js "github.com/gopherjs/gopherjs/js"
	gop "github.com/goplus/gop"
	qspec "github.com/goplus/gop/exec.spec"
)

func execFormatNode(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := lute.FormatNode(args[0].(*ast.Node), args[1].(*parse.Options), args[2].(*render.Options))
	p.Ret(3, ret0)
}

func execmLuteHTML2Markdown(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(*lute.Lute).HTML2Markdown(args[1].(string))
	p.Ret(2, ret0, ret1)
}

func execmLuteHTML2Tree(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*lute.Lute).HTML2Tree(args[1].(string))
	p.Ret(2, ret0)
}

func execmLuteMarkdown(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*lute.Lute).Markdown(args[1].(string), args[2].([]byte))
	p.Ret(3, ret0)
}

func execmLuteMarkdownStr(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*lute.Lute).MarkdownStr(args[1].(string), args[2].(string))
	p.Ret(3, ret0)
}

func execmLuteFormat(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*lute.Lute).Format(args[1].(string), args[2].([]byte))
	p.Ret(3, ret0)
}

func execmLuteFormatStr(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*lute.Lute).FormatStr(args[1].(string), args[2].(string))
	p.Ret(3, ret0)
}

func execmLuteTextBundle(_ int, p *gop.Context) {
	args := p.GetArgs(4)
	ret0, ret1 := args[0].(*lute.Lute).TextBundle(args[1].(string), args[2].([]byte), args[3].([]string))
	p.Ret(4, ret0, ret1)
}

func execmLuteTextBundleStr(_ int, p *gop.Context) {
	args := p.GetArgs(4)
	ret0, ret1 := args[0].(*lute.Lute).TextBundleStr(args[1].(string), args[2].(string), args[3].([]string))
	p.Ret(4, ret0, ret1)
}

func execmLuteHTML2Text(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*lute.Lute).HTML2Text(args[1].(string))
	p.Ret(2, ret0)
}

func execmLuteRenderJSON(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*lute.Lute).RenderJSON(args[1].(string))
	p.Ret(2, ret0)
}

func execmLuteSpace(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*lute.Lute).Space(args[1].(string))
	p.Ret(2, ret0)
}

func execmLuteIsValidLinkDest(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*lute.Lute).IsValidLinkDest(args[1].(string))
	p.Ret(2, ret0)
}

func execmLuteGetEmojis(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*lute.Lute).GetEmojis()
	p.Ret(1, ret0)
}

func execmLutePutEmojis(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*lute.Lute).PutEmojis(args[1].(map[string]string))
	p.PopN(2)
}

func execmLuteRemoveEmoji(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*lute.Lute).RemoveEmoji(args[1].(string))
	p.Ret(2, ret0)
}

func execmLuteGetTerms(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*lute.Lute).GetTerms()
	p.Ret(1, ret0)
}

func execmLutePutTerms(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*lute.Lute).PutTerms(args[1].(map[string]string))
	p.PopN(2)
}

func execmLuteTree2HTML(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*lute.Lute).Tree2HTML(args[1].(*parse.Tree), args[2].(*render.Options))
	p.Ret(3, ret0)
}

func execmLuteSetGFMTable(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*lute.Lute).SetGFMTable(args[1].(bool))
	p.PopN(2)
}

func execmLuteSetGFMTaskListItem(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*lute.Lute).SetGFMTaskListItem(args[1].(bool))
	p.PopN(2)
}

func execmLuteSetGFMTaskListItemClass(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*lute.Lute).SetGFMTaskListItemClass(args[1].(string))
	p.PopN(2)
}

func execmLuteSetGFMStrikethrough(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*lute.Lute).SetGFMStrikethrough(args[1].(bool))
	p.PopN(2)
}

func execmLuteSetGFMAutoLink(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*lute.Lute).SetGFMAutoLink(args[1].(bool))
	p.PopN(2)
}

func execmLuteSetSoftBreak2HardBreak(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*lute.Lute).SetSoftBreak2HardBreak(args[1].(bool))
	p.PopN(2)
}

func execmLuteSetCodeSyntaxHighlight(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*lute.Lute).SetCodeSyntaxHighlight(args[1].(bool))
	p.PopN(2)
}

func execmLuteSetCodeSyntaxHighlightDetectLang(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*lute.Lute).SetCodeSyntaxHighlightDetectLang(args[1].(bool))
	p.PopN(2)
}

func execmLuteSetCodeSyntaxHighlightInlineStyle(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*lute.Lute).SetCodeSyntaxHighlightInlineStyle(args[1].(bool))
	p.PopN(2)
}

func execmLuteSetCodeSyntaxHighlightLineNum(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*lute.Lute).SetCodeSyntaxHighlightLineNum(args[1].(bool))
	p.PopN(2)
}

func execmLuteSetCodeSyntaxHighlightStyleName(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*lute.Lute).SetCodeSyntaxHighlightStyleName(args[1].(string))
	p.PopN(2)
}

func execmLuteSetFootnotes(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*lute.Lute).SetFootnotes(args[1].(bool))
	p.PopN(2)
}

func execmLuteSetToC(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*lute.Lute).SetToC(args[1].(bool))
	p.PopN(2)
}

func execmLuteSetHeadingID(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*lute.Lute).SetHeadingID(args[1].(bool))
	p.PopN(2)
}

func execmLuteSetAutoSpace(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*lute.Lute).SetAutoSpace(args[1].(bool))
	p.PopN(2)
}

func execmLuteSetFixTermTypo(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*lute.Lute).SetFixTermTypo(args[1].(bool))
	p.PopN(2)
}

func execmLuteSetEmoji(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*lute.Lute).SetEmoji(args[1].(bool))
	p.PopN(2)
}

func execmLuteSetEmojis(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*lute.Lute).SetEmojis(args[1].(map[string]string))
	p.PopN(2)
}

func execmLuteSetEmojiSite(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*lute.Lute).SetEmojiSite(args[1].(string))
	p.PopN(2)
}

func execmLuteSetHeadingAnchor(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*lute.Lute).SetHeadingAnchor(args[1].(bool))
	p.PopN(2)
}

func execmLuteSetTerms(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*lute.Lute).SetTerms(args[1].(map[string]string))
	p.PopN(2)
}

func execmLuteSetVditorWYSIWYG(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*lute.Lute).SetVditorWYSIWYG(args[1].(bool))
	p.PopN(2)
}

func execmLuteSetVditorIR(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*lute.Lute).SetVditorIR(args[1].(bool))
	p.PopN(2)
}

func execmLuteSetVditorSV(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*lute.Lute).SetVditorSV(args[1].(bool))
	p.PopN(2)
}

func execmLuteSetInlineMathAllowDigitAfterOpenMarker(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*lute.Lute).SetInlineMathAllowDigitAfterOpenMarker(args[1].(bool))
	p.PopN(2)
}

func execmLuteSetLinkPrefix(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*lute.Lute).SetLinkPrefix(args[1].(string))
	p.PopN(2)
}

func execmLuteSetLinkBase(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*lute.Lute).SetLinkBase(args[1].(string))
	p.PopN(2)
}

func execmLuteGetLinkBase(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*lute.Lute).GetLinkBase()
	p.Ret(1, ret0)
}

func execmLuteSetVditorCodeBlockPreview(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*lute.Lute).SetVditorCodeBlockPreview(args[1].(bool))
	p.PopN(2)
}

func execmLuteSetVditorMathBlockPreview(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*lute.Lute).SetVditorMathBlockPreview(args[1].(bool))
	p.PopN(2)
}

func execmLuteSetVditorHTMLBlockPreview(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*lute.Lute).SetVditorHTMLBlockPreview(args[1].(bool))
	p.PopN(2)
}

func execmLuteSetRenderListStyle(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*lute.Lute).SetRenderListStyle(args[1].(bool))
	p.PopN(2)
}

func execmLuteSetSanitize(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*lute.Lute).SetSanitize(args[1].(bool))
	p.PopN(2)
}

func execmLuteSetImageLazyLoading(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*lute.Lute).SetImageLazyLoading(args[1].(string))
	p.PopN(2)
}

func execmLuteSetChineseParagraphBeginningSpace(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*lute.Lute).SetChineseParagraphBeginningSpace(args[1].(bool))
	p.PopN(2)
}

func execmLuteSetYamlFrontMatter(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*lute.Lute).SetYamlFrontMatter(args[1].(bool))
	p.PopN(2)
}

func execmLuteSetBlockRef(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*lute.Lute).SetBlockRef(args[1].(bool))
	p.PopN(2)
}

func execmLuteSetMark(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*lute.Lute).SetMark(args[1].(bool))
	p.PopN(2)
}

func execmLuteSetKramdownIAL(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*lute.Lute).SetKramdownIAL(args[1].(bool))
	p.PopN(2)
}

func execmLuteSetKramdownBlockIAL(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*lute.Lute).SetKramdownBlockIAL(args[1].(bool))
	p.PopN(2)
}

func execmLuteSetKramdownSpanIAL(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*lute.Lute).SetKramdownSpanIAL(args[1].(bool))
	p.PopN(2)
}

func execmLuteSetKramdownIALIDRenderName(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*lute.Lute).SetKramdownIALIDRenderName(args[1].(string))
	p.PopN(2)
}

func execmLuteSetTag(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*lute.Lute).SetTag(args[1].(bool))
	p.PopN(2)
}

func execmLuteSetImgPathAllowSpace(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*lute.Lute).SetImgPathAllowSpace(args[1].(bool))
	p.PopN(2)
}

func execmLuteSetSuperBlock(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*lute.Lute).SetSuperBlock(args[1].(bool))
	p.PopN(2)
}

func execmLuteSetSup(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*lute.Lute).SetSup(args[1].(bool))
	p.PopN(2)
}

func execmLuteSetSub(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*lute.Lute).SetSub(args[1].(bool))
	p.PopN(2)
}

func execmLuteSetGitConflict(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*lute.Lute).SetGitConflict(args[1].(bool))
	p.PopN(2)
}

func execmLuteSetJSRenderers(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*lute.Lute).SetJSRenderers(args[1].(map[string]map[string]*js.Object))
	p.PopN(2)
}

func execmLuteSpinVditorIRDOM(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*lute.Lute).SpinVditorIRDOM(args[1].(string))
	p.Ret(2, ret0)
}

func execmLuteHTML2VditorIRDOM(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*lute.Lute).HTML2VditorIRDOM(args[1].(string))
	p.Ret(2, ret0)
}

func execmLuteVditorIRDOM2HTML(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*lute.Lute).VditorIRDOM2HTML(args[1].(string))
	p.Ret(2, ret0)
}

func execmLuteMd2VditorIRDOM(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*lute.Lute).Md2VditorIRDOM(args[1].(string))
	p.Ret(2, ret0)
}

func execmLuteVditorIRDOM2Md(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*lute.Lute).VditorIRDOM2Md(args[1].(string))
	p.Ret(2, ret0)
}

func execmLuteSpinVditorIRBlockDOM(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*lute.Lute).SpinVditorIRBlockDOM(args[1].(string))
	p.Ret(2, ret0)
}

func execmLuteHTML2VditorIRBlockDOM(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*lute.Lute).HTML2VditorIRBlockDOM(args[1].(string))
	p.Ret(2, ret0)
}

func execmLuteVditorIRBlockDOM2HTML(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*lute.Lute).VditorIRBlockDOM2HTML(args[1].(string))
	p.Ret(2, ret0)
}

func execmLuteMd2VditorIRBlockDOM(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*lute.Lute).Md2VditorIRBlockDOM(args[1].(string))
	p.Ret(2, ret0)
}

func execmLuteVditorIRBlockDOM2Md(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*lute.Lute).VditorIRBlockDOM2Md(args[1].(string))
	p.Ret(2, ret0)
}

func execmLuteVditorIRBlockDOM2StdMd(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*lute.Lute).VditorIRBlockDOM2StdMd(args[1].(string))
	p.Ret(2, ret0)
}

func execmLuteVditorIRBlockDOM2Text(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*lute.Lute).VditorIRBlockDOM2Text(args[1].(string))
	p.Ret(2, ret0)
}

func execmLuteVditorIRBlockDOM2TextLen(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*lute.Lute).VditorIRBlockDOM2TextLen(args[1].(string))
	p.Ret(2, ret0)
}

func execmLuteTree2VditorIRBlockDOM(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*lute.Lute).Tree2VditorIRBlockDOM(args[1].(*parse.Tree), args[2].(*render.Options))
	p.Ret(3, ret0)
}

func execmLuteVditorIRBlockDOM2Tree(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(*lute.Lute).VditorIRBlockDOM2Tree(args[1].(string))
	p.Ret(2, ret0, ret1)
}

func execmLuteSpinVditorSVDOM(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*lute.Lute).SpinVditorSVDOM(args[1].(string))
	p.Ret(2, ret0)
}

func execmLuteHTML2VditorSVDOM(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*lute.Lute).HTML2VditorSVDOM(args[1].(string))
	p.Ret(2, ret0)
}

func execmLuteMd2VditorSVDOM(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*lute.Lute).Md2VditorSVDOM(args[1].(string))
	p.Ret(2, ret0)
}

func execmLuteMd2HTML(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*lute.Lute).Md2HTML(args[1].(string))
	p.Ret(2, ret0)
}

func execmLuteSpinVditorDOM(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*lute.Lute).SpinVditorDOM(args[1].(string))
	p.Ret(2, ret0)
}

func execmLuteHTML2VditorDOM(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*lute.Lute).HTML2VditorDOM(args[1].(string))
	p.Ret(2, ret0)
}

func execmLuteVditorDOM2HTML(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*lute.Lute).VditorDOM2HTML(args[1].(string))
	p.Ret(2, ret0)
}

func execmLuteMd2VditorDOM(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*lute.Lute).Md2VditorDOM(args[1].(string))
	p.Ret(2, ret0)
}

func execmLuteVditorDOM2Md(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*lute.Lute).VditorDOM2Md(args[1].(string))
	p.Ret(2, ret0)
}

func execmLuteRenderEChartsJSON(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*lute.Lute).RenderEChartsJSON(args[1].(string))
	p.Ret(2, ret0)
}

func execmLuteRenderKityMinderJSON(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*lute.Lute).RenderKityMinderJSON(args[1].(string))
	p.Ret(2, ret0)
}

func execmLuteHTML2Md(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*lute.Lute).HTML2Md(args[1].(string))
	p.Ret(2, ret0)
}

func toSlice0(args []interface{}) []lute.ParseOption {
	ret := make([]lute.ParseOption, len(args))
	for i, arg := range args {
		ret[i] = arg.(lute.ParseOption)
	}
	return ret
}

func execNew(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	ret0 := lute.New(toSlice0(args)...)
	p.Ret(arity, ret0)
}

func execRenderNodeVditorIRBlockDOM(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := lute.RenderNodeVditorIRBlockDOM(args[0].(*ast.Node), args[1].(*parse.Options), args[2].(*render.Options))
	p.Ret(3, ret0)
}

// I is a Go package instance.
var I = gop.NewGoPackage("github.com/88250/lute")

func init() {
	I.RegisterFuncs(
		I.Func("FormatNode", lute.FormatNode, execFormatNode),
		I.Func("(*Lute).HTML2Markdown", (*lute.Lute).HTML2Markdown, execmLuteHTML2Markdown),
		I.Func("(*Lute).HTML2Tree", (*lute.Lute).HTML2Tree, execmLuteHTML2Tree),
		I.Func("(*Lute).Markdown", (*lute.Lute).Markdown, execmLuteMarkdown),
		I.Func("(*Lute).MarkdownStr", (*lute.Lute).MarkdownStr, execmLuteMarkdownStr),
		I.Func("(*Lute).Format", (*lute.Lute).Format, execmLuteFormat),
		I.Func("(*Lute).FormatStr", (*lute.Lute).FormatStr, execmLuteFormatStr),
		I.Func("(*Lute).TextBundle", (*lute.Lute).TextBundle, execmLuteTextBundle),
		I.Func("(*Lute).TextBundleStr", (*lute.Lute).TextBundleStr, execmLuteTextBundleStr),
		I.Func("(*Lute).HTML2Text", (*lute.Lute).HTML2Text, execmLuteHTML2Text),
		I.Func("(*Lute).RenderJSON", (*lute.Lute).RenderJSON, execmLuteRenderJSON),
		I.Func("(*Lute).Space", (*lute.Lute).Space, execmLuteSpace),
		I.Func("(*Lute).IsValidLinkDest", (*lute.Lute).IsValidLinkDest, execmLuteIsValidLinkDest),
		I.Func("(*Lute).GetEmojis", (*lute.Lute).GetEmojis, execmLuteGetEmojis),
		I.Func("(*Lute).PutEmojis", (*lute.Lute).PutEmojis, execmLutePutEmojis),
		I.Func("(*Lute).RemoveEmoji", (*lute.Lute).RemoveEmoji, execmLuteRemoveEmoji),
		I.Func("(*Lute).GetTerms", (*lute.Lute).GetTerms, execmLuteGetTerms),
		I.Func("(*Lute).PutTerms", (*lute.Lute).PutTerms, execmLutePutTerms),
		I.Func("(*Lute).Tree2HTML", (*lute.Lute).Tree2HTML, execmLuteTree2HTML),
		I.Func("(*Lute).SetGFMTable", (*lute.Lute).SetGFMTable, execmLuteSetGFMTable),
		I.Func("(*Lute).SetGFMTaskListItem", (*lute.Lute).SetGFMTaskListItem, execmLuteSetGFMTaskListItem),
		I.Func("(*Lute).SetGFMTaskListItemClass", (*lute.Lute).SetGFMTaskListItemClass, execmLuteSetGFMTaskListItemClass),
		I.Func("(*Lute).SetGFMStrikethrough", (*lute.Lute).SetGFMStrikethrough, execmLuteSetGFMStrikethrough),
		I.Func("(*Lute).SetGFMAutoLink", (*lute.Lute).SetGFMAutoLink, execmLuteSetGFMAutoLink),
		I.Func("(*Lute).SetSoftBreak2HardBreak", (*lute.Lute).SetSoftBreak2HardBreak, execmLuteSetSoftBreak2HardBreak),
		I.Func("(*Lute).SetCodeSyntaxHighlight", (*lute.Lute).SetCodeSyntaxHighlight, execmLuteSetCodeSyntaxHighlight),
		I.Func("(*Lute).SetCodeSyntaxHighlightDetectLang", (*lute.Lute).SetCodeSyntaxHighlightDetectLang, execmLuteSetCodeSyntaxHighlightDetectLang),
		I.Func("(*Lute).SetCodeSyntaxHighlightInlineStyle", (*lute.Lute).SetCodeSyntaxHighlightInlineStyle, execmLuteSetCodeSyntaxHighlightInlineStyle),
		I.Func("(*Lute).SetCodeSyntaxHighlightLineNum", (*lute.Lute).SetCodeSyntaxHighlightLineNum, execmLuteSetCodeSyntaxHighlightLineNum),
		I.Func("(*Lute).SetCodeSyntaxHighlightStyleName", (*lute.Lute).SetCodeSyntaxHighlightStyleName, execmLuteSetCodeSyntaxHighlightStyleName),
		I.Func("(*Lute).SetFootnotes", (*lute.Lute).SetFootnotes, execmLuteSetFootnotes),
		I.Func("(*Lute).SetToC", (*lute.Lute).SetToC, execmLuteSetToC),
		I.Func("(*Lute).SetHeadingID", (*lute.Lute).SetHeadingID, execmLuteSetHeadingID),
		I.Func("(*Lute).SetAutoSpace", (*lute.Lute).SetAutoSpace, execmLuteSetAutoSpace),
		I.Func("(*Lute).SetFixTermTypo", (*lute.Lute).SetFixTermTypo, execmLuteSetFixTermTypo),
		I.Func("(*Lute).SetEmoji", (*lute.Lute).SetEmoji, execmLuteSetEmoji),
		I.Func("(*Lute).SetEmojis", (*lute.Lute).SetEmojis, execmLuteSetEmojis),
		I.Func("(*Lute).SetEmojiSite", (*lute.Lute).SetEmojiSite, execmLuteSetEmojiSite),
		I.Func("(*Lute).SetHeadingAnchor", (*lute.Lute).SetHeadingAnchor, execmLuteSetHeadingAnchor),
		I.Func("(*Lute).SetTerms", (*lute.Lute).SetTerms, execmLuteSetTerms),
		I.Func("(*Lute).SetVditorWYSIWYG", (*lute.Lute).SetVditorWYSIWYG, execmLuteSetVditorWYSIWYG),
		I.Func("(*Lute).SetVditorIR", (*lute.Lute).SetVditorIR, execmLuteSetVditorIR),
		I.Func("(*Lute).SetVditorSV", (*lute.Lute).SetVditorSV, execmLuteSetVditorSV),
		I.Func("(*Lute).SetInlineMathAllowDigitAfterOpenMarker", (*lute.Lute).SetInlineMathAllowDigitAfterOpenMarker, execmLuteSetInlineMathAllowDigitAfterOpenMarker),
		I.Func("(*Lute).SetLinkPrefix", (*lute.Lute).SetLinkPrefix, execmLuteSetLinkPrefix),
		I.Func("(*Lute).SetLinkBase", (*lute.Lute).SetLinkBase, execmLuteSetLinkBase),
		I.Func("(*Lute).GetLinkBase", (*lute.Lute).GetLinkBase, execmLuteGetLinkBase),
		I.Func("(*Lute).SetVditorCodeBlockPreview", (*lute.Lute).SetVditorCodeBlockPreview, execmLuteSetVditorCodeBlockPreview),
		I.Func("(*Lute).SetVditorMathBlockPreview", (*lute.Lute).SetVditorMathBlockPreview, execmLuteSetVditorMathBlockPreview),
		I.Func("(*Lute).SetVditorHTMLBlockPreview", (*lute.Lute).SetVditorHTMLBlockPreview, execmLuteSetVditorHTMLBlockPreview),
		I.Func("(*Lute).SetRenderListStyle", (*lute.Lute).SetRenderListStyle, execmLuteSetRenderListStyle),
		I.Func("(*Lute).SetSanitize", (*lute.Lute).SetSanitize, execmLuteSetSanitize),
		I.Func("(*Lute).SetImageLazyLoading", (*lute.Lute).SetImageLazyLoading, execmLuteSetImageLazyLoading),
		I.Func("(*Lute).SetChineseParagraphBeginningSpace", (*lute.Lute).SetChineseParagraphBeginningSpace, execmLuteSetChineseParagraphBeginningSpace),
		I.Func("(*Lute).SetYamlFrontMatter", (*lute.Lute).SetYamlFrontMatter, execmLuteSetYamlFrontMatter),
		I.Func("(*Lute).SetBlockRef", (*lute.Lute).SetBlockRef, execmLuteSetBlockRef),
		I.Func("(*Lute).SetMark", (*lute.Lute).SetMark, execmLuteSetMark),
		I.Func("(*Lute).SetKramdownIAL", (*lute.Lute).SetKramdownIAL, execmLuteSetKramdownIAL),
		I.Func("(*Lute).SetKramdownBlockIAL", (*lute.Lute).SetKramdownBlockIAL, execmLuteSetKramdownBlockIAL),
		I.Func("(*Lute).SetKramdownSpanIAL", (*lute.Lute).SetKramdownSpanIAL, execmLuteSetKramdownSpanIAL),
		I.Func("(*Lute).SetKramdownIALIDRenderName", (*lute.Lute).SetKramdownIALIDRenderName, execmLuteSetKramdownIALIDRenderName),
		I.Func("(*Lute).SetTag", (*lute.Lute).SetTag, execmLuteSetTag),
		I.Func("(*Lute).SetImgPathAllowSpace", (*lute.Lute).SetImgPathAllowSpace, execmLuteSetImgPathAllowSpace),
		I.Func("(*Lute).SetSuperBlock", (*lute.Lute).SetSuperBlock, execmLuteSetSuperBlock),
		I.Func("(*Lute).SetSup", (*lute.Lute).SetSup, execmLuteSetSup),
		I.Func("(*Lute).SetSub", (*lute.Lute).SetSub, execmLuteSetSub),
		I.Func("(*Lute).SetGitConflict", (*lute.Lute).SetGitConflict, execmLuteSetGitConflict),
		I.Func("(*Lute).SetJSRenderers", (*lute.Lute).SetJSRenderers, execmLuteSetJSRenderers),
		I.Func("(*Lute).SpinVditorIRDOM", (*lute.Lute).SpinVditorIRDOM, execmLuteSpinVditorIRDOM),
		I.Func("(*Lute).HTML2VditorIRDOM", (*lute.Lute).HTML2VditorIRDOM, execmLuteHTML2VditorIRDOM),
		I.Func("(*Lute).VditorIRDOM2HTML", (*lute.Lute).VditorIRDOM2HTML, execmLuteVditorIRDOM2HTML),
		I.Func("(*Lute).Md2VditorIRDOM", (*lute.Lute).Md2VditorIRDOM, execmLuteMd2VditorIRDOM),
		I.Func("(*Lute).VditorIRDOM2Md", (*lute.Lute).VditorIRDOM2Md, execmLuteVditorIRDOM2Md),
		I.Func("(*Lute).SpinVditorIRBlockDOM", (*lute.Lute).SpinVditorIRBlockDOM, execmLuteSpinVditorIRBlockDOM),
		I.Func("(*Lute).HTML2VditorIRBlockDOM", (*lute.Lute).HTML2VditorIRBlockDOM, execmLuteHTML2VditorIRBlockDOM),
		I.Func("(*Lute).VditorIRBlockDOM2HTML", (*lute.Lute).VditorIRBlockDOM2HTML, execmLuteVditorIRBlockDOM2HTML),
		I.Func("(*Lute).Md2VditorIRBlockDOM", (*lute.Lute).Md2VditorIRBlockDOM, execmLuteMd2VditorIRBlockDOM),
		I.Func("(*Lute).VditorIRBlockDOM2Md", (*lute.Lute).VditorIRBlockDOM2Md, execmLuteVditorIRBlockDOM2Md),
		I.Func("(*Lute).VditorIRBlockDOM2StdMd", (*lute.Lute).VditorIRBlockDOM2StdMd, execmLuteVditorIRBlockDOM2StdMd),
		I.Func("(*Lute).VditorIRBlockDOM2Text", (*lute.Lute).VditorIRBlockDOM2Text, execmLuteVditorIRBlockDOM2Text),
		I.Func("(*Lute).VditorIRBlockDOM2TextLen", (*lute.Lute).VditorIRBlockDOM2TextLen, execmLuteVditorIRBlockDOM2TextLen),
		I.Func("(*Lute).Tree2VditorIRBlockDOM", (*lute.Lute).Tree2VditorIRBlockDOM, execmLuteTree2VditorIRBlockDOM),
		I.Func("(*Lute).VditorIRBlockDOM2Tree", (*lute.Lute).VditorIRBlockDOM2Tree, execmLuteVditorIRBlockDOM2Tree),
		I.Func("(*Lute).SpinVditorSVDOM", (*lute.Lute).SpinVditorSVDOM, execmLuteSpinVditorSVDOM),
		I.Func("(*Lute).HTML2VditorSVDOM", (*lute.Lute).HTML2VditorSVDOM, execmLuteHTML2VditorSVDOM),
		I.Func("(*Lute).Md2VditorSVDOM", (*lute.Lute).Md2VditorSVDOM, execmLuteMd2VditorSVDOM),
		I.Func("(*Lute).Md2HTML", (*lute.Lute).Md2HTML, execmLuteMd2HTML),
		I.Func("(*Lute).SpinVditorDOM", (*lute.Lute).SpinVditorDOM, execmLuteSpinVditorDOM),
		I.Func("(*Lute).HTML2VditorDOM", (*lute.Lute).HTML2VditorDOM, execmLuteHTML2VditorDOM),
		I.Func("(*Lute).VditorDOM2HTML", (*lute.Lute).VditorDOM2HTML, execmLuteVditorDOM2HTML),
		I.Func("(*Lute).Md2VditorDOM", (*lute.Lute).Md2VditorDOM, execmLuteMd2VditorDOM),
		I.Func("(*Lute).VditorDOM2Md", (*lute.Lute).VditorDOM2Md, execmLuteVditorDOM2Md),
		I.Func("(*Lute).RenderEChartsJSON", (*lute.Lute).RenderEChartsJSON, execmLuteRenderEChartsJSON),
		I.Func("(*Lute).RenderKityMinderJSON", (*lute.Lute).RenderKityMinderJSON, execmLuteRenderKityMinderJSON),
		I.Func("(*Lute).HTML2Md", (*lute.Lute).HTML2Md, execmLuteHTML2Md),
		I.Func("RenderNodeVditorIRBlockDOM", lute.RenderNodeVditorIRBlockDOM, execRenderNodeVditorIRBlockDOM),
	)
	I.RegisterFuncvs(
		I.Funcv("New", lute.New, execNew),
	)
	I.RegisterTypes(
		I.Type("Lute", reflect.TypeOf((*lute.Lute)(nil)).Elem()),
		I.Type("ParseOption", reflect.TypeOf((*lute.ParseOption)(nil)).Elem()),
	)
	I.RegisterConsts(
		I.Const("Version", qspec.ConstBoundString, lute.Version),
	)
}
