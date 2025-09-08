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
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/wengoldx/xcore/logger"
	"github.com/wengoldx/xcore/mvc/provider/mysql"
	"github.com/wengoldx/xcore/utils"
	"gopkg.in/ini.v1"
)

// Test case datas for multiple testing.
type TestCase struct {
	Case   string // Current excute unit test case.
	User   string // Http request author uuid.
	Want   int    // Http response code for pass current test.
	Params any    // Http request params, body of POST method, form data of GET method.
}

// Url params form datas type.
type TestForm map[string]any

// Return test case object for easy multipe testing.
func NewCase(c string, want int, param any) *TestCase {
	return &TestCase{Case: c, Want: want, Params: param}
}

// Return test case object with authed user for easy multipe testing.
func UserCase(user, c string, want int, param any) *TestCase {
	return &TestCase{Case: c, User: user, Want: want, Params: param}
}

// Multiple testing from given test cases.
func TestMults(t *testing.T, api, method string, cases []*TestCase) {
	t.Helper()
	for _, c := range cases {
		TestMain(t, c.Case, c.User, api, method, c.Want, c.Params)
		time.Sleep(20 * time.Millisecond) // wait 20ms for next
	}
}

// A single testing to simulate send http request and logout test response if exist.
func TestMain(t *testing.T, c, uid, api, method string, want int, params any) {
	resp := httptest.NewRecorder()
	t.Helper()
	t.Run(c, func(t *testing.T) {
		contentType := "application/json"
		url := "/" + beego.BConfig.AppName + "/" + api

		// perpare http request params.
		var requestBody io.Reader
		switch method {
		case http.MethodPost:
			if params != nil && params != struct{}{} {
				paramsJson, _ := json.Marshal(params)
				requestBody = strings.NewReader(string(paramsJson))
			}
		case http.MethodGet:
			contentType = "application/x-www-form-urlencoded"
			if params != nil && params != struct{}{} {
				if forms := _t.parseForms(params.(TestForm)); forms != "" {
					url += "?" + forms
				}
			}
		default:
			t.Fatalf("Unsupport HTTP method: %s for test !!", method)
		}

		// create http request and set auth headers.
		req, _ := http.NewRequest(method, url, requestBody)
		req.Header.Add("Content-Type", contentType)
		if uid != "" {
			req.Header.Add("Author", _t.author)
			req.Header.Add("Token", _t.getToken(uid))
		}

		beego.BeeApp.Handlers.ServeHTTP(resp, req)
		if resp.Code != want {
			t.Fatalf("Unexpected value:%v, want is %v", resp.Code, want)
		}
	})

	// Logout response datas if exist.
	if rst := resp.Body.String(); rst != "" && rst != "<nil>" && rst != "null" {
		t.Log("Test response:", rst)
	}
}

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
//	  |    |- .test  // -> Enable test mode, delete it for disable.
//	  ...
//
// WARNING: DO NOT use beego.BConfig.AppName when unexist app.conf!
func CheckTestMode(app string) string {
	if pwd, err := os.Getwd(); err == nil {
		if length := len(app); length > 0 {
			if start := strings.Index(pwd, app); start > 0 {
				env := pwd[:start+length] + "/conf/.test"
				if utils.IsFile(env) {
					fmt.Println()
					fmt.Println("\t+--------------------+")
					fmt.Println("\t+ LAUNCHED TEST MODE +")
					fmt.Println("\t+--------------------+")
					fmt.Println()
					return env
				}
			}
		}
	}
	return ""
}

// Open database for testing by given .test env file.
//
//	opts := mysql.Options{
//		Host: "localhost:3306", Database: "testdb",
//		User: "user", Password: "****",
//	}
func OpenTestDatabase(env string) {
	UseDebugLogger()

	opts := readTestEnv(env)
	if opts.Host == "" {
		panic("Empty database host !!")
	} else if err := mysql.OpenWithOptions(opts, "utf8mb4"); err != nil {
		panic("Failed Open test database: " + err.Error())
	}
	logger.I("Opened test database...")
}

// Read test env configs from .test file.
//
//	[DATABASE]
//	Host="localhost:3306"
//	Database="testdb"
//	User="user"
//	Password="****"
func readTestEnv(env string) mysql.Options {
	opts := mysql.Options{}
	if !utils.IsFile(env) {
		return opts
	}

	if cfg, err := ini.Load(env); err != nil {
		panic("Failed read test env:" + err.Error())
	} else if section := cfg.Section("DATABASE"); section != nil {
		opts.Host = section.Key("Host").String()
		opts.Database = section.Key("Database").String()
		opts.User = section.Key("User").String()
		opts.Password = section.Key("Password").String()
	}
	return opts
}
