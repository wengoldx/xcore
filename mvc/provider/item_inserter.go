// Copyright (c) 2018-Now Dunyu All Rights Reserved.
//
// Author      : https://www.wengold.net
// Email       : support@wengold.net
//
// Prismy.No | Date       | Modified by. | Description
// -------------------------------------------------------------------
// 00001       2026/04/18   yangping       New version
// -------------------------------------------------------------------

package pd

import (
	"database/sql"
)

// A interface implement for insert multiple rows into target table,
// and only for TableProvider.Trans() method.
//
// Call NewInserter() to create a inserter to execute inserts.
type Inserter interface {
	DoInserts(tx *sql.Tx, query string) error
}

// A callback for format insert rows as string to insert record.
type InsertsCallback[T any] func(iv T) string

// SQL transaction inserter for cache insert datas.
type TxInserter [T any] struct {
	rows []T				// Multiple rows datas/
	cb   InsertsCallback[T] // Row datas format string callback.
}

var _ Inserter = (*TxInserter[any])(nil)

// Create a transaction inserter to insert multiple rows datas.
func NewInserter[T any](rows []T, cb InsertsCallback[T]) *TxInserter[T] {
	return &TxInserter[T]{rows:rows, cb:cb}
}

// Excute transaction step to insert multiple records.
//
//	err := h.Trans(
//		func (t *pd.Traner) error { return t.Inserts(query, pd.NewInserter(datas, func(iv *MyStruct) string {
//			return fmt.Sprintf("(%v, '%v')", iv.D1, iv.D2)
//		}, ...)
func (i *TxInserter[T]) DoInserts(tx *sql.Tx, query string) error {
	return TxInserts(tx, query, i.rows, i.cb)
}

