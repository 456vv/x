// export by github.com/goplus/igop/cmd/qexp

//go:build igop_lib
// +build igop_lib

package sqltable

import (
	q "github.com/456vv/x/sqltable"

	"reflect"

	"github.com/goplus/igop"
)

func init() {
	igop.RegisterPackage(&igop.Package{
		Name: "sqltable",
		Path: "github.com/456vv/x/sqltable",
		Deps: map[string]string{
			"fmt":     "fmt",
			"reflect": "reflect",
			"strconv": "strconv",
			"strings": "strings",
		},
		Interfaces: map[string]reflect.Type{},
		NamedTypes: map[string]igop.NamedType{
			"SQLTable": {reflect.TypeOf((*q.SQLTable)(nil)).Elem(), "", "And,Ands,Args,ColumnMark,Excluded,ExtArgs,Or,Ors,Prepare,SQL,Set,SetValues,Sets,ToTypes,Value,ValueValues,Values,Where,addColumnMark,assemble,assembleStatement,excluded,more,pairArgs,progress,set,values,where"},
		},
		AliasTypes: map[string]reflect.Type{},
		Vars:       map[string]reflect.Value{},
		Funcs: map[string]reflect.Value{
			"NewSQLTable": reflect.ValueOf(q.NewSQLTable),
			"Prepare":     reflect.ValueOf(q.Prepare),
		},
		TypedConsts:   map[string]igop.TypedConst{},
		UntypedConsts: map[string]igop.UntypedConst{},
	})
}
