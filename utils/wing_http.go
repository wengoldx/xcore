// Copyright (c) 2018-Now Dunyu All Rights Reserved.
//
// Author      : https://www.wengold.net
// Email       : support@wengold.net
//
// Prismy.No | Date       | Modified by. | Description
// -------------------------------------------------------------------
// 00001       2019/05/22   yangping       New version
// -------------------------------------------------------------------

package utils

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"

	"github.com/wengoldx/xcore/invar"
	"github.com/wengoldx/xcore/logger"
)

const (
	// ContentTypeJson json content type
	ContentTypeJson = "application/json;charset=UTF-8"

	// ContentTypeForm form content type
	ContentTypeForm = "application/x-www-form-urlencoded"

	// ContentTypeFile file upload data type
	ContentTypeFile = "multipart/form-data"
)

// A middleware method called before execute http.Client.Do to set http request
// headers which such as username and passord of authentications or others.
//
//	@param req Http requester
//	@return - bool  Retrun true for using TLS, or false by default not verify.
//			- error Exception message
//
//	@See more [content-types](https://www.runoob.com/http/http-content-type.html).
type SetRequest func(req *http.Request) (bool, error)

// A utils struct for http accesss.
type httpUtils struct{ silent bool }

// The global singleton of httpUtils for easy handle http GET or POST request.
var (
	HttpUtils  = &httpUtils{silent: false}
	HttpSlient = &httpUtils{silent: true}
)

// Read response body after executed request, it should return invar.ErrInvalidState
// when response code is not http.StatusOK (200).
func (u *httpUtils) readResponse(resp *http.Response, parse bool) ([]byte, error) {
	if resp.StatusCode != http.StatusOK {
		logger.E("Http request failed, code:", resp.StatusCode)
		return nil, invar.ErrInvalidState
	}

	// parse response data if require.
	if parse {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			logger.E("Read response, err:", err)
			return nil, err
		} else if !u.silent {
			logger.D("Response:", string(body))
		}
		return body, nil
	}

	// return success without parse response
	return nil, nil
}

// Unmarshal response body after execute request, it not check the body whether empty.
func (u *httpUtils) unmarshalResponse(body []byte, out any) error {
	if err := json.Unmarshal(body, out); err != nil {
		logger.E("Unmarshal response, err:", err)
		return err
	}
	u.log("Response struct:", out)
	return nil
}

// Post http request with json params.
func (u *httpUtils) postJson(tagurl string, datas any, parse bool) ([]byte, error) {
	params, err := json.Marshal(datas)
	if err != nil {
		logger.E("Marshal post datas, err:", err)
		return nil, err
	}

	resp, err := http.Post(tagurl, ContentTypeJson, bytes.NewReader(params))
	if err != nil {
		logger.E("Http post, err:", err)
		return nil, err
	}

	defer resp.Body.Close()
	return u.readResponse(resp, parse)
}

// Post http request with form valus as url.Values.
func (u *httpUtils) postForm(tagurl string, datas url.Values, parse bool) ([]byte, error) {
	resp, err := http.PostForm(tagurl, datas)
	if err != nil {
		logger.E("Http post, err:", err)
		return nil, err
	}

	defer resp.Body.Close()
	return u.readResponse(resp, parse)
}

// Excute http.Client.Do with request header set callback, and return response results.
func (u *httpUtils) clientDo(req *http.Request, setRequestFunc SetRequest) ([]byte, error) {
	client := &http.Client{}

	// use middle-ware to set request header
	if setRequestFunc != nil {
		ignoreTLS, err := setRequestFunc(req)
		if err != nil {
			logger.E("Set http header, err:", err)
			return nil, err
		}

		logger.I("httpUtils: ignore TLS >", ignoreTLS)
		client.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: ignoreTLS,
			},
		}
	}

	// execute http request
	resp, err := client.Do(req)
	if err != nil {
		logger.E("Execute client.DO, err:", err)
		return nil, err
	}

	defer resp.Body.Close()
	return u.readResponse(resp, true)
}

// Handle http GET method request and parse response data if required.
func (u *httpUtils) handleGet(tagurl string, parse bool, params ...any) ([]byte, error) {
	if len(params) > 0 {
		tagurl = fmt.Sprintf(tagurl, params...)
	}

	rawurl := EncodeUrl(tagurl)
	u.log("Http Get:", rawurl)

	resp, err := http.Get(rawurl)
	if err != nil {
		logger.E("Failed http get, err:", err)
		return nil, err
	}
	defer resp.Body.Close()
	return u.readResponse(resp, parse)
}

// Handle http POST method request and parse response data if required.
func (u *httpUtils) handlePost(tagurl string, datas any, parse bool, contentType ...string) ([]byte, error) {
	ct := VarString(contentType, ContentTypeJson)
	u.log("Http Post:", tagurl, "ContentType:", ct)

	switch ct {
	case ContentTypeJson:
		return u.postJson(tagurl, datas, parse)
	case ContentTypeForm:
		return u.postForm(tagurl, datas.(url.Values), parse)
	}
	return nil, invar.ErrInvalidParams
}

// Output debug logs if require not slient.
func (u *httpUtils) log(msgs ...any) {
	if !u.silent {
		logger.D(msgs...)
	}
}

// Handle http GET method and return original response bytes.
//
//	USAGE:
//
//	params := "key=Value&id=123"
//	resp, err := utils.HttpUtils.Get(tagurl, params)
func (u *httpUtils) Get(tagurl string, params ...any) ([]byte, error) {
	return u.handleGet(tagurl, true, params...)
}

// Handle http GET method without parse any response datas.
func (u *httpUtils) GEmit(tagurl string, params ...any) (e error) {
	_, e = u.handleGet(tagurl, false, params...)
	return
}

// Handle http GET method and return response datas as string.
func (u *httpUtils) GString(tagurl string, params ...any) (string, error) {
	resp, err := u.Get(tagurl, params...)
	if err != nil {
		return "", err
	}
	return strings.Trim(string(resp), "\""), nil
}

// Handle http GET method and parse response datas to given struct object.
func (u *httpUtils) GStruct(tagurl string, out any, params ...any) error {
	body, err := u.Get(tagurl, params...)
	if err != nil {
		return err
	}
	return u.unmarshalResponse(body, out)
}

// Handle http POST method and return original response bytes,
// the content-type header can be set as utils.ContentTypeJson, utils.ContentTypeForm,
// utils.ContentTypeFile or others which you want.
//
//	USAGE:
//
//	// set post datas as json string.
//	datas := struct {"key": "Value", "id": "123"}
//	resp, err := utils.HttpUtils.Post(tagurl, data)
//
//	// set post datas as form string.
//	datas := "key=Value&id=123"
//	resp, err := utils.HttpUtils.Post(tagurl, datas, comm.ContentTypeForm)
func (u *httpUtils) Post(tagurl string, datas any, contentType ...string) ([]byte, error) {
	return u.handlePost(tagurl, datas, true, contentType...)
}

// Handle http POST method without parse any response datas.
func (u *httpUtils) PEmit(tagurl string, datas any, contentType ...string) (e error) {
	_, e = u.handlePost(tagurl, datas, false, contentType...)
	return
}

// Handle http POST method and return response datas as string.
func (u *httpUtils) PString(tagurl string, datas any, contentType ...string) (string, error) {
	resp, err := u.Post(tagurl, datas, contentType...)
	if err != nil {
		return "", err
	}
	return strings.Trim(string(resp), "\""), nil
}

// Handle http POST method and parse response datas to given struct object.
func (u *httpUtils) PStruct(tagurl string, datas, out any, contentType ...string) error {
	body, err := u.Post(tagurl, datas, contentType...)
	if err != nil {
		return err
	}
	return u.unmarshalResponse(body, out)
}

// Handle http GET method by http.Client and return original response bytes,
// use the setRequstFunc middleware callback to set request headers, or ignore
// TLS verfiy of https auth.
//
//	USAGE:
//
//	resp, err := utils.HttpUtils.CGet(tagurl, func(req *http.Request) (bool, error) {
//		req.Header.Set("Content-Type", "application/json;charset=UTF-8")
//		req.SetBasicAuth("username", "password") // set auther header
//		return true, nil                         // true is ignore TLS verify of https url.
//	}, "get-form-params");
func (u *httpUtils) CGet(tagurl string, setRequestFunc SetRequest, params ...any) ([]byte, error) {
	if len(params) > 0 {
		tagurl = fmt.Sprintf(tagurl, params...)
	}

	rawurl := EncodeUrl(tagurl)
	u.log("Http Client Get:", rawurl)

	// generate new request instanse
	req, err := http.NewRequest(http.MethodGet, rawurl, http.NoBody)
	if err != nil {
		logger.E("Create http request err:", err)
		return nil, err
	}
	return u.clientDo(req, setRequestFunc)
}

// Handle http GET method by http.Client and parse response datas to given struct object.
func (u *httpUtils) CGStruct(tagurl string, setRequestFunc SetRequest, out any, params ...any) error {
	body, err := u.CGet(tagurl, setRequestFunc, params...)
	if err != nil {
		return err
	}
	return u.unmarshalResponse(body, out)
}

// Handle http POST method by http.Client and return original response bytes,
// use the setRequstFunc middleware callback to set request headers, or ignore
// TLS verfiy of https auth.
//
//	USAGE:
//
//	resp, err := utils.HttpUtils.CPost(tagurl, func(req *http.Request) (bool, error) {
//		req.Header.Set("Content-Type", "application/json;charset=UTF-8")
//		req.SetBasicAuth("username", "password") // set auther header
//		return true, nil                         // true is ignore TLS verify of https auth.
//	}, "post-data")
func (u *httpUtils) CPost(tagurl string, setRequestFunc SetRequest, datas ...any) ([]byte, error) {
	var body io.Reader
	if len(datas) > 0 {
		params, err := json.Marshal(datas[0])
		if err != nil {
			logger.E("Marshal post data err:", err)
			return nil, err
		}
		body = bytes.NewReader(params)
	} else {
		body = http.NoBody
	}

	u.log("Http Client Post:", tagurl)
	// generate new request instanse
	req, err := http.NewRequest(http.MethodPost, tagurl, body)
	if err != nil {
		logger.E("Create http request err:", err)
		return nil, err
	}

	// set json as default content type
	req.Header.Set("Content-Type", ContentTypeJson)
	return u.clientDo(req, setRequestFunc)
}

// Handle http POST method by http.Client and parse response datas to given struct object.
func (u *httpUtils) CPStruct(tagurl string, setRequestFunc SetRequest, out any, datas ...any) error {
	body, err := u.CPost(tagurl, setRequestFunc, datas...)
	if err != nil {
		return err
	}
	return u.unmarshalResponse(body, out)
}

// ----------------------------------------

// GetIP get just ip not port from controller.Ctx.Request.RemoteAddr of beego
func GetIP(remoteaddr string) string {
	ip, _, _ := net.SplitHostPort(remoteaddr)
	if ip == "::1" {
		ip = "127.0.0.1"
	}
	logger.I("Got ip [", ip, "] from [", remoteaddr, "]")
	return ip
}

// GetLocalIPs get all the loacl IP of current deploy machine
func GetLocalIPs() ([]string, error) {
	netfaces, err := net.Interfaces()
	if err != nil {
		logger.E("Get ip interfaces err:", err)
		return nil, err
	}

	ips := []string{}
	for _, netface := range netfaces {
		addrs, err := netface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			switch v := addr.(type) {
			case *net.IPNet:
				if v.IP.IsGlobalUnicast() {
					ips = append(ips, v.IP.String())
				}
			}
		}
	}

	// Check the result list is empty
	if len(ips) == 0 {
		return nil, invar.ErrNotFound
	}

	return ips, nil
}

// EncodeUrl encode url params
func EncodeUrl(rawurl string) string {
	enurl, err := url.Parse(rawurl)
	if err != nil {
		logger.E("Encode urlm err:", err)
		return rawurl
	}
	enurl.RawQuery = enurl.Query().Encode()
	return enurl.String()
}

// ----------------------------------------

// Deprecated: Handle http get method.
func HttpGet(tagurl string, params ...any) ([]byte, error) {
	return HttpUtils.Get(tagurl, params...)
}

// Deprecated: Handle http post method.
func HttpPost(tagurl string, datas any, contentType ...string) ([]byte, error) {
	return HttpUtils.Post(tagurl, datas, contentType...)
}

// Deprecated: Call HttpGet and trim " char both begin and end.
func HttpGetString(tagurl string, params ...any) (string, error) {
	return HttpUtils.GString(tagurl, params...)
}

// Deprecated: Call HttpPost and trim " char both begin and end.
func HttpPostString(tagurl string, datas any, contentType ...string) (string, error) {
	return HttpUtils.PString(tagurl, datas, contentType...)
}

// Deprecated: Handle http get method and unmarshal data to struct object.
func HttpGetStruct(tagurl string, out any, params ...any) error {
	return HttpUtils.GStruct(tagurl, out, params...)
}

// Deprecated: Handle http post method and unmarshal data to struct object.
func HttpPostStruct(tagurl string, datas, out any, contentType ...string) error {
	return HttpUtils.PStruct(tagurl, datas, out, contentType...)
}

// Deprecated: Handle http get by http.Client.
func HttpClientGet(tagurl string, setRequestFunc SetRequest, params ...any) ([]byte, error) {
	return HttpUtils.CGet(tagurl, setRequestFunc, params...)
}

// Deprecated: Handle https post by http.Client.
func HttpClientPost(tagurl string, setRequestFunc SetRequest, datas ...any) ([]byte, error) {
	return HttpUtils.CPost(tagurl, setRequestFunc, datas...)
}

// Deprecated: Handle http get method and unmarshal data to struct object.
func HttpClientGetStruct(tagurl string, setRequestFunc SetRequest, out any, params ...any) error {
	return HttpUtils.CGStruct(tagurl, setRequestFunc, out, params...)
}

// Deprecated: Handle http post method and unmarshal data to struct object.
func HttpClientPostStruct(tagurl string, setRequestFunc SetRequest, out any, datas ...any) error {
	return HttpUtils.CPStruct(tagurl, setRequestFunc, out, datas...)
}

// Deprecated: Call httpUtils.Get() on silent state.
func SilentGet(tagurl string, params ...any) ([]byte, error) {
	return HttpSlient.Get(tagurl, params...)
}

// Deprecated: all httpUtils.Post() on silent state.
func SilentPost(tagurl string, datas any, contentType ...string) ([]byte, error) {
	return HttpSlient.Post(tagurl, datas, contentType...)
}

// Deprecated: Call httpUtils.GString() on silent state.
func SilentGetString(tagurl string, params ...any) (string, error) {
	return HttpSlient.GString(tagurl, params...)
}

// Deprecated: Call httpUtils.PString() on silent state.
func SilentPostString(tagurl string, datas any, contentType ...string) (string, error) {
	return HttpSlient.PString(tagurl, datas, contentType...)
}

// Deprecated: Call httpUtils.GStruct() on silent state.
func SilentGetStruct(tagurl string, out any, params ...any) error {
	return HttpSlient.GStruct(tagurl, out, params...)
}

// Deprecated: Call httpUtils.PStruct() on silent state.
func SilentPostStruct(tagurl string, datas, out any, contentType ...string) error {
	return HttpSlient.PStruct(tagurl, datas, out, contentType...)
}

// Deprecated: Call httpUtils.CGet() on silent state.
func SilentClientGet(tagurl string, setRequestFunc SetRequest, params ...any) ([]byte, error) {
	return HttpSlient.CGet(tagurl, setRequestFunc, params...)
}

// Deprecated: Call httpUtils.CPost() on silent state.
func SilentClientPost(tagurl string, setRequestFunc SetRequest, datas ...any) ([]byte, error) {
	return HttpSlient.CPost(tagurl, setRequestFunc, datas...)
}

// Deprecated: Call httpUtils.CGStruct() on silent state.
func SilentClientGetStruct(tagurl string, setRequestFunc SetRequest, out any, params ...any) error {
	return HttpSlient.CGStruct(tagurl, setRequestFunc, out, params...)
}

// Deprecated: Call httpUtils.CPStruct() on silent state.
func SilentClientPostStruct(tagurl string, setRequestFunc SetRequest, out any, datas ...any) error {
	return HttpSlient.CPStruct(tagurl, setRequestFunc, out, datas...)
}
