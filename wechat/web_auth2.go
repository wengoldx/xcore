// Copyright (c) 2018-Now Dunyu All Rights Reserved.
//
// Author      : https://www.wengold.net
// Email       : support@wengold.net
//
// Prismy.No | Date       | Modified by. | Description
// -------------------------------------------------------------------
// 00001       2019/05/22   yangping       New version
// -------------------------------------------------------------------

package wechat

import (
	"strings"
)

// WxIFAgent interfaces agent to using Wechat Official Account AppID and AppSecret to
// authenticate wechat user and get user profiles.
//
// DESCRIPTION FROM [Wechat](https://developers.weixin.qq.com/doc/offiaccount/OA_Web_Apps/Wechat_webpage_authorization.html) at 2019/07/05
//
// There are four steps in the process for webpage authorization with wechar OAuth2.0
//
// #### Step 1
//
// Sending the user to the authorization page to consent to authorization, obtain code,
//
//	`http` : GET (please use https as header),
//	`URL`  : "https://open.weixin.qq.com/connect/oauth2/authorize?appid=APPID&redirect_uri=REDIRECT_URI&response_type=code&scope=SCOPE&state=STATE#wechat_redirect"
//	`Description` : If the user agrees to authorization, the page will jump to redirect_uri/?code=CODE&state=STATE.
//
// #### Step2
//
// Use the code in exchange for the access_token of the webpage authorization (different from the access_token found in the basic support),
//
//	`http` : GET (please use https as header)
//	`URL`  : "https://api.weixin.qq.com/sns/oauth2/access_token?appid=APPID&secret=SECRET&code=CODE&grant_type=authorization_code"
//	`Description` : An accurate return JSON data includes the following: {
//		"access_token"  : "ACCESS_TOKEN",
//		"expires_in"    : 7200,
//		"refresh_token" : "REFRESH_TOKEN",
//		"openid"        : "OPENID",
//		"scope"         : "SCOPE"
//	}
//	Wechat will return JSON data as follows when there is an error: {
//		"errcode"       : 40029,
//		"errmsg"        : "invalid code"
//	}
//
// #### Step 3
//
// If necessary, the developer can refresh the webpage authorization access_token prevent it from expiring.
//
//	`http` : GET (please use https as header)
//	`URL`  : "https://api.weixin.qq.com/sns/oauth2/refresh_token?appid=APPID&grant_type=refresh_token&refresh_token=REFRESH_TOKEN"
//	`Description` : An accurate return JSON data includes the following: {
//		"access_token"  : "ACCESS_TOKEN",
//		"expires_in"    : 7200,
//		"refresh_token" : "REFRESH_TOKEN",
//		"openid"        : "OPENID",
//		"scope"         : "SCOPE"
//	}
//	Wechat will return JSON data as follows when there is an error: {
//		"errcode"       : 40029,
//		"errmsg"        : "invalid code"
//	}
//
// #### Step 4
//
// Use the webpage authorization access_token and openid to obtain basic user information.
//
//	`http` : GET (please use https as header)
//	`URL`  : "https://api.weixin.qq.com/sns/userinfo?access_token=ACCESS_TOKEN&openid=OPENID&lang=zh_CN"
//	`Description` : An accurate return JSON data includes the following: {
//	    "openid"     : "OPENID",
//	    "nickname"   : "NICKNAME",
//	    "sex"        : 1,
//	    "province"   : "PROVINCE"
//	    "city"       : "CITY",
//	    "country"    : "COUNTRY",
//	    "headimgurl" : "http://thirdwx.qlogo.cn/mmopen/g3MonUZtNHkdmzicIlibx6iaFqAc56vxLSUfpb6n5WKSYVY0ChQKkiaJSgQ1dZuTOgvLLrhJbERQQ4eMsv84eavHiaiceqxibJxCfHe/46",
//	    "privilege"  : ["PRIVILEGE1" "PRIVILEGE2"],
//	    "unionid"    : "o6_bmasdasdsad6_2sgVt7hMZOPfL"
//	}
//	Wechat will return JSON data as follows when there is an error: {
//	    "errcode" : 40003,
//	    "errmsg"  : "invalid openid"
//	}
//
// `Additional` :
//
// Testing the validity of the authorization certificate (access_token)
//
//	`http` : GET (please use https as header)
//	`URL`  : https://api.weixin.qq.com/sns/auth?access_token=ACCESS_TOKEN&openid=OPENID
//	`Description` : Accurate JSON return results: {
//		"errcode" : 0, "errmsg" : "ok"
//	}
//	Example of JSON returns when there are errors: {
//		"errcode" : 40003, "errmsg" : "invalid openid"
//	}
type WxIFAgent struct {
	AppID     string `json:"appid"`     // Wechat Official Account App ID
	AppSecret string `json:"appsecret"` // Wechat Official Account App Securet
	Scope     string `json:"scope"`     // Scope key of 'snsapi_base' or 'snsapi_userinfo'
	IsWxApp   bool   `json:"isapp"`     // Indicate wechat app or not, true is app
}

// WxToken wechat access and refresh tokens
type WxToken struct {
	AccessToken  string `json:"access_token"`
	Expires      int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	OpenID       string `json:"openid"`
	Scope        string `json:"scope"`
}

// WxUserInfo wechat user informations
type WxUserInfo struct {
	OpenID     string   `json:"openid"`
	Nickname   string   `json:"nickname"`
	Sex        int      `json:"sex"`
	Province   string   `json:"province"`
	City       string   `json:"city"`
	Country    string   `json:"country"`
	Headimgurl string   `json:"headimgurl"`
	Privilege  []string `json:"privilege"`
	UnionID    string   `json:"unionid"`
}

// WxResult request result return from server
type WxResult struct {
	ErrCode int    `json:"errcode"`
	Message string `json:"errmsg"`
}

// Agents map
var WxAgentsConfig map[string]*WxIFAgent

const (
	wxauth2OpenUrlDomain = "https://open.weixin.qq.com"
	wxauth2ApisUrlDomain = "https://api.weixin.qq.com/sns"
)

// `Step 1` : Bind redirect url and return wechat url to get request code
func (w *WxIFAgent) ToWxCodeUrl(redirecturl string, state ...string) string {
	codeurl := wxauth2OpenUrlDomain +
		"/connect/oauth2/authorize?appid=APPID&redirect_uri=REDIRECT_URI&response_type=code&scope=SCOPE&state=STATE#wechat_redirect"
	codeurl = strings.Replace(codeurl, "APPID", w.AppID, -1)
	codeurl = strings.Replace(codeurl, "REDIRECT_URI", redirecturl, -1)
	codeurl = strings.Replace(codeurl, "SCOPE", w.Scope, -1)

	// replace the STATE field by given state as optional param
	if len(state) > 0 && len(state[0]) > 0 {
		codeurl = strings.Replace(codeurl, "STATE", state[0], -1)
	}
	return codeurl
}

// `Step 2` : Bind request code and return wechat url to get access token
func (w *WxIFAgent) ToWxTokenUrl(requestcode string) string {
	tokenurl := wxauth2ApisUrlDomain +
		"/oauth2/access_token?appid=APPID&secret=SECRET&code=CODE&grant_type=authorization_code"
	tokenurl = strings.Replace(tokenurl, "APPID", w.AppID, -1)
	tokenurl = strings.Replace(tokenurl, "SECRET", w.AppSecret, -1)
	return strings.Replace(tokenurl, "CODE", requestcode, -1)
}

// `Step 3` : Bind expired access toke and return wechat url to refresh it
func (w *WxIFAgent) ToWxRefreshUrl(accesscode string) string {
	refreshurl := wxauth2ApisUrlDomain +
		"/oauth2/refresh_token?appid=APPID&grant_type=refresh_token&refresh_token=REFRESH_TOKEN"
	refreshurl = strings.Replace(refreshurl, "APPID", w.AppID, -1)
	return strings.Replace(refreshurl, "REFRESH_TOKEN", accesscode, -1)
}

// `Step 4` : Bind access token and openid, than return wechat url to get user informations
func (w *WxIFAgent) ToWxUserUrl(accesstoken, openid string) string {
	infourl := wxauth2ApisUrlDomain + "/userinfo?access_token=ACCESS_TOKEN&openid=OPENID&lang=zh_CN"
	infourl = strings.Replace(infourl, "ACCESS_TOKEN", accesstoken, -1)
	return strings.Replace(infourl, "OPENID", openid, -1)
}

// `Additional` : Bind access token and openid, than return wechat url to check access token expires
func (w *WxIFAgent) ToWxVerifyUrl(accesstoken, openid string) string {
	viaurl := wxauth2ApisUrlDomain + "/auth?access_token=ACCESS_TOKEN&openid=OPENID"
	viaurl = strings.Replace(viaurl, "ACCESS_TOKEN", accesstoken, -1)
	return strings.Replace(viaurl, "OPENID", openid, -1)
}
