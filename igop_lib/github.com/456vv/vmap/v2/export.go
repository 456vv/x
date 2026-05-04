// export by github.com/goplus/ixgo/cmd/qexp

//go:build igop_lib
// +build igop_lib

package vmap

import (
	q "github.com/456vv/vmap/v2"

	"reflect"

	"github.com/goplus/ixgo"
)

func init() {
	ixgo.RegisterPackage(&ixgo.Package{
		Name: "vmap",
		Path: "github.com/456vv/vmap/v2",
		Deps: map[string]string{
			"encoding/json": "json",
			"fmt":           "fmt",
			"reflect":       "reflect",
			"sync":          "sync",
			"time":          "time",
		},
		Interfaces: map[string]reflect.Type{},
		NamedTypes: map[string]reflect.Type{
			"Map": reflect.TypeOf((*q.Map)(nil)).Elem(),
		},
		AliasTypes: map[string]reflect.Type{},
		Vars:       map[string]reflect.Value{},
		Funcs: map[string]reflect.Value{
			"NewMap": reflect.ValueOf(q.NewMap),
		},
		TypedConsts:   map[string]ixgo.TypedConst{},
		UntypedConsts: map[string]ixgo.UntypedConst{},
	})
}
