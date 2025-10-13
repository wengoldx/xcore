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
	"net/url"

	pd "github.com/wengoldx/xcore/mvc/provider"
	"github.com/wengoldx/xcore/mvc/provider/mysql"
	"github.com/wengoldx/xcore/utils"
	"github.com/wengoldx/xcore/utils/httpx"
)

// Restful api tester runtime configs.
type tester struct {
	tabler   *pd.TableProvider // TableProvider instance.
	tokens   map[string]string // Auth token of test user, format as {uuid:token}.
	envs     map[string]any    // Env caches, user ut.Get(key) and ut.Set(key, value) to get and set datas.
	author   string            // Author header, such as 'WENGOLD-V1.1', 'WENGOLD-V1.2', 'WENGOLD-V2.0'
	tokenApi string            // Rest4 API to get user token, like 'http://192.168.1.100:8000/server/token?id=%s'
}

// Global tester signleton for test rest4 apis.
//
// Setup it as the follow code.
//
//	func init() {
//		// logger.SilentLoggers() // silent logger if comment out.
//		ut.SetupTester(ut.WithAuthor("WENGOLD-V2.0"),
//			ut.WithTokenApi("http://192.168.1.100:8000/server/token?id=%s"),
//		)
//	}
var _t = &tester{
	tokens: make(map[string]string),
	envs:   make(map[string]any),
}

// Setup tester instance with options and table provider.
func SetupTester(opts ...Option) {
	_t.tabler = mysql.GetTabler()
	for _, optFunc := range opts {
		optFunc(_t)
	}
}

// Return tester instance.
func Tester() *tester { return _t }

/* ------------------------------------------------------------------- */
/* SQL Querier, Inserter, Updater, Deleter                             */
/* ------------------------------------------------------------------- */

func Querier(t ...string) *pd.QueryBuilder   { return _t.tabler.Querier().Table(utils.Variable(t, "")) }
func Inserter(t ...string) *pd.InsertBuilder { return _t.tabler.Inserter().Table(utils.Variable(t, "")) }
func Updater(t ...string) *pd.UpdateBuilder  { return _t.tabler.Updater().Table(utils.Variable(t, "")) }
func Deleter(t ...string) *pd.DeleteBuilder  { return _t.tabler.Deleter().Table(utils.Variable(t, "")) }

/* ------------------------------------------------------------------- */
/* Env datas getter and setter                                         */
/* ------------------------------------------------------------------- */

func Get[T any](key string) T        { return _t.envs[key].(T) } // Get a env value, call it as iv := ut.Get[int]().
func Set[T any](key string, value T) { _t.envs[key] = value }    // Set a env value, call it as ut.Set[int](iv).

// Transform url params map to url.Values for http GET method.
func (t *tester) parseForms(params TestForm) string {
	forms := url.Values{}
	for param, value := range params {
		forms[param] = []string{fmt.Sprintf("%v", value)}
	}
	return forms.Encode()
}

// Get user access token from caches, or request by rest4 api.
func (t *tester) getToken(uid string) string {
	if token, ok := t.tokens[uid]; ok {
		return token
	} else if t.tokenApi != "" {
		var token string
		if err := httpx.Get(t.tokenApi, &token, uid); err == nil {
			t.tokens[uid] = token
			return token
		}
	}
	return ""
}
