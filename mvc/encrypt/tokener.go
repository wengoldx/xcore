// Copyright (c) 2018-Now Dunyu All Rights Reserved.
//
// Author      : https://www.quantkernel.com
// Email       : ping.yang@wengold.net
//
// Prismy.No | Date       | Modified by. | Description
// -------------------------------------------------------------------
// 00001       2024/11/18   youhei         New version
// -------------------------------------------------------------------

package enc

import (
	"errors"
	"time"

	"github.com/wengoldx/xcore/secure"
)

// JWT token utils.
type Tokener string

// ----------------------------------------------
// JWT Token Utils
// ----------------------------------------------

// Global default tokener instance for easy used.
var _def_tokener *Tokener

// New tokener instance for generate and verify JWT tokens.
//
//	// Useage 1:
//	enc.NewTokener("xxxx", true)
//	enc.VerifyToken(jwttoken)     // use default tokener.
//	enc.NewToken(keyword, expire) // use default tokener.
//
//	// Useage 2:
//	custom := enc.NewTokener("xxxx")
//	custom.VerifyToken(jwttoken)     // use custom tokener.
//	custom.NewToken(keyword, expire) // use custom tokener.
func NewTokener(salt string, asdef ...bool) *Tokener {
	tokener := (Tokener)(salt)
	if len(asdef) > 0 && asdef[0] {
		_def_tokener = &tokener
	}
	return &tokener
}

// Use default tokener to verify JWT token.
//
// # WARING:
//	- The default tokener must init by enc.NewTokener(xx, ture)!
func VerifyToken(jwttoken string) (string, error) {
	if _def_tokener != nil {
		return _def_tokener.VerifyToken(jwttoken)
	}
	return "", errors.New("Not inited!")
}

// Use default tokener to generate a new JWT token.
//
// # WARING:
//	- The default tokener must init by enc.NewTokener(xx, ture)!
func NewToken(keyword string, expire time.Duration) (string, error) {
	if _def_tokener != nil {
		return _def_tokener.NewToken(keyword, expire)
	}
	return "", errors.New("Not inited!")
}

// Verify JWT token and return JWT keywords.
func (t *Tokener) VerifyToken(jwttoken string) (string, error) {
	return secure.ViaJwtToken(jwttoken, string(*t))
}

// Create a new JWT token with given expire time.
func (t *Tokener) NewToken(keyword string, expire time.Duration) (string, error) {
	return secure.NewJwtToken(keyword, string(*t), expire)
}
