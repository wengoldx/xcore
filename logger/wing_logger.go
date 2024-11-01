// Copyright (c) 2018-Now Dunyu All Rights Reserved.
//
// Author      : https://www.wengold.net
// Email       : support@wengold.net
//
// Prismy.No | Date       | Modified by. | Description
// -------------------------------------------------------------------
// 00001       2019/05/22   yangping       New version
// -------------------------------------------------------------------

package logger

import (
	"runtime"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

const (
	logConfigLevel   = "logger::level"   // configs key of logger level
	logConfigMaxDays = "logger::maxdays" // configs key of logger max days

	// LevelDebug debug level of logger
	LevelDebug = "debug"

	// LevelInfo info level of logger
	LevelInfo = "info"

	// LevelWarn warn level of logger
	LevelWarn = "warn"

	// LevelError error level of logger
	LevelError = "error"
)

// init initialize app logger
//
// `NOTICE` : you must config logger params in /conf/app.config file as:
//
// ---
//
//	[logger]
//	level = "debug"
//	maxdays = "7"
//
// ---
//
// - the level values range in : [debug, info, warn, error], default is info.
//
// - maxdays is the max days to hold logs cache, default is 7 days.
func init() {
	setupFileLogger()
	logs.SetLogFuncCall(true) // use the default func depth
	logs.Async(3)             // allow 3 asynchronous chanels

	// set application logger level
	switch beego.AppConfig.String(logConfigLevel) {
	case LevelDebug:
		logs.SetLevel(beego.LevelDebug)
	case LevelInfo:
		logs.SetLevel(beego.LevelInformational)
	case LevelWarn:
		logs.SetLevel(beego.LevelWarning)
	case LevelError:
		logs.SetLevel(beego.LevelError)
	}
}

// setupFileLogger init and set logger output to file
func setupFileLogger() {
	app := beego.BConfig.AppName
	if app == "" || app == "beego" {
		return
	}

	maxdays := beego.AppConfig.String(logConfigMaxDays)
	if maxdays == "" {
		maxdays = "7"
	}
	config := "{\"filename\":\"logs/" + app + ".log\", \"daily\":true, \"maxdays\":" + maxdays + "}"
	logs.SetLogger(logs.AdapterFile, config)
}

// Return log format string like '%v %v %v' when n is 3.
// Here will set logger perfix as outside caller function name.
func logFormatString(n int) string {

	// append runtime calling function name as logger prefix, it format like :
	// -----------------------------------------------------------------------
	// 2023/05/31 10:56:36.609 [I] [code_file.go:89] FuncName() Log messages
	// -----------------------------------------------------------------------

	/* Fixed the call skipe on 2 to filter inner functions name */
	if pc, _, _, ok := runtime.Caller(2); ok {
		if funcptr := runtime.FuncForPC(pc); funcptr != nil {
			if funname := funcptr.Name(); funname != "" {
				fns := strings.SplitAfter(funname, ".")
				logs.SetPrefix(fns[len(fns)-1] + "()")
			}
		}
	}

	return strings.Repeat("%v ", n)
}

// SetOutputLogger close console logger on prod mode and only remain file logger.
func SetOutputLogger() {
	if beego.BConfig.RunMode != "dev" && GetLevel() != LevelDebug {
		beego.BeeLogger.DelLogger(logs.AdapterConsole)
	}
}

// GetLevel return current logger output level
func GetLevel() string {
	switch beego.BeeLogger.GetLevel() {
	case beego.LevelDebug:
		return LevelDebug
	case beego.LevelInformational:
		return LevelInfo
	case beego.LevelWarning:
		return LevelWarn
	case beego.LevelError:
		return LevelError
	default:
		return ""
	}
}

// EM logs a message at emergency level.
func EM(v ...any) {
	logs.Emergency(logFormatString(len(v)), v...)
}

// AL logs a message at alert level.
func AL(v ...any) {
	logs.Alert(logFormatString(len(v)), v...)
}

// CR logs a message at critical level.
func CR(v ...any) {
	logs.Critical(logFormatString(len(v)), v...)
}

// E logs a message at error level.
func E(v ...any) {
	logs.Error(logFormatString(len(v)), v...)
}

// W logs a message at warning level.
func W(v ...any) {
	logs.Warn(logFormatString(len(v)), v...)
}

// N logs a message at notice level.
func N(v ...any) {
	logs.Notice(logFormatString(len(v)), v...)
}

// I logs a message at info level.
func I(v ...any) {
	logs.Info(logFormatString(len(v)), v...)
}

// D logs a message at debug level.
func D(v ...any) {
	logs.Debug(logFormatString(len(v)), v...)
}
