// Code generated by 'yaegi extract github.com/456vv/vforward'. DO NOT EDIT.

//go:build yaegi_lib
// +build yaegi_lib

package yaegi_lib

import (
	"github.com/456vv/vforward"
	"reflect"
)

func init() {
	Symbols["github.com/456vv/vforward/vforward"] = map[string]reflect.Value{
		// function, constant and variable definitions
		"DefaultReadBufSize": reflect.ValueOf(vforward.DefaultReadBufSize),

		// type definitions
		"Addr":    reflect.ValueOf((*vforward.Addr)(nil)),
		"D2D":     reflect.ValueOf((*vforward.D2D)(nil)),
		"D2DSwap": reflect.ValueOf((*vforward.D2DSwap)(nil)),
		"L2D":     reflect.ValueOf((*vforward.L2D)(nil)),
		"L2DSwap": reflect.ValueOf((*vforward.L2DSwap)(nil)),
		"L2L":     reflect.ValueOf((*vforward.L2L)(nil)),
		"L2LSwap": reflect.ValueOf((*vforward.L2LSwap)(nil)),
	}
}
