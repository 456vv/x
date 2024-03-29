// Code generated by 'yaegi extract github.com/456vv/vconn'. DO NOT EDIT.

//go:build yaegi_lib
// +build yaegi_lib

package yaegi_lib

import (
	"github.com/456vv/vconn"
	"reflect"
)

func init() {
	Symbols["github.com/456vv/vconn/vconn"] = map[string]reflect.Value{
		// function, constant and variable definitions
		"New":     reflect.ValueOf(vconn.New),
		"NewConn": reflect.ValueOf(vconn.NewConn),

		// type definitions
		"CloseNotifier": reflect.ValueOf((*vconn.CloseNotifier)(nil)),
		"Conn":          reflect.ValueOf((*vconn.Conn)(nil)),

		// interface wrapper definitions
		"_CloseNotifier": reflect.ValueOf((*_github_com_456vv_vconn_CloseNotifier)(nil)),
	}
}

// _github_com_456vv_vconn_CloseNotifier is an interface wrapper for CloseNotifier type
type _github_com_456vv_vconn_CloseNotifier struct {
	IValue       interface{}
	WCloseNotify func() <-chan error
}

func (W _github_com_456vv_vconn_CloseNotifier) CloseNotify() <-chan error {
	return W.WCloseNotify()
}
