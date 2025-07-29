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
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"os"
	"time"

	"github.com/wengoldx/xcore/invar"
)

// ### 1. How to encrypt and decrypt by RSA
//
// - (1). use secure.NewRSAKeys() to generate RSA keys, and set content bits length.
//
// - (2). use secure.RSAEncrypt() to encrypt original data with given public key.
//
// - (3). use secure.RSADecrypt() to decrypt ciphertext with given private key.
//
// `USAGE`
//
//	// Use the pubkey to encrypt and use the prikey to decrypt
//	prikey, pubkey, _ := secure.NewRSAKeys(1024)
//	logger.I("public  key:", pubkey, "private key:", prikey)
//
//	ciphertext, _ := secure.RSAEncrypt([]byte(pubkey), []byte("original-content"))
//	ciphertextBase64 := secure.EncodeBase64(string(ciphertext))
//	logger.I("ciphertext base64 string:", ciphertextBase64)
//
//	original, _ := secure.RSADecrypt([]byte(prikey), ciphertext)
//	logger.I("original string:", string(original))	// print 'original-content'
//
//
// ----
//
//
// ### 2. How to digital signature and verify by RSA
//
// - (1). use secure.NewRSAKeys() to generate RSA keys, and set content bits length.
//
// - (2). use secure.RSASign() to make digital signature with given private key.`
//
// - (3). use secure.RSAVerify() to verify data's integrity with given public key and digital signature
//
// `USAGE`
//
//	// Use the private key to create digital signature and use pubkey to verify it
//	prikey, pubkey, _ := secure.NewRSAKeys(1024)
//	logger.I("public  key:", pubkey, "private key:", prikey)
//
//	original := []byte("original-content")
//	signature, _ := secure.RSASign([]byte(prikey), original)
//	logger.I("original string:", string(original))
//	logger.I("signature string:", string(signature))
//
//	if err := secure.RSAVerify([]byte(pubkey), original, signature); err != nil {
//		logger.E("Verify failed with err:", err)
//		return
//	}
//	logger.I("Verify success")
const RSA_UTIL_DESCRIPTION = 0 /* just use for description */

const (
	RSA_PEM_PRI_HEADER = "RSA PRIVATE KEY" // private key pem file header
	RSA_PEM_PUB_HEADER = "RSA PUBLIC KEY"  // public  key pem file header
)

// Load RSA private or public key content from the given pem file,
// and the input buffer size of bits must larger than pem file size
// by call NewRSAKeys to set bits.
func LoadRSAKey(filepath string, bits ...int) (string, error) {
	if len(bits) > 0 && bits[0] > 0 {
		pemfile, err := os.Open(filepath)
		if err != nil {
			return "", err
		}
		defer pemfile.Close()

		keybuf := make([]byte, bits[0])
		num, err := pemfile.Read(keybuf)
		if err != nil {
			return "", err
		}
		return string(keybuf[:num]), nil
	} else {
		pemfile, err := os.ReadFile(filepath)
		if err != nil {
			return "", err
		}
		return string(pemfile), nil
	}
}

// Parse PKCS#1 or PKCS#8 RSA private key from pem file, set pkcs8
// param to true for use #PCSC#8, or false as default for use PKCS#1
// to parse RSA private key.
func ParsePriKey(prikey string, pkcs8 ...bool) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(prikey))
	if block == nil {
		return nil, invar.ErrBadPrivateKey
	}

	if len(pkcs8) > 0 && pkcs8[0] {
		priif, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		return priif.(*rsa.PrivateKey), nil
	}

	pri, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return pri, nil
}

// Parse RSA private key from pem file whatever PKSC#1 or PKCS#8 format.
func ParsePubKey(pubkey string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(pubkey))
	if block == nil {
		return nil, invar.ErrBadPublicKey
	}

	pubif, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return pubif.(*rsa.PublicKey), nil
}

// -------------------------------------------------------------------
// Create a RSA key as PKCS#1, ASN.1 format and encrypt by public
// key, than decrypt by private key.
// -------------------------------------------------------------------

// Create RSA private and public keys in PKCS#1, ASN.1 DER format,
// and limit bits length of key cert.
//	@param bits Limit bits length of key cert
//	@return - string Private key original string
//			- string Public key original string
//			- error Exception message
func NewRSAKeys(bits ...int) (string, string, error) {
	return newRSAKeysByType("PKCS1", bits...)
}

// Create RSA private and public keys, then save to target pem files.
func NewRSAPemFiles(prifile, pubfile string, bits ...int) error {
	if prikey, pubkey, err := NewRSAKeys(bits...); err != nil {
		return err
	} else if err := os.WriteFile(prifile, []byte(prikey), 0666); err != nil {
		return err
	} else if err := os.WriteFile(pubfile, []byte(pubkey), 0666); err != nil {
		if info, err := os.Stat(prifile); err == nil && info != nil {
			os.Remove(prifile)
		}
		return err
	}
	return nil
}

// Using RSA public key to encrypt original data.
//
//	prikey, pubkey, _ := secure.NewRSAKeys(1024)
//	ciphertext, _ := secure.RSAEncrypt(pubkey, "original-content")
//	original, _ := secure.RSADecrypt(prikey, ciphertext)
//	// original == 'original-content'
func RSAEncrypt(pubkey, original string) ([]byte, error) {
	pub, err := ParsePubKey(pubkey)
	if err != nil {
		return nil, err
	}
	return rsa.EncryptPKCS1v15(rand.Reader, pub, []byte(original))
}

// Using RSA public key to encrypt original data,
// then format to base64 form.
//
//	prikey, pubkey, _ := secure.NewRSAKeys(1024)
//	ciphertextBase64, _ := secure.RSAEncryptB64(pubkey, "original-content")
//	ciphertext, _ := secure.DecodeBase64(ciphertextBase64)
//	original, _ := secure.RSADecrypt(prikey, ciphertext)
//	// original == 'original-content'
func RSAEncryptB64(pubkey, original string) (string, error) {
	buf, err := RSAEncrypt(pubkey, original)
	if err != nil {
		return "", nil
	}
	return ByteToBase64(buf), nil
}

// Using RSA public key file to encrypt original data.
func RSAEncrypt4F(pubfile, original string) ([]byte, error) {
	pubkey, err := LoadRSAKey(pubfile)
	if err != nil {
		return nil, err
	}
	return RSAEncrypt(pubkey, original)
}

// Using RSA public key file to encrypt original data,
// then format to base64 form.
func RSAEncrypt4FB64(pubfile, original string) (string, error) {
	buf, err := RSAEncrypt4F(pubfile, original)
	if err != nil {
		return "", err
	}
	return ByteToBase64(buf), nil
}

// Using RSA private key to decrypt ciphertext.
//
//	prikey, pubkey, _ := secure.NewRSAKeys(1024)
//	ciphertext, _ := secure.RSAEncrypt(pubkey, "original-content")
//	original, _ := secure.RSADecrypt(prikey, ciphertext)
//	// original == 'original-content'
func RSADecrypt(prikey string, ciphertext []byte) ([]byte, error) {
	pri, err := ParsePriKey(prikey)
	if err != nil {
		return nil, err
	}
	return rsa.DecryptPKCS1v15(rand.Reader, pri, ciphertext)
}

// Using RSA private key file to decrypt ciphertext.
func RSADecrypt4F(prifile string, ciphertext []byte) ([]byte, error) {
	prikey, err := LoadRSAKey(prifile)
	if err != nil {
		return nil, err
	}
	return RSADecrypt(string(prikey), ciphertext)
}

// -------------------------------------------------------------------
// Sign the given string by RSA private key as PKCS#1, ASN.1 format
// and verify by public key.
// -------------------------------------------------------------------

// Using RSA private key to make digital signature,
// the private key in PKCS#1, ASN.1 DER form.
//
//	original := "original-content"
//	prikey, pubkey, _ := secure.NewRSAKeys(1024)
//	signature, _ := secure.RSASign(prikey, original)
//	err := secure.RSAVerify(pubkey, original, signature)
//	// success : err != nil
func RSASign(prikey, original string) ([]byte, error) {
	pri, err := ParsePriKey(prikey)
	if err != nil {
		return nil, err
	}

	hashed := HashSHA256([]byte(original))
	return rsa.SignPKCS1v15(rand.Reader, pri, crypto.SHA256, hashed)
}

// Using RSA private key file to make digital signature,
// the private key in PKCS#1, ASN.1 DER form.
func RSASign4F(prifile string, original string) ([]byte, error) {
	prikey, err := LoadRSAKey(prifile)
	if err != nil {
		return nil, err
	}
	return RSASign(prikey, original)
}

// Using RSA private key to make digital signature,
// then format to base64 form, the private key in PKCS#1, ASN.1 DER form.
//
//	original := "original-content"
//	prikey, pubkey, _ := secure.NewRSAKeys(1024)
//	signatureBase64, _ := secure.RSASignB64(prikey, original)
//	signature, _ := secure.DecodeBase64(signatureBase64)
//	err := secure.RSAVerify(pubkey, original, signature)
//	// success : err != nil
func RSASignB64(prikey, original string) (string, error) {
	buf, err := RSASign(prikey, original)
	if err != nil {
		return "", err
	}
	return ByteToBase64(buf), nil
}

// Using RSA private key file to make digital signature,
// then format to base64 form, the private key in PKCS#1, ASN.1 DER form.
func RSASign4FB64(prifile, original string) (string, error) {
	buf, err := RSASign4F(prifile, original)
	if err != nil {
		return "", err
	}
	return ByteToBase64(buf), nil
}

// ------- PKIX : Verify PKCS#1 v1.5 signature data

// Using RSA public key to verify PKCS#1 v1.5 signatured data.
//
//	original := "original-content"
//	prikey, pubkey, _ := secure.NewRSAKeys(1024)
//	signature, _ := secure.RSASign(prikey, original)
//	err := secure.RSAVerify(pubkey, original, signature)
//	// success : err != nil
func RSAVerify(pubkey, original string, signature []byte) error {
	pub, err := ParsePubKey(pubkey)
	if err != nil {
		return err
	}
	hashed := HashSHA256([]byte(original))
	return rsa.VerifyPKCS1v15(pub, crypto.SHA256, hashed[:], signature)
}

// Using RSA public key file to verify PKCS#1 v1.5 signatured data.
func RSAVerify4F(pubfile, original string, signature []byte) error {
	pubkey, err := LoadRSAKey(pubfile)
	if err != nil {
		return err
	}
	return RSAVerify(pubkey, original, signature)
}

// ------- ASN : Verify ASN.1 signature data

// Using RSA private key to make digital signature,
// the private key in PKCS#1, ASN.1 DER form.
//
//	original := "original-content"
//	prikey, pubkey, _ := secure.NewRSAKeys(1024)
//	signature, _ := secure.RSASignASN(prikey, original)
//	err := secure.RSAVerifyASN(pubkey, original, signature)
//	// success : err != nil
func RSASignASN(prikey, original string) ([]byte, error) {
	pri, err := ParsePriKey(prikey)
	if err != nil {
		return nil, err
	}

	hash, digect := crypto.SHA256, HashSHA256([]byte(original))
	opts := &rsa.PSSOptions{SaltLength: rsa.PSSSaltLengthAuto, Hash: hash}
	signature, err := rsa.SignPSS(rand.Reader, pri, hash, digect, opts)
	if err != nil {
		return nil, err
	}
	return signature, nil
}

// Using RSA public key to verify ASN.1 signatured data.
//
//	original := "original-content"
//	prikey, pubkey, _ := secure.NewRSAKeys(1024)
//	signature, _ := secure.RSASignASN(prikey, original)
//	err := secure.RSAVerifyASN(pubkey, original, signature)
//	// success : err != nil
func RSAVerifyASN(pubkey, original string, signature []byte) error {
	pub, err := ParsePubKey(pubkey)
	if err != nil {
		return err
	}

	hash, digect := crypto.SHA256, HashSHA256([]byte(original))
	opts := &rsa.PSSOptions{SaltLength: rsa.PSSSaltLengthAuto, Hash: hash}
	return rsa.VerifyPSS(pub, hash, digect, signature, opts)
}

// -------------------------------------------------------------------
// Sign the given string by RSA private key as PKCS#8, ASN.1 format
// and verify by public key.
// -------------------------------------------------------------------

// Create RSA private and public keys in PKCS#8, ASN.1 DER format,
// and limit bits length of key cert.
func NewRSA8Keys(bits ...int) (string, string, error) {
	return newRSAKeysByType("PKCS8", bits...)
}

// Using RSA private key to make digital signature,
// the private key in PKCS#8, ASN.1 DER form.
func RSA8Sign(prikey, original string) ([]byte, error) {
	pri, err := ParsePriKey(prikey, true)
	if err != nil {
		return nil, err
	}

	hashed := HashSHA256([]byte(original))
	return rsa.SignPKCS1v15(rand.Reader, pri, crypto.SHA256, hashed)
}

// Using RSA private key file to make digital signature,
// then format to base64 form, the private key in PKCS#8, ASN.1 DER form.
func RSA8SignB64(prikey, original string) (string, error) {
	buf, err := RSA8Sign(prikey, original)
	if err != nil {
		return "", err
	}
	return ByteToBase64(buf), nil
}

// Using RSA private key file to make digital signature.
// the private key in PKCS#8, ASN.1 DER form.
func RSA8Sign4F(prifile, original string) ([]byte, error) {
	prikey, err := LoadRSAKey(prifile)
	if err != nil {
		return nil, err
	}
	return RSA8Sign(prikey, original)
}

// Using RSA private key file to make digital signature,
// then format to base64 form, the private key in PKCS#8, ASN.1 DER form.
func RSA8Sign4FB64(prifile, original string) (string, error) {
	buf, err := RSA8Sign4F(prifile, original)
	if err != nil {
		return "", err
	}
	return ByteToBase64(buf), nil
}

// Using RSA public key to verify PKCS#8, ASN.1 signatured data.
func RSA8Verify(pubkey, original string, signature []byte) error {
	return RSAVerify(pubkey, original, signature)
}

// Using RSA public key to verify PKCS#8, ASN.1 signatured data.
func RSA8Verify4F(pubfile, original string, signature []byte) error {
	return RSAVerify4F(pubfile, original, signature)
}

// -------------------------------------------------------------------
// Create a cert by given PKCS#1 or PKCS#8 RSA private key with target
// organization and expire days. So, call RSA3Sign() and RSA3Verify()
// to sign and verify source datas base on this cert.
// -------------------------------------------------------------------

// Create a serianl number for generate cert pem file data
// by call NewRSACert().
func NewSerialNumber() (*big.Int, error) {
	return rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
}

// Create a cert from the given PKSC#1 or PKCS#8 RSA private key,
// with the target organization and expire days, default using PKCS#1,
// set pkcs8 = true to use PKCS#8 parse RSA private key.
//
//	// USAGE 1: For PKSC#1 RSA Private Key.
//	prikey, _, _ := secure.NewRSAKeys(1024)
//	serialnum, _ := secure.NewSerialNumber()
//	cert, _ := secure.NewRSACert(prikey, serialnum, "Your Organization", 365)
//
//	sign, _ := secure.RSASign(prikey, "original text content")
//	err := secure.RSACertVerify(cert, "original text content", sign)
//	// check err if verify success.
//
//
//	// USAGE 2: For PKSC#8 RSA Private Key.
//	prikey, _, _ := secure.NewRSA8Keys(1024)
//	serialnum, _ := secure.NewSerialNumber()
//	cert, _ := secure.NewRSACert(prikey, serialnum, "Your Organization", 365, true)
//
//	sign, _ := secure.RSA8Sign(prikey, "original text content")
//	err := secure.RSACertVerify(cert, "original text content", sign)
//	// check err if verify success.
func NewRSACert(prikey string, serialnum *big.Int, organization string, days int, pkcs8 ...bool) (string, error) {
	if serialnum == nil || organization == "" || days <= 0 {
		return "", invar.ErrInvalidParams
	}

	pri, err := ParsePriKey(prikey, pkcs8...)
	if err != nil {
		return "", err
	}

	certtmp := x509.Certificate{
		SerialNumber:          serialnum,
		Subject:               pkix.Name{Organization: []string{organization}},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(time.Duration(days) * 24 * time.Hour),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	certbytes, err := x509.CreateCertificate(rand.Reader, &certtmp, &certtmp, &pri.PublicKey, pri)
	if err != nil {
		return "", err
	}

	cert := pem.EncodeToMemory(&pem.Block{
		Type: "CERTIFICATE", Bytes: certbytes,
	})
	return string(cert), nil
}

// Using RSA cert to verify RSA signatured data, call NewRSACert()
// to create PKCS#1 or PKCS#8 cert pem data, then sign source by
// RSASign() for PKCS#1 cert, RSA8Sign() for PKCS#8 cert.
func RSACertVerify(certpem, original string, signature []byte) error {
	block, _ := pem.Decode([]byte(certpem))
	if block == nil {
		return invar.ErrBadPublicKey
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return err
	}

	pub := cert.PublicKey.(*rsa.PublicKey)
	hashed := HashSHA256([]byte(original))
	return rsa.VerifyPKCS1v15(pub, crypto.SHA256, hashed[:], signature)
}

// Load RSA cert pem data from given file and verify the signatured data,
// the cert pem data maybe PKSC#1 or PKSC#8 formated.
func RSACertVerify4F(certfile, original string, signature []byte) error {
	certpem, err := LoadRSAKey(certfile)
	if err != nil {
		return err
	}
	return RSACertVerify(certpem, original, signature)
}

// -------------------------------------------------------------------
// Private methods define.
// -------------------------------------------------------------------

// Create RSA private and public keys by given PKCS# type and cert length,
// the pkcs input param can valued one of 'PKCS1' or 'PKCS8'.
func newRSAKeysByType(pkcs string, bits ...int) (string, string, error) {
	certlen := 1024 // default cert length
	if len(bits) > 0 && bits[0] > 0 {
		certlen = bits[0]
	}

	// generate private key
	prikey, err := rsa.GenerateKey(rand.Reader, certlen)
	if err != nil {
		return "", "", err
	}

	// marshal private key by given type
	var derstream []byte
	switch pkcs {
	case "PKCS1":
		derstream = x509.MarshalPKCS1PrivateKey(prikey)
	case "PKCS8":
		cs8stream, err := x509.MarshalPKCS8PrivateKey(prikey)
		if err != nil {
			return "", "", err
		}
		derstream = cs8stream
	default:
		return "", "", invar.ErrInvalidOptions
	}

	// create buffer to save private pem content
	pribuff := new(bytes.Buffer)
	block := &pem.Block{Type: RSA_PEM_PRI_HEADER, Bytes: derstream}
	if err = pem.Encode(pribuff, block); err != nil {
		return "", "", err
	}

	pubkey := &prikey.PublicKey
	derpkix, err := x509.MarshalPKIXPublicKey(pubkey)
	if err != nil {
		return "", "", err
	}

	pubbuff := new(bytes.Buffer)
	block = &pem.Block{Type: RSA_PEM_PUB_HEADER, Bytes: derpkix}
	if err = pem.Encode(pubbuff, block); err != nil {
		return "", "", err
	}

	return pribuff.String(), pubbuff.String(), nil
}
