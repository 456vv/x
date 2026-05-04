package vweb_dynamic

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/goplus/ixgo"
	"github.com/issue9/assert/v4"
)

func Test_igop_1(t *testing.T) {
	as := assert.New(t, true)
	ssa := &Ixgo{}

	pRoot := "./testdata/wwwroot/igop"
	pPath := "pkg1.go"
	ssa.SetPath(pRoot, pPath)

	fPath := filepath.Join(pRoot, pPath)
	err := ssa.ParseFile(fPath)
	as.NotError(err)

	var wg sync.WaitGroup
	for i := range 100 {
		as.Go(func(as *assert.Assertion) {
			defer wg.Done()
			buf := bytes.NewBuffer(nil)
			err = ssa.Execute(buf, i)
			as.NotError(err).Equal(buf.String(), fmt.Sprintf("%d+receive", i))
		})
		wg.Add(1)
	}
	wg.Wait()
}

func Test_igop_2(t *testing.T) {
	as := assert.New(t, true)
	ctx := ixgo.NewContext(0)
	ctx.Lookup = func(root, path string) (dir string, found bool) {
		dir = filepath.Join(root, "src", path)
		if info, err := os.Stat(dir); err == nil && info.IsDir() {
			found = true
		}
		return
	}
	_, err := ctx.Run("./testdata/wwwroot/igop/pkg1.go", nil)
	as.NotError(err)
}
