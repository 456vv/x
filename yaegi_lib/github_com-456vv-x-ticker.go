// Code generated by 'yaegi extract github.com/456vv/x/ticker'. DO NOT EDIT.

//go:build yaegi_lib
// +build yaegi_lib

package yaegi_lib

import (
	"github.com/456vv/x/ticker"
	"reflect"
)

func init() {
	Symbols["github.com/456vv/x/ticker/ticker"] = map[string]reflect.Value{
		// function, constant and variable definitions
		"NewTicker": reflect.ValueOf(ticker.NewTicker),

		// type definitions
		"Ticker": reflect.ValueOf((*ticker.Ticker)(nil)),
	}
}