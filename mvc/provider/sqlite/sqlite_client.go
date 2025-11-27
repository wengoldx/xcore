// Copyright (c) 2018-Now Dunyu All Rights Reserved.
//
// Author      : https://www.wengold.net
// Email       : support@wengold.net
//
// Prismy.No | Date       | Modified by. | Description
// -------------------------------------------------------------------
// 00001       2025/07/01   yangping       New version
// -------------------------------------------------------------------

package sqlite

import (
	"database/sql"

	"github.com/wengoldx/xcore/invar"
	"github.com/wengoldx/xcore/logger"
	pd "github.com/wengoldx/xcore/mvc/provider"
	"github.com/wengoldx/xcore/utils"
)

// ----------------------------------------
// NOTIC :
//
// import the follow driver for Sqlite database access.
//
// _ "github.com/mattn/go-sqlite3"   // use for sqlite
//
// ----------------------------------------

// Sqlite client for access target sqlite database.
type Sqlite struct {
	options Options
	conn    *sql.DB
}

var _ pd.DBClient = (*Sqlite)(nil)

// Sqlite clients pool for cache multiple connected clients.
var _sqliteClients = make(map[string]pd.DBClient)

const (
	// Sqlite Server driver name.
	_sqliteDriver = "sqlite3"

	// Sqlite default config session name.
	_sqliteDefSession = "sqlite"

	// Sqlite database mode for memroy.
	_sqliteMemDB = ":memory:"
)

// Create a Sqlite client, set the options by sqlite.WithXxxx(x) setters.
//
// # USAGE:
//
// 1. Create and connect with testdb file database.
//
//	client := sqlite.New(
//		sqlite.WithSession("mysqlite"),
//		sqlite.WithDatabase("testdb"),
//	)
//
// 2. Create and connect with memory database.
//
//	client := sqlite.New(
//		sqlite.WithSession("mysqlite"),
//		sqlite.WithIsMemory(true),
//	)
//
// # NOTICE:
//
// 1. This method create a Sqlite instance by WithXxxx(x) options setters,
// not load from ./conf/app.conf file.
//
// 2. The created Sqlite instance must call Connect() to connect the target
// database before use query, insert, update and delete methods. it will
// cache the instance to clients map, so use Select() to get the default
// or target instance by given session is safly.
func New(opts ...Option) *Sqlite {
	client := &Sqlite{options: DefaultOptions()}
	for _, optFunc := range opts {
		optFunc(client)
	}

	session := client.options.Session
	if _, ok := _sqliteClients[session]; ok {
		logger.W("Override exist Sqlite client:", session)
	}
	_sqliteClients[session] = client
	return client
}

// Create a Sqlite client and connect with options which loaded from app.conf file.
//
//	[sqlite]
//	database = "sample.db"  ; only for file database.
//	memory = false          ; set true for memory sqlite database.
//
// # NOTICE:
//	- This method useful for beego project easy to connect a Sqlite database.
func Open(session ...string) error {
	return OpenWithOptions(LoadOptions(session...))
}

// Create a Sqlite client by given options, and connect with database.
func OpenWithOptions(opts Options) error {
	if !opts.IsMemory && opts.Database == "" {
		return invar.ErrInvalidConfigs
	}

	client := &Sqlite{options: opts}
	_sqliteClients[opts.Session] = client
	return client.Connect()
}

// Find and return the exist Sqlite instance by given session.
func Select(session ...string) pd.DBClient {
	return _sqliteClients[utils.Variable(session, _sqliteDefSession)]
}

// Close and remove the target Sqlite client.
func Close(session ...string) error {
	s := utils.Variable(session, _sqliteDefSession)
	if client := _sqliteClients[s]; client != nil {
		defer delete(_sqliteClients, s)
		return client.Close()
	}
	return nil
}

// Create and return a BaseProvider instance with Sqlite client.
//
// # WARNING
//
// This method maybe init the nil DBClient client when sqlite.Open(), or
// sqlite.OpenWithOptions() not called, So call sqlite.SetupTables() later
// to set valid DBClient client for all tables!
func NewBase(session ...string) *pd.BaseProvider {
	return pd.NewBaseProvider(Select(session...))
}

// Create and return a TableProvider instance with Sqlite client.
//
// # WARNING
//
// This method maybe init the nil DBClient client when sqlite.Open(), or
// sqlite.OpenWithOptions() not called, So call sqlite.SetupTables() later
// to set valid DBClient client for all tables!
func NewTable(table string, debug bool, session ...string) *pd.TableProvider {
	return pd.NewTableProvider(Select(session...), pd.WithTable(table), pd.WithDebug(debug))
}

// Bind tables with the DBClient client of default session.
//
// # WARNING
//
// Call sqlite.Open(), or sqlite.OpenWithOptions() first to ensure the
// DBClient client inited (not nil), later call this method to set tables
// DBClient client if need!
func SetClient(tables ...pd.ClientStub) {
	client := Select() // use the default session.
	for _, table := range tables {
		table.SetClient(client)
	}
}

// Return Sqlite database client, maybe nil when not call Connect() before.
func (m *Sqlite) DB() *sql.DB { return m.conn }

// Connect sqlite database and cache the client to Sqlite clients pool.
func (m *Sqlite) Connect() error {
	dsn := utils.Condition(m.options.IsMemory, _sqliteMemDB, m.options.Database)
	logger.I("Connect Sqlite dabase", dsn)

	// open and connect database.
	conn, err := sql.Open(_sqliteDriver, dsn)
	if err != nil {
		return err
	}

	// check database validable.
	if err = conn.Ping(); err != nil {
		return err
	}

	conn.SetMaxIdleConns(1)
	conn.SetMaxOpenConns(20)
	m.conn = conn
	return nil
}

// Close the Sqlite client and remove from cache pool.
func (m *Sqlite) Close() error {
	if m.conn != nil {
		if err := m.conn.Close(); err != nil {
			logger.E("Close Sqlite err:", err)
			return err
		}
	}

	// remove the cached Sqlite instance.
	if o := m.options; o.Session != "" {
		delete(_sqliteClients, o.Session)
	}
	return nil
}

// Execute tables stmt string to create tables for database on connecte status
// if unexist, the stmt string like follow (sqlite3 driver):
//
//	const _table_stmt_settings = `CREATE TABLE IF NOT EXISTS settings (
//	    name    varchar (64)    PRIMARY KEY,    -- settings name.
//	    value   text                            -- settings value.
//	);`
func (m *Sqlite) CreateTables(tables ...string) error {
	if m.conn == nil {
		return invar.ErrBadDBConnect
	}

	for _, stmt := range tables {
		if _, err := m.conn.Exec(stmt); err != nil {
			return err
		}
	}
	return nil
}
