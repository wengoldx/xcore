// Copyright (c) 2018-Now Dunyu All Rights Reserved.
//
// Author      : https://www.wengold.net
// Email       : support@wengold.net
//
// Prismy.No | Date       | Modified by. | Description
// -------------------------------------------------------------------
// 00001       2025/07/01   yangping       New version
// -------------------------------------------------------------------

package mssql

import (
	"database/sql"
	"fmt"

	"github.com/wengoldx/xcore/invar"
	"github.com/wengoldx/xcore/logger"
	pd "github.com/wengoldx/xcore/mvc/provider"
	"github.com/wengoldx/xcore/mvc/provider/provider"
	"github.com/wengoldx/xcore/utils"
)

// ----------------------------------------
// NOTIC :
//
// import the follow driver for MSSQL database access.
//
// _ "github.com/denisenkom/go-mssqldb" // use for sql server 2017 ~ 2019
//
// ----------------------------------------

// MSSQL client for access target Microsoft SQL Server database.
type MSSQL struct {
	options Options
	conn    *sql.DB
}

var _ pd.DBClient = (*MSSQL)(nil)

// MSSQL clients pool for cache multiple connected clients.
var _mssqlClients = make(map[string]pd.DBClient)

const (
	// Microsoft SQL Server driver name.
	_mssqlDriver = "mssql" // 'mssql' for processQueryText=true, 'sqlserver' for false

	// Microsoft SQL Server database source name.
	_mssqlDsn = "server=%s;port=%d;database=%s;user id=%s;password=%s;connection timeout=%d;dial timeout=%d;"
)

// Create a MSSQL client, set the options by mssql.WithXxxx(x) setters.
//
//	client := mssql.NewMSSQL(
//		mssql.WithSession("mssql"),
//		mssql.WithHost("127.0.0.1"),
//		mssql.WithPort(1433),
//		mssql.WithUser("sa"),
//		mssql.WithPassword("123456"),
//		mssql.WithDatabase("TestDB"),
//		mssql.WithTimeout(30),
//		mssql.WithMaxIdles(100),
//		mssql.WithMaxOpens(100),
//	)
//
// # NOTICE:
//
// 1. This method create a Sqlite instance by WithXxxx(x) options setters,
// not load from ./conf/app.conf file.
//
// 2. The created MSSQL instance must call Connect() to connect the target
// database before use query, insert, update and delete methods. it will
// cache the instance to clients map, so use Select() to get the default
// or target instance by given session is safly.
func New(opts ...Option) *MSSQL {
	client := &MSSQL{options: DefaultOptions(_mssqlDriver)}
	for _, optFunc := range opts {
		optFunc(client)
	}

	session := client.options.Session
	if _, ok := _mssqlClients[session]; ok {
		logger.W("Override exist MSSQL client:", session)
	}
	_mssqlClients[session] = client
	return client
}

// Find and return the exist MSSQL instance by given session.
func Select(session ...string) pd.DBClient {
	return _mssqlClients[utils.Variable(session, _mssqlDriver)]
}

// Close and remove the target MSSQL client.
func Close(session ...string) error {
	s := utils.Variable(session, _mssqlDriver)
	if client := Select(s); client != nil {
		defer delete(_mssqlClients, s)
		return client.Close()
	}
	return nil
}

/* ------------------------------------------------------------------- */
/* Setup Provider                                                      */
/* ------------------------------------------------------------------- */

// Create and return a BaseProvider instance with MSSQL client.
//
// # USAGE:
//
//	type MyTable struct{ provider.BaseProvider }
//	var MyTableIns = MyTable{ mssql.NewBase()}
//	// Call mssql.New(), or mssql.Open() to create mssql client here!
//	mssql.SetClient(MyTableIns)
//
// # WARNING:
//
// This method maybe init the nil DBClient client when mssql.Open(), or
// mssql.OpenWithOptions() not called, So call mssql.SetupTables() later
// to set valid DBClient client for all tables!
func NewBase(session ...string) *provider.BaseProvider {
	return provider.NewBaseProvider(Select(session...))
}

// Create and return a TableProvider instance with MSSQL client.
//
// # USAGE:
//
//	type MyTable struct{ pd.TableProvider }
//	var MyTableIns = MyTable{ mssql.NewTable("mytable", _logsql)}
//	// Call mssql.New(), or mssql.Open() to create mssql client here!
//	mssql.SetClient(MyTableIns)
//
// # WARNING:
//
// This method maybe init the nil DBClient client when mssql.Open(), or
// mssql.OpenWithOptions() not called, So call mssql.SetupTables() later
// to set valid DBClient client for all tables!
func NewTable(table string, debug bool, session ...string) pd.TableProvider {
	return provider.NewTableProvider(Select(session...),
		provider.WithTable(table), provider.WithDebug(debug))
}

// Bind tables with the DBClient client.
//
// # WARNING:
//
// Call mssql.Open(), or mssql.OpenWithOptions() first to ensure the
// DBClient client inited (not nil), later call this method to set tables
// DBClient client if need!
func SetClient(tables ...pd.Provider) {
	client := Select() // use the default session.
	for _, table := range tables {
		table.SetClient(client)
	}
}

/* ------------------------------------------------------------------- */
/* Create & Connect From app.conf                                      */
/* ------------------------------------------------------------------- */

// Create a MSSQL client and connect with options which loaded from app.conf file.
//
//	[mssql]
//	host    = "192.168.100.102"
//	port    = 1433
//	name    = "sampledb"
//	user    = "sa"
//	pwd     = "123456"
//	timeout = 30
//
// # NOTICE:
//	- This method useful for beego project easy to connect a mssql database.
func Open(session ...string) error {
	return OpenWithOptions(LoadOptions(session...))
}

// Create a MSSQL client by given options, and connect with database.
func OpenWithOptions(opts Options) error {
	if opts.Database == "" || opts.User == "" || opts.Password == "" {
		return invar.ErrInvalidConfigs
	} else if opts.Timeout <= 0 {
		opts.Timeout = 30 // fix the dial timeout over 30s
	}

	client := &MSSQL{options: opts}
	_mssqlClients[opts.Session] = client
	return client.Connect()
}

/* ------------------------------------------------------------------- */
/* DBClient Interface Implements                                       */
/* ------------------------------------------------------------------- */

// Return MSSQL database client, maybe nil when not call Connect() before.
func (m *MSSQL) DB() *sql.DB { return m.conn }

// Connect mssql database and cache the client to MSSQL clients pool.
func (m *MSSQL) Connect() error {
	o := m.options
	dsn := fmt.Sprintf(_mssqlDsn, o.Host, o.Port, o.Database, o.User, o.Password, o.Timeout, o.Timeout+5)
	logger.I("Connect MSSQL from session:", o.Session)

	// open and connect database.
	conn, err := sql.Open(_mssqlDriver, dsn)
	if err != nil {
		return err
	}

	// check database validable.
	if err = conn.Ping(); err != nil {
		return err
	}

	conn.SetMaxIdleConns(o.MaxIdles)
	conn.SetMaxOpenConns(o.MaxOpens)
	m.conn = conn
	return nil
}

// Close the MSSQL client and remove from cache pool.
func (m *MSSQL) Close() error {
	if m.conn != nil {
		if err := m.conn.Close(); err != nil {
			logger.E("Close MSSQL err:", err)
			return err
		}
	}

	// remove the cached MSSQL instance.
	if o := m.options; o.Session != "" {
		delete(_mssqlClients, o.Session)
	}
	return nil
}
