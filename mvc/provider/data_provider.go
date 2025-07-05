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

// A callback for scan query result records, it will interrupt
// scanning when callback return error.
type ScanCallback func(rows *sql.Rows) error

// A callback for format insert values as string to insert record.
type InsertCallback func(index int) string

// A callback for handle transaction by call provider.Trans().
type TransCallback func(tx *sql.Tx) error

// A interface for data provider export util methods.
type DataProvider interface {
	Has(query string, args ...any) (bool, error)
	Count(query string, args ...any) (int, error)
	Exec(query string, args ...any) error
	ExecResult(query string, args ...any) (int64, error)
	One(query string, cb ScanCallback, args ...any) error
	Query(query string, cb ScanCallback, args ...any) error
	Insert(query string, args ...any) (int64, error)
	Inserts(query string, cnt int, cb InsertCallback) error
	Update(query string, args ...any) error
	Delete(query string, args ...any) error
	Clear(table string) error
	Tran(query string, args ...any) error
	Trans(cbs ...TransCallback) error

	Affected(result sql.Result) (int64, error)
	Affects(result sql.Result) int64
	LastID(result sql.Result) int64

	MysqlTable(table string, print ...bool) *Table
	MssqlTable(table string, print ...bool) *Table
	PrintTable(table *Table)
}
