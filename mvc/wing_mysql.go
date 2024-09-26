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

// WingProvider content provider to support database utils
type WingProvider struct {
	Conn *sql.DB
}

// ScanCallback use for scan query result from rows
type ScanCallback func(rows *sql.Rows) error

// InsertCallback format query values to string for Inserts().
type InsertCallback func(index int) string

// TransactionCallback transaction callback for Transactions().
type TransactionCallback func(tx *sql.Tx) (sql.Result, error)

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
	// it will nil before mvc.OpenMySQL() called.
	WingHelper *WingProvider

	// Cache all mysql providers into pool for multiple databases server connect.
	connPool = make(map[string]*WingProvider)
)

// readMySQLCofnigs read mysql database params from config file,
// than verify them if empty except host.
func readMySQLCofnigs(session string) (string, string, string, string, error) {
	user := beego.AppConfig.String(fmt.Sprintf(mysqlConfigUser, session))
	pwd := beego.AppConfig.String(fmt.Sprintf(mysqlConfigPwd, session))
	host := beego.AppConfig.String(fmt.Sprintf(mysqlConfigHost, session))
	name := beego.AppConfig.String(fmt.Sprintf(mysqlConfigName, session))

	if user == "" || pwd == "" || name == "" {
		return "", "", "", "", invar.ErrInvalidConfigs
	}
	return user, pwd, host, name, nil
}

// openMySQLPool open mysql and cached to connection pool by given session keys.
func openMySQLPool(charset string, sessions []string) error {
	for _, session := range sessions {
		// combine develop session key on dev mode
		if beego.BConfig.RunMode == "dev" {
			session = session + "-dev"
		}

		// load configs by session key
		dbuser, dbpwd, dbhost, dbname, err := readMySQLCofnigs(session)
		if err != nil {
			return err
		}

		dsn := ""
		if len(dbhost) > 0 /* check database host whether using TCP to connect */ {
			// conntect with remote host database server
			dsn = fmt.Sprintf(mysqldsnTcp, dbuser, dbpwd, dbhost, dbname, charset)
		} else {
			// just connect local database server
			dsn = fmt.Sprintf(mysqldsnLocal, dbuser, dbpwd, dbname, charset)
		}
		logger.I("Open MySQL on {", session, ":", dsn, "}")

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
	if len(sessions) > 0 {
		if err := openMySQLPool(charset, sessions); err != nil {
			return err
		}
		WingHelper = Select(sessions[0]) // using the first connection as primary helper
	} else {
		session := "mysql"
		if err := openMySQLPool(charset, []string{session}); err != nil {
			return err
		}
		WingHelper = Select(session)
	}
	return nil
}

// Select mysql Connection by request key words
// if mode is dev, the key will auto splice '-dev'
func Select(session string) *WingProvider {
	if beego.BConfig.RunMode == "dev" {
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

// Query directory call sql.Query() from database connection.
func (w *WingProvider) Query(query string, args ...any) (*sql.Rows, error) {
	return w.Conn.Query(query, args...)
}

// QueryOne call sql.Query() to query the top one record.
func (w *WingProvider) QueryOne(query string, cb ScanCallback, args ...any) error {
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

// QueryArray call sql.Query() to query multiple records.
func (w *WingProvider) QueryArray(query string, cb ScanCallback, args ...any) error {
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
// `@see` Use Inserts() to insert multiple values in once database operation.
func (w *WingProvider) Insert(query string, args ...any) (int64, error) {
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
// this method can provide high-performance than call Insert() one by one.
//
// ---
//
//	query := "INSERT sametable (field1, field2) VALUES"
//	err := mvc.Inserts(query, len(vs), func(index int) string {
//		return fmt.Sprintf("(%v, %v)", v1, vs[index])
//
//		// For string values like follows:
//		// return fmt.Sprintf("(\"%s\", \"%s\")", v1, vs[index])
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
// this method can provide high-performance same as call Insert() without callback.
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

// Updates format the values map and update result sets.
//
// ---
//
//	values := map[string]any{ "Age": 16, "Name": "ZhangSan" }
//	query := "UPDATE person SET %s WHERE id=?"
//	err := mvc.updates(query, values, "id-123456")
func (w *WingProvider) Updates(query string, values map[string]any, args ...any) error {
	sets, err := w.FormatSets(values)
	if err != nil {
		return err
	}
	return w.Execute(fmt.Sprintf(query, sets), args...)
}

// Execute call sql.Prepare() and stmt.Exec() to insert, update or delete records
// without any result datas to return as silent.
//
// `@see` Use Execute2() return results.
func (w *WingProvider) Execute(query string, args ...any) error {
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
// `@see` Use Execute() on silent, use Inserts() to multiple insert.
func (w *WingProvider) Execute2(query string, args ...any) (int64, error) {
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

// Transaction execute one sql transaction, it will rollback when operate failed.
//
// `@see` Use Transactions() to excute multiple transaction as once.
func (w *WingProvider) Transaction(query string, args ...any) error {
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

// Transactions excute multiple transactions, it will rollback when case any error.
//
// ---
//
//	// Excute 3 transactions in callback with different query1 ~ 3
//	err := mvc.Transactions(
//		func(tx *sql.Tx) (sql.Result, error) { return tx.Exec(query1, args...) },
//		func(tx *sql.Tx) (sql.Result, error) { return tx.Exec(query2, args...) },
//		func(tx *sql.Tx) (sql.Result, error) { return tx.Exec(query3, args...) })
func (w *WingProvider) Transactions(cbs ...TransactionCallback) error {
	if tx, err := w.Conn.Begin(); err != nil {
		return err
	} else {
		defer tx.Rollback()

		// start excute multiple transactions in callback
		for _, cb := range cbs {
			if _, err := cb(tx); err != nil {
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

// AppendLike append like keyword end of sql string, DON'T call it when exist limit key in sql string.
func (w *WingProvider) AppendLike(query, filed, keyword string, and ...bool) string {
	if len(and) > 0 && and[0] {
		return query + " AND " + filed + " LIKE '%%" + keyword + "%%'"
	}
	return query + " WHERE " + filed + " LIKE '%%" + keyword + "%%'"
}

// Affected get update or delete record counts.
func (w *WingProvider) Affected(result sql.Result) (int64, error) {
	row, err := result.RowsAffected()
	if err != nil || row == 0 {
		return 0, invar.ErrNotChanged
	}
	return row, nil
}

// JoinIDs join int64 ids as string '1,2,3', or append to query strings as formart:
//
// - `query`` : "SELECT * FROM tablename WHERE id IN (%s)",
//
// - `ids``   : []int64{1, 2, 3}
//
// The result is "SELECT * FROM tablename WHERE id IN (1,2,3)".
func (w *WingProvider) JoinIDs(query string, ids []int64) string {
	if len(ids) > 0 {
		vs := []string{}
		for _, id := range ids {
			if v := strconv.FormatInt(id, 10); v != "" {
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
//	// => Age=16, Male=true, Name="ZhangSan", Height=176.8
func (w *WingProvider) FormatSets(values map[string]any) (string, error) {
	sets := []string{}
	for key, value := range values {
		if key == "" && value == nil {
			continue
		}

		v := reflect.ValueOf(value)
		if v.Kind() == reflect.String {
			sets = append(sets, fmt.Sprintf(key+"=\"%s\"", value))
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
			items = append(items, fmt.Sprintf("(\"%s\")", item))
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
					fields = append(fields, fmt.Sprintf("\"%s\"", itv))
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
