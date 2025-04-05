package dbgo

import (
	"strconv"
	"strings"
)

type WhereBuilder struct {
	whrCols    []string
	whrColVals []any

	whrBetweenCols    []string
	whrBetweenColVals []any

	ordrByCols []string
	ordrByDesc bool

	limitRowCount  *int
	limitRowOffset *int
}

func NewWhereBuilder() *WhereBuilder {
	o := WhereBuilder{}

	o.whrCols = []string{}
	o.whrColVals = []any{}

	o.whrBetweenCols = []string{}
	o.whrBetweenColVals = []any{}

	o.ordrByCols = []string{}

	return &o
}

func (o *WhereBuilder) Add(colName string, colVal any) {
	o.whrCols = append(o.whrCols, colName)
	o.whrColVals = append(o.whrColVals, colVal)
}

func (o *WhereBuilder) AddBetween(colName string, from, to any) {
	o.whrBetweenCols = append(o.whrBetweenCols, colName)
	o.whrBetweenColVals = append(o.whrBetweenColVals, from)
	o.whrBetweenColVals = append(o.whrBetweenColVals, to)
}

func (o *WhereBuilder) AddOrdrByCols(colNames ...string) {
	o.ordrByCols = append(o.ordrByCols, colNames...)
}

func (o *WhereBuilder) SetOrdrByDesc() {
	o.ordrByDesc = true
}

func (o *WhereBuilder) SetLimitOffset(rowCount, rowOffset int) {
	o.limitRowCount = &rowCount
	o.limitRowOffset = &rowOffset
}

func (o *WhereBuilder) SetLimit(rowCount int) {
	o.limitRowCount = &rowCount
}

func (o *WhereBuilder) SetPage(pageNumber, pageSize int) {
	rowCount := pageSize
	rowOffset := pageSize * pageNumber
	o.limitRowCount = &rowCount
	o.limitRowOffset = &rowOffset
}

func (o *WhereBuilder) GenLimit() (sqlPartial string) {
	sb := strings.Builder{}
	if o.limitRowCount != nil {
		sb.WriteString(" LIMIT ")
		sb.WriteString(strconv.Itoa(*o.limitRowCount))

		if o.limitRowOffset != nil {
			sb.WriteString(" OFFSET ")
			sb.WriteString(strconv.Itoa(*o.limitRowOffset))
		}
	}
	sqlPartial = sb.String()
	return
}

func (o *WhereBuilder) gen(joiner string) (sqlPartial string, vals []any) {
	vals = []any{}

	sb := strings.Builder{}

	// Where conditions:
	if len(o.whrCols) != 0 {
		sb.WriteString(" WHERE ")

		sb.WriteString(strings.Join(o.whrCols, "=?"+joiner))
		sb.WriteString("=?")
		vals = append(vals, o.whrColVals...)
	}

	if len(o.whrBetweenCols) != 0 {
		if len(o.whrCols) != 0 {
			sb.WriteString(joiner)
		} else {
			sb.WriteString(" WHERE ")
		}
		for i, v := range o.whrBetweenCols {
			if i != 0 {
				sb.WriteString(joiner)
			}
			sb.WriteString(v)
			sb.WriteString(" BETWEEN ? AND ?")
		}
		vals = append(vals, o.whrBetweenColVals...)
	}

	// Order By:
	if len(o.ordrByCols) > 0 {
		sb.WriteString(" ORDER BY ")
		sb.WriteString(strings.Join(o.ordrByCols, ", "))

		if o.ordrByDesc {
			sb.WriteString(" DESC")
		}
	}

	// Pagination:
	sb.WriteString(o.GenLimit())

	sqlPartial = sb.String()

	return
}

func (o *WhereBuilder) GenAnd() (sqlPartial string, vals []any) {
	return o.gen(" AND ")
}

func (o *WhereBuilder) GenOr() (sqlPartial string, vals []any) {
	return o.gen(" OR ")
}
