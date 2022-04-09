// export by github.com/goplus/gossa/cmd/qexp

//go:build gossa_lib
// +build gossa_lib

package lute

import (
	q "github.com/88250/lute"

	"go/constant"
	"reflect"

	"github.com/goplus/gossa"
)

func init() {
	gossa.RegisterPackage(&gossa.Package{
		Name: "lute",
		Path: "github.com/88250/lute",
		Deps: map[string]string{
			"bytes":                           "bytes",
			"github.com/88250/lute/ast":       "ast",
			"github.com/88250/lute/html":      "html",
			"github.com/88250/lute/html/atom": "atom",
			"github.com/88250/lute/lex":       "lex",
			"github.com/88250/lute/parse":     "parse",
			"github.com/88250/lute/render":    "render",
			"github.com/88250/lute/util":      "util",
			"github.com/gopherjs/gopherjs/js": "js",
			"strconv":                         "strconv",
			"strings":                         "strings",
			"unicode":                         "unicode",
		},
		Interfaces: map[string]reflect.Type{},
		NamedTypes: map[string]gossa.NamedType{
			"Lute":        {reflect.TypeOf((*q.Lute)(nil)).Elem(), "", "BlockDOM2HTML,BlockDOM2Md,BlockDOM2StdMd,BlockDOM2Text,BlockDOM2TextLen,BlockDOM2Tree,Blocks2Blockquote,Blocks2Hs,Blocks2OLs,Blocks2Ps,Blocks2TLs,Blocks2ULs,BlocksMergeSuperBlock,CancelList,CancelSuperBlock,Format,FormatStr,GetEmojis,GetLinkBase,GetTerms,H2P,HLevel,HTML2BlockDOM,HTML2Markdown,HTML2Md,HTML2Text,HTML2Tree,HTML2VditorDOM,HTML2VditorIRDOM,HTML2VditorSVDOM,InlineMd2BlockDOM,IsValidLinkDest,Markdown,MarkdownStr,Md2BlockDOM,Md2HTML,Md2VditorDOM,Md2VditorIRDOM,Md2VditorSVDOM,OL2TL,OL2UL,P2H,ProtylePreview,PutEmojis,PutTerms,RemoveEmoji,RenderEChartsJSON,RenderJSON,RenderKityMinderJSON,SetAutoSpace,SetBlockRef,SetChineseParagraphBeginningSpace,SetCodeSyntaxHighlight,SetCodeSyntaxHighlightDetectLang,SetCodeSyntaxHighlightInlineStyle,SetCodeSyntaxHighlightLineNum,SetCodeSyntaxHighlightStyleName,SetEmoji,SetEmojiSite,SetEmojis,SetFixTermTypo,SetFootnotes,SetGFMAutoLink,SetGFMStrikethrough,SetGFMTable,SetGFMTaskListItem,SetGFMTaskListItemClass,SetGitConflict,SetHeadingAnchor,SetHeadingID,SetImageLazyLoading,SetImgPathAllowSpace,SetIndentCodeBlock,SetInlineMathAllowDigitAfterOpenMarker,SetJSRenderers,SetKramdownBlockIAL,SetKramdownIAL,SetKramdownIALIDRenderName,SetKramdownSpanIAL,SetLinkBase,SetLinkPrefix,SetLinkRef,SetMark,SetProtyleWYSIWYG,SetRenderListStyle,SetSanitize,SetSetext,SetSoftBreak2HardBreak,SetSub,SetSup,SetSuperBlock,SetTag,SetTerms,SetToC,SetVditorCodeBlockPreview,SetVditorHTMLBlockPreview,SetVditorIR,SetVditorMathBlockPreview,SetVditorSV,SetVditorWYSIWYG,SetYamlFrontMatter,Space,SpinBlockDOM,SpinVditorDOM,SpinVditorIRDOM,SpinVditorSVDOM,TL2OL,TL2UL,TextBundle,TextBundleStr,Tree2BlockDOM,Tree2HTML,UL2OL,UL2TL,VditorDOM2HTML,VditorDOM2Md,VditorIRDOM2HTML,VditorIRDOM2Md,adjustVditorDOM,adjustVditorDOMListItemInP,adjustVditorDOMListList,adjustVditorDOMListTight0,blockDOM2Md,domAttrValue,domChild,domCode,domCode0,domCustomAttrs,domHTML,domText,domText0,forwardNextBlock,genASTByBlockDOM,genASTByDOM,genASTByVditorDOM,genASTByVditorIRDOM,genASTContenteditable,hasAttr,hljsSpans,isEmptyText,isInline,isTightList,listItemEnter,mergeSameSpan,mergeVditorDOMList0,parentIs,parseHTML,removeDOMAttr,removeEmptyNodes,removeHighlightJSSpans,searchEmptyNodes,setBlockIAL,setDOMAttrValue,setSpanIAL,startsWithNewline,vditorDOM2Md,vditorIRDOM2Md"},
			"ParseOption": {reflect.TypeOf((*q.ParseOption)(nil)).Elem(), "", ""},
		},
		AliasTypes: map[string]reflect.Type{},
		Vars:       map[string]reflect.Value{},
		Funcs: map[string]reflect.Value{
			"FormatNode":         reflect.ValueOf(q.FormatNode),
			"New":                reflect.ValueOf(q.New),
			"RenderNodeBlockDOM": reflect.ValueOf(q.RenderNodeBlockDOM),
		},
		TypedConsts: map[string]gossa.TypedConst{},
		UntypedConsts: map[string]gossa.UntypedConst{
			"Version": {"untyped string", constant.MakeString(string(q.Version))},
		},
	})
}
