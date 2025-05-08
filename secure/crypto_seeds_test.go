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
	"testing"
)

// -------------------------------------------------------------------
// USAGE: Enter ~/xcore/secure, and excute command to test.
//
//	go test -v -cover
// -------------------------------------------------------------------

// Test ViaSignCode, GetSignCode, CreateSeeds.
func TestViaSignCode(t *testing.T) {
	cases := []struct {
		Case      string
		Plaintext string
	}{
		{"Verify sign code", "this is a plaintext to sign by ECC"},
		{"Verify chinese", "中文编码字符串签名测试"},
	}

	signSeeds := "0aAbBcC1dDeEfF2gGhHiI3jJkKlL4mMnNoO5pPqQrR6sStTuU7vVwWxX8yYzZ9"
	CreateSeeds(signSeeds)

	for _, c := range cases {
		t.Run(c.Case, func(t *testing.T) {
			prikey, _ := NewEccPriKey()
			sign, _ := EccSign(c.Plaintext, prikey)

			if code := GetSignCode(sign); code == "" {
				t.Fatal("Failed get sign code!")
			} else if !ViaSignCode(sign, code) {
				t.Fatal("Verify sign&code failed!")
			}
		})
	}
}
