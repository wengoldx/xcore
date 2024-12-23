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
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"math/big"

	"github.com/wengoldx/xcore/invar"
)

const (
	ECC_PEM_PRI_HEADER = "ECDSA PRIVATE KEY" // private key pem file header
	ECC_PEM_PUB_HEADER = "ECDSA PUBLIC KEY"  // public  key pem file header
)

// A string streaming to write string as writer.
type Stringer struct {
	data string
}

// Write data to stringer cache param.
func (s *Stringer) Write(p []byte) (n int, err error) {
	s.data += string(p)
	return len(p), nil
}

// Generate a ECC random private key, then you can get the pair
// public key from prikey.PublicKey param.
//
//	prikey, _ := secure.GenEccPriKey()
//	pubkey := &prikey.PublicKey // get public key
func GenEccPriKey() (*ecdsa.PrivateKey, error) {
	curve := elliptic.P256()
	prikey, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		return nil, err
	}
	return prikey, nil
}

// Format ECC private key to pem string, it can be save to file directly.
func EccPriString(prikey *ecdsa.PrivateKey) (string, error) {
	dertext, err := x509.MarshalECPrivateKey(prikey)
	if err != nil {
		return "", err
	}

	stringer := &Stringer{}
	block := &pem.Block{Type: ECC_PEM_PRI_HEADER, Bytes: dertext}
	if err := pem.Encode(stringer, block); err != nil {
		return "", err
	}
	return stringer.data, nil
}

// Format ECC public key to pem string, it can be save to file directly.
//
//	prikey, _ := secure.GenEccPriKey()
//	pubkey := &prikey.PublicKey              // get public key
//	pubstr, _ := secure.EccPubString(pubkey) // format public key to pem string
func EccPubString(pubkey *ecdsa.PublicKey) (string, error) {
	dertext, err := x509.MarshalPKIXPublicKey(&pubkey)
	if err != nil {
		return "", err
	}

	stringer := &Stringer{}
	block := &pem.Block{Type: ECC_PEM_PUB_HEADER, Bytes: dertext}
	if err := pem.Encode(stringer, block); err != nil {
		return "", err
	}
	return stringer.data, nil
}

// Generate ECC private key, and format to private and public key as pem string.
func GenEccKeys() (string, string, error) {
	prikey, err := GenEccPriKey()
	if err != nil {
		return "", "", err
	}

	pripem, err := EccPriString(prikey)
	if err != nil {
		return "", "", err
	}

	pubkey := &prikey.PublicKey
	pubpem, err := EccPubString(pubkey)
	if err != nil {
		return "", "", err
	}
	return pripem, pubpem, nil
}

// Get ECC private key from private pem string.
//
//	prikey, _ := GenEccPriKey()
//	pripem, _ := EccPriString(prikey)
//	newkey, _ := EccPriKey(pripem) // prikey == newkey
func EccPriKey(pripem string) (*ecdsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(pripem))
	pri, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return pri, nil
}

// Get ECC public key from public pem string.
//
//	prikey, _ := GenEccPriKey()
//	pubpem, _ := EccPubString(&prikey.PublicKey)
//	newkey, _ := EccPubKey(pubpem) // prikey.PublicKey == newkey
func EccPubKey(pubkey string) (*ecdsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(pubkey))
	pubif, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	pub, success := pubif.(*ecdsa.PublicKey)
	if !success {
		return nil, invar.ErrBadPublicKey
	}
	return pub, nil
}

// Parse ECC digital signs from signed string, to veriry plaintext.
//
//	prikey, _ := GenEccPriKey()
//	plaintext := "This is a plainttext to sign and verfiy!"
//	signb64, _ := EccSign(plaintext, prikey)
//	valid, _ : EccVerify(plaintext, signb64, &prikey.PublicKey)
//	fmt.Println("ECC verify result:", valid)
func EccDigitalSigns(sign []byte) (*big.Int, *big.Int) {
	if len(sign) != 64 {
		zero := big.NewInt(0)
		return zero, zero
	}

	rb, sb := make([]byte, 32), make([]byte, 32)
	for i := 0; i < len(sign); i++ {
		if i < 32 {
			rb[i] = sign[i]
		} else {
			sb[i-32] = sign[i]
		}
	}

	r, s := new(big.Int), new(big.Int)
	return r.SetBytes(rb), s.SetBytes(sb)
}

// Sign the given plaintext by ECC private key, and return the signed code
// on base64 format.
func EccSign(plaintext string, prikey *ecdsa.PrivateKey) (string, error) {
	hash := sha256.Sum256([]byte(plaintext))
	r, s, err := ecdsa.Sign(rand.Reader, prikey, hash[:])
	if err != nil {
		return "", err
	}

	sign, sb := r.Bytes(), s.Bytes()
	sign = append(sign, sb...)
	return ByteToBase64(sign), nil
}

// Verify the given plaintext by ECC public key and base64 formated sign code.
func EccVerify(plaintext, signb64 string, pubkey *ecdsa.PublicKey) (bool, error) {
	signs, err := Base64ToByte(signb64)
	if err != nil {
		return false, err
	}

	r, s := EccDigitalSigns(signs)
	hash := sha256.Sum256([]byte(plaintext))
	valid := ecdsa.Verify(pubkey, hash[:], r, s)
	return valid, nil
}
