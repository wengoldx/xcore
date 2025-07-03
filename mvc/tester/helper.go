// Copyright (c) 2018-Now Dunyu All Rights Reserved.
//
// Author      : https://www.wengold.net
// Email       : support@wengold.net
//
// Prismy.No | Date       | Modified by. | Description
// -------------------------------------------------------------------
// 00001       2019/05/22   yangping       New version
// -------------------------------------------------------------------

package tester

import (
	"database/sql"
	"fmt"

	pd "github.com/wengoldx/xcore/mvc/provider"
	"github.com/wengoldx/xcore/mvc/provider/mysqlc"
)

// Database helper for unit test.
type helper struct {
	pd.BaseProvider

	tag   string    // target output field name, such as id, uid...
	where pd.Wheres // where conditions field-value pairs, default empty.
	order string    // order by field name, auto append 'ORDER BY' perfix.
	desc  string    // order type, default DESC, one of 'DESC', 'ASC'.
	limit int       // limit number, auto append 'LIMIT' perfix.
}

// Create a data helper to query table datas.
func NewHelper() *helper {
	return &helper{
		BaseProvider: mysqlc.GetProvider(),
		desc:         "DESC", limit: 0,
	}
}

// Format order by condition to string.
func (t *helper) formatOrder() string {
	if t.order != "" {
		return fmt.Sprintf("ORDER BY %s %s", t.order, t.desc)
	}
	return ""
}

// Format limit condition to string.
func (t *helper) formatLimit() string {
	if t.limit > 0 {
		return fmt.Sprintf("LIMIT %d", t.limit)
	}
	return ""
}

// Format multiple value as IN condition in where.
func (t *helper) formatWhereIns(tag string, ins []string) string {
	if tag != "" && len(ins) > 0 {
		return tag + " " + t.Builder.JoinStrings(ins, "IN (%s)")
	}
	return ""
}

const (
	_sql_ut_get = "SELECT %s FROM %s %s %s %s" // table, where, order, limit
	_sql_ut_add = "INSERT %s (%s) VALUE (%s)"  // table, target fields, (?,?,...)
	_sql_ut_del = "DELETE FROM %s %s %s"       // table, where, in (%s)
)

// Get target field string value from given table, or with options.
//
//	SQL: SELECT %s FROM %s WHERE %s ORDER BY %s DESC LIMIT 1
//	            ^       ^  -------^ ----------^----- ------^
//	          tag   table     where       order        limit
//
//	@param table   Target table name.
//	@param tag     Target filed name to output query result.
//	@param where   Where conditions, the map key must like 'created>=?'.
//	@param options Setter for set order by field name, limit number.
//	@return out - any Output result value, like &int, &int64, &float64, &string...
func (t *helper) Target(table, tag string, where pd.Wheres, out any, options ...Option) error {
	t.tag, t.where = tag, where
	applyOptions(t, options...)

	wheres, values := t.FormatWheres(t.where) // format wheres sting and input values.
	order := t.formatOrder()                  // format order by string.
	limit := t.formatLimit()                  // format limit string.

	// SELECT tag FROM table wheres order limit
	query := fmt.Sprintf(_sql_ut_get, t.tag, table, wheres, order, limit)
	return t.One(query, func(rows *sql.Rows) error {
		if err := rows.Scan(out); err != nil {
			return err
		}
		return nil
	}, values...)
}

// Get last id from given table, or with options.
//
//	USAGE:
//
//	(1). mvc.UTest().LastID("account")
//	-> SELECT id FROM account ORDER BY id DESC LIMIT 1
//
//	(2). mvc.UTest().LastID("account", mvc.WithID("userid"))
//	-> SELECT userid FROM account ORDER BY id DESC LIMIT 1
//
//	(3). mvc.UTest().LastID("account", mvc.WithWhere({"acc=?" : "nickname"}))
//	-> SELECT id FROM account WHERE acc='nickname' ORDER BY id DESC LIMIT 1
//
//	(4). mvc.UTest().LastID("account", mvc.WithOrder("created"))
//	-> SELECT id FROM account ORDER BY created DESC LIMIT 1
//
//	@param table   Target table name.
//	@param options Setter for set id and order field name, where conditions, limit number.
//	@return id - Last record id.
//
//	See Target() for more sql query format infos.
func (t *helper) LastID(table string, options ...Option) (id int64, e error) {
	t.tag, t.order, t.limit = "id", "id", 1
	return id, t.Target(table, t.tag, t.where, &id, options...)
}

// Get last uid from given table, or with options.
//
//	See Target(), LastID for more query format or usage infos.
func (t *helper) LastUID(table string, options ...Option) (uid string, e error) {
	t.tag, t.order, t.limit = "uid", "created", 1
	return uid, t.Target(table, t.tag, t.where, &uid, options...)
}

// Insert a record into target table by given values.
//
//	SQL: INSERT %s (%s) VALUE (%s)
//	             ^   ^          ^
//	         table   tags    args
//
//	@param table Target table name.
//	@param ins   Target fields name and insert values.
func (t *helper) Add(table string, ins pd.KValues) error {
	fields, args, values := t.Builder.FormatInserts(ins)
	return t.Execute(fmt.Sprintf(_sql_ut_add, table, fields, args), values...)
}

// Insert a record into target table and return the inserted id.
//
//	See Add() for insert without record id.
func (t *helper) AddWithID(table string, ins pd.KValues) (int64, error) {
	fields, args, values := t.Builder.FormatInserts(ins)
	return t.Insert(fmt.Sprintf(_sql_ut_add, table, fields, args), values...)
}

// Deleta records by given where conditions.
//
//	SQL: DELETE FROM %s WHERE %s
//	                  ^ -------^
//	              table    where
//
//	@param table Target table name.
//	@param where Field name as where condition like: field = value.
func (t *helper) DeleteBy(table string, where pd.Wheres) error {
	t.where = where
	wheres, values := t.FormatWheres(t.where)
	return t.Execute(fmt.Sprintf(_sql_ut_del, table, wheres, ""), values...)
}

// Deleta records by given where condition and target fields vlaues.
//
//	SQL: DELETE FROM %s WHERE %s AND %s IN (%s)
//	                  ^ -------^---- ^-------^-
//	              table    where     tag     values
//
//	@param table Target table name.
//	@param where Field name as where condition like: field IN (values).
//	@param value Where condition values to query.
func (t *helper) DeleteIns(table, tag string, ins []string, options ...Option) error {
	applyOptions(t, options...)
	wheres, values := t.FormatWheres(t.where) // format wheres sting and input values.
	instring := t.formatWhereIns(tag, ins)    // format in values.
	if wheres != "" && instring != "" {
		wheres += " AND "
	}
	return t.Execute(fmt.Sprintf(_sql_ut_del, table, wheres, instring), values...)
}
