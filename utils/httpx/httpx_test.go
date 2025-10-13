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
	"testing"
)

func TestGet(t *testing.T) {
	req := "http://192.168.1.192:3103/accservice/debug/token?uuid=20000680"
	var token []byte
	if err := Get(req, &token); err != nil {
		t.Fatal(err)
	}
	t.Log("Token string:", string(token))
	t.Log("Token bytes:", token)
}

func TestGetString(t *testing.T) {
	req := "http://192.168.1.192:3103/accservice/debug/token?uuid=%s"
	params := "20000680" // FIXME: maybe changed!

	var token string
	if err := Get(req, &token, params); err != nil {
		t.Fatal(err)
	}
	t.Log("Token string:", token)
}

type MyStruct struct {
	OpenID  string `json:"openid"`
	Token   string `json:"token"`
	UnionID string `json:"unionid"`
	UUID    string `json:"uuid"`
}

func TestGetStruct(t *testing.T) {
	req := "http://192.168.1.192:3103/accservice/v4/wx/login/unionid?id=%s"
	params := "oRWA1645rQoiHAkp7CXODTggEpIY" // FIXME: maybe changed!

	var out MyStruct
	if err := Get(req, &out, params); err != nil {
		t.Fatal(err)
	}
	t.Log("Out struct:", out)
}

func TestGetStruct2(t *testing.T) {
	req := "http://192.168.1.192:3103/accservice/v4/wx/login/unionid?id=%s"
	params := "oRWA1645rQoiHAkp7CXODTggEpIY" // FIXME: maybe changed!

	var out struct {
		OpenID  string `json:"openid"`
		Token   string `json:"token"`
		UnionID string `json:"unionid"`
		UUID    string `json:"uuid"`
	}
	if err := Get(req, &out, params); err != nil {
		t.Fatal(err)
	}
	t.Log("Out struct:", out)
}
