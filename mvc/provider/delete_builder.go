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

var _ SQLBuilder = (*DeleteBuilder)(nil)

// Create a DeleteBuilder instance to build a query string.
func NewDelete(table string) *DeleteBuilder {
	return &DeleteBuilder{table: table}
}

// Specify the target table for query.
func (b *DeleteBuilder) Table(table string) *DeleteBuilder {
	b.table = table
	return b
}

// Specify the where conditions and args for query.
func (b *DeleteBuilder) Wheres(where Wheres) *DeleteBuilder {
	b.wheres = where
	return b
}

// Specify the where in condition with field and args for query.
func (b *DeleteBuilder) WhereIn(field string, args []any) *DeleteBuilder {
	b.ins = b.FormatWhereIn(field, args)
	return b
}

// Specify the like condition for query.
func (b *DeleteBuilder) Like(field, filter string) *DeleteBuilder {
	b.like = b.FormatLike(field, filter)
	return b
}

// Specify the limit result for query.
func (b *DeleteBuilder) Limit(limit int) *DeleteBuilder {
	b.limit = limit
	return b
}

// Build and output query string and args for DataProvider execute delete action.
func (b *DeleteBuilder) Build() (string, []any) {
	where, args := b.BuildWheres(b.wheres, b.ins, b.like) // WHERE wheres AND field IN (v1,v2...) AND field2 LIKE '%%filter%%'
	limit := b.FormatLimit(b.limit)                       // LIMIT n

	query := "DELETE FROM %s %s %s"
	query = fmt.Sprintf(query, b.table, where, limit)
	query = strings.TrimSuffix(query, " ")
	return query, args
}

// Reset builder datas for next prepare and build.
func (b *DeleteBuilder) Reset() {
	clear(b.wheres)
	b.ins, b.like = "", ""
	b.limit = 0
}
