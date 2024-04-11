package vweb_dynamic

import (
	"bytes"
	"fmt"
	"path/filepath"
	"testing"
	"time"
)

func Test_Yaegi(t *testing.T) {
	yaegi := &Yaegi{}
	yaegi.SetPath("./testdata/wwwroot", "yagei/y1.go")

	if err := yaegi.ParseFile(filepath.Join(yaegi.rootPath, yaegi.pagePath)); err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 1000; i++ {
		go func(i int) {
			buf := bytes.NewBuffer(nil)
			if err := yaegi.Execute(buf, i); err != nil {
				panic(err)
			}
			if buf.String() != fmt.Sprint(i) {
				panic("error")
			}
		}(i)
	}
	time.Sleep(time.Second * 2)
}
