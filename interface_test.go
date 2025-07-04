package dbgo

import "testing"

func TestInterface(t *testing.T) {
	db := DBGo{}
	fn := func(o DBOps) {}
	fn(&db)

	tx := Tx{}
	fn(&tx)
}
