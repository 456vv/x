//go:build vweb_lib
// +build vweb_lib

package vweb_lib

import (
	"github.com/456vv/x/vweb_dynamic"
)

func init(){
	//给template模板增加模块包
	for name, pkg := range templatePackage() {
		vweb_dynamic.ExtendPackage(name, pkg)
	}
}
