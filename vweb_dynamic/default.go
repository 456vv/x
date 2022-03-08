package vweb_dynamic

import (
	"github.com/456vv/vweb/v2"
)

func DefaultPlus() map[string]vweb.DynamicTemplateFunc {
	plus := map[string]vweb.DynamicTemplateFunc{
		"ank": vweb.DynamicTemplateFunc(func(D *vweb.ServerHandlerDynamic) vweb.DynamicTemplater {return &Anko{}}),
		"yaegi": vweb.DynamicTemplateFunc(func(D *vweb.ServerHandlerDynamic) vweb.DynamicTemplater {return &Yaegi{}}),
		"template": vweb.DynamicTemplateFunc(func(D *vweb.ServerHandlerDynamic) vweb.DynamicTemplater {return &Template{}}),
		"gossa": vweb.DynamicTemplateFunc(func(D *vweb.ServerHandlerDynamic) vweb.DynamicTemplater {return &Gossa{}}),
	}
	return plus
}