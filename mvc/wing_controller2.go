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
	"strings"

	"github.com/wengoldx/xcore/logger"
)

// WAuthController the extend controller base on WingController to support
// auth account from http headers, the client caller must append two headers
// before post request if expect the controller method enable execute token
// authentication from header.
//
// * Author : It must fixed keyword as WENGOLD-V1.2
//
// * Token : Authenticate JWT token responsed by login success
//
// `Optional headers` :
//
// * Location : Optional value of client indicator, global location
//
// * Authoration : The old version keyword for WENGOLD-V1.1
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
//		mvc.WAuthController
//	}
//
//	func init() {
//		mvc.GAuthHandlerFunc = func(token string) (string, string) {
//			// decode and verify token string, than return indecated
//			// account uuid and password.
//			return "account uuid", "account password"
//		}
//	}
//
// `USECASE 1. Auth account and Parse input params`
//
//	//	@Description Restful api bind with /login on POST method
//	//	@Param Author header string true "WENGOLD-V1.2"
//	//	@Param Token  header string true "Authentication token"
//	//	@Param data body types.Accout true "input param description"
//	//	@Success 200 {string} "response data description"
//	//	@router /login [post]
//	func (c *AccController) AccLogin() {
//		ps := &types.Accout{}
//		c.DoAfterValidated(ps, func(uuid string) (int, any) {
//		// Or get authed account password as :
//		// c.DoAfterAuthValidated(ps, func(uuid, pwd string) (int, any) {
//			// do same business with input NO-EMPTY account uuid,
//			// directe use c and ps param in this methed.
//			// ...
//			return http.StatusOK, "Done business"
//		} , false /* not limit error message even code is 40x */)
//	}
//
// `USECASE 2. Auth account on GET http method`
//
//	//	@Description Restful api bind with /detail on GET method
//	//	@Param Author header string true "WENGOLD-V1.2"
//	//	@Param Token  header string true "Authentication token"
//	//	@Success 200 {types.Detail} "response data description"
//	//	@router /detail [get]
//	func (c *AccController) AccDetail() {
//		if uuid := c.AuthRequestHeader(); uuid != "" {
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
type WAuthController struct {
	WingController
}

// NextFunc2 do action after input params validated, it decode token to get account uuid.
type NextFunc2 func(uuid string) (int, any)

// NextFunc3 do action after input params validated, it decode token to get account uuid and password.
type NextFunc3 func(uuid, pwd string) (int, any)

// AuthHandlerFunc auth request token from http header and returen account secures.
type AuthHandlerFunc func(token string) (string, string)

// RoleHandlerFunc verify role access permission from account service.
type RoleHandlerFunc func(sub, obj, act string) bool

// Global handler function to auth token from http header
var GAuthHandlerFunc AuthHandlerFunc

// Global handler function to verify role from http header
var GRoleHandlerFunc RoleHandlerFunc

// Get authoration and token from http header, than verify it and return account secures.
//	@Return 401, 403 codes returned on error.
func (c *WAuthController) AuthRequestHeader(hidelog ...bool) string {
	uuid, _ := c.innerAuthHeader(len(hidelog) > 0 && hidelog[1])
	return uuid
}

// DoAfterValidated do bussiness action after success validate the given json data.
//	@Return 400, 401, 403, 404 codes returned on error.
func (c *WAuthController) DoAfterValidated(ps any, nextFunc2 NextFunc2, fs ...bool) {
	protect, hidelog := !(len(fs) > 0 && !fs[0]), (len(fs) > 1 && fs[1])
	if uuid, _ := c.innerAuthHeader(hidelog); uuid != "" {
		c.doAfterValidatedInner("json", ps, nextFunc2, uuid, true, protect, hidelog)
	}
}

// DoAfterUnmarshal do bussiness action after success unmarshaled the given json data.
//	@Return 400, 401, 403, 404 codes returned on error.
func (c *WAuthController) DoAfterUnmarshal(ps any, nextFunc2 NextFunc2, fs ...bool) {
	protect, hidelog := !(len(fs) > 0 && !fs[0]), (len(fs) > 1 && fs[1])
	if uuid, _ := c.innerAuthHeader(hidelog); uuid != "" {
		c.doAfterValidatedInner("json", ps, nextFunc2, uuid, false, protect, hidelog)
	}
}

// DoAfterValidatedXml do bussiness action after success validate the given xml data.
//	@Return 400, 401, 403, 404 codes returned on error.
func (c *WAuthController) DoAfterValidatedXml(ps any, nextFunc2 NextFunc2, fs ...bool) {
	protect, hidelog := !(len(fs) > 0 && !fs[0]), (len(fs) > 1 && fs[1])
	if uuid, _ := c.innerAuthHeader(hidelog); uuid != "" {
		c.doAfterValidatedInner("xml", ps, nextFunc2, uuid, true, protect, hidelog)
	}
}

// DoAfterUnmarshalXml do bussiness action after success unmarshaled the given xml data.
//	@Return 400, 401, 403, 404 codes returned on error.
func (c *WAuthController) DoAfterUnmarshalXml(ps any, nextFunc2 NextFunc2, fs ...bool) {
	protect, hidelog := !(len(fs) > 0 && !fs[0]), (len(fs) > 1 && fs[1])
	if uuid, _ := c.innerAuthHeader(hidelog); uuid != "" {
		c.doAfterValidatedInner("xml", ps, nextFunc2, uuid, false, protect, hidelog)
	}
}

// ----------------------------------------

// DoAfterAuthValidated do bussiness action after success validate the given json data.
//	@Return 400, 401, 403, 404 codes returned on error.
func (c *WAuthController) DoAfterAuthValidated(ps any, nextFunc3 NextFunc3, fs ...bool) {
	protect, hidelog := !(len(fs) > 0 && !fs[0]), (len(fs) > 1 && fs[1])
	if uuid, pwd := c.innerAuthHeader(hidelog); uuid != "" {
		c.doAfterValidatedInner3("json", ps, nextFunc3, uuid, pwd, true, protect, hidelog)
	}
}

// DoAfterAuthUnmarshal do bussiness action after success unmarshaled the given json data.
//	@Return 400, 401, 403, 404 codes returned on error.
func (c *WAuthController) DoAfterAuthUnmarshal(ps any, nextFunc3 NextFunc3, fs ...bool) {
	protect, hidelog := !(len(fs) > 0 && !fs[0]), (len(fs) > 1 && fs[1])
	if uuid, pwd := c.innerAuthHeader(hidelog); uuid != "" {
		c.doAfterValidatedInner3("json", ps, nextFunc3, uuid, pwd, false, protect, hidelog)
	}
}

// DoAfterAuthValidatedXml do bussiness action after success validate the given xml data.
//	@Return 400, 401, 403, 404 codes returned on error.
func (c *WAuthController) DoAfterAuthValidatedXml(ps any, nextFunc3 NextFunc3, fs ...bool) {
	protect, hidelog := !(len(fs) > 0 && !fs[0]), (len(fs) > 1 && fs[1])
	if uuid, pwd := c.innerAuthHeader(hidelog); uuid != "" {
		c.doAfterValidatedInner3("xml", ps, nextFunc3, uuid, pwd, true, protect, hidelog)
	}
}

// DoAfterAuthUnmarshalXml do bussiness action after success unmarshaled the given xml data.
//	@Return 400, 401, 403, 404 codes returned on error.
func (c *WAuthController) DoAfterAuthUnmarshalXml(ps any, nextFunc3 NextFunc3, fs ...bool) {
	protect, hidelog := !(len(fs) > 0 && !fs[0]), (len(fs) > 1 && fs[1])
	if uuid, pwd := c.innerAuthHeader(hidelog); uuid != "" {
		c.doAfterValidatedInner3("xml", ps, nextFunc3, uuid, pwd, false, protect, hidelog)
	}
}

// ----------------------------------------

// Get authoration and token from http header, than verify it and return account secures.
//	@return 401: Unsupport author header or auth token failed.
//	@return 403: Denied permission of user access the rest4 API.
func (c *WAuthController) innerAuthHeader(hidelog bool) (string, string) {
	if GAuthHandlerFunc == nil || GRoleHandlerFunc == nil {
		c.E401Unauthed("Controller not set global handlers!")
		return "", ""
	}

	// check authoration secure key on right version
	header := c.Ctx.Request.Header
	if author := strings.ToUpper(header.Get("Author")); author == "" {
		if author = strings.ToUpper(header.Get("Authoration")); author != "WENGOLD-V1.1" {
			c.E401Unauthed("Unsupport v1 author: " + author)
			return "", ""
		}
	} else if author != "WENGOLD-V1.2" {
		c.E401Unauthed("Unsupport v2 author: " + author)
		return "", ""
	}

	// get token from header and verify it and user role
	if token := header.Get("Token"); token != "" {
		if uuid, pwd := GAuthHandlerFunc(token); uuid == "" {
			c.E401Unauthed("Unauthed header token!")
			return "", ""
		} else {
			if !GRoleHandlerFunc(uuid, c.Ctx.Input.URL(), c.Ctx.Request.Method) {
				c.E403Denind("Role permission denied for " + uuid)
				return "", ""
			}

			if !hidelog {
				logger.D("Authenticated account:", uuid)
			}
			return uuid, pwd
		}
	}

	// token is empty or invalid, response unauthed
	c.E401Unauthed("Unauthed header token!")
	return "", ""
}

// doAfterValidatedInner do bussiness action after success unmarshal params or
// validate the unmarshaled json data.
//	@See validatrParams() for more 400, 404 error code returned.
func (c *WAuthController) doAfterValidatedInner(datatype string,
	ps any, nextFunc2 NextFunc2, uuid string, validate, protect, hidelog bool) {
	if !c.validatrParams(datatype, ps, validate) {
		return
	}

	// execute business function after unmarshal and validated
	if status, resp := nextFunc2(uuid); resp != nil {
		c.responCheckState(datatype, protect, hidelog, status, resp)
	} else {
		c.responCheckState(datatype, protect, hidelog, status)
	}
}

// doAfterValidatedInner3 do bussiness action after success unmarshal params or
// validate the unmarshaled json data.
//	@See validatrParams() for more 400, 404 error code returned.
func (c *WAuthController) doAfterValidatedInner3(datatype string,
	ps any, nextFunc3 NextFunc3, uuid, pwd string, validate, protect, hidelog bool) {
	if !c.validatrParams(datatype, ps, validate) {
		return
	}

	// execute business function after unmarshal and validated
	if status, resp := nextFunc3(uuid, pwd); resp != nil {
		c.responCheckState(datatype, protect, hidelog, status, resp)
	} else {
		c.responCheckState(datatype, protect, hidelog, status)
	}
}
