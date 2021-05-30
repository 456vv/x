// Package watch provide Go+ "github.com/456vv/x/watch" package, as "github.com/456vv/x/watch" package in Go.
package watch

import (
	reflect "reflect"

	watch "github.com/456vv/x/watch"
	gop "github.com/goplus/gop"
)

func execNewWatch(_ int, p *gop.Context) {
	ret0, ret1 := watch.NewWatch()
	p.Ret(0, ret0, ret1)
}

func execmWatchRemove(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*watch.Watch).Remove(args[1].(string))
	p.PopN(2)
}

func execmWatchMonitor(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*watch.Watch).Monitor(args[1].(string), args[2].(watch.WatchEventFun))
	p.Ret(3, ret0)
}

func execmWatchClose(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*watch.Watch).Close()
	p.Ret(1, ret0)
}

// I is a Go package instance.
var I = gop.NewGoPackage("github.com/456vv/x/watch")

func init() {
	I.RegisterFuncs(
		I.Func("NewWatch", watch.NewWatch, execNewWatch),
		I.Func("(*Watch).Remove", (*watch.Watch).Remove, execmWatchRemove),
		I.Func("(*Watch).Monitor", (*watch.Watch).Monitor, execmWatchMonitor),
		I.Func("(*Watch).Close", (*watch.Watch).Close, execmWatchClose),
	)
	I.RegisterTypes(
		I.Type("Watch", reflect.TypeOf((*watch.Watch)(nil)).Elem()),
		I.Type("WatchEventFun", reflect.TypeOf((*watch.WatchEventFun)(nil)).Elem()),
	)
}
