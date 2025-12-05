// Copyright (c) 2018-Now Dunyu All Rights Reserved.
//
// Author      : https://www.wengold.net
// Email       : support@wengold.net
//
// Prismy.No | Date       | Modified by. | Description
// -------------------------------------------------------------------
// 00001       2019/05/22   yangping       New version
// -------------------------------------------------------------------

package xt

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

/* ------------------------------------------------------------------- */
/* For Test Case & Multiples                                           */
/* ------------------------------------------------------------------- */

// Test case datas for multiple testing.
type TestCase struct {
	Case   string
	Params any
	Want   any
}

// Test handler for function test.
type TestHandler func(param any) any

// Return test case object for easy multipe testing.
func NewCase(label string, want any, param any) *TestCase {
	return &TestCase{Case: label, Want: want, Params: param}
}

// Test helper to execute multiple cases and check the result.
//
// Example:
//
//	cases := []*wt.TestCase{
//		wt.NewCase("Check 1", "1 \\2\\3", "/  1 /2\\3\\    "),
//		wt.NewCase("Check 2", ".", ""),
//	}
//
//	wt.TestMults(t, cases, func(param any) any {
//		return NormalizePath(param.(string))
//	})
func TestMults(t *testing.T, cases []*TestCase, callback TestHandler) {
	t.Helper()

	LogI("Start Testing Cases...")
	for _, c := range cases {
		if want := callback(c.Params); want != c.Want {
			t.Fatal("Failed, want:", c.Want, "but result:", want)
		}
		LogI("[OK]", c.Case, "-", c.Want)
	}
	LogI("Finished Tests!")
}

/* ------------------------------------------------------------------- */
/* For Test Utils                                                      */
/* ------------------------------------------------------------------- */

func LogI(msg ...any) { fmt.Println(append([]any{"[I]"}, msg...)...) }
func LogE(msg ...any) { fmt.Println(append([]any{"[E]"}, msg...)...) }

// Set 'dev' runmode and fix debug logger.
func UseDebugLogger() {
	beego.BConfig.RunMode = "dev"
	logs.SetLevel(beego.LevelDebug)
}

// Check app whether running on test mode, it just check the .test
// file whether exist in ~/{app}/conf/ folder.
//
//	~/{app}
//	  |- bin
//	  |- conf
//	  |  |- .test  // -> Enable test mode, delete it for disable.
//	  ...
//
// # WARNING:
//	- DO NOT use beego.BConfig.AppName when unexist app.conf!
//	- Return empty when ~/conf/.test unexist.
func GetTestEnv(app string) string {
	length := len(app)
	if pwd, err := os.Getwd(); err == nil && length > 0 {
		if start := strings.Index(pwd, app); start > 0 {
			env := pwd[:start+length] + "/conf/.test"
			if existFile(env) {
				return env
			}
		}
	}
	return ""
}

// Return server root dir: /home/.../{server} on test model.
//
// # WARNING:
//	- Use xt.GetTestEnv() get abstract path and trim conf paths.
//	- Return empty when ~/conf/.test unexist.
func GetServRoot(app string) string {
	if test := GetTestEnv(app); test != "" {
		return strings.TrimSuffix(test, "/conf/.test")
	}
	return ""
}

func existFile(env string) bool {
	if info, err := os.Stat(env); err == nil {
		return !info.IsDir()
	}
	return false
}
