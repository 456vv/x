// export by github.com/goplus/igop/cmd/qexp

//go:build igop_lib
// +build igop_lib

package db

import (
	q "github.com/456vv/x/db"

	"reflect"

	"github.com/goplus/igop"
)

func init() {
	igop.RegisterPackage(&igop.Package{
		Name: "db",
		Path: "github.com/456vv/x/db",
		Deps: map[string]string{
			"context":                  "context",
			"database/sql":             "sql",
			"errors":                   "errors",
			"fmt":                      "fmt",
			"github.com/456vv/vweb/v2": "vweb",
			"github.com/lib/pq":        "pq",
			"log":                      "log",
			"reflect":                  "reflect",
			"strings":                  "strings",
			"time":                     "time",
		},
		Interfaces: map[string]reflect.Type{
			"Rower": reflect.TypeOf((*q.Rower)(nil)).Elem(),
		},
		NamedTypes: map[string]reflect.Type{
			"DB": reflect.TypeOf((*q.DB)(nil)).Elem(),
		},
		AliasTypes: map[string]reflect.Type{},
		Vars: map[string]reflect.Value{
			"ErrNoRows": reflect.ValueOf(&q.ErrNoRows),
			"ErrRows":   reflect.ValueOf(&q.ErrRows),
		},
		Funcs:         map[string]reflect.Value{},
		TypedConsts:   map[string]igop.TypedConst{},
		UntypedConsts: map[string]igop.UntypedConst{},
	})
}
