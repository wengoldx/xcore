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

// Test ViaLoginToken, EncLoginToken.
func TestViaLoginToken(t *testing.T) {
	cases := []struct {
		Case string
		Acc  string
		Pwd  string
	}{
		{"Verify full params", "acc-1234", "pwd-321"},
		{"Verify only password", "", "pwd-321"},
		{"Verify only account", "acc-1234", ""},
		{"Verify empty data", "", ""},
	}

	for _, c := range cases {
		t.Run(c.Case, func(t *testing.T) {
			token := EncLoginToken(c.Acc, c.Pwd)
			if valid, err := ViaLoginToken(c.Acc, c.Pwd, token, 5); err != nil {
				t.Fatal("Verify custom token, err:", err)
			} else {
				t.Log("Verified:", valid)
			}
		})
	}
}
