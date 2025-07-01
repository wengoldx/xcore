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
)

// ----------------------------------------
// NOTIC :
//
// import the follow database driver when using WingProvider.
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

// Create a MySQL instance, set the options by using mysqlc.WithXxxx(x) functions.
func NewMySQL(opts ...Option) *MySQL {
	client := &MySQL{options: DefaultOptions(_mysqlDriver)}
	for _, optFunc := range opts {
		optFunc(client)
	}
	return client
}

// Create and open a MySQL client by load options from app.conf file.
//
// The function useful for beego backend project to connect mysql database.
func OpenMySQL(charset string, session ...string) error {
	opts := LoadOptions(session...)
	return OpenWithOptions(charset, opts)
}

// Create and open a MySQL client by exist options.
func OpenWithOptions(charset string, opts Options) error {
	if opts.Database == "" || opts.User == "" || opts.Password == "" {
		return invar.ErrInvalidConfigs
	}

	client := NewMySQL(
		WithSession(opts.Session),
		WithHost(opts.Host), // maybe empty for localhost.
		WithUser(opts.User),
		WithPassword(opts.Password),
		WithDatabase(opts.Database),
		WithCharset(charset),
		WithMaxIdles(opts.MaxIdles),
		WithMaxOpens(opts.MaxOpens),
		WithMaxLifetime(opts.MaxLifetime),
	)
	_mysqlClients[opts.Session] = client
	return client.Connect()
}

// Find and return the exist MySQL instance by given session.
func Select(session string) pd.DBClient {
	return _mysqlClients[session]
}

// Close and remove the target MySQL client.
func Close(session string) error {
	if client := Select(session); client != nil {
		defer delete(_mysqlClients, session)
		return client.Close()
	}
	return nil
}

// Create and return BaseProvider instance with MySQL client.
func GetProvider() pd.BaseProvider {
	return *pd.NewProvider(Select(_mysqlDriver))
}

// Return MySQL database client, maybe nil when not call Connect() before.
func (m *MySQL) DB() *sql.DB {
	return m.client
}

// Connect mysql database and cache the client to MySQL clients pool.
func (m *MySQL) Connect() error {
	dsn := ""
	if len(m.options.Host) > 0 {
		// conntect with remote host database server.
		dsn = fmt.Sprintf(_mysqlDsnTcp, m.options.User, m.options.Password,
			m.options.Host, m.options.Database, m.options.Charset)
	} else {
		// just connect local database server.
		dsn = fmt.Sprintf(_mysqlDsnLocal, m.options.User, m.options.Password,
			m.options.Database, m.options.Charset)
	}
	logger.I("Connect MySQL from", m.options.Session)

	// open and connect database.
	con, err := sql.Open(_mysqlDriver, dsn)
	if err != nil {
		return err
	}

	// check database validable.
	if err = con.Ping(); err != nil {
		return err
	}

	con.SetMaxIdleConns(m.options.MaxIdles)
	con.SetMaxOpenConns(m.options.MaxOpens)
	con.SetConnMaxLifetime(m.options.MaxLifetime)
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
	if m.options.Session != "" {
		delete(_mysqlClients, m.options.Session)
	}
	return nil
}
