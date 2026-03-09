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
	"os"
	"testing"
)

// -------------------------------------------------------------------
// USAGE: Enter ~/xcore/secure, and excute command to test.
//
//	go test -v -cover
// -------------------------------------------------------------------

func TestNewSeedSign(t *testing.T) {
	signSeeds := "0aAbBcC1dDeEfF2gGhHiI3jJkKlL4mMnNoO5pPqQrR6sStTuU7vVwWxX8yYzZ9"
	s1 := NewSeedSign(signSeeds)
	s2 := NewSeedSign(signSeeds)

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

	for _, c := range cases {
		t.Run(c.Case, func(t *testing.T) {
			ss := SeedSign{}
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
		Case   string
		Data   string
		Extras []string
		Want   string
	}{
		{"Append plaintexts     ", "data", []string{"123", "abc", "ABC"}, "data\n123\nabc\nABC"},
		{"Append empty data     ", "", []string{"123", "abc", "ABC"}, "123\nabc\nABC"},
		{"Append empty plaintext", "data", []string{"123", "", "ABC"}, "data\n123\nABC"},
		{"Append empty start    ", "data", []string{"", "abc", "ABC"}, "data\nabc\nABC"},
		{"Append empty end      ", "data", []string{"123", "abc", ""}, "data\n123\nabc"},
		{"Append emptys         ", "", []string{"", "", ""}, ""},
		{"All emptys            ", "", []string{}, ""},
	}

	for _, c := range cases {
		t.Run(c.Case, func(t *testing.T) {
			code := SignPlaintext(c.Data, c.Extras...)
			if code != c.Want {
				t.Fatal("Failed sign plaintexts!")
			} else if SignPlaintext(c.Data) != c.Data {
				t.Fatal("Failed sign data!")
			}
		})
	}
}

func TestSignCode(t *testing.T) {
	cases := []struct {
		Case      string
		Type      string
		Plaintext string
	}{
		{"Verify ECC sign code", "ECC", "this is a plaintext to sign by ECC"},
		{"Verify ECC chinese  ", "ECC", "中文编码字符串签名测试"},
		{"Verify RSA sign code", "RSA", "this is a plaintext to sign by ECC"},
		{"Verify RSA chinese  ", "RSA", "中文编码字符串签名测试"},
	}

	signSeeds := "0aAbBcC1dDeEfF2gGhHiI3jJkKlL4mMnNoO5pPqQrR6sStTuU7vVwWxX8yYzZ9"
	ss := NewSeedSign(signSeeds)

	for _, c := range cases {
		t.Run(c.Case, func(t *testing.T) {
			sign := ""
			if c.Type == "ECC" {
				prikey, _ := NewEccPriKey()
				sign, _ = EccSign(c.Plaintext, prikey)
			} else {
				pri, _, _ := NewRSAKeys(2048)
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
	plaintext := "This a plaintext!"
	ss := SeedSign{}

	NewEccPemFile(crtfile)
	defer os.Remove(crtfile)
	sign, _ := ss.EccSign(plaintext, crtfile)

	// check output signs whether same!
	for i := 0; i < 10; i++ {
		sign2, _ := ss.EccSign(plaintext, crtfile)
		if sign == sign2 {
			t.Fatal("Can not out the same signs!")
		}
	}

	// verify sign and plaintext.
	prikey, _ := LoadEccPemFile(crtfile)
	pubpem, _ := EccPubString(&prikey.PublicKey)
	valid, _ := ss.EccVerify(plaintext, sign, pubpem)
	if !valid {
		t.Fatal("Failed verify ECC sign!")
	}
	t.Log("Passed ECC sign & verify!")
}

func TestRsaSignVerify(t *testing.T) {
	plaintext := "This a plaintext!"
	ss := SeedSign{}

	pri, pub, _ := NewRSAKeys(2048)
	sign, _ := ss.RsaSign(plaintext, pri)

	// check output signs whether same!
	for i := 0; i < 10; i++ {
		sign2, _ := ss.RsaSign(plaintext, pri)
		if sign != sign2 {
			t.Fatal("Can not out the same signs!")
		}
	}

	// verify sign and plaintext.
	valid, _ := ss.RsaVerify(plaintext, sign, pub)
	if !valid {
		t.Fatal("Failed verify RSA sign!")
	}
	t.Log("Passed ECC sign & verify!")
}

func TestViaSignOne(t *testing.T) {
	signSeeds := "0aAbBcC1dDeEfF2gGhHiI3jJkKlL4mMnNoO5pPqQrR6sStTuU7vVwWxX8yYzZ9"
	ss := NewSeedSign(signSeeds)

	t.Run("Verify Manual2", func(t *testing.T) {
		sign := "ghdWBIEJFuiKgKtL89dfNBfNX7hXKAQj85hP40UcbgC+rPIujfCcac1w6fz/wcdzr1dTAvR2zXfn1yegPnsYDCA="
		code := "5hUjz/MA"

		if !ss.ViaCode(sign, code) {
			t.Fatal("Verify sign&code failed!")
		}
	})
}
