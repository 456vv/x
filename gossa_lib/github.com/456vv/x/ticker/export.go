// export by github.com/goplus/gossa/cmd/qexp

//go:build gossa_lib
// +build gossa_lib

package ticker

import (
	q "github.com/456vv/x/ticker"

	"reflect"

	"github.com/goplus/gossa"
)

func init() {
	gossa.RegisterPackage(&gossa.Package{
		Name: "ticker",
		Path: "github.com/456vv/x/ticker",
		Deps: map[string]string{
			"time": "time",
		},
		Interfaces: map[string]reflect.Type{},
		NamedTypes: map[string]gossa.NamedType{
			"Ticker": {reflect.TypeOf((*q.Ticker)(nil)).Elem(), "", "Func,Reset,Stop"},
		},
		AliasTypes: map[string]reflect.Type{},
		Vars:       map[string]reflect.Value{},
		Funcs: map[string]reflect.Value{
			"NewTicker": reflect.ValueOf(q.NewTicker),
		},
		TypedConsts:   map[string]gossa.TypedConst{},
		UntypedConsts: map[string]gossa.UntypedConst{},
	})
}
