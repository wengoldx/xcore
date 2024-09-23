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
	"encoding/json"
	"encoding/xml"
	"strings"

	"github.com/astaxie/beego"
	"github.com/go-playground/validator/v10"
	"github.com/wengoldx/xcore/invar"
	"github.com/wengoldx/xcore/logger"
)

// WingController the base controller to support bee http functions.
//
// `USAGE` :
//
// Notice that you should register the field level validator for the input data's struct,
// then use it in struct describetion label as validate target.
//
// ---
//
// `types.go`
//
//	// define restful api router input param struct
//	type struct Accout {
//		Acc string `json:"acc" validate:"required,IsVaildUuid"`
//		PWD string `json:"pwd" validate:"required_without"`
//		Num int    `json:"num"`
//	}
//
//	// define custom validator on struct field level
//	func isVaildUuid(fl validator.FieldLevel) bool {
//		m, _ := regexp.Compile("^[0-9a-zA-Z]*$")
//		str := fl.Field().String()
//		return m.MatchString(str)
//	}
//
//	func init() {
//		// register router input params struct validators
//		mvc.RegisterFieldValidator("IsVaildUuid", isVaildUuid)
//	}
//
// ---
//
// `controller.go`
//
//	//	@Description Restful api bind with /login on POST method
//	//	@Param data body types.Accout true "input param description"
//	//	@Success 200 {string} "response data description"
//	//	@router /login [post]
//	func (c *AccController) AccLogin() {
//		ps := &types.Accout{}
//		c.DoAfterValidated(ps, func() (int, any) {
//			// do same business function
//			// directe use c and ps param in this methed.
//			// ...
//			return http.StatusOK, "Done business"
//		} , false /* not limit error message even code is 40x */)
//	}
type WingController struct {
	beego.Controller
}

// NextFunc do action after input params validated.
type NextFunc func() (int, any)

// Validator use for verify the input params on struct level
var Validator *validator.Validate

// ensureValidatorGenerated generat the validator instance if need
func ensureValidatorGenerated() {
	if Validator == nil {
		Validator = validator.New()
	}
}

// RegisterValidators register struct field validators from given map
func RegisterValidators(valmap map[string]validator.Func) {
	for tag, valfunc := range valmap {
		RegisterFieldValidator(tag, valfunc)
	}
}

// RegisterFieldValidator register validators on struct field level
func RegisterFieldValidator(tag string, valfunc validator.Func) {
	ensureValidatorGenerated()
	if err := Validator.RegisterValidation(tag, valfunc); err != nil {
		logger.E("Register validator:"+tag+", err:", err)
		return
	}
	logger.I("Registered validator:", tag)
}

// ----------------------------------------

// ResponJSON sends a json response to client on status check mode.
func (c *WingController) ResponJSON(state int, data ...any) {
	c.responCheckState("json", true, false, state, data...)
}

// ResponJSONP sends a jsonp response to client on status check mode.
func (c *WingController) ResponJSONP(state int, data ...any) {
	c.responCheckState("jsonp", true, false, state, data...)
}

// ResponXML sends xml response to client on status check mode.
func (c *WingController) ResponXML(state int, data ...any) {
	c.responCheckState("xml", true, false, state, data...)
}

// ResponYAML sends yaml response to client on status check mode.
func (c *WingController) ResponYAML(state int, data ...any) {
	c.responCheckState("yaml", true, false, state, data...)
}

// UncheckJSON sends a json response to client witchout status check.
func (c *WingController) UncheckJSON(state int, dataORerr ...any) {
	c.responCheckState("json", false, false, state, dataORerr...)
}

// UncheckJSONP sends a jsonp response to client witchout status check.
func (c *WingController) UncheckJSONP(state int, dataORerr ...any) {
	c.responCheckState("jsonp", false, false, state, dataORerr...)
}

// UncheckXML sends xml response to client witchout status check.
func (c *WingController) UncheckXML(state int, dataORerr ...any) {
	c.responCheckState("xml", false, false, state, dataORerr...)
}

// UncheckYAML sends yaml response to client witchout status check.
func (c *WingController) UncheckYAML(state int, dataORerr ...any) {
	c.responCheckState("yaml", false, false, state, dataORerr...)
}

// SilentJSON sends a json response to client without output ok log.
func (c *WingController) SilentJSON(state int, data ...any) {
	c.responCheckState("json", true, true, state, data...)
}

// SilentJSONP sends a jsonp response to client without output ok log.
func (c *WingController) SilentJSONP(state int, data ...any) {
	c.responCheckState("jsonp", true, true, state, data...)
}

// SilentXML sends xml response to client without output ok log.
func (c *WingController) SilentXML(state int, data ...any) {
	c.responCheckState("xml", true, true, state, data...)
}

// SilentYAML sends yaml response to client without output ok log.
func (c *WingController) SilentYAML(state int, data ...any) {
	c.responCheckState("yaml", true, true, state, data...)
}

// SlientData sends JSON, JSONNP, XML, YAML, response data to client without output ok log.
func (c *WingController) SlientData(state int, data ...map[any]any) {
	if state != invar.StatusOK {
		c.ErrorState(state)
		return
	}

	c.Ctx.Output.Status = state
	if len(data) > 0 {
		c.Data = data[0]
	}
	c.ServeFormatted()
}

// ResponData sends JSON, JSONP, XML, YAML response datas to client, the data type depending
// on the value of the Accept header.
func (c *WingController) ResponData(state int, data ...map[any]any) {
	if state != invar.StatusOK {
		c.ErrorState(state)
		return
	}

	ctl, act := c.GetControllerAndAction()
	logger.I("Respone OK-DATA >", ctl+"."+act)

	c.Ctx.Output.Status = state
	if len(data) > 0 {
		c.Data = data[0]
	}
	c.ServeFormatted()
}

// ResponOK sends a empty success response to client
func (c *WingController) ResponOK(hidelog ...bool) {
	if !(len(hidelog) > 0 && hidelog[0]) {
		ctl, act := c.GetControllerAndAction()
		logger.I("Respone OK >", ctl+"."+act)
	}

	w := c.Ctx.ResponseWriter
	w.WriteHeader(invar.StatusOK)
	// FIXME: here maybe not set content type when response error
	// w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(""))
}

// ResponExErr sends a extend error as response data on 202 status code
func (c *WingController) ResponExErr(errmsg invar.WExErr) {
	ctl, act := c.GetControllerAndAction()
	logger.E("Respone ERR-EX >", ctl+"."+act, "err:", errmsg.Message)

	c.Ctx.Output.Status = invar.StatusExError
	c.Data["json"] = errmsg
	c.ServeJSON()
}

// ErrorState response error state to client
func (c *WingController) ErrorState(state int, err ...string) {
	ctl, act := c.GetControllerAndAction()
	errmsg := invar.StatusText(state)
	if len(err) > 0 {
		errmsg += ", " + err[0]
	}
	logger.E("Respone ERR:", state, ">", ctl+"."+act, errmsg)

	w := c.Ctx.ResponseWriter
	w.WriteHeader(state)
	// FIXME here maybe not set content type when response error
	// w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(""))
}

// E400Params response 400 invalid params error state to client
func (c *WingController) E400Params(err ...string) {
	c.ErrorState(invar.E400ParseParams, err...)
}

// E400rUnmarshal response 400 unmarshal params error state to client
func (c *WingController) E400Unmarshal(err ...string) {
	c.ErrorState(invar.E400ParseParams, err...)
}

// E400Validate response 400 invalid params error state to client, then print
// the params data and validate error
func (c *WingController) E400Validate(ps any, err ...string) {
	logger.E("Invalid input params:", ps)
	c.ErrorState(invar.E400ParseParams, err...)
}

// E401Unauthed response 401 unauthenticated error state to client
func (c *WingController) E401Unauthed(err ...string) {
	c.ErrorState(invar.E401Unauthorized, err...)
}

// E403Denind response 403 permission denind error state to client
func (c *WingController) E403Denind(err ...string) {
	c.ErrorState(invar.E403PermissionDenied, err...)
}

// E404Exception response 404 not found error state to client
func (c *WingController) E404Exception(err ...string) {
	c.ErrorState(invar.E404Exception, err...)
}

// E405Disabled response 405 function disabled error state to client
func (c *WingController) E405Disabled(err ...string) {
	c.ErrorState(invar.E405FuncDisabled, err...)
}

// E406Input response 406 invalid inputs error state to client
func (c *WingController) E406Input(err ...string) {
	c.ErrorState(invar.E406InputParams, err...)
}

// E409Duplicate response 409 duplicate error state to client
func (c *WingController) E409Duplicate(err ...string) {
	c.ErrorState(invar.E409Duplicate, err...)
}

// E410Gone response 410 gone error state to client
func (c *WingController) E410Gone(err ...string) {
	c.ErrorState(invar.E410Gone, err...)
}

// E426UpgradeRequired response 426 upgrade required error state to client
func (c *WingController) E426UpgradeRequired(err ...string) {
	c.ErrorState(invar.E426UpgradeRequired, err...)
}

// ClientFrom return client ip from who requested
func (c *WingController) ClientFrom() string {
	return c.Ctx.Request.RemoteAddr
}

// BindValue bind value with key from url, the dest container must pointer
func (c *WingController) BindValue(key string, dest any) error {
	if err := c.Ctx.Input.Bind(dest, key); err != nil {
		logger.E("Parse", key, "from url, err:", err)
		return invar.ErrInvalidData
	}
	return nil
}

// DoAfterValidated do bussiness action after success validate the given json data.
//	@Return 400, 404 codes returned on error.
func (c *WingController) DoAfterValidated(ps any, nextFunc NextFunc, fs ...bool) {
	protect, hidelog := !(len(fs) > 0 && !fs[0]), (len(fs) > 1 && fs[1])
	c.doAfterValidatedInner("json", ps, nextFunc, true, protect, hidelog)
}

// DoAfterUnmarshal do bussiness action after success unmarshaled the given json data.
//	@Return 400, 404 codes returned on error.
func (c *WingController) DoAfterUnmarshal(ps any, nextFunc NextFunc, fs ...bool) {
	protect, hidelog := !(len(fs) > 0 && !fs[0]), (len(fs) > 1 && fs[1])
	c.doAfterValidatedInner("json", ps, nextFunc, false, protect, hidelog)
}

// DoAfterValidatedXml do bussiness action after success validate the given xml data.
//	@Return 400, 404 codes returned on error.
func (c *WingController) DoAfterValidatedXml(ps any, nextFunc NextFunc, fs ...bool) {
	protect, hidelog := !(len(fs) > 0 && !fs[0]), (len(fs) > 1 && fs[1])
	c.doAfterValidatedInner("xml", ps, nextFunc, true, protect, hidelog)
}

// DoAfterUnmarshalXml do bussiness action after success unmarshaled the given xml data.
//	@Return 400, 404 codes returned on error.
func (c *WingController) DoAfterUnmarshalXml(ps any, nextFunc NextFunc, fs ...bool) {
	protect, hidelog := !(len(fs) > 0 && !fs[0]), (len(fs) > 1 && fs[1])
	c.doAfterValidatedInner("xml", ps, nextFunc, false, protect, hidelog)
}

// ----------------------------------------

// responCheckState check respon state and print out log, the datatype must
// range in ['json', 'jsonp', 'xml', 'yaml'], if out of range current controller
// just return blank string to close http connection.
// the protect param set true by default, by can be change from input flags.
func (c *WingController) responCheckState(datatype string, protect, hidelog bool, state int, data ...any) {
	dt := strings.ToUpper(datatype)
	if state != invar.StatusOK {
		/* ------------------------------------------------------------
		 * Not response error message to frontend when protect is true!
		 * ------------------------------------------------------------ */
		if state != invar.StatusExError && protect {
			c.ErrorState(state)
			return
		}

		/*
		 * Not Protect mode, response error code and message to frontend,
		 * Here dispathed 4xx http request errors and 202 extend error!
		 */
		errmsg := invar.StatusText(state)
		ctl, act := c.GetControllerAndAction()
		logger.E("["+dt+"] Respone ERR:", state, ">", ctl+"."+act, errmsg)
	}

	// Output simple ok response usually, but can hide by input flag.
	if !hidelog {
		ctl, act := c.GetControllerAndAction()
		logger.I("["+dt+"] Respone OK >", ctl+"."+act)
	}

	c.Ctx.Output.Status = state
	if len(data) > 0 {
		c.Data[datatype] = data[0]
	}

	switch datatype {
	case "json":
		c.ServeJSON()
	case "jsonp":
		c.ServeJSONP()
	case "xml":
		c.ServeXML()
	case "yaml":
		c.ServeYAML()
	default:
		// just return blank string to close http connection
		logger.W("Invalid response data tyep:" + datatype)
		c.Ctx.ResponseWriter.Write([]byte(""))
	}
}

// doAfterValidatedInner do bussiness action after success unmarshal params or
// validate the unmarshaled json data.
//	@See validatrParams() for more 400, 404 error code returned.
func (c *WingController) doAfterValidatedInner(datatype string, ps any, nextFunc NextFunc, validate, protect, hidelog bool) {
	if !c.validatrParams(datatype, ps, validate) {
		return
	}

	// execute business function after unmarshal and validated
	if status, resp := nextFunc(); resp != nil {
		c.responCheckState(datatype, protect, hidelog, status, resp)
	} else {
		c.responCheckState(datatype, protect, hidelog, status)
	}
}

// validatrParams do bussiness action after success unmarshal params or validate the unmarshaled json data.
//	@Return 400: Invalid input params(Unmarshal error or invalid params value).
//	@Return 404: Internale server error(not support content type unless json and xml).
func (c *WingController) validatrParams(datatype string, ps any, validate bool) bool {
	switch datatype {
	case "json":
		if err := json.Unmarshal(c.Ctx.Input.RequestBody, ps); err != nil {
			c.E400Unmarshal(err.Error())
			return false
		}
	case "xml":
		if err := xml.Unmarshal(c.Ctx.Input.RequestBody, ps); err != nil {
			c.E400Unmarshal(err.Error())
			return false
		}
	default: // current not support the jsonp and yaml parse
		c.E404Exception("Invalid data type:" + datatype)
		return false
	}

	// validate input params if need
	if validate {
		ensureValidatorGenerated()
		if err := Validator.Struct(ps); err != nil {
			c.E400Validate(ps, err.Error())
			return false
		}
	}
	return true
}
