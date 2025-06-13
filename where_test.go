package dbgo

import (
	"log"
	"testing"
	"time"
)

func TestWhereBuilderAnd(t *testing.T) {
	wb := NewWhereBuilder(AND)
	wb.Add("acc_id", 101)
	wb.Add("usr_id", 102)
	wb.Add("type", "create")
	wb.AddBetween("order_time", time.Now(), time.Now())
	partialQry, vals := wb.Gen()
	log.Println(partialQry, "|", len(vals))
	if partialQry != " WHERE acc_id=? AND usr_id=? AND type=? AND order_time BETWEEN ? AND ?" {
		t.Fail()
	}
}

func TestWhereBuilderOr(t *testing.T) {
	wb := NewWhereBuilder(OR)
	wb.Add("acc_id", 101)
	wb.Add("usr_id", 102)
	wb.Add("type", "create")
	wb.AddBetween("order_time", time.Now(), time.Now())
	partialQry, vals := wb.Gen()
	log.Println(partialQry, "|", len(vals))
	if partialQry != " WHERE acc_id=? OR usr_id=? OR type=? OR order_time BETWEEN ? AND ?" {
		t.Fail()
	}
}

func TestWhereBuilderOnlyBetween(t *testing.T) {
	wb := NewWhereBuilder(AND)
	wb.AddBetween("order_time", time.Now(), time.Now())
	partialQry, vals := wb.Gen()
	log.Println(partialQry, "|", len(vals))
	if partialQry != " WHERE order_time BETWEEN ? AND ?" {
		t.Fail()
	}
}

func TestWhereBuilderOnlyOrdrBy(t *testing.T) {
	wb := NewWhereBuilder(AND)
	wb.AddOrdrByCols("salary", "age")
	partialQry, vals := wb.Gen()
	log.Println(partialQry, "|", len(vals))
	if partialQry != " ORDER BY salary, age" {
		t.Fail()
	}
}

func TestWhereBuilderPage(t *testing.T) {
	wb := NewWhereBuilder(AND)
	wb.AddOrdrByCols("salary", "age")
	wb.SetPage(2, 20)
	partialQry, vals := wb.Gen()
	log.Println(partialQry, "|", len(vals))
	if partialQry != " ORDER BY salary, age LIMIT 20 OFFSET 40" {
		t.Fail()
	}
}

func TestGroupAnd(t *testing.T) {
	wb := NewWhereBuilder(AND)
	wb.Add("acc_id", 101)
	wb.Add("type", "create")
	grp := NewGroup(OR)
	grp.Add("usr_id", 202)
	grp.Add("usr_id", 303)
	{
		g := NewGroup(AND)
		g.Add("age", 20)
		g.Add("age", 30)
		grp.AddGroup(g)
	}
	wb.AddGroup(grp)
	partialQry, vals := wb.Gen()
	log.Println(partialQry, "|", len(vals))
	exp := " WHERE acc_id=? AND type=? AND (usr_id = ? OR usr_id = ? OR (age = ? AND age = ?))"
	if partialQry != exp {
		t.Log("", partialQry, "|", len(vals))
		t.Fail()
	}
	if len(vals) != 6 {
		t.Log("len(vals) != 6")
		t.Fail()
	}
}

func TestIsNull(t *testing.T) {
	wb := NewWhereBuilder(AND)
	wb.Add("acc_id", 101)
	wb.Add("type", "create")
	wb.AddIsNull("usr_id")
	wb.AddIsNotNull("age")
	partialQry, vals := wb.Gen()
	log.Println(partialQry, "|", len(vals))
	exp := " WHERE acc_id=? AND type=? AND usr_id IS NULL AND age IS NOT NULL"
	if partialQry != exp {
		t.Log("", partialQry, "|", len(vals))
		t.Fail()
	}
	if len(vals) != 2 {
		t.Log("len(vals) != 2")
		t.Fail()
	}
}
