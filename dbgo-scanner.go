package dbgo

// Scanner is an interface that defines the Scan method that is part of sql.Rows and sql.Row.
// This helps modularize the scanning code to one function that can be reused across different types.
type Scanner interface {
	Scan(dest ...any) error
}
