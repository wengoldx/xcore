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
)

// WRoleController the extend controller base on WingController to support
// auth account from http headers, the client caller must append auth hander
// before post request if expect the controller method enable execute token
// authentication from header.
//
// * Author : It must fixed keyword as WENGOLD-V2.0
//
// * Token : Authenticate JWT token responsed by login success
//
// `USAGE` :
//
// The validator register code of input params struct see WingController description,
// but the restful auth api of router method as follow usecase 1 and 2.
//
// ---
//
// `controller.go`
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
// `USECASE 1. Auth account and Parse input params`
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
// `USECASE 2. Auth account on GET http method`
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
// `USECASE 3. No-Auth and Use WingController`
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
// `USECASE 4. No-Auth and Custom code`
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
	UID  int64  // Account unique id, maybe 0.
	Acc  string // Account unique name or string unique id, maybe empty.
	Pwd  string // Account password plaintext, maybe empty.
	Role string // Account role, maybe empty.
}

// NextHander do action after input params validated, it decode token to get account secures.
type NextHander func(a *WAuths) (int, any)

// ValidateHandlerFunc auth request token from http header and returen account secures.
type ValidateHandlerFunc func(token, router, method string) *WAuths

// Global handler function to verify token and role from http header
var ValidateHandler ValidateHandlerFunc

// Get authoration and token from http header, than verify it and return account secures.
//
//	This function just suport version of 'WENGOLD-V2.0' header without any input params.
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
			if len(silent) > 0 && silent[0] {
				logger.D("Authed account:", s.UID, s.Acc)
			}
			return s // account secures
		}
	}

	// token is empty or invalid, response unauthed
	c.E401Unauthed("Unauthed header token!")
	return nil
}

// Parse and validate input params, then do api action after success validated.
//
//	This function only suport 'WENGOLD-V2.0' header with input params.
//
//	@return 401: Invalid author header or permission denied.
//	@Return 400: Failed parse input params or validate error.
//	@Return 404: Case exception in server.
func (c *WRoleController) DoAfterValidated(ps any, nextHander NextHander, opts ...bool) {
	silent := len(opts) > 0 && opts[0]
	if s := c.AuthRequestHeader(silent); s != nil {
		c.doAfterValidatedInner(ps, nextHander, s, true, silent)
	}
}

// Parse input params, then do api action after success unmarshaled.
//
//	This function only suport 'WENGOLD-V2.0' header with input params.
//
//	@return 401: Invalid author header or permission denied.
//	@Return 400: Failed parse input params.
//	@Return 404: Case exception in server.
func (c *WRoleController) DoAfterUnmarshal(ps any, nextHander NextHander, opts ...bool) {
	silent := len(opts) > 0 && opts[0]
	if s := c.AuthRequestHeader(silent); s != nil {
		c.doAfterValidatedInner(ps, nextHander, s, false, silent)
	}
}

// Parse input param, validate if need, then call api hander method and response result.
func (c *WRoleController) doAfterValidatedInner(ps any, nextHander NextHander, s *WAuths, validate, silent bool) {
	datatype := "json"
	if !c.validatrParams(datatype, ps, validate) {
		return
	}

	// execute business function after unmarshal and validated
	if status, resp := nextHander(s); resp != nil {
		c.responCheckState(datatype, true, silent, status, resp)
	} else {
		c.responCheckState(datatype, true, silent, status)
	}
}
