// export by github.com/goplus/ixgo/cmd/qexp

//go:build igop_lib
// +build igop_lib

package watch

import (
	q "github.com/456vv/x/watch"

	"reflect"

	"github.com/goplus/ixgo"
)

func init() {
	ixgo.RegisterPackage(&ixgo.Package{
		Name: "watch",
		Path: "github.com/456vv/x/watch",
		Deps: map[string]string{
			"fmt":                          "fmt",
			"github.com/fsnotify/fsnotify": "fsnotify",
			"path/filepath":                "filepath",
			"strings":                      "strings",
			"sync":                         "sync",
			"time":                         "time",
		},
		Interfaces: map[string]reflect.Type{},
		NamedTypes: map[string]reflect.Type{
			"Watch":         reflect.TypeOf((*q.Watch)(nil)).Elem(),
			"WatchEventFun": reflect.TypeOf((*q.WatchEventFun)(nil)).Elem(),
		},
		AliasTypes: map[string]reflect.Type{},
		Vars:       map[string]reflect.Value{},
		Funcs: map[string]reflect.Value{
			"NewWatch": reflect.ValueOf(q.NewWatch),
		},
		TypedConsts:   map[string]ixgo.TypedConst{},
		UntypedConsts: map[string]ixgo.UntypedConst{},
	})
}
