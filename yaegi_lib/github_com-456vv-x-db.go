// Code generated by 'yaegi extract github.com/456vv/x/db'. DO NOT EDIT.

//go:build yaegi_lib
// +build yaegi_lib

package yaegi_lib

import (
	"github.com/456vv/x/db"
	"reflect"
)

func init() {
	Symbols["github.com/456vv/x/db/db"] = map[string]reflect.Value{
		// function, constant and variable definitions
		"ErrNoRows": reflect.ValueOf(&db.ErrNoRows).Elem(),
		"ErrRows":   reflect.ValueOf(&db.ErrRows).Elem(),

		// type definitions
		"DB": reflect.ValueOf((*db.DB)(nil)),
	}
}