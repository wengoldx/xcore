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
	"reflect"
	"strconv"
	"strings"

	"github.com/wengoldx/xcore/invar"
	"github.com/wengoldx/xcore/logger"
)

// Base provider for simple access database datas.
type BaseProvider struct {
	client DBClient
}

var _ DataProvider = (*BaseProvider)(nil)

// Create a BaseProvider with given database client.
func NewProvider(client DBClient) *BaseProvider {
	return &BaseProvider{client: client}
}

/* ------------------------------------------------------------------- */
/* Util Methods For Database Access                                    */
/* ------------------------------------------------------------------- */

// Call sql.Query() to check target data if empty.
func (p *BaseProvider) IsEmpty(query string, args ...any) (bool, error) {
	if p.client == nil || p.client.DB() == nil {
		return false, invar.ErrBadDBConnect
	}

	db := p.client.DB()
	rows, err := db.Query(query, args...)
	if err != nil {
		return false, err
	}
	defer rows.Close()
	return !rows.Next(), nil
}

// Call sql.Query() to check target data if exist.
func (p *BaseProvider) IsExist(query string, args ...any) (bool, error) {
	empty, err := p.IsEmpty(query, args...)
	return !empty, err
}

// Call sql.Query() to count results.
func (p *BaseProvider) Count(query string, args ...any) (int, error) {
	if p.client == nil || p.client.DB() == nil {
		return 0, invar.ErrBadDBConnect
	}

	db := p.client.DB()
	if rows, err := db.Query(query, args...); err != nil {
		return 0, err
	} else {
		defer rows.Close()
		if !rows.Next() {
			return 0, invar.ErrNotFound
		}
		rows.Columns()

		counts := 0
		if err := rows.Scan(&counts); err != nil {
			return 0, err
		}
		return counts, nil
	}
}

// Call sql.Query() to query the top one record.
func (p *BaseProvider) One(query string, cb ScanCallback, args ...any) error {
	if p.client == nil || p.client.DB() == nil {
		return invar.ErrBadDBConnect
	}

	db := p.client.DB()
	if rows, err := db.Query(query, args...); err != nil {
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

// Call sql.Query() to query multiple records.
func (p *BaseProvider) Query(query string, cb ScanCallback, args ...any) error {
	if p.client == nil || p.client.DB() == nil {
		return invar.ErrBadDBConnect
	}

	db := p.client.DB()
	if rows, err := db.Query(query, args...); err != nil {
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

// Call sql.Prepare() and stmt.Exec() to insert a new record, and return the inserted id.
//
//	- Use provider.Inserts() to insert multiple values in once request.
func (p *BaseProvider) Insert(query string, args ...any) (int64, error) {
	if p.client == nil || p.client.DB() == nil {
		return -1, invar.ErrBadDBConnect
	}

	db := p.client.DB()
	if stmt, err := db.Prepare(query); err != nil {
		return -1, err
	} else {
		defer stmt.Close()

		result, err := stmt.Exec(args...)
		if err != nil {
			return -1, err
		}
		return result.LastInsertId()
	}
}

// Insert the format and combine multiple values at once.
//
// This method can provide high-performance than call provider.Insert() one by one.
//
//	query := "INSERT sametable (field1, field2) VALUES"
//	err := provider.Inserts(query, len(vs), func(index int) string {
//		return fmt.Sprintf("(%v, %v)", v1, vs[index])
//		// return fmt.Sprintf("('%s', '%s')", v1, vs[index])
//	})
func (p *BaseProvider) Inserts(query string, cnt int, cb InsertCallback) error {
	values := []string{}
	for i := 0; i < cnt; i++ {
		value := strings.TrimSpace(cb(i))
		if value != "" {
			values = append(values, value)
		}
	}
	query = query + " " + strings.Join(values, ",")
	return p.Execute(query)
}

// Insert the format and combine slice values at once.
//
// This method can provide high-performance same as call provider.Insert() without callback.
//
//	values := []Person{
//		{Age: 16, Male: true,  Name: "ZhangSan"},
//		{Age: 22, Male: false, Name: "LiXiang"},
//	}
//	query := "INSERT person (age, male, name) VALUES"
//	err := provider.Inserts2(query, values)
func (p *BaseProvider) Inserts2(query string, values any) error {
	items, err := p.FormatInserts(values)
	if err != nil {
		return err
	}
	return p.Execute(query + " " + items)
}

// Call sql.Prepare() and stmt.Exec() to update record, then check the
// updated result if return invar.ErrNotChanged error when not changed any one.
//
//	- Use provider.Updates() to update mapping values on silent.
//	- Use provider.Execute() to update record on silent.
func (p *BaseProvider) Update(query string, args ...any) error {
	rows, err := p.Execute2(query, args...)
	if rows == 0 {
		return invar.ErrNotChanged
	}
	return err /* nil or error */
}

// Update record from mapping values as colmun sets, it not check the
// updated result whatever changed or not.
//
//	values := map[string]any{ "Age": 16, "Name": "ZhangSan" }
//	query := "UPDATE person SET %s WHERE id=?"
//	err := provider.updates(query, values, "id-123456")
//
//	- Use provider.Update() to update record and check result.
func (p *BaseProvider) Update2(query string, values map[string]any, args ...any) error {
	sets, err := p.FormatSets(values)
	if err != nil {
		return err
	}
	return p.Execute(fmt.Sprintf(query, sets), args...)
}

// Call sql.Prepare() and stmt.Exec() to delete record, then check the
// deleted result if return invar.ErrNotChanged error when none delete.
//
//	- Use provider.Execute() to delete record on silent.
func (p *BaseProvider) Delete(query string, args ...any) error {
	rows, err := p.Execute2(query, args...)
	if rows == 0 {
		return invar.ErrNotChanged
	}
	return err /* nil or error */
}

// Call sql.Prepare() and stmt.Exec() to insert, update or delete records
// without any result datas to return as silent.
//
//	- Use provider.Execute2() return results.
func (p *BaseProvider) Execute(query string, args ...any) error {
	if p.client == nil || p.client.DB() == nil {
		return invar.ErrBadDBConnect
	}

	db := p.client.DB()
	if stmt, err := db.Prepare(query); err != nil {
		return err
	} else {
		defer stmt.Close()
		if _, err := stmt.Exec(args...); err != nil {
			return err
		}
		return nil
	}
}

// Call sql.Prepare() and stmt.Exec() to update or delete records (but not
// for multiple inserts) with result counts to return.
//
//	- Use provider.Execute() on silent, use provider.Inserts() to multiple insert.
func (p *BaseProvider) Execute2(query string, args ...any) (int64, error) {
	if p.client == nil || p.client.DB() == nil {
		return 0, invar.ErrBadDBConnect
	}

	db := p.client.DB()
	if stmt, err := db.Prepare(query); err != nil {
		return 0, err
	} else {
		defer stmt.Close()

		result, err := stmt.Exec(args...)
		if err != nil {
			return 0, err
		}
		return p.Affected(result)
	}
}

// Execute single sql transaction, it will rollback when operate failed.
//
//	- Use provider.Trans() to excute multiple transaction as once.
func (p *BaseProvider) TranRoll(query string, args ...any) error {
	if p.client == nil || p.client.DB() == nil {
		return invar.ErrBadDBConnect
	}

	db := p.client.DB()
	if tx, err := db.Begin(); err != nil {
		return err
	} else {
		defer tx.Rollback()

		if _, err := tx.Exec(query, args...); err != nil {
			return err
		}

		if err := tx.Commit(); err != nil {
			return err
		}
	}
	return nil
}

// Excute multiple transactions, it will rollback when case any error.
//
//	// Excute 3 transactions in callback with different query1 ~ 3
//	err := provider.Trans(
//		func(tx *sql.Tx) error { return provider.TxQuery(tx, query1, func(rows *sql.Rows) error {
//				// Fetch all rows to get result datas...
//			}, args...) },
//		func(tx *sql.Tx) error { return provider.TxExec(tx, query2, args...) },
//		func(tx *sql.Tx) error { return provider.TxExec(tx, query3, args...) })
func (p *BaseProvider) Trans(cbs ...TransCallback) error {
	if p.client == nil || p.client.DB() == nil {
		return invar.ErrBadDBConnect
	}

	db := p.client.DB()
	if tx, err := db.Begin(); err != nil {
		return err
	} else {
		defer tx.Rollback()

		// start excute multiple transactions in callback
		for _, cb := range cbs {
			if err := cb(tx); err != nil {
				return err
			}
		}

		if err := tx.Commit(); err != nil {
			return err
		}
	}
	return nil
}

/* ------------------------------------------------------------------- */
/* Util Methods For Simple Query                                       */
/* ------------------------------------------------------------------- */

// Join int64 numbers as string '1,2,3', or append to query strings as formart:
//
//	- `query` : "SELECT * FROM tablename WHERE id IN (%s)"
//	- `nums`  : []int64{1, 2, 3}
//
// The result is "SELECT * FROM tablename WHERE id IN (1,2,3)".
func (p *BaseProvider) JoinInts(query string, nums []int64) string {
	if len(nums) > 0 {
		vs := []string{}
		for _, num := range nums {
			if v := strconv.FormatInt(num, 10); v != "" {
				vs = append(vs, v)
			}
		}

		// Append ids into none-empty query string
		if query != "" {
			return fmt.Sprintf(query, strings.Join(vs, ","))
		}
		return strings.Join(vs, ",")
	}
	return query
}

// Join strings with ',', then insert into the given format string;
//
//	- `query ` : "SELECT * FROM account WHERE uuid IN (%s)"
//	- `values` : []string{"D23", "4R", "A34"}
//
// The result is "SELECT * FROM account WHERE uuid IN ('D23','4R','A34')"
func (p *BaseProvider) JoinStrings(query string, values []string) string {
	if query != "" {
		return fmt.Sprintf(query, "'"+strings.Join(values, "','")+"'")
	}
	return "'" + strings.Join(values, "','") + "'"
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

// Get inserted record id without error check.
func (p *BaseProvider) LastID(result sql.Result) int64 {
	id, _ := result.LastInsertId()
	return id
}

// -------------------------------------

// Format update sets for sql update.
//
//	values := map[string]any{
//		"":       123456,   // Filter out empty field
//		"Age":    16,
//		"Male":   true,
//		"Name":   "ZhangSan",
//		"Height": 176.8,
//		"Secure": nil,      // Filter out nil value
//	}
//	// => Age=16, Male=true, Name='ZhangSan', Height=176.8
func (p *BaseProvider) FormatSets(values map[string]any) (string, error) {
	sets := []string{}
	for key, value := range values {
		if key == "" && value == nil {
			continue
		}

		v := reflect.ValueOf(value)
		if v.Kind() == reflect.String {
			sets = append(sets, fmt.Sprintf(key+"='%s'", value))
		} else if v.Kind() == reflect.Bool || v.CanInt() || v.CanFloat() || v.CanUint() {
			sets = append(sets, fmt.Sprintf(key+"=%v", value))
		}
	}

	if len(sets) == 0 {
		return "", invar.ErrEmptyData
	}
	return strings.Join(sets, ","), nil
}

// Format insert values for sql multiple insert.
//
// ---
//
// `Usecase 1` : For struct objects.
//
//	values := []Person{
//		{Age: 16, Male: true,  Name: "ZhangSan"},
//		{Age: 22, Male: false, Name: "LiXiang"},
//	}
//	// => (16,true,'ZhangSan'),(22,false,'LiXiang')
//
// `Usecase 2` : For struct pointers, it will filter nil datas.
//
//	values := []*Person{
//		{Age: 16, Male: true,  Name: "ZhangSan"},
//		{Age: 22, Male: false, Name: "LiXiang"},
//		nil,
//	}
//	// => (16,true,'ZhangSan'),(22,false,'LiXiang')
//
// `Usecase 3` : For no-struct single value array
//
//	values := []string{"ZhangSan", "LiXiang"} // => ('ZhangSan'),('LiXiang')
//	values := []bool{true, false}             // => (true),(false)
//	values := []float64{1.6, 22}              // => (1.6),(22)
//	values := []int{16, -22}                  // => (16),(-22)
//
// ---
//
// `WARNING` : DO NOT define sliice item or struct field as pointer type like follows.
//
//	str:="123"; values := []*string{&str}     // Error input params
//	type Person struct {
//		Age  *int                             // Error struct field type define
//		Male bool
//		Name string
//	}
func (p *BaseProvider) FormatInserts(values any) (string, error) {
	pv := reflect.ValueOf(values)
	if pv.Kind() != reflect.Slice {
		return "", invar.ErrInvalidData
	}

	items := []string{}
	for i, cnt := 0, pv.Len(); i < cnt; i++ {
		item := pv.Index(i) // fetch values array item
		switch item.Kind() {
		case reflect.Struct:
			item = reflect.ValueOf(item.Interface())
		case reflect.Pointer:
			item = item.Elem()
		case reflect.String: // for string values array
			items = append(items, fmt.Sprintf("('%s')", item))
			continue
		default: // for basic data types values array
			if item.Kind() == reflect.Bool || item.CanInt() || item.CanFloat() || item.CanUint() {
				items = append(items, fmt.Sprintf("(%v)", item))
				continue
			}
			return "", invar.ErrInvalidData
		}

		// for struct or struct pointer array to parse fields
		if item.IsValid() && item.Kind() == reflect.Struct && item.NumField() > 0 {
			fields := []string{}
			for j, vs := 0, item.NumField(); j < vs; j++ {
				itv := item.Field(j) // fetch value fields
				if itv.Kind() == reflect.String {
					fields = append(fields, fmt.Sprintf("'%s'", itv))
				} else if itv.Kind() == reflect.Bool || itv.CanInt() || itv.CanFloat() || itv.CanUint() {
					fields = append(fields, fmt.Sprintf("%v", itv))
				} else {
					return "", invar.ErrInvalidData
				}
			}

			// join fields as '(1, "2", 3.4, -5, true, ...)'
			if its := strings.Join(fields, ","); its != "" {
				items = append(items, "("+its+")")
			}
		}
	}

	// check parse result and json items
	if len(items) == 0 {
		return "", invar.ErrEmptyData
	}
	return strings.Join(items, ","), nil
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
	if p.client == nil || p.client.DB() == nil {
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
	if p.client == nil || p.client.DB() == nil {
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
		for i := 0; i < max-cnt; i++ {
			text += " "
		}
	}
	return text
}

// Get max length divider as '---'.
func asDivider(max int) string {
	devider := ""
	for i := 0; i < max; i++ {
		devider += "-"
	}
	return devider
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
