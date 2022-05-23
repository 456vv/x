package vweb_dynamic

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/issue9/assert/v2"
)

func Test_gossa(t *testing.T) {
	as := assert.New(t, true)
	ssa := &Gossa{}

	pRoot := "./test/wwwroot/gossa"
	pPath := "pkg1.go"
	ssa.SetPath(pRoot, pPath)

	fPath := filepath.Join(pRoot, pPath)
	err := ssa.ParseFile(fPath)
	as.NotError(err)
	err = ssa.Execute(os.Stdout, nil)
	as.NotError(err)
}
