package dbgo

import "time"

// Scanner is an interface that defines the Scan method that is part of sql.Rows and sql.Row.
// This helps modularize the scanning code to one function that can be reused across different types.
type Scanner interface {
	Scan(dest ...any) error
}

func ScanInt(row Scanner) (*int64, error) {
	var num int64
	err := row.Scan(&num)
	return &num, err
}

func ScanFloat(row Scanner) (*float64, error) {
	var f float64
	err := row.Scan(&f)
	return &f, err
}

func ScanString(row Scanner) (*string, error) {
	var str string
	err := row.Scan(&str)
	return &str, err
}

func ScanBool(row Scanner) (*bool, error) {
	var b bool
	err := row.Scan(&b)
	return &b, err
}

func ScanTime(row Scanner) (*time.Time, error) {
	var t time.Time
	err := row.Scan(&t)
	return &t, err
}

func ScanBytes(row Scanner) (*[]byte, error) {
	var b []byte
	err := row.Scan(&b)
	return &b, err
}
