package sqltable

import (
	"testing"

	"github.com/issue9/assert/v2"
)

func Test_NewSQLTable_Where(t *testing.T) {
	as := assert.New(t, true)
	ssql1 := Prepare(`select * from "B" $Where$`).Ors(`"G1"=?`, 1, `"G2"=?`, 2)
	ssql := Prepare(`select * from "A" $Where$`).And(`"B1"=?`, 1).And(`"B3"=?`, 3).Where("and", "B4=?", ssql1)

	as.Equal(ssql.Args(), []interface{}{1, 3, 1, 2})
	as.Equal(ssql.SQL(), "select * from \"A\" where \"B1\"=$1 and \"B3\"=$2 and B4=(select * from \"B\" where \"G1\"=$3 or \"G2\"=$4)")
}

func Test_NewSQLTable_Values(t *testing.T) {
	as := assert.New(t, true)

	ssql1 := Prepare(`select * from "B" $Where$`).Or(`"G1"=?`, 1).Or(`"G2"=?`, 2)
	ssql := Prepare(`insert into "B" $Values$ $Where$`).Values(`A1`, 1, `A2`, 1).Ands(`"B3"=?`, 3, `"B1"=?`, 1).Where("and", "B4=?", ssql1)
	as.Equal(ssql.Args(), []interface{}{3, 1, 1, 2, 1, 1})
	as.Equal(ssql.SQL(), "insert into \"B\" (\"A1\",\"A2\") values($5,$6) where \"B3\"=$1 and \"B1\"=$2 and B4=(select * from \"B\" where \"G1\"=$3 or \"G2\"=$4)")
}

func Test_NewSQLTable_Set(t *testing.T) {
	as := assert.New(t, true)
	ssql1 := Prepare(`select * from "B" $Where$ limit ?`).Or(`"G1"=?`, 6).Or(`"G2"=?`, 7)
	ssql1 = ssql1.ExtArgs(4)
	ssql := Prepare(`update into "A" $Set$ $Where$ limit ?`).Sets("A1", 2, "A2", 3).And(`"B1"=?`, 4).And(`"B3"=?`, 5).Where("and", "B4=?", ssql1)
	as.Equal(ssql.Args(1), []interface{}{1, 4, 5, 4, 6, 7, 2, 3})
	as.Equal(ssql.SQL(), "update into \"A\" set \"A1\"=$7,\"A2\"=$8 where \"B1\"=$2 and \"B3\"=$3 and B4=(select * from \"B\" where \"G1\"=$5 or \"G2\"=$6 limit $4) limit $1")
}

func Test_NewSQLTable_progress(t *testing.T) {
	as := assert.New(t, true)
	t1 := Prepare(`select id from t1 where a=? and b=?`)
	t1 = t1.ExtArgs(1)
	t1_1 := t1.ExtArgs(2)
	t1_2 := t1.ExtArgs(3)

	t2 := Prepare(`inster into t2 $Values$`)
	t2.Values("a", t1_1, "b", t1_2)

	as.Equal(t2.Args(), []interface{}{1, 2, 1, 3}).Equal(t2.SQL(), "inster into t2 (a,b) values((select id from t1 where a=$1 and b=$2),(select id from t1 where a=$3 and b=$4))")

	as.Equal(t1.Args(4), []interface{}{1, 4}).Equal(t1.SQL(), "select id from t1 where a=$1 and b=$2")

	as.Equal(t1_1.Args(), []interface{}{1, 2}).Equal(t1_1.SQL(), "select id from t1 where a=$1 and b=$2")

	as.Equal(t1_2.Args(), []interface{}{1, 3}).Equal(t1_2.SQL(), "select id from t1 where a=$1 and b=$2")

	t3 := Prepare(`select * from t3 $Where$`)
	t3.Ands("a=?", t2, "b=?", 5)

	as.Equal(t3.Args(), []interface{}{1, 2, 1, 3, 5}).Equal(t3.SQL(), "select * from t3 where a=(inster into t2 (a,b) values((select id from t1 where a=$1 and b=$2),(select id from t1 where a=$3 and b=$4))) and b=$5")
}

func Test_NewSQLTable_Copy(t *testing.T) {
	as := assert.New(t, true)
	t1 := Prepare(`select id from t1 $Where$ limit ?`)
	t1.And("a=?", 2)
	t1.Args(1)

	t11 := t1.Copy()
	t11.And("b=?", 3)
	t11.Args(4)

	as.Length(t1.swhere, 1)
	as.Length(t1.args, 2)

	as.Length(t11.swhere, 2)
	as.Length(t11.args, 3)
}
