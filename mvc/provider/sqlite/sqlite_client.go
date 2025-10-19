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
func New(opts ...Option) *Sqlite {
	client := &Sqlite{options: DefaultOptions()}
	for _, optFunc := range opts {
		optFunc(client)
	}
	return client
}

// Create a Sqlite client and connect with options which loaded from app.conf file.
//
// This method useful for beego project easy to connect a Sqlite database.
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

// Create and return BaseProvider instance with Sqlite client.
func GetProvider(session ...string) *pd.BaseProvider {
	return pd.NewProvider(Select(session...))
}

// Create and return BaseProvider instance with Sqlite client.
func GetTabler(opts ...pd.Option) *pd.TableProvider {
	return pd.NewTabler(Select(), opts...)
}

// Setup tables with name and provider.
func SetupTables(tables map[string]pd.TableSetup, debug ...bool) {
	isdebug := utils.Variable(debug, false)
	for name, table := range tables {
		table.Setup(Select(),
			pd.WithTable(name),
			pd.WithDebug(isdebug))
	}
}

// Return Sqlite database client, maybe nil when not call Connect() before.
func (m *Sqlite) DB() *sql.DB {
	return m.conn
}

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
	conn.SetMaxOpenConns(1)
	m.conn = conn
	return nil
}

// Create the given tables for sqlite database if not exist.
func (m *Sqlite) CreateTables(tables []string) error {
	if m.conn == nil {
		return invar.ErrBadDBConnect
	}

	for index, stmt := range tables {
		if _, err := m.conn.Exec(stmt); err != nil {
			logger.E("Create table at", index, "err:", err)
			return err
		}
	}
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
