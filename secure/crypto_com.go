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
	"crypto/hmac"
	"crypto/md5"
	crypto "crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/wengoldx/xcore/invar"
	"golang.org/x/crypto/scrypt"
)

const (
	oauthCodeSeedsNum   = "0123456789"
	oauthCodeSeedsLower = "abcdefghijklmnopqrstuvwxyz"
	oauthCodeSeedsUpper = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	oauthCodeSeedsChars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	radixCodeCharLoNum  = "0123456789abcdefghijklmnopqrstuvwxyz"
	radixCodeCharUpNum  = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	radixCodeCharMap    = "01234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	passwordHashBytes   = 64 // default password hash length
)

// For generate uuid string
var uuidNode *snowflake.Node
var rander *rand.Rand

// init uuid generater
func init() {
	rander = rand.New(rand.NewSource(time.Now().UnixNano()))
	if uuidNode == nil {
		node, err := snowflake.NewNode(1)
		if err != nil {
			panic(err)
		}
		uuidNode = node
	}
}

// Generate a code from given chars mapping, params src must over 0, mapping not empty.
func genCodeFromMapping(src int64, mapping string) string {
	radix := (int64)(len(mapping))
	if src <= 0 || radix == 0 {
		return "" // invalid input params
	}

	// encode by given chars mapping
	code := []byte{}
	for v := src; v > 0; v /= radix {
		i := v % radix
		code = append(code, mapping[i])
	}

	// reverse the chars order
	for i, l := 0, len(code); i < l/2; i++ {
		code[i], code[l-i-1] = code[l-i-1], code[i]
	}
	return (string)(code)
}

// Generate a new uuid in int64
func GenUUID() int64 {
	return uuidNode.Generate().Int64()
}

// Generate a new uuid in string
func GenUUIDString() string {
	return uuidNode.Generate().String()
}

// Generate a random number uuid with specified digits
func GenRandUUID(buflen ...int) string {
	length := passwordHashBytes
	if len(buflen) > 0 && buflen[0] > 0 {
		length = buflen[0]
	}

	letters := []rune(oauthCodeSeedsChars)
	buf, letlen := make([]rune, length), len(letters)
	for i := range buf {
		buf[i] = letters[rand.Intn(letlen)]
	}
	return string(buf)
}

// Generate a code just as current nano seconds time, e.g. 1693359476235899600
func GenNano() string {
	return strconv.FormatInt(time.Now().UnixNano(), 10)
}

// Generate a code by using current nanosecond, e.g. M25eNdE4rF5
func GenCode() string {
	return genCodeFromMapping(time.Now().UnixNano(), radixCodeCharMap)
}

// Generate a code from given int64 data, e.g. M25eNdE4rF5
func GenCodeFrom(src int64) string {
	return genCodeFromMapping(src, radixCodeCharMap)
}

// Generate a code formated only lower chars, e.g. mabendecrfdme
func GenLowCode() string {
	return genCodeFromMapping(time.Now().UnixNano(), oauthCodeSeedsLower)
}

// Generate a code formated only upper chars, e.g. MABENDECRFDME
func GenUpCode() string {
	return genCodeFromMapping(time.Now().UnixNano(), oauthCodeSeedsUpper)
}

// Generate a code formated only number and lower chars, e.g. m25ende4rf5m
func GenLowNum() string {
	return genCodeFromMapping(time.Now().UnixNano(), radixCodeCharLoNum)
}

// Generate a code formated only number and upper chars, e.g. M25ENDE4RF5M
func GenUpNum() string {
	return genCodeFromMapping(time.Now().UnixNano(), radixCodeCharUpNum)
}

// Generate a code by using current nanosecond and append random suffix, e.g. M25eNdE4rF50987
func GenRandCode() string {
	rander.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%s%04d", GenCode(), rander.Intn(1000))
}

// Generate a code from given int64 data and append random suffix, e.g. M25eNdE4rF50987
func GenRandCodeFrom(src int64) string {
	rander.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%s%04d", GenCodeFrom(src), rand.Intn(1000))
}

// Convert to lower string and encode by base64 -> md5
func GenToken(original string) string {
	return EncodeB64MD5(strings.ToLower(original))
}

// Generate a random num and convert to string
func GenNonce() string {
	res := make([]byte, 32)
	seeds := [][]int{{10, 48}, {26, 97}, {26, 65}}

	rander.Seed(time.Now().UnixNano())
	for i := 0; i < 32; i++ {
		v := seeds[rand.Intn(3)]
		res[i] = uint8(v[1] + rand.Intn(v[0]))
	}
	return string(res)
}

// Generate a random OAuth code
func GenOAuthCode(length int, randomType string) (string, error) {
	// fill random seeds chars
	buf := bytes.Buffer{}
	if strings.Contains(randomType, "0") {
		buf.WriteString(oauthCodeSeedsNum)
	}
	if strings.Contains(randomType, "a") {
		buf.WriteString(oauthCodeSeedsLower)
	}
	if strings.Contains(randomType, "A") {
		buf.WriteString(oauthCodeSeedsUpper)
	}

	// check random seeds if empty
	str := buf.String()
	len := len(str)
	if len == 0 {
		return "", invar.ErrUnkownCharType
	}

	// random OAuth code
	buf.Reset()
	rander.Seed(time.Now().UnixNano())
	for i := 0; i < length; i++ {
		buf.WriteByte(str[rand.Intn(len)])
	}
	return buf.String(), nil
}

// Generates a random salt, default length is 64 * 2,
// you may set buffer length by buflen input param, and return
// (buflen * 2) length salt string.
func GenSalt(buflen ...int) (string, error) {
	length := passwordHashBytes
	if len(buflen) > 0 && buflen[0] > 0 {
		length = buflen[0]
	}

	buf := make([]byte, length)
	if _, err := io.ReadFull(crypto.Reader, buf); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", buf), nil
}

// Hash the given source with salt, default length is 64 * 2,
// you may set buffer length by buflen input param, and return
// (buflen * 2) length hash string.
func GenHash(src, salt string, buflen ...int) (string, error) {
	length := passwordHashBytes
	if len(buflen) > 0 && buflen[0] > 0 {
		length = buflen[0]
	}

	hex, err := scrypt.Key([]byte(src), []byte(salt), 16384, 8, 1, length)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", hex), nil
}

// Hash string by md5, it ignore write buffer errors
func HashMD5(original []byte) []byte {
	h := md5.New()
	h.Write(original)
	return h.Sum(nil)
}

// Hash string by md5 and check write buffer errors
func HashMD5Check(original []byte) ([]byte, error) {
	h := md5.New()
	if _, err := h.Write(original); err != nil {
		return nil, err
	}
	return h.Sum(nil), nil
}

// Hash byte array by sha256
func HashSHA256(original []byte) []byte {
	// h := sha256.New()
	// h.Write(original)
	// hashed := h.Sum(nil)
	hashed := sha256.Sum256(original)
	return hashed[:]
}

// Hash byte array by sha256 then encode to hex
func HashSHA256Hex(original []byte) string {
	return hex.EncodeToString(HashSHA256(original))
}

// Hash string by sha256
func HashSHA256String(original string) []byte {
	return HashSHA256([]byte(original))
}

// Use HmacSHA1 to calculate the signature,
// and format as base64 string
func SignSHA1(securekey string, src string) string {
	mac := hmac.New(sha1.New, []byte(securekey))
	mac.Write([]byte(src))
	return ByteToBase64(mac.Sum(nil))
}

// Use HmacSHA256 to calculate the signature,
// and format as base64 string
func SignSHA256(securekey string, src string) string {
	mac := hmac.New(sha256.New, []byte(securekey))
	mac.Write([]byte(src))
	return ByteToBase64(mac.Sum(nil))
}

// Decode base64 string to byte array
func Base64ToByte(ciphertext string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(ciphertext)
}

// Encode byte array to base64 string
func ByteToBase64(original []byte) string {
	return base64.StdEncoding.EncodeToString(original)
}

// Decode from base64 string
func DecodeBase64(ciphertext string) (string, error) {
	original, err := Base64ToByte(ciphertext)
	if err != nil {
		return "", err
	}
	return string(original), nil
}

// Encode string by base64
func EncodeBase64(original string) string {
	return ByteToBase64([]byte(original))
}

// Hash string by sha256 and than to base64 string
func HashThenBase64(data string) string {
	return ByteToBase64(HashSHA256String(data))
}

// Hash byte array by sha256 and than to base64 string
func HashByteThenBase64(data []byte) string {
	return ByteToBase64(HashSHA256(data))
}

// Encode string by md5, it ignore write buffer errors
func EncodeMD5(original string) string {
	return hex.EncodeToString(HashMD5([]byte(original)))
}

// Encode string by md5 and check write buffer errors
func EncodeMD5Check(original string) (string, error) {
	cipher, err := HashMD5Check([]byte(original))
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(cipher), nil
}

// Encode string to base64, and then encode by md5
func EncodeB64MD5(original string) string {
	return EncodeMD5(EncodeBase64(original))
}

// Encode string to md5, and then encode by base64
func EncodeMD5B64(original string) string {
	return EncodeBase64(EncodeMD5(original))
}

// Encode multi-input to md5 one string,
// it same as EncodeMD5 when input only one string.
func ToMD5Hex(input ...string) string {
	h := md5.New()
	if len(input) > 0 {
		for _, v := range input {
			io.WriteString(h, v)
		}
	}
	cipher := h.Sum(nil)
	return hex.EncodeToString(cipher)
}

// Encode string to md5 and then transform to uppers.
func ToMD5Upper(original string) (string, error) {
	md5sign, err := EncodeMD5Check(original)
	if err != nil {
		return "", err
	}
	return strings.ToUpper(md5sign), nil
}

// Encode string to md5 and then transform to lowers.
func ToMD5Lower(original string) (string, error) {
	md5sign, err := EncodeMD5Check(original)
	if err != nil {
		return "", err
	}
	return strings.ToLower(md5sign), nil
}

// Encode string to md5 and then transform to uppers without check error.
func MD5Upper(original string) string {
	return strings.ToUpper(EncodeMD5(original))
}

// Encode string to md5 and then transform to lowers without check error.
func MD5Lower(original string) string {
	return strings.ToLower(EncodeMD5(original))
}
