// Copyright (c) 2018-Now Dunyu All Rights Reserved.
//
// Author      : https://www.wengold.net
// Email       : support@wengold.net
//
// Prismy.No | Date       | Modified by. | Description
// -------------------------------------------------------------------
// 00001       2019/05/22   yangping       New version
// -------------------------------------------------------------------

package builder

import (
	"fmt"
	"strings"

	"github.com/wengoldx/xcore/logger"
	pd "github.com/wengoldx/xcore/mvc/provider"
	"github.com/wengoldx/xcore/utils"
)

// Build a query string for sql insert.
//
//	`MySQL & MSSQL`: INSERT table (tags) VALUES (?, ?, ?)...
//	`SQLITE`       : INSERT INTO table (tags) VALUES (?, ?, ?)...
type InsertBuilder struct {
	BaseBuilder
	rows []pd.KValues // Target row records to insert.
}

var _ pd.SQLBuilder = (*InsertBuilder)(nil)

// Create a InsertBuilder instance to build a query string.
func NewInsert(table string, provider ...pd.ProviderUtils) *InsertBuilder {
	return &InsertBuilder{BaseBuilder: *NewBuilder(table, provider...)}
}

/* ------------------------------------------------------------------- */
/* For Provider Insert Utils                                           */
/* ------------------------------------------------------------------- */

// Insert a new record without check.
//
//	h.Inserter().Values(pd.KValues{"role": admin, ...}).Exec()
func (b *InsertBuilder) Exec() error { return b.provider.Exec(b) }

// Insert records and check inserted row id or counts.
//
//	id, err := h.Inserter().Values(pd.KValues{"role": admin, ...}).Insert()
func (b *InsertBuilder) Insert() (int64, error) { return b.provider.Insert(b) }

// Insert records and check result.
//
//	h.Inserter().Values(pd.KValues{"role": admin, ...}).InsertCheck()
//	// check insert id or rows count.
func (b *InsertBuilder) InsertCheck() error { return b.provider.InsertCheck(b) }

// Insert records without result check.
//
//	h.Inserter().Values(pd.KValues{"role": admin, ...}).InsertUncheck()
//	// not check insert id and rows count.
func (b *InsertBuilder) InsertUncheck() error { return b.provider.InsertUncheck(b) }

/* ------------------------------------------------------------------- */
/* For SQL String Build Utils                                          */
/* ------------------------------------------------------------------- */

// Specify the values of row to insert.
//
// 1. Set signle one row for provider.Insert() insert a record.
//
//	row := pd.KValues{
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
//	rows := []pd.KValues{
//		{ "Age": 16, "Name": "ZhangSan", "Height": 176.8 },
//		{ "Age": 15, "Name": "LiXu", "Height": 168.5 },
//	}
//	// => (Age=16,Name="ZhangSan",Height=176.8),(Age=15,Name="LiXu",Height=168.5)
func (b *InsertBuilder) Values(row ...pd.KValues) *InsertBuilder {
	b.rows = row
	return b
}

// Reset builder datas for next prepare and build.
func (b *InsertBuilder) Reset() *InsertBuilder {
	clear(b.rows)
	return b
}

// Return rows count which insert to table later.
func (b *InsertBuilder) ValRows() int {
	return len(b.rows)
}

/* ------------------------------------------------------------------- */
/* For SQLBuilder interface                                            */
/* ------------------------------------------------------------------- */

// Build the insert action sql string and args for provider to insert datas.
//
//	`MySQL & MSSQL`: INSERT table (tags) VALUES (?, ?, ?)...
//	`SQLITE`       : INSERT INTO table (tags) VALUES (?, ?, ?)...
//
// # WARNING:
//
// The InsertBuild not well insert nil value by arg for single row
// insert, but good for insert nil value as NULL for multiple rows insert.
//
// And, it use the first row args key as the column headers.
func (b *InsertBuilder) Build(debug ...bool) (string, []any) {
	// diff := utils.Condition(b.driver == "sqlite", "INTO ", "")
	if cnt := len(b.rows); cnt == 1 {
		// INSERT INTO table (v1, v2...) VALUES (?,?...)'
		fields, holders, args := b.FormatInserts(b.rows[0])

		query := "INSERT INTO %s (%s) VALUES (%s)"
		query = fmt.Sprintf(query, b.table, fields, holders)
		if utils.Variable(debug, false) {
			logger.D("[INSERT] SQL:", query, "|", args)
		}
		return query, args
	} else if cnt > 1 {
		// INSERT [INTO] table (v1, v2...) VALUES (1,2...),(3,4...)...'
		headers := []string{}
		for key := range b.rows[0] {
			if key != "" { //fetch valid headers.
				headers = append(headers, key)
			}
		}

		rows := []string{}
		for _, row := range b.rows { //fetch rows
			vs := []string{}
			for _, h := range headers { // fetch colmuns
				if value, ok := row[h]; ok {
					// FIXME: Translate nil arg to NULL value
					// for multiple rows insert!
					if value == nil {
						vs = append(vs, "NULL")
						continue
					}

					switch v := value.(type) {
					case string:
						vs = append(vs, "'"+v+"'")
					default:
						vs = append(vs, fmt.Sprintf("%v", v))
					}
				}
			}
			// append row values: (1,'2',3.45,true,NULL,...)
			rows = append(rows, "("+strings.Join(vs, ",")+")")
		}

		fields := strings.Join(headers, ", ")
		values := strings.Join(rows, ", ")

		query := "INSERT INTO %s (%s) VALUES %s"
		query = fmt.Sprintf(query, b.table, fields, values)
		if utils.Variable(debug, false) {
			logger.D("[INSERT-S] SQL:", query)
		}
		return query, nil
	}
	return "", nil
}
