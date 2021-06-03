package sqltable

import (
	"strings"
	"fmt"
	"strconv"
)
//使用方法：
//
//NewSQLTable().ColumnMark(`"`).Prepare(`select "A1" from "A" $Where$ order by "ID" limit 10`).And("A1=?", value)
//NewSQLTable().ColumnMark(`"`).Prepare(`insert into "A"$Values$ $Where$`).Values("A1",value).And("A1=?", value)
//NewSQLTable().ColumnMark(`"`).Prepare(`update "A" $Set$ $Where$`).Set("A1",value).And("A1=?", value)
//NewSQLTable().ColumnMark(`"`).Prepare(`delete from "A" $Where$`).And("A1=?", value)
//NewSQLTable().ColumnMark(`"`).Prepare(`insert into "A"$Values$ on conflict("ID") do update set $Excluded$  $Where$`).Values("A1",value).And("A1=?", value)
//
//Args(args ...interface{}) []interface{}
//SQL() string
//Prepare(s string) *SQLTable
//Excluded(a ...string) *SQLTable
//Sets(a ...interface{}) *SQLTable
//Set(column string, value interface{}) *SQLTable
//Where(symbol, column string, value interface{}) *SQLTable
//Ors(a ...interface{}) *SQLTable
//Or(column string, value interface{}) *SQLTable
//Ands(a ...interface{}) *SQLTable
//And(column string, value interface{}) *SQLTable
//Values(a ...interface{}) *SQLTable
//Value(column string, value interface{}) *SQLTable
//ColumnMark(q string) *SQLTable


type SQLTable struct{
	ssql		string								//未处理的sql语句
	rsqlt		string								//最终sql语句
	svalues		map[string]interface{}
	swhere		[]*sqlTableWhere
	sset 		map[string]interface{}
	sexcluded 	[]string
	count		int
	args		[]interface{}						//参数
	prog		func(val interface{}) string
	columnMark	string
}
func NewSQLTable() *SQLTable{
	return &SQLTable{
		svalues	: make(map[string]interface{}),
		sset	: make(map[string]interface{}),
	}
}
func (T *SQLTable) reset(){
	if T.count == 0 {
		return
	}
	T.rsqlt = ""
	T.count = 0
	T.args	= []interface{}{}
}
func (T *SQLTable) progress(val interface{}) string {
	if T.prog == nil {
		T.count++
		T.args = append(T.args, val)
		return "$"+strconv.Itoa(T.count)
	}
	return T.prog(val)
}

//可以对字段名加引号或者其它什么的
func (T *SQLTable) ColumnMark(q string) *SQLTable  {
	T.columnMark=q
	return T
}
func (T *SQLTable) addColumnMark(col string) string {
	mark := T.columnMark
	if mark == "" {
		mark = `"`
	}
	for _, c := range col {
		if c >= 'A' && c <= 'Z' {
			return mark+ col +mark
		}
	}
	return col
}

//他会替换语句中的 $Values$ 符号
func (T *SQLTable) Value(column string, value interface{}) *SQLTable {
	T.reset()
	T.svalues[column]=value
	return T
}

//他会替换语句中的 $Values$ 符号
func (T *SQLTable) Values(a ...interface{}) *SQLTable {
	return T.more(T.Value, a...)
}
func (T *SQLTable) values(s string) string {
	var (
		conumns []string
		values []string
	)
	for key, val := range T.svalues {
		conumns = append(conumns,  T.addColumnMark(key) )
		values	= append(values, T.progress(val))
	}
	if len(conumns) >0 {
		s = strings.Replace(s, "$Values$", fmt.Sprintf(`(%s) values(%s)`, strings.Join(conumns, ","), strings.Join(values, ",")), 1)
	}
	
	return s
}

//他会替换语句中的 $Where$ 符号
func (T *SQLTable) And(column string, value interface{}) *SQLTable {
	T.reset()
	return T.Where("and", column, value)
}

//他会替换语句中的 $Where$ 符号
func (T *SQLTable) Ands(a ...interface{}) *SQLTable {
	return T.more(T.And, a...)
}

//他会替换语句中的 $Where$ 符号
func (T *SQLTable) Or(column string, value interface{}) *SQLTable {
	T.reset()
	return T.Where("or", column, value)
}

//他会替换语句中的 $Where$ 符号
func (T *SQLTable) Ors(a ...interface{}) *SQLTable {
	return T.more(T.Or, a...)
}
//他会替换语句中的 $Where$ 符号
func (T *SQLTable) Where(symbol, column string, value interface{}) *SQLTable {
	T.reset()
	//if !strings.Contains(column, "?") {
	//	column = fmt.Sprintf("%s=?" T.columnMark + column + T.columnMark)
	//}
	T.swhere = append(T.swhere, &sqlTableWhere{symbol, column, value})
	return T
}
func (T *SQLTable) where(s string) string {
	var (
	 	swhere []string
	 	isOptional bool = strings.Contains(s, "*Where*")
	 ) 
	for _, w := range T.swhere {
		if len(swhere)>0 || isOptional {
			swhere = append(swhere, w.symbol)
		}
		if sqlTable, ok := w.val.(*SQLTable); ok {
			sqlTable.prog = T.progress
			swhere = append(swhere, strings.Replace(w.key, "?", fmt.Sprintf("(%s)", sqlTable.SQL()), 1))
			continue
		}
		swhere = append(swhere, strings.Replace(w.key, "?", T.progress(w.val), 1))
	}
	
	where := strings.Join(swhere, " ")
	s = strings.Replace(s, "*Where*", where, 1)
	if len(swhere) != 0 {
		where = "where "+where
	}
	s = strings.Replace(s, "$Where$", where, 1)
	return s

}

//他会替换语句中的 $Set$ 符号
func (T *SQLTable) Set(column string, value interface{}) *SQLTable {
	T.reset()
	T.sset[column]=value
	return T
}

//他会替换语句中的 $Set$ 符号
func (T *SQLTable) Sets(a ...interface{}) *SQLTable {
	return T.more(T.Set, a...)
}
func (T *SQLTable) more(f func(string, interface{}) *SQLTable ,a ...interface{}) *SQLTable {
	var l = len(a)
	if l%2 != 0 {
		return T
	}
	for i:=0;i<l;i+=2 {
		c,ok := a[i].(string)
		if !ok {
			return T
		}
		f(c, a[i+1])
	}
	return T
}
func (T *SQLTable) set(s string) string {
	var sets []string
	for key, val  := range T.sset {
		sets = append(sets, fmt.Sprintf(`%s=%s`, T.addColumnMark(key), T.progress(val)))
	}
	
	set := strings.Join(sets,",")
	s = strings.Replace(s, "*Set*", set, 1)
	if len(sets) > 0 {
		set = "set "+set
	}
	s = strings.Replace(s, "$Set$", set, 1)
	return s
}

//他会替换语句中的 $Excluded$ 符号
//excluded 用于在插入存在，则更新时，排除的例名
func (T *SQLTable) Excluded(a ...string) *SQLTable {
	T.sexcluded = append(T.sexcluded, a...)
	return T
}
func (T *SQLTable) excluded(s string) string {
	var sets []string
	for _, val  := range T.sexcluded {
		sets = append(sets, fmt.Sprintf(`%[1]s=excluded.%[1]s`, T.addColumnMark(val)))
	}
	if len(sets) >0 {
		s = strings.Replace(s,"$Excluded$", strings.Join(sets,","), 1)
	}
	return s
}
type sqlTableWhere struct {
	symbol	string
	key		string
	val		interface{}
}

//解析sql语句
//	s string 待解析的sql句子
func (T *SQLTable) Prepare(s string) *SQLTable {
	T.ssql = s
	return T
}

func (T *SQLTable) parse() {
	T.rsqlt = T.ssql
	T.rsqlt = T.where(T.rsqlt)
	T.rsqlt = T.values(T.rsqlt)
	T.rsqlt = T.set(T.rsqlt)
	T.rsqlt = T.excluded(T.rsqlt)
}

//返回sql语句，需要先调用.Args方法再调用该方法
//	string 格式后的SQL句子
func (T *SQLTable) SQL() string {
	if T.rsqlt == ""{
		T.parse()
	}
	return T.rsqlt
}

//载入参数，他会替换语句中的？符号
//	args ...interface{}		额外参数，在非解析SQL句子里的？
//	[]interface{}			返回所有参数
func (T *SQLTable) Args(args ...interface{}) []interface{}{
	if T.rsqlt == ""{
		T.parse()
	}
	for _, arg := range args {
		T.count++
		T.args = append(T.args, arg)
		T.rsqlt = strings.Replace(T.rsqlt, "?", "$"+strconv.Itoa(T.count), 1)
	}
	return T.args
}
