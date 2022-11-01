package vweb_dynamic

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/goplus/igop"
	"github.com/issue9/assert/v2"
)

func Test_igop_1(t *testing.T) {
	as := assert.New(t, true)
	ssa := &Igop{}

	pRoot := "./testdata/wwwroot/igop"
	pPath := "pkg1.go"
	ssa.SetPath(pRoot, pPath)

	fPath := filepath.Join(pRoot, pPath)
	err := ssa.ParseFile(fPath)
	as.NotError(err)
	err = ssa.Execute(os.Stdout, nil)
	as.NotError(err)
}

func Test_igop_2(t *testing.T) {
	as := assert.New(t, true)
	ctx := igop.NewContext(0)
	ctx.Lookup = func(root, path string) (dir string, found bool) {
		dir = filepath.Join(root, "vendor", path)
		if info, err := os.Stat(dir); err == nil && info.IsDir() {
			found = true
		}
		return
	}
	_, err := ctx.Run("./testdata/wwwroot/igop/pkg1.go", nil)
	as.NotError(err)
}
