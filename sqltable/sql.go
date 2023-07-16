package sqltable

import (
	"fmt"
	"strconv"
	"strings"
)

//使用方法：
//
//Prepare(`select "A1" from "A" $Where$ order by "ID" limit 10`).And("A1=?", value)
//Prepare(`insert into "A"$Values$ $Where$`).Values("A1",value).And("A1=?", value)
//Prepare(`update "A" $Set$ $Where$`).Set("A1",value).And("A1=?", value)
//Prepare(`delete from "A" $Where$`).And("A1=?", value)
//Prepare(`insert into "A"$Values$ on conflict("ID") do update set $Excluded$  $Where$`).Values("A1",value).And("A1=?", value)
//
//Args(args ...interface{}) []interface{}
//SQL() string
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

type SQLTable struct {
	source     string // 未处理的sql语句
	result     string // 最终sql语句
	swhere     []*sqlTableWhere
	sexcluded  []string
	count      int
	args       []interface{} // 参数
	extArgs    []interface{}
	columnMark string
	valuesKeys []string
	valuesVals []interface{}
	setKeys    []string
	setVals    []interface{}
	types      map[string]string
	funcs      map[string]string
}

func Prepare(s string) *SQLTable {
	st := &SQLTable{
		types: make(map[string]string),
		funcs: make(map[string]string),
	}

	for i := 1; true; i++ {
		dollar := fmt.Sprint("$", i)
		if !strings.Contains(s, dollar) {
			break
		}
		s = strings.ReplaceAll(s, dollar, fmt.Sprint("@", i))
	}
	st.source = s
	return st
}

func (T *SQLTable) progress(val interface{}) string {
	if t, ok := val.(*SQLTable); ok {
		statement := t.clone()
		statement.count = T.count
		s := "(" + statement.SQL() + ")"
		T.count = statement.count
		T.args = append(T.args, statement.args...)

		return s
	}
	T.count++
	T.args = append(T.args, val)
	return "$" + strconv.Itoa(T.count)
}

// 支持设置转换类型
func (T *SQLTable) ToTypes(a ...interface{}) *SQLTable {
	T.pairArgs(func(i int, k, v interface{}) bool {
		key, ok := k.(string)
		if !ok {
			panic(fmt.Sprintf("SQLTable.ToTypes in %d, 键名是%#v, 不是字符串!", i, k))
		}
		val, ok := v.(string)
		if !ok {
			panic(fmt.Sprintf("SQLTable.ToTypes in %d, 键值是%#v, 不是字符串!", i, v))
		}
		keyName := T.addColumnMark(key)
		T.types[keyName] = val
		return false
	}, a...)
	return T
}

// 支持设置函数转转
func (T *SQLTable) ToFuncs(a ...interface{}) *SQLTable {
	T.pairArgs(func(i int, k, v interface{}) bool {
		key, ok := k.(string)
		if !ok {
			panic(fmt.Sprintf("SQLTable.ToFuncs in %d, 键名是%#v, 不是字符串!", i, k))
		}
		val, ok := v.(string)
		if !ok {
			panic(fmt.Sprintf("SQLTable.ToFuncs in %d, 键值是%#v, 不是字符串!", i, v))
		}
		keyName := T.addColumnMark(key)
		T.funcs[keyName] = val
		return false
	}, a...)
	return T
}

func (T *SQLTable) pairArgs(a func(i int, k, v interface{}) bool, b ...interface{}) {
	l := len(b)
	if l%2 != 0 {
		return
	}
	for i := 0; i < l; i += 2 {
		if a(i, b[i], b[i+1]) {
			return
		}
	}
}

// 可以对字段名加引号或者其它什么的
func (T *SQLTable) ColumnMark(q string) *SQLTable {
	T.columnMark = q
	return T
}

func (T *SQLTable) addColumnMark(col string) string {
	mark := T.columnMark
	if mark == "" {
		mark = `"`
	}
	for _, c := range col {
		if c >= 'A' && c <= 'Z' {
			return mark + col + mark
		}
	}
	return col
}

// 他会替换语句中的 $Values$ 符号
func (T *SQLTable) Value(column string, value interface{}) *SQLTable {
	keyName := T.addColumnMark(column)
	T.valuesKeys = append(T.valuesKeys, keyName)
	T.valuesVals = append(T.valuesVals, value)
	return T
}

// 读取所有值
func (T *SQLTable) ValueValues(v ...interface{}) []interface{} {
	vals := append(T.valuesVals, v...)
	return quoteValue(vals)
}

// 他会替换语句中的 $Values$ 符号
func (T *SQLTable) Values(a ...interface{}) *SQLTable {
	return T.more(T.Value, a...)
}

func (T *SQLTable) values(s string) string {
	var values []string
	for index, val := range T.valuesVals {
		keyName := T.valuesKeys[index]
		keyVal := T.progress(val)
		if typ, ok := T.types[keyName]; ok {
			keyVal += "::" + typ
		}
		if funcs, ok := T.funcs[keyName]; ok {
			keyVal = strings.ReplaceAll(funcs, "?", keyVal)
		}
		values = append(values, keyVal)
	}
	if len(T.valuesKeys) > 0 {
		valuesKeys := strings.Join(T.valuesKeys, ",")
		valuesSerial := strings.Join(values, ",")
		s = strings.Replace(s, "$Values$", fmt.Sprintf(`(%s) values(%s)`, valuesKeys, valuesSerial), 1)
		s = strings.Replace(s, "$ValuesKeys$", valuesKeys, 1)
		s = strings.Replace(s, "$ValuesSerial$", valuesSerial, 1)
	}
	return s
}

// 他会替换语句中的 $Where$ 符号
func (T *SQLTable) And(column string, value interface{}) *SQLTable {
	return T.Where("and", column, value)
}

// 他会替换语句中的 $Where$ 符号
func (T *SQLTable) Ands(a ...interface{}) *SQLTable {
	return T.more(T.And, a...)
}

// 他会替换语句中的 $Where$ 符号
func (T *SQLTable) Or(column string, value interface{}) *SQLTable {
	return T.Where("or", column, value)
}

// 他会替换语句中的 $Where$ 符号
func (T *SQLTable) Ors(a ...interface{}) *SQLTable {
	return T.more(T.Or, a...)
}

// 他会替换语句中的 $Where$ 符号
func (T *SQLTable) Where(symbol, column string, value interface{}) *SQLTable {
	T.swhere = append(T.swhere, &sqlTableWhere{symbol, column, value})
	return T
}

func (T *SQLTable) where(s string) string {
	var (
		swhere     []string
		isOptional bool = strings.Contains(s, "*Where*")
	)
	for _, w := range T.swhere {
		if len(swhere) > 0 || isOptional {
			swhere = append(swhere, w.symbol)
		}
		swhere = append(swhere, strings.Replace(w.key, "?", T.progress(w.val), 1))
	}

	where := strings.Join(swhere, " ")
	s = strings.Replace(s, "*Where*", where, 1)
	if len(swhere) != 0 {
		where = "where " + where
	}
	s = strings.Replace(s, "$Where$", where, 1)
	return s
}

// 他会替换语句中的 $Set$ 符号
func (T *SQLTable) Set(column string, value interface{}) *SQLTable {
	keyName := T.addColumnMark(column)
	T.setKeys = append(T.setKeys, keyName)
	T.setVals = append(T.setVals, value)
	return T
}

// 他会替换语句中的 $Set$ 符号
func (T *SQLTable) Sets(a ...interface{}) *SQLTable {
	return T.more(T.Set, a...)
}

func (T *SQLTable) more(f func(string, interface{}) *SQLTable, a ...interface{}) *SQLTable {
	l := len(a)
	if l%2 != 0 {
		return T
	}
	for i := 0; i < l; i += 2 {
		c, ok := a[i].(string)
		if !ok {
			return T
		}
		val := a[i+1]
		f(c, val)
	}
	return T
}

func (T *SQLTable) set(s string) string {
	var (
		sets       []string
		setsSerial []string
	)
	for index, keyName := range T.setKeys {
		val := T.setVals[index]
		keyVal := T.progress(val)
		if typ, ok := T.types[keyName]; ok {
			keyVal += "::" + typ
		}
		if funcs, ok := T.funcs[keyName]; ok {
			keyVal = strings.ReplaceAll(funcs, "?", keyVal)
		}
		sets = append(sets, fmt.Sprintf(`%s=%s`, keyName, keyVal))
		setsSerial = append(setsSerial, keyVal)
	}

	set := strings.Join(sets, ",")
	s = strings.Replace(s, "*Set*", set, 1)
	if len(sets) > 0 {
		set = "set " + set
	}
	s = strings.Replace(s, "$Set$", set, 1)
	s = strings.Replace(s, "$SetKeys$", strings.Join(T.setKeys, ","), 1)
	s = strings.Replace(s, "$SetSerial$", strings.Join(setsSerial, ","), 1)
	return s
}

// 读取所有值
func (T *SQLTable) SetValues(v ...interface{}) []interface{} {
	vals := append(T.setVals, v...)
	return quoteValue(vals)
}

// 他会替换语句中的 $Excluded$ 符号
// excluded 用于在插入存在，则更新时，排除的例名
func (T *SQLTable) Excluded(a ...string) *SQLTable {
	T.sexcluded = append(T.sexcluded, a...)
	return T
}

func (T *SQLTable) excluded(s string) string {
	var sets []string
	for _, val := range T.sexcluded {
		val = T.addColumnMark(val)
		if sliceContainString(T.valuesKeys, val) {
			sets = append(sets, fmt.Sprintf(`%[1]s=excluded.%[1]s`, val))
		}
	}
	if len(sets) > 0 {
		s = strings.Replace(s, "$Excluded$", strings.Join(sets, ","), 1)
	}
	return s
}

type sqlTableWhere struct {
	symbol string
	key    string
	val    interface{}
}

func (T *SQLTable) assemble(source string) string {
	source = T.where(source)
	source = T.values(source)
	source = T.set(source)
	source = T.excluded(source)
	return source
}

func (T *SQLTable) assembleStatement(args ...interface{}) string {
	if T.result != "" {
		return T.result
	}
	source := T.source
	for index, arg := range append(T.extArgs, args...) {
		dollar := T.progress(arg)
		essential := fmt.Sprint("@", index+1)
		if strings.Contains(source, essential) {
			source = strings.ReplaceAll(source, essential, dollar)
			continue
		}
		source = strings.Replace(source, "?", dollar, 1)
	}
	T.result = T.assemble(source)
	return T.result
}

// 返回sql语句，需要先调用.Args方法再调用该方法
//
//	string 格式后的SQL句子
func (T *SQLTable) SQL() string {
	return T.assembleStatement()
}

func (T *SQLTable) clone() *SQLTable {
	t := new(SQLTable)
	*t = *T
	t.result = ""
	t.count = 0
	t.args = nil
	return t
}

func (T *SQLTable) Copy() *SQLTable {
	return T.clone()
}

// 扩展参数
//
//	args ...interface{}		额外参数，在非解析SQL句子里的？
func (T *SQLTable) ExtArgs(args ...interface{}) *SQLTable {
	t := T.clone()
	t.extArgs = append(T.extArgs, args...)
	return t
}

// 载入参数，他会替换语句中的？符号
//
//	args ...interface{}		额外参数，在非解析SQL句子里的？
//	[]interface{}			返回所有参数
func (T *SQLTable) Args(args ...interface{}) []interface{} {
	T.assembleStatement(args...)
	return T.args
}
