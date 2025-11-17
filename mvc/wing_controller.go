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
	"mime/multipart"
	"reflect"
	"strings"

	"github.com/astaxie/beego"
	"github.com/go-playground/validator/v10"
	"github.com/wengoldx/xcore/invar"
	"github.com/wengoldx/xcore/logger"
)

// WingController the base controller to support bee http functions.
//
// # USAGE:
//
// Notice that you should register the field level validator for the input data's struct,
// then use it in struct describetion label as validate target.
//
//	`types.go`
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
//	`controller.go`
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

// FileFunc do action after get multipart file param.
type FileFunc func(file multipart.File, header *multipart.FileHeader)

// FilesFunc do action after get multipart files param.
type FilesFunc func(headers []*multipart.FileHeader)

/* ------------------------------------------------------------------- */
/* For Validator                                                       */
/* ------------------------------------------------------------------- */

// Validator use for verify the input params on struct level
var Validator *validator.Validate

// Generat the validator instance if need
func ensureValidatorGenerated() {
	if Validator == nil {
		Validator = validator.New()
	}
}

// Register struct field validators from given map
func RegisterValidators(valmap map[string]validator.Func) {
	for tag, valfunc := range valmap {
		RegisterFieldValidator(tag, valfunc)
	}
}

// Register validators on struct field level
func RegisterFieldValidator(tag string, valfunc validator.Func) {
	ensureValidatorGenerated()
	if err := Validator.RegisterValidation(tag, valfunc); err != nil {
		logger.E("Register validator:"+tag+", err:", err)
		return
	}
}

/* ------------------------------------------------------------------- */
/* For Response Status & Datas                                         */
/* ------------------------------------------------------------------- */

// Sends a json response to client on status check mode.
func (c *WingController) ResponJSON(state int, data ...any) {
	c.responCheckState(newOptions(true, false), state, data...)
}

// UncheckJSON sends a json response to client witchout status check.
func (c *WingController) UncheckJSON(state int, dataORerr ...any) {
	c.responCheckState(newOptions(false, false), state, dataORerr...)
}

// SilentJSON sends a json response to client without output ok log.
func (c *WingController) SilentJSON(state int, data ...any) {
	c.responCheckState(newOptions(true, true), state, data...)
}

// Sends a ['json', 'jsonp', 'xml', 'yarm'] response to client on status check mode.
func (c *WingController) ResponAsType(datatype string, state int, data ...any) {
	c.responCheckState(newOptions(true, false).outType(datatype), state, data...)
}

// Sends a ['json', 'jsonp', 'xml', 'yarm'] response to client witchout status check.
func (c *WingController) UncheckAsType(datatype string, state int, dataORerr ...any) {
	c.responCheckState(newOptions(false, false).outType(datatype), state, dataORerr...)
}

// SilentXML sends xml response to client without output ok log.
func (c *WingController) SilentAsType(datatype string, state int, data ...any) {
	c.responCheckState(newOptions(true, true).outType(datatype), state, data...)
}

// SilentData sends JSON, JSONNP, XML, YAML, response data to client without output ok log.
func (c *WingController) SilentData(state int, data ...map[any]any) {
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
func (c *WingController) ResponOK(slient ...bool) {
	if !(len(slient) > 0 && slient[0]) {
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

/* ------------------------------------------------------------------- */
/* For Response Errors                                                 */
/* ------------------------------------------------------------------- */

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

/* ------------------------------------------------------------------- */
/* For Export Utils Methods                                            */
/* ------------------------------------------------------------------- */

// Get number value by target key from url params with default value, it
// will just return default value when target key value not a [int], [int64],
// [float64] number type.
//
//	See c.GetString() to get a string value.
func GetParam[T any](c *beego.Controller, key string, def T) T {
	switch v := (any(def)).(type) {
	case int:
		if val, err := c.GetInt(key, v); err == nil {
			return (any(val)).(T)
		}
	case int64:
		if val, err := c.GetInt64(key, v); err == nil {
			return (any(val)).(T)
		}
	case float64:
		if val, err := c.GetFloat(key, v); err == nil {
			return (any(val)).(T)
		}
	default:
		logger.E("Unspport url params key:", key)
	}
	return def
}

// BindValue bind value with key from url, the dest container must pointer
func (c *WingController) BindValue(key string, dest any) error {
	if err := c.Ctx.Input.Bind(dest, key); err != nil {
		logger.E("Parse", key, "from url, err:", err)
		return invar.ErrInvalidData
	}
	return nil
}

// OutHeader set none-empty response header as key:value to frontend.
func (c *WingController) OutHeader(key, value string) {
	if key != "" && value != "" {
		c.Ctx.Output.Header(key, value)
	}
}

// OutRole set role as response header to frontend.
func (c *WingController) OutRole(role string, status int) {
	if status == invar.StatusOK && role != "" {
		c.Ctx.Output.Header("role", role)
	}
}

// ClientFrom return client ip from who requested
func (c *WingController) ClientFrom() string {
	return c.Ctx.Request.RemoteAddr
}

/* ------------------------------------------------------------------- */
/* For Swagger Rest4 API Utils                                         */
/* ------------------------------------------------------------------- */

// Do next action on 'dev' runmode.
func (c *WingController) RunDevMode(next func()) {
	if beego.BConfig.RunMode != "dev" {
		c.E403Denind("Only For Debug!")
		return
	}
	next()
}

// Do next action after got the opened multipart file, and close it when exit method.
func (c *WingController) GetSingleFile(key string, next FileFunc) {
	file, header, err := c.GetFile(key)
	if err != nil {
		logger.E("Get file by:", key, "err:", err)
		c.E400Params()
		return
	}
	defer file.Close()
	next(file, header)
}

// Do next action after got the multipart files header.
func (c *WingController) GetMultiFiles(key string, next FilesFunc) {
	headers, err := c.GetFiles(key)
	if err != nil {
		logger.E("Get files by:", key, "err:", err)
		c.E400Params()
		return
	}
	next(headers)
}

// Read multipart file content and return buffer datas.
func (c *WingController) ReadFile(header multipart.FileHeader) ([]byte, error) {
	buf := make([]byte, header.Size)
	if tf, err := header.Open(); err != nil {
		return nil, err
	} else {
		defer tf.Close()
		if _, err = tf.Read(buf); err != nil {
			return nil, err
		}
		return buf, nil
	}
}

// Do bussiness action after parsed url params and success validate.
//
//	@Return 400 code returned on error.
//
// # NOTICE:
//	- This function not check 'WENGOLD-V*' header.
//	- Use AuthController, RoleController to check header and token.
//
// # WARNING:
//	- The out param must create as a struct pointer for this methoed!
//	- This method only support types: bool, string, int, int32, int64, uint, uint32, uint64, float32, float64.
func (c *WingController) DoAfterParsed(ps any, nextFunc NextFunc, opts ...Option) {
	c.doAfterParsedInner(ps, nextFunc, parseOptions(true /* no-use */, opts...))
}

// Do bussiness action after success validate the given json or xml data.
//
//	@Return 400, 404 codes returned on error.
//
// # NOTICE:
//	- This function not check 'WENGOLD-V*' header.
//	- Use AuthController, RoleController to check header and token.
func (c *WingController) DoAfterValidated(ps any, nextFunc NextFunc, opts ...Option) {
	c.doAfterValidatedInner(ps, nextFunc, parseOptions(true, opts...))
}

// Do bussiness action after success unmarshaled the given json or xml data.
//
//	@Return 400, 404 codes returned on error.
//
// # NOTICE:
//	- This function not check 'WENGOLD-V*' header.
//	- Use AuthController, RoleController to check header and token.
func (c *WingController) DoAfterUnmarshal(ps any, nextFunc NextFunc, opts ...Option) {
	c.doAfterValidatedInner(ps, nextFunc, parseOptions(false, opts...))
}

/* ------------------------------------------------------------------- */
/* For Internal Utils Methods                                          */
/* ------------------------------------------------------------------- */

// Check respon state and print out log, the datatype must range in
// ['json', 'jsonp', 'xml', 'yaml'], if out of range current controller
// just return blank string to close http connection.
//
// The protect param set true by default, by can be change from input flags.
func (c *WingController) responCheckState(opts *Options, state int, data ...any) {
	dt := strings.ToUpper(opts.datatype)
	if state != invar.StatusOK {
		/* ------------------------------------------------------------
		 * Not response error message to frontend when protect is true!
		 * ------------------------------------------------------------ */
		if state != invar.StatusExError && opts.Protect {
			c.ErrorState(state)
			return
		}

		/*
		 * Unprotect mode, response error code and message to frontend,
		 * Here dispathed 4xx http request errors and 202 extend error!
		 */
		errmsg := invar.StatusText(state)
		ctl, act := c.GetControllerAndAction()
		logger.E("["+dt+"] Respone ERR:", state, ">", ctl+"."+act, errmsg)
	}

	// Output simple ok response usually, but can hide by input flag.
	if !opts.Silent {
		ctl, act := c.GetControllerAndAction()
		logger.I("["+dt+"] Respone OK >", ctl+"."+act)
	}

	c.Ctx.Output.Status = state
	if len(data) > 0 && data[0] != nil {
		c.Data[opts.datatype] = data[0]
	}

	switch opts.datatype {
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
		logger.W("Unsupport response type:" + dt)
		c.Ctx.ResponseWriter.Write([]byte(""))
	}
}

// Do bussiness action after success validate the given json or xml data.
//
//	@See validateUrlParams() for more 400 error code returned.
func (c *WingController) doAfterParsedInner(ps any, nextFunc NextFunc, opts *Options) {
	if !c.validateUrlParams(ps, true) { // fixed validate true!
		return
	}

	// execute business function after validated.
	status, resp := nextFunc()
	c.responCheckState(opts, status, resp)
}

// Do bussiness action after success unmarshal params or validate the unmarshaled json or xml data.
//
//	@See validateParams() for more 400, 404 error code returned.
func (c *WingController) doAfterValidatedInner(ps any, nextFunc NextFunc, opts *Options) {
	if !c.validateParams(ps, opts) {
		return
	}

	// execute business function after unmarshal and validated
	status, resp := nextFunc()
	c.responCheckState(opts, status, resp)
}

// Do bussiness action after success unmarshal params or validate the unmarshaled json or xml data.
//
//	@Return 400: Invalid input params(Unmarshal error or invalid params value).
//	@Return 404: Internale server error(not support content type unless json and xml).
func (c *WingController) validateParams(ps any, opts *Options) bool {
	switch opts.datatype {
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
		c.E404Exception("Invalid data type:" + opts.datatype)
		return false
	}

	// validate input params if need
	if opts.validate {
		ensureValidatorGenerated()
		if err := Validator.Struct(ps); err != nil {
			c.E400Validate(ps, err.Error())
			return false
		}
	}
	return true
}

// Do bussiness action after success parsed and validate the params.
//
//	@Return 400: Invalid url params(Parse error or invalid params value).
func (c *WingController) validateUrlParams(ps any, validate bool) bool {
	if !c.parseUrlParams(ps) {
		c.E400Params("Failed parse url params!")
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

// Parse and save the input params from http request url for GET method.
//
// # NOTICE:
//
// This method only support simple value types of bool, int, float, string
// for input struct field, and filter out the others value types.
//
//	param := &MyStruct{
//		Name string `json:"name"` // get 'name' value from url and set to Name.
//		Aga  int                  // none json tag, filter out.
//	}
//	parseUrlParams(param)
//
// # WARNING:
//	- The sample param must create as a struct pointer for this methoed!
//	- This method only support types: bool, string, int, int32, int64, uint, uint32, uint64, float32, float64.
func (c *WingController) parseUrlParams(ps any) bool {
	rv := reflect.ValueOf(ps)
	if !rv.IsValid() || rv.Kind() != reflect.Ptr || rv.IsNil() {
		return false
	}

	rv = rv.Elem()  // get param struct value
	rt := rv.Type() // get param struct types
	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		name, tag := field.Name, field.Tag.Get("json")
		if name == "" || tag == "" {
			continue // filter none json tag fields.
		}

		v := rv.FieldByName(name)
		if v.IsValid() && v.CanSet() {
			switch field.Type.Kind() {
			case reflect.Bool:
				pv, _ := c.GetBool(tag, false)
				v.SetBool(pv)
			case reflect.String:
				v.SetString(c.GetString(tag, ""))
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				pv, _ := c.GetInt64(tag, 0)
				v.SetInt(pv)
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				pv, _ := c.GetUint64(tag, 0)
				v.SetUint(pv)
			case reflect.Float32, reflect.Float64:
				pv, _ := c.GetFloat(tag, 0)
				v.SetFloat(pv)
			}
		}
	}
	return true
}
