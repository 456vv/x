// Code generated by 'yaegi extract github.com/456vv/x/watch'. DO NOT EDIT.

//go:build yaegi_lib
// +build yaegi_lib

package yaegi_lib

import (
	"github.com/456vv/x/watch"
	"reflect"
)

func init() {
	Symbols["github.com/456vv/x/watch/watch"] = map[string]reflect.Value{
		// function, constant and variable definitions
		"NewWatch": reflect.ValueOf(watch.NewWatch),

		// type definitions
		"Watch":         reflect.ValueOf((*watch.Watch)(nil)),
		"WatchEventFun": reflect.ValueOf((*watch.WatchEventFun)(nil)),
	}
}
