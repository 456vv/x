package	db
import(
	"github.com/456vv/vbody"
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

var	ErrRows	= errors.New("sql: have	rows in	result set")
var ErrNoRows = sql.ErrNoRows

type sqlRower interface{
	Scan(dest ...interface{}) error
	Err() error
}
type dbRow struct{
	tx	*sql.Tx
	r	*sql.Row
	err	error
}
func (T	*dbRow)	Scan(dest ...interface{}) error	{
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
	ErrorLog	*log.Logger
	Debug		bool
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

//日志
func (T *DB) logf(format string, a ...interface{}) {
	if T.ErrorLog != nil {
		T.ErrorLog.Output(2, fmt.Sprintf(format+"\n", a...))
	}
	log.Printf(format+"\n", a...)
}
//打印调试信息
func (T *DB) debugPrint(sqlStr string, args interface{}){
	if T.Debug {
		sqlStr =  strings.ReplaceAll(sqlStr, "\n", " ")
		sqlStr =  strings.ReplaceAll(sqlStr, "\t", "")
		sqlStr =  strings.ReplaceAll(sqlStr, "  ", " ")
		T.logf("%s %#v\n", sqlStr, args)
	}
}
//支持写入(postgresql)
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
	
	if tx == nil {
		tx,	err	= T.DB.Begin()
		if err != nil {
			return nil,	err
		}
		defer func(){
			if err != nil {
				e := tx.Rollback()
				if e != nil {
					T.logf(e.Error())
				}
				return
			}
			err	= tx.Commit()
			if err != nil {
				T.logf(err.Error())
			}
		}()
	}
	T.debugPrint(sqlstr, args)
	result,	err	= tx.ExecContext(ctx, sqlstr, args...)
	return
}

//支持读取/写入
func (T	*DB) QueryRow(tx *sql.Tx, sqlstr string, args ...interface{}) sqlRower {
	ctx	:= context.Background()
	var	cancel context.CancelFunc
	if T.Timeout !=	0 {
		ctx, cancel	= context.WithTimeout(ctx, T.Timeout)
		defer cancel()
	}
	return T.QueryRowContext(tx, ctx, sqlstr, args...)
}
func (T	*DB) QueryRowContext(tx	*sql.Tx, ctx context.Context, sqlstr string, args ...interface{}) sqlRower {
	if T.DB	== nil {
		return &dbRow{err:sql.ErrConnDone}
	}
	sqlstrTL :=	strings.TrimLeft(sqlstr, " \t")
	sqlstrL	:= strings.ToLower(sqlstrTL[:6])
	
	//查询，用不上回滚
	if strings.HasPrefix(sqlstrL, "select")	{
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
	
	T.debugPrint(sqlstr, args)
	sqlRow := tx.QueryRowContext(ctx, sqlstr, args...)
	return &dbRow{tx:sqlTx,	r:sqlRow}
}

//支持读取/写入
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

	sqlstrTL :=	strings.ToLower(strings.TrimLeft(sqlstr, " \t\r\n"))
	
	var	rows *sql.Rows
	if strings.HasPrefix(sqlstrTL, "select") {
		//查询，有返回
		T.debugPrint(sqlstr, args)
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
			defer func(){
				if err != nil {
					e := tx.Rollback()
					if e != nil {
						T.logf(e.Error())
					}
					return
				}
				err	= tx.Commit()
				if err != nil {
					T.logf(err.Error())
				}
			}()
		}
		
		T.debugPrint(sqlstr, args)
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
	var	values	= make([]interface{}, len(columnTypes))
	for	i, tp := range columnTypes {
		st := tp.ScanType()
		if st == nil {
			st = reflect.TypeOf((*interface{})(nil))
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

//支持读取/写入
//line={"A":*vbody.Reader, "B":*vbody.Reader}
func (T	*DB) QueryReader(tx	*sql.Tx, id	string,	sqlstr string, args	...interface{})	(line map[interface{}]*vbody.Reader, err error)	{
	ctx	:= context.Background()
	var	cancel context.CancelFunc
	if T.Timeout !=	0 {
		ctx, cancel	= context.WithTimeout(ctx, T.Timeout)
		defer cancel()
	}
	return T.QueryReaderContext(tx,	ctx, id, sqlstr, args...)
}
func (T	*DB) QueryReaderContext(tx *sql.Tx,	ctx	context.Context, id	string,	sqlstr string, args	...interface{})	(line map[interface{}]*vbody.Reader, err error)	{
	if T.DB	== nil {
		return nil,	sql.ErrConnDone
	}
	
	sqlstrTL :=	strings.TrimLeft(sqlstr, " \t\r\n")
	sqlstrL	:= strings.ToLower(sqlstrTL[:6])
	if T.FormatColumnName {
		id = keyToLower(id)
	}
	var	rows *sql.Rows
	if strings.HasPrefix(sqlstrL, "select")	{
		T.debugPrint(sqlstr, args)
		rows, err =	T.DB.QueryContext(ctx, sqlstr, args...)
		if err != nil {
			return nil,	err
		}
	}else{
		if tx == nil {
			tx,	err	= T.DB.Begin()
			if err != nil {
				return nil,	err
			}
			defer func(){
				if err != nil {
					e := tx.Rollback()
					if e != nil {
						T.logf(e.Error())
					}
					return
				}
				err	= tx.Commit()
				if err != nil {
					T.logf(err.Error())
				}
			}()
		}
		
		T.debugPrint(sqlstr, args)
		rows, err =	tx.QueryContext(ctx, sqlstr, args...)
		if err != nil {
			return nil,	err
		}
		defer rows.Close()
	}
	
	//字段名称
	columns, err :=	rows.Columns()
	if err != nil {
		return nil,	err
	}
	if T.FormatColumnName {
		for	index, col := range	columns	{
			//字段名称转小写开头
			columns[index]=keyToLower(col)
		}
	}
	
	//字段类型
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
		values[i] =	reflect.New(st).Interface()
	}
	
	//读取行
	for	rows.Next()	{
		//扫描值放在values
		err	= rows.Scan(values...)
		if err != nil {
			return nil,	err
		}
		if line	== nil {
			line = make(map[interface{}]*vbody.Reader)
		}
		data	:= make(map[string]interface{})
		for	i, value :=	range values {
			data[columns[i]] = reflect.Indirect( reflect.ValueOf(value)	).Interface()
		}
		line[data[id]]=vbody.NewReader(data)
	}
	if line	== nil {
		err	= sql.ErrNoRows
	}
	return
}

//判断数据存在
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
