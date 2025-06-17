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

// -------------------------------------------------------------------
// Define methods as database helper.
// -------------------------------------------------------------------

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
type rest4Tester struct {
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
var Wter = &rest4Tester{
	tokens: make(map[string]string),
	Env:    make(map[string]any),
}

// Transform url params map to url.Values for http GET method.
func (t *rest4Tester) parseForms(params TestForm) string {
	forms := url.Values{}
	for param, value := range params {
		forms[param] = []string{fmt.Sprintf("%v", value)}
	}
	return forms.Encode()
}

// Get test token of target user from cachs map, or request by restful api.
func (t *rest4Tester) getToken(uid string) string {
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

// -------------------------------------------------------------------
// Define methods as database helper.
// -------------------------------------------------------------------

// Database helper for unit test.
type utestHelper struct {
	WingProvider

	tag   string  // target output field name, such as id, uid...
	where KValues // where conditions field-value pairs, default empty.
	order string  // order by field name, auto append 'ORDER BY' perfix.
	desc  string  // order type, default DESC, one of 'DESC', 'ASC'.
	limit int     // limit number, auto append 'LIMIT' perfix.
}

// Create a data helper to query table datas.
func UTest() *utestHelper {
	return &utestHelper{
		WingProvider: *WingHelper, desc: "DESC", limit: 0,
	}
}

// Typed key-value map.
type KValues map[string]any

// Unit test helper options setter.
type Option func(t *utestHelper)

func WithTag(tag string) Option                     { return func(u *utestHelper) { u.tag = tag } }
func WithWhere(where KValues) Option                { return func(u *utestHelper) { u.where = where } }
func WithLimit(limit int) Option                    { return func(u *utestHelper) { u.limit = limit } }
func WithOrder(order string, desc ...string) Option { return func(u *utestHelper) { u.order = order } }

// Apply all unit test helper options settngs.
func (t *utestHelper) applyOptions(options ...Option) {
	for _, option := range options {
		option(t)
	}
}

// Format wheres condition to string, and with values.
func (t *utestHelper) formatWheres() (string, []any) {
	wheres, wvals := "", []any{}
	if len(t.where) > 0 {
		wkeys := []string{}
		for wkey, wval := range t.where {
			wkeys = append(wkeys, wkey)
			wvals = append(wvals, wval)
		}
		wheres = "WHERE " + strings.Join(wkeys, " AND ")
	}
	return wheres, wvals
}

// Format order by condition to string.
func (t *utestHelper) formatOrder() string {
	if t.order != "" {
		return fmt.Sprintf("ORDER BY %s %s", t.order, t.desc)
	}
	return ""
}

// Format limit condition to string.
func (t *utestHelper) formatLimit() string {
	if t.limit > 0 {
		return fmt.Sprintf("LIMIT %d", t.limit)
	}
	return ""
}

// Format multiple value as IN condition in where.
func (t *utestHelper) formatWhereIns(tag string, ins []string) string {
	if tag != "" && len(ins) > 0 {
		return tag + " " + t.JoinStrings("IN (%s)", ins)
	}
	return ""
}

// Format multiple value as insert sql string and values.
func (t *utestHelper) formatInsertValus(ins KValues) (string, string, []any) {
	fields, params, vals := "", "", []any{}
	if len(ins) > 0 {
		keys, args := []string{}, []string{}
		for key, val := range ins {
			keys = append(keys, key)
			vals = append(vals, val)
			args = append(args, "?")
		}
		fields = strings.Join(keys, ", ")
		params = strings.Join(args, ", ")
	}
	return fields, params, vals
}

const (
	_sql_ut_get = "SELECT %s FROM %s %s %s %s" // table, where, order, limit
	_sql_ut_add = "INSERT %s (%s) VALUE (%s)"  // table, target fields, (?,?,...)
	_sql_ut_del = "DELETE FROM %s %s %s"       // table, where, in (%s)
)

// Get target field string value from given table, or with options.
//
//	SQL: SELECT %s FROM %s WHERE %s ORDER BY %s DESC LIMIT 1
//	            ^       ^  -------^ ----------^----- ------^
//	          tag   table     where       order        limit
//
//	@param table   Target table name.
//	@param tag     Target filed name to output query result.
//	@param where   Where conditions, the map key must like 'created>=?'.
//	@param options Setter for set order by field name, limit number.
//	@return out - any Output result value, like &int, &int64, &float64, &string...
func (t *utestHelper) Target(table, tag string, where KValues, out any, options ...Option) error {
	t.tag, t.where = tag, where
	t.applyOptions(options...)

	wheres, values := t.formatWheres() // format wheres sting and input values.
	order := t.formatOrder()           // format order by string.
	limit := t.formatLimit()           // format limit string.

	// SELECT tag FROM table wheres order limit
	query := fmt.Sprintf(_sql_ut_get, t.tag, table, wheres, order, limit)
	return t.One(query, func(rows *sql.Rows) error {
		if err := rows.Scan(out); err != nil {
			return err
		}
		return nil
	}, values...)
}

// Get last id from given table, or with options.
//
//	USAGE:
//
//	(1). mvc.UTest().LastID("account")
//	-> SELECT id FROM account ORDER BY id DESC LIMIT 1
//
//	(2). mvc.UTest().LastID("account", mvc.WithID("userid"))
//	-> SELECT userid FROM account ORDER BY id DESC LIMIT 1
//
//	(3). mvc.UTest().LastID("account", mvc.WithWhere({"acc=?" : "nickname"}))
//	-> SELECT id FROM account WHERE acc='nickname' ORDER BY id DESC LIMIT 1
//
//	(4). mvc.UTest().LastID("account", mvc.WithOrder("created"))
//	-> SELECT id FROM account ORDER BY created DESC LIMIT 1
//
//	@param table   Target table name.
//	@param options Setter for set id and order field name, where conditions, limit number.
//	@return id - Last record id.
//
//	See Target() for more sql query format infos.
func (t *utestHelper) LastID(table string, options ...Option) (id int64, e error) {
	t.tag, t.order, t.limit = "id", "id", 1
	return id, t.Target(table, t.tag, t.where, &id, options...)
}

// Get last uid from given table, or with options.
//
//	See Target(), LastID for more query format or usage infos.
func (t *utestHelper) LastUID(table string, options ...Option) (uid string, e error) {
	t.tag, t.order, t.limit = "uid", "created", 1
	return uid, t.Target(table, t.tag, t.where, &uid, options...)
}

// Insert a record into target table by given values.
//
//	SQL: INSERT %s (%s) VALUE (%s)
//	             ^   ^          ^
//	         table   tags    args
//
//	@param table Target table name.
//	@param ins   Target fields name and insert values.
func (t *utestHelper) Add(table string, ins KValues) error {
	fields, args, values := t.formatInsertValus(ins)
	return t.Execute(fmt.Sprintf(_sql_ut_add, table, fields, args), values...)
}

// Insert a record into target table and return the inserted id.
//
//	See Add() for insert without record id.
func (t *utestHelper) AddWithID(table string, ins KValues) (int64, error) {
	fields, args, values := t.formatInsertValus(ins)
	return t.Insert(fmt.Sprintf(_sql_ut_add, table, fields, args), values...)
}

// Deleta records by given where conditions.
//
//	SQL: DELETE FROM %s WHERE %s
//	                  ^ -------^
//	              table    where
//
//	@param table Target table name.
//	@param where Field name as where condition like: field = value.
func (t *utestHelper) DeleteBy(table string, where KValues) error {
	t.where = where
	wheres, values := t.formatWheres()
	return t.Execute(fmt.Sprintf(_sql_ut_del, table, wheres, ""), values...)
}

// Deleta records by given where condition and target fields vlaues.
//
//	SQL: DELETE FROM %s WHERE %s AND %s IN (%s)
//	                  ^ -------^---- ^-------^-
//	              table    where     tag     values
//
//	@param table Target table name.
//	@param where Field name as where condition like: field IN (values).
//	@param value Where condition values to query.
func (t *utestHelper) DeleteIns(table, tag string, ins []string, options ...Option) error {
	t.applyOptions(options...)
	wheres, values := t.formatWheres()     // format wheres sting and input values.
	instring := t.formatWhereIns(tag, ins) // format in values.
	if wheres != "" && instring != "" {
		wheres += " AND "
	}
	return t.Execute(fmt.Sprintf(_sql_ut_del, table, wheres, instring), values...)
}
