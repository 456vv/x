package vweb_dynamic

import (
	"bytes"
	"strings"
	"testing"
	"text/template"

	"github.com/456vv/vweb/v2"
)

func Test_Template_separation(t *testing.T) {
	tests := []struct {
		content []byte
		err     bool
	}{
		{
			content: []byte("//file=./2.tmpl\r\n" +
				"//file=./3.tmpl\r\n" +
				"//file=/5.tmpl\r\n" +
				"//DelimLeft={{\r\n" +
				"//DelimRight=}}\r\n" +
				"\r\n" +
				"1234567890"),
		}, {
			err: true, content: []byte("//file=./2.tmpl\r\n" +
				"//file:./3.tmpl\r\n" + // 不正确
				"//file=/5.tmpl\r\n" +
				"//DelimLeft={{\r\n" +
				"//DelimRight=}}\r\n" +
				"\r\n" +
				"1234567890"),
		}, {
			content: []byte("//file=./2.tmpl\r\n" +
				"//File=./3.tmpl\r\n" + // 被忽略了
				"//file=/5.tmpl\r\n" +
				"//DelimLeft={{\r\n" +
				"//DelimRight=}}\r\n" +
				"\r\n" +
				"1234567890"),
		}, {
			err: true, content: []byte("//file=./2.tmpl\r\n" +
				"//file=./3.tmpl\r\n" +
				"//file=\r\n" + // 不正确
				"//DelimLeft={{\r\n" +
				"//DelimRight=}}\r\n" +
				"\r\n" +
				"1234567890"),
		}, {
			content: []byte("file=./2.tmpl\r\n" + // 不正确
				"file=./3.tmpl\r\n" + // 不正确
				"//DelimLeft={{\r\n" +
				"//DelimRight=}}\r\n" +
				"\r\n" +
				"1234567890"),
		}, {
			content: []byte("//file=./2.tmpl\r\n" +
				"//file=./3.tmpl\r\n" +
				"//DelimLeft={{\r\n" +
				"//DelimRight=}}\r\n" +
				"\r\n"), // 无内容
		}, {
			err: true, content: []byte("//file=./2.tmpl\r\n" +
				"//file=./3.tmpl\r\n" +
				"//DelimLeft={{\r\n" +
				"//DelimRight=}}\r\n"), // 不正确格式,无内容
		},
	}
	for _, v := range tests {
		shdt := Template{
			rootPath: "./testdata/wwwroot",
			pagePath: "/template/t.bw",
		}
		bytesBuffer := bytes.NewBuffer(v.content)
		h, _, err := TemplateSeparation(bytesBuffer)
		if err != nil {
			if !v.err {
				t.Fatal(err)
			}
			continue
		}

		_, err = h.OpenFile(shdt.rootPath, shdt.pagePath)
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
		shdth  TemplateHeader
		length int
	}{
		{shdth: TemplateHeader{File: []string{"./2.tmpl", "./3.tmpl", "/5.tmpl"}}, length: 3},
		{shdth: TemplateHeader{File: []string{"./2.tmpl", "./3.tmpl", "/6.tmpl"}}, length: 0},    // "/6.tmpl" 该文件不存在
		{shdth: TemplateHeader{File: []string{"./2.tmpl", "/../3.tmpl", "/5.tmpl"}}, length: 0},  // "/../3.tmpl" 等于 "/3.tmpl" ，该文件不存在
		{shdth: TemplateHeader{File: []string{"./2.tmpl", "./../5.tmpl", "/5.tmpl"}}, length: 2}, // "./../5.tmpl" 等于 "/5.tmpl"
		{shdth: TemplateHeader{File: []string{"./2.tmpl", "../5.tmpl", "/5.tmpl"}}, length: 2},   // "../5.tmpl" 等于 "/5.tmpl"
		{shdth: TemplateHeader{File: []string{"./2.tmpl", "../5.tmpl", "/"}}, length: 0},         // "/" 表示是根目录 "./test/wwwroot"，不是文件。
		{shdth: TemplateHeader{File: []string{"./2.tmpl", "../5.tmpl", "../../"}}, length: 0},    // "../../" 表示是根目录 "./test/wwwroot"，因为不能跨越根目录。同时也不是一个有效的文件。
		{shdth: TemplateHeader{File: []string{"./2.tmpl", "3.tmpl", "t.bw"}}, length: 3},
	}
	for index, v := range tests {
		m, err := v.shdth.OpenFile(rootPath, pagePath)
		if len(m) != v.length {
			t.Fatalf("%d %v", index, err)
		}
	}
}

func Test_Template_format(t *testing.T) {
	tests := []struct {
		shdth   TemplateHeader
		content string
		result  string
	}{
		{
			shdth:   TemplateHeader{DelimLeft: "{{", DelimRight: "}}"},
			content: "{{\r\n.\r\n}}1234{{\r\n.\r\n}}",
			result:  "{{.}}1234{{.}}",
		}, {
			shdth:   TemplateHeader{DelimLeft: "{{", DelimRight: "}}"},
			content: "{{\r\n.\r\n}}1234\r\n{{.}}",
			result:  "{{.}}1234\r\n{{.}}",
		}, {
			shdth:   TemplateHeader{DelimLeft: "{{", DelimRight: "}}"},
			content: "{{\r\n.\r\n}}1234{{.}}",
			result:  "{{.}}1234{{.}}",
		}, {
			shdth:   TemplateHeader{DelimLeft: "{{", DelimRight: "}}"},
			content: "{{\r\n.\r\n}}\r\n1234\r\n{{.}}",
			result:  "{{.}}\r\n1234\r\n{{.}}",
		}, {
			shdth:   TemplateHeader{DelimLeft: "#*", DelimRight: "*#"},
			content: "111#*\r\n.\r\n*#3333",
			result:  "111#*.*#3333",
		},
	}
	shdt := &Template{}
	for index, v := range tests {
		content := shdt.format(v.shdth.DelimLeft, v.shdth.DelimRight, v.content)
		if content != v.result {
			t.Fatalf("%d %v", index, content)
		}
	}
}

func Test_Template_loadTmpl(t *testing.T) {
	tests := []struct {
		shdth   TemplateHeader
		content map[string]string
		result  string
		err     bool
	}{
		{
			shdth:   TemplateHeader{DelimLeft: "{{", DelimRight: "}}"},
			content: map[string]string{"1.tmpl": "{{define \"1.tmpl\"}}1111111{{end}}", "2.tmpl": "{{define \"2.tmpl\"}}222222{{end}}"},
		}, {
			shdth:   TemplateHeader{DelimLeft: "{{", DelimRight: "}}"},
			content: map[string]string{"1.tmpl": "{{define \"1.tmpl\"}}1111111{{end}}", "2.tmpl": "{{define \"2.tmpl\"}}222222"},
			err:     true,
		}, {
			shdth:   TemplateHeader{DelimLeft: "{{", DelimRight: "}}"},
			content: map[string]string{"1.tmpl": "{{define \"1.tmpl\"}}1111111{{end}}", "2.tmpl": "222222222"},
		},
	}
	shdt := Template{}
	for index, v := range tests {
		t1 := template.New("test")
		t1.Delims(v.shdth.DelimLeft, v.shdth.DelimRight)
		t1, err := shdt.loadTmpl(t1, v.shdth.DelimLeft, v.shdth.DelimRight, v.content)

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
	err := shdt.ParseText("test", "\r\n{{define \"func\"}}123456{{end}}{{$t := .Context.Value \"Template\"}}{{$f := $t.NewFunc \"func\"}}{{print (NotError $f)}}")
	if err != nil {
		t.Fatal(err)
	}
	buf := bytes.NewBuffer(nil)
	in := new(vweb.TemplateDot)

	err = shdt.Execute(buf, in)
	if err != nil {
		t.Fatal(err)
	}
	if buf.String() != "true" {
		t.Fatalf("错误的结果，true == %s", buf.String())
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
	in := new(vweb.TemplateDot)

	if err := shdt.Execute(buf, in); err != nil {
		t.Fatal(err)
	}
	result := strings.ReplaceAll(buf.String(), "\n", "")
	if result != "[1 2 3 4 5 6]" {
		t.Fatalf("错误的结果，[1 2 3 4 5 6] == %s", result)
	}
}
