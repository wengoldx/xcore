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

// Test NewRSAKeys.
func TestNewRSAKeys(t *testing.T) {
	cases := []struct {
		Case string
		Bits int
	}{
		{"New RSA keys on 1024 bits", 1024},
		{"New RSA keys on 2048 bits", 2048},
	}

	for _, c := range cases {
		t.Run(c.Case, func(t *testing.T) {
			if pri, pub, err := NewRSAKeys(c.Bits); err != nil {
				t.Fatal("New RSA keys, err:", err)
			} else {
				t.Log("PriKey string:", "\n"+pri)
				t.Log("PubKey string:", "\n"+pub)
			}
		})
	}
}

// Test RSAEncryptB64, NewRSAKeys, RSAEncrypt.
func TestRSAEncryptB64(t *testing.T) {
	cases := []struct {
		Case string
		Bits int
	}{
		{"RSA encrypt on 1024 bits", 1024},
		{"RSA encrypt on 2048 bits", 2048},
	}

	for _, c := range cases {
		t.Run(c.Case, func(t *testing.T) {
			if _, pub, err := NewRSAKeys(c.Bits); err != nil {
				t.Fatal("New RSA keys, err:", err)
			} else if ciphertext, err := RSAEncryptB64([]byte(pub), []byte(c.Case)); err != nil {
				t.Fatal("RSA encrypt, err:", err)
			} else {
				t.Log("RSA encrypted:", "\n"+ciphertext)
			}
		})
	}
}

// Test RSADecrypt, NewRSAKeys, RSAEncrypt.
func TestRSADecrypt(t *testing.T) {
	cases := []struct {
		Case string
		Bits int
	}{
		{"RSA decrypt on 1024 bits", 1024},
		{"RSA decrypt on 2048 bits", 2048},
	}

	for _, c := range cases {
		t.Run(c.Case, func(t *testing.T) {
			if pri, pub, err := NewRSAKeys(c.Bits); err != nil {
				t.Fatal("New RSA keys, err:", err)
			} else if ciphertext, err := RSAEncrypt([]byte(pub), []byte(c.Case)); err != nil {
				t.Fatal("RSA encrypt, err:", err)
			} else if plaintext, err := RSADecrypt([]byte(pri), ciphertext); err != nil {
				t.Fatal("RSA decrypt, err:", err)
			} else if string(plaintext) != c.Case {
				t.Fatal("Failed verifid!")
			}
		})
	}
}

// Test RSASignB64, NewRSAKeys, RSAEncrypt.
func TestRSASignB64(t *testing.T) {
	cases := []struct {
		Case string
		Bits int
	}{
		{"RSA sign on 1024 bits", 1024},
		{"RSA sign on 2048 bits", 2048},
	}

	for _, c := range cases {
		t.Run(c.Case, func(t *testing.T) {
			if pri, _, err := NewRSAKeys(c.Bits); err != nil {
				t.Fatal("New RSA keys, err:", err)
			} else if sign, err := RSASignB64([]byte(pri), []byte(c.Case)); err != nil {
				t.Fatal("RSA sign, err:", err)
			} else {
				t.Log("RSA signed:", "\n"+sign)
			}
		})
	}
}

// Test RSAVerify, NewRSAKeys, RSASign.
func TestRSAVerify(t *testing.T) {
	cases := []struct {
		Case string
		Bits int
	}{
		{"RSA verify on 1024 bits", 1024},
		{"RSA verify on 2048 bits", 2048},
	}

	for _, c := range cases {
		t.Run(c.Case, func(t *testing.T) {
			if pri, pub, err := NewRSAKeys(c.Bits); err != nil {
				t.Fatal("New RSA keys, err:", err)
			} else if sign, err := RSASign([]byte(pri), []byte(c.Case)); err != nil {
				t.Fatal("RSA sign, err:", err)
			} else if err := RSAVerify([]byte(pub), []byte(c.Case), sign); err != nil {
				t.Fatal("RSA verify, err:", err)
			}
		})
	}
}

// Test NewRSAPemFiles, NewRSAKeys.
func TestNewRSAPemFiles(t *testing.T) {
	cases := []struct {
		Case    string
		PriFile string
		PubFile string
		Bits    int
	}{
		{"New RSA pem file on 1024 bits", "./rsa-pri.pem", "./rsa-pub.pem", 1024},
		{"New RSA pem file on 2048 bits", "./rsa-pri.pem", "./rsa-pub.pem", 2048},
	}

	for _, c := range cases {
		t.Run(c.Case, func(t *testing.T) {
			if err := NewRSAPemFiles(c.PriFile, c.PubFile, c.Bits); err != nil {
				t.Fatal("New RSA pem files, err:", err)
			}
			os.Remove(c.PriFile)
			os.Remove(c.PubFile)
		})
	}
}

// Test LoadRSAKey, NewRSAPemFiles, NewRSAKeys.
func TestLoadRSAKey(t *testing.T) {
	cases := []struct {
		Case    string
		PriFile string
		PubFile string
		Bits    int
	}{
		{"New RSA pem file on 1024 bits", "./rsa-pri.pem", "./rsa-pub.pem", 1024},
		{"New RSA pem file on 2048 bits", "./rsa-pri.pem", "./rsa-pub.pem", 2048},
	}

	for _, c := range cases {
		t.Run(c.Case, func(t *testing.T) {
			if err := NewRSAPemFiles(c.PriFile, c.PubFile, c.Bits); err != nil {
				t.Fatal("New RSA pem files, err:", err)
			}

			if pri, err := LoadRSAKey(c.PriFile, c.Bits); err != nil {
				t.Fatal("Load RSA pri pem file, err:", err)
			} else if pub, err := LoadRSAKey(c.PubFile, c.Bits); err != nil {
				t.Fatal("Load RSA pri pem file, err:", err)
			} else if len(pri) == 0 || len(pub) == 0 {
				t.Fatal("Invalid RSA pems!")
			}

			os.Remove(c.PriFile)
			os.Remove(c.PubFile)
		})
	}
}

// Test RSAEncrypt4FB64, RSAEncrypt4F, RSAEncrypt, LoadRSAKey.
func TestRSAEncrypt4FB64(t *testing.T) {
	cases := []struct {
		Case    string
		PriFile string
		PubFile string
		Bits    int
	}{
		{"RSA encrypt from pem file on 1024 bits", "./rsa-pri.pem", "./rsa-pub.pem", 1024},
		{"RSA encrypt from pem file on 2048 bits", "./rsa-pri.pem", "./rsa-pub.pem", 2048},
	}

	for _, c := range cases {
		t.Run(c.Case, func(t *testing.T) {
			if err := NewRSAPemFiles(c.PriFile, c.PubFile, c.Bits); err != nil {
				t.Fatal("New RSA pem files, err:", err)
			}

			if ciphertext, err := RSAEncrypt4FB64(c.PubFile, []byte(c.Case)); err != nil {
				t.Fatal("RSA encrypt from pem file, err:", err)
			} else {
				t.Log("RSA encrypted:", "\n"+ciphertext)
			}

			os.Remove(c.PriFile)
			os.Remove(c.PubFile)
		})
	}
}

// Test RSADecrypt4F, RSADecrypt, LoadRSAKey, RSAEncrypt4F, RSAEncrypt.
func TestRSADecrypt4F(t *testing.T) {
	cases := []struct {
		Case    string
		PriFile string
		PubFile string
		Bits    int
	}{
		{"RSA decrypt from pem file on 1024 bits", "./rsa-pri.pem", "./rsa-pub.pem", 1024},
		{"RSA decrypt from pem file on 2048 bits", "./rsa-pri.pem", "./rsa-pub.pem", 2048},
	}

	for _, c := range cases {
		t.Run(c.Case, func(t *testing.T) {
			if err := NewRSAPemFiles(c.PriFile, c.PubFile, c.Bits); err != nil {
				t.Fatal("New RSA pem files, err:", err)
			}

			if ciphertext, err := RSAEncrypt4F(c.PubFile, []byte(c.Case)); err != nil {
				t.Fatal("RSA encrypt from pem file, err:", err)
			} else if plaintext, err := RSADecrypt4F(c.PriFile, ciphertext); err != nil {
				t.Fatal("RSA decrypt from pem file, err:", err)
			} else if string(plaintext) != c.Case {
				t.Fatal("Failed verifid!")
			}

			os.Remove(c.PriFile)
			os.Remove(c.PubFile)
		})
	}
}

// Test RSAVerify4F, NewRSAPemFiles, RSASign4F.
func TestRSAVerify4F(t *testing.T) {
	cases := []struct {
		Case    string
		PriFile string
		PubFile string
		Bits    int
	}{
		{"RSA decrypt from pem file on 1024 bits", "./rsa-pri.pem", "./rsa-pub.pem", 1024},
		{"RSA decrypt from pem file on 2048 bits", "./rsa-pri.pem", "./rsa-pub.pem", 2048},
	}

	for _, c := range cases {
		t.Run(c.Case, func(t *testing.T) {
			if err := NewRSAPemFiles(c.PriFile, c.PubFile, c.Bits); err != nil {
				t.Fatal("New RSA pem files, err:", err)
			}

			if sign, err := RSASign4F(c.PriFile, []byte(c.Case)); err != nil {
				t.Fatal("RSA sign from pem file, err:", err)
			} else if err := RSAVerify4F(c.PubFile, []byte(c.Case), sign); err != nil {
				t.Fatal("RSA verify from pem file, err:", err)
			}

			os.Remove(c.PriFile)
			os.Remove(c.PubFile)
		})
	}
}

// Test RSAVerifyASN, NewRSAKeys, RSASign.
func TestRSAVerifyASN(t *testing.T) {
	cases := []struct {
		Case string
		Bits int
	}{
		{"RSA verify (ASN) on 1024 bits", 1024},
		{"RSA verify (ASN) on 2048 bits", 2048},
	}

	for _, c := range cases {
		t.Run(c.Case, func(t *testing.T) {
			if pri, pub, err := NewRSAKeys(c.Bits); err != nil {
				t.Fatal("New RSA keys, err:", err)
			} else if sign, err := RSASignASN([]byte(pri), []byte(c.Case)); err != nil {
				t.Fatal("RSA sign, err:", err)
			} else if err := RSAVerifyASN([]byte(pub), []byte(c.Case), sign); err != nil {
				t.Fatal("RSA verify (ASN), err:", err)
			}
		})
	}
}

// Test RSA2Verify, NewRSA2Keys, RSA2Sign.
func TestRSA2Verify(t *testing.T) {
	cases := []struct {
		Case string
		Bits int
	}{
		{"RSA verify (PKCS8) on 1024 bits", 1024},
		{"RSA verify (PKCS8) on 2048 bits", 2048},
	}

	for _, c := range cases {
		t.Run(c.Case, func(t *testing.T) {
			if pri, pub, err := NewRSA2Keys(c.Bits); err != nil {
				t.Fatal("New RSA PKCS8 keys, err:", err)
			} else if sign, err := RSA2Sign([]byte(pri), []byte(c.Case)); err != nil {
				t.Fatal("RSA PKCS8 sign, err:", err)
			} else if err := RSA2Verify([]byte(pub), []byte(c.Case), sign); err != nil {
				t.Fatal("RSA verify (PKCS8), err:", err)
			}
		})
	}
}
