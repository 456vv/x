// export by github.com/goplus/igop/cmd/qexp

//go:build igop_lib
// +build igop_lib

package vmap

import (
	q "github.com/456vv/vmap/v2"

	"reflect"

	"github.com/goplus/igop"
)

func init() {
	igop.RegisterPackage(&igop.Package{
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
		NamedTypes: map[string]igop.NamedType{
			"Map": {reflect.TypeOf((*q.Map)(nil)).Elem(), "", "Copy,Del,Dels,Get,GetHas,GetNewMap,GetNewMaps,GetOrDefault,Has,Index,IndexHas,Keys,Len,MarshalJSON,New,Range,ReadAll,ReadFrom,Reset,Set,SetExpired,SetExpiredCall,String,UnmarshalJSON,Values,WriteTo,afterFunc,marshalJSON,readFrom,unmarshalJSON,writeTo"},
		},
		AliasTypes: map[string]reflect.Type{},
		Vars:       map[string]reflect.Value{},
		Funcs: map[string]reflect.Value{
			"NewMap": reflect.ValueOf(q.NewMap),
		},
		TypedConsts:   map[string]igop.TypedConst{},
		UntypedConsts: map[string]igop.UntypedConst{},
	})
}
