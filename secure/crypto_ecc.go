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
	"os"

	"github.com/wengoldx/xcore/invar"
)

/**
 * It's a ECC secure utils to support generate keys and sign data, then verify it.
 *
 * Here have 3 importent datas:
 * - The ecc private key, it contain a pair public key.
 * - The sign string, use private key and plaintext to signed result.
 * - The plaintext data to sign and verify compare check.
 *
 * ---
 *
 * `WARNING`:
 *
 * ECC not best for encript/decript, but better for sign/verify, if you want
 * encript/decript data with high performence, please use RSA to implement them.
 *
 * `USAGE`:
 *
 * 1. Call secure.NewEccPriKey() create a ecc private key.
 * 2. Call secure.EccKeysString(prikey) return private and public keys pem datas to save.
 * 3. Call secure.EccSign(plaintext, prikey) sign plaintext.
 * 4. Call secure.EccVerify(plaintext, signstring, pubkey) to verify valid.
 *
 * `Extend`:
 *
 * - Call secure.EccPriKey(pripem) return private key from pem data.
 * - Call secure.EccPubKey(pubpem) return public key from pem data.
 * - Call secure.EccDigitalSigns(signstring) return R and S digitals.
 */

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

// Deprecated: use utils.NewEccPriKey instead it.
func GenEccPriKey(sign ...string) (*ecdsa.PrivateKey, error) { return NewEccPriKey(sign...) }

// Create a ECC random private key with curve type one of P224, P256, P384,
// P521, or use P256 curve as default, then you can get the pair public key
// from prikey.PublicKey param.
//
//	prikey, _ := secure.NewEccPriKey() // same as secure.NewEccPriKey("P256")
//	pubkey := &prikey.PublicKey        // get public key
func NewEccPriKey(sign ...string) (*ecdsa.PrivateKey, error) {
	curvetype := "P256"
	if len(sign) > 0 && sign[0] != "" {
		curvetype = sign[0]
	}

	var curve elliptic.Curve
	switch curvetype {
	case "P224":
		curve = elliptic.P224() // FIXME: Sign<->Verify Error!
	case "P384":
		curve = elliptic.P384() // FIXME: Sign<->Verify Error!
	case "P521":
		curve = elliptic.P521() // FIXME: Sign<->Verify Error!
	default: // P256 as default
		curve = elliptic.P256()
	}

	// generate random ecc private key
	prikey, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		return nil, err
	}
	return prikey, nil
}

// Deprecated: use utils.NewEccKeys instead it.
func GenEccKeys(sign ...string) (string, string, error) { return NewEccKeys(sign...) }

// Create ECC private key, and format private and public keys as pem strings,
// by default it create P256 curve to sign data, or you can create other keys
// for set sign param as P224, P384, P521.
//
// @see secure.NewEccPriKey()
func NewEccKeys(sign ...string) (string, string, error) {
	prikey, err := NewEccPriKey(sign...)
	if err != nil {
		return "", "", err
	}
	return EccKeysString(prikey)
}

// Create ECC private key and save to target pem file.
func NewEccPemFile(outfile string, sign ...string) error {
	if prikey, err := NewEccPriKey(sign...); err != nil {
		return err
	} else if pripem, err := EccPriString(prikey); err != nil {
		return err
	} else {
		return os.WriteFile(outfile, []byte(pripem), 0666)
	}
}

// Load ECC private pem file and return private key.
func LoadEccPemFile(pemfile string) (*ecdsa.PrivateKey, error) {
	pripem, err := os.ReadFile(pemfile)
	if err != nil {
		return nil, err
	}
	return EccPriKey(string(pripem))
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
//	prikey, _ := secure.NewEccPriKey()
//	pubkey := &prikey.PublicKey              // get public key
//	pubstr, _ := secure.EccPubString(pubkey) // format public key to pem string
func EccPubString(pubkey *ecdsa.PublicKey) (string, error) {
	dertext, err := x509.MarshalPKIXPublicKey(pubkey)
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

// Format ECC private and public keys to pem strings.
func EccKeysString(prikey *ecdsa.PrivateKey) (string, string, error) {
	if pripem, err := EccPriString(prikey); err != nil {
		return "", "", err
	} else if pubpem, err := EccPubString(&prikey.PublicKey); err != nil {
		return "", "", err
	} else {
		return pripem, pubpem, nil
	}
}

// Get ECC private key from private pem string.
//
//	prikey, _ := NewEccPriKey()
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
//	prikey, _ := NewEccPriKey()
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
//	prikey, _ := NewEccPriKey()
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
