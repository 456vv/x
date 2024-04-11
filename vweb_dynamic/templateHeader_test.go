package vweb_dynamic

import (
	"testing"
)

func Test_templateHeader1(t *testing.T) {
	lines := []string{
		"file=a1.v",
		"file=b1.v",
	}
	fh := templateHeader(lines)
	if len(fh.File) != 2 {
		t.Fatalf("引入文件长度不正确")
	}
}

func Test_TemplateHeader_OpenFile(t *testing.T) {
	th := &TemplateHeader{
		File: []string{
			"a1.v",
			"b1.v",
		},
	}
	fbs, err := th.OpenFile("./testdata/wwwroot", "/template/index.v")
	if err != nil {
		t.Fatal(err)
	}
	if len(fbs) != 2 {
		t.Fatal("读取模板头文件数量不正确")
	}
}
