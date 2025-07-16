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
	"fmt"

	"github.com/astaxie/beego"
	"github.com/wengoldx/xcore/utils"
)

const (
	_mssqlOptionUser    = "%s::user"    // configs key of mssql database user
	_mssqlOptionPwd     = "%s::pwd"     // configs key of mssql database password
	_mssqlOptionHost    = "%s::host"    // configs key of mssql database server host
	_mssqlOptionPort    = "%s::port"    // configs key of mssql database port
	_mssqlOptionName    = "%s::name"    // configs key of mssql database name
	_mssqlOptionTimeout = "%s::timeout" // configs key of mssql database connect timeout
)

// MSSQL client options.
type Options struct {
	Session  string // Session name for load options from app.conf file.
	Host     string // Database host address.
	Port     int    // Database server port.
	User     string // Database connect auth user.
	Password string // Database connect auth password.
	Database string // Database name to connect with.
	Timeout  int    // Database connect timeout.
	MaxIdles int    // Maximums idle connect chains, default 100.
	MaxOpens int    // Maximums opening connections, default 100.
}

// Create a Options with default values.
func DefaultOptions(session string) Options {
	return Options{
		Session:  session,
		Host:     "127.0.0.1",
		Port:     1433,
		Timeout:  30,
		MaxIdles: 100,
		MaxOpens: 100,
	}
}

// Load MSSQL options from app.conf configs file.
//
// By default return the configs from 'mssql' session on prod mode,
// or from 'mssql-dev' session on dev mode.
//
// The app.conf configs like:
//
//	; MSSQL configs for prod mode.
//	[mssql]
//	host    = "192.168.100.102"
//	port    = 1433
//	name    = "sampledb"
//	user    = "sa"
//	pwd     = "123456"
//	timeout = 30
//
//	; MSSQl configs for dev mode.
//	[mssql-dev]
//	host    = "127.0.0.1"
//	port    = 1433
//	name    = "sampledb"
//	user    = "sa"
//	pwd     = "123456"
//	timeout = 30
func LoadOptions(session ...string) Options {
	s := utils.Variable(session, _mssqlDriver)
	opts := DefaultOptions(s)

	// auto append suffix for dev mode.
	if beego.BConfig.RunMode == "dev" {
		s += "-dev"
	}

	opts.User = beego.AppConfig.String(fmt.Sprintf(_mssqlOptionUser, s))
	opts.Password = beego.AppConfig.String(fmt.Sprintf(_mssqlOptionPwd, s))
	opts.Host = beego.AppConfig.DefaultString(fmt.Sprintf(_mssqlOptionHost, s), "127.0.0.1")
	opts.Port = beego.AppConfig.DefaultInt(fmt.Sprintf(_mssqlOptionPort, s), 1433)
	opts.Database = beego.AppConfig.String(fmt.Sprintf(_mssqlOptionName, s))
	opts.Timeout = beego.AppConfig.DefaultInt(fmt.Sprintf(_mssqlOptionTimeout, s), 30) // seconds
	return opts
}

// The setter for set Options fields.
type Option func(*MSSQL)

// Specify the session name.
func WithSession(session string) Option {
	return func(m *MSSQL) { m.options.Session = session }
}

// Specify the MSSQL server host.
func WithHost(host string) Option {
	return func(m *MSSQL) { m.options.Host = host }
}

// Specify the MSSQL server port.
func WithPort(port int) Option {
	return func(m *MSSQL) { m.options.Port = port }
}

// Specify the MSSQL connect user.
func WithUser(user string) Option {
	return func(m *MSSQL) { m.options.User = user }
}

// Specify the MSSQL connect password.
func WithPassword(password string) Option {
	return func(m *MSSQL) { m.options.Password = password }
}

// Specify the MSSQL database to assess.
func WithDatabase(database string) Option {
	return func(m *MSSQL) { m.options.Database = database }
}

// Specify the connect timeout duration seconds.
func WithTimeout(timeout int) Option {
	return func(m *MSSQL) { m.options.Timeout = timeout }
}

// Specify the maximums idle connect chains.
func WithMaxIdles(idles int) Option {
	return func(m *MSSQL) { m.options.MaxIdles = idles }
}

// Specify the maximums opening connections.
func WithMaxOpens(opens int) Option {
	return func(m *MSSQL) { m.options.MaxOpens = opens }
}
