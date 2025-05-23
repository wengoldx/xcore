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
	"crypto/aes"
	"crypto/cipher"
	crand "crypto/rand"
	"io"
	"math/rand"
	"time"

	"github.com/wengoldx/xcore/invar"
)

// For other languages, you can use the follows example code to encrypt or decrypt AES.
//
// `AES for java (Android)`
//
// ----
//
//	public String encryptByAES(String secretkey, String original) {
//	    try {
//	        // use md5 value as the real key
//	        byte[] b = secretkey.getBytes();
//	        MessageDigest md = MessageDigest.getInstance("MD5");
//	        byte[] hashed = md.digest(b);
//
//	        // create an 16-byte initialization vector
//	        byte[] iv = new byte[] {
//	            0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f
//	        };
//	        AlgorithmParameterSpec spec = new IvParameterSpec(iv);
//	        SecretKeySpec keyspec = new SecretKeySpec(hashed), "AES");
//
//	        // create cipher and initialize CBC vector
//	        Cipher ecipher = Cipher.getInstance("AES/CBC/PKCS5Padding");
//	        ecipher.init(Cipher.ENCRYPT_MODE, keyspec, spec);
//
//	        byte[] plaintext = original.getBytes();
//	        byte[] ciphertext = ecipher.doFinal(plaintext, 0, plaintext.length);
//
//	        return Base64.encodeToString(ciphertext, Base64.DEFAULT);
//	    } catch (Exception e) {
//	        e.printStackTrace();
//	    }
//	    return null;
//	}
//
//	public String decryptByAES(String secretkey, String ciphertextb64) {
//	    try {
//	        // use md5 value as the real key
//	        byte[] b = secretkey.getBytes();
//	        MessageDigest md = MessageDigest.getInstance("MD5");
//	        byte[] hashed = md.digest(b);
//
//	        // create an 16-byte initialization vector
//	        byte[] iv = new byte[] {
//	            0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f
//	        };
//	        AlgorithmParameterSpec spec = new IvParameterSpec(iv);
//	        SecretKeySpec keyspec = new SecretKeySpec(hashed), "AES");
//
//	        // create cipher and initialize CBC vector
//	        Cipher dcipher = Cipher.getInstance("AES/CBC/PKCS5Padding");
//	        dcipher.init(Cipher.DECRYPT_MODE, keyspec, spec);
//
//	        byte[] ciphertext = Base64.decode(ciphertextb64, Base64.DEFAULT);
//	        byte[] original = dcipher.doFinal(ciphertext, 0, ciphertext.length);
//
//	        return new String(original);
//	    } catch (Exception e) {
//	        e.printStackTrace();
//	    }
//	    return null;
//	}
//
// `AES for node.js`
//
// ----
//
//	let iv = [ 0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f ];
//
//	function encrypt_by_aes(secretkey, original) {
//	    let md5 = crypto.createHash('md5').update(secretkey).digest('hex');
//	    const ecipher = crypto.createCipheriv(
//	        'aes-128-cbc',
//	        new Buffer(md5, 'hex'),
//	        new Buffer(iv)
//	    );
//	    // ecipher.setAutoPadding(true);
//	    var ciphertextb64 = ecipher.update(original, 'utf8', 'base64');
//	    ciphertextb64 += ecipher.final('base64');
//	    console.log('ciphertextb64: ' + ciphertextb64);
//	    return ciphertextb64;
//	}
//
//	function decrypt_by_aes(secretkey, ciphertextb64) {
//	    let md5 = crypto.createHash('md5').update(secretkey).digest('hex');
//	    const dcipher = crypto.createDecipheriv(
//	        'aes-128-cbc',
//	        new Buffer(md5, 'hex'),
//	        new Buffer(iv)
//	    );
//	    var original = dcipher.update(ciphertextb64, 'base64', 'utf8');
//	    original += dcipher.final('utf8');
//	    console.log('original: ' + original);
//	    return original;
//	}
const AES_UTIL_DESCRIPTION = 0 /* just use for description */

var (
	// aesKeySeeds use to create secret key for ase crypto
	aesKeySeeds = []rune("abcdefghijklmnopqrstuvwxyz1234567890")

	// aesInitVector initialization vector for ase crypto
	aesInitVector = []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f}

	// aesKeyLength key string length of AES
	aesKeyLength = len(aesInitVector)
)

// Create a new AES key range chars in [0-9a-z]{16}
func NewAESKey() string {
	rander.Seed(time.Now().UnixNano())
	secretkey := make([]rune, aesKeyLength)
	sendslen := (int32)(len(aesKeySeeds))
	for i := 0; i < aesKeyLength; i++ {
		j := rand.Int31n(sendslen)
		secretkey[i] = aesKeySeeds[j]
	}
	return string(secretkey)
}

// ------- AES-CBC : Cipher block chaining (CBC) mode.

// Using CBC formated secret key to encrypt original data
//	@param secretkey Secure key buffer
//	@param original Original datas buffer to encrypt
//	@return - string Encrypted ciphertext formated by base64
//			- error Exception message
//
// ----
//
//	secretkey := secure.NewAESKey()
//	original := []byte("original-content")
//	ciphertext, _ := secure.AESEncrypt([]byte(secretkey), original)
//
//	encrypted, _ := secure.AESDecrypt([]byte(secretkey), ciphertext)
//	logger.I("encrypted original string: ", encrypted)
func AESEncrypt(secretkey, original []byte) (string, error) {
	if len(secretkey) != aesKeyLength {
		return "", invar.ErrKeyLenSixteen
	}

	hashed := HashMD5(secretkey)
	block, err := aes.NewCipher(hashed)
	if err != nil {
		return "", err
	}

	enc := cipher.NewCBCEncrypter(block, aesInitVector)
	content := pkcs5Padding(original, block.BlockSize())
	crypted := make([]byte, len(content))
	enc.CryptBlocks(crypted, content)
	return EncodeBase64(string(crypted)), nil
}

// Using CBC formated secret key to decrypt ciphertext
//	@param secretkey Secure key buffer
//	@param ciphertextb64 Ciphertext formated by base64
//	@return - string Decrypted plaintext string
//			- error Exception message
func AESDecrypt(secretkey []byte, ciphertextb64 string) (string, error) {
	if len(secretkey) != aesKeyLength {
		return "", invar.ErrKeyLenSixteen
	}

	hashed := HashMD5(secretkey)
	block, err := aes.NewCipher(hashed)
	if err != nil {
		return "", err
	}

	ciphertext, err := DecodeBase64(ciphertextb64)
	if err != nil {
		return "", err
	}

	dec := cipher.NewCBCDecrypter(block, aesInitVector)
	decrypted := make([]byte, len(ciphertext))
	dec.CryptBlocks(decrypted, []byte(ciphertext))
	return string(pkcs5Unpadding(decrypted)), nil
}

// Use to padding the space of data
func pkcs5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

// Use to unpadding the space of data
func pkcs5Unpadding(encrypt []byte) []byte {
	padding := encrypt[len(encrypt)-1]
	return encrypt[:len(encrypt)-int(padding)]
}

// ------- AES-256-GCM

// Using AES-256-GCM to encrypt the given original text
//	@param secretkey Secure key buffer
//	@param original Original datas buffer to encrypt
//	@param additional Additional datas
//	@return - string Encrypted ciphertext formated by base64
//			- string Nonce string
//			- error Exception message
//
// `NOTICE` :
//
// You can use secure.NewAESKey() to generate a AES-256-GSM secret key
// to as secretkey input param, or use hex.EncodeToString() encode any
// secret string, but use hex.DecodeString() decode the encode hash key
// before call this function.
//
// ----
//
//	// use secure.NewAESKey() generate a secretkey
//	secretkey := secure.NewAESKey()
//	ciphertex, noncestr, err := secure.GCMEncrypt(secretkey, original)
//	ciphertex, noncestr, err := secure.GCMEncrypt(secretkey, original, additional)
//
//	// use hex.EncodeToString() and hex.DecodeString()
//	hashkey := hex.EncodeToString(secretkey)
//	// do samething with hashkey...
//	secretkey, err := hex.DecodeString(hashkey)
//	ciphertex, noncestr, err := secure.GCMEncrypt(secretkey, original)
//	ciphertex, noncestr, err := secure.GCMEncrypt(secretkey, original, additional)
func GCMEncrypt(secretkey, original []byte, additional ...[]byte) (string, string, error) {
	block, err := aes.NewCipher(secretkey)
	if err != nil {
		return "", "", err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", "", err
	}

	nonce := make([]byte, aesgcm.NonceSize())
	if _, err := io.ReadFull(crand.Reader, nonce); err != nil {
		return "", "", err
	}

	var additionalData []byte = nil
	if len(additional) > 0 && len(additional[0]) > 0 {
		additionalData = additional[0]
	}

	ciphertext := aesgcm.Seal(nil, nonce, original, additionalData)
	return ByteToBase64(ciphertext), string(nonce), nil
}

// Using AES-256-GCM to decrypt ciphertext
//	@param secretkey Secure key buffer
//	@param ciphertextb64 Ciphertext formated by base64
//	@param noncestr Nonce string which generated when encrypt
//	@param additional additional datas used by encrypt, it maybe null
//	@return - string Decrypted plaintext string
//			- error Exception message
func GCMDecrypt(secretkey []byte, ciphertextb64, noncestr string, additional ...[]byte) (string, error) {
	block, err := aes.NewCipher(secretkey)
	if err != nil {
		return "", err
	}

	nonce := []byte(noncestr)
	aesgcm, err := cipher.NewGCMWithNonceSize(block, len(nonce))
	if err != nil {
		return "", err
	}

	ciphertext, err := Base64ToByte(ciphertextb64)
	if err != nil {
		return "", err
	}

	var additionalData []byte = nil
	if len(additional) > 0 && len(additional[0]) > 0 {
		additionalData = additional[0]
	}

	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, additionalData)
	if err != nil {
		return "", err
	}
	return string(plaintext), err
}
