// export by github.com/goplus/gossa/cmd/qexp

//go:build gossa_lib
// +build gossa_lib

package db

import (
	q "github.com/456vv/x/db"

	"reflect"

	"github.com/goplus/gossa"
)

func init() {
	gossa.RegisterPackage(&gossa.Package{
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
		NamedTypes: map[string]gossa.NamedType{
			"DB": {reflect.TypeOf((*q.DB)(nil)).Elem(), "", "Close,Exec,ExecContext,Has,HasContext,Open,Pexec,PexecContext,Query,QueryContext,QueryRow,QueryRowContext,debugFormat,debugPrint,logf,txCommit"},
		},
		AliasTypes: map[string]reflect.Type{},
		Vars: map[string]reflect.Value{
			"ErrNoRows": reflect.ValueOf(&q.ErrNoRows),
			"ErrRows":   reflect.ValueOf(&q.ErrRows),
		},
		Funcs:         map[string]reflect.Value{},
		TypedConsts:   map[string]gossa.TypedConst{},
		UntypedConsts: map[string]gossa.UntypedConst{},
	})
}
