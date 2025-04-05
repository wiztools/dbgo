package dbgo

import (
	"log"
	"testing"
	"time"
)

func TestWhereBuilderAnd(t *testing.T) {
	wb := NewWhereBuilder()
	wb.Add("acc_id", 101)
	wb.Add("usr_id", 102)
	wb.Add("type", "create")
	wb.AddBetween("order_time", time.Now(), time.Now())
	partialQry, vals := wb.GenAnd()
	log.Println(partialQry, "|", len(vals))
	if partialQry != " WHERE acc_id=? AND usr_id=? AND type=? AND order_time BETWEEN ? AND ?" {
		t.Fail()
	}
}

func TestWhereBuilderOr(t *testing.T) {
	wb := NewWhereBuilder()
	wb.Add("acc_id", 101)
	wb.Add("usr_id", 102)
	wb.Add("type", "create")
	wb.AddBetween("order_time", time.Now(), time.Now())
	partialQry, vals := wb.GenOr()
	log.Println(partialQry, "|", len(vals))
	if partialQry != " WHERE acc_id=? OR usr_id=? OR type=? OR order_time BETWEEN ? AND ?" {
		t.Fail()
	}
}

func TestWhereBuilderOnlyBetween(t *testing.T) {
	wb := NewWhereBuilder()
	wb.AddBetween("order_time", time.Now(), time.Now())
	partialQry, vals := wb.GenAnd()
	log.Println(partialQry, "|", len(vals))
	if partialQry != " WHERE order_time BETWEEN ? AND ?" {
		t.Fail()
	}
}

func TestWhereBuilderOnlyOrdrBy(t *testing.T) {
	wb := NewWhereBuilder()
	wb.AddOrdrByCols("salary", "age")
	partialQry, vals := wb.GenAnd()
	log.Println(partialQry, "|", len(vals))
	if partialQry != " ORDER BY salary, age" {
		t.Fail()
	}
}

func TestWhereBuilderPage(t *testing.T) {
	wb := NewWhereBuilder()
	wb.AddOrdrByCols("salary", "age")
	wb.SetPage(2, 20)
	partialQry, vals := wb.GenAnd()
	log.Println(partialQry, "|", len(vals))
	if partialQry != " ORDER BY salary, age LIMIT 20 OFFSET 40" {
		t.Fail()
	}
}
