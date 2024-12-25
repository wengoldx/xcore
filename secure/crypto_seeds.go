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
	"math/rand"
	"strings"
)

// Seeds map for get ecd sign code and verify it,
// call secure.CreateSeeds() to setup it before use.
var seedsMap map[int]string

// Seeds map radix for get seed string in valid range,
// it inited when call secure.CreateSeeds().
var seedRadix = 0

// Create seeds values and cached as a map, search by index.
//
// `WARNING`:
//
// The src param string MUST all unique chars:
//
//	- OK : '1234567890' // good src, all chars uniqued
//	- NG : '1234467890' // have double 4 chars
func CreateSeeds(src string) {
	if seedsMap == nil {
		seedsMap = make(map[int]string)

		// setup seeds maps
		sl, index := len(src), 0
		for i := 0; i < sl; i++ {
			for j := i + 1; j < sl; j++ {
				start, end := src[i:i+1], src[j:j+1]
				per, tail := src[:i]+src[i+1:j], src[j+1:]

				seedsMap[index] = per + end + start + tail
				index++
			}
		}

		// init seeds radix
		seedRadix = len(seedsMap)
	}
}

// Get sign string verify seed by sign string sum number.
func getSignSeed(sign string) string {
	sum := 0
	for _, char := range sign {
		sum += int(char)
	}
	return seedsMap[sum%seedRadix]
}

// Return a random code from sign string.
func GetSignCode(sign string) string {
	sl, seed := len(sign), getSignSeed(sign)
	if radix := len(seed); sl > radix {
		sl = radix
	}

	code := "" // random 4 group segements
	for seg := 0; seg < 4; seg++ {
		pos := rand.Intn(sl)
		code += seed[pos:pos+1] + sign[pos:pos+1]
	}
	return code
}

// Verify the auth code with sign string if valid.
//
// ---
//
//	seeds := "1234567890"
//	secure.CreateSeeds(seeds)
//
//	// get ecc sign string whereever came from
//	sign := "mhdWY0hJZmBLO4PnxTSWeUd2yqDUUgHyoAbnMjnOZpjo5IVlayfdrkDLsTquvj7nEkpqlSZCKIWx1OhuDq1ZLg=="
//	code := GetSignCode(sign)
//	rst  := VerifySignCode(sign, code) // rst == true
func ViaSignCode(sign, code string) bool {
	sl, cl := len(sign), len(code)
	if sl == 0 || cl%2 != 0 || sl < (cl/2) {
		return false
	}

	seed := getSignSeed(sign)
	for index := 0; index < cl-1; index += 2 {
		poschar, digital := code[index], code[index+1]
		pos := strings.Index(seed, string(poschar))
		if pos < 0 || pos >= sl {
			return false
		}

		if sign[pos] != digital {
			return false
		}
	}
	return true
}
