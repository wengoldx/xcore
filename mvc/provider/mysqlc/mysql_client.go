// Copyright (c) 2018-Now Dunyu All Rights Reserved.
//
// Author      : https://www.wengold.net
// Email       : support@wengold.net
//
// Prismy.No | Date       | Modified by. | Description
// -------------------------------------------------------------------
// 00001       2025/07/01   yangping       New version
// -------------------------------------------------------------------

package mysqlc

import (
	"database/sql"
	"fmt"

	"github.com/wengoldx/xcore/invar"
	"github.com/wengoldx/xcore/logger"
	pd "github.com/wengoldx/xcore/mvc/provider"
	"github.com/wengoldx/xcore/utils"
)

// ----------------------------------------
// NOTIC :
//
// import the follow driver for MySQL database access.
//
// _ "github.com/go-sql-driver/mysql"   // use for mysql
//
// ----------------------------------------

// Mysql client for access target mysql database.
type MySQL struct {
	options Options
	client  *sql.DB
}

var _ pd.DBClient = (*MySQL)(nil)

// MySQL clients pool for cache multiple connected clients.
var _mysqlClients = make(map[string]pd.DBClient)

const (
	// Mysql Server driver name.
	_mysqlDriver = "mysql"

	// Mysql Server database source name for local connection.
	_mysqlDsnLocal = "%s:%s@/%s?charset=%s"

	// Mysql Server database source name for tcp connection.
	_mysqlDsnTcp = "%s:%s@tcp(%s)/%s?charset=%s"
)

// Create a MySQL client, set the options by mysqlc.WithXxxx(x) setters.
//
//	client := mysqlc.NewMySQL(
//		mysqlc.WithSession("mysql"),
//		mysqlc.WithHost("127.0.0.1:3306"), // maybe empty for localhost.
//		mysqlc.WithUser("user"),
//		mysqlc.WithPassword("123456"),
//		mysqlc.WithDatabase("testdb"),
//		mysqlc.WithCharset("utf8mb4"),
//		mysqlc.WithMaxIdles(100),
//		mysqlc.WithMaxOpens(100),
//		mysqlc.WithMaxLifetime(28740),
//	)
func NewMySQL(opts ...Option) *MySQL {
	client := &MySQL{options: DefaultOptions(_mysqlDriver)}
	for _, optFunc := range opts {
		optFunc(client)
	}
	return client
}

// Create a MySQL client and connect with options which loaded from app.conf file.
//
// This method useful for beego project easy to connect a mysql database.
func OpenMySQL(charset string, session ...string) error {
	return OpenWithOptions(charset, LoadOptions(session...))
}

// Create a MySQL client by given options, and connect with database.
func OpenWithOptions(charset string, opts Options) error {
	if opts.Database == "" || opts.User == "" || opts.Password == "" {
		return invar.ErrInvalidConfigs
	}

	client := &MySQL{options: opts}
	_mysqlClients[opts.Session] = client
	return client.Connect()
}

// Find and return the exist MySQL instance by given session.
func Select(session ...string) pd.DBClient {
	return _mysqlClients[utils.Variable(session, _mysqlDriver)]
}

// Close and remove the target MySQL client.
func Close(session ...string) error {
	s := utils.Variable(session, _mysqlDriver)
	if client := Select(s); client != nil {
		defer delete(_mysqlClients, s)
		return client.Close()
	}
	return nil
}

// Create and return BaseProvider instance with MySQL client.
func GetProvider() *pd.BaseProvider {
	return pd.NewProvider(Select())
}

// Create and return BaseProvider instance with MySQL client.
func GetSimpler(opts ...pd.Option) *pd.SimpleProvider {
	return pd.NewSimpler(Select(), opts...)
}

// Return MySQL database client, maybe nil when not call Connect() before.
func (m *MySQL) DB() *sql.DB {
	return m.client
}

// Connect mysql database and cache the client to MySQL clients pool.
func (m *MySQL) Connect() error {
	dsn, o := "", m.options
	if len(o.Host) > 0 {
		// conntect with remote host database server.
		dsn = fmt.Sprintf(_mysqlDsnTcp, o.User, o.Password, o.Host, o.Database, o.Charset)
	} else {
		// just connect local database server.
		dsn = fmt.Sprintf(_mysqlDsnLocal, o.User, o.Password, o.Database, o.Charset)
	}
	logger.I("Connect MySQL from session:", o.Session)

	// open and connect database.
	con, err := sql.Open(_mysqlDriver, dsn)
	if err != nil {
		return err
	}

	// check database validable.
	if err = con.Ping(); err != nil {
		return err
	}

	con.SetMaxIdleConns(o.MaxIdles)
	con.SetMaxOpenConns(o.MaxOpens)
	con.SetConnMaxLifetime(o.MaxLifetime)
	m.client = con
	return nil
}

// Close the MySQL client and remove from cache pool.
func (m *MySQL) Close() error {
	if m.client != nil {
		if err := m.client.Close(); err != nil {
			logger.E("Close MySQL err:", err)
			return err
		}
	}

	// remove the cached MySQL instance.
	if o := m.options; o.Session != "" {
		delete(_mysqlClients, o.Session)
	}
	return nil
}
