package dbgo

import "database/sql"

type DB interface {
	Exec(query string, args ...any) (sql.Result, error)
	ExecGetLastInsertId(query string, args ...any) (*int64, error)
	ExecGetRowsAffected(query string, args ...any) (*int64, error)
	Query(query string, args ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) *sql.Row
}
