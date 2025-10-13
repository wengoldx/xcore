// Copyright (c) 2018-Now Dunyu All Rights Reserved.
//
// Author      : https://www.wengold.net
// Email       : support@wengold.net
//
// Prismy.No | Date       | Modified by. | Description
// -------------------------------------------------------------------
// 00001       2019/05/22   yangping       New version
// -------------------------------------------------------------------

package httpx

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
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

// Read response body after executed request, it should return invar.ErrInvalidState
// when response code is not http.StatusOK (200).
func readResponse(resp *http.Response, parse bool) ([]byte, error) {
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
		} 
		// log("Response:", string(body))
		return body, nil
	}

	// return success without parse response
	return nil, nil
}

// Unmarshal response body after execute request, it not check the body whether empty.
func unmarshalResponse(body []byte, out any) error {
	if err := json.Unmarshal(body, out); err != nil {
		logger.E("Unmarshal response, err:", err)
		return err
	}
	// logger.D("Response struct:", out)
	return nil
}

// Post http request with json params.
func postJson(tagurl string, datas any, parse bool) ([]byte, error) {
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
	return readResponse(resp, parse)
}

// Post http request with form valus as url.Values.
func postForm(tagurl string, datas url.Values, parse bool) ([]byte, error) {
	resp, err := http.PostForm(tagurl, datas)
	if err != nil {
		logger.E("Http post, err:", err)
		return nil, err
	}

	defer resp.Body.Close()
	return readResponse(resp, parse)
}

// Handle http GET method request and parse response data if required.
//
// # WARNING:
//	- The tagurl must contain format marks such as '%s', '%d', '%v' when params not empty!
func handleGet(tagurl string, parse bool, params ...any) ([]byte, error) {
	if len(params) > 0 {
		tagurl = fmt.Sprintf(tagurl, params...)
	}

	rawurl := EncodeUrl(tagurl)
	logger.D("Http Get:", rawurl)

	resp, err := http.Get(rawurl)
	if err != nil {
		logger.E("Failed http get, err:", err)
		return nil, err
	}
	defer resp.Body.Close()
	return readResponse(resp, parse)
}

// Handle http POST method request and parse response data if required.
func handlePost(tagurl string, datas any, parse bool, contentType ...string) ([]byte, error) {
	ct :=ContentTypeJson
	if len(contentType) > 0 && contentType[0] != "" {
		ct = contentType[0]
	}
	logger.D("Http Post:", tagurl, "ContentType:", ct)

	switch ct {
	case ContentTypeJson:
		return postJson(tagurl, datas, parse)
	case ContentTypeForm:
		return postForm(tagurl, datas.(url.Values), parse)
	}
	return nil, invar.ErrInvalidParams
}

// Excute http.Client.Do with request header set callback, and return response results.
func clientDo[T any](req *http.Request, setRequestFunc SetRequest, out *T) error {
	client := &http.Client{}

	// use middle-ware to set request header
	if setRequestFunc != nil {
		ignoreTLS, err := setRequestFunc(req)
		if err != nil {
			logger.E("Set http header, err:", err)
			return err
		}

		logger.I("Ignore TLS >", ignoreTLS)
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
		return err
	}

	defer resp.Body.Close()
	datas, err :=  readResponse(resp, true)
	if err != nil {
		return err
	}
	return parseResp(datas, out)
}

func parseResp[T any](resp []byte, out *T) error {
	if out != nil {
		switch v := any(out).(type) {
		case *[]byte :
			*v = resp
			return nil
		case *string:
			*v = strings.Trim(string(resp), "\"")
			return nil
		case any:
			vt := reflect.TypeOf(*out)
			if vt.Kind() == reflect.Struct {
				return unmarshalResponse(resp, out)
			}
		}
	}
	return invar.ErrUnsupportFormat
}

// Handle http GET method without parse any response datas.
//
// # WARNING:
//	- The tagurl must contain format marks such as '%s', '%d', '%v' when params not empty!
func GEmit(tagurl string, params ...any) (e error) {
	_, e = handleGet(tagurl, false, params...)
	return
}

// Handle http POST method without parse any response datas.
func PEmit(tagurl string, datas any, contentType ...string) (e error) {
	_, e = handlePost(tagurl, datas, false, contentType...)
	return
}

// Handle http GET method and return response datas.
//
// # WARNING:
//	- The tagurl must contain format marks such as '%s', '%d', '%v' when params not empty!
//
// # USAGE:
//
//	var outbytes []byte    // get bytes response.
//	err := httputil.Get(tagurl, &outbytes)
//
//	tagurl_marks := "http://192.168.1.100/acc?key=%s&id=%d"
//	var outstring string   // get string response.
//	err := httputil.Get(tagurl_marks, &outstring, "key", 123)
//
//	var outstruct MyStruct // get srtuct response.
//	err := httputil.Get(tagurl, &outstruct)
func Get[T any](tagurl string, out *T, params ...any) error {
	resp, err := handleGet(tagurl, true, params...)
	if err != nil {
		return err
	} else if out == nil {
		return nil
	}
	return parseResp(resp, out)
}

// Handle http POST method and return original response bytes,
// the content-type header can be set as httputil.ContentTypeJson, 
// httputil.ContentTypeForm, httputil.ContentTypeFile or others which you want.
//
// # USAGE:
//
//	// set post datas as json string.
//	datas := struct {"key": "Value", "id": "123"}
//
//	var outbytes []byte    // get bytes response.
//	err := httputil.Post(tagurl, data, &outbytes)
//
//	var outstring string   // get string response.
//	err := httputil.Post(tagurl, data, &outstring, comm.ContentTypeForm)
//
//	var outstruct MyStruct // get srtuct response.
//	err := httputil.Post(tagurl, datas, &outstruct)
func Post[T any](tagurl string, datas any, out *T, contentType ...string) error {
	resp, err := handlePost(tagurl, datas, true, contentType...)
	if err != nil {
		return err
	} else if out == nil {
		return nil
	}
	return parseResp(resp, out)
}

// Handle http GET method by http.Client and return original response bytes,
// use the setRequstFunc middleware callback to set request headers, or ignore
// TLS verfiy of https auth.
//
// # USAGE:
//
//	resp, err := utils.HttpUtils.CGet(tagurl, func(req *http.Request) (bool, error) {
//		req.Header.Set("Content-Type", "application/json;charset=UTF-8")
//		req.SetBasicAuth("username", "password") // set auther header
//		return true, nil                         // true is ignore TLS verify of https url.
//	}, "get-form-params");
func ClientGet[T any](tagurl string, setRequestFunc SetRequest, out *T, params ...any) error {
	if len(params) > 0 {
		tagurl = fmt.Sprintf(tagurl, params...)
	}

	rawurl := EncodeUrl(tagurl)
	logger.D("Http Client Get:", rawurl)

	// generate new request instanse
	req, err := http.NewRequest(http.MethodGet, rawurl, http.NoBody)
	if err != nil {
		logger.E("Create http request err:", err)
		return err
	}
	return clientDo(req, setRequestFunc, out)
}

// Handle http POST method by http.Client and return original response bytes,
// use the setRequstFunc middleware callback to set request headers, or ignore
// TLS verfiy of https auth.
//
// # USAGE:
//
//	resp, err := utils.HttpUtils.CPost(tagurl, func(req *http.Request) (bool, error) {
//		req.Header.Set("Content-Type", "application/json;charset=UTF-8")
//		req.SetBasicAuth("username", "password") // set auther header
//		return true, nil                         // true is ignore TLS verify of https auth.
//	}, "post-data")
func ClientPost[T any](tagurl string, setRequestFunc SetRequest, out *T, datas ...any) error {
	var body io.Reader
	if len(datas) > 0 {
		params, err := json.Marshal(datas[0])
		if err != nil {
			logger.E("Marshal post data err:", err)
			return err
		}
		body = bytes.NewReader(params)
	} else {
		body = http.NoBody
	}

	logger.D("Http Client Post:", tagurl)
	// generate new request instanse
	req, err := http.NewRequest(http.MethodPost, tagurl, body)
	if err != nil {
		logger.E("Create http request err:", err)
		return err
	}

	// set json as default content type
	req.Header.Set("Content-Type", ContentTypeJson)
	return clientDo(req, setRequestFunc, out)
}

