// Copyright (c) 2018-Now Dunyu All Rights Reserved.
//
// Author      : https://www.wengold.net
// Email       : support@wengold.net
//
// Prismy.No | Date       | Modified by. | Description
// -------------------------------------------------------------------
// 00001       2019/05/22   yangping       New version
// -------------------------------------------------------------------

package pd

import (
	"fmt"
	"strings"
)

// Build a query string for sql insert.
//
//	INSERT table (tags) VALUES (?, ?, ?)...
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

/* ------------------------------------------------------------------- */
/* SQL Action Utils By Using master Provider                           */
/* ------------------------------------------------------------------- */

func (b *InsertBuilder) Exec() error            { return b.master.Exec(b) }          // Insert a new record without check.
func (b *InsertBuilder) Insert() (int64, error) { return b.master.Insert(b) }        // Insert records and check inserted row id or counts.
func (b *InsertBuilder) InsertUncheck() error   { return b.master.InsertUncheck(b) } // Insert records without result check.

/* ------------------------------------------------------------------- */
/* SQL Action Builder Methonds                                         */
/* ------------------------------------------------------------------- */

// Specify master provider.
func (b *InsertBuilder) Master(master *TableProvider) *InsertBuilder {
	b.master = master
	return b
}

// Specify the target table for query.
func (b *InsertBuilder) Table(table string) *InsertBuilder {
	b.table = table
	return b
}

// Specify the values of row to insert.
//
// 1. Set signle one row for provider.Insert() insert a record.
//
//	row := KValues{
//		"":       123456,   // Filter out empty field
//		"Age":    16,
//		"Male":   true,
//		"Name":   "ZhangSan",
//		"Height": 176.8,
//		"Secure": nil,      // Filter out nil value
//	}
//	// => Age=?, Male=?, Name=?, Height=?
//	// => ?,?,?,?
//	// => []any{16, true, "ZhangSan", 176.8}
//
// 2. Set multiple rows for provider.Inserts() to insert records at one time.
//
//	rows := []KValues{
//		{ "Age": 16, "Name": "ZhangSan", "Height": 176.8 },
//		{ "Age": 15, "Name": "LiXu", "Height": 168.5 },
//	}
//	// => (Age=16,Name="ZhangSan",Height=176.8),(Age=15,Name="LiXu",Height=168.5)
func (b *InsertBuilder) Values(row ...KValues) *InsertBuilder {
	b.rows = row
	return b
}

// Build the insert action sql string and args for provider to insert datas.
//
//	INSERT table (tags) VALUES (?, ?, ?)...
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
func (b *InsertBuilder) Reset() *InsertBuilder {
	clear(b.rows)
	return b
}
