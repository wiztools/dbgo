package dbgo

import "testing"

func TestInterface(t *testing.T) {
	db := DBGo{}
	fn := func(o DB) {}
	fn(&db)

	tx := Tx{}
	fn(&tx)
}
