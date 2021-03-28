// Package sqlite3 provide Go+ "github.com/mattn/go-sqlite3" package, as "github.com/mattn/go-sqlite3" package in Go.
package sqlite3

import (
        context "context"
        driver "database/sql/driver"
        reflect "reflect"

        gop "github.com/goplus/gop"
        sqlite3 "github.com/mattn/go-sqlite3"
)

func execCryptEncoderSHA1(_ int, p *gop.Context) {
        args := p.GetArgs(2)
        ret0 := sqlite3.CryptEncoderSHA1(args[0].([]byte), args[1])
        p.Ret(2, ret0)
}

func execCryptEncoderSHA256(_ int, p *gop.Context) {
        args := p.GetArgs(2)
        ret0 := sqlite3.CryptEncoderSHA256(args[0].([]byte), args[1])
        p.Ret(2, ret0)
}

func execCryptEncoderSHA384(_ int, p *gop.Context) {
        args := p.GetArgs(2)
        ret0 := sqlite3.CryptEncoderSHA384(args[0].([]byte), args[1])
        p.Ret(2, ret0)
}

func execCryptEncoderSHA512(_ int, p *gop.Context) {
        args := p.GetArgs(2)
        ret0 := sqlite3.CryptEncoderSHA512(args[0].([]byte), args[1])
        p.Ret(2, ret0)
}

func execCryptEncoderSSHA1(_ int, p *gop.Context) {
        args := p.GetArgs(1)
        ret0 := sqlite3.CryptEncoderSSHA1(args[0].(string))
        p.Ret(1, ret0)
}

func execCryptEncoderSSHA256(_ int, p *gop.Context) {
        args := p.GetArgs(1)
        ret0 := sqlite3.CryptEncoderSSHA256(args[0].(string))
        p.Ret(1, ret0)
}

func execCryptEncoderSSHA384(_ int, p *gop.Context) {
        args := p.GetArgs(1)
        ret0 := sqlite3.CryptEncoderSSHA384(args[0].(string))
        p.Ret(1, ret0)
}

func execCryptEncoderSSHA512(_ int, p *gop.Context) {
        args := p.GetArgs(1)
        ret0 := sqlite3.CryptEncoderSSHA512(args[0].(string))
        p.Ret(1, ret0)
}

func execmErrNoError(_ int, p *gop.Context) {
        args := p.GetArgs(1)
        ret0 := args[0].(sqlite3.ErrNo).Error()
        p.Ret(1, ret0)
}

func execmErrNoExtend(_ int, p *gop.Context) {
        args := p.GetArgs(2)
        ret0 := args[0].(sqlite3.ErrNo).Extend(args[1].(int))
        p.Ret(2, ret0)
}

func execmErrNoExtendedError(_ int, p *gop.Context) {
        args := p.GetArgs(1)
        ret0 := args[0].(sqlite3.ErrNoExtended).Error()
        p.Ret(1, ret0)
}

func execmErrorError(_ int, p *gop.Context) {
        args := p.GetArgs(1)
        ret0 := args[0].(sqlite3.Error).Error()
        p.Ret(1, ret0)
}

func execmSQLiteBackupStep(_ int, p *gop.Context) {
        args := p.GetArgs(2)
        ret0, ret1 := args[0].(*sqlite3.SQLiteBackup).Step(args[1].(int))
        p.Ret(2, ret0, ret1)
}

func execmSQLiteBackupRemaining(_ int, p *gop.Context) {
        args := p.GetArgs(1)
        ret0 := args[0].(*sqlite3.SQLiteBackup).Remaining()
        p.Ret(1, ret0)
}

func execmSQLiteBackupPageCount(_ int, p *gop.Context) {
        args := p.GetArgs(1)
        ret0 := args[0].(*sqlite3.SQLiteBackup).PageCount()
        p.Ret(1, ret0)
}

func execmSQLiteBackupFinish(_ int, p *gop.Context) {
        args := p.GetArgs(1)
        ret0 := args[0].(*sqlite3.SQLiteBackup).Finish()
        p.Ret(1, ret0)
}

func execmSQLiteBackupClose(_ int, p *gop.Context) {
        args := p.GetArgs(1)
        ret0 := args[0].(*sqlite3.SQLiteBackup).Close()
        p.Ret(1, ret0)
}

func toType0(v interface{}) context.Context {
        if v == nil {
                return nil
        }
        return v.(context.Context)
}

func execmSQLiteConnPing(_ int, p *gop.Context) {
        args := p.GetArgs(2)
        ret0 := args[0].(*sqlite3.SQLiteConn).Ping(toType0(args[1]))
        p.Ret(2, ret0)
}

func execmSQLiteConnQueryContext(_ int, p *gop.Context) {
        args := p.GetArgs(4)
        ret0, ret1 := args[0].(*sqlite3.SQLiteConn).QueryContext(toType0(args[1]), args[2].(string), args[3].([]driver.NamedValue))
        p.Ret(4, ret0, ret1)
}

func execmSQLiteConnExecContext(_ int, p *gop.Context) {
        args := p.GetArgs(4)
        ret0, ret1 := args[0].(*sqlite3.SQLiteConn).ExecContext(toType0(args[1]), args[2].(string), args[3].([]driver.NamedValue))
        p.Ret(4, ret0, ret1)
}

func execmSQLiteConnPrepareContext(_ int, p *gop.Context) {
        args := p.GetArgs(3)
        ret0, ret1 := args[0].(*sqlite3.SQLiteConn).PrepareContext(toType0(args[1]), args[2].(string))
        p.Ret(3, ret0, ret1)
}

func execmSQLiteConnBeginTx(_ int, p *gop.Context) {
        args := p.GetArgs(3)
        ret0, ret1 := args[0].(*sqlite3.SQLiteConn).BeginTx(toType0(args[1]), args[2].(driver.TxOptions))
        p.Ret(3, ret0, ret1)
}

func execmSQLiteConnRegisterPreUpdateHook(_ int, p *gop.Context) {
        args := p.GetArgs(2)
        args[0].(*sqlite3.SQLiteConn).RegisterPreUpdateHook(args[1].(func(sqlite3.SQLitePreUpdateData)))
        p.PopN(2)
}

func execmSQLiteConnBackup(_ int, p *gop.Context) {
        args := p.GetArgs(4)
        ret0, ret1 := args[0].(*sqlite3.SQLiteConn).Backup(args[1].(string), args[2].(*sqlite3.SQLiteConn), args[3].(string))
        p.Ret(4, ret0, ret1)
}

func execmSQLiteConnRegisterCollation(_ int, p *gop.Context) {
        args := p.GetArgs(3)
        ret0 := args[0].(*sqlite3.SQLiteConn).RegisterCollation(args[1].(string), args[2].(func(string, string) int))
        p.Ret(3, ret0)
}

func execmSQLiteConnRegisterCommitHook(_ int, p *gop.Context) {
        args := p.GetArgs(2)
        args[0].(*sqlite3.SQLiteConn).RegisterCommitHook(args[1].(func() int))
        p.PopN(2)
}

func execmSQLiteConnRegisterRollbackHook(_ int, p *gop.Context) {
        args := p.GetArgs(2)
        args[0].(*sqlite3.SQLiteConn).RegisterRollbackHook(args[1].(func()))
        p.PopN(2)
}

func execmSQLiteConnRegisterUpdateHook(_ int, p *gop.Context) {
        args := p.GetArgs(2)
        args[0].(*sqlite3.SQLiteConn).RegisterUpdateHook(args[1].(func(int, string, string, int64)))
        p.PopN(2)
}

func execmSQLiteConnRegisterAuthorizer(_ int, p *gop.Context) {
        args := p.GetArgs(2)
        args[0].(*sqlite3.SQLiteConn).RegisterAuthorizer(args[1].(func(int, string, string, string) int))
        p.PopN(2)
}

func execmSQLiteConnRegisterFunc(_ int, p *gop.Context) {
        args := p.GetArgs(4)
        ret0 := args[0].(*sqlite3.SQLiteConn).RegisterFunc(args[1].(string), args[2], args[3].(bool))
        p.Ret(4, ret0)
}

func execmSQLiteConnRegisterAggregator(_ int, p *gop.Context) {
        args := p.GetArgs(4)
        ret0 := args[0].(*sqlite3.SQLiteConn).RegisterAggregator(args[1].(string), args[2], args[3].(bool))
        p.Ret(4, ret0)
}

func execmSQLiteConnAutoCommit(_ int, p *gop.Context) {
        args := p.GetArgs(1)
        ret0 := args[0].(*sqlite3.SQLiteConn).AutoCommit()
        p.Ret(1, ret0)
}

func execmSQLiteConnExec(_ int, p *gop.Context) {
        args := p.GetArgs(3)
        ret0, ret1 := args[0].(*sqlite3.SQLiteConn).Exec(args[1].(string), args[2].([]driver.Value))
        p.Ret(3, ret0, ret1)
}

func execmSQLiteConnQuery(_ int, p *gop.Context) {
        args := p.GetArgs(3)
        ret0, ret1 := args[0].(*sqlite3.SQLiteConn).Query(args[1].(string), args[2].([]driver.Value))
        p.Ret(3, ret0, ret1)
}

func execmSQLiteConnBegin(_ int, p *gop.Context) {
        args := p.GetArgs(1)
        ret0, ret1 := args[0].(*sqlite3.SQLiteConn).Begin()
        p.Ret(1, ret0, ret1)
}

func execmSQLiteConnClose(_ int, p *gop.Context) {
        args := p.GetArgs(1)
        ret0 := args[0].(*sqlite3.SQLiteConn).Close()
        p.Ret(1, ret0)
}

func execmSQLiteConnPrepare(_ int, p *gop.Context) {
        args := p.GetArgs(2)
        ret0, ret1 := args[0].(*sqlite3.SQLiteConn).Prepare(args[1].(string))
        p.Ret(2, ret0, ret1)
}

func execmSQLiteConnGetFilename(_ int, p *gop.Context) {
        args := p.GetArgs(2)
        ret0 := args[0].(*sqlite3.SQLiteConn).GetFilename(args[1].(string))
        p.Ret(2, ret0)
}

func execmSQLiteConnGetLimit(_ int, p *gop.Context) {
        args := p.GetArgs(2)
        ret0 := args[0].(*sqlite3.SQLiteConn).GetLimit(args[1].(int))
        p.Ret(2, ret0)
}

func execmSQLiteConnSetLimit(_ int, p *gop.Context) {
        args := p.GetArgs(3)
        ret0 := args[0].(*sqlite3.SQLiteConn).SetLimit(args[1].(int), args[2].(int))
        p.Ret(3, ret0)
}

func execmSQLiteConnLoadExtension(_ int, p *gop.Context) {
        args := p.GetArgs(3)
        ret0 := args[0].(*sqlite3.SQLiteConn).LoadExtension(args[1].(string), args[2].(string))
        p.Ret(3, ret0)
}

func execmSQLiteConnAuthenticate(_ int, p *gop.Context) {
        args := p.GetArgs(3)
        ret0 := args[0].(*sqlite3.SQLiteConn).Authenticate(args[1].(string), args[2].(string))
        p.Ret(3, ret0)
}

func execmSQLiteConnAuthUserAdd(_ int, p *gop.Context) {
        args := p.GetArgs(4)
        ret0 := args[0].(*sqlite3.SQLiteConn).AuthUserAdd(args[1].(string), args[2].(string), args[3].(bool))
        p.Ret(4, ret0)
}

func execmSQLiteConnAuthUserChange(_ int, p *gop.Context) {
        args := p.GetArgs(4)
        ret0 := args[0].(*sqlite3.SQLiteConn).AuthUserChange(args[1].(string), args[2].(string), args[3].(bool))
        p.Ret(4, ret0)
}

func execmSQLiteConnAuthUserDelete(_ int, p *gop.Context) {
        args := p.GetArgs(2)
        ret0 := args[0].(*sqlite3.SQLiteConn).AuthUserDelete(args[1].(string))
        p.Ret(2, ret0)
}

func execmSQLiteConnAuthEnabled(_ int, p *gop.Context) {
        args := p.GetArgs(1)
        ret0 := args[0].(*sqlite3.SQLiteConn).AuthEnabled()
        p.Ret(1, ret0)
}

func execmSQLiteContextResultBool(_ int, p *gop.Context) {
        args := p.GetArgs(2)
        args[0].(*sqlite3.SQLiteContext).ResultBool(args[1].(bool))
        p.PopN(2)
}

func execmSQLiteContextResultBlob(_ int, p *gop.Context) {
        args := p.GetArgs(2)
        args[0].(*sqlite3.SQLiteContext).ResultBlob(args[1].([]byte))
        p.PopN(2)
}

func execmSQLiteContextResultDouble(_ int, p *gop.Context) {
        args := p.GetArgs(2)
        args[0].(*sqlite3.SQLiteContext).ResultDouble(args[1].(float64))
        p.PopN(2)
}

func execmSQLiteContextResultInt(_ int, p *gop.Context) {
        args := p.GetArgs(2)
        args[0].(*sqlite3.SQLiteContext).ResultInt(args[1].(int))
        p.PopN(2)
}

func execmSQLiteContextResultInt64(_ int, p *gop.Context) {
        args := p.GetArgs(2)
        args[0].(*sqlite3.SQLiteContext).ResultInt64(args[1].(int64))
        p.PopN(2)
}

func execmSQLiteContextResultNull(_ int, p *gop.Context) {
        args := p.GetArgs(1)
        args[0].(*sqlite3.SQLiteContext).ResultNull()
        p.PopN(1)
}

func execmSQLiteContextResultText(_ int, p *gop.Context) {
        args := p.GetArgs(2)
        args[0].(*sqlite3.SQLiteContext).ResultText(args[1].(string))
        p.PopN(2)
}

func execmSQLiteContextResultZeroblob(_ int, p *gop.Context) {
        args := p.GetArgs(2)
        args[0].(*sqlite3.SQLiteContext).ResultZeroblob(args[1].(int))
        p.PopN(2)
}

func execmSQLiteDriverOpen(_ int, p *gop.Context) {
        args := p.GetArgs(2)
        ret0, ret1 := args[0].(*sqlite3.SQLiteDriver).Open(args[1].(string))
        p.Ret(2, ret0, ret1)
}

func execmSQLiteResultLastInsertId(_ int, p *gop.Context) {
        args := p.GetArgs(1)
        ret0, ret1 := args[0].(*sqlite3.SQLiteResult).LastInsertId()
        p.Ret(1, ret0, ret1)
}

func execmSQLiteResultRowsAffected(_ int, p *gop.Context) {
        args := p.GetArgs(1)
        ret0, ret1 := args[0].(*sqlite3.SQLiteResult).RowsAffected()
        p.Ret(1, ret0, ret1)
}

func execmSQLiteRowsClose(_ int, p *gop.Context) {
        args := p.GetArgs(1)
        ret0 := args[0].(*sqlite3.SQLiteRows).Close()
        p.Ret(1, ret0)
}

func execmSQLiteRowsColumns(_ int, p *gop.Context) {
        args := p.GetArgs(1)
        ret0 := args[0].(*sqlite3.SQLiteRows).Columns()
        p.Ret(1, ret0)
}

func execmSQLiteRowsDeclTypes(_ int, p *gop.Context) {
        args := p.GetArgs(1)
        ret0 := args[0].(*sqlite3.SQLiteRows).DeclTypes()
        p.Ret(1, ret0)
}

func execmSQLiteRowsNext(_ int, p *gop.Context) {
        args := p.GetArgs(2)
        ret0 := args[0].(*sqlite3.SQLiteRows).Next(args[1].([]driver.Value))
        p.Ret(2, ret0)
}

func execmSQLiteRowsColumnTypeDatabaseTypeName(_ int, p *gop.Context) {
        args := p.GetArgs(2)
        ret0 := args[0].(*sqlite3.SQLiteRows).ColumnTypeDatabaseTypeName(args[1].(int))
        p.Ret(2, ret0)
}

func execmSQLiteRowsColumnTypeScanType(_ int, p *gop.Context) {
        args := p.GetArgs(2)
        ret0 := args[0].(*sqlite3.SQLiteRows).ColumnTypeScanType(args[1].(int))
        p.Ret(2, ret0)
}

func execmSQLiteStmtQueryContext(_ int, p *gop.Context) {
        args := p.GetArgs(3)
        ret0, ret1 := args[0].(*sqlite3.SQLiteStmt).QueryContext(toType0(args[1]), args[2].([]driver.NamedValue))
        p.Ret(3, ret0, ret1)
}

func execmSQLiteStmtExecContext(_ int, p *gop.Context) {
        args := p.GetArgs(3)
        ret0, ret1 := args[0].(*sqlite3.SQLiteStmt).ExecContext(toType0(args[1]), args[2].([]driver.NamedValue))
        p.Ret(3, ret0, ret1)
}

func execmSQLiteStmtClose(_ int, p *gop.Context) {
        args := p.GetArgs(1)
        ret0 := args[0].(*sqlite3.SQLiteStmt).Close()
        p.Ret(1, ret0)
}

func execmSQLiteStmtNumInput(_ int, p *gop.Context) {
        args := p.GetArgs(1)
        ret0 := args[0].(*sqlite3.SQLiteStmt).NumInput()
        p.Ret(1, ret0)
}

func execmSQLiteStmtQuery(_ int, p *gop.Context) {
        args := p.GetArgs(2)
        ret0, ret1 := args[0].(*sqlite3.SQLiteStmt).Query(args[1].([]driver.Value))
        p.Ret(2, ret0, ret1)
}

func execmSQLiteStmtExec(_ int, p *gop.Context) {
        args := p.GetArgs(2)
        ret0, ret1 := args[0].(*sqlite3.SQLiteStmt).Exec(args[1].([]driver.Value))
        p.Ret(2, ret0, ret1)
}

func execmSQLiteTxCommit(_ int, p *gop.Context) {
        args := p.GetArgs(1)
        ret0 := args[0].(*sqlite3.SQLiteTx).Commit()
        p.Ret(1, ret0)
}

func execmSQLiteTxRollback(_ int, p *gop.Context) {
        args := p.GetArgs(1)
        ret0 := args[0].(*sqlite3.SQLiteTx).Rollback()
        p.Ret(1, ret0)
}

func execVersion(_ int, p *gop.Context) {
        ret0, ret1, ret2 := sqlite3.Version()
        p.Ret(0, ret0, ret1, ret2)
}

// I is a Go package instance.
var I = gop.NewGoPackage("github.com/mattn/go-sqlite3")

func init() {
        I.RegisterFuncs(
                I.Func("CryptEncoderSHA1", sqlite3.CryptEncoderSHA1, execCryptEncoderSHA1),
                I.Func("CryptEncoderSHA256", sqlite3.CryptEncoderSHA256, execCryptEncoderSHA256),
                I.Func("CryptEncoderSHA384", sqlite3.CryptEncoderSHA384, execCryptEncoderSHA384),
                I.Func("CryptEncoderSHA512", sqlite3.CryptEncoderSHA512, execCryptEncoderSHA512),
                I.Func("CryptEncoderSSHA1", sqlite3.CryptEncoderSSHA1, execCryptEncoderSSHA1),
                I.Func("CryptEncoderSSHA256", sqlite3.CryptEncoderSSHA256, execCryptEncoderSSHA256),
                I.Func("CryptEncoderSSHA384", sqlite3.CryptEncoderSSHA384, execCryptEncoderSSHA384),
                I.Func("CryptEncoderSSHA512", sqlite3.CryptEncoderSSHA512, execCryptEncoderSSHA512),
                I.Func("(ErrNo).Error", (sqlite3.ErrNo).Error, execmErrNoError),
                I.Func("(ErrNo).Extend", (sqlite3.ErrNo).Extend, execmErrNoExtend),
                I.Func("(ErrNoExtended).Error", (sqlite3.ErrNoExtended).Error, execmErrNoExtendedError),
                I.Func("(Error).Error", (sqlite3.Error).Error, execmErrorError),
                I.Func("(*SQLiteBackup).Step", (*sqlite3.SQLiteBackup).Step, execmSQLiteBackupStep),
                I.Func("(*SQLiteBackup).Remaining", (*sqlite3.SQLiteBackup).Remaining, execmSQLiteBackupRemaining),
                I.Func("(*SQLiteBackup).PageCount", (*sqlite3.SQLiteBackup).PageCount, execmSQLiteBackupPageCount),
                I.Func("(*SQLiteBackup).Finish", (*sqlite3.SQLiteBackup).Finish, execmSQLiteBackupFinish),
                I.Func("(*SQLiteBackup).Close", (*sqlite3.SQLiteBackup).Close, execmSQLiteBackupClose),
                I.Func("(*SQLiteConn).Ping", (*sqlite3.SQLiteConn).Ping, execmSQLiteConnPing),
                I.Func("(*SQLiteConn).QueryContext", (*sqlite3.SQLiteConn).QueryContext, execmSQLiteConnQueryContext),
                I.Func("(*SQLiteConn).ExecContext", (*sqlite3.SQLiteConn).ExecContext, execmSQLiteConnExecContext),
                I.Func("(*SQLiteConn).PrepareContext", (*sqlite3.SQLiteConn).PrepareContext, execmSQLiteConnPrepareContext),
                I.Func("(*SQLiteConn).BeginTx", (*sqlite3.SQLiteConn).BeginTx, execmSQLiteConnBeginTx),
                I.Func("(*SQLiteConn).RegisterPreUpdateHook", (*sqlite3.SQLiteConn).RegisterPreUpdateHook, execmSQLiteConnRegisterPreUpdateHook),
                I.Func("(*SQLiteConn).Backup", (*sqlite3.SQLiteConn).Backup, execmSQLiteConnBackup),
                I.Func("(*SQLiteConn).RegisterCollation", (*sqlite3.SQLiteConn).RegisterCollation, execmSQLiteConnRegisterCollation),
                I.Func("(*SQLiteConn).RegisterCommitHook", (*sqlite3.SQLiteConn).RegisterCommitHook, execmSQLiteConnRegisterCommitHook),
                I.Func("(*SQLiteConn).RegisterRollbackHook", (*sqlite3.SQLiteConn).RegisterRollbackHook, execmSQLiteConnRegisterRollbackHook),
                I.Func("(*SQLiteConn).RegisterUpdateHook", (*sqlite3.SQLiteConn).RegisterUpdateHook, execmSQLiteConnRegisterUpdateHook),
                I.Func("(*SQLiteConn).RegisterAuthorizer", (*sqlite3.SQLiteConn).RegisterAuthorizer, execmSQLiteConnRegisterAuthorizer),
                I.Func("(*SQLiteConn).RegisterFunc", (*sqlite3.SQLiteConn).RegisterFunc, execmSQLiteConnRegisterFunc),
                I.Func("(*SQLiteConn).RegisterAggregator", (*sqlite3.SQLiteConn).RegisterAggregator, execmSQLiteConnRegisterAggregator),
                I.Func("(*SQLiteConn).AutoCommit", (*sqlite3.SQLiteConn).AutoCommit, execmSQLiteConnAutoCommit),
                I.Func("(*SQLiteConn).Exec", (*sqlite3.SQLiteConn).Exec, execmSQLiteConnExec),
                I.Func("(*SQLiteConn).Query", (*sqlite3.SQLiteConn).Query, execmSQLiteConnQuery),
                I.Func("(*SQLiteConn).Begin", (*sqlite3.SQLiteConn).Begin, execmSQLiteConnBegin),
                I.Func("(*SQLiteConn).Close", (*sqlite3.SQLiteConn).Close, execmSQLiteConnClose),
                I.Func("(*SQLiteConn).Prepare", (*sqlite3.SQLiteConn).Prepare, execmSQLiteConnPrepare),
                I.Func("(*SQLiteConn).GetFilename", (*sqlite3.SQLiteConn).GetFilename, execmSQLiteConnGetFilename),
                I.Func("(*SQLiteConn).GetLimit", (*sqlite3.SQLiteConn).GetLimit, execmSQLiteConnGetLimit),
                I.Func("(*SQLiteConn).SetLimit", (*sqlite3.SQLiteConn).SetLimit, execmSQLiteConnSetLimit),
                I.Func("(*SQLiteConn).LoadExtension", (*sqlite3.SQLiteConn).LoadExtension, execmSQLiteConnLoadExtension),
                I.Func("(*SQLiteConn).Authenticate", (*sqlite3.SQLiteConn).Authenticate, execmSQLiteConnAuthenticate),
                I.Func("(*SQLiteConn).AuthUserAdd", (*sqlite3.SQLiteConn).AuthUserAdd, execmSQLiteConnAuthUserAdd),
                I.Func("(*SQLiteConn).AuthUserChange", (*sqlite3.SQLiteConn).AuthUserChange, execmSQLiteConnAuthUserChange),
                I.Func("(*SQLiteConn).AuthUserDelete", (*sqlite3.SQLiteConn).AuthUserDelete, execmSQLiteConnAuthUserDelete),
                I.Func("(*SQLiteConn).AuthEnabled", (*sqlite3.SQLiteConn).AuthEnabled, execmSQLiteConnAuthEnabled),
                I.Func("(*SQLiteContext).ResultBool", (*sqlite3.SQLiteContext).ResultBool, execmSQLiteContextResultBool),
                I.Func("(*SQLiteContext).ResultBlob", (*sqlite3.SQLiteContext).ResultBlob, execmSQLiteContextResultBlob),
                I.Func("(*SQLiteContext).ResultDouble", (*sqlite3.SQLiteContext).ResultDouble, execmSQLiteContextResultDouble),
                I.Func("(*SQLiteContext).ResultInt", (*sqlite3.SQLiteContext).ResultInt, execmSQLiteContextResultInt),
                I.Func("(*SQLiteContext).ResultInt64", (*sqlite3.SQLiteContext).ResultInt64, execmSQLiteContextResultInt64),
                I.Func("(*SQLiteContext).ResultNull", (*sqlite3.SQLiteContext).ResultNull, execmSQLiteContextResultNull),
                I.Func("(*SQLiteContext).ResultText", (*sqlite3.SQLiteContext).ResultText, execmSQLiteContextResultText),
                I.Func("(*SQLiteContext).ResultZeroblob", (*sqlite3.SQLiteContext).ResultZeroblob, execmSQLiteContextResultZeroblob),
                I.Func("(*SQLiteDriver).Open", (*sqlite3.SQLiteDriver).Open, execmSQLiteDriverOpen),
                I.Func("(*SQLiteResult).LastInsertId", (*sqlite3.SQLiteResult).LastInsertId, execmSQLiteResultLastInsertId),
                I.Func("(*SQLiteResult).RowsAffected", (*sqlite3.SQLiteResult).RowsAffected, execmSQLiteResultRowsAffected),
                I.Func("(*SQLiteRows).Close", (*sqlite3.SQLiteRows).Close, execmSQLiteRowsClose),
                I.Func("(*SQLiteRows).Columns", (*sqlite3.SQLiteRows).Columns, execmSQLiteRowsColumns),
                I.Func("(*SQLiteRows).DeclTypes", (*sqlite3.SQLiteRows).DeclTypes, execmSQLiteRowsDeclTypes),
                I.Func("(*SQLiteRows).Next", (*sqlite3.SQLiteRows).Next, execmSQLiteRowsNext),
                I.Func("(*SQLiteRows).ColumnTypeDatabaseTypeName", (*sqlite3.SQLiteRows).ColumnTypeDatabaseTypeName, execmSQLiteRowsColumnTypeDatabaseTypeName),
                I.Func("(*SQLiteRows).ColumnTypeScanType", (*sqlite3.SQLiteRows).ColumnTypeScanType, execmSQLiteRowsColumnTypeScanType),
                I.Func("(*SQLiteStmt).QueryContext", (*sqlite3.SQLiteStmt).QueryContext, execmSQLiteStmtQueryContext),
                I.Func("(*SQLiteStmt).ExecContext", (*sqlite3.SQLiteStmt).ExecContext, execmSQLiteStmtExecContext),
                I.Func("(*SQLiteStmt).Close", (*sqlite3.SQLiteStmt).Close, execmSQLiteStmtClose),
                I.Func("(*SQLiteStmt).NumInput", (*sqlite3.SQLiteStmt).NumInput, execmSQLiteStmtNumInput),
                I.Func("(*SQLiteStmt).Query", (*sqlite3.SQLiteStmt).Query, execmSQLiteStmtQuery),
                I.Func("(*SQLiteStmt).Exec", (*sqlite3.SQLiteStmt).Exec, execmSQLiteStmtExec),
                I.Func("(*SQLiteTx).Commit", (*sqlite3.SQLiteTx).Commit, execmSQLiteTxCommit),
                I.Func("(*SQLiteTx).Rollback", (*sqlite3.SQLiteTx).Rollback, execmSQLiteTxRollback),
                I.Func("Version", sqlite3.Version, execVersion),
        )
        I.RegisterVars(
                I.Var("ErrAbort", &sqlite3.ErrAbort),
                I.Var("ErrAbortRollback", &sqlite3.ErrAbortRollback),
                I.Var("ErrAuth", &sqlite3.ErrAuth),
                I.Var("ErrBusy", &sqlite3.ErrBusy),
                I.Var("ErrBusyRecovery", &sqlite3.ErrBusyRecovery),
                I.Var("ErrBusySnapshot", &sqlite3.ErrBusySnapshot),
                I.Var("ErrCantOpen", &sqlite3.ErrCantOpen),
                I.Var("ErrCantOpenConvPath", &sqlite3.ErrCantOpenConvPath),
                I.Var("ErrCantOpenFullPath", &sqlite3.ErrCantOpenFullPath),
                I.Var("ErrCantOpenIsDir", &sqlite3.ErrCantOpenIsDir),
                I.Var("ErrCantOpenNoTempDir", &sqlite3.ErrCantOpenNoTempDir),
                I.Var("ErrConstraint", &sqlite3.ErrConstraint),
                I.Var("ErrConstraintCheck", &sqlite3.ErrConstraintCheck),
                I.Var("ErrConstraintCommitHook", &sqlite3.ErrConstraintCommitHook),
                I.Var("ErrConstraintForeignKey", &sqlite3.ErrConstraintForeignKey),
                I.Var("ErrConstraintFunction", &sqlite3.ErrConstraintFunction),
                I.Var("ErrConstraintNotNull", &sqlite3.ErrConstraintNotNull),
                I.Var("ErrConstraintPrimaryKey", &sqlite3.ErrConstraintPrimaryKey),
                I.Var("ErrConstraintRowID", &sqlite3.ErrConstraintRowID),
                I.Var("ErrConstraintTrigger", &sqlite3.ErrConstraintTrigger),
                I.Var("ErrConstraintUnique", &sqlite3.ErrConstraintUnique),
                I.Var("ErrConstraintVTab", &sqlite3.ErrConstraintVTab),
                I.Var("ErrCorrupt", &sqlite3.ErrCorrupt),
                I.Var("ErrCorruptVTab", &sqlite3.ErrCorruptVTab),
                I.Var("ErrEmpty", &sqlite3.ErrEmpty),
                I.Var("ErrError", &sqlite3.ErrError),
                I.Var("ErrFormat", &sqlite3.ErrFormat),
                I.Var("ErrFull", &sqlite3.ErrFull),
                I.Var("ErrInternal", &sqlite3.ErrInternal),
                I.Var("ErrInterrupt", &sqlite3.ErrInterrupt),
                I.Var("ErrIoErr", &sqlite3.ErrIoErr),
                I.Var("ErrIoErrAccess", &sqlite3.ErrIoErrAccess),
                I.Var("ErrIoErrBlocked", &sqlite3.ErrIoErrBlocked),
                I.Var("ErrIoErrCheckReservedLock", &sqlite3.ErrIoErrCheckReservedLock),
                I.Var("ErrIoErrClose", &sqlite3.ErrIoErrClose),
                I.Var("ErrIoErrConvPath", &sqlite3.ErrIoErrConvPath),
                I.Var("ErrIoErrDelete", &sqlite3.ErrIoErrDelete),
                I.Var("ErrIoErrDeleteNoent", &sqlite3.ErrIoErrDeleteNoent),
                I.Var("ErrIoErrDirClose", &sqlite3.ErrIoErrDirClose),
                I.Var("ErrIoErrDirFsync", &sqlite3.ErrIoErrDirFsync),
                I.Var("ErrIoErrFstat", &sqlite3.ErrIoErrFstat),
                I.Var("ErrIoErrFsync", &sqlite3.ErrIoErrFsync),
                I.Var("ErrIoErrGetTempPath", &sqlite3.ErrIoErrGetTempPath),
                I.Var("ErrIoErrLock", &sqlite3.ErrIoErrLock),
                I.Var("ErrIoErrMMap", &sqlite3.ErrIoErrMMap),
                I.Var("ErrIoErrNoMem", &sqlite3.ErrIoErrNoMem),
                I.Var("ErrIoErrRDlock", &sqlite3.ErrIoErrRDlock),
                I.Var("ErrIoErrRead", &sqlite3.ErrIoErrRead),
                I.Var("ErrIoErrSHMLock", &sqlite3.ErrIoErrSHMLock),
                I.Var("ErrIoErrSHMMap", &sqlite3.ErrIoErrSHMMap),
                I.Var("ErrIoErrSHMOpen", &sqlite3.ErrIoErrSHMOpen),
                I.Var("ErrIoErrSHMSize", &sqlite3.ErrIoErrSHMSize),
                I.Var("ErrIoErrSeek", &sqlite3.ErrIoErrSeek),
                I.Var("ErrIoErrShortRead", &sqlite3.ErrIoErrShortRead),
                I.Var("ErrIoErrTruncate", &sqlite3.ErrIoErrTruncate),
                I.Var("ErrIoErrUnlock", &sqlite3.ErrIoErrUnlock),
                I.Var("ErrIoErrWrite", &sqlite3.ErrIoErrWrite),
                I.Var("ErrLocked", &sqlite3.ErrLocked),
                I.Var("ErrLockedSharedCache", &sqlite3.ErrLockedSharedCache),
                I.Var("ErrMismatch", &sqlite3.ErrMismatch),
                I.Var("ErrMisuse", &sqlite3.ErrMisuse),
                I.Var("ErrNoLFS", &sqlite3.ErrNoLFS),
                I.Var("ErrNomem", &sqlite3.ErrNomem),
                I.Var("ErrNotADB", &sqlite3.ErrNotADB),
                I.Var("ErrNotFound", &sqlite3.ErrNotFound),
                I.Var("ErrNotice", &sqlite3.ErrNotice),
                I.Var("ErrNoticeRecoverRollback", &sqlite3.ErrNoticeRecoverRollback),
                I.Var("ErrNoticeRecoverWAL", &sqlite3.ErrNoticeRecoverWAL),
                I.Var("ErrPerm", &sqlite3.ErrPerm),
                I.Var("ErrProtocol", &sqlite3.ErrProtocol),
                I.Var("ErrRange", &sqlite3.ErrRange),
                I.Var("ErrReadonly", &sqlite3.ErrReadonly),
                I.Var("ErrReadonlyCantLock", &sqlite3.ErrReadonlyCantLock),
                I.Var("ErrReadonlyDbMoved", &sqlite3.ErrReadonlyDbMoved),
                I.Var("ErrReadonlyRecovery", &sqlite3.ErrReadonlyRecovery),
                I.Var("ErrReadonlyRollback", &sqlite3.ErrReadonlyRollback),
                I.Var("ErrSchema", &sqlite3.ErrSchema),
                I.Var("ErrTooBig", &sqlite3.ErrTooBig),
                I.Var("ErrWarning", &sqlite3.ErrWarning),
                I.Var("ErrWarningAutoIndex", &sqlite3.ErrWarningAutoIndex),
                I.Var("SQLiteTimestampFormats", &sqlite3.SQLiteTimestampFormats),
        )
        I.RegisterTypes(
                I.Type("ErrNo", reflect.TypeOf((*sqlite3.ErrNo)(nil)).Elem()),
                I.Type("ErrNoExtended", reflect.TypeOf((*sqlite3.ErrNoExtended)(nil)).Elem()),
                I.Type("Error", reflect.TypeOf((*sqlite3.Error)(nil)).Elem()),
                I.Type("SQLiteBackup", reflect.TypeOf((*sqlite3.SQLiteBackup)(nil)).Elem()),
                I.Type("SQLiteConn", reflect.TypeOf((*sqlite3.SQLiteConn)(nil)).Elem()),
                I.Type("SQLiteContext", reflect.TypeOf((*sqlite3.SQLiteContext)(nil)).Elem()),
                I.Type("SQLiteDriver", reflect.TypeOf((*sqlite3.SQLiteDriver)(nil)).Elem()),
                I.Type("SQLitePreUpdateData", reflect.TypeOf((*sqlite3.SQLitePreUpdateData)(nil)).Elem()),
                I.Type("SQLiteResult", reflect.TypeOf((*sqlite3.SQLiteResult)(nil)).Elem()),
                I.Type("SQLiteRows", reflect.TypeOf((*sqlite3.SQLiteRows)(nil)).Elem()),
                I.Type("SQLiteStmt", reflect.TypeOf((*sqlite3.SQLiteStmt)(nil)).Elem()),
                I.Type("SQLiteTx", reflect.TypeOf((*sqlite3.SQLiteTx)(nil)).Elem()),
        )
}