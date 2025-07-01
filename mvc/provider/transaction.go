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
	"database/sql"
	"strings"

	"github.com/wengoldx/xcore/invar"
)

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
	} else if rid, err := rst.LastInsertId(); err != nil {
		return err
	} else if out != nil {
		*out = rid
		return nil
	}
	return invar.ErrNotInserted
}

// Excute transaction step to insert multiple records.
//
//	query := "INSERT sametable (field1, field2) VALUES"
//	err := mvc.TxInserts(tx, query, len(vs), func(index int) string {
//		return fmt.Sprintf("(%v, %v)", v1, vs[index])
//		// return fmt.Sprintf("('%s', '%s')", v1, vs[index])
//	})
func TxInserts(tx *sql.Tx, query string, cnt int, cb InsertCallback) error {
	values := []string{}
	for i := 0; i < cnt; i++ {
		value := strings.TrimSpace(cb(i))
		if value != "" {
			values = append(values, value)
		}
	}
	query = query + " " + strings.Join(values, ",")
	_, err := tx.Exec(query)
	return err
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
