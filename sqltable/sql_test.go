package sqltable

import (
	"testing"
)

func Test_NewSQLTable_Where(t *testing.T){
	var ssql1 = NewSQLTable().Prepare(`select * from "B" $Where$`).Ors(`"G1"=?`,1,`"G2"=?`,2)
	var ssql = NewSQLTable().Prepare(`select * from "A" $Where$`).And(`"B1"=?`,1).And(`"B3"=?`,3).Where("and","B4=?",ssql1)
	t.Log(ssql.SQL())
}

func Test_NewSQLTable_Values(t *testing.T){
	var ssql1 = NewSQLTable().Prepare(`select * from "B" $Where$`).Or(`"G1"=?`,1).Or(`"G2"=?`,2)
	var ssql = NewSQLTable().Prepare(`insert into "B" $Values$ $Where$`).Values(`A1`,1,`A2`,1).Ands(`"B3"=?`,3,`"B1"=?`,1).Where("and","B4=?",ssql1)
	t.Log(ssql.SQL())
}

func Test_NewSQLTable_Set(t *testing.T){
	var ssql1 = NewSQLTable().Prepare(`select * from "B" $Where$ limit ?`).Or(`"G1"=?`,1).Or(`"G2"=?`,2)
	ssql1.ExtArgs(1)
	var ssql = NewSQLTable().Prepare(`update into "A" $Set$ $Where$ limit ?`).Sets("A1",1,"A2",2).And(`"B1"=?`,1).And(`"B3"=?`,3).Where("and","B4=?",ssql1)
	t.Log(ssql.Args(2), ssql.SQL())
}


