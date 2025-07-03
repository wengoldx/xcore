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

// Build a query string for sql delete.
//
//	DELETE FROM table
//		WHERE wherer AND field IN (v1,v2...) AND field2 LIKE '%%filter%%'
//		LIMIT limit.
//
// WARNING: This builder only for single table access.
type DeleteBuilder struct {
	BaseBuilder

	table  string // Table name for query
	wheres Wheres // Where conditions and values.
	ins    string // Where in conditions.
	like   string // Like conditions string.
	limit  int    // Limit number.
}

// Create a DeleteBuilder instance to build a query string.
func NewDelete(table string) *DeleteBuilder {
	return &DeleteBuilder{table: table}
}

// Specify the target table for query.
func (q *DeleteBuilder) Table(table string) *DeleteBuilder {
	q.table = table
	return q
}

// Specify the where conditions and args for query.
func (q *DeleteBuilder) Wheres(where Wheres) *DeleteBuilder {
	q.wheres = where
	return q
}

// Specify the where in condition with field and args for query.
func (q *DeleteBuilder) WhereIn(field string, args []any) *DeleteBuilder {
	q.ins = q.FormatWhereIn(field, args)
	return q
}

// Specify the like condition for query.
func (q *DeleteBuilder) Like(field, filter string) *DeleteBuilder {
	q.like = q.FormatLike(field, filter)
	return q
}

// Specify the limit result for query.
func (q *DeleteBuilder) Limit(limit int) *DeleteBuilder {
	q.limit = limit
	return q
}

// Build and output query string and args for DataProvider execute delete action.
func (q *DeleteBuilder) Build() (string, []any) {
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

	query := "DELETE FROM %s %s %s"
	query = fmt.Sprintf(query, q.table, where, limit)
	query = strings.TrimSuffix(query, " ")
	return query, args
}

// Reset builder datas for next prepare and build.
func (b *DeleteBuilder) Reset() {
	clear(b.wheres)
	b.ins, b.like = "", ""
	b.limit = 0
}
