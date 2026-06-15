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

const (
	wxauth2OpenUrlDomain = "https://open.weixin.qq.com"
	wxauth2ApisUrlDomain = "https://api.weixin.qq.com"
)

// Wechat access and refresh tokens.
type WxToken struct {
	AccessToken  string `json:"access_token"`  // Wechat app login access token.
	Expires      int    `json:"expires_in"`    // Access token expire time.
	RefreshToken string `json:"refresh_token"` // Refreshed token.
	OpenID       string `json:"openid"`        // Wechat account openid.
	Scope        string `json:"scope"`         // Requesnt scope string.
}

// Wechat user informations.
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

// Request result return from server.
type WxResult struct {
	ErrCode int    `json:"errcode"`
	Message string `json:"errmsg"`
}

// Wechat app id and secure datas.
type WxSecret struct {
	AppID     string `json:"appid"`      // Wechat app id as 'APPID'.
	Secret    string `json:"secret"`     // Wechat app secret as 'APPSECRET'.
	GrantType string `json:"grant_type"` // Grant type, maybe fixed 'client_credential'.
}
