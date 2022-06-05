// export by github.com/goplus/igop/cmd/qexp

//go:build igop_lib
// +build igop_lib

package vbody

import (
	q "github.com/456vv/vbody"

	"reflect"

	"github.com/goplus/igop"
)

func init() {
	igop.RegisterPackage(&igop.Package{
		Name: "vbody",
		Path: "github.com/456vv/vbody",
		Deps: map[string]string{
			"bytes":         "bytes",
			"encoding/json": "json",
			"errors":        "errors",
			"fmt":           "fmt",
			"io":            "io",
			"net/http":      "http",
			"reflect":       "reflect",
			"sync":          "sync",
			"text/template": "template",
		},
		Interfaces: map[string]reflect.Type{},
		NamedTypes: map[string]igop.NamedType{
			"Change": {reflect.TypeOf((*q.Change)(nil)).Elem(), "", "Delete,Set,Update"},
			"Reader": {reflect.TypeOf((*q.Reader)(nil)).Elem(), "", "Array,Bool,BoolAnyEqual,Change,Err,Float64,Float64AnyEqual,Has,Index,IndexArray,IndexFloat64,IndexInt64,IndexString,Int64,Int64AnyEqual,Interface,IsNil,MarshalJSON,NewArray,NewIndex,NewIndexArray,NewInterface,NewSlice,NoZero,ReadFrom,Reset,Slice,String,StringAnyEqual,UnmarshalJSON,has,isNil"},
			"Writer": {reflect.TypeOf((*q.Writer)(nil)).Elem(), "", "Message,Messagef,Result,SetResult,Status,String,WriteTo,ready"},
		},
		AliasTypes: map[string]reflect.Type{},
		Vars:       map[string]reflect.Value{},
		Funcs: map[string]reflect.Value{
			"NewReader": reflect.ValueOf(q.NewReader),
			"NewWriter": reflect.ValueOf(q.NewWriter),
		},
		TypedConsts:   map[string]igop.TypedConst{},
		UntypedConsts: map[string]igop.UntypedConst{},
	})
}
