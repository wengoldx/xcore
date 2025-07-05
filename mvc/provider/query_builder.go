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
// See InsertBuilder, UpdateBuilder, DeleteBuilder.
type QueryBuilder struct {
	BaseBuilder

	table  string   // Table name for query
	tags   []string // Target fields for output values.
	wheres Wheres   // Where conditions and args values.
	ins    string   // Where in conditions.
	like   string   // Like conditions string.
	order  string   // Keyword for order by condition.
	limit  int      // Limit number.
}

var _ SQLBuilder = (*QueryBuilder)(nil)

// Create a QueryBuilder instance to build a query string.
func NewQuery(table string) *QueryBuilder {
	return &QueryBuilder{table: table}
}

// Specify the target table for query.
func (b *QueryBuilder) Table(table string) *QueryBuilder {
	b.table = table
	return b
}

// Specify the target output fields for query.
func (b *QueryBuilder) Tags(tag ...string) *QueryBuilder {
	b.tags = tag
	return b
}

// Specify the where conditions and args for query.
func (b *QueryBuilder) Wheres(where Wheres) *QueryBuilder {
	b.wheres = where
	return b
}

// Specify the where in condition with field and args for query.
func (b *QueryBuilder) WhereIn(field string, args []any) *QueryBuilder {
	b.ins = b.FormatWhereIn(field, args)
	return b
}

// Specify the order by condition for query.
func (b *QueryBuilder) OrderBy(field string, desc bool) *QueryBuilder {
	b.order = b.FormatOrder(field, desc)
	return b
}

// Specify the like condition for query.
func (b *QueryBuilder) Like(field, filter string) *QueryBuilder {
	b.like = b.FormatLike(field, filter)
	return b
}

// Specify the limit result for query.
func (b *QueryBuilder) Limit(limit int) *QueryBuilder {
	b.limit = limit
	return b
}

// Build and output query string and args for DataProvider execute query action.
func (b *QueryBuilder) Build() (string, []any) {
	tags := strings.Join(b.tags, ",")                     // out1,out2,out3...
	where, args := b.BuildWheres(b.wheres, b.ins, b.like) // WHERE wheres AND field IN (v1,v2...) AND field2 LIKE '%%filter%%'
	limit := b.FormatLimit(b.limit)                       // LIMIT n

	query := "SELECT %s FROM %s %s %s %s"
	query = fmt.Sprintf(query, tags, b.table, where, b.order, limit)
	query = strings.TrimSuffix(query, " ")
	return query, args
}

// Reset builder datas for next prepare and build.
func (b *QueryBuilder) Reset() {
	clear(b.tags)
	clear(b.wheres)
	b.ins, b.like, b.order = "", "", ""
	b.limit = 0
}
