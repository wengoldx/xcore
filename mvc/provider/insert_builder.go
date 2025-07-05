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

// Build a query string for sql insert.
//
//	INSERT table (tags) VALUES (?, ?, ?)
//
// See QueryBuilder, UpdateBuilder, DeleteBuilder.
type InsertBuilder struct {
	BaseBuilder

	table  string  // Table name for insert
	values KValues // Target fields and value to insert.
}

var _ SQLBuilder = (*InsertBuilder)(nil)

// Create a InsertBuilder instance to build a query string.
func NewInsert(table string) *InsertBuilder {
	return &InsertBuilder{table: table}
}

// Specify the target table for query.
func (b *InsertBuilder) Table(table string) *InsertBuilder {
	b.table = table
	return b
}

// Specify the values of row to insert.
func (b *InsertBuilder) Values(row KValues) *InsertBuilder {
	b.values = row
	return b
}

// Build and output query string and args for DataProvider execute insert action.
func (b *InsertBuilder) Build() (string, []any) {
	fields, holders, args := b.FormatInserts(b.values) // INSERT table (v1, v2...) VALUES (?,?...)'

	query := "INSERT %s (%s) VALUES (%s)"
	query = fmt.Sprintf(query, b.table, fields, holders)
	query = strings.TrimSuffix(query, " ")
	return query, args
}

// Reset builder datas for next prepare and build.
func (b *InsertBuilder) Reset() {
	clear(b.values)
}
