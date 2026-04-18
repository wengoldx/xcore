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
	"strings"

	"github.com/wengoldx/xcore/invar"
)

/* ------------------------------------------------------------------- */
/* For TableProvider Trans() Method Utils                              */
/* ------------------------------------------------------------------- */

// SQL translate utils, alias from sql.Tx.
//
//	tx -> tr convert: `tr, ok := any(tx).(*Traner)` OR `tr = (*Traner)(tx)`
//	tr -> tx convert: `tx, ok := any(tr).(*sql.Tx)` OR `tx = (*sql.Tx)(tr)`
type Traner sql.Tx

// Excute transaction step to update, insert, or delete datas without check result.
func (t *Traner) Exec(query string, args ...any) error {
	return TxExec((*sql.Tx)(t), query, args...)
}

// Excute transaction step to check if data exist, it wil return
// invar.ErrNotFound if unexist any records, or return nil when exist results.
func (t *Traner) Exist(query string, args ...any) error {
	return TxExist((*sql.Tx)(t), query, args...)
}

// Excute transaction step to query single data and get result in scan callback.
func (t *Traner) One(query string, cb ScanCallback, args ...any) error {
	return TxOne((*sql.Tx)(t), query, cb, args...)
}

// Excute transaction step to query datas, and fetch result in scan callback.
func (t *Traner) Query(query string, cb ScanCallback, args ...any) error {
	return TxQuery((*sql.Tx)(t), query, cb, args...)
}

// Excute transaction step to insert a new record and return inserted id.
func (t *Traner) Insert(query string, out *int64, args ...any) error {
	return TxInsert((*sql.Tx)(t), query, out, args...)
}

// Excute transaction step to insert multiple records.
//
//	// type MyStruct {D1 int, D2 string}
//	// datas := []*MyStruct{{1,"2"}, {3,"4"}}
//	query := "INSERT sametable (field1, field2) VALUES"
//	err := h.Trans(
//		func (t *pd.Traner) error { return t.Query(query_str..., args...)},
//		func (t *pd.Traner) error { return t.Inserts(query, pd.NewInserter(datas, func(iv *MyStruct) string {
//			return fmt.Sprintf("(%v, '%v')", iv.D1, iv.D2)
//		})
func (t *Traner) Inserts(query string, inserter Inserter) error {
	return inserter.DoInserts((*sql.Tx)(t), query)
}

// Excute transaction step to delete record and check result,
// set 'out = nil' for not output inserted record id.
func (t *Traner) Delete(query string, out *int64, args ...any) error {
	return TxDelete((*sql.Tx)(t), query, args...)
}

/* ------------------------------------------------------------------- */
/* For Global Transaction Utils                                        */
/* ------------------------------------------------------------------- */

// Excute transaction step to update, insert, or delete datas without check result.
func TxExec(tx *sql.Tx, query string, args ...any) error {
	_, err := tx.Exec(query, args...)
	return err
}

// Excute transaction step to check if data exist, it wil return
// invar.ErrNotFound if unexist any records, or return nil when exist results.
func TxExist(tx *sql.Tx, query string, args ...any) error {
	if rows, err := tx.Query(query, args...); err != nil {
		return err
	} else {
		defer rows.Close()
		if !rows.Next() {
			return invar.ErrNotFound
		}
	}
	return nil
}

// Excute transaction step to query single data and get result in scan callback.
func TxOne(tx *sql.Tx, query string, cb ScanCallback, args ...any) error {
	if rows, err := tx.Query(query, args...); err != nil {
		return err
	} else {
		defer rows.Close()

		if !rows.Next() {
			return invar.ErrNotFound
		}
		rows.Columns()
		return cb(rows)
	}
}

// Excute transaction step to query datas, and fetch result in scan callback.
func TxQuery(tx *sql.Tx, query string, cb ScanCallback, args ...any) error {
	if rows, err := tx.Query(query, args...); err != nil {
		return err
	} else {
		defer rows.Close()

		for rows.Next() {
			rows.Columns()
			if err := cb(rows); err != nil {
				return err
			}
		}
	}
	return nil
}

// Excute transaction step to insert a new record and return inserted id.
func TxInsert(tx *sql.Tx, query string, out *int64, args ...any) error {
	if rst, err := tx.Exec(query, args...); err != nil {
		return err
	} else if out != nil {
		rid, err := rst.LastInsertId()
		if err != nil {
			return err
		} else if rid == 0 {
			return invar.ErrNotInserted
		}
		*out = rid
	}
	return nil
}

// Excute transaction step to insert multiple records.
//
//	// type MyStruct {D1 int, D2 string}
//	// datas := []*MyStruct{{1,"2"}, {3,"4"}}
//	query := "INSERT sametable (field1, field2) VALUES"
//	err := pd.TxInserts(tx, query, datas, func(iv *MyStruct) string {
//		return fmt.Sprintf("(%v, '%v')", iv.D1, iv.D2)
//	})
func TxInserts[T any](tx *sql.Tx, query string, rows []T, cb InsertsCallback[T]) error {
	if cnt := len(rows); cnt > 0 {
		values := []string{}
		for i := 0; i < cnt; i++ {
			if row := strings.TrimSpace(cb(rows[i])); row != "" {
				values = append(values, row)
			}
		}
		query = query + " " + strings.Join(values, ",")
		_, err := tx.Exec(query)
		return err
	}
	return nil
}

// Excute transaction step to delete record and check result.
func TxDelete(tx *sql.Tx, query string, args ...any) error {
	if rst, err := tx.Exec(query, args...); err != nil {
		return err
	} else if cnt, err := rst.RowsAffected(); err != nil {
		return err
	} else if cnt == 0 {
		return invar.ErrNotChanged
	}
	return nil
}
