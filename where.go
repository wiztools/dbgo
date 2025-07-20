package dbgo

import (
	"fmt"
	"strings"
)

type FuseType string

const (
	AND FuseType = "AND"
	OR  FuseType = "OR"
)

type Group interface {
	Wheres
	gen(sb *strings.Builder, params *[]any)
}

type Wheres interface {
	SetFuse(fuseType FuseType)
	AddRaw(conditions ...string)
	Add(colName string, colVal any)
	AddIsNull(colName string)
	AddIsNotNull(colName string)
	AddBetween(colName string, from, to any)
	AddGroup(group Group)
}

type GroupImpl struct {
	fuseType FuseType

	whrRawCols []string

	whrCols    []string
	whrColVals []any

	whrBetweenCols    []string
	whrBetweenColVals []any

	whrIsNullCols    []string
	whrIsNotNullCols []string

	groups []Group
}

func NewGroup(fuseType FuseType) Group {
	return &GroupImpl{
		fuseType:          fuseType,
		whrRawCols:        []string{},
		whrCols:           []string{},
		whrColVals:        []any{},
		whrBetweenCols:    []string{},
		whrBetweenColVals: []any{},
		whrIsNullCols:     []string{},
		whrIsNotNullCols:  []string{},
		groups:            []Group{},
	}
}

func (o *GroupImpl) SetFuse(ft FuseType) {
	o.fuseType = ft
}

func (o *GroupImpl) AddRaw(conditions ...string) {
	o.whrRawCols = append(o.whrRawCols, conditions...)
}

func (o *GroupImpl) Add(colName string, colVal any) {
	o.whrCols = append(o.whrCols, colName)
	o.whrColVals = append(o.whrColVals, colVal)
}

func (o *GroupImpl) AddIsNull(colName string) {
	o.whrCols = append(o.whrCols, colName)
	o.whrColVals = append(o.whrColVals, nil)
}

func (o *GroupImpl) AddIsNotNull(colName string) {
	o.whrCols = append(o.whrCols, colName)
	o.whrColVals = append(o.whrColVals, nil)
}

func (o *GroupImpl) AddBetween(colName string, from, to any) {
	o.whrBetweenCols = append(o.whrBetweenCols, colName)
	o.whrBetweenColVals = append(o.whrBetweenColVals, from, to)
}

func (o *GroupImpl) AddGroup(group Group) {
	o.groups = append(o.groups, group)
}

func (o *GroupImpl) gen(sb *strings.Builder, params *[]any) {
	fuse := fmt.Sprintf(" %s ", o.fuseType)

	sb.WriteString("(")

	fuseNeeded := false

	if len(o.whrColVals) > 0 {
		for i, col := range o.whrCols {
			if i > 0 {
				sb.WriteString(fuse)
			}
			sb.WriteString(col)
			sb.WriteString(" = ?")
		}
		*params = append(*params, o.whrColVals...)
		fuseNeeded = true
	}

	if len(o.whrBetweenCols) > 0 {
		if fuseNeeded {
			sb.WriteString(fuse)
		}
		for i, col := range o.whrBetweenCols {
			if i > 0 {
				sb.WriteString(fuse)
			}
			sb.WriteString(col)
			sb.WriteString(" BETWEEN ? AND ?")
		}
		*params = append(*params, o.whrBetweenColVals...)
		fuseNeeded = true
	}

	if len(o.whrIsNullCols) > 0 {
		if fuseNeeded {
			sb.WriteString(fuse)
		}
		for i, col := range o.whrIsNullCols {
			if i > 0 {
				sb.WriteString(fuse)
			}
			sb.WriteString(col)
			sb.WriteString(" IS NULL")
		}
		fuseNeeded = true
	}

	if len(o.whrIsNotNullCols) > 0 {
		if fuseNeeded {
			sb.WriteString(fuse)
		}
		for i, col := range o.whrIsNullCols {
			if i > 0 {
				sb.WriteString(fuse)
			}
			sb.WriteString(col)
			sb.WriteString(" IS NULL")
		}
		fuseNeeded = true
	}

	if len(o.groups) != 0 {
		if fuseNeeded {
			sb.WriteString(fuse)
		}
		for i, group := range o.groups {
			if i > 0 {
				sb.WriteString(fuse)
			}
			group.gen(sb, params)
		}
		fuseNeeded = true
	}

	sb.WriteString(")")
}

type WhereBuilderImpl struct {
	fuseType FuseType

	whrRawCols []string

	whrCols    []string
	whrColVals []any

	whrBetweenCols    []string
	whrBetweenColVals []any

	whrIsNullCols    []string
	whrIsNotNullCols []string

	groups []Group

	ordrByCols []string
	ordrByDesc bool

	limitRowCount  *int32
	limitRowOffset *int32
}

type WhereBuilder interface {
	Wheres
	AddOrdrByCols(colNames ...string)
	SetOrdrByDesc()
	SetLimitOffset(rowCount, rowOffset int32)
	SetLimit(rowCount int32)
	SetPage(pageNumber, pageSize int32)
	Gen() (string, []any)
}

func NewWhereBuilder(ft FuseType) WhereBuilder {
	o := WhereBuilderImpl{}

	o.fuseType = ft

	o.whrRawCols = []string{}

	o.whrCols = []string{}
	o.whrColVals = []any{}

	o.whrBetweenCols = []string{}
	o.whrBetweenColVals = []any{}

	o.whrIsNullCols = []string{}
	o.whrIsNotNullCols = []string{}

	o.groups = []Group{}

	o.ordrByCols = []string{}

	return &o
}

func (o *WhereBuilderImpl) SetFuse(ft FuseType) {
	o.fuseType = ft
}

func (o *WhereBuilderImpl) AddRaw(conditions ...string) {
	o.whrRawCols = append(o.whrRawCols, conditions...)
}

func (o *WhereBuilderImpl) Add(colName string, colVal any) {
	o.whrCols = append(o.whrCols, colName)
	o.whrColVals = append(o.whrColVals, colVal)
}

func (o *WhereBuilderImpl) AddBetween(colName string, from, to any) {
	o.whrBetweenCols = append(o.whrBetweenCols, colName)
	o.whrBetweenColVals = append(o.whrBetweenColVals, from)
	o.whrBetweenColVals = append(o.whrBetweenColVals, to)
}

func (o *WhereBuilderImpl) AddIsNull(colName string) {
	o.whrIsNullCols = append(o.whrIsNullCols, colName)
}

func (o *WhereBuilderImpl) AddIsNotNull(colName string) {
	o.whrIsNotNullCols = append(o.whrIsNotNullCols, colName)
}

func (o *WhereBuilderImpl) AddGroup(group Group) {
	o.groups = append(o.groups, group)
}

func (o *WhereBuilderImpl) AddOrdrByCols(colNames ...string) {
	o.ordrByCols = append(o.ordrByCols, colNames...)
}

func (o *WhereBuilderImpl) SetOrdrByDesc() {
	o.ordrByDesc = true
}

func (o *WhereBuilderImpl) SetLimitOffset(rowCount, rowOffset int32) {
	o.limitRowCount = &rowCount
	o.limitRowOffset = &rowOffset
}

func (o *WhereBuilderImpl) SetLimit(rowCount int32) {
	o.limitRowCount = &rowCount
}

func (o *WhereBuilderImpl) SetPage(pageNumber, pageSize int32) {
	rowCount := pageSize
	rowOffset := pageSize * pageNumber
	o.limitRowCount = &rowCount
	o.limitRowOffset = &rowOffset
}

func (o *WhereBuilderImpl) GenLimit() (sqlPartial string) {
	sb := strings.Builder{}
	if o.limitRowCount != nil {
		sb.WriteString(" LIMIT ")
		sb.WriteString(fmt.Sprintf("%d", *o.limitRowCount))

		if o.limitRowOffset != nil {
			sb.WriteString(" OFFSET ")
			sb.WriteString(fmt.Sprintf("%d", *o.limitRowOffset))
		}
	}
	sqlPartial = sb.String()
	return
}

func (o *WhereBuilderImpl) gen() (sqlPartial string, vals []any) {
	vals = []any{}

	sb := strings.Builder{}

	fuse := fmt.Sprintf(" %s ", o.fuseType)

	whereWritten := false

	// Raw columns:
	if len(o.whrRawCols) != 0 {
		sb.WriteString(" WHERE ")
		whereWritten = true
		sb.WriteString(strings.Join(o.whrRawCols, fuse))
	}

	// Where conditions:
	if len(o.whrCols) != 0 {
		if whereWritten {
			sb.WriteString(fuse)
		} else {
			sb.WriteString(" WHERE ")
			whereWritten = true
		}

		sb.WriteString(strings.Join(o.whrCols, "=?"+fuse))
		sb.WriteString("=?")
		vals = append(vals, o.whrColVals...)
	}

	if len(o.whrBetweenCols) != 0 {
		if whereWritten {
			sb.WriteString(fuse)
		} else {
			sb.WriteString(" WHERE ")
			whereWritten = true
		}
		for i, v := range o.whrBetweenCols {
			if i != 0 {
				sb.WriteString(fuse)
			}
			sb.WriteString(v)
			sb.WriteString(" BETWEEN ? AND ?")
		}
		vals = append(vals, o.whrBetweenColVals...)
	}

	if len(o.whrIsNullCols) != 0 {
		if whereWritten {
			sb.WriteString(fuse)
		} else {
			sb.WriteString(" WHERE ")
			whereWritten = true
		}
		for i, v := range o.whrIsNullCols {
			if i != 0 {
				sb.WriteString(fuse)
			}
			sb.WriteString(v)
			sb.WriteString(" IS NULL")
		}
	}

	if len(o.whrIsNotNullCols) != 0 {
		if whereWritten {
			sb.WriteString(fuse)
		} else {
			sb.WriteString(" WHERE ")
			whereWritten = true
		}
		for i, v := range o.whrIsNotNullCols {
			if i != 0 {
				sb.WriteString(fuse)
			}
			sb.WriteString(v)
			sb.WriteString(" IS NOT NULL")
		}
	}

	if len(o.groups) != 0 {
		if whereWritten {
			sb.WriteString(fuse)
		} else {
			sb.WriteString(" WHERE ")
			whereWritten = true
		}
		for i, group := range o.groups {
			if i != 0 {
				sb.WriteString(fuse)
			}
			group.gen(&sb, &vals)
		}
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

func (o *WhereBuilderImpl) Gen() (sqlPartial string, vals []any) {
	return o.gen()
}
