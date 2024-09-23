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
)

// SetRequest use for set http request before execute http.Client.Do,
// you can use this middle-ware to set auth as username and passord, and so on.
//	@param req Http requester
//	@return - bool If current request ignore TLS verify or not, false is verify by default.
//			- error Exception message
type SetRequest func(req *http.Request) (bool, error)

// httpUtils inner http utils struct
type httpUtils struct{ silent bool }

// getHttpUtils return http utils instanse with silent status
func getHttpUtils(silent bool) *httpUtils { return &httpUtils{silent: silent} }

// readResponse read response body after executed request, it should return
// invar.ErrInvalidState when response code is not http.StatusOK.
func (u *httpUtils) readResponse(resp *http.Response) ([]byte, error) {
	if resp.StatusCode != http.StatusOK {
		logger.E("Failed http client, status:", resp.StatusCode)
		return nil, invar.ErrInvalidState
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.E("Failed read response, err:", err)
		return nil, err
	} else if !u.silent {
		logger.D("Response:", string(body))
	}
	return body, nil
}

// unmarshalResponse unmarshal response body after execute request,
// it may not check the given body if empty.
func (u *httpUtils) unmarshalResponse(body []byte, out any) error {
	if err := json.Unmarshal(body, out); err != nil {
		logger.E("Unmarshal body to struct err:", err)
		return err
	} else if !u.silent {
		logger.D("Response struct:", out)
	}
	return nil
}

// httpPostJson http post method, you can set post data as json struct.
func (u *httpUtils) httpPostJson(tagurl string, postdata any) ([]byte, error) {
	params, err := json.Marshal(postdata)
	if err != nil {
		logger.E("Marshal post data err:", err)
		return nil, err
	}

	resp, err := http.Post(tagurl, ContentTypeJson, bytes.NewReader(params))
	if err != nil {
		logger.E("Http post json err:", err)
		return nil, err
	}

	defer resp.Body.Close()
	return u.readResponse(resp)
}

// httpPostForm http post method, you can set post data as url.Values.
func (u *httpUtils) httpPostForm(tagurl string, postdata url.Values) ([]byte, error) {
	resp, err := http.PostForm(tagurl, postdata)
	if err != nil {
		logger.E("Http post form err:", err)
		return nil, err
	}

	defer resp.Body.Close()
	return u.readResponse(resp)
}

// httpClientDo handle http client DO method, and return response.
func (u *httpUtils) httpClientDo(req *http.Request, setRequestFunc SetRequest) ([]byte, error) {
	client := &http.Client{}

	// use middle-ware to set request header
	if setRequestFunc != nil {
		ignoreTLS, err := setRequestFunc(req)
		if err != nil {
			logger.E("Set http request err:", err)
			return nil, err
		}

		logger.I("httpClientDo: ignore TLS:", ignoreTLS)
		client.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: ignoreTLS,
			},
		}
	}

	// execute http request
	resp, err := client.Do(req)
	if err != nil {
		logger.E("Execute client DO, err:", err)
		return nil, err
	}

	defer resp.Body.Close()
	return u.readResponse(resp)
}

func (u *httpUtils) Get(tagurl string, params ...any) ([]byte, error) {
	if len(params) > 0 {
		tagurl = fmt.Sprintf(tagurl, params...)
	}

	rawurl := EncodeUrl(tagurl)
	if !u.silent {
		logger.D("Http Get:", rawurl)
	}

	resp, err := http.Get(rawurl)
	if err != nil {
		logger.E("Failed http get, err:", err)
		return nil, err
	}
	defer resp.Body.Close()
	return u.readResponse(resp)
}

func (u *httpUtils) Post(tagurl string, postdata any, contentType ...string) ([]byte, error) {
	ct := ContentTypeJson
	if len(contentType) > 0 {
		ct = contentType[0]
	} else if !u.silent {
		logger.D("Http Post:", tagurl, "ContentType:", ct)
	}

	switch ct {
	case ContentTypeJson:
		return u.httpPostJson(tagurl, postdata)
	case ContentTypeForm:
		return u.httpPostForm(tagurl, postdata.(url.Values))
	}
	return nil, invar.ErrInvalidParams
}

func (u *httpUtils) GetString(tagurl string, params ...any) (string, error) {
	resp, err := u.Get(tagurl, params...)
	if err != nil {
		return "", err
	}
	return strings.Trim(string(resp), "\""), nil
}

func (u *httpUtils) PostString(tagurl string, postdata any, contentType ...string) (string, error) {
	resp, err := u.Post(tagurl, postdata, contentType...)
	if err != nil {
		return "", err
	}
	return strings.Trim(string(resp), "\""), nil
}

func (u *httpUtils) GetStruct(tagurl string, out any, params ...any) error {
	body, err := u.Get(tagurl, params...)
	if err != nil {
		return err
	}
	return u.unmarshalResponse(body, out)
}

func (u *httpUtils) PostStruct(tagurl string, postdata, out any, contentType ...string) error {
	body, err := u.Post(tagurl, postdata, contentType...)
	if err != nil {
		return err
	}
	return u.unmarshalResponse(body, out)
}

func (u *httpUtils) ClientGet(tagurl string, setRequestFunc SetRequest, params ...any) ([]byte, error) {
	if len(params) > 0 {
		tagurl = fmt.Sprintf(tagurl, params...)
	}

	rawurl := EncodeUrl(tagurl)
	if !u.silent {
		logger.D("Http Client Get:", rawurl)
	}

	// generate new request instanse
	req, err := http.NewRequest(http.MethodGet, rawurl, http.NoBody)
	if err != nil {
		logger.E("Create http request err:", err)
		return nil, err
	}

	return u.httpClientDo(req, setRequestFunc)
}

func (u *httpUtils) ClientPost(tagurl string, setRequestFunc SetRequest, postdata ...any) ([]byte, error) {
	var body io.Reader
	if len(postdata) > 0 {
		params, err := json.Marshal(postdata[0])
		if err != nil {
			logger.E("Marshal post data err:", err)
			return nil, err
		}
		body = bytes.NewReader(params)
	} else {
		body = http.NoBody
	}

	if !u.silent {
		logger.D("Http Client Post:", tagurl)
	}

	// generate new request instanse
	req, err := http.NewRequest(http.MethodPost, tagurl, body)
	if err != nil {
		logger.E("Create http request err:", err)
		return nil, err
	}

	// set json as default content type
	req.Header.Set("Content-Type", ContentTypeJson)
	return u.httpClientDo(req, setRequestFunc)
}

func (u *httpUtils) ClientGetStruct(tagurl string, setRequestFunc SetRequest, out any, params ...any) error {
	body, err := u.ClientGet(tagurl, setRequestFunc, params...)
	if err != nil {
		return err
	}
	return u.unmarshalResponse(body, out)
}

func (u *httpUtils) ClientPostStruct(tagurl string, setRequestFunc SetRequest, out any, postdata ...any) error {
	body, err := u.ClientPost(tagurl, setRequestFunc, postdata...)
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

// HttpGet handle http get method
func HttpGet(tagurl string, params ...any) ([]byte, error) {
	return getHttpUtils(false).Get(tagurl, params...)
}

// HttpPost handle http post method, you can set content type as
// comm.ContentTypeJson or comm.ContentTypeForm, or other you need set.
//
// ---
//
//	// set post data as json string
//	data := struct {"key": "Value", "id": "123"}
//	resp, err := comm.HttpPost(tagurl, data)
//
//	// set post data as form string
//	data := "key=Value&id=123"
//	resp, err := comm.HttpPost(tagurl, data, comm.ContentTypeForm)
func HttpPost(tagurl string, postdata any, contentType ...string) ([]byte, error) {
	return getHttpUtils(false).Post(tagurl, postdata, contentType...)
}

// HttpGetString call HttpGet and trim " char both begin and end
func HttpGetString(tagurl string, params ...any) (string, error) {
	return getHttpUtils(false).GetString(tagurl, params...)
}

// HttpPostString call HttpPost and trim " char both begin and end.
func HttpPostString(tagurl string, postdata any, contentType ...string) (string, error) {
	return getHttpUtils(false).PostString(tagurl, postdata, contentType...)
}

// HttpGetStruct handle http get method and unmarshal data to struct object
func HttpGetStruct(tagurl string, out any, params ...any) error {
	return getHttpUtils(false).GetStruct(tagurl, out, params...)
}

// HttpPostStruct handle http post method and unmarshal data to struct object
func HttpPostStruct(tagurl string, postdata, out any, contentType ...string) error {
	return getHttpUtils(false).PostStruct(tagurl, postdata, out, contentType...)
}

// HttpClientGet handle http get by http.Client, you can set request headers or
// ignore TLS verfiy of https url by setRequstFunc middle-ware function as :
//
// ---
//
//	resp, err := comm.HttpClientGet(tagurl, func(req *http.Request) (bool, error) {
//		req.Header.Set("Content-Type", "application/json;charset=UTF-8")
//		req.SetBasicAuth("username", "password") // set auther header
//		return true, nil  // true is ignore TLS verify of https url
//	}, "same-params");
func HttpClientGet(tagurl string, setRequestFunc SetRequest, params ...any) ([]byte, error) {
	return getHttpUtils(false).ClientGet(tagurl, setRequestFunc, params...)
}

// HttpClientPost handle https post by http.Client, you can set request headers or
// ignore TLS verfiy of https url by setRequstFunc middle-ware function as :
//
// ---
//
//	resp, err := comm.HttpClientPost(tagurl, func(req *http.Request) (bool, error) {
//		req.Header.Set("Content-Type", "application/json;charset=UTF-8")
//		req.SetBasicAuth("username", "password") // set auther header
//		return true, nil  // true is ignore TLS verify of https url
//	}, "post-data")
func HttpClientPost(tagurl string, setRequestFunc SetRequest, postdata ...any) ([]byte, error) {
	return getHttpUtils(false).ClientPost(tagurl, setRequestFunc, postdata...)
}

// HttpClientGetStruct handle http get method and unmarshal data to struct object
func HttpClientGetStruct(tagurl string, setRequestFunc SetRequest, out any, params ...any) error {
	return getHttpUtils(false).ClientGetStruct(tagurl, setRequestFunc, out, params...)
}

// HttpClientPostStruct handle http post method and unmarshal data to struct object
func HttpClientPostStruct(tagurl string, setRequestFunc SetRequest, out any, postdata ...any) error {
	return getHttpUtils(false).ClientPostStruct(tagurl, setRequestFunc, out, postdata...)
}

// ----------------------------------------

// SilentGet call httpUtils.Get() on silent state
func SilentGet(tagurl string, params ...any) ([]byte, error) {
	return getHttpUtils(true).Get(tagurl, params...)
}

// SilentPost call httpUtils.Post() on silent state
func SilentPost(tagurl string, postdata any, contentType ...string) ([]byte, error) {
	return getHttpUtils(true).Post(tagurl, postdata, contentType...)
}

// SilentGetString call httpUtils.GetString() on silent state
func SilentGetString(tagurl string, params ...any) (string, error) {
	return getHttpUtils(true).GetString(tagurl, params...)
}

// SilentPostString call httpUtils.PostString() on silent state
func SilentPostString(tagurl string, postdata any, contentType ...string) (string, error) {
	return getHttpUtils(true).PostString(tagurl, postdata, contentType...)
}

// SilentGetStruct call httpUtils.GetStruct() on silent state
func SilentGetStruct(tagurl string, out any, params ...any) error {
	return getHttpUtils(true).GetStruct(tagurl, out, params...)
}

// SilentPostStruct call httpUtils.PostStruct() on silent state
func SilentPostStruct(tagurl string, postdata, out any, contentType ...string) error {
	return getHttpUtils(true).PostStruct(tagurl, postdata, out, contentType...)
}

// SilentClientGet call httpUtils.ClientGet() on silent state
func SilentClientGet(tagurl string, setRequestFunc SetRequest, params ...any) ([]byte, error) {
	return getHttpUtils(true).ClientGet(tagurl, setRequestFunc, params...)
}

// SilentClientPost call httpUtils.ClientPost() on silent state
func SilentClientPost(tagurl string, setRequestFunc SetRequest, postdata ...any) ([]byte, error) {
	return getHttpUtils(true).ClientPost(tagurl, setRequestFunc, postdata...)
}

// SilentClientGetStruct call httpUtils.ClientGetStruct() on silent state
func SilentClientGetStruct(tagurl string, setRequestFunc SetRequest, out any, params ...any) error {
	return getHttpUtils(true).ClientGetStruct(tagurl, setRequestFunc, out, params...)
}

// SilentClientPostStruct call httpUtils.ClientPostStruct() on silent state
func SilentClientPostStruct(tagurl string, setRequestFunc SetRequest, out any, postdata ...any) error {
	return getHttpUtils(true).ClientPostStruct(tagurl, setRequestFunc, out, postdata...)
}
