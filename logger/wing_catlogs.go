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
	"strings"

	"github.com/astaxie/beego/logs"
)

// Object logger with optional settings to format output logs.
type wingLogger struct {
	cat string // Category mark output as [CAT] in logs
	tag string // Indicater key output as 'Tag:' prefix in logs
}

// Create a logger instance to output logs with ' [CAT] ' mark, notice
// that the mark will auto append '[]' and change to upper strings.
//
//	---------------------------------------------------------------------------
//	2023/05/31 10:56:36.609 [I] [code_file.go:89] [CAT] FuncName() ...
//	---------------------------------------------------------------------------
func CatLogger(cat string) *wingLogger {
	return NewLogger(cat, "")
}

// Create a logger instance to output logs with ' Tag:' perfix in logs message,
// notice that the tag will auto tail ':' if unexist the char.
//
//	---------------------------------------------------------------------------
//	2023/05/31 10:56:36.609 [I] [code_file.go:89] FuncName() Tag: xxx ...
//	---------------------------------------------------------------------------
func TagLogger(tag string) *wingLogger {
	return NewLogger("", tag)
}

// Create a logger instance to output logs with ' [CAT] ' mark and ' Tag:' perfix
// if set any string value, notice that the mark will auto append '[]' and change
// to upper strings, the tag will tail ':' if unexist end of target key.
//
//	---------------------------------------------------------------------------
//	2023/05/31 10:56:36.609 [I] [code_file.go:89] [CAT] FuncName() Tag: xxx ...
//	---------------------------------------------------------------------------
//
//	Call logger.CatLogger() to create logger output category mark only.
//	Call logger.TagLogger() to create logger output target perfix key only.
func NewLogger(cat, tag string) *wingLogger {
	cat, tag = strings.TrimSpace(cat), strings.TrimSpace(tag)
	if cat != "" {
		cat = "[" + strings.ToUpper(cat) + "] "
	}

	if tag != "" {
		if !strings.HasSuffix(tag, ":") {
			tag = tag + ":"
		}
		tag = " " + tag // ensure formated as ' Tag:'
	}
	return &wingLogger{cat: cat, tag: tag}
}

// E logs a message at error level.
func (l *wingLogger) E(v ...any) {
	logs.Error(logFormatString(len(v), l.cat, l.tag), v...)
	logs.SetPrefix("")
}

// W logs a message at warning level.
func (l *wingLogger) W(v ...any) {
	logs.Warn(logFormatString(len(v), l.cat, l.tag), v...)
	logs.SetPrefix("")
}

// I logs a message at info level.
func (l *wingLogger) I(v ...any) {
	logs.Info(logFormatString(len(v), l.cat, l.tag), v...)
	logs.SetPrefix("")
}

// D logs a message at debug level.
func (l *wingLogger) D(v ...any) {
	logs.Debug(logFormatString(len(v), l.cat, l.tag), v...)
	logs.SetPrefix("")
}

// E logs a message at error level.
func (l *wingLogger) Ef(f string, v ...any) {
	logs.Error(f+logFormatString(0, l.cat, l.tag), v...)
	logs.SetPrefix("")
}

// W logs a message at warning level.
func (l *wingLogger) Wf(f string, v ...any) {
	logs.Warn(f+logFormatString(0, l.cat, l.tag), v...)
	logs.SetPrefix("")
}

// I logs a message at info level.
func (l *wingLogger) If(f string, v ...any) {
	logs.Info(f+logFormatString(0, l.cat, l.tag), v...)
	logs.SetPrefix("")
}

// D logs a message at debug level.
func (l *wingLogger) Df(f string, v ...any) {
	logs.Debug(f+logFormatString(0, l.cat, l.tag), v...)
	logs.SetPrefix("")
}
