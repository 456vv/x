package vweb_dynamic

import (
	"bytes"
	"strings"
	"testing"
	"text/template"

	"github.com/456vv/vweb/v2"
)

func Test_templateHeader(t *testing.T) {
	tests := []struct {
		content []string
		length  int
	}{
		{
			length: 3, content: []string{"file=./2.tmpl", "file=./3.tmpl", "file=/5.tmpl"},
		}, {
			length: 2, content: []string{"file=./2.tmpl", "file=/5.tmpl", "//File=./3.tmpl"},
		}, {
			length: 2, content: []string{"file=./2.tmpl", "file=./3.tmpl", "file="},
		},
	}
	for _, v := range tests {
		th := templateHeader(v.content)
		if len(th.File) != v.length {
			t.Fatal("error")
		}

		tl := Template{
			rootPath: "./testdata/wwwroot",
			pagePath: "/template/t.bw",
		}
		_, err := th.OpenFile(tl.rootPath, tl.pagePath)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func Test_shdtHeader_openFile(t *testing.T) {
	var (
		rootPath = "./testdata/wwwroot"
		pagePath = "/template/1.tmpl"
	)
	tests := []struct {
		th     TemplateHeader
		length int
	}{
		{th: TemplateHeader{File: []string{"./2.tmpl", "./3.tmpl", "/5.tmpl"}}, length: 3},
		{th: TemplateHeader{File: []string{"./2.tmpl", "./3.tmpl", "/6.tmpl"}}, length: 0},    // "/6.tmpl" 该文件不存在
		{th: TemplateHeader{File: []string{"./2.tmpl", "/../3.tmpl", "/5.tmpl"}}, length: 0},  // "/../3.tmpl" 等于 "/3.tmpl" ，该文件不存在
		{th: TemplateHeader{File: []string{"./2.tmpl", "./../5.tmpl", "/5.tmpl"}}, length: 2}, // "./../5.tmpl" 等于 "/5.tmpl"
		{th: TemplateHeader{File: []string{"./2.tmpl", "../5.tmpl", "/5.tmpl"}}, length: 2},   // "../5.tmpl" 等于 "/5.tmpl"
		{th: TemplateHeader{File: []string{"./2.tmpl", "../5.tmpl", "/"}}, length: 0},         // "/" 表示是根目录 "./test/wwwroot"，不是文件。
		{th: TemplateHeader{File: []string{"./2.tmpl", "../5.tmpl", "../../"}}, length: 0},    // "../../" 表示是根目录 "./test/wwwroot"，因为不能跨越根目录。同时也不是一个有效的文件。
		{th: TemplateHeader{File: []string{"./2.tmpl", "3.tmpl", "t.bw"}}, length: 3},
	}
	for index, v := range tests {
		m, err := v.th.OpenFile(rootPath, pagePath)
		if len(m) != v.length {
			t.Fatalf("%d %v", index, err)
		}
	}
}

func Test_Template_format(t *testing.T) {
	tests := []struct {
		th      TemplateHeader
		content string
		result  string
	}{
		{
			th:      TemplateHeader{DelimLeft: "{{", DelimRight: "}}"},
			content: "{{\n.\n}}1234{{\n.\n}}",
			result:  "{{.}}1234{{.}}",
		}, {
			th:      TemplateHeader{DelimLeft: "{{", DelimRight: "}}"},
			content: "{{\n.\n}}1234\n{{.}}",
			result:  "{{.}}1234\n{{.}}",
		}, {
			th:      TemplateHeader{DelimLeft: "{{", DelimRight: "}}"},
			content: "{{\n.\n}}1234{{.}}",
			result:  "{{.}}1234{{.}}",
		}, {
			th:      TemplateHeader{DelimLeft: "{{", DelimRight: "}}"},
			content: "{{\r\n.\r\n}}\r\n1234\n{{.}}",
			result:  "{{.}}\r\n1234\n{{.}}",
		}, {
			th:      TemplateHeader{DelimLeft: "#*", DelimRight: "*#"},
			content: "111#*\n.\n*#3333",
			result:  "111#*.*#3333",
		},
	}
	tl := &Template{}
	for index, v := range tests {
		content := tl.format(v.th.DelimLeft, v.th.DelimRight, v.content)
		if content != v.result {
			t.Fatalf("%d %v", index, content)
		}
	}
}

func Test_Template_loadTmpl(t *testing.T) {
	tests := []struct {
		th      TemplateHeader
		content map[string]string
		result  string
		err     bool
	}{
		{
			th:      TemplateHeader{DelimLeft: "{{", DelimRight: "}}"},
			content: map[string]string{"1.tmpl": "{{define \"1.tmpl\"}}1111111{{end}}", "2.tmpl": "{{define \"2.tmpl\"}}222222{{end}}"},
		}, {
			th:      TemplateHeader{DelimLeft: "{{", DelimRight: "}}"},
			content: map[string]string{"1.tmpl": "{{define \"1.tmpl\"}}1111111{{end}}", "2.tmpl": "{{define \"2.tmpl\"}}222222"},
			err:     true,
		}, {
			th:      TemplateHeader{DelimLeft: "{{", DelimRight: "}}"},
			content: map[string]string{"1.tmpl": "{{define \"1.tmpl\"}}1111111{{end}}", "2.tmpl": "222222222"},
		},
	}
	shdt := Template{}
	for index, v := range tests {
		t1 := template.New("test")
		t1.Delims(v.th.DelimLeft, v.th.DelimRight)
		t1, err := shdt.loadTmpl(t1, v.th.DelimLeft, v.th.DelimRight, v.content)

		if err != nil && !v.err {
			t.Fatalf("%d 加载模板(%s)，错误：%v", index, v.content, err)
		}
		if err != nil {
			continue
		}
		ts := t1.Templates()
		if len(ts) != len(v.content) {
			t.Fatalf("%d %v", index, ts)
		}
	}
}

func Test_TemplateExtend_NewFunc(t *testing.T) {
	// 仅支持本地测试,需要替换text/template 中的文件，在本目录下的patch目录可以找到有关文件'
	return
	var shdt Template
	err := shdt.ParseText("test", "{{define \"func\"}}123456{{end}}{{$t := .Context.Value \"Template\"}}{{$f := $t.NewFunc \"func\"}}{{print (NotError $f)}}")
	if err != nil {
		t.Fatal(err)
	}
	buf := bytes.NewBuffer(nil)
	in := new(vweb.Dot)

	err = shdt.Execute(buf, in)
	if err != nil {
		t.Fatal(err)
	}
	if buf.String() != "true" {
		t.Fatalf("错误的结果，true != %s", buf.String())
	}
}

func Test_TemplateExtend_Call(t *testing.T) {
	// 仅支持本地测试,需要替换text/template 中的文件，在本目录下的patch目录可以找到有关文件
	return
	text := `
{{define "func"}}{{CallMethod . "Result" (.Args -1)}}{{end}}{{$t := .Context.Value "Template"}}{{$f := $t.NewFunc "func"}}{{$rets := $t.Call $f 1 2 3 4 5 6}}{{print $rets}}`

	var shdt Template
	if err := shdt.ParseText("test", text); err != nil {
		t.Fatal(err)
	}
	buf := bytes.NewBuffer(nil)
	in := new(vweb.Dot)

	if err := shdt.Execute(buf, in); err != nil {
		t.Fatal(err)
	}
	result := strings.ReplaceAll(buf.String(), "\n", "")
	if result != "[1 2 3 4 5 6]" {
		t.Fatalf("错误的结果，[1 2 3 4 5 6] == %s", result)
	}
}
