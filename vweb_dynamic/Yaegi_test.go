package vweb_dynamic

import (
	"testing"
	"path/filepath"
	"os"
)

func Test_Yaegi(t *testing.T){
	
	yaegi := &Yaegi{}
	yaegi.SetPath("./test/wwwroot", "yagei/y1.y")
	
	if err := yaegi.ParseFile(filepath.Join(yaegi.rootPath, yaegi.pagePath)); err != nil {
		t.Fatal(err)
	}
	
	var in = map[string]string{
		"a":"a1",
	}
	if err := yaegi.Execute(os.Stdout, in); err != nil {
		t.Fatal(err)
	}
}