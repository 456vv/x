package	db
import(
	"github.com/456vv/vweb/v2"
	"database/sql"
	_ "github.com/lib/pq"
	"context"
	"strings"
	"reflect"
	"time"
	"errors"
	"log"
	"fmt"
)
/*
	DB.Open(driverName,	dataSourceName string) (*DB, error)
	DB.Close() error
	DB.Pexec(tx *sql.Tx, sqlstr string, args ...interface{}) error
	DB.PexecContext(tx	*sql.Tx, ctx context.Context, sqlstr string, args ...interface{}) (err error)
	DB.Exec(tx *sql.Tx,	sqlstr string, args	...interface{})	(result	sql.Result,	err	error)
	DB.ExecContext(tx *sql.Tx, ctx context.Context,	sqlstr string, args	...interface{})	(result	sql.Result,	err	error)
	DB.QueryRow(tx *sql.Tx,	sqlstr string, args	...interface{})	sqlRower
	DB.QueryRowContext(tx *sql.Tx, ctx context.Context,	sqlstr string, args	...interface{})	sqlRower
	DB.Query(tx	*sql.Tx, sqlstr	string,	args ...interface{}) (data []map[string]interface{}, err error)
	DB.QueryContext(tx *sql.Tx,	ctx	context.Context, sqlstr	string,	args ...interface{}) (data []map[string]interface{}, err error)
	DB.QueryReader(tx *sql.Tx, id string, sqlstr string, args ...interface{}) (line	map[interface{}]*vbody.Reader, err error)
	DB.QueryReaderContext(tx *sql.Tx, ctx context.Context, id string, sqlstr string, args ...interface{}) (line	map[interface{}]*vbody.Reader, err error)
	DB.Has(sqlstr string, args ...interface{}) error
	DB.HasDataContext(ctx context.Context, sqlstr string, args ...interface{}) error
*/

//支持于HasContext支持
var	ErrRows	= errors.New("sql: have	rows in	result set")
//支持于 HasContext，QueryReaderContext，QueryRowContext，PexecContext
var ErrNoRows = sql.ErrNoRows
	
var errDebugResult = errors.New("The debug result type does not match the returned result type")

type Rower interface{
	Scan(dest ...interface{}) error
	Err() error
}
type dbRow struct{
	tx	*sql.Tx
	r	*sql.Row
	err	error
	cancel context.CancelFunc
}
func (T	*dbRow)	Scan(dest ...interface{}) error	{
	if T.cancel != nil {
		defer T.cancel()
	}
	if T.err !=	nil	{
		return T.err
	}
	err	:= T.r.Scan(dest...)
	if T.tx	== nil {
		return err
	}
	if err != nil {
		T.tx.Rollback()
		return err
	}
	return T.tx.Commit()
}
func (T	*dbRow)	Err() error{
	return T.err
}

//说明：https://github.com/lib/pq/blob/master/doc.go
//postgresql://[user[:password]@][netloc][:port][,...][/dbname][?param1=value1&...]
type DB	struct{
	*sql.DB
	Timeout	time.Duration
	FormatColumnName bool			//将例名 AbcDf	=> abcDf
	TypeTo		func(ct *sql.ColumnType) reflect.Type
	ErrorLog	*log.Logger
	DebugPrint	bool
	DebugResult func(tx	*sql.Tx, ctx context.Context, sqlstr string, args ... interface{}) interface{}
	driverName	string
}


func (T	*DB) Open(driverName, dataSourceName string) (*DB, error) {
	var	err	error
	if T.DB	== nil {
		T.driverName=driverName
		T.DB, err =	sql.Open(driverName, dataSourceName)
	}
	return T, err
}
func (T	*DB) Close() error {
	if T.DB	== nil {
		return sql.ErrConnDone
	}
	err	:= T.DB.Close()
	T.DB = nil
	return err
}

func (T *DB) txCommit(tx *sql.Tx, e *error) {
	if *e != nil {
		if err := tx.Rollback(); err != nil {
			T.logf(err.Error())
		}
	}else if err := tx.Commit(); err != nil {
		T.logf(err.Error())
	}
}

//日志
func (T *DB) logf(format string, a ...interface{}) {
	if T.ErrorLog != nil {
		T.ErrorLog.Output(2, fmt.Sprintf(format+"\n", a...))
	}
	log.Printf(format+"\n", a...)
}

func (T *DB) debugFormat(sqlStr string) string {
	sqlStr =  strings.ReplaceAll(sqlStr, "\n", " ")
	sqlStr =  strings.ReplaceAll(sqlStr, "\t", "")
	sqlStr =  strings.ReplaceAll(sqlStr, "  ", " ")
	return sqlStr
}
//打印调试信息
func (T *DB) debugPrint(sqlStr string, args interface{}){
	if T.DebugPrint {
		sqlStr = T.debugFormat(sqlStr)
		rv := reflect.ValueOf(args)
		rv = reflect.Indirect(rv)
		var argsStr []string
		if rv.Kind() == reflect.Slice {
			for i:=0;i<rv.Len();i++ {
				srv := rv.Index(i)
				argsStr = append(argsStr, fmt.Sprintf("'%v'",vweb.InDirect(srv)))
			}
		}
		T.logf("\nprepare test as %s;\nexecute test(%v);\ndeallocate test;\n", sqlStr, strings.Join(argsStr,","))
	}
}

//支持写入(postgresql)，支持返回错误：error, ErrNoRows
func (T	*DB) Pexec(tx *sql.Tx,	sqlstr string, args	...interface{})	error {
	ctx	:= context.Background()
	var	cancel context.CancelFunc
	if T.Timeout !=	0 {
		ctx, cancel	= context.WithTimeout(ctx, T.Timeout)
		defer cancel()
	}
	return T.PexecContext(tx, ctx,	sqlstr,	args...)
}
func (T	*DB) PexecContext(tx *sql.Tx, ctx context.Context,	sqlstr string, args	...interface{})	(err error)	{
	result,	err	:= T.ExecContext(tx, ctx, sqlstr, args...)
	if err != nil {
		return
	}
	
	affected, err := result.RowsAffected()
	if err != nil {
		return 
	}
	if affected	== 0 {
		return ErrNoRows
	}
	return
}

//支持写入
func (T	*DB) Exec(tx *sql.Tx, sqlstr string, args ...interface{}) (result sql.Result, err error) {
	ctx	:= context.Background()
	var	cancel context.CancelFunc
	if T.Timeout !=	0 {
		ctx, cancel	= context.WithTimeout(ctx, T.Timeout)
		defer cancel()
	}
	return T.ExecContext(tx, ctx, sqlstr, args...)
}
func (T	*DB) ExecContext(tx	*sql.Tx, ctx context.Context, sqlstr string, args ...interface{}) (result sql.Result, err error) {
	if T.DB	== nil {
		return nil,	sql.ErrConnDone
	}
	
	T.debugPrint(sqlstr, args)
	if T.DebugResult != nil {
		if v, ok := T.DebugResult(tx, ctx, T.debugFormat(sqlstr), args...).(sql.Result); ok {
			return v, nil
		}
	}
	
	if tx == nil {
		tx,	err	= T.DB.Begin()
		if err != nil {
			return nil,	err
		}
		defer T.txCommit(tx, &err)
	}
	result,	err	= tx.ExecContext(ctx, sqlstr, args...)
	return
}

//支持读取/写入，支持返回错误：error, ErrNoRows
func (T	*DB) QueryRow(tx *sql.Tx, sqlstr string, args ...interface{}) Rower {
	ctx	:= context.Background()
	return T.QueryRowContext(tx, ctx, sqlstr, args...)
}
func (T	*DB) QueryRowContext(tx	*sql.Tx, ctx context.Context, sqlstr string, args ...interface{}) Rower {
	if T.DB	== nil {
		return &dbRow{err:sql.ErrConnDone}
	}
	
	T.debugPrint(sqlstr, args)
	if T.DebugResult != nil {
		if v, ok := T.DebugResult(tx, ctx, T.debugFormat(sqlstr), args...).(Rower); ok {
			return v
		}
	}
	
	//查询，用不上回滚
	var sqlstrTL = strings.ToLower(strings.TrimLeft(sqlstr, " \t\r\n"))
	if strings.HasPrefix(sqlstrTL, "select")	{
		return T.DB.QueryRowContext(ctx, sqlstr, args...)
	}
	
	//执行，需要回滚
	var	sqlTx *sql.Tx
	if tx == nil {
		var	err	error
		tx,	err	= T.DB.Begin()
		if err != nil {
			return &dbRow{err:err}
		}
		sqlTx =	tx
	}
	
	var cancel context.CancelFunc = func(){}
	if T.Timeout !=	0 {
		//这个上下文没有设置超时
		_, ok := ctx.Deadline()
		if !ok {
			ctx, cancel	= context.WithTimeout(ctx, T.Timeout)
		}
	}
	sqlRow := tx.QueryRowContext(ctx, sqlstr, args...)
	return &dbRow{tx:sqlTx,	r:sqlRow, cancel:cancel}
}

//支持读取/写入，支持返回错误：error, ErrNoRows
func (T	*DB) Query(tx *sql.Tx, sqlstr string, args ...interface{}) (data []map[string]interface{}, err error) {
	ctx	:= context.Background()
	var	cancel context.CancelFunc
	if T.Timeout !=	0 {
		ctx, cancel	= context.WithTimeout(ctx, T.Timeout)
		defer cancel()
	}
	return T.QueryContext(tx, ctx, sqlstr, args...)
}
func (T	*DB) QueryContext(tx *sql.Tx, ctx context.Context, sqlstr string, args ...interface{}) (data []map[string]interface{}, err error) {
	if T.DB	== nil {
		return nil,	sql.ErrConnDone
	}
	
	T.debugPrint(sqlstr, args)
	if T.DebugResult != nil {
		if v, ok := T.DebugResult(tx, ctx, T.debugFormat(sqlstr), args...).([]map[string]interface{}); ok {
			return v, nil
		}
	}
	
	var (
		sqlstrTL = strings.ToLower(strings.TrimLeft(sqlstr, " \t\r\n"))
		rows *sql.Rows
	)	
	if strings.HasPrefix(sqlstrTL, "select") {
		//查询，有返回
		rows, err =	T.DB.QueryContext(ctx, sqlstr, args...)
		if err != nil {
			return nil,	err
		}
	}else if strings.Index(sqlstrTL, "returning") == -1	{
		//插入或更新，没有返回
		return nil,	T.PexecContext(tx,	ctx, sqlstr, args...)
	}else{
		//插入或更新，有返回
		if tx == nil {
			tx,	err	= T.DB.Begin()
			if err != nil {
				return nil,	err
			}
			defer T.txCommit(tx, &err)
		}
		
		rows, err =	tx.QueryContext(ctx, sqlstr, args...)
		if err != nil {
			return nil,	err
		}
		defer rows.Close()
	}

	//列名处理
	columns, err :=	rows.Columns()
	if err != nil {
		return nil,	err
	}
	if T.FormatColumnName {
		for	index, col := range	columns	{
			columns[index]=keyToLower(col)
		}
	}
	
	//值处理
	columnTypes, err :=	rows.ColumnTypes();
	if err != nil {
		return nil,	err
	}
	var	values = make([]interface{}, len(columnTypes))
	for	i, tp := range columnTypes {
		st := tp.ScanType()
		if st == nil {
			st = reflect.TypeOf((*interface{})(nil))
		}
		if T.TypeTo != nil {
			if tto := T.TypeTo(tp); tto != nil {
				st = tto
			}
			
		}
		values[i] =	reflect.New(st).Interface()
	}
	
	//读出例值
	for	rows.Next()	{
		err	= rows.Scan(values...)
		if err != nil {
			return nil,	err
		}
		column := make(map[string]interface{})
		data = append(data,	column)
		for	i, value :=	range values {
			column[columns[i]] = reflect.Indirect( reflect.ValueOf(value) ).Interface()
		}
	}
	if data	== nil {
		err	= sql.ErrNoRows
	}
	return
}

//判断数据存在，支持返回错误：error, ErrRows, ErrNoRows
func (T	*DB) Has(sqlstr	string,	args ...interface{}) error {
	ctx	:= context.Background()
	var	cancel context.CancelFunc
	if T.Timeout !=	0 {
		ctx, cancel	= context.WithTimeout(ctx, T.Timeout)
		defer cancel()
	}
	return T.HasContext(ctx, sqlstr, args...)
}
func (T	*DB) HasContext(ctx	context.Context, sqlstr	string,	args ...interface{}) error {
	if T.DB	== nil {
		return sql.ErrConnDone
	}
	T.debugPrint(sqlstr, args)
	if T.DebugResult != nil {
		if v, ok := T.DebugResult(nil, ctx, T.debugFormat(sqlstr), args...).(error); ok {
			return v
		}
	}
	
	var	err	error
	var	rows *sql.Rows
	rows, err =	T.DB.QueryContext(ctx, sqlstr, args...)
	if err != nil {
		return err
	}
	defer rows.Close()
	if rows.Next() {
		return ErrRows
	}
	return ErrNoRows
}