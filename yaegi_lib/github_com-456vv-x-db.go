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
		"ColumnArray":      reflect.ValueOf(db.ColumnArray),
		"ColumnArrayNil":   reflect.ValueOf(db.ColumnArrayNil),
		"ColumnFilter":     reflect.ValueOf(db.ColumnFilter),
		"DataToContextKey": reflect.ValueOf(&db.DataToContextKey).Elem(),
		"TypeToContextKey": reflect.ValueOf(&db.TypeToContextKey).Elem(),

		// type definitions
		"DB":    reflect.ValueOf((*db.DB)(nil)),
		"Rower": reflect.ValueOf((*db.Rower)(nil)),

		// interface wrapper definitions
		"_Rower": reflect.ValueOf((*_github_com_456vv_x_db_Rower)(nil)),
	}
}

// _github_com_456vv_x_db_Rower is an interface wrapper for Rower type
type _github_com_456vv_x_db_Rower struct {
	IValue interface{}
	WErr   func() error
	WScan  func(dest ...interface{}) error
}

func (W _github_com_456vv_x_db_Rower) Err() error {
	return W.WErr()
}
func (W _github_com_456vv_x_db_Rower) Scan(dest ...interface{}) error {
	return W.WScan(dest...)
}
