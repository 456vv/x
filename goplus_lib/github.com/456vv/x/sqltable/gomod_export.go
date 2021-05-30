// Package sqltable provide Go+ "github.com/456vv/x/sqltable" package, as "github.com/456vv/x/sqltable" package in Go.
package sqltable

import (
	reflect "reflect"

	sqltable "github.com/456vv/x/sqltable"
	gop "github.com/goplus/gop"
)

func execNewSQLTable(_ int, p *gop.Context) {
	ret0 := sqltable.NewSQLTable()
	p.Ret(0, ret0)
}

func execmSQLTableColumnMark(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*sqltable.SQLTable).ColumnMark(args[1].(string))
	p.Ret(2, ret0)
}

func execmSQLTableValue(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*sqltable.SQLTable).Value(args[1].(string), args[2])
	p.Ret(3, ret0)
}

func execmSQLTableValues(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	ret0 := args[0].(*sqltable.SQLTable).Values(args[1:]...)
	p.Ret(arity, ret0)
}

func execmSQLTableAnd(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*sqltable.SQLTable).And(args[1].(string), args[2])
	p.Ret(3, ret0)
}

func execmSQLTableAnds(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	ret0 := args[0].(*sqltable.SQLTable).Ands(args[1:]...)
	p.Ret(arity, ret0)
}

func execmSQLTableOr(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*sqltable.SQLTable).Or(args[1].(string), args[2])
	p.Ret(3, ret0)
}

func execmSQLTableOrs(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	ret0 := args[0].(*sqltable.SQLTable).Ors(args[1:]...)
	p.Ret(arity, ret0)
}

func execmSQLTableWhere(_ int, p *gop.Context) {
	args := p.GetArgs(4)
	ret0 := args[0].(*sqltable.SQLTable).Where(args[1].(string), args[2].(string), args[3])
	p.Ret(4, ret0)
}

func execmSQLTableSet(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*sqltable.SQLTable).Set(args[1].(string), args[2])
	p.Ret(3, ret0)
}

func execmSQLTableSets(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	ret0 := args[0].(*sqltable.SQLTable).Sets(args[1:]...)
	p.Ret(arity, ret0)
}

func execmSQLTableExcluded(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	ret0 := args[0].(*sqltable.SQLTable).Excluded(gop.ToStrings(args[1:])...)
	p.Ret(arity, ret0)
}

func execmSQLTablePrepare(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*sqltable.SQLTable).Prepare(args[1].(string))
	p.Ret(2, ret0)
}

func execmSQLTableSQL(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*sqltable.SQLTable).SQL()
	p.Ret(1, ret0)
}

func execmSQLTableArgs(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	ret0 := args[0].(*sqltable.SQLTable).Args(args[1:]...)
	p.Ret(arity, ret0)
}

// I is a Go package instance.
var I = gop.NewGoPackage("github.com/456vv/x/sqltable")

func init() {
	I.RegisterFuncs(
		I.Func("NewSQLTable", sqltable.NewSQLTable, execNewSQLTable),
		I.Func("(*SQLTable).ColumnMark", (*sqltable.SQLTable).ColumnMark, execmSQLTableColumnMark),
		I.Func("(*SQLTable).Value", (*sqltable.SQLTable).Value, execmSQLTableValue),
		I.Func("(*SQLTable).And", (*sqltable.SQLTable).And, execmSQLTableAnd),
		I.Func("(*SQLTable).Or", (*sqltable.SQLTable).Or, execmSQLTableOr),
		I.Func("(*SQLTable).Where", (*sqltable.SQLTable).Where, execmSQLTableWhere),
		I.Func("(*SQLTable).Set", (*sqltable.SQLTable).Set, execmSQLTableSet),
		I.Func("(*SQLTable).Prepare", (*sqltable.SQLTable).Prepare, execmSQLTablePrepare),
		I.Func("(*SQLTable).SQL", (*sqltable.SQLTable).SQL, execmSQLTableSQL),
	)
	I.RegisterFuncvs(
		I.Funcv("(*SQLTable).Values", (*sqltable.SQLTable).Values, execmSQLTableValues),
		I.Funcv("(*SQLTable).Ands", (*sqltable.SQLTable).Ands, execmSQLTableAnds),
		I.Funcv("(*SQLTable).Ors", (*sqltable.SQLTable).Ors, execmSQLTableOrs),
		I.Funcv("(*SQLTable).Sets", (*sqltable.SQLTable).Sets, execmSQLTableSets),
		I.Funcv("(*SQLTable).Excluded", (*sqltable.SQLTable).Excluded, execmSQLTableExcluded),
		I.Funcv("(*SQLTable).Args", (*sqltable.SQLTable).Args, execmSQLTableArgs),
	)
	I.RegisterTypes(
		I.Type("SQLTable", reflect.TypeOf((*sqltable.SQLTable)(nil)).Elem()),
	)
}
