package vweb_dynamic

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func Test_Yaegi(t *testing.T) {
	yaegi := &Yaegi{}
	yaegi.SetPath("./testdata/wwwroot", "yagei/y1.yg")

	if err := yaegi.ParseFile(filepath.Join(yaegi.rootPath, yaegi.pagePath)); err != nil {
		t.Fatal(err)
	}

	go func() {
		a := 11
		if err := yaegi.Execute(os.Stdout, &a); err != nil {
			panic(err)
		}
	}()
	go func() {
		var a int = 12
		if err := yaegi.Execute(os.Stdout, &a); err != nil {
			panic(err)
		}
	}()
	go func() {
		var a int = 13
		if err := yaegi.Execute(os.Stdout, &a); err != nil {
			panic(err)
		}
	}()
	time.Sleep(time.Second * 2)
}
