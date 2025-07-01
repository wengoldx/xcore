// Copyright (c) 2018-Now Dunyu All Rights Reserved.
//
// Author      : https://www.wengold.net
// Email       : support@wengold.net
//
// Prismy.No | Date       | Modified by. | Description
// -------------------------------------------------------------------
// 00001       2019/05/22   yangping       New version
// -------------------------------------------------------------------

package tester

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/astaxie/beego"
	"github.com/wengoldx/xcore/utils"
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
				if forms := Wter.parseForms(params.(TestForm)); forms != "" {
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
			req.Header.Add("Author", Wter.Author)
			req.Header.Add("Token", Wter.getToken(uid))
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

// Restful api tester runtime configs.
type tester struct {
	tokens   map[string]string // Auth token of test user, format as {uuid:token}.
	Author   string            // Author header, such as 'WENGOLD-V1.1', 'WENGOLD-V1.2', 'WENGOLD-V2.0'
	TokenApi string            // Rest4 API to get user token, like 'http://192.168.1.100:8000/server/token?id=%s'
	User     string            // User uuid for testing

	// Env params for testing, set param by code 'Wter.Env["param-name"] = param-vaule'
	// and used as 'value := Wter.Env["param-name"].(string)' to get string value.
	Env map[string]any
}

// Global Restful API tester signleton.
//
//	USAGE: Init mvc.Wter configs before use it as follow:
//
//	func init() {
//		// logger.SilentLoggers() // silent logger if comment out.
//		mvc.Wter.Author = "WENGOLD-V2.0"
//		mvc.Wter.TokenApi = "http://192.168.1.100:8000/server/token?id=%s"
//		mvc.Wter.User = "12345678"
//	}
var Wter = &tester{
	tokens: make(map[string]string),
	Env:    make(map[string]any),
}

// Transform url params map to url.Values for http GET method.
func (t *tester) parseForms(params TestForm) string {
	forms := url.Values{}
	for param, value := range params {
		forms[param] = []string{fmt.Sprintf("%v", value)}
	}
	return forms.Encode()
}

// Get test token of target user from cachs map, or request by restful api.
func (t *tester) getToken(uid string) string {
	if token, ok := t.tokens[uid]; ok {
		return token
	} else if t.TokenApi == "" {
		return ""
	}

	// request user token from remote server by given debug api.
	if token, err := utils.HttpUtils.GString(t.TokenApi, uid); err == nil {
		t.tokens[uid] = token
		return token
	}
	return ""
}
