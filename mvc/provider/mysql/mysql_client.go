// Copyright (c) 2018-Now Dunyu All Rights Reserved.
//
// Author      : https://www.wengold.net
// Email       : support@wengold.net
//
// Prismy.No | Date       | Modified by. | Description
// -------------------------------------------------------------------
// 00001       2025/07/01   yangping       New version
// -------------------------------------------------------------------

package mysql

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
	conn    *sql.DB
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

	// default charset for database access
	_defCharset = "utf8mb4"
)

// Create a MySQL client, set the options by mysql.WithXxxx(x) setters.
//
//	client := mysql.New(
//		mysql.WithSession("mysql"),
//		mysql.WithHost("127.0.0.1:3306"), // maybe empty for localhost.
//		mysql.WithUser("user"),
//		mysql.WithPassword("123456"),
//		mysql.WithDatabase("testdb"),
//		mysql.WithCharset("utf8mb4"),
//		mysql.WithMaxIdles(100),
//		mysql.WithMaxOpens(100),
//		mysql.WithMaxLifetime(28740),
//	)
//
// # NOTICE:
//
// 1. This method create a Sqlite instance by WithXxxx(x) options setters,
// not load from ./conf/app.conf file.
//
// 2. The created MySQL instance must call Connect() to connect the target
// database before use query, insert, update and delete methods. it will
// cache the instance to clients map, so use Select() to get the default
// or target instance by given session is safly.
func New(opts ...Option) *MySQL {
	client := &MySQL{options: DefaultOptions(_mysqlDriver)}
	for _, optFunc := range opts {
		optFunc(client)
	}

	session := client.options.Session
	if _, ok := _mysqlClients[session]; ok {
		logger.W("Override exist MySQL client:", session)
	}
	_mysqlClients[session] = client
	return client
}

// Create a MySQL client and connect with options which loaded from app.conf file.
//
//	[mysql]
//	host = "192.168.100.102:3306"
//	name = "sampledb"
//	user = "root"
//	pwd  = "123456"
//
// # NOTICE:
//	- This method useful for beego project easy to connect a mysql database.
func Open(charset string, session ...string) error {
	return OpenWithOptions(LoadOptions(session...), charset)
}

// Create a MySQL client by given options, and connect with database.
func OpenWithOptions(opts Options, charset ...string) error {
	if opts.Database == "" || opts.User == "" || opts.Password == "" {
		return invar.ErrInvalidConfigs
	}

	opts.Charset = utils.Variable(charset, opts.Charset)
	opts.Session = utils.Condition(opts.Session == "", _mysqlDriver, opts.Session)
	opts.Charset = utils.Condition(opts.Charset == "", _defCharset, opts.Charset)

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

// Create and return a BaseProvider instance with MySQL client.
func NewBase(session ...string) *pd.BaseProvider {
	return pd.NewBaseProvider(Select(session...))
}

// Create and return a TableProvider instance with MySQL client.
func NewTable(table string, debug bool, session ...string) *pd.TableProvider {
	return pd.NewTableProvider(Select(session...), pd.WithTable(table), pd.WithDebug(debug))
}

// Return MySQL database client, maybe nil when not call Connect() before.
func (m *MySQL) DB() *sql.DB { return m.conn }

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
	conn, err := sql.Open(_mysqlDriver, dsn)
	if err != nil {
		return err
	}

	// check database validable.
	if err = conn.Ping(); err != nil {
		return err
	}

	conn.SetMaxIdleConns(o.MaxIdles)
	conn.SetMaxOpenConns(o.MaxOpens)
	conn.SetConnMaxLifetime(o.MaxLifetime)
	m.conn = conn
	return nil
}

// Close the MySQL client and remove from cache pool.
func (m *MySQL) Close() error {
	if m.conn != nil {
		if err := m.conn.Close(); err != nil {
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
