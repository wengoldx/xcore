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
	// _ "github.com/go-sql-driver/mysql"   // use for mysql
	//
	// ----------------------------------------
)

// MyConfs mysql database connect configs.
type MyConfs struct {
	Host string // Database host address and port
	User string // Database connect auth user
	Pwd  string // Database connect auth password
	Name string // Target database name
}

// MySQL database configs
const (
	mysqlConfigUser = "%s::user" // configs key of mysql database user
	mysqlConfigPwd  = "%s::pwd"  // configs key of mysql database password
	mysqlConfigHost = "%s::host" // configs key of mysql database host and port
	mysqlConfigName = "%s::name" // configs key of mysql database name

	// Mysql Server database source name for local connection
	mysqldsnLocal = "%s:%s@/%s?charset=%s"

	// Mysql Server database source name for tcp connection
	mysqldsnTcp = "%s:%s@tcp(%s)/%s?charset=%s"
)

// WingHelper content provider to hold mysql database connections,
// the WingHelper.Conn pointer is nil before mvc.OpenMySQL() called.
var WingHelper = &WingProvider{}

// Cache all mysql providers into pool for multiple databases server connect.
var connPool = make(map[string]*WingProvider)

// connectMySQL connect target mysql database with configs.
func connectMySQL(charset, session string, c *MyConfs) error {
	dsn := ""
	if len(c.Host) > 0 /* check database host whether using TCP to connect */ {
		// conntect with remote host database server
		dsn = fmt.Sprintf(mysqldsnTcp, c.User, c.Pwd, c.Host, c.Name, charset)
	} else {
		// just connect local database server
		dsn = fmt.Sprintf(mysqldsnLocal, c.User, c.Pwd, c.Name, charset)
	}
	logger.I("Open MySQL from session:", session)

	// open and connect database
	con, err := sql.Open("mysql", dsn)
	if err != nil {
		return err
	}

	// check database validable
	if err = con.Ping(); err != nil {
		return err
	}

	con.SetMaxIdleConns(100)
	con.SetMaxOpenConns(100)
	con.SetConnMaxLifetime(28740)
	connPool[session] = &WingProvider{con}
	return nil
}

// MySqlConfigs read mysql database params from config file.
func MySqlConfigs(session string) *MyConfs {
	return &MyConfs{
		User: beego.AppConfig.String(fmt.Sprintf(mysqlConfigUser, session)),
		Pwd:  beego.AppConfig.String(fmt.Sprintf(mysqlConfigPwd, session)),
		Host: beego.AppConfig.String(fmt.Sprintf(mysqlConfigHost, session)),
		Name: beego.AppConfig.String(fmt.Sprintf(mysqlConfigName, session)),
	}
}

// OpenMySQL connect database and check ping result, the connection holded
// by mvc.WingHelper object if signle connect, or cached connections in connPool map
// if multiple connect and select same one by given sessions of input params.
// the datatable charset maybe 'utf8' or 'utf8mb4' same as database set.
//
// `USAGE`
//
// you must config database params in /conf/app.config file as follows
//
// ---
//
// #### Case 1 : For signle connect on prod mode.
//
//	[mysql]
//	host = "127.0.0.1:3306"
//	name = "sampledb"
//	user = "root"
//	pwd  = "123456"
//
// #### Case 2 : For signle connect on dev mode.
//
//	[mysql-dev]
//	host = "127.0.0.1:3306"
//	name = "sampledb"
//	user = "root"
//	pwd  = "123456"
//
// #### Case 3 : For both dev and prod mode, you can config all of up cases.
//
// #### Case 4 : For multi-connections to set custom session keywords.
//
//	[mysql-a]
//	... same as use Case 1.
//
//	[mysql-a-dev]
//	... same as use Case 2.
//
//	[mysql-x]
//	... same as use Case 1.
//
//	[mysql-x-dev]
//	... same as use Case 2.
func OpenMySQL(charset string, sessions ...string) error {
	if len(sessions) == 0 {
		sessions = []string{"mysql"}
	}

	// connect all mysql from sessions
	for _, session := range sessions {
		// combine develop session key on dev mode
		if beego.BConfig.RunMode == "dev" {
			session = session + "-dev"
		}

		// load configs by session key
		confs := MySqlConfigs(session)
		if confs.User == "" || confs.Pwd == "" || confs.Name == "" {
			return invar.ErrInvalidConfigs
		}

		if err := connectMySQL(charset, session, confs); err != nil {
			return err
		}
	}

	// using the first connection as primary helper
	WingHelper = Select(sessions[0])
	return nil
}

// OpenMySQL2 connect database with given configs, not from app.conf file.
func OpenMySQL2(charset string, confs *MyConfs) error {
	if confs.User == "" || confs.Pwd == "" || confs.Name == "" {
		return invar.ErrInvalidConfigs
	}

	session := "mysql"
	if err := connectMySQL(charset, session, confs); err != nil {
		return err
	}

	WingHelper = Select(session, true)
	return nil
}

// Select mysql Connection by request key words
// if mode is dev, the key will auto splice '-dev'
func Select(session string, fix ...bool) *WingProvider {
	auto := !(len(fix) > 0 && fix[0])
	if auto && beego.BConfig.RunMode == "dev" {
		session = session + "-dev"
	}
	return connPool[session]
}
