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
)

// ----------------------------------------
// NOTIC :
//
// import the follow database driver when using WingProvider.
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
	_mssqlDriver = "mssql"

	// Microsoft SQL Server database source name.
	_mssqlDsn = "server=%s;port=%d;database=%s;user id=%s;password=%s;connection timeout=%d;dial timeout=%d;"
)

// Create a MSSQL instance, set the options by using mssqlc.WithXxxx(x) functions.
func NewMSSQL(options ...Option) *MSSQL {
	client := &MSSQL{options: DefaultOptions(_mssqlDriver)}
	for _, optFunc := range options {
		optFunc(client)
	}
	return client
}

// Create and open a MSSQL client by load options from app.conf file.
//
// The function useful for beego backend project to connect mmsql database.
func OpenMSSQL(charset string, session ...string) error {
	opts := LoadOptions(session...)
	return OpenWithOptions(charset, opts)
}

// Create and open a MSSQL client by exist options.
func OpenWithOptions(charset string, opts Options) error {
	if opts.Database == "" || opts.User == "" || opts.Password == "" {
		return invar.ErrInvalidConfigs
	} else if opts.Timeout <= 0 {
		opts.Timeout = 30 // fix the dial timeout over 30s
	}

	client := NewMSSQL(
		WithSession(opts.Session),
		WithHost(opts.Host),
		WithPort(opts.Port),
		WithUser(opts.User),
		WithPassword(opts.Password),
		WithDatabase(opts.Database),
		WithTimeout(opts.Timeout),
		WithMaxIdles(opts.MaxIdles),
		WithMaxOpens(opts.MaxOpens),
	)
	_mssqlClients[opts.Session] = client
	return client.Connect()
}

// Find and return the exist MSSQL instance by given session.
func Select(session string) pd.DBClient {
	return _mssqlClients[session]
}

// Close and remove the target MSSQL client.
func Close(session string) error {
	if client := Select(session); client != nil {
		defer delete(_mssqlClients, session)
		return client.Close()
	}
	return nil
}

// Return MSSQL database client, maybe nil when not call Connect() before.
func (m *MSSQL) DB() *sql.DB {
	return m.client
}

// Connect mssql database and cache the client to MSSQL clients pool.
func (m *MSSQL) Connect() error {
	// driver := "mssql" // mssql for processQueryText=true, sqlserver for false
	dsn := fmt.Sprintf(_mssqlDsn, m.options.Host, m.options.Port, m.options.Database,
		m.options.User, m.options.Password, m.options.Timeout, m.options.Timeout+5)
	logger.I("Connect MSSQL from", m.options.Session)

	// open and connect database.
	con, err := sql.Open(_mssqlDriver, dsn)
	if err != nil {
		return err
	}

	// check database validable.
	if err = con.Ping(); err != nil {
		return err
	}

	con.SetMaxIdleConns(m.options.MaxIdles)
	con.SetMaxOpenConns(m.options.MaxOpens)
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
	if m.options.Session != "" {
		delete(_mssqlClients, m.options.Session)
	}
	return nil
}
