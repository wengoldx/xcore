// Copyright (c) 2018-Now Dunyu All Rights Reserved.
//
// Author      : https://www.wengold.net
// Email       : support@wengold.net
//
// Prismy.No | Date       | Modified by. | Description
// -------------------------------------------------------------------
// 00001       2019/05/22   yangping       New version
// -------------------------------------------------------------------

package mvc

import (
	"database/sql"
	"fmt"
	"strings"
)

// Database helper for unit test.
type utestHelper struct {
	WingProvider
}

// Create a data helper to query table datas.
func UTest() *utestHelper {
	return &utestHelper{*WingHelper}
}

const (
	_sql_ut_last_id    = "SELECT id FROM %s %s ORDER BY %s DESC LIMIT 1"
	_sql_ut_last_ids   = "SELECT id FROM %s WHERE %s >= ?"
	_sql_ut_last_field = "SELECT %s FROM %s ORDER BY %s DESC LIMIT 1"
	_sql_ut_get_target = "SELECT %s FROM %s WHERE %s = ? ORDER BY id DESC LIMIT 1"
	_sql_ut_get_datas  = "SELECT %s FROM %s WHERE %s IN (%s)"
	_sql_ut_del_one    = "DELETE FROM %s WHERE %s = ?"
	_sql_ut_del_multis = "DELETE FROM %s WHERE %s IN (%s)"
	_sql_ut_clear      = "DELETE FROM %s %s"
)

// Get last id from given table, or with time confitions.
func (t *utestHelper) LastID(table, field string, times ...string) (id int64, e error) {
	where := ""
	if len(times) >= 2 {
		where = fmt.Sprintf("WHERE %s >= \"%s\"", times[0], times[1])
	}
	query := fmt.Sprintf(_sql_ut_last_id, table, where, field)

	return id, t.One(query, func(rows *sql.Rows) error {
		if e = rows.Scan(&id); e != nil {
			return e
		}
		return nil
	})
}

// Get last ids from given table and query time, or compare condition value.
func (t *utestHelper) LastIDs(table, field string, value any) ([]int64, error) {
	ids, query := []int64{}, fmt.Sprintf(_sql_ut_last_ids, table, field)
	return ids, t.Query(query, func(rows *sql.Rows) error {
		var id int64
		if err := rows.Scan(&(id)); err != nil {
			return err
		}

		ids = append(ids, id)
		return nil
	}, value)
}

// Query the target field value by the top most given order field.
func (t *utestHelper) LastField(table string, target string, order string) (v string, e error) {
	query := fmt.Sprintf(_sql_ut_last_field, target, table, order)
	return v, t.One(query, func(rows *sql.Rows) error {
		if e = rows.Scan(&v); e != nil {
			return e
		}
		return nil
	})
}

// Query the target field last value by given condition field and value.
func (t *utestHelper) Target(table, target, field, value string) (v string, e error) {
	query := fmt.Sprintf(_sql_ut_get_target, target, table, field)
	return v, t.One(query, func(rows *sql.Rows) error {
		if e = rows.Scan(&v); e != nil {
			return e
		}
		return nil
	}, value)
}

// Query the target field values by given condition field and values.
func (t *utestHelper) Datas(table string, target string, field string, values string) ([]string, error) {
	rsts, query := []string{}, fmt.Sprintf(_sql_ut_get_datas, target, table, field, values)
	return rsts, t.Query(query, func(rows *sql.Rows) error {
		rst := ""
		if err := rows.Scan(&rst); err != nil {
			return err
		}
		rsts = append(rsts, rst)
		return nil
	})
}

// Deleta records by target field on equal condition.
func (t *utestHelper) DelOne(table, field string, value any) {
	t.Execute(fmt.Sprintf(_sql_ut_del_one, table, field), value)
}

// Deleta records by target field on in range condition.
func (t *utestHelper) DelMults(table, field, values string) {
	t.Execute(fmt.Sprintf(_sql_ut_del_multis, table, field, values))
}

// Clear the target table all datas, or ranged datas of in given conditions.
func (t *utestHelper) Clear(table string, wheres ...string) {
	t.Execute(fmt.Sprintf(_sql_ut_clear, table, strings.Join(wheres, " ")))
}
