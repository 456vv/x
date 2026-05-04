// export by github.com/goplus/ixgo/cmd/qexp

//go:build igop_lib
// +build igop_lib

package ticker

import (
	q "github.com/456vv/x/ticker"

	"reflect"

	"github.com/goplus/ixgo"
)

func init() {
	ixgo.RegisterPackage(&ixgo.Package{
		Name: "ticker",
		Path: "github.com/456vv/x/ticker",
		Deps: map[string]string{
			"time": "time",
		},
		Interfaces: map[string]reflect.Type{},
		NamedTypes: map[string]reflect.Type{
			"Ticker": reflect.TypeOf((*q.Ticker)(nil)).Elem(),
		},
		AliasTypes: map[string]reflect.Type{},
		Vars:       map[string]reflect.Value{},
		Funcs: map[string]reflect.Value{
			"NewTicker": reflect.ValueOf(q.NewTicker),
		},
		TypedConsts:   map[string]ixgo.TypedConst{},
		UntypedConsts: map[string]ixgo.UntypedConst{},
	})
}
