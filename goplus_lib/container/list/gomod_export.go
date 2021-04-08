// Package list provide Go+ "container/list" package, as "container/list" package in Go.
package list

import (
	list "container/list"
	reflect "reflect"

	gop "github.com/goplus/gop"
)

func execmElementNext(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*list.Element).Next()
	p.Ret(1, ret0)
}

func execmElementPrev(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*list.Element).Prev()
	p.Ret(1, ret0)
}

func execmListInit(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*list.List).Init()
	p.Ret(1, ret0)
}

func execmListLen(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*list.List).Len()
	p.Ret(1, ret0)
}

func execmListFront(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*list.List).Front()
	p.Ret(1, ret0)
}

func execmListBack(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*list.List).Back()
	p.Ret(1, ret0)
}

func execmListRemove(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*list.List).Remove(args[1].(*list.Element))
	p.Ret(2, ret0)
}

func execmListPushFront(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*list.List).PushFront(args[1])
	p.Ret(2, ret0)
}

func execmListPushBack(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*list.List).PushBack(args[1])
	p.Ret(2, ret0)
}

func execmListInsertBefore(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*list.List).InsertBefore(args[1], args[2].(*list.Element))
	p.Ret(3, ret0)
}

func execmListInsertAfter(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*list.List).InsertAfter(args[1], args[2].(*list.Element))
	p.Ret(3, ret0)
}

func execmListMoveToFront(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*list.List).MoveToFront(args[1].(*list.Element))
	p.PopN(2)
}

func execmListMoveToBack(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*list.List).MoveToBack(args[1].(*list.Element))
	p.PopN(2)
}

func execmListMoveBefore(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	args[0].(*list.List).MoveBefore(args[1].(*list.Element), args[2].(*list.Element))
	p.PopN(3)
}

func execmListMoveAfter(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	args[0].(*list.List).MoveAfter(args[1].(*list.Element), args[2].(*list.Element))
	p.PopN(3)
}

func execmListPushBackList(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*list.List).PushBackList(args[1].(*list.List))
	p.PopN(2)
}

func execmListPushFrontList(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*list.List).PushFrontList(args[1].(*list.List))
	p.PopN(2)
}

func execNew(_ int, p *gop.Context) {
	ret0 := list.New()
	p.Ret(0, ret0)
}

// I is a Go package instance.
var I = gop.NewGoPackage("container/list")

func init() {
	I.RegisterFuncs(
		I.Func("(*Element).Next", (*list.Element).Next, execmElementNext),
		I.Func("(*Element).Prev", (*list.Element).Prev, execmElementPrev),
		I.Func("(*List).Init", (*list.List).Init, execmListInit),
		I.Func("(*List).Len", (*list.List).Len, execmListLen),
		I.Func("(*List).Front", (*list.List).Front, execmListFront),
		I.Func("(*List).Back", (*list.List).Back, execmListBack),
		I.Func("(*List).Remove", (*list.List).Remove, execmListRemove),
		I.Func("(*List).PushFront", (*list.List).PushFront, execmListPushFront),
		I.Func("(*List).PushBack", (*list.List).PushBack, execmListPushBack),
		I.Func("(*List).InsertBefore", (*list.List).InsertBefore, execmListInsertBefore),
		I.Func("(*List).InsertAfter", (*list.List).InsertAfter, execmListInsertAfter),
		I.Func("(*List).MoveToFront", (*list.List).MoveToFront, execmListMoveToFront),
		I.Func("(*List).MoveToBack", (*list.List).MoveToBack, execmListMoveToBack),
		I.Func("(*List).MoveBefore", (*list.List).MoveBefore, execmListMoveBefore),
		I.Func("(*List).MoveAfter", (*list.List).MoveAfter, execmListMoveAfter),
		I.Func("(*List).PushBackList", (*list.List).PushBackList, execmListPushBackList),
		I.Func("(*List).PushFrontList", (*list.List).PushFrontList, execmListPushFrontList),
		I.Func("New", list.New, execNew),
	)
	I.RegisterTypes(
		I.Type("Element", reflect.TypeOf((*list.Element)(nil)).Elem()),
		I.Type("List", reflect.TypeOf((*list.List)(nil)).Elem()),
	)
}
