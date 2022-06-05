package vweb_dynamic

import (
	"github.com/456vv/viot/v2"
	"github.com/456vv/vweb/v2"
)

func WEBModule() map[string]vweb.DynamicTemplateFunc {
	module := map[string]vweb.DynamicTemplateFunc{
		"ank":      vweb.DynamicTemplateFunc(func(D *vweb.ServerHandlerDynamic) vweb.DynamicTemplater { return &Anko{} }),
		"yaegi":    vweb.DynamicTemplateFunc(func(D *vweb.ServerHandlerDynamic) vweb.DynamicTemplater { return &Yaegi{} }),
		"template": vweb.DynamicTemplateFunc(func(D *vweb.ServerHandlerDynamic) vweb.DynamicTemplater { return &Template{} }),
		"igop":     vweb.DynamicTemplateFunc(func(D *vweb.ServerHandlerDynamic) vweb.DynamicTemplater { return &Igop{} }),
	}
	return module
}

func IOTModule() map[string]viot.DynamicTemplateFunc {
	module := map[string]viot.DynamicTemplateFunc{
		"ank":      viot.DynamicTemplateFunc(func(D *viot.ServerHandlerDynamic) viot.DynamicTemplater { return &Anko{} }),
		"yaegi":    viot.DynamicTemplateFunc(func(D *viot.ServerHandlerDynamic) viot.DynamicTemplater { return &Yaegi{} }),
		"template": viot.DynamicTemplateFunc(func(D *viot.ServerHandlerDynamic) viot.DynamicTemplater { return &Template{} }),
		"igop":     viot.DynamicTemplateFunc(func(D *viot.ServerHandlerDynamic) viot.DynamicTemplater { return &Igop{} }),
	}
	return module
}
