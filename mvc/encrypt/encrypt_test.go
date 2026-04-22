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
	"fmt"
	"testing"
	"time"

	"github.com/wengoldx/xcore/secure"
)

func TestEncoder(t *testing.T) {
	plaintext := "12345678"
	Setup(secure.NewAESKey(), "1234567890")
	encode := Encrypt(plaintext)
	decode := Decrypt(encode)
	if plaintext != decode {
		t.Fatal("Encode & Decode unmatched!")
	}
	fmt.Println("Plaintext:", plaintext, "> encoded:", encode, "- decoded:", decode)
}

func TestTokener(t *testing.T) {
	keyword := "user:12345"
	Setup(secure.NewAESKey(), "1234567890")
	token, _ := NewToken(keyword, 300*time.Second)
	claims, _ := VerifyToken(token)
	if keyword != claims {
		t.Fatal("Keyword & Claims unmathed!")
	}
	fmt.Println("Keyword:", keyword, "- Claims:", claims)
}
