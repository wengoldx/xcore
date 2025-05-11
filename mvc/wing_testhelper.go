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
			req.Header.Add("Author", "WENGOLD-V1.2")
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
	authTokens map[string]string // Auth token of test user, format as {uuid:token}.
	TokenApi   string            // Rest4 API to get user token, like 'http://192.168.1.100:8000/server/token?id=%s'
	User       string            // User uuid for testing

	// Env params for testing, set param by code 'Wter.Env["param-name"] = param-vaule'
	// and used as 'value := Wter.Env["param-name"].(string)' to get string value.
	Env map[string]any
}

// Global Restful API tester signleton.
//
//	USAGE: Init mvc.Wter configs before use it as follow:
//
//	func init() {
//		// logger.SilentLoggers() // slient logger if comment out.
//		mvc.Wter.TokenApi = "http://192.168.1.100:8000/server/token?id=%s"
//		mvc.Wter.User = "12345678"
//	}
var Wter = &rest4Tester{
	authTokens: make(map[string]string),
	Env:        make(map[string]any),
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
	if token, ok := t.authTokens[uid]; ok {
		return token
	} else if t.TokenApi == "" {
		return ""
	}

	// request user token from remote server by given debug api.
	if token, err := utils.HttpGetString(t.TokenApi, uid); err == nil {
		t.authTokens[uid] = token
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
}

// Create a data helper to query table datas.
func UTest() *utestHelper {
	return &utestHelper{*WingHelper}
}

const (
	_sql_ut_last_id    = "SELECT id FROM %s %s ORDER BY %s DESC LIMIT 1"
	_sql_ut_last_ids   = "SELECT id FROM %s WHERE %s >= ?"
	_sql_ut_last_field = "SELECT %s FROM %s ORDER BY %s DESC LIMIT 1"
	_sql_ut_get_target = "SELECT %s FROM %s WHERE %s = ? ORDER BY id DESC LIMIT 1"
	_sql_ut_get_tgt_id = "SELECT id FROM %s WHERE %s = ? ORDER BY id DESC LIMIT 1"
	_sql_ut_get_datas  = "SELECT %s FROM %s WHERE %s IN (%s)"
	_sql_ut_del_one    = "DELETE FROM %s WHERE %s = ?"
	_sql_ut_del_multis = "DELETE FROM %s WHERE %s IN (%s)"
	_sql_ut_clear      = "DELETE FROM %s %s"
)

// Get last id from given table, or with time confitions.
//
//	SQL: SELECT id FROM %s ORDER BY %s DESC LIMIT 1
//	                    ^           ^
//	                    table       order
//
//	SQL: SELECT id FROM %s WHERE %s >= "%s" ORDER BY %s DESC LIMIT 1
//	                    ^        ^      ^            ^
//	                 table wheres[0,    1]           order
//
//	@param table  Target table name.
//	@param order  Field name to order by.
//	@param wheres Where conditions contain field name and value.
func (t *utestHelper) LastID(table, order string, wheres ...string) (id int64, e error) {
	where := ""
	if len(wheres) >= 2 {
		where = fmt.Sprintf("WHERE %s >= \"%s\"", wheres[0], wheres[1])
	}
	query := fmt.Sprintf(_sql_ut_last_id, table, where, order)

	return id, t.One(query, func(rows *sql.Rows) error {
		if e = rows.Scan(&id); e != nil {
			return e
		}
		return nil
	})
}

// Get last ids from given table and query time, or compare condition value.
//
//	SQL: SELECT id FROM %s WHERE %s >= ?
//	                    ^        ^     ^
//	                    table    where value
//
//	@param table Target table name.
//	@param where Field name as where condition like: field > value.
//	@param value Where condition value to filter.
func (t *utestHelper) LastIDs(table, where string, value any) ([]int64, error) {
	ids, query := []int64{}, fmt.Sprintf(_sql_ut_last_ids, table, where)
	return ids, t.Query(query, func(rows *sql.Rows) error {
		var id int64
		if err := rows.Scan(&(id)); err != nil {
			return err
		}

		ids = append(ids, id)
		return nil
	}, value)
}

// Query the target field value by the top most given order field.
//
//	SQL: SELECT %s FROM %s ORDER BY %s DESC LIMIT 1
//	            ^       ^           ^
//	            target  table       order
//
//	@param table  Target table name.
//	@param target Target field name to output query result.
//	@param order  Field name to order by.
func (t *utestHelper) LastField(table string, target string, order string) (v string, e error) {
	query := fmt.Sprintf(_sql_ut_last_field, target, table, order)
	return v, t.One(query, func(rows *sql.Rows) error {
		if e = rows.Scan(&v); e != nil {
			return e
		}
		return nil
	})
}

// Query the target field last value by given condition field and value.
//
//	SQL: SELECT %s FROM %s WHERE %s = ? ORDER BY id DESC LIMIT 1
//	            ^       ^        ^    ^
//	            target  table  where  value
//
//	@param table  Target table name.
//	@param target Target field name to output query result.
//	@param where  Field name as where condition like: field = value.
//	@param value  Where condition value to query.
func (t *utestHelper) Target(table, target, where, value string) (v string, e error) {
	query := fmt.Sprintf(_sql_ut_get_target, target, table, where)
	return v, t.One(query, func(rows *sql.Rows) error {
		if e = rows.Scan(&v); e != nil {
			return e
		}
		return nil
	}, value)
}

// Query the target field last id by given condition field and value.
//
//	SQL: SELECT id FROM %s WHERE %s = ? ORDER BY id DESC LIMIT 1
//	                    ^        ^    ^
//	                    table  where  value
//
//	@param table  Target table name.
//	@param where  Field name as where condition like: field = value.
//	@param value  Where condition value to query.
func (t *utestHelper) TagID(table, where, value string) (id int64, e error) {
	query := fmt.Sprintf(_sql_ut_get_tgt_id, table, where)
	return id, t.One(query, func(rows *sql.Rows) error {
		if e = rows.Scan(&id); e != nil {
			return e
		}
		return nil
	}, value)
}

// Query the target field values by given condition field and values.
//
//	SQL: SELECT %s FROM %s WHERE %s IN (%s)
//	            ^       ^        ^      ^
//	            target  table    where  values
//
//	@param table  Target table name.
//	@param target Target field name to output query results.
//	@param where  Field name as where condition like: field IN (values).
//	@param values Where condition values to query.
func (t *utestHelper) Datas(table string, target string, field string, values string) ([]string, error) {
	rsts, query := []string{}, fmt.Sprintf(_sql_ut_get_datas, target, table, field, values)
	return rsts, t.Query(query, func(rows *sql.Rows) error {
		rst := ""
		if err := rows.Scan(&rst); err != nil {
			return err
		}
		rsts = append(rsts, rst)
		return nil
	})
}

// Deleta records by target field on equal condition.
//
//	SQL: DELETE FROM %s WHERE %s = ?
//	                 ^        ^    ^
//	                 table  where  value
//
//	@param table Target table name.
//	@param where Field name as where condition like: field = value.
//	@param value Where condition value to query.
func (t *utestHelper) DelOne(table, where string, value any) {
	t.Execute(fmt.Sprintf(_sql_ut_del_one, table, where), value)
}

// Deleta records by target field on in range condition.
//
//	SQL: DELETE FROM %s WHERE %s IN (%s)
//	                 ^        ^      ^
//	                 table    where  values
//
//	@param table Target table name.
//	@param where Field name as where condition like: field IN (values).
//	@param value Where condition values to query.
func (t *utestHelper) DelMults(table, where, values string) {
	t.Execute(fmt.Sprintf(_sql_ut_del_multis, table, where, values))
}

// Clear the target table all datas, or ranged datas of in given conditions.
//
//	SQL: DELETE FROM %s %s
//	                 ^  ^
//	             table  wheres
//
//	@param table  Target table name.
//	@param wheres Where conditions append to sql command tails if exist.
func (t *utestHelper) Clear(table string, wheres ...string) {
	t.Execute(fmt.Sprintf(_sql_ut_clear, table, strings.Join(wheres, " ")))
}
