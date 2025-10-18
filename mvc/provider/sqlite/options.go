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
	"fmt"

	"github.com/astaxie/beego"
	"github.com/wengoldx/xcore/utils"
)

const (
	_sqliteOptionDBFile = "%s::database" // Configs key of sqlite database filepath.
	_sqliteOptionMemory = "%s::memory"   // Configs key of sqlite database memory mode.
)

// Sqlite client options.
type Options struct {
	Session  string // Session name for load options from app.conf file.
	Database string // Database filepath to connect with, not used for memory database.
	IsMemory bool   // Indicate the sqlite database whether on memory mode.
}

// Create a Options with default values.
func DefaultOptions() Options {
	return Options{Session: _sqliteDefSession, IsMemory: false}
}

// Load Sqlite options from app.conf configs file,
// By default return the configs from 'sqlite' session.
//
// The app.conf configs like:
//
//	; Sqlite database options.
//	[sqlite]
//	database = "sample.db"  ; only for file database.
//	memory = false          ; set true for memory sqlite database.
func LoadOptions(session ...string) Options {
	s := utils.Variable(session, _sqliteDefSession)
	opts := DefaultOptions()
	opts.Database = beego.AppConfig.String(fmt.Sprintf(_sqliteOptionDBFile, s))
	opts.IsMemory = beego.AppConfig.DefaultBool(fmt.Sprintf(_sqliteOptionMemory, s), false)
	return opts
}

// The setter for set Options fields.
type Option func(*Sqlite)

// Specify the session name.
func WithSession(session string) Option {
	return func(m *Sqlite) { m.options.Session = session }
}

// Specify the Sqlite database to assess.
func WithDatabase(database string) Option {
	return func(m *Sqlite) { m.options.Database = database }
}

// Specify the Sqlite database on memory mode.
func WithIsMemory(ismemory bool) Option {
	return func(m *Sqlite) { m.options.IsMemory = ismemory }
}
