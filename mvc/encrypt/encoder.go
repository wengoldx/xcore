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

	"github.com/wengoldx/xcore/logger"
	"github.com/wengoldx/xcore/secure"
)

// Data encrypt & decrypt utils.
type Encoder []byte

// ----------------------------------------------
// Encript & Decript Datas By AES Secure Methods
// ----------------------------------------------

// Global default encoder instance for easy used.
var _def_encoder *Encoder

// For logout encrypt or decrypt logs.
var eclog = logger.CatLogger("Encrypter")

// New encoder instance for data encrypt and decrypt by AES keys.
//
//	// Useage 1:
//	enc.NewEncoder("xxxx", true)
//	enc.Encrypt(plaintext) // use default encoder.
//	enc.Decrypt(plaintext) // use default encoder.
//
//	// Useage 2:
//	custom := enc.NewEncoder("xxxx")
//	custom.Encrypt(plaintext) // use custom encoder.
//	custom.Decrypt(plaintext) // use custom encoder.
func NewEncoder(secure string, asdef ...bool) *Encoder {
	encoder := (Encoder)([]byte(secure))
	if len(asdef) > 0 && asdef[0] {
		_def_encoder = &encoder
	}
	return &encoder
}

// Use default encoder to encrypt plaintext.
//
// # WARING:
//	- The default encoder must init by enc.NewEncoder(xx, ture)!
func Encrypt(plaintext string) string {
	if _def_encoder != nil {
		return _def_encoder.Encrypt(plaintext)
	}
	return ""
}

// Use default encoder to decrypt ciphertext.
//
// # WARING:
//	- The default encoder must init by enc.NewEncoder(xx, ture)!
func Decrypt(ciphertext string) string {
	if _def_encoder != nil {
		return _def_encoder.Decrypt(ciphertext)
	}
	return ""
}

// Encrypt plaintext data then return AES ciphertext, it may not encrypt
// when the given data is empty or encrypt error, and return empty string
// on invalid status.
//
// # NOTICE:
//
// This method will trim space chars from the input plaintext
// both start and end, example " ab c " string will trimed as "ab c".
func (e *Encoder) Encrypt(plaintext string) string {
	if plaintext == "" {
		return ""
	}

	plaintext = strings.TrimSpace(plaintext)
	if plaintext != "" {
		secure_key, datas := ([]byte)(*e), []byte(plaintext)
		if ciphertext, err := secure.AESEncrypt(secure_key, datas); err != nil {
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
func (e *Encoder) Decrypt(ciphertext string) string {
	if ciphertext != "" {
		secure_key := ([]byte)(*e)
		if plaintext, err := secure.AESDecrypt(secure_key, ciphertext); err != nil {
			eclog.E("Decrypt data:", ciphertext, "err:", err)
		} else {
			return plaintext
		}
	}
	return "" // not decrypt for empty data
}
