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

// `Step 2` : Bind request code and return wechat app url to get access token,
// please use [wx.login](https://developers.weixin.qq.com/miniprogram/dev/api/open-api/login/wx.login.html) get requestcode first.
//
// see more links
//
// - [Wechat app login follow](https://developers.weixin.qq.com/miniprogram/dev/framework/open-ability/login.html)
// - [Wechat app login API](https://developers.weixin.qq.com/miniprogram/dev/api-backend/open-api/login/auth.code2Session.html)
func (w *WxIFAgent) ToWxAppTokenUrl(requestcode string) string {
	tokenurl := wxauth2ApisUrlDomain +
		"/sns/jscode2session?appid=APPID&secret=SECRET&js_code=CODE&grant_type=authorization_code"
	tokenurl = strings.Replace(tokenurl, "APPID", w.AppID, -1)
	tokenurl = strings.Replace(tokenurl, "SECRET", w.AppSecret, -1)
	return strings.Replace(tokenurl, "CODE", requestcode, -1)
}

/* ------------------------------------------------------------------- */
/* Global API Access Token                                             */
/* ------------------------------------------------------------------- */

// Return API access token request url of target wechat app, this api maybe not stable,
// use the utils.GetApiTokenUrl() instead it.
//
// see more links
//
// - [getAccessToke](https://developers.weixin.qq.com/miniprogram/dev/server/API/mp-access-token/api_getaccesstoken.html)
// - [Access Token Usage](https://developers.weixin.qq.com/doc/oplatform/developers/dev/AccessToken.html)
func (w *WxIFAgent) GetAccessTokenUrl() string {
	tokenurl := wxauth2ApisUrlDomain +
		"/cgi-bin/token?appid=AppID&secret=AppSecret&grant_type=client_credential"
	tokenurl = strings.Replace(tokenurl, "APPID", w.AppID, -1)
	tokenurl = strings.Replace(tokenurl, "SECRET", w.AppSecret, -1)
	return tokenurl
}

// Return API access token request url of target wechat app,
// the wechat app id and secret params input as POST payload.
//
//	{
//	    "grant_type": "client_credential",
//	    "appid": "APPID",
//	    "secret": "APPSECRET"
//	}
//
// see more links
//
// - [getStableAccessToken](https://developers.weixin.qq.com/miniprogram/dev/server/API/mp-access-token/api_getstableaccesstoken.html)
func (w *WxIFAgent) GetApiTokenUrl() string {
	return wxauth2ApisUrlDomain + "/cgi-bin/stable_token"
}
