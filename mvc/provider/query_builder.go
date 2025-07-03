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

// Build a query string for sql query.
//
//	SELECT outs FROM table
//		WHERE wherer AND field IN (v1,v2...) AND field2 LIKE '%%filter%%'
//		ORDER BY order DESC
//		LIMIT limit.
//
// WARNING: This builder only for single table access.
type QueryBuilder struct {
	BaseBuilder

	table  string   // Table name for query
	outs   []string // Fields for output values.
	wheres Wheres   // Where conditions and values.
	ins    string   // Where in conditions.
	like   string   // Like conditions string.
	order  string   // Keyword for order by condition.
	limit  int      // Limit number.
}

// Create a QueryBuilder instance to build a query string.
func NewQuery(table string) *QueryBuilder {
	return &QueryBuilder{table: table}
}

// Special the target table for query.
func (q *QueryBuilder) Table(table string) *QueryBuilder {
	q.table = table
	return q
}

// Special the output fields for query.
func (q *QueryBuilder) Outs(field ...string) *QueryBuilder {
	q.outs = field
	return q
}

// Special the where conditions and args for query.
func (q *QueryBuilder) Wheres(where Wheres) *QueryBuilder {
	q.wheres = where
	return q
}

// Special the where in condition with field and args for query.
func (q *QueryBuilder) WhereIn(field string, args []any) *QueryBuilder {
	q.ins = q.FormatWhereIn(field, args)
	return q
}

// Special the order by condition for query.
func (q *QueryBuilder) OrderBy(field string, desc bool) *QueryBuilder {
	q.order = q.FormatOrder(field, desc)
	return q
}

// Special the like condition for query.
func (q *QueryBuilder) Like(field, filter string) *QueryBuilder {
	q.like = q.FormatLike(field, filter)
	return q
}

// Special the limit result for query.
func (q *QueryBuilder) Limit(limit int) *QueryBuilder {
	q.limit = limit
	return q
}

// Build and output query string and args for DataProvider execute query action.
func (q *QueryBuilder) Build() (string, []any) {
	outs := strings.Join(q.outs, ",")       // out1,out2,out3...
	where, args := q.FormatWheres(q.wheres) // WHERE wheres
	if where != "" {
		// WHERE wheres AND field IN (v1,v2...)
		if q.ins != "" {
			where += " AND " + q.ins
		}

		// WHERE wheres AND field IN (v1,v2...) AND field2 LIKE '%%filter%%'
		if q.like != "" {
			where += " AND " + q.like
		}
	} else {
		if q.ins != "" {
			// WHERE field IN (v1,v2...) AND field2 LIKE '%%filter%%'
			where = "WHERE " + q.ins
			if q.like != "" {
				where += " AND " + q.like
			}
		} else if q.like != "" {
			// WHERE field LIKE '%%filter%%'
			where = "WHERE " + q.like
		}
	}
	limit := q.FormatLimit(q.limit) // LIMIT n

	query := "SELECT %s FROM %s %s %s %s"
	query = fmt.Sprintf(query, outs, q.table, where, q.order, limit)
	query = strings.TrimSuffix(query, " ")
	return query, args
}

// Reset builder datas for next prepare and build.
func (b *QueryBuilder) Reset() {
	clear(b.outs)
	clear(b.wheres)
	b.ins, b.like, b.order = "", "", ""
	b.limit = 0
}
