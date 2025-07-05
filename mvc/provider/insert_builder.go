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

	table string    // Table name for insert
	rows  []KValues // Target row records to insert.
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
//
//	- Set signle one row for provider.Insert().
//	- Set multiple rows for provider.Inserts().
func (b *InsertBuilder) Values(row ...KValues) *InsertBuilder {
	b.rows = row
	return b
}

// Build and output query string and args for DataProvider execute insert action.
func (b *InsertBuilder) Build() (string, []any) {
	if cnt := len(b.rows); cnt == 1 {
		// INSERT table (v1, v2...) VALUES (?,?...)'
		fields, holders, args := b.FormatInserts(b.rows[0])

		query := "INSERT %s (%s) VALUES (%s)"
		query = fmt.Sprintf(query, b.table, fields, holders)
		return query, args
	} else if cnt > 1 {
		// INSERT table (v1, v2...) VALUES (1,2...),(3,4...)...'
		headers := []string{}
		for key, value := range b.rows[0] { //fetch headers
			if key != "" && value != nil {
				headers = append(headers, key)
			}
		}

		rows := []string{}
		for _, row := range b.rows { //fetch rows
			vs := []string{}
			for _, h := range headers { // fetch colmuns
				if value, ok := row[h]; ok {
					switch v := value.(type) {
					case string:
						vs = append(vs, "'"+v+"'")
					default:
						vs = append(vs, fmt.Sprintf("%v", v))
					}
				}
			}
			// append row values: (1,'2',3.45,true,...)
			rows = append(rows, "("+strings.Join(vs, ",")+")")
		}

		fields := strings.Join(headers, ", ")
		values := strings.Join(rows, ", ")

		query := "INSERT %s (%s) VALUES %s"
		query = fmt.Sprintf(query, b.table, fields, values)
		return query, nil
	}
	return "", nil
}

// Reset builder datas for next prepare and build.
func (b *InsertBuilder) Reset() {
	clear(b.rows)
}
