// export by github.com/goplus/igop/cmd/qexp

//go:build igop_lib
// +build igop_lib

package watch

import (
	q "github.com/456vv/x/watch"

	"reflect"

	"github.com/goplus/igop"
)

func init() {
	igop.RegisterPackage(&igop.Package{
		Name: "watch",
		Path: "github.com/456vv/x/watch",
		Deps: map[string]string{
			"github.com/fsnotify/fsnotify": "fsnotify",
			"path/filepath":                "filepath",
			"strings":                      "strings",
			"sync":                         "sync",
			"time":                         "time",
		},
		Interfaces: map[string]reflect.Type{},
		NamedTypes: map[string]igop.NamedType{
			"Watch":         {reflect.TypeOf((*q.Watch)(nil)).Elem(), "", "Close,Monitor,Remove,event"},
			"WatchEventFun": {reflect.TypeOf((*q.WatchEventFun)(nil)).Elem(), "", ""},
		},
		AliasTypes: map[string]reflect.Type{},
		Vars:       map[string]reflect.Value{},
		Funcs: map[string]reflect.Value{
			"NewWatch": reflect.ValueOf(q.NewWatch),
		},
		TypedConsts:   map[string]igop.TypedConst{},
		UntypedConsts: map[string]igop.UntypedConst{},
	})
}
