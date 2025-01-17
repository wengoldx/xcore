// Copyright (c) 2018-Now Dunyu All Rights Reserved.
//
// Author      : https://www.wengold.net
// Email       : support@wengold.net
//
// Prismy.No | Date       | Modified by. | Description
// -------------------------------------------------------------------
// 00001       2019/05/22   yangping       New version
// -------------------------------------------------------------------

package secure

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/wengoldx/xcore/invar"
)

// Encode account and password as a simple login token.
//
// ---
//
//  account   password
//      |- + -|
//         |
//      base64          current nanosecode
//         |                    |
//        md5                base64
//         +------- "."---------|
//                   |
//                base64 => token
func EncLoginToken(acc, pwd string) string {
	timestamp := fmt.Sprintf("%v", time.Now().UnixNano())
	origin := EncodeB64MD5(acc+"."+pwd) + "." + EncodeBase64(timestamp)
	return EncodeBase64(origin)
}

// Deprecated: use utils.EncLoginToken instead it.
func GenLoginToken(acc, pwd string) string { return EncLoginToken(acc, pwd) }

// Verify login token.
//
// ---
//        token => base64
//                   |
//         +------- "."---------|
//        md5                base64
//         |                    |
//      base64          current nanosecode
//         |
//      |- + -|
//  account   password
func ViaLoginToken(acc, pwd, token string, duration int64) (bool, error) {
	origin, err := DecodeBase64(token)
	if err != nil {
		return false, err
	}

	segments := strings.Split(string(origin), ".")
	if len(segments) == 2 {
		if segments[0] != EncodeB64MD5(acc+"."+pwd) {
			return false, nil
		}

		latestByte, err := DecodeBase64(segments[1])
		if err != nil {
			return false, err
		}
		latest, err := strconv.ParseInt(string(latestByte), 10, 64)
		if err != nil {
			return false, err
		}

		// check token period
		if time.Now().UnixNano()-latest <= duration {
			return true, nil
		}
		return false, invar.ErrTokenExpired
	}
	return false, nil
}
