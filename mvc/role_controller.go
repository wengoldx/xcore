// Copyright (c) 2018-Now Dunyu All Rights Reserved.
//
// Author      : https://www.wengold.net
// Email       : support@wengold.net
//
// Prismy.No | Date       | Modified by. | Description
// -------------------------------------------------------------------
// 00001       2025/05/11   yangping       New version
// -------------------------------------------------------------------

package mvc

import (
	"strings"

	"github.com/wengoldx/xcore/logger"
	"github.com/wengoldx/xcore/utils"
)

// WRoleController the extend controller base on WingController to support
// auth account from http headers, the client caller must append auth hander
// before post request if expect the controller method enable execute token
// authentication from header.
//
//	- Author : It must fixed keyword as WENGOLD-V2.0
//	- Token  : Authenticate JWT token responsed by login success.
//
// # USAGE:
//
// The validator register code of input params struct see WingController description,
// but the restful auth api of router method as follow usecase 1 and 2.
//
//	`controller.go`
//
//	// define custom controller using header auth function
//	type AccController struct {
//		mvc.WRoleController
//	}
//
//	func init() {
//		mvc.ValidateHandler = func(token, router, method string) *mvc.WAuths {
//			secures := &mvc.WAuths{}
//			// decode and verify token string, than return indecated
//			// account id, string uuid, password plaintext and role.
//			return secures
//		}
//	}
//
//	`USECASE 1. Auth account and Parse input params`
//
//	//	@Description Restful api bind with /login on POST method
//	//	@Param Author header string true "WENGOLD-V2.0"
//	//	@Param Token  header string true "Authentication token"
//	//	@Param data   body   types.Accout true "input param description"
//	//	@Success 200 {string} "response data description"
//	//	@router /login [post]
//	func (c *AccController) AccLogin() {
//		ps := &types.Accout{}
//		c.DoAfterValidated(ps, func(s *WAuths) (int, any) {
//			// do api action with input NOT-NULL account secures,
//			// directe use c and ps param in this methed.
//			// ...
//			return http.StatusOK, "Done business"
//		} , false /* silent flag, true is not output success logs */)
//	}
//
//	`USECASE 2. Auth account on GET http method`
//
//	//	@Description Restful api bind with /detail on GET method
//	//	@Param Author header string true "WENGOLD-V2.0"
//	//	@Param Token  header string true "Authentication token"
//	//	@Success 200 {types.Detail} "response data description"
//	//	@router /detail [get]
//	func (c *AccController) AccDetail() {
//		if s := c.AuthRequestHeader(); s != nil {
//			// use c.BindValue("fieldkey", out) parse params from url
//			c.ResponJSON(service.AccDetail(uuid))
//		}
//	}
//
//	`USECASE 3. No-Auth and Use WingController`
//
//	//	@Description Restful api bind with /update on POST method
//	//	@Param data body types.UserInfo true "input param description"
//	//	@Success 200
//	//	@router /update [post]
//	func (c *AccController) AccUpdate() {
//		ps := &types.UserInfo{}
//		c.WingController.DoAfterValidated(ps, func() (int, any) {
//			// directe use c and ps param in this methed.
//			// ...
//			return http.StatusOK, nil
//		} , false /* not limit error message even code is 40x */)
//	}
//
//	`USECASE 4. No-Auth and Custom code`
//
//	//	@Description Restful api bind with /list on GET method
//	//	@Success 200 {object} []types.Account "response data description"
//	//	@router /list [get]
//	func (c *AccController) AccList() {
//		// do same business without auth and input params
//		c.ResponJSON(service.AccList())
//	}
type WRoleController struct {
	WingController
}

// Account secure datas, set after pass validation.
type WAuths struct {
	ID   int64  // Account id of int64,  set when using number id, default -1.
	UID  string // Account id of string, set when using string id, defalut empty.
	Pwd  string // Account password plaintext, maybe empty.
	Role string // Account role, maybe empty.
}

// Do action after input params validated, it decode token to get account secures.
type NextHander func(a *WAuths) (int, any)

// Auth request token from http header and returen account secures.
type ValidateHandlerFunc func(token, router, method string) *WAuths

// Global handler function to verify token and role from http header
var ValidateHandler ValidateHandlerFunc

// Get authoration and token from http header, than verify it and return account secures.
//
// # WARNING:
//
// This function only suport 'WENGOLD-V2.0' header for GET http method
// without any input params.
//
//	@return 401: Invalid author header or permission denied.
func (c *WRoleController) AuthRequestHeader(silent ...bool) *WAuths {
	if ValidateHandler == nil {
		c.E401Unauthed("Controller not set auth handler!")
		return nil
	}

	// check authoration secure key on right version
	header := c.Ctx.Request.Header
	if author := strings.ToUpper(header.Get("Author")); author != "WENGOLD-V2.0" {
		c.E401Unauthed("Unsupport v2 author: " + author)
		return nil
	}

	// get token from header and verify it and user role
	if token := header.Get("Token"); token != "" {
		router, method := c.Ctx.Input.URL(), c.Ctx.Request.Method
		if s := ValidateHandler(token, router, method); s == nil {
			c.E401Unauthed("Unauthed account!")
			return nil
		} else {
			if !utils.Variable(silent, false) {
				logger.Df("Authed account: %d:%s", s.ID, s.UID)
			}
			return s // account secures
		}
	}

	// token is empty or invalid, response unauthed
	c.E401Unauthed("Unauthed header token!")
	return nil
}

// Parse url params for GET method, then do api action after success parsed.
//
// # WARNING:
//
// This function only suport 'WENGOLD-V2.0' header for 'GET' http method,
// and parse simple input params from url.
//
//	@return 401: Invalid author header or permission denied.
//	@Return 400: Failed parse url params.
//	@Return 404: Case exception in server.
func (c *WRoleController) DoAfterParsed(ps any, nextHander NextHander, opt ...Option) {
	opts := parseOptions(true, opt...)
	if s := c.AuthRequestHeader(opts.Silent); s != nil {
		c.doAfterParsedInner(ps, nextHander, s, opts)
	}
}

// Parse url param, validate if need, then call api hander method and response result.
//
// # WARNING:
//
// This function not check 'WENGOLD-V2.0' header.
//
//	@Return 400: Failed parse url params.
//	@Return 404: Case exception in server.
func (c *WRoleController) DoParsedInsecure(ps any, nextFunc NextFunc, opt ...Option) {
	c.WingController.DoAfterParsed(ps, nextFunc, opt...)
}

// Parse and validate input params, then do api action after success validated.
//
// # WARNING:
//
// This function only suport 'WENGOLD-V2.0' header for POST http method.
//
//	@return 401: Invalid author header or permission denied.
//	@Return 400: Failed parse input params or validate error.
//	@Return 404: Case exception in server.
func (c *WRoleController) DoAfterValidated(ps any, nextHander NextHander, opt ...Option) {
	opts := parseOptions(true, opt...)
	if s := c.AuthRequestHeader(opts.Silent); s != nil {
		c.doAfterValidatedInner(ps, nextHander, s, opts)
	}
}

// Parse input params, then do api action after success unmarshaled.
//
// # WARNING:
//
// This function only suport 'WENGOLD-V2.0' header for POST http method.
//
//	@return 401: Invalid author header or permission denied.
//	@Return 400: Failed parse input params.
//	@Return 404: Case exception in server.
func (c *WRoleController) DoAfterUnmarshal(ps any, nextHander NextHander, opt ...Option) {
	opts := parseOptions(false, opt...)
	if s := c.AuthRequestHeader(opts.Silent); s != nil {
		c.doAfterValidatedInner(ps, nextHander, s, opts)
	}
}

// Parse input param, validate if need, then call api hander method and response result.
func (c *WRoleController) doAfterValidatedInner(ps any, nextHander NextHander, s *WAuths, opts *Options) {
	if !c.validateParams(ps, opts) {
		return
	}

	// execute business function after validated.
	status, resp := nextHander(s)
	c.responCheckState(opts, status, resp)
}

// Parse url param, validate if need, then call api hander method and response result.
func (c *WRoleController) doAfterParsedInner(ps any, nextHander NextHander, s *WAuths, opts *Options) {
	if !c.validateUrlParams(ps, true) { // fixed validate true!
		return
	}

	// execute business function after validated.
	status, resp := nextHander(s)
	c.responCheckState(opts, status, resp)
}
