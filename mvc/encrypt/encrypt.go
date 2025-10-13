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
	"strings"
	"time"

	"github.com/wengoldx/xcore/logger"
	"github.com/wengoldx/xcore/secure"
)

// For logout encrypt or decrypt logs.
var eclog = logger.CatLogger("Encrypter")

// ----------------------------------------------
// Encript & Decript Datas By AES Secure Methods
// ----------------------------------------------

/* !!! DO CHANGE THIS SECURE KEYS !!! */
var _secure = []byte("1234567890")
var _slat = "xxxx-xxxx-xxxx-xxxx"

// Setup the secure and slat datas before use the encrypter.
func Setup(secure, slat string) {
	_secure, _slat = []byte(secure), slat
}

// Encrypt plaintext data then return AES ciphertext, it may not encrypt
// when the given data is empty or encrypt error, and return empty string
// on invalid status.
//
// # NOTICE:
//
// This method will trim space chars from the input plaintext
// both start and end, example " ab c " string will trimed as "ab c".
func Encrypt(plaintext string) string {
	if plaintext == "" {
		return ""
	}

	plaintext = strings.TrimSpace(plaintext)
	if plaintext != "" {
		if ciphertext, err := secure.AESEncrypt(_secure, []byte(plaintext)); err != nil {
			eclog.E("Encrypt data:", plaintext, "err:", err)
		} else {
			return ciphertext
		}
	}
	return "" // not encrypt for empty data
}

// Decrypt AES ciphertext data then return plaintext, it may not decrypt
// when the given data is empty or decrypt error, and return empty string
// on invalid status.
func Decrypt(ciphertext string) string {
	if ciphertext != "" {
		if plaintext, err := secure.AESDecrypt(_secure, ciphertext); err != nil {
			eclog.E("Decrypt data:", ciphertext, "err:", err)
		} else {
			return plaintext
		}
	}
	return "" // not decrypt for empty data
}

// Verify JWT token and return JWT keywords.
func VerifyToken(jwttoken string) (string, error) {
	return secure.ViaJwtToken(jwttoken, _slat)
}

// Create a new JWT token with given expire time.
func NewToken(keyword string, expire time.Duration) (string, error) {
	return secure.NewJwtToken(keyword, _slat, expire)
}
