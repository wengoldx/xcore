// Copyright (c) 2018-Now Dunyu All Rights Reserved.
//
// Author      : https://www.wengold.net
// Email       : support@wengold.net
//
// Prismy.No | Date       | Modified by. | Description
// -------------------------------------------------------------------
// 00001       2019/05/22   yangping       New version
// -------------------------------------------------------------------

package provider

import (
	"fmt"
	"strings"
)

type QueryBuilder struct {
	BaseBuilder

	table  string   // Table name for query
	outs   []string // Fields for output values.
	wheres Wheres   // Where conditions and values.
	order  string   // Keyword for order by condition.
	limit  int      // Limit number.
}

func NewQuery(table string) *QueryBuilder {
	return &QueryBuilder{table: table}
}

func (q *QueryBuilder) Outs(field ...string) *QueryBuilder {
	q.outs = field
	return q
}

func (q *QueryBuilder) Wheres(where Wheres) *QueryBuilder {
	q.wheres = where
	return q
}

func (q *QueryBuilder) OrderBy(field string, desc bool) *QueryBuilder {
	q.order = q.FormatOrder(field, desc)
	return q
}

func (q *QueryBuilder) Limit(limit int) *QueryBuilder {
	q.limit = limit
	return q
}

func (q *QueryBuilder) Build() (string, []any) {
	outs := strings.Join(q.outs, ",")
	where, args := q.FormatWheres(q.wheres)
	limit := q.FormatLimit(q.limit)

	query := "SELECT %s FROM %s %s %s %s"
	query = fmt.Sprintf(query, outs, q.table, where, q.order, limit)
	return query, args
}
