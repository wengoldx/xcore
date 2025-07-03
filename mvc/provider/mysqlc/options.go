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
	"fmt"
	"time"

	"github.com/astaxie/beego"
	"github.com/wengoldx/xcore/utils"
)

const (
	_mysqlOptionUser = "%s::user" // configs key of mysql database user
	_mysqlOptionPwd  = "%s::pwd"  // configs key of mysql database password
	_mysqlOptionHost = "%s::host" // configs key of mysql database host and port
	_mysqlOptionName = "%s::name" // configs key of mysql database name
)

// MySQL client options.
type Options struct {
	Session     string        // Session name for load options from app.conf file.
	Host        string        // Database host address and port.
	User        string        // Database connect auth user.
	Password    string        // Database connect auth password.
	Database    string        // Database name to connect with.
	Charset     string        // Database charset, one of 'utf8', 'utf8mb4'...
	MaxIdles    int           // Maximums idle connect chains, default 100.
	MaxOpens    int           // Maximums opening connections, default 100.
	MaxLifetime time.Duration // Maximums lifetime of connection, default 28740s.
}

// Create a Options with default values.
func DefaultOptions(session string) Options {
	return Options{
		Session:     session,
		Charset:     "utf8mb4",
		MaxIdles:    100,
		MaxOpens:    100,
		MaxLifetime: 28740,
	}
}

// Load MySQL options from app.conf configs file.
//
// By default return the configs from 'mysql' session on prod mode,
// or from 'mysql-dev' session on dev mode.
//
// The app.conf configs like:
//
//	; MySQl configs for prod mode.
//	[mysql]
//	host = "192.168.100.102:3306"
//	name = "sampledb"
//	user = "root"
//	pwd  = "123456"
//
//	; MySQl configs for dev mode.
//	[mysql-dev]
//	host = "127.0.0.1:3306"
//	name = "sampledb"
//	user = "root"
//	pwd  = "123456"
func LoadOptions(session ...string) Options {
	s := utils.VarString(session, _mysqlDriver)
	opts := DefaultOptions(s)

	// auto append suffix for dev mode.
	if beego.BConfig.RunMode == "dev" {
		s += "-dev"
	}

	opts.User = beego.AppConfig.String(fmt.Sprintf(_mysqlOptionUser, s))
	opts.Password = beego.AppConfig.String(fmt.Sprintf(_mysqlOptionPwd, s))
	opts.Host = beego.AppConfig.String(fmt.Sprintf(_mysqlOptionHost, s))
	opts.Database = beego.AppConfig.String(fmt.Sprintf(_mysqlOptionName, s))
	return opts
}

// The setter for set Options fields.
type Option func(*MySQL)

// Specify the session name.
func WithSession(session string) Option {
	return func(m *MySQL) { m.options.Session = session }
}

// Specify the MySQL server host and port.
func WithHost(host string) Option {
	return func(m *MySQL) { m.options.Host = host }
}

// Specify the MySQL connect user.
func WithUser(user string) Option {
	return func(m *MySQL) { m.options.User = user }
}

// Specify the MySQL connect password.
func WithPassword(password string) Option {
	return func(m *MySQL) { m.options.Password = password }
}

// Specify the MySQL database to assess.
func WithDatabase(database string) Option {
	return func(m *MySQL) { m.options.Database = database }
}

// Specify the MySQL database charset.
func WithCharset(charset string) Option {
	return func(m *MySQL) { m.options.Charset = charset }
}

// Specify the maximums idle connect chains.
func WithMaxIdles(idles int) Option {
	return func(m *MySQL) { m.options.MaxIdles = idles }
}

// Specify the maximums opening connections.
func WithMaxOpens(opens int) Option {
	return func(m *MySQL) { m.options.MaxOpens = opens }
}

// Specify the maximums lifetime of connection.
func WithMaxLifetime(lifetime time.Duration) Option {
	return func(m *MySQL) { m.options.MaxLifetime = lifetime }
}
