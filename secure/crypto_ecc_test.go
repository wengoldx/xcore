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

// Test NewEccKeys, NewEccPriKey, EccKeysString.
func TestNewEccKeys(t *testing.T) {
	cases := []struct {
		Case  string
		Curve string
	}{
		{"New ECC secure keys as default", ""},
		{"New ECC secure keys by P224", "P224"},
		{"New ECC secure keys by P384", "P384"},
		{"New ECC secure keys by P521", "P521"},
		{"New ECC secure keys by ????", "????"}, // as default P256
	}

	for _, c := range cases {
		t.Run(c.Case, func(t *testing.T) {
			if pri, pub, err := NewEccKeys(c.Curve); err != nil {
				t.Fatal("New ECC secure keys, err:", err)
			} else {
				t.Log("PriKey String:", "\n"+pri)
				t.Log("PubKey String:", "\n"+pub)
			}
		})
	}
}

// Test NewEccPemFile, LoadEccPemFile, EccPriKey.
func TestLoadEccPemFile(t *testing.T) {
	cases := []struct {
		Case    string
		OutFile string
		Curve   string
	}{
		{"New ECC pem file as default", "./pri_def.pem", ""},
		{"New ECC pem file by P224", "./pri_P224.pem", "P224"},
		{"New ECC pem file by P384", "./pri_P384.pem", "P384"},
		{"New ECC pem file by P521", "./pri_P521.pem", "P521"},
		{"New ECC pem file by ????", "./pri_err.pem", "????"}, // as default P256
	}

	for _, c := range cases {
		t.Run(c.Case, func(t *testing.T) {
			if err := NewEccPemFile(c.OutFile, c.Curve); err != nil {
				t.Fatal("New ECC pem file, err:", err)
			}

			if _, err := LoadEccPemFile(c.OutFile); err != nil {
				t.Fatal("Load ECC pem file, err:", err)
			}
			os.Remove(c.OutFile)
		})
	}
}

// Test EccPubKey, NewEccPriKey, EccPubString.
func TestEccPubKey(t *testing.T) {
	cases := []struct {
		Case  string
		Curve string
	}{
		{"New ECC PubKey as default", ""},
		{"New ECC PubKey by P224", "P224"},
		{"New ECC PubKey by P384", "P384"},
		{"New ECC PubKey by P521", "P521"},
		{"New ECC PubKey by ????", "????"}, // as default P256
	}

	for _, c := range cases {
		t.Run(c.Case, func(t *testing.T) {
			if pri, err := NewEccPriKey(c.Curve); err != nil {
				t.Fatal("New ECC PriKey, err:", err)
			} else if pem, err := EccPubString(&pri.PublicKey); err != nil {
				t.Fatal("Trans PubKey to PubPem, err:", err)
			} else if _, err := EccPubKey(pem); err != nil {
				t.Fatal("Trans PubPem to PubKey, err:", err)
			}
		})
	}
}

// Test NewEccPriKey, EccSign, EccVerify, EccDigitalSigns.
func TestEccVerify(t *testing.T) {
	cases := []struct {
		Case  string
		Curve string
	}{
		{"Verify sign text by default key", ""},
		{"Verify sign text by P224", "P224"},
		{"Verify sign text by P384", "P384"},
		{"Verify sign text by P521", "P521"},
		{"Verify sign text by ????", "????"}, // as default P256
	}

	for _, c := range cases {
		t.Run(c.Case, func(t *testing.T) {
			if prikey, err := NewEccPriKey(c.Curve); err != nil {
				t.Fatal("New ECC PriKey, err:", err)
			} else if sign, err := EccSign(c.Case, prikey); err != nil {
				t.Fatal("Sign plaintext, err:", err)
			} else if verify, err := EccVerify(c.Case, sign, &prikey.PublicKey); err != nil {
				t.Fatal("Verify sign text, err:", err)
			} else if !verify {
				t.Fatal("Failed verifid!")
			}
		})
	}
}
