// Package db provide Go+ "github.com/456vv/x/db" package, as "github.com/456vv/x/db" package in Go.
package db

import (
	context "context"
	sql "database/sql"
	reflect "reflect"

	db "github.com/456vv/x/db"
	gop "github.com/goplus/gop"
)

func execmDBOpen(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0, ret1 := args[0].(*db.DB).Open(args[1].(string), args[2].(string))
	p.Ret(3, ret0, ret1)
}

func execmDBClose(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*db.DB).Close()
	p.Ret(1, ret0)
}

func execmDBPexec(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	ret0 := args[0].(*db.DB).Pexec(args[1].(*sql.Tx), args[2].(string), args[3:]...)
	p.Ret(arity, ret0)
}

func toType0(v interface{}) context.Context {
	if v == nil {
		return nil
	}
	return v.(context.Context)
}

func execmDBPexecContext(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	ret0 := args[0].(*db.DB).PexecContext(args[1].(*sql.Tx), toType0(args[2]), args[3].(string), args[4:]...)
	p.Ret(arity, ret0)
}

func execmDBExec(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	ret0, ret1 := args[0].(*db.DB).Exec(args[1].(*sql.Tx), args[2].(string), args[3:]...)
	p.Ret(arity, ret0, ret1)
}

func execmDBExecContext(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	ret0, ret1 := args[0].(*db.DB).ExecContext(args[1].(*sql.Tx), toType0(args[2]), args[3].(string), args[4:]...)
	p.Ret(arity, ret0, ret1)
}

func execmDBQueryRow(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	ret0 := args[0].(*db.DB).QueryRow(args[1].(*sql.Tx), args[2].(string), args[3:]...)
	p.Ret(arity, ret0)
}

func execmDBQueryRowContext(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	ret0 := args[0].(*db.DB).QueryRowContext(args[1].(*sql.Tx), toType0(args[2]), args[3].(string), args[4:]...)
	p.Ret(arity, ret0)
}

func execmDBQuery(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	ret0, ret1 := args[0].(*db.DB).Query(args[1].(*sql.Tx), args[2].(string), args[3:]...)
	p.Ret(arity, ret0, ret1)
}

func execmDBQueryContext(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	ret0, ret1 := args[0].(*db.DB).QueryContext(args[1].(*sql.Tx), toType0(args[2]), args[3].(string), args[4:]...)
	p.Ret(arity, ret0, ret1)
}

func execmDBQueryReader(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	ret0, ret1 := args[0].(*db.DB).QueryReader(args[1].(*sql.Tx), args[2].(string), args[3].(string), args[4:]...)
	p.Ret(arity, ret0, ret1)
}

func execmDBQueryReaderContext(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	ret0, ret1 := args[0].(*db.DB).QueryReaderContext(args[1].(*sql.Tx), toType0(args[2]), args[3].(string), args[4].(string), args[5:]...)
	p.Ret(arity, ret0, ret1)
}

func execmDBHas(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	ret0 := args[0].(*db.DB).Has(args[1].(string), args[2:]...)
	p.Ret(arity, ret0)
}

func execmDBHasContext(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	ret0 := args[0].(*db.DB).HasContext(toType0(args[1]), args[2].(string), args[3:]...)
	p.Ret(arity, ret0)
}

// I is a Go package instance.
var I = gop.NewGoPackage("github.com/456vv/x/db")

func init() {
	I.RegisterFuncs(
		I.Func("(*DB).Open", (*db.DB).Open, execmDBOpen),
		I.Func("(*DB).Close", (*db.DB).Close, execmDBClose),
	)
	I.RegisterFuncvs(
		I.Funcv("(*DB).Pexec", (*db.DB).Pexec, execmDBPexec),
		I.Funcv("(*DB).PexecContext", (*db.DB).PexecContext, execmDBPexecContext),
		I.Funcv("(*DB).Exec", (*db.DB).Exec, execmDBExec),
		I.Funcv("(*DB).ExecContext", (*db.DB).ExecContext, execmDBExecContext),
		I.Funcv("(*DB).QueryRow", (*db.DB).QueryRow, execmDBQueryRow),
		I.Funcv("(*DB).QueryRowContext", (*db.DB).QueryRowContext, execmDBQueryRowContext),
		I.Funcv("(*DB).Query", (*db.DB).Query, execmDBQuery),
		I.Funcv("(*DB).QueryContext", (*db.DB).QueryContext, execmDBQueryContext),
		I.Funcv("(*DB).QueryReader", (*db.DB).QueryReader, execmDBQueryReader),
		I.Funcv("(*DB).QueryReaderContext", (*db.DB).QueryReaderContext, execmDBQueryReaderContext),
		I.Funcv("(*DB).Has", (*db.DB).Has, execmDBHas),
		I.Funcv("(*DB).HasContext", (*db.DB).HasContext, execmDBHasContext),
	)
	I.RegisterVars(
		I.Var("ErrNoRows", &db.ErrNoRows),
		I.Var("ErrRows", &db.ErrRows),
	)
	I.RegisterTypes(
		I.Type("DB", reflect.TypeOf((*db.DB)(nil)).Elem()),
	)
}
