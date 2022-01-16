//go:build vweb_lib
// +build vweb_lib

package vweb_dynamic

import (
	"github.com/456vv/x/vweb_lib"
)
func init(){
	for name, pkg := range vweb_lib.Symbols {
		ExtendPackage(name, pkg)
	}
}