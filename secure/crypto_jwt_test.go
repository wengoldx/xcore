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
	"strings"
	"testing"
	"time"
)

// -------------------------------------------------------------------
// USAGE: Enter ~/xcore/secure, and excute command to test.
//
//	go test -v -cover
// -------------------------------------------------------------------

// Test ViaJwtToken, NewJwtToken, EncClaims, DecClaims, NewSalt.
func TestViaJwtToken(t *testing.T) {
	cases := []struct {
		Case   string
		UUID   string
		Params string
		Dur    time.Duration
	}{
		{"Verify valid token", "12345678", "params,1,abc", time.Minute},
		{"Verify empty uuid", "", "params,1,abc", time.Minute},
		{"Verify empty params", "12345678", "", time.Minute},
		{"Verify invalid duration", "12345678", "params,1,abc", time.Millisecond},
	}

	for _, c := range cases {
		t.Run(c.Case, func(t *testing.T) {
			salt, _ := NewSalt()
			claims := EncClaims(c.UUID, c.Params)
			if token, err := NewJwtToken(claims, salt, c.Dur); err != nil {
				t.Fatal("New Jwt token, err:", err)
			} else if decode, err := ViaJwtToken(token, salt); err != nil {
				if strings.HasPrefix(err.Error(), "token is expired by") {
					t.Log("Verified expire Jwt token!")
				} else {
					t.Fatal("Verify Jwt token, err:", err)
				}
			} else if keywords, err := DecClaims(decode, 2); err != nil {
				t.Fatal("Decode claims, err:", err)
			} else if keywords[0] != c.UUID || keywords[1] != c.Params {
				t.Fatal("Jwt keywords invalid!")
			}
		})
	}
}
