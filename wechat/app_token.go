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
		"/jscode2session?appid=APPID&secret=SECRET&js_code=CODE&grant_type=authorization_code"
	tokenurl = strings.Replace(tokenurl, "APPID", w.AppID, -1)
	tokenurl = strings.Replace(tokenurl, "SECRET", w.AppSecret, -1)
	return strings.Replace(tokenurl, "CODE", requestcode, -1)
}
