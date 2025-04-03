package dbgo

import "database/sql"

type DBGo struct {
	db *sql.DB
}

func New(db *sql.DB) *DBGo {
	return &DBGo{db: db}
}

func (o *DBGo) Exec(qry string, args ...any) (sql.Result, error) {
	var err error
	var res sql.Result
	if res, err = o.db.Exec(qry, args...); err != nil {
		return nil, err
	}
	return res, nil
}

func (o *DBGo) ExecGetLastInsertId(qry string, args ...any) (int64, error) {
	if res, err := o.Exec(qry, args...); err != nil {
		return 0, err
	} else {
		return res.LastInsertId()
	}
}

func (o *DBGo) ExecGetRowsAffected(qry string, args ...any) (int64, error) {
	if res, err := o.Exec(qry, args...); err != nil {
		return 0, err
	} else {
		return res.RowsAffected()
	}
}
