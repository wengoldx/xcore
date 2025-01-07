// Copyright (c) 2018-Now Dunyu All Rights Reserved.
//
// Author      : https://www.wengold.net
// Email       : support@wengold.net
//
// Prismy.No | Date       | Modified by. | Description
// -------------------------------------------------------------------
// 00001       2019/05/22   yangping       New version
// -------------------------------------------------------------------

package mvc

import (
	"database/sql"
	"fmt"

	"github.com/astaxie/beego"
	"github.com/wengoldx/xcore/invar"
	"github.com/wengoldx/xcore/logger"
	// ----------------------------------------
	// NOTIC :
	//
	// import the follow database driver when using WingProvider.
	//
	// _ "github.com/denisenkom/go-mssqldb" // use for sql server 2017 ~ 2019
	//
	// ----------------------------------------
)

// Microsoft SQL Server configs
const (
	mssqlConfigUser = "%s::user"    // configs key of mssql database user
	mssqlConfigPwd  = "%s::pwd"     // configs key of mssql database password
	mssqlConfigHost = "%s::host"    // configs key of mssql database server host
	mssqlConfigPort = "%s::port"    // configs key of mssql database port
	mssqlConfigName = "%s::name"    // configs key of mssql database name
	mssqlConfigTout = "%s::timeout" // configs key of mssql database connect timeout

	// Microsoft SQL Server database source name
	mssqldsn = "server=%s;port=%d;database=%s;user id=%s;password=%s;connection timeout=%d;dial timeout=%d;"
)

// MssqlHelper content provider to hold mssql database connections,
// it will nil before mvc.OpenMssql() called.
var MssqlHelper *WingProvider

// readMssqlCofnigs read mssql database params from config file,
// than verify them if empty.
func readMssqlCofnigs(session string) (string, string, string, int, string, int, error) {
	user := beego.AppConfig.String(fmt.Sprintf(mssqlConfigUser, session))
	pwd := beego.AppConfig.String(fmt.Sprintf(mssqlConfigPwd, session))
	host := beego.AppConfig.DefaultString(fmt.Sprintf(mssqlConfigHost, session), "127.0.0.1")
	port := beego.AppConfig.DefaultInt(fmt.Sprintf(mssqlConfigPort, session), 1433)
	name := beego.AppConfig.String(fmt.Sprintf(mssqlConfigName, session))
	timeout := beego.AppConfig.DefaultInt(fmt.Sprintf(mssqlConfigTout, session), 30) // seconds

	if user == "" || pwd == "" || name == "" {
		return "", "", "", 0, "", 0, invar.ErrInvalidConfigs
	}
	return user, pwd, host, port, name, timeout, nil
}

// OpenMssql connect mssql database and check ping result,
// the connections holded by mvc.MssqlHelper object,
// the charset maybe 'utf8' or 'utf8mb4' same as database set.
//
// `NOTICE`
//
// you must config database params in /conf/app.config file as:
//
// ---
//
// #### Case 1 For connect on prod mode.
//
//	[mssql]
//	host    = "127.0.0.1"
//	port    = 1433
//	name    = "sampledb"
//	user    = "sa"
//	pwd     = "123456"
//	timeout = 30
//
// #### Case 2 For connect on dev mode.
//
//	[mssql-dev]
//	host    = "127.0.0.1"
//	port    = 1433
//	name    = "sampledb"
//	user    = "sa"
//	pwd     = "123456"
//	timeout = 30
//
// #### Case 3 For both dev and prod mode, you can config all of up cases.
func OpenMssql(charset string) error {
	session := "mssql"
	if beego.BConfig.RunMode == "dev" {
		session = session + "-dev"
	}

	user, pwd, server, port, dbn, to, err := readMssqlCofnigs(session)
	if err != nil {
		return err
	} else if to <= 0 { // check connection timeout
		to = 30 // fix the dial timeout over 5s
	}

	driver := "mssql" // mssql for processQueryText=true, sqlserver for false
	dsn := fmt.Sprintf(mssqldsn, server, port, dbn, user, pwd, to, to+5)
	logger.I("Open MSSQL from session:", session)

	// open and connect database
	con, err := sql.Open(driver, dsn)
	if err != nil {
		return err
	}

	// check database validable
	if err = con.Ping(); err != nil {
		return err
	}

	con.SetMaxIdleConns(100)
	con.SetMaxOpenConns(100)
	MssqlHelper = &WingProvider{con}
	return nil
}
