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

// Build a query string for sql update.
//
//	UPDATE table
//		SET v1=?, v2=?, v3=?...
//		WHERE wherer AND field IN (v1,v2...) AND field2 LIKE '%%filter%%'
//
// See QueryBuilder, InsertBuilder, DeleteBuilder.
type UpdateBuilder struct {
	BaseBuilder

	table  string  // Table name for update
	values KValues // Target fields and values to update.
	wheres Wheres  // Where conditions and args values.
	ins    string  // Where in conditions.
	like   string  // Like conditions string.
}

var _ SQLBuilder = (*UpdateBuilder)(nil)

// Create a UpdateBuilder instance to build a query string.
func NewUpdate(table string) *UpdateBuilder {
	return &UpdateBuilder{table: table}
}

// Specify the target table for query.
func (b *UpdateBuilder) Table(table string) *UpdateBuilder {
	b.table = table
	return b
}

// Specify the values of row to update.
func (b *UpdateBuilder) Values(row KValues) *UpdateBuilder {
	b.values = row
	return b
}

// Specify the where conditions and args for query.
func (b *UpdateBuilder) Wheres(where Wheres) *UpdateBuilder {
	b.wheres = where
	return b
}

// Specify the where in condition with field and args for query.
func (b *UpdateBuilder) WhereIn(field string, args []any) *UpdateBuilder {
	b.ins = b.FormatWhereIn(field, args)
	return b
}

// Specify the like condition for query.
func (b *UpdateBuilder) Like(field, filter string) *UpdateBuilder {
	b.like = b.FormatLike(field, filter)
	return b
}

// Build and output query string and args for DataProvider execute update action.
func (b *UpdateBuilder) Build() (string, []any) {
	tags, args := b.FormatSets(b.values)                 // SET v1=?,v2=?...
	where, wvs := b.BuildWheres(b.wheres, b.ins, b.like) // WHERE wheres AND field IN (v1,v2...) AND field2 LIKE '%%filter%%'
	args = append(args, wvs)

	query := "UPDATE %s SET %s %s"
	query = fmt.Sprintf(query, b.table, tags, where)
	query = strings.TrimSuffix(query, " ")
	return query, args
}

// Reset builder datas for next prepare and build.
func (b *UpdateBuilder) Reset() {
	clear(b.values)
	clear(b.wheres)
	b.ins, b.like = "", ""
}
