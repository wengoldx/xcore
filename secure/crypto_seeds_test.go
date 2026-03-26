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
	"os"
	"testing"
)

// -------------------------------------------------------------------
// USAGE: Enter ~/xcore/secure, and excute command to test.
//
//	go test -v -cover
// -------------------------------------------------------------------

const _test_seed = "0aAbBcC1dDeEfF2gGhHiI3jJkKlL4mMnNoO5pPqQrR6sStTuU7vVwWxX8yYzZ9" // "0123456789"
const _test_sign = "ghdWBIEJFuiKgKtL89dfNBfNX7hXKAQj85hP40UcbgC+rPIujfCcac1w6fz/wcdzr1dTAvR2zXfn1yegPnsYDCA="

func TestNewSeedSigns(t *testing.T) {
	s1 := NewSeedSign(_test_seed)
	s2 := NewSeedSign(_test_seed)

	for i, seed := range s1.seeds {
		if seed2, ok := s2.seeds[i]; !ok || seed != seed2 {
			t.Fatal("Not matched seed, at index:", i)
		}
		t.Log("Matched:", i, "-", seed)
	}
	t.Log("Seed Count:", len(s1.seeds))
}

func TestFilterDupChars(t *testing.T) {
	cases := []struct {
		Case string
		Src  string
		Out  string
	}{
		{"OK Src ", "1234567890", "1234567890"},
		{"Bad Src", "1234467890", "123467890"},
	}

	ss := SeedSign{}
	for _, c := range cases {
		t.Run(c.Case, func(t *testing.T) {
			out := ss.filterDupChars(c.Src)
			if out != c.Out {
				t.Fatal("Want:", c.Out, "but output:", out)
			}
			t.Log("Src:", c.Case, ">", out)
		})
	}
}

func TestSignPlaintext(t *testing.T) {
	cases := []struct {
		Case  string
		Texts []string
		Want  string
	}{
		{"Append plaintexts     ", []string{"123", "abc", "ABC"}, "123\nabc\nABC"},
		{"Append empty start    ", []string{"", "abc", "ABC"}, "abc\nABC"},
		{"Append empty plaintext", []string{"123", "", "ABC"}, "123\nABC"},
		{"Append empty end      ", []string{"123", "abc", ""}, "123\nabc"},
		{"Append emptys         ", []string{"", "", ""}, ""},
		{"All emptys            ", []string{}, ""},
	}

	ss := SeedSign{}
	for _, c := range cases {
		t.Run(c.Case, func(t *testing.T) {
			code := ss.SignPlaintext(c.Texts...)
			if code != c.Want {
				t.Fatal("Failed sign plaintexts!")
			}
		})
	}
}

func TestSignAndVerifyCode(t *testing.T) {
	cases := []struct {
		Case      string
		Type      string
		Bits      int
		Plaintext string
	}{
		{"Verify ECC sign code       ", "ECC", 0, "this is a plaintext to sign by ECC"},
		{"Verify ECC chinese         ", "ECC", 0, "中文编码字符串签名测试"},
		{"Verify RSA sign code [1024]", "RSA", 1024, "this is a plaintext to sign by RSA"},
		{"Verify RSA chinese   [1024]", "RSA", 1024, "中文编码字符串签名测试"},
		{"Verify RSA sign code [2048]", "RSA", 2048, "this is a plaintext to sign by RSA"},
		{"Verify RSA chinese   [2048]", "RSA", 2048, "中文编码字符串签名测试"},
	}

	ss := NewSeedSign(_test_seed)
	for _, c := range cases {
		t.Run(c.Case, func(t *testing.T) {
			sign := ""
			if c.Type == "ECC" {
				prikey, _ := NewEccPriKey()
				sign, _ = EccSign(c.Plaintext, prikey)
			} else {
				pri, _, _ := NewRSAKeys(c.Bits)
				sign, _ = RSASignB64(pri, c.Plaintext)
			}

			if code := ss.SignCode(sign); code == "" {
				t.Fatal("Failed get sign code!")
			} else if !ss.ViaCode(sign, code) {
				t.Fatal("Verify sign&code failed!")
			} else {
				t.Log("SignType:", c.Type, ">", code, "-", sign)
			}
		})
	}
}

func TestEccSignVerify(t *testing.T) {
	crtfile := "./_test_ecc.pem"
	plaintext := []string{"This a plaintext!", "Second text"}
	ss := SeedSign{}

	NewEccPemFile(crtfile)
	defer os.Remove(crtfile)
	sign, _ := ss.EccSign(crtfile, plaintext...)

	// check output signs whether same!
	for i := 0; i < 10; i++ {
		sign2, _ := ss.EccSign(crtfile, plaintext...)
		if sign == sign2 {
			t.Fatal("Exist same signs (ECC)!!")
		}
	}

	// verify sign and plaintext.
	prikey, _ := LoadEccPemFile(crtfile)
	pubpem, _ := EccPubString(&prikey.PublicKey)
	valid, _ := ss.EccVerify(sign, pubpem, plaintext...)
	if !valid {
		t.Fatal("Failed verify ECC sign!")
	}
	t.Log("Passed ECC sign & verify!")
}

func TestRsaSignVerify(t *testing.T) {
	plaintext := []string{"This a plaintext!", "Second text"}
	ss := SeedSign{}

	pri, pub, _ := NewRSAKeys(1024)
	sign, _ := ss.RsaSign(pri, plaintext...)

	// check output signs whether same!
	for i := 0; i < 10; i++ {
		sign2, _ := ss.RsaSign(pri, plaintext...)
		if sign != sign2 {
			t.Fatal("Exist different signs (RSA)!!")
		}
	}

	// verify sign and plaintext.
	valid, _ := ss.RsaVerify(sign, pub, plaintext...)
	if !valid {
		t.Fatal("Failed verify RSA sign!")
	}
	t.Log("Passed ECC sign & verify!")
}

func TestGenSeedCodes(t *testing.T) {
	ss := NewSeedSign(_test_seed)
	codes, conflicts := make(map[string]struct{}), 0
	for i, cnt := 0, 0; i < 10000; i++ { // test 10000 times.
		code := ss.SignCode(_test_sign)
		if _, ok := codes[code]; ok {
			conflicts++
			continue
		}
		cnt++
		codes[code] = struct{}{}
		fmt.Println("[", cnt, "]", "Code:", code)
	}
	fmt.Println("Conflicted", conflicts)
}

func TestViaSignOne(t *testing.T) {
	ss := NewSeedSign(_test_seed)
	for i := 0; i < 10000; i++ {
		code := ss.SignCode(_test_sign)
		if !ss.ViaCode(_test_sign, code) {
			t.Fatal("Verify sign&code failed!")
		}
		fmt.Println("Signed Code:", code)
	}
}

func TestSignToNum(t *testing.T) {
	ss := NewSeedSign(_test_seed)

	badchars := "!@#$%^&*()_+}{\":?><~}`-=[]\\;',./ "
	num, sign := "", _test_sign+badchars
	fmt.Println("sign src:", sign)
	for i := 0; i < 10; i++ { // test 10 times.
		if num == "" {
			_, num = ss.signSeedNum(sign)
			fmt.Println("sign num:", num)
			continue
		} else if _, n := ss.signSeedNum(sign); num != n {
			t.Fatal("Sign number not matched!")
		}
	}
}
