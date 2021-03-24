// Package ast provide Go+ "github.com/88250/lute/ast" package, as "github.com/88250/lute/ast" package in Go.
package ast

import (
	reflect "reflect"

	ast "github.com/88250/lute/ast"
	gop "github.com/goplus/gop"
	qspec "github.com/goplus/gop/exec.spec"
)

func execNewNodeID(_ int, p *gop.Context) {
	ret0 := ast.NewNodeID()
	p.Ret(0, ret0)
}

func execmNodeRemoveIALAttr(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*ast.Node).RemoveIALAttr(args[1].(string))
	p.PopN(2)
}

func execmNodeSetIALAttr(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	args[0].(*ast.Node).SetIALAttr(args[1].(string), args[2].(string))
	p.PopN(3)
}

func execmNodeIALAttr(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*ast.Node).IALAttr(args[1].(string))
	p.Ret(2, ret0)
}

func execmNodeTokensStr(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*ast.Node).TokensStr()
	p.Ret(1, ret0)
}

func execmNodeLastDeepestChild(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*ast.Node).LastDeepestChild()
	p.Ret(1, ret0)
}

func execmNodeFirstDeepestChild(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*ast.Node).FirstDeepestChild()
	p.Ret(1, ret0)
}

func execmNodeChildByType(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*ast.Node).ChildByType(args[1].(ast.NodeType))
	p.Ret(2, ret0)
}

func execmNodeChildrenByType(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*ast.Node).ChildrenByType(args[1].(ast.NodeType))
	p.Ret(2, ret0)
}

func execmNodeText(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*ast.Node).Text()
	p.Ret(1, ret0)
}

func execmNodeTextLen(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*ast.Node).TextLen()
	p.Ret(1, ret0)
}

func execmNodeTokenLen(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*ast.Node).TokenLen()
	p.Ret(1, ret0)
}

func execmNodeNextNodeText(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*ast.Node).NextNodeText()
	p.Ret(1, ret0)
}

func execmNodePreviousNodeText(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*ast.Node).PreviousNodeText()
	p.Ret(1, ret0)
}

func execmNodeUnlink(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	args[0].(*ast.Node).Unlink()
	p.PopN(1)
}

func execmNodeAppendTokens(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*ast.Node).AppendTokens(args[1].([]byte))
	p.PopN(2)
}

func execmNodeInsertAfter(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*ast.Node).InsertAfter(args[1].(*ast.Node))
	p.PopN(2)
}

func execmNodeInsertBefore(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*ast.Node).InsertBefore(args[1].(*ast.Node))
	p.PopN(2)
}

func execmNodeAppendChild(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*ast.Node).AppendChild(args[1].(*ast.Node))
	p.PopN(2)
}

func execmNodePrependChild(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*ast.Node).PrependChild(args[1].(*ast.Node))
	p.PopN(2)
}

func execmNodeList(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*ast.Node).List()
	p.Ret(1, ret0)
}

func toSlice0(args []interface{}) []ast.NodeType {
	ret := make([]ast.NodeType, len(args))
	for i, arg := range args {
		ret[i] = arg.(ast.NodeType)
	}
	return ret
}

func execmNodeParentIs(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	ret0 := args[0].(*ast.Node).ParentIs(args[1].(ast.NodeType), toSlice0(args[2:])...)
	p.Ret(arity, ret0)
}

func execmNodeIsBlock(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*ast.Node).IsBlock()
	p.Ret(1, ret0)
}

func execmNodeIsContainerBlock(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*ast.Node).IsContainerBlock()
	p.Ret(1, ret0)
}

func execmNodeIsMarker(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*ast.Node).IsMarker()
	p.Ret(1, ret0)
}

func execmNodeAcceptLines(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*ast.Node).AcceptLines()
	p.Ret(1, ret0)
}

func execmNodeCanContain(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*ast.Node).CanContain(args[1].(ast.NodeType))
	p.Ret(2, ret0)
}

func execmNodeTypeString(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(ast.NodeType).String()
	p.Ret(1, ret0)
}

func execStr2NodeType(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := ast.Str2NodeType(args[0].(string))
	p.Ret(1, ret0)
}

func execWalk(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ast.Walk(args[0].(*ast.Node), args[1].(ast.Walker))
	p.PopN(2)
}

// I is a Go package instance.
var I = gop.NewGoPackage("github.com/88250/lute/ast")

func init() {
	I.RegisterFuncs(
		I.Func("NewNodeID", ast.NewNodeID, execNewNodeID),
		I.Func("(*Node).RemoveIALAttr", (*ast.Node).RemoveIALAttr, execmNodeRemoveIALAttr),
		I.Func("(*Node).SetIALAttr", (*ast.Node).SetIALAttr, execmNodeSetIALAttr),
		I.Func("(*Node).IALAttr", (*ast.Node).IALAttr, execmNodeIALAttr),
		I.Func("(*Node).TokensStr", (*ast.Node).TokensStr, execmNodeTokensStr),
		I.Func("(*Node).LastDeepestChild", (*ast.Node).LastDeepestChild, execmNodeLastDeepestChild),
		I.Func("(*Node).FirstDeepestChild", (*ast.Node).FirstDeepestChild, execmNodeFirstDeepestChild),
		I.Func("(*Node).ChildByType", (*ast.Node).ChildByType, execmNodeChildByType),
		I.Func("(*Node).ChildrenByType", (*ast.Node).ChildrenByType, execmNodeChildrenByType),
		I.Func("(*Node).Text", (*ast.Node).Text, execmNodeText),
		I.Func("(*Node).TextLen", (*ast.Node).TextLen, execmNodeTextLen),
		I.Func("(*Node).TokenLen", (*ast.Node).TokenLen, execmNodeTokenLen),
		I.Func("(*Node).NextNodeText", (*ast.Node).NextNodeText, execmNodeNextNodeText),
		I.Func("(*Node).PreviousNodeText", (*ast.Node).PreviousNodeText, execmNodePreviousNodeText),
		I.Func("(*Node).Unlink", (*ast.Node).Unlink, execmNodeUnlink),
		I.Func("(*Node).AppendTokens", (*ast.Node).AppendTokens, execmNodeAppendTokens),
		I.Func("(*Node).InsertAfter", (*ast.Node).InsertAfter, execmNodeInsertAfter),
		I.Func("(*Node).InsertBefore", (*ast.Node).InsertBefore, execmNodeInsertBefore),
		I.Func("(*Node).AppendChild", (*ast.Node).AppendChild, execmNodeAppendChild),
		I.Func("(*Node).PrependChild", (*ast.Node).PrependChild, execmNodePrependChild),
		I.Func("(*Node).List", (*ast.Node).List, execmNodeList),
		I.Func("(*Node).IsBlock", (*ast.Node).IsBlock, execmNodeIsBlock),
		I.Func("(*Node).IsContainerBlock", (*ast.Node).IsContainerBlock, execmNodeIsContainerBlock),
		I.Func("(*Node).IsMarker", (*ast.Node).IsMarker, execmNodeIsMarker),
		I.Func("(*Node).AcceptLines", (*ast.Node).AcceptLines, execmNodeAcceptLines),
		I.Func("(*Node).CanContain", (*ast.Node).CanContain, execmNodeCanContain),
		I.Func("(NodeType).String", (ast.NodeType).String, execmNodeTypeString),
		I.Func("Str2NodeType", ast.Str2NodeType, execStr2NodeType),
		I.Func("Walk", ast.Walk, execWalk),
	)
	I.RegisterFuncvs(
		I.Funcv("(*Node).ParentIs", (*ast.Node).ParentIs, execmNodeParentIs),
	)
	I.RegisterVars(
		I.Var("Testing", &ast.Testing),
	)
	I.RegisterTypes(
		I.Type("ListData", reflect.TypeOf((*ast.ListData)(nil)).Elem()),
		I.Type("Node", reflect.TypeOf((*ast.Node)(nil)).Elem()),
		I.Type("NodeType", reflect.TypeOf((*ast.NodeType)(nil)).Elem()),
		I.Type("WalkStatus", reflect.TypeOf((*ast.WalkStatus)(nil)).Elem()),
		I.Type("Walker", reflect.TypeOf((*ast.Walker)(nil)).Elem()),
	)
	I.RegisterConsts(
		I.Const("NodeBackslash", qspec.Int, ast.NodeBackslash),
		I.Const("NodeBackslashContent", qspec.Int, ast.NodeBackslashContent),
		I.Const("NodeBang", qspec.Int, ast.NodeBang),
		I.Const("NodeBlockEmbed", qspec.Int, ast.NodeBlockEmbed),
		I.Const("NodeBlockEmbedID", qspec.Int, ast.NodeBlockEmbedID),
		I.Const("NodeBlockEmbedSpace", qspec.Int, ast.NodeBlockEmbedSpace),
		I.Const("NodeBlockEmbedText", qspec.Int, ast.NodeBlockEmbedText),
		I.Const("NodeBlockEmbedTextTplRenderResult", qspec.Int, ast.NodeBlockEmbedTextTplRenderResult),
		I.Const("NodeBlockQueryEmbed", qspec.Int, ast.NodeBlockQueryEmbed),
		I.Const("NodeBlockQueryEmbedScript", qspec.Int, ast.NodeBlockQueryEmbedScript),
		I.Const("NodeBlockRef", qspec.Int, ast.NodeBlockRef),
		I.Const("NodeBlockRefID", qspec.Int, ast.NodeBlockRefID),
		I.Const("NodeBlockRefSpace", qspec.Int, ast.NodeBlockRefSpace),
		I.Const("NodeBlockRefText", qspec.Int, ast.NodeBlockRefText),
		I.Const("NodeBlockRefTextTplRenderResult", qspec.Int, ast.NodeBlockRefTextTplRenderResult),
		I.Const("NodeBlockquote", qspec.Int, ast.NodeBlockquote),
		I.Const("NodeBlockquoteMarker", qspec.Int, ast.NodeBlockquoteMarker),
		I.Const("NodeCloseBrace", qspec.Int, ast.NodeCloseBrace),
		I.Const("NodeCloseBracket", qspec.Int, ast.NodeCloseBracket),
		I.Const("NodeCloseParen", qspec.Int, ast.NodeCloseParen),
		I.Const("NodeCodeBlock", qspec.Int, ast.NodeCodeBlock),
		I.Const("NodeCodeBlockCode", qspec.Int, ast.NodeCodeBlockCode),
		I.Const("NodeCodeBlockFenceCloseMarker", qspec.Int, ast.NodeCodeBlockFenceCloseMarker),
		I.Const("NodeCodeBlockFenceInfoMarker", qspec.Int, ast.NodeCodeBlockFenceInfoMarker),
		I.Const("NodeCodeBlockFenceOpenMarker", qspec.Int, ast.NodeCodeBlockFenceOpenMarker),
		I.Const("NodeCodeSpan", qspec.Int, ast.NodeCodeSpan),
		I.Const("NodeCodeSpanCloseMarker", qspec.Int, ast.NodeCodeSpanCloseMarker),
		I.Const("NodeCodeSpanContent", qspec.Int, ast.NodeCodeSpanContent),
		I.Const("NodeCodeSpanOpenMarker", qspec.Int, ast.NodeCodeSpanOpenMarker),
		I.Const("NodeDocument", qspec.Int, ast.NodeDocument),
		I.Const("NodeEmA6kCloseMarker", qspec.Int, ast.NodeEmA6kCloseMarker),
		I.Const("NodeEmA6kOpenMarker", qspec.Int, ast.NodeEmA6kOpenMarker),
		I.Const("NodeEmU8eCloseMarker", qspec.Int, ast.NodeEmU8eCloseMarker),
		I.Const("NodeEmU8eOpenMarker", qspec.Int, ast.NodeEmU8eOpenMarker),
		I.Const("NodeEmoji", qspec.Int, ast.NodeEmoji),
		I.Const("NodeEmojiAlias", qspec.Int, ast.NodeEmojiAlias),
		I.Const("NodeEmojiImg", qspec.Int, ast.NodeEmojiImg),
		I.Const("NodeEmojiUnicode", qspec.Int, ast.NodeEmojiUnicode),
		I.Const("NodeEmphasis", qspec.Int, ast.NodeEmphasis),
		I.Const("NodeFootnotesDef", qspec.Int, ast.NodeFootnotesDef),
		I.Const("NodeFootnotesDefBlock", qspec.Int, ast.NodeFootnotesDefBlock),
		I.Const("NodeFootnotesRef", qspec.Int, ast.NodeFootnotesRef),
		I.Const("NodeGitConflict", qspec.Int, ast.NodeGitConflict),
		I.Const("NodeGitConflictCloseMarker", qspec.Int, ast.NodeGitConflictCloseMarker),
		I.Const("NodeGitConflictContent", qspec.Int, ast.NodeGitConflictContent),
		I.Const("NodeGitConflictOpenMarker", qspec.Int, ast.NodeGitConflictOpenMarker),
		I.Const("NodeHTMLBlock", qspec.Int, ast.NodeHTMLBlock),
		I.Const("NodeHTMLEntity", qspec.Int, ast.NodeHTMLEntity),
		I.Const("NodeHardBreak", qspec.Int, ast.NodeHardBreak),
		I.Const("NodeHeading", qspec.Int, ast.NodeHeading),
		I.Const("NodeHeadingC8hMarker", qspec.Int, ast.NodeHeadingC8hMarker),
		I.Const("NodeHeadingID", qspec.Int, ast.NodeHeadingID),
		I.Const("NodeImage", qspec.Int, ast.NodeImage),
		I.Const("NodeInlineHTML", qspec.Int, ast.NodeInlineHTML),
		I.Const("NodeInlineMath", qspec.Int, ast.NodeInlineMath),
		I.Const("NodeInlineMathCloseMarker", qspec.Int, ast.NodeInlineMathCloseMarker),
		I.Const("NodeInlineMathContent", qspec.Int, ast.NodeInlineMathContent),
		I.Const("NodeInlineMathOpenMarker", qspec.Int, ast.NodeInlineMathOpenMarker),
		I.Const("NodeKramdownBlockIAL", qspec.Int, ast.NodeKramdownBlockIAL),
		I.Const("NodeKramdownSpanIAL", qspec.Int, ast.NodeKramdownSpanIAL),
		I.Const("NodeLink", qspec.Int, ast.NodeLink),
		I.Const("NodeLinkDest", qspec.Int, ast.NodeLinkDest),
		I.Const("NodeLinkRefDef", qspec.Int, ast.NodeLinkRefDef),
		I.Const("NodeLinkRefDefBlock", qspec.Int, ast.NodeLinkRefDefBlock),
		I.Const("NodeLinkSpace", qspec.Int, ast.NodeLinkSpace),
		I.Const("NodeLinkText", qspec.Int, ast.NodeLinkText),
		I.Const("NodeLinkTitle", qspec.Int, ast.NodeLinkTitle),
		I.Const("NodeList", qspec.Int, ast.NodeList),
		I.Const("NodeListItem", qspec.Int, ast.NodeListItem),
		I.Const("NodeMark", qspec.Int, ast.NodeMark),
		I.Const("NodeMark1CloseMarker", qspec.Int, ast.NodeMark1CloseMarker),
		I.Const("NodeMark1OpenMarker", qspec.Int, ast.NodeMark1OpenMarker),
		I.Const("NodeMark2CloseMarker", qspec.Int, ast.NodeMark2CloseMarker),
		I.Const("NodeMark2OpenMarker", qspec.Int, ast.NodeMark2OpenMarker),
		I.Const("NodeMathBlock", qspec.Int, ast.NodeMathBlock),
		I.Const("NodeMathBlockCloseMarker", qspec.Int, ast.NodeMathBlockCloseMarker),
		I.Const("NodeMathBlockContent", qspec.Int, ast.NodeMathBlockContent),
		I.Const("NodeMathBlockOpenMarker", qspec.Int, ast.NodeMathBlockOpenMarker),
		I.Const("NodeOpenBrace", qspec.Int, ast.NodeOpenBrace),
		I.Const("NodeOpenBracket", qspec.Int, ast.NodeOpenBracket),
		I.Const("NodeOpenParen", qspec.Int, ast.NodeOpenParen),
		I.Const("NodeParagraph", qspec.Int, ast.NodeParagraph),
		I.Const("NodeSoftBreak", qspec.Int, ast.NodeSoftBreak),
		I.Const("NodeStrikethrough", qspec.Int, ast.NodeStrikethrough),
		I.Const("NodeStrikethrough1CloseMarker", qspec.Int, ast.NodeStrikethrough1CloseMarker),
		I.Const("NodeStrikethrough1OpenMarker", qspec.Int, ast.NodeStrikethrough1OpenMarker),
		I.Const("NodeStrikethrough2CloseMarker", qspec.Int, ast.NodeStrikethrough2CloseMarker),
		I.Const("NodeStrikethrough2OpenMarker", qspec.Int, ast.NodeStrikethrough2OpenMarker),
		I.Const("NodeStrong", qspec.Int, ast.NodeStrong),
		I.Const("NodeStrongA6kCloseMarker", qspec.Int, ast.NodeStrongA6kCloseMarker),
		I.Const("NodeStrongA6kOpenMarker", qspec.Int, ast.NodeStrongA6kOpenMarker),
		I.Const("NodeStrongU8eCloseMarker", qspec.Int, ast.NodeStrongU8eCloseMarker),
		I.Const("NodeStrongU8eOpenMarker", qspec.Int, ast.NodeStrongU8eOpenMarker),
		I.Const("NodeSub", qspec.Int, ast.NodeSub),
		I.Const("NodeSubCloseMarker", qspec.Int, ast.NodeSubCloseMarker),
		I.Const("NodeSubOpenMarker", qspec.Int, ast.NodeSubOpenMarker),
		I.Const("NodeSup", qspec.Int, ast.NodeSup),
		I.Const("NodeSupCloseMarker", qspec.Int, ast.NodeSupCloseMarker),
		I.Const("NodeSupOpenMarker", qspec.Int, ast.NodeSupOpenMarker),
		I.Const("NodeSuperBlock", qspec.Int, ast.NodeSuperBlock),
		I.Const("NodeSuperBlockCloseMarker", qspec.Int, ast.NodeSuperBlockCloseMarker),
		I.Const("NodeSuperBlockLayoutMarker", qspec.Int, ast.NodeSuperBlockLayoutMarker),
		I.Const("NodeSuperBlockOpenMarker", qspec.Int, ast.NodeSuperBlockOpenMarker),
		I.Const("NodeTable", qspec.Int, ast.NodeTable),
		I.Const("NodeTableCell", qspec.Int, ast.NodeTableCell),
		I.Const("NodeTableHead", qspec.Int, ast.NodeTableHead),
		I.Const("NodeTableRow", qspec.Int, ast.NodeTableRow),
		I.Const("NodeTag", qspec.Int, ast.NodeTag),
		I.Const("NodeTagCloseMarker", qspec.Int, ast.NodeTagCloseMarker),
		I.Const("NodeTagOpenMarker", qspec.Int, ast.NodeTagOpenMarker),
		I.Const("NodeTaskListItemMarker", qspec.Int, ast.NodeTaskListItemMarker),
		I.Const("NodeText", qspec.Int, ast.NodeText),
		I.Const("NodeThematicBreak", qspec.Int, ast.NodeThematicBreak),
		I.Const("NodeToC", qspec.Int, ast.NodeToC),
		I.Const("NodeTypeMaxVal", qspec.Int, ast.NodeTypeMaxVal),
		I.Const("NodeVditorCaret", qspec.Int, ast.NodeVditorCaret),
		I.Const("NodeYamlFrontMatter", qspec.Int, ast.NodeYamlFrontMatter),
		I.Const("NodeYamlFrontMatterCloseMarker", qspec.Int, ast.NodeYamlFrontMatterCloseMarker),
		I.Const("NodeYamlFrontMatterContent", qspec.Int, ast.NodeYamlFrontMatterContent),
		I.Const("NodeYamlFrontMatterOpenMarker", qspec.Int, ast.NodeYamlFrontMatterOpenMarker),
		I.Const("WalkContinue", qspec.ConstUnboundInt, ast.WalkContinue),
		I.Const("WalkSkipChildren", qspec.ConstUnboundInt, ast.WalkSkipChildren),
		I.Const("WalkStop", qspec.ConstUnboundInt, ast.WalkStop),
	)
}
