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
	"reflect"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
	"github.com/wengoldx/xcore/invar"
	"github.com/wengoldx/xcore/logger"
	// ----------------------------------------
	// NOTIC :
	//
	// import the follow database driver when using WingProvider.
	//
	// _ "github.com/go-sql-driver/mysql"   // use for mysql
	//
	// ----------------------------------------
)

// WingProvider content provider to support database utils,
// you can implement custom helper to access any table with mvc.WingHelper
// like the follow codes:
//
//	types CustomHelper struct {
//		mvc.WingProvider
//	}
//
//	func HStub() *CustomHelper {
//		return &CustomHelper{*mvc.WingHelper}
//	}
//
// `WARNING`: The Conn instance maybe not inited when database configs
// invalid by call mvc.OpenMySQL(), OpenMysqlName().
type WingProvider struct {
	Conn *sql.DB
}

// ScanCallback use for scan query result from rows
type ScanCallback func(rows *sql.Rows) error

// InsertCallback format query values to string for Inserts().
type InsertCallback func(index int) string

// TransCallback transaction callback for Trans().
type TransCallback func(tx *sql.Tx) error

// MySQL database configs
const (
	mysqlConfigUser = "%s::user" // configs key of mysql database user
	mysqlConfigPwd  = "%s::pwd"  // configs key of mysql database password
	mysqlConfigHost = "%s::host" // configs key of mysql database host and port
	mysqlConfigName = "%s::name" // configs key of mysql database name

	// Mysql Server database source name for local connection
	mysqldsnLocal = "%s:%s@/%s?charset=%s"

	// Mysql Server database source name for tcp connection
	mysqldsnTcp = "%s:%s@tcp(%s)/%s?charset=%s"
)

var (
	// WingHelper content provider to hold database connections,
	// the WingHelper.Conn pointer will nil before mvc.OpenMySQL() called.
	WingHelper = &WingProvider{} // empty as default

	// Cache all mysql providers into pool for multiple databases server connect.
	connPool = make(map[string]*WingProvider)
)

// readMySQLCofnigs read mysql database params from config file,
// than verify them if empty except host.
func readMySQLCofnigs(session string, check bool) (string, string, string, string, error) {
	user := beego.AppConfig.String(fmt.Sprintf(mysqlConfigUser, session))
	pwd := beego.AppConfig.String(fmt.Sprintf(mysqlConfigPwd, session))
	host := beego.AppConfig.String(fmt.Sprintf(mysqlConfigHost, session))
	name := beego.AppConfig.String(fmt.Sprintf(mysqlConfigName, session))

	if user == "" || pwd == "" || (check && name == "") {
		return "", "", "", "", invar.ErrInvalidConfigs
	}
	return user, pwd, host, name, nil
}

// openMySQLPool open mysql and cached to connection pool by given session keys.
func openMySQLPool(charset, name string, fix bool, sessions []string) error {
	for _, session := range sessions {
		// combine develop session key on dev mode
		if !fix && beego.BConfig.RunMode == "dev" {
			session = session + "-dev"
		}

		// check database name config when not set as param
		check := name == ""

		// load configs by session key
		user, pwd, host, dbn, err := readMySQLCofnigs(session, check)
		if err != nil {
			return err
		} else if check {
			name = dbn // use database name from config file
		}

		dsn := ""
		if len(host) > 0 /* check database host whether using TCP to connect */ {
			// conntect with remote host database server
			dsn = fmt.Sprintf(mysqldsnTcp, user, pwd, host, name, charset)
		} else {
			// just connect local database server
			dsn = fmt.Sprintf(mysqldsnLocal, user, pwd, name, charset)
		}
		logger.I("Open MySQL from session:", session)

		// open and connect database
		con, err := sql.Open("mysql", dsn)
		if err != nil {
			return err
		}

		// check database validable
		if err = con.Ping(); err != nil {
			return err
		}

		con.SetMaxIdleConns(100)
		con.SetMaxOpenConns(100)
		con.SetConnMaxLifetime(28740)
		connPool[session] = &WingProvider{con}
	}
	return nil
}

// OpenMySQL connect database and check ping result, the connection holded
// by mvc.WingHelper object if signle connect, or cached connections in connPool map
// if multiple connect and select same one by given sessions of input params.
// the datatable charset maybe 'utf8' or 'utf8mb4' same as database set.
//
// `USAGE`
//
// you must config database params in /conf/app.config file as follows
//
// ---
//
// #### Case 1 : For signle connect on prod mode.
//
//	[mysql]
//	host = "127.0.0.1:3306"
//	name = "sampledb"
//	user = "root"
//	pwd  = "123456"
//
// #### Case 2 : For signle connect on dev mode.
//
//	[mysql-dev]
//	host = "127.0.0.1:3306"
//	name = "sampledb"
//	user = "root"
//	pwd  = "123456"
//
// #### Case 3 : For both dev and prod mode, you can config all of up cases.
//
// #### Case 4 : For multi-connections to set custom session keywords.
//
//	[mysql-a]
//	... same as use Case 1.
//
//	[mysql-a-dev]
//	... same as use Case 2.
//
//	[mysql-x]
//	... same as use Case 1.
//
//	[mysql-x-dev]
//	... same as use Case 2.
func OpenMySQL(charset string, sessions ...string) error {
	if len(sessions) == 0 {
		sessions = []string{"mysql"}
	}

	// connect all mysql from sessions
	if err := openMySQLPool(charset, "", false, sessions); err != nil {
		return err
	}

	// using the first connection as primary helper
	WingHelper = Select(sessions[0])
	return nil
}

// OpenMySQLName connect database with target name, set fix=ture for ignore
// server runmode, it will read configs from fixed 'mysql' session, not from
// auto append 'mysql-dev' session when runmode on 'dev'.
func OpenMySQLName(charset, name string, fix bool) error {
	sessions := []string{"mysql"}
	if err := openMySQLPool(charset, name, fix, sessions); err != nil {
		return err
	}

	WingHelper = Select(sessions[0], fix)
	return nil
}

// Select mysql Connection by request key words
// if mode is dev, the key will auto splice '-dev'
func Select(session string, fix ...bool) *WingProvider {
	auto := !(len(fix) > 0 && fix[0])
	if auto && beego.BConfig.RunMode == "dev" {
		session = session + "-dev"
	}
	return connPool[session]
}

// Stub return content provider connection.
func (w *WingProvider) Stub() *sql.DB {
	return w.Conn
}

// IsEmpty call sql.Query() to check target data if empty.
func (w *WingProvider) IsEmpty(query string, args ...any) (bool, error) {
	if w.Conn == nil {
		return false, invar.ErrBadDBConnect
	}

	rows, err := w.Conn.Query(query, args...)
	if err != nil {
		return false, err
	}
	defer rows.Close()
	return !rows.Next(), nil
}

// IsExist call sql.Query() to check target data if exist.
func (w *WingProvider) IsExist(query string, args ...any) (bool, error) {
	empty, err := w.IsEmpty(query, args...)
	return !empty, err
}

// Count call sql.Query() to count results.
func (w *WingProvider) Count(query string, args ...any) (int, error) {
	if w.Conn == nil {
		return 0, invar.ErrBadDBConnect
	}

	if rows, err := w.Conn.Query(query, args...); err != nil {
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

// One call sql.Query() to query the top one record.
func (w *WingProvider) One(query string, cb ScanCallback, args ...any) error {
	if w.Conn == nil {
		return invar.ErrBadDBConnect
	}

	if rows, err := w.Conn.Query(query, args...); err != nil {
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

// Query call sql.Query() to query multiple records.
func (w *WingProvider) Query(query string, cb ScanCallback, args ...any) error {
	if w.Conn == nil {
		return invar.ErrBadDBConnect
	}

	if rows, err := w.Conn.Query(query, args...); err != nil {
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

// Insert call sql.Prepare() and stmt.Exec() to insert a new record.
//
// `@see` Use mvc.Inserts() to insert multiple values in once request.
func (w *WingProvider) Insert(query string, args ...any) (int64, error) {
	if w.Conn == nil {
		return -1, invar.ErrBadDBConnect
	}

	if stmt, err := w.Conn.Prepare(query); err != nil {
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

// Inserts format and combine multiple values to insert at once,
// this method can provide high-performance than call mvc.Insert() one by one.
//
// ---
//
//	query := "INSERT sametable (field1, field2) VALUES"
//	err := mvc.Inserts(query, len(vs), func(index int) string {
//		return fmt.Sprintf("(%v, %v)", v1, vs[index])
//		// return fmt.Sprintf("('%s', '%s')", v1, vs[index])
//	})
func (w *WingProvider) Inserts(query string, cnt int, cb InsertCallback) error {
	values := []string{}
	for i := 0; i < cnt; i++ {
		value := strings.TrimSpace(cb(i))
		if value != "" {
			values = append(values, value)
		}
	}
	query = query + " " + strings.Join(values, ",")
	return w.Execute(query)
}

// Inserts2 format and combine slice values to insert multiple items at once,
// this method can provide high-performance same as call mvc.Insert() without callback.
//
// ---
//
//	values := []Person{
//		{Age: 16, Male: true,  Name: "ZhangSan"},
//		{Age: 22, Male: false, Name: "LiXiang"},
//	}
//	query := "INSERT person (age, male, name) VALUES"
//	err := mvc.Inserts2(query, values)
func (w *WingProvider) Inserts2(query string, values any) error {
	items, err := w.FormatInserts(values)
	if err != nil {
		return err
	}
	return w.Execute(query + " " + items)
}

// Update call sql.Prepare() and stmt.Exec() to update record, then check the
// updated result if return invar.ErrNotChanged error when not changed any one.
//
// `@see` Use mvc.Updates() to update mapping values on slient.
//
// `@see` Use mvc.Execute() to update record on silent.
func (w *WingProvider) Update(query string, args ...any) error {
	rows, err := w.Execute2(query, args...)
	if rows == 0 {
		return invar.ErrNotChanged
	}
	return err /* nil or error */
}

// Updates update record from mapping values as colmun sets, it not check the
// updated result whatever changed or not.
//
// ---
//
//	values := map[string]any{ "Age": 16, "Name": "ZhangSan" }
//	query := "UPDATE person SET %s WHERE id=?"
//	err := mvc.updates(query, values, "id-123456")
//
// `@see` Use mvc.Update() to update record and check result.
func (w *WingProvider) Update2(query string, values map[string]any, args ...any) error {
	sets, err := w.FormatSets(values)
	if err != nil {
		return err
	}
	return w.Execute(fmt.Sprintf(query, sets), args...)
}

// Delete call sql.Prepare() and stmt.Exec() to delete record, then check the
// deleted result if return invar.ErrNotChanged error when none delete.
//
// `@see` Use mvc.Execute() to delete record on silent.
func (w *WingProvider) Delete(query string, args ...any) error {
	rows, err := w.Execute2(query, args...)
	if rows == 0 {
		return invar.ErrNotChanged
	}
	return err /* nil or error */
}

// Execute call sql.Prepare() and stmt.Exec() to insert, update or delete records
// without any result datas to return as silent.
//
// `@see` Use mvc.Execute2() return results.
func (w *WingProvider) Execute(query string, args ...any) error {
	if w.Conn == nil {
		return invar.ErrBadDBConnect
	}

	if stmt, err := w.Conn.Prepare(query); err != nil {
		return err
	} else {
		defer stmt.Close()
		if _, err := stmt.Exec(args...); err != nil {
			return err
		}
		return nil
	}
}

// Execute2 call sql.Prepare() and stmt.Exec() to update or delete records (but not
// for multiple inserts) with result counts to return.
//
// `@see` Use mvc.Execute() on silent, use mvc.Inserts() to multiple insert.
func (w *WingProvider) Execute2(query string, args ...any) (int64, error) {
	if w.Conn == nil {
		return 0, invar.ErrBadDBConnect
	}

	if stmt, err := w.Conn.Prepare(query); err != nil {
		return 0, err
	} else {
		defer stmt.Close()

		result, err := stmt.Exec(args...)
		if err != nil {
			return 0, err
		}
		return w.Affected(result)
	}
}

// TranRoll execute one sql transaction, it will rollback when operate failed.
//
// `@see` Use mvc.Trans() to excute multiple transaction as once.
func (w *WingProvider) TranRoll(query string, args ...any) error {
	if w.Conn == nil {
		return invar.ErrBadDBConnect
	}

	if tx, err := w.Conn.Begin(); err != nil {
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

// Trans excute multiple transactions, it will rollback when case any error.
//
// ---
//
//	// Excute 3 transactions in callback with different query1 ~ 3
//	err := mvc.Trans(
//		func(tx *sql.Tx) error { return mvc.TxQuery(tx, query1, func(rows *sql.Rows) error {
//				// Fetch all rows to get result datas...
//			}, args...) },
//		func(tx *sql.Tx) error { return mvc.TxExec(tx, query2, args...) },
//		func(tx *sql.Tx) error { return mvc.TxExec(tx, query3, args...) })
func (w *WingProvider) Trans(cbs ...TransCallback) error {
	if w.Conn == nil {
		return invar.ErrBadDBConnect
	}

	if tx, err := w.Conn.Begin(); err != nil {
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

// ----------------------------------------

// Affected get update or delete record counts.
func (w *WingProvider) Affected(result sql.Result) (int64, error) {
	row, err := result.RowsAffected()
	if err != nil || row == 0 {
		return 0, invar.ErrNotChanged
	}
	return row, nil
}

// Affects get update or delete record counts without error check.
func (w *WingProvider) Affects(result sql.Result) int64 {
	rows, _ := result.RowsAffected()
	return rows
}

// LastID get inserted record id without error check.
func (w *WingProvider) LastID(result sql.Result) int64 {
	id, _ := result.LastInsertId()
	return id
}

// AddLike append like field and keyword into given query or just return like string.
//
// - `query` : "SELECT * FROM tablename WHERE status=? %s ORDER BY id DESC",
//
// - `field` : "title", `keyword` : "Hello", `appendand` : true
//
// The result is "SELECT * FROM tablename WHERE status=? AND title LIKE '%%Hello%%' ORDER BY id DESC".
func (w *WingProvider) AddLike(query, field, keyword string, appendand ...bool) string {
	like := field + " LIKE '%%" + keyword + "%%'"
	if len(appendand) > 0 {
		like = "AND " + like
	}

	if query != "" {
		return fmt.Sprintf(query, like)
	}
	return like
}

// TailLimit tail limit condition into query string, default limit 1 record.
//
// - `query`  : "SELECT * FROM tablename WHERE status=?",
//
// - `limits` : 2
//
// The result is "SELECT * FROM tablename WHERE status=? LIMIT 2".
func (w *WingProvider) TailLimit(query string, limits ...int) string {
	num := "1"
	if len(limits) > 0 && limits[0] > 0 {
		num = strconv.Itoa(limits[0])
	}
	return query + " LIMIT " + num
}

// JoinInts join int64 numbers as string '1,2,3', or append to query strings as formart:
//
// - `query` : "SELECT * FROM tablename WHERE id IN (%s)",
//
// - `nums`  : []int64{1, 2, 3}
//
// The result is "SELECT * FROM tablename WHERE id IN (1,2,3)".
func (w *WingProvider) JoinInts(query string, nums []int64) string {
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
// - `query ` : "SELECT * FROM account WHERE uuid IN (%s)"
//
// - `values` : []string{"D23", "4R", "A34"}
//
// The result is "SELECT * FROM account WHERE uuid IN ('D23','4R','A34')"
func (w *WingProvider) JoinStrings(query string, values []string) string {
	if query != "" {
		return fmt.Sprintf(query, "'"+strings.Join(values, "','")+"'")
	}
	return "'" + strings.Join(values, "','") + "'"
}

// FormatSets format update sets for sql update.
//
// ---
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
func (w *WingProvider) FormatSets(values map[string]any) (string, error) {
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

// FormatInserts Format insert values for sql multiple insert.
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
func (w *WingProvider) FormatInserts(values any) (string, error) {
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

// ----------------------------------------

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
// ---
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

// -------------------------------------------

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
func (w *WingProvider) MysqlTable(table string, print ...bool) *Table {
	if w.Conn == nil {
		return nil
	}

	rows, err := w.Conn.Query("DESCRIBE " + table + ";")
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
func (w *WingProvider) MssqlTable(table string, print ...bool) *Table {
	if w.Conn == nil {
		return nil
	}

	query := "SELECT column_name, data_type, is_nullable, column_default FROM INFORMATION_SCHEMA.COLUMNS WHERE table_name='" + table + "';"
	rows, err := w.Conn.Query(query)
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
// `USAGE`
//
//	table := mvc.MysqlTable("config", true)
//	mvc.PrintTable(table)
func (w *WingProvider) PrintTable(table *Table) {
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
