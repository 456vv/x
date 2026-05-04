// export by github.com/goplus/ixgo/cmd/qexp

//go:build igop_lib
// +build igop_lib

package sqltable

import (
	q "github.com/456vv/x/sqltable"

	"reflect"

	"github.com/goplus/ixgo"
)

func init() {
	ixgo.RegisterPackage(&ixgo.Package{
		Name: "sqltable",
		Path: "github.com/456vv/x/sqltable",
		Deps: map[string]string{
			"fmt":     "fmt",
			"reflect": "reflect",
			"strconv": "strconv",
			"strings": "strings",
		},
		Interfaces: map[string]reflect.Type{},
		NamedTypes: map[string]reflect.Type{
			"SQLTable": reflect.TypeOf((*q.SQLTable)(nil)).Elem(),
		},
		AliasTypes: map[string]reflect.Type{},
		Vars:       map[string]reflect.Value{},
		Funcs: map[string]reflect.Value{
			"Prepare": reflect.ValueOf(q.Prepare),
		},
		TypedConsts:   map[string]ixgo.TypedConst{},
		UntypedConsts: map[string]ixgo.UntypedConst{},
	})
}
