// Copyright (c) 2018-Now Dunyu All Rights Reserved.
//
// Author      : https://www.wengold.net
// Email       : support@wengold.net
//
// Prismy.No | Date       | Modified by. | Description
// -------------------------------------------------------------------
// 00001       2019/05/22   yangping       New version
// -------------------------------------------------------------------

package tu

import (
	"fmt"
	"os"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/wengoldx/xcore/mvc/provider/mysql"
	"github.com/wengoldx/xcore/mvc/provider/sqlite"
	"gopkg.in/ini.v1"
)

// Set 'dev' runmode and fix debug logger.
func UseDebugLogger() {
	beego.BConfig.RunMode = "dev"
	logs.SetLevel(beego.LevelDebug)
}

/* ------------------------------------------------------------------- */
/* For MySQL                                                           */
/* ------------------------------------------------------------------- */

// Open database for testing by given .test env file.
//
//	opts := mysql.Options{
//		Host: "localhost:3306", Database: "testdb",
//		User: "user", Password: "****",
//	}
func OpenTestMysql(env string) {
	UseDebugLogger()

	opts := readMysqlEnv(env)
	if opts.Host == "" {
		panic("Empty database host !!")
	} else if err := mysql.OpenWithOptions(opts, "utf8mb4"); err != nil {
		panic("Failed Open test database: " + err.Error())
	}
	fmt.Println("[I] Opened test database...")
}

// Read test env configs from .test file.
//
//	[DATABASE]
//	Host="localhost:3306"
//	Database="testdb"
//	User="user"
//	Password="****"
func readMysqlEnv(env string) mysql.Options {
	opts := mysql.Options{}
	info, err := os.Stat(env)
	if err == nil && !info.IsDir() {
		if cfg, err := ini.Load(env); err != nil {
			panic("Failed read test env:" + err.Error())
		} else if section := cfg.Section("DATABASE"); section != nil {
			opts.Host = section.Key("Host").String()
			opts.Database = section.Key("Database").String()
			opts.User = section.Key("User").String()
			opts.Password = section.Key("Password").String()
		}
	}
	return opts
}

/* ------------------------------------------------------------------- */
/* For Sqlite                                                          */
/* ------------------------------------------------------------------- */

// Open database for testing by given .test env file.
//
//	opts := sqlite.Options{
//		Database: "testdb", Memory: false,
//	}
func OpenTestSqlite(env string) {
	UseDebugLogger()

	opts := readSqliteEnv(env)
	if err := sqlite.OpenWithOptions(opts); err != nil {
		panic("Failed Open test database: " + err.Error())
	}
	fmt.Println("[I] Opened test database...")
}

// Read test env configs from .test file.
//
//	[DATABASE]
//	Database="testdb"
//	Memory=false
func readSqliteEnv(env string) sqlite.Options {
	opts := sqlite.Options{}
	info, err := os.Stat(env)
	if err == nil && !info.IsDir() {
		if cfg, err := ini.Load(env); err != nil {
			panic("Failed read test env:" + err.Error())
		} else if section := cfg.Section("DATABASE"); section != nil {
			opts.Database = section.Key("Database").String()
			opts.IsMemory, _ = section.Key("Memory").Bool()
		}
	}
	return opts
}
