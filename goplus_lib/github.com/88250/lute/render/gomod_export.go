// Package render provide Go+ "github.com/88250/lute/render" package, as "github.com/88250/lute/render" package in Go.
package render

import (
	reflect "reflect"

	ast "github.com/88250/lute/ast"
	parse "github.com/88250/lute/parse"
	render "github.com/88250/lute/render"
	gop "github.com/goplus/gop"
)

func execmBaseRendererLinkPath(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*render.BaseRenderer).LinkPath(args[1].([]byte))
	p.Ret(2, ret0)
}

func execmBaseRendererPrefixPath(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*render.BaseRenderer).PrefixPath(args[1].([]byte))
	p.Ret(2, ret0)
}

func execmBaseRendererRelativePath(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*render.BaseRenderer).RelativePath(args[1].([]byte))
	p.Ret(2, ret0)
}

func execmBaseRendererRender(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*render.BaseRenderer).Render()
	p.Ret(1, ret0)
}

func execmBaseRendererWriteByte(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*render.BaseRenderer).WriteByte(args[1].(byte))
	p.PopN(2)
}

func execmBaseRendererWrite(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*render.BaseRenderer).Write(args[1].([]byte))
	p.PopN(2)
}

func execmBaseRendererWriteString(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*render.BaseRenderer).WriteString(args[1].(string))
	p.PopN(2)
}

func execmBaseRendererNewline(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	args[0].(*render.BaseRenderer).Newline()
	p.PopN(1)
}

func execmBaseRendererTextAutoSpacePrevious(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*render.BaseRenderer).TextAutoSpacePrevious(args[1].(*ast.Node))
	p.PopN(2)
}

func execmBaseRendererTextAutoSpaceNext(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*render.BaseRenderer).TextAutoSpaceNext(args[1].(*ast.Node))
	p.PopN(2)
}

func execmBaseRendererLinkTextAutoSpacePrevious(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*render.BaseRenderer).LinkTextAutoSpacePrevious(args[1].(*ast.Node))
	p.PopN(2)
}

func execmBaseRendererLinkTextAutoSpaceNext(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*render.BaseRenderer).LinkTextAutoSpaceNext(args[1].(*ast.Node))
	p.PopN(2)
}

func execmBaseRendererTag(_ int, p *gop.Context) {
	args := p.GetArgs(4)
	args[0].(*render.BaseRenderer).Tag(args[1].(string), args[2].([][]string), args[3].(bool))
	p.PopN(4)
}

func execmBaseRendererNodeID(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*render.BaseRenderer).NodeID(args[1].(*ast.Node))
	p.Ret(2, ret0)
}

func execmBaseRendererNodeAttrs(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*render.BaseRenderer).NodeAttrs(args[1].(*ast.Node))
	p.Ret(2, ret0)
}

func execmBaseRendererNodeAttrsStr(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*render.BaseRenderer).NodeAttrsStr(args[1].(*ast.Node))
	p.Ret(2, ret0)
}

func execmBaseRendererSpace(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*render.BaseRenderer).Space(args[1].([]byte))
	p.Ret(2, ret0)
}

func execmBaseRendererFixTermTypo(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*render.BaseRenderer).FixTermTypo(args[1].([]byte))
	p.Ret(2, ret0)
}

func execHeadingID(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := render.HeadingID(args[0].(*ast.Node))
	p.Ret(1, ret0)
}

func execmHtmlRendererRender(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*render.HtmlRenderer).Render()
	p.Ret(1, ret0)
}

func execmHtmlRendererRenderFootnotes(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*render.HtmlRenderer).RenderFootnotes()
	p.Ret(1, ret0)
}

func execNewBaseRenderer(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := render.NewBaseRenderer(args[0].(*parse.Tree), args[1].(*render.Options))
	p.Ret(2, ret0)
}

func execNewEChartsJSONRenderer(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := render.NewEChartsJSONRenderer(args[0].(*parse.Tree), args[1].(*render.Options))
	p.Ret(2, ret0)
}

func execNewFormatRenderer(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := render.NewFormatRenderer(args[0].(*parse.Tree), args[1].(*render.Options))
	p.Ret(2, ret0)
}

func execNewHtmlRenderer(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := render.NewHtmlRenderer(args[0].(*parse.Tree), args[1].(*render.Options))
	p.Ret(2, ret0)
}

func execNewJSONRenderer(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := render.NewJSONRenderer(args[0].(*parse.Tree), args[1].(*render.Options))
	p.Ret(2, ret0)
}

func execNewKityMinderJSONRenderer(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := render.NewKityMinderJSONRenderer(args[0].(*parse.Tree), args[1].(*render.Options))
	p.Ret(2, ret0)
}

func execNewOptions(_ int, p *gop.Context) {
	ret0 := render.NewOptions()
	p.Ret(0, ret0)
}

func execNewTerms(_ int, p *gop.Context) {
	ret0 := render.NewTerms()
	p.Ret(0, ret0)
}

func execNewTextBundleRenderer(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := render.NewTextBundleRenderer(args[0].(*parse.Tree), args[1].([]string), args[2].(*render.Options))
	p.Ret(3, ret0)
}

func execNewVditorIRBlockRenderer(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := render.NewVditorIRBlockRenderer(args[0].(*parse.Tree), args[1].(*render.Options))
	p.Ret(2, ret0)
}

func execNewVditorIRRenderer(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := render.NewVditorIRRenderer(args[0].(*parse.Tree), args[1].(*render.Options))
	p.Ret(2, ret0)
}

func execNewVditorRenderer(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := render.NewVditorRenderer(args[0].(*parse.Tree), args[1].(*render.Options))
	p.Ret(2, ret0)
}

func execNewVditorSVRenderer(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := render.NewVditorSVRenderer(args[0].(*parse.Tree), args[1].(*render.Options))
	p.Ret(2, ret0)
}

func execRenderHeadingText(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := render.RenderHeadingText(args[0].(*ast.Node))
	p.Ret(1, ret0)
}

func execiRendererRender(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(render.Renderer).Render()
	p.Ret(1, ret0)
}

func execSpace0(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := render.Space0(args[0].(string))
	p.Ret(1, ret0)
}

func execSubStr(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := render.SubStr(args[0].(string), args[1].(int))
	p.Ret(2, ret0)
}

func execmTextBundleRendererRender(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := args[0].(*render.TextBundleRenderer).Render()
	p.Ret(1, ret0, ret1)
}

func execmVditorIRBlockRendererNodeID(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*render.VditorIRBlockRenderer).NodeID(args[1].(*ast.Node))
	p.Ret(2, ret0)
}

func execmVditorIRBlockRendererText(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*render.VditorIRBlockRenderer).Text(args[1].(*ast.Node))
	p.Ret(2, ret0)
}

func execmVditorIRRendererText(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*render.VditorIRRenderer).Text(args[1].(*ast.Node))
	p.Ret(2, ret0)
}

func execmVditorSVRendererWriteByte(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*render.VditorSVRenderer).WriteByte(args[1].(byte))
	p.PopN(2)
}

func execmVditorSVRendererWrite(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*render.VditorSVRenderer).Write(args[1].([]byte))
	p.PopN(2)
}

func execmVditorSVRendererWriteString(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*render.VditorSVRenderer).WriteString(args[1].(string))
	p.PopN(2)
}

func execmVditorSVRendererNewline(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	args[0].(*render.VditorSVRenderer).Newline()
	p.PopN(1)
}

func execmVditorSVRendererText(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*render.VditorSVRenderer).Text(args[1].(*ast.Node))
	p.Ret(2, ret0)
}

// I is a Go package instance.
var I = gop.NewGoPackage("github.com/88250/lute/render")

func init() {
	I.RegisterFuncs(
		I.Func("(*BaseRenderer).LinkPath", (*render.BaseRenderer).LinkPath, execmBaseRendererLinkPath),
		I.Func("(*BaseRenderer).PrefixPath", (*render.BaseRenderer).PrefixPath, execmBaseRendererPrefixPath),
		I.Func("(*BaseRenderer).RelativePath", (*render.BaseRenderer).RelativePath, execmBaseRendererRelativePath),
		I.Func("(*BaseRenderer).Render", (*render.BaseRenderer).Render, execmBaseRendererRender),
		I.Func("(*BaseRenderer).WriteByte", (*render.BaseRenderer).WriteByte, execmBaseRendererWriteByte),
		I.Func("(*BaseRenderer).Write", (*render.BaseRenderer).Write, execmBaseRendererWrite),
		I.Func("(*BaseRenderer).WriteString", (*render.BaseRenderer).WriteString, execmBaseRendererWriteString),
		I.Func("(*BaseRenderer).Newline", (*render.BaseRenderer).Newline, execmBaseRendererNewline),
		I.Func("(*BaseRenderer).TextAutoSpacePrevious", (*render.BaseRenderer).TextAutoSpacePrevious, execmBaseRendererTextAutoSpacePrevious),
		I.Func("(*BaseRenderer).TextAutoSpaceNext", (*render.BaseRenderer).TextAutoSpaceNext, execmBaseRendererTextAutoSpaceNext),
		I.Func("(*BaseRenderer).LinkTextAutoSpacePrevious", (*render.BaseRenderer).LinkTextAutoSpacePrevious, execmBaseRendererLinkTextAutoSpacePrevious),
		I.Func("(*BaseRenderer).LinkTextAutoSpaceNext", (*render.BaseRenderer).LinkTextAutoSpaceNext, execmBaseRendererLinkTextAutoSpaceNext),
		I.Func("(*BaseRenderer).Tag", (*render.BaseRenderer).Tag, execmBaseRendererTag),
		I.Func("(*BaseRenderer).NodeID", (*render.BaseRenderer).NodeID, execmBaseRendererNodeID),
		I.Func("(*BaseRenderer).NodeAttrs", (*render.BaseRenderer).NodeAttrs, execmBaseRendererNodeAttrs),
		I.Func("(*BaseRenderer).NodeAttrsStr", (*render.BaseRenderer).NodeAttrsStr, execmBaseRendererNodeAttrsStr),
		I.Func("(*BaseRenderer).Space", (*render.BaseRenderer).Space, execmBaseRendererSpace),
		I.Func("(*BaseRenderer).FixTermTypo", (*render.BaseRenderer).FixTermTypo, execmBaseRendererFixTermTypo),
		I.Func("HeadingID", render.HeadingID, execHeadingID),
		I.Func("(*HtmlRenderer).Render", (*render.HtmlRenderer).Render, execmHtmlRendererRender),
		I.Func("(*HtmlRenderer).RenderFootnotes", (*render.HtmlRenderer).RenderFootnotes, execmHtmlRendererRenderFootnotes),
		I.Func("NewBaseRenderer", render.NewBaseRenderer, execNewBaseRenderer),
		I.Func("NewEChartsJSONRenderer", render.NewEChartsJSONRenderer, execNewEChartsJSONRenderer),
		I.Func("NewFormatRenderer", render.NewFormatRenderer, execNewFormatRenderer),
		I.Func("NewHtmlRenderer", render.NewHtmlRenderer, execNewHtmlRenderer),
		I.Func("NewJSONRenderer", render.NewJSONRenderer, execNewJSONRenderer),
		I.Func("NewKityMinderJSONRenderer", render.NewKityMinderJSONRenderer, execNewKityMinderJSONRenderer),
		I.Func("NewOptions", render.NewOptions, execNewOptions),
		I.Func("NewTerms", render.NewTerms, execNewTerms),
		I.Func("NewTextBundleRenderer", render.NewTextBundleRenderer, execNewTextBundleRenderer),
		I.Func("NewVditorIRBlockRenderer", render.NewVditorIRBlockRenderer, execNewVditorIRBlockRenderer),
		I.Func("NewVditorIRRenderer", render.NewVditorIRRenderer, execNewVditorIRRenderer),
		I.Func("NewVditorRenderer", render.NewVditorRenderer, execNewVditorRenderer),
		I.Func("NewVditorSVRenderer", render.NewVditorSVRenderer, execNewVditorSVRenderer),
		I.Func("RenderHeadingText", render.RenderHeadingText, execRenderHeadingText),
		I.Func("(Renderer).Render", (render.Renderer).Render, execiRendererRender),
		I.Func("Space0", render.Space0, execSpace0),
		I.Func("SubStr", render.SubStr, execSubStr),
		I.Func("(*TextBundleRenderer).Render", (*render.TextBundleRenderer).Render, execmTextBundleRendererRender),
		I.Func("(*VditorIRBlockRenderer).NodeID", (*render.VditorIRBlockRenderer).NodeID, execmVditorIRBlockRendererNodeID),
		I.Func("(*VditorIRBlockRenderer).Text", (*render.VditorIRBlockRenderer).Text, execmVditorIRBlockRendererText),
		I.Func("(*VditorIRRenderer).Text", (*render.VditorIRRenderer).Text, execmVditorIRRendererText),
		I.Func("(*VditorSVRenderer).WriteByte", (*render.VditorSVRenderer).WriteByte, execmVditorSVRendererWriteByte),
		I.Func("(*VditorSVRenderer).Write", (*render.VditorSVRenderer).Write, execmVditorSVRendererWrite),
		I.Func("(*VditorSVRenderer).WriteString", (*render.VditorSVRenderer).WriteString, execmVditorSVRendererWriteString),
		I.Func("(*VditorSVRenderer).Newline", (*render.VditorSVRenderer).Newline, execmVditorSVRendererNewline),
		I.Func("(*VditorSVRenderer).Text", (*render.VditorSVRenderer).Text, execmVditorSVRendererText),
	)
	I.RegisterVars(
		I.Var("NewlineSV", &render.NewlineSV),
	)
	I.RegisterTypes(
		I.Type("BaseRenderer", reflect.TypeOf((*render.BaseRenderer)(nil)).Elem()),
		I.Type("EChartsJSONRenderer", reflect.TypeOf((*render.EChartsJSONRenderer)(nil)).Elem()),
		I.Type("ExtRendererFunc", reflect.TypeOf((*render.ExtRendererFunc)(nil)).Elem()),
		I.Type("FormatRenderer", reflect.TypeOf((*render.FormatRenderer)(nil)).Elem()),
		I.Type("Heading", reflect.TypeOf((*render.Heading)(nil)).Elem()),
		I.Type("HtmlRenderer", reflect.TypeOf((*render.HtmlRenderer)(nil)).Elem()),
		I.Type("JSONRenderer", reflect.TypeOf((*render.JSONRenderer)(nil)).Elem()),
		I.Type("KityMinderJSONRenderer", reflect.TypeOf((*render.KityMinderJSONRenderer)(nil)).Elem()),
		I.Type("Options", reflect.TypeOf((*render.Options)(nil)).Elem()),
		I.Type("Renderer", reflect.TypeOf((*render.Renderer)(nil)).Elem()),
		I.Type("RendererFunc", reflect.TypeOf((*render.RendererFunc)(nil)).Elem()),
		I.Type("TextBundleRenderer", reflect.TypeOf((*render.TextBundleRenderer)(nil)).Elem()),
		I.Type("VditorIRBlockRenderer", reflect.TypeOf((*render.VditorIRBlockRenderer)(nil)).Elem()),
		I.Type("VditorIRRenderer", reflect.TypeOf((*render.VditorIRRenderer)(nil)).Elem()),
		I.Type("VditorRenderer", reflect.TypeOf((*render.VditorRenderer)(nil)).Elem()),
		I.Type("VditorSVRenderer", reflect.TypeOf((*render.VditorSVRenderer)(nil)).Elem()),
	)
}
