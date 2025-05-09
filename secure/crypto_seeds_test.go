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

// Test SignPlaintext.
func TestSignPlaintext(t *testing.T) {
	cases := []struct {
		Case   string
		Inputs []string
		Want   string
	}{
		{"Append plaintexts", []string{"123", "abc", "ABC"}, "123\nabc\nABC"},
		{"Append empty plaintext", []string{"123", "", "ABC"}, "123\nABC"},
		{"Append empty start", []string{"", "abc", "ABC"}, "abc\nABC"},
		{"Append empty end", []string{"123", "abc", ""}, "123\nabc"},
	}

	for _, c := range cases {
		t.Run(c.Case, func(t *testing.T) {
			code := SignPlaintext(c.Inputs...)
			if code != c.Want {
				t.Fatal("Failed sign plaintexts!")
			}
		})
	}
}

func TestViaSignOne(t *testing.T) {
	signSeeds := "0aAbBcC1dDeEfF2gGhHiI3jJkKlL4mMnNoO5pPqQrR6sStTuU7vVwWxX8yYzZ9"
	CreateSeeds(signSeeds)

	t.Run("Verify Manual2", func(t *testing.T) {
		sign := "ghdWBIEJFuiKgKtL89dfNBfNX7hXKAQj85hP40UcbgC+rPIujfCcac1w6fz/wcdzr1dTAvR2zXfn1yegPnsYDCA="
		code := "5hUjz/MA"

		if !ViaSignCode(sign, code) {
			t.Fatal("Verify sign&code failed!")
		}
	})
}
