// Code generated by 'yaegi extract github.com/456vv/vbody/v2'. DO NOT EDIT.

//go:build yaegi_lib
// +build yaegi_lib

package yaegi_lib

import (
	"github.com/456vv/vbody/v2"
	"reflect"
)

func init() {
	Symbols["github.com/456vv/vbody/v2/vbody"] = map[string]reflect.Value{
		// function, constant and variable definitions
		"NewReader": reflect.ValueOf(vbody.NewReader),
		"NewWriter": reflect.ValueOf(vbody.NewWriter),

		// type definitions
		"Change": reflect.ValueOf((*vbody.Change)(nil)),
		"Reader": reflect.ValueOf((*vbody.Reader)(nil)),
		"Writer": reflect.ValueOf((*vbody.Writer)(nil)),
	}
}
