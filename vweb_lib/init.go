package vweb_lib

import (
	"github.com/456vv/vweb/v2"
)

func init(){
	//给template模板增加模块包
	for name, pkg := range templatePackage() {
		vweb.ExtendTemplatePackage(name, pkg)
	}
	for name, pkg := range reflectxTemplatePackage() {
		vweb.ExtendTemplatePackage(name, pkg)
	}
}
