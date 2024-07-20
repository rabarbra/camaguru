package orm

import (
	"fmt"
)

// Filter

type FilterOperation string

const (
	OpEqual        FilterOperation = "="
	OpNotEqual     FilterOperation = "<>"
	OpGreaterThan  FilterOperation = ">"
	OpLessThan     FilterOperation = "<"
	OpGreaterEqual FilterOperation = ">="
	OpLessEqual    FilterOperation = "<="
	OpLike         FilterOperation = "LIKE"
	Operation      FilterOperation = "IN"
)

type Filter struct {
	Key       string
	Value     any
	Operation FilterOperation
}

func (f Filter) ToSQL(filters ...Filter) string {
	sql := "WHERE "
	for i, item := range filters {
		sql += fmt.Sprintf("%s %s %s", item.Key, item.Operation, ToSQL(item.Value))
		if i < len(filters)-1 {
			sql += " AND "
		}
	}
	return sql
}

// Pagination

type Pagination struct {
	Limit  int
	Offset int
}

func (p Pagination) ToSQL() string {
	if p.Limit == 0 {
		return ""
	}
	return fmt.Sprintf("LIMIT %d OFFSET %d", p.Limit, p.Offset)
}

// Sorting

type OrderDirection string

const (
	ASC  OrderDirection = "ASC"
	DESC OrderDirection = "DESC"
)

type Sort struct {
	Key       string
	Direction OrderDirection
}

func (s Sort) ToSQL(sorts ...Sort) string {
	sql := "ORDER BY "
	for i, item := range sorts {
		sql += fmt.Sprintf("%s %s", item.Key, item.Direction)
		if i < len(sorts)-1 {
			sql += ", "
		}
	}
	return sql
}
