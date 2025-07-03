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

import "database/sql"

// ScanCallback use for scan query result from rows
type ScanCallback func(rows *sql.Rows) error

// InsertCallback format query values to string for Inserts().
type InsertCallback func(index int) string

// TransCallback transaction callback for Trans().
type TransCallback func(tx *sql.Tx) error

// A interface for data provider export util methods.
type DataProvider interface {
	Has(query string, args ...any) (bool, error)
	Count(query string, args ...any) (int, error)
	Exec(query string, args ...any) error
	Delete(query string, args ...any) error

	None(builder *QueryBuilder) (bool, error)
	Counts(builder *QueryBuilder) (int, error)
	Deletes(builder *DeleteBuilder) error

	One(query string, cb ScanCallback, args ...any) error
	Query(query string, cb ScanCallback, args ...any) error
	Insert(query string, args ...any) (int64, error)
	Inserts(query string, cnt int, cb InsertCallback) error
	Inserts2(query string, values any) error
	Update(query string, args ...any) error
	Update2(query string, values map[string]any, args ...any) error
	Execute2(query string, args ...any) (int64, error)
	TranRoll(query string, args ...any) error
	Trans(cbs ...TransCallback) error

	Affected(result sql.Result) (int64, error)
	Affects(result sql.Result) int64
	LastID(result sql.Result) int64
	FormatSets(values map[string]any) (string, error)
	FormatInserts(values any) (string, error)

	MysqlTable(table string, print ...bool) *Table
	MssqlTable(table string, print ...bool) *Table
	PrintTable(table *Table)
}
