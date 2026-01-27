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
	"fmt"
	"strings"

	"github.com/wengoldx/xcore/invar"
	"github.com/wengoldx/xcore/logger"
	pd "github.com/wengoldx/xcore/mvc/provider"
	"github.com/wengoldx/xcore/mvc/provider/builder"
)

// Base provider for simple access database datas.
type BaseProvider struct {
	client pd.DBClient // Database conncet client.

	/* Only for BaseProvider as build utils! */
	Builder pd.BaseBuilder // Base builder as utils tools.
}

// Create a BaseProvider with given database client.
func NewBaseProvider(client pd.DBClient) *BaseProvider {
	// FIXME: the client maybe nil!
	return &BaseProvider{client, &builder.BuilderImpl{}}
}

var _ pd.Provider = (*BaseProvider)(nil)

// Set provider database client.
func (p *BaseProvider) SetClient(client pd.DBClient) {
	if client == nil {
		logger.E("@@ DBClient is nil!")
	}
	p.client = client
}

/* ------------------------------------------------------------------- */
/* Direct Use Query String To Access Database                          */
/* ------------------------------------------------------------------- */

// Execute query string to check target record whether exist, it will
// auto append 'LIMIT 1' at query tail if not specified.
//
// Use the QueryBuilder to build a query string and args.
func (p *BaseProvider) Has(query string, args ...any) (bool, error) {
	if !p.prepared() || query == "" {
		return false, invar.ErrBadDBConnect
	}

	query = p.Builder.CheckLimit(query)
	rows, err := p.client.DB().Query(query, args...)
	if err != nil {
		return false, err
	}
	defer rows.Close()
	return rows.Next(), nil
}

// Execute query string to count records on given conditions, it will
// return 0 when notfound anyone.
//
//	Use the QueryBuilder to build a query string and args.
func (p *BaseProvider) Count(query string, args ...any) (int, error) {
	if !p.prepared() || query == "" {
		return 0, invar.ErrBadDBConnect
	}

	rows, err := p.client.DB().Query(query, args...)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	counts := 0
	if rows.Next() {
		rows.Columns()
		if err := rows.Scan(&counts); err != nil {
			return 0, err
		}
	}
	return counts, nil
}

// Execute query string without any results to get, it useful for
// update or delete record datas.
//
//	Use UpdateBuilder or DeleteBuilder to build a query string and args.
func (p *BaseProvider) Exec(query string, args ...any) error {
	if !p.prepared() || query == "" {
		return invar.ErrBadDBConnect
	}

	stmt, err := p.client.DB().Prepare(query)
	if err != nil {
		return err
	}

	defer stmt.Close()
	if _, err := stmt.Exec(args...); err != nil {
		return err
	}
	return nil
}

// Execute query string and return the affected rows count, it will return
// invar.ErrNotChanged error when none updated or deleted.
//
//	Use UpdateBuilder or DeleteBuilder to build a query string and args.
func (p *BaseProvider) ExecResult(query string, args ...any) (int64, error) {
	if !p.prepared() || query == "" {
		return 0, invar.ErrBadDBConnect
	}

	stmt, err := p.client.DB().Prepare(query)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(args...)
	if err != nil {
		return 0, err
	}
	return p.Affected(result)
}

// Execute query string to get the top one record, it will auto append
// 'LIMIT 1' as tail in query string for high-performance.
//
//	Use QueryBuilder to build a query string and agrs.
func (p *BaseProvider) One(query string, cb pd.ScanCallback, args ...any) error {
	if !p.prepared() || query == "" || cb == nil {
		return invar.ErrBadDBConnect
	}

	query = p.Builder.CheckLimit(query)
	rows, err := p.client.DB().Query(query, args...)
	if err != nil {
		return err
	}

	defer rows.Close()
	if !rows.Next() {
		return invar.ErrNotFound
	}
	rows.Columns()
	return cb(rows)
}

// Execute query string to get the top one record with outs non-nil params,
// it will auto append 'LIMIT 1' as tail in query string for high-performance.
//
//	Use QueryBuilder to build a query string and agrs.
func (p *BaseProvider) OneDone(query string, outs []any, done pd.DoneCallback, args ...any) error {
	if !p.prepared() || query == "" || len(outs) <= 0 {
		return invar.ErrBadDBConnect
	}

	query = p.Builder.CheckLimit(query)
	rows, err := p.client.DB().Query(query, args...)
	if err != nil {
		return err
	}

	defer rows.Close()
	if !rows.Next() {
		return invar.ErrNotFound
	}
	rows.Columns()
	if err := rows.Scan(outs...); err != nil {
		return err
	} else if done != nil {
		done()
	}
	return nil
}

// Execute query string with scan callback to read result records.
//
//	Use QueryBuilder to build a query string and agrs.
func (p *BaseProvider) Query(query string, cb pd.ScanCallback, args ...any) error {
	if !p.prepared() || query == "" || cb == nil {
		return invar.ErrBadDBConnect
	}

	rows, err := p.client.DB().Query(query, args...)
	if err != nil {
		return err
	}

	defer rows.Close()
	for rows.Next() {
		rows.Columns()
		if err := cb(rows); err != nil {
			return err
		}
	}
	return nil
}

// Execute query string to insert a row into target table which contain
// the 'auto increment' field of id as primary key.
//
//	Use InsertBuilder to build a query string and args.
func (p *BaseProvider) Insert(query string, args ...any) (int64, error) {
	if !p.prepared() || query == "" {
		return -1, invar.ErrBadDBConnect
	}

	stmt, err := p.client.DB().Prepare(query)
	if err != nil {
		return -1, err
	}

	defer stmt.Close()
	result, err := stmt.Exec(args...)
	if err != nil {
		return -1, err
	}
	return result.LastInsertId()
}

// Execute query string to insert multiple rows into target table with
// a callback for format row values, it insert all rows at one time to
// provide high-performance than call provider.Insert() one by one.
//
//	query := "INSERT table (field1, field2) VALUES"
//	err := provider.Inserts(query, len(vs), func(index int) string {
//		return fmt.Sprintf("(%v, %v)", v1, vs[index])
//		// return fmt.Sprintf("('%s', '%s')", v1, vs[index])
//	})
//	// => INSERT table (field1, field2) VALUES (1,2),(3,4)..
//	// => INSERT table (field1, field2) VALUES ('1','2'),('3','4')..
func (p *BaseProvider) Inserts(query string, cnt int, cb pd.InsertCallback) error {
	values := []string{}
	for i := 0; i < cnt; i++ {
		value := strings.TrimSpace(cb(i))
		if value != "" {
			values = append(values, value)
		}
	}
	query = query + " " + strings.Join(values, ",")
	return p.Exec(query)
}

// Execute query string to update target records by where condition, it will
// return invar.ErrNotChanged error when none updated.
//
//	Use UpdateBuilder to build a query string and args.
func (p *BaseProvider) Update(query string, args ...any) error {
	rows, err := p.ExecResult(query, args...)
	if err == nil && rows == 0 {
		return invar.ErrNotChanged
	}
	return err /* nil or error */
}

// Execute query string to delete records on given conditions, it will
// return invar.ErrNotChanged error when none deleted.
//
//	Use the DeleteBuilder to build a query string and args.
func (p *BaseProvider) Delete(query string, args ...any) error {
	rows, err := p.ExecResult(query, args...)
	if err == nil && rows == 0 {
		return invar.ErrNotChanged
	}
	return err /* nil or error */
}

// Clear all records for the given table.
func (p *BaseProvider) Clear(table string) error {
	if !p.prepared() || table == "" {
		return invar.ErrBadDBConnect
	}
	query := fmt.Sprintf("DELETE FROM %s", table)
	return p.Exec(query)
}

// Execute query string for single transaction, it will rollback when handle failed.
//
//	Use the anyone builder to build a query string and args.
func (p *BaseProvider) Tran(query string, args ...any) error {
	if !p.prepared() || query == "" {
		return invar.ErrBadDBConnect
	}

	tx, err := p.client.DB().Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()
	if _, err := tx.Exec(query, args...); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

// Excute multiple transactions, it will rollback when cased one error.
//
//	// Excute 3 transactions in callback with different query1 ~ 3
//	err := provider.Trans(
//		func(tx *sql.Tx) error { return provider.TxQuery(tx, query1, func(rows *sql.Rows) error {
//				// Fetch all rows to get result datas...
//			}, args...) },
//		func(tx *sql.Tx) error { return provider.TxExec(tx, query2, args...) },
//		func(tx *sql.Tx) error { return provider.TxExec(tx, query3, args...) })
func (p *BaseProvider) Trans(cbs ...pd.TransCallback) error {
	if !p.prepared() || len(cbs) == 0 {
		return invar.ErrBadDBConnect
	}

	tx, err := p.client.DB().Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()
	for _, cb := range cbs {
		if err := cb(tx); err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

/* ------------------------------------------------------------------- */
/* Helper Methods For Construct Query or Parse Results                 */
/* ------------------------------------------------------------------- */

// Check the database client whther prepared and connected.
func (p *BaseProvider) prepared() bool {
	return p.client != nil && p.client.DB() != nil
}

// Get update or delete record counts.
func (p *BaseProvider) Affected(result sql.Result) (int64, error) {
	rows, err := result.RowsAffected()
	if err != nil || rows == 0 {
		return 0, invar.ErrNotChanged
	}
	return rows, nil
}

// Get update or delete record counts without error check.
func (p *BaseProvider) Affects(result sql.Result) int64 {
	rows, _ := result.RowsAffected()
	return rows
}

// Get the last inserted record id without error check.
func (p *BaseProvider) LastID(result sql.Result) int64 {
	id, _ := result.LastInsertId()
	return id
}

/* ------------------------------------------------------------------- */
/* Util Methods For Access Table                                       */
/* ------------------------------------------------------------------- */

// Table datas for describe table structures.
type Table struct {
	Columns []*Column // Table column infos
	Spans   [6]int    // spans lenght for print table
}

// Table column datas.
type Column struct {
	Field string // Column name
	Type  string // Field value type
	Null  string // Flag for indicate field if null
	Def   string // Field default value
	Key   string // [Only MySQL] Primary key, foreign key or normal field
	Extra string // [Only MySQL] Extra infos
}

// Get target table structs by name from mysql databse.
func (p *BaseProvider) MysqlTable(table string, print ...bool) *Table {
	if !p.prepared() {
		return nil
	}

	db := p.client.DB()
	rows, err := db.Query("DESCRIBE " + table + ";")
	if err != nil {
		logger.E("Describe table:", table, "err:", err)
		return nil
	}
	defer rows.Close()

	cs, spans := []*Column{}, defPaddings()
	for rows.Next() {
		var def *string
		c := &Column{Def: "NULL"}
		if err := rows.Scan(&(c.Field), &(c.Type), &(c.Null), &(c.Key), &def, &(c.Extra)); err != nil {
			logger.E("Scane table:", table, "struct, err:", err)
			return nil
		}

		if def != nil {
			c.Def = *def
		}

		// calculate spans for format print
		if len(print) > 0 && print[0] {
			spans = calculatePaddings(c, spans)
		}
		cs = append(cs, c)
	}
	return &Table{Columns: cs, Spans: spans}
}

// Get target table structs by name from mssql database.
func (p *BaseProvider) MssqlTable(table string, print ...bool) *Table {
	if !p.prepared() {
		return nil
	}

	db := p.client.DB()
	query := "SELECT column_name, data_type, is_nullable, column_default FROM INFORMATION_SCHEMA.COLUMNS WHERE table_name='" + table + "';"
	rows, err := db.Query(query)
	if err != nil {
		logger.E("Describe table:", table, "err:", err)
		return nil
	}
	defer rows.Close()

	cs, spans := []*Column{}, defPaddings()
	for rows.Next() {
		var def *string
		c := &Column{Def: "NULL"}
		if err := rows.Scan(&(c.Field), &(c.Type), &(c.Null), &def); err != nil {
			logger.E("Scane table:", table, "struct, err:", err)
			return nil
		}

		if def != nil {
			c.Def = *def
		}

		// calculate spans for format print
		if len(print) > 0 && print[0] {
			spans = calculatePaddings(c, spans)
		}
		cs = append(cs, c)
	}
	return &Table{Columns: cs, Spans: spans}
}

// Print target table structs.
//
//	table := provider.MysqlTable("config", true)
//	provider.PrintTable(table)
func (p *BaseProvider) PrintTable(table *Table) {
	ps, cnt := table.Spans, len(table.Columns)
	for i, c := range table.Columns {
		if i == 0 {
			printHeader(1, ps) // +------------------------------------------------+
			printHeader(2, ps) // | FIELD | TYPE | IS NULL | DEFAULT | KEY | EXTRA |
			printHeader(3, ps) // |-------+------+---------+-----+---------+-------|
		}

		fmt.Printf("| %s | %s | %s | %s | %s | %s |\n",
			withSpan(c.Field, ps[0]), withSpan(c.Type, ps[1]), withSpan(c.Null, ps[2]),
			withSpan(c.Def, ps[3]), withSpan(c.Key, ps[4]), withSpan(c.Extra, ps[5]))

		if i == cnt-1 {
			printHeader(1, ps) // +------------------------------------------------+
		}
	}
}

// Calculate padding spans to print table struct as formated.
func calculatePaddings(c *Column, paddings [6]int) [6]int {
	fields := []string{c.Field, c.Type, c.Null, c.Def, c.Key, c.Extra}
	for i, field := range fields {
		if flen := len(field); flen > paddings[i] {
			paddings[i] = flen
		}
	}
	return paddings
}

// Return default header paddings.
//
// ------------------------------------------------------
// | FIELD | TYPE | IS NULLABLE | DEFAULT | KEY | EXTRA |
// ------------------------------------------------------
func defPaddings() [6]int { return [6]int{5, 4, 11, 7, 3, 5} }

// Tial ' ' chars into given text if length over max.
func withSpan(text string, max int) string {
	if cnt := len(text); cnt < max {
		text += strings.Repeat(" ", max-cnt)
	}
	return text
}

// Get max length divider as '---'.
func asDivider(max int) string {
	if max > 0 {
		return strings.Repeat("-", max)
	}
	return ""
}

// Print table columns labels on formated.
func printHeader(header int, ps [6]int) {
	switch header {
	case 1: // the table start and end line
		fmt.Printf("+-%s---%s---%s---%s---%s---%s-+\n",
			asDivider(ps[0]), asDivider(ps[1]), asDivider(ps[2]),
			asDivider(ps[3]), asDivider(ps[4]), asDivider(ps[5]))

	case 2: // the table header label line
		fmt.Printf("| %s | %s | %s | %s | %s | %s |\n",
			withSpan("FIELD", ps[0]), withSpan("TYPE", ps[1]), withSpan("IS NULL", ps[2]),
			withSpan("DEFAULT", ps[3]), withSpan("KEY", ps[4]), withSpan("EXTRA", ps[5]))

	case 3: // the header and content diliver line
		fmt.Printf("|-%s-+-%s-+-%s-+-%s-+-%s-+-%s-|\n",
			asDivider(ps[0]), asDivider(ps[1]), asDivider(ps[2]),
			asDivider(ps[3]), asDivider(ps[4]), asDivider(ps[5]))
	}
}
