// Copyright (c) 2018-Now Dunyu All Rights Reserved.
//
// Author      : https://www.wengold.net
// Email       : support@wengold.net
//
// Prismy.No | Date       | Modified by. | Description
// -------------------------------------------------------------------
// 00001       2025/07/01   yangping       New version
// -------------------------------------------------------------------

package mssqlc

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
// import the follow driver for MSSQL database access.
//
// _ "github.com/denisenkom/go-mssqldb" // use for sql server 2017 ~ 2019
//
// ----------------------------------------

// MSSQL client for access target Microsoft SQL Server database.
type MSSQL struct {
	options Options
	client  *sql.DB
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

// Create a MSSQL client, set the options by mssqlc.WithXxxx(x) setters.
//
//	client := mssqlc.NewMSSQL(
//		mssqlc.WithSession("mssql"),
//		mssqlc.WithHost("127.0.0.1"),
//		mssqlc.WithPort(1433),
//		mssqlc.WithUser("sa"),
//		mssqlc.WithPassword("123456"),
//		mssqlc.WithDatabase("TestDB"),
//		mssqlc.WithTimeout(30),
//		mssqlc.WithMaxIdles(100),
//		mssqlc.WithMaxOpens(100),
//	)
func NewMSSQL(opts ...Option) *MSSQL {
	client := &MSSQL{options: DefaultOptions(_mssqlDriver)}
	for _, optFunc := range opts {
		optFunc(client)
	}
	return client
}

// Create a MSSQL client and connect with options which loaded from app.conf file.
//
// This method useful for beego project easy to connect a mssql database.
func OpenMSSQL(session ...string) error {
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

// Create and return BaseProvider instance with MSSQL client.
func GetProvider() *pd.BaseProvider {
	return pd.NewProvider(Select())
}

// Create and return TableProvider instance with MSSQL client.
func GetTabler(opts ...pd.Option) *pd.TableProvider {
	return pd.NewTabler(Select(), opts...)
}

// Setup tables with name and provider.
func SetupTables(tables map[string]pd.TableSetup) {
	for name, table := range tables {
		table.Setup(Select(), pd.WithTable(name))
	}
}

// Return MSSQL database client, maybe nil when not call Connect() before.
func (m *MSSQL) DB() *sql.DB {
	return m.client
}

// Connect mssql database and cache the client to MSSQL clients pool.
func (m *MSSQL) Connect() error {
	o := m.options
	dsn := fmt.Sprintf(_mssqlDsn, o.Host, o.Port, o.Database, o.User, o.Password, o.Timeout, o.Timeout+5)
	logger.I("Connect MSSQL from session:", o.Session)

	// open and connect database.
	con, err := sql.Open(_mssqlDriver, dsn)
	if err != nil {
		return err
	}

	// check database validable.
	if err = con.Ping(); err != nil {
		return err
	}

	con.SetMaxIdleConns(o.MaxIdles)
	con.SetMaxOpenConns(o.MaxOpens)
	m.client = con
	return nil
}

// Close the MSSQL client and remove from cache pool.
func (m *MSSQL) Close() error {
	if m.client != nil {
		if err := m.client.Close(); err != nil {
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
