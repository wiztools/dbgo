package dbgo

import (
	"fmt"
	"strconv"
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
	Add(colName string, colVal any)
	AddBetween(colName string, from, to any)
	AddGroup(group Group)
}

type GroupImpl struct {
	fuseType FuseType

	whrCols    []string
	whrColVals []any

	whrBetweenCols    []string
	whrBetweenColVals []any

	groups []Group
}

func NewGroup(fuseType FuseType) Group {
	return &GroupImpl{
		fuseType:          fuseType,
		whrCols:           []string{},
		whrColVals:        []any{},
		whrBetweenCols:    []string{},
		whrBetweenColVals: []any{},
		groups:            []Group{},
	}
}

func (o *GroupImpl) SetFuse(ft FuseType) {
	o.fuseType = ft
}

func (o *GroupImpl) Add(colName string, colVal any) {
	o.whrCols = append(o.whrCols, colName)
	o.whrColVals = append(o.whrColVals, colVal)
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

	for i, col := range o.whrCols {
		if i > 0 {
			sb.WriteString(fuse)
		}
		sb.WriteString(col)
		sb.WriteString(" = ?")
	}
	*params = append(*params, o.whrColVals...)

	for i, col := range o.whrBetweenCols {
		if i > 0 {
			sb.WriteString(fuse)
		}
		sb.WriteString(col)
		sb.WriteString(" BETWEEN ? AND ?")
	}
	*params = append(*params, o.whrBetweenColVals...)

	if len(o.groups) != 0 {
		if len(o.whrCols) > 0 || len(o.whrBetweenCols) > 0 {
			sb.WriteString(fuse)
		}
		for i, group := range o.groups {
			if i > 0 {
				sb.WriteString(fuse)
			}
			group.gen(sb, params)
		}
	}

	sb.WriteString(")")
}

type WhereBuilderImpl struct {
	fuseType FuseType

	whrCols    []string
	whrColVals []any

	whrBetweenCols    []string
	whrBetweenColVals []any

	groups []Group

	ordrByCols []string
	ordrByDesc bool

	limitRowCount  *int
	limitRowOffset *int
}

type WhereBuilder interface {
	Wheres
	AddOrdrByCols(colNames ...string)
	SetOrdrByDesc()
	SetLimitOffset(rowCount, rowOffset int)
	SetLimit(rowCount int)
	SetPage(pageNumber, pageSize int)
	Gen() (string, []any)
}

func NewWhereBuilder(ft FuseType) WhereBuilder {
	o := WhereBuilderImpl{}

	o.fuseType = ft

	o.whrCols = []string{}
	o.whrColVals = []any{}

	o.whrBetweenCols = []string{}
	o.whrBetweenColVals = []any{}

	o.groups = []Group{}

	o.ordrByCols = []string{}

	return &o
}

func (o *WhereBuilderImpl) SetFuse(ft FuseType) {
	o.fuseType = ft
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

func (o *WhereBuilderImpl) AddGroup(group Group) {
	o.groups = append(o.groups, group)
}

func (o *WhereBuilderImpl) AddOrdrByCols(colNames ...string) {
	o.ordrByCols = append(o.ordrByCols, colNames...)
}

func (o *WhereBuilderImpl) SetOrdrByDesc() {
	o.ordrByDesc = true
}

func (o *WhereBuilderImpl) SetLimitOffset(rowCount, rowOffset int) {
	o.limitRowCount = &rowCount
	o.limitRowOffset = &rowOffset
}

func (o *WhereBuilderImpl) SetLimit(rowCount int) {
	o.limitRowCount = &rowCount
}

func (o *WhereBuilderImpl) SetPage(pageNumber, pageSize int) {
	rowCount := pageSize
	rowOffset := pageSize * pageNumber
	o.limitRowCount = &rowCount
	o.limitRowOffset = &rowOffset
}

func (o *WhereBuilderImpl) GenLimit() (sqlPartial string) {
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

func (o *WhereBuilderImpl) gen() (sqlPartial string, vals []any) {
	vals = []any{}

	sb := strings.Builder{}

	fuse := fmt.Sprintf(" %s ", o.fuseType)

	// Where conditions:
	if len(o.whrCols) != 0 {
		sb.WriteString(" WHERE ")

		sb.WriteString(strings.Join(o.whrCols, "=?"+fuse))
		sb.WriteString("=?")
		vals = append(vals, o.whrColVals...)
	}

	if len(o.whrBetweenCols) != 0 {
		if len(o.whrCols) != 0 {
			sb.WriteString(fuse)
		} else {
			sb.WriteString(" WHERE ")
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

	if len(o.groups) != 0 {
		if len(o.whrCols) != 0 || len(o.whrBetweenCols) != 0 {
			sb.WriteString(fuse)
		} else {
			sb.WriteString(" WHERE ")
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
