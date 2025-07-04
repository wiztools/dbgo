package dbgo

import "database/sql"

// Tx is a wrapper around sql.Tx that provides additional functionality.
type Tx struct {
	tx *sql.Tx
}

func (o *DBGo) TxBegin() (*Tx, error) {
	tx, err := o.db.Begin()
	if err != nil {
		return nil, err
	}
	return &Tx{tx: tx}, nil
}

func (t *Tx) Commit() error {
	return t.tx.Commit()
}

func (t *Tx) Rollback() error {
	return t.tx.Rollback()
}

func (t *Tx) Exec(query string, args ...any) (sql.Result, error) {
	return t.tx.Exec(query, args...)
}

func (t *Tx) ExecGetLastInsertId(query string, args ...any) (*int64, error) {
	res, err := t.Exec(query, args...)
	if err != nil {
		return nil, err
	}
	o, err := res.LastInsertId()
	return &o, err
}

func (t *Tx) ExecGetRowsAffected(query string, args ...any) (*int64, error) {
	res, err := t.Exec(query, args...)
	if err != nil {
		return nil, err
	}
	o, err := res.RowsAffected()
	return &o, err
}

func (t *Tx) Query(query string, args ...any) (*sql.Rows, error) {
	return t.tx.Query(query, args...)
}

func (t *Tx) QueryRow(query string, args ...any) *sql.Row {
	return t.tx.QueryRow(query, args...)
}
