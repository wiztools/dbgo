package dbgo

import "strings"

type WhereBuilder struct {
	colNames []string
	colVals  []any

	betweenColNames []string
	betweenColVals  []any
}

func NewWhereBuilder() *WhereBuilder {
	o := WhereBuilder{}

	o.colNames = []string{}
	o.colVals = []any{}

	o.betweenColNames = []string{}
	o.betweenColVals = []any{}

	return &o
}

func (o *WhereBuilder) Add(colName string, colVal any) {
	o.colNames = append(o.colNames, colName)
	o.colVals = append(o.colVals, colVal)
}

func (o *WhereBuilder) AddBetween(colName string, from, to any) {
	o.betweenColNames = append(o.betweenColNames, colName)
	o.betweenColVals = append(o.betweenColVals, from)
	o.betweenColVals = append(o.betweenColVals, to)
}

func (o *WhereBuilder) gen(joiner string) (sqlPartial string, vals []any) {
	vals = []any{}

	sb := strings.Builder{}

	if len(o.colNames) != 0 {
		sb.WriteString(" WHERE ")
	}

	sb.WriteString(strings.Join(o.colNames, "=?"+joiner))
	sb.WriteString("=?")
	vals = append(vals, o.colVals...)

	if len(o.betweenColNames) != 0 {
		if len(o.colNames) != 0 {
			sb.WriteString(joiner)
		} else {
			sb.WriteString(" WHERE ")
		}
		for i, v := range o.betweenColNames {
			if i != 0 {
				sb.WriteString(joiner)
			}
			sb.WriteString(v)
			sb.WriteString(" BETWEEN ? AND ?")
		}
		vals = append(vals, o.betweenColVals...)
	}

	sqlPartial = sb.String()

	return
}

func (o *WhereBuilder) GenAnd() (sqlPartial string, vals []any) {
	return o.gen(" AND ")
}

func (o *WhereBuilder) GenOr() (sqlPartial string, vals []any) {
	return o.gen(" OR ")
}
