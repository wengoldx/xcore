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

// Create a seeds mapping from the given distinct string, The SeedSign
// enable to generate 8 chars random code for verify anywhere without
// cached any ECC, RSA secure certs.
//
// But for security, You can use ECC or RSA certs to encrypt the input
// sign string which output by SignPlaintext().
//
//	+-------+                                       +-----------+
//	| data1 |  -->  secure.SignPlaintext(...)  -->  | plaintext |
//	| data2 |                                       +-----------+
//	| ...   |                                             |
//	+-------+           ECC.EccSign() or RSA.RSASignB64() |
//	                                                      v
//	          +------+                                 +------+
//	     +--   | code |  <--  SeedSign.SignCode()  <--  | sign |
//	     |    +------+                                 +------+
//	     |       |                                         |
//	     |       +--------------+           +--------------+
//	     |                       \_________/
//	like 'vWUABrs2'                    |
//	8 chars string           SeedSign.ViaCode()
//
// The SeedSign better for both side all own the same plaintext datas.
type SeedSign struct {
	// Seeds map for get ecd sign code and verify it,
	// call secure.CreateSeeds() to setup it before use.
	seeds map[int]string

	// Seeds map radix for get seed string in valid range,
	// it inited when call secure.CreateSeeds().
	radix int
}

// The default global singleton, call AsDefault() to init it.
var _defsign *SeedSign

// Create a SeedSign object by distinct src string, it will filter
// out all duplicate chars before create valid seeds mapping.
//
//	See SeedSign.filterDupChars() for filter duplicate chars.
func NewSeedSign(src string) *SeedSign {
	seedsign := &SeedSign{seeds: make(map[int]string)}
	seedsign.createSeeds(src)
	return seedsign
}

// Return the default glocal SeedSign singleton after called
// secure.NewSeedSign().AsDefault() to init it.
func DefSeedSign() *SeedSign { return _defsign }

// Use current SeedSign as the global singleton.
func (s *SeedSign) AsDefault() { _defsign = s }

// Create seeds values and cached as a map, search by index.
//
// # WARNING:
//
// The src param string MUST all unique chars:
//
//	- OK : '1234567890' // good src, all chars uniqued
//	- NG : '1234467890' // have double 4 chars
func (s *SeedSign) createSeeds(src string) {
	src = s.filterDupChars(src)

	// setup seeds maps
	sl, index := len(src), 0
	for i := 0; i < sl; i++ {
		for j := i + 1; j < sl; j++ {
			start, end := src[i:i+1], src[j:j+1]
			per, tail := src[:i]+src[i+1:j], src[j+1:]

			s.seeds[index] = per + end + start + tail
			index++
		}
	}

	// init seeds radix
	s.radix = len(s.seeds)
}

// Filter out all duplicate chars, remain the first one and remove
// others, then return the distinct chars string.
//
//	'1234567890' -> '1234567890' OK
//	'1234467890' -> '123467890'  Filter out duplicate '4'.
func (s *SeedSign) filterDupChars(src string) string {
	chars := []rune{}
	distinct := make(map[rune]struct{})
	for _, char := range src {
		if _, ok := distinct[char]; !ok {
			distinct[char] = struct{}{}
			chars = append(chars, char)
		}
	}
	return string(chars)
}

// Get verify seed by sign string sum number.
func (s *SeedSign) getSignSeed(sign string) string {
	sum := 0
	for _, char := range sign {
		sum += int(char)
	}
	return s.seeds[sum%s.radix]
}

// Return a random code from sign string of ECC or RSA.
//
// # For ECC signature.
//
//	pri, pub, _ := secure.NewEccKeys()
//	plaintext := secure.SignPlaintext(data, data1, ...)
//	sign, _ := secure.EccSign(plaintext, pri)
//	// verify by secure.EccVerify(plaintext, sign, pubkey)
//
// # For RSA signature.
//
//	pri, pub, _ := NewRSAKeys(2048)
//	plaintext := secure.SignPlaintext(data, data1, ...)
//	sign, err := secure.RSASignB64(pri, plaintext) // sign base64 string.
//	// signbytes, _ := secure.Base64ToByte(sign)
//	// verify by secure.EccVerify(plaintext, signbytes, pubkey)
//
// Call secure.ViaSignCode() to verify sign and code whether matched.
func (s *SeedSign) SignCode(sign string) string {
	sl, seed := len(sign), s.getSignSeed(sign)
	if radix := len(seed); sl > radix {
		sl = radix
	}

	code := "" // random 4 group segements
	for seg := 0; seg < 4; seg++ {
		pos := rand.Intn(sl)
		code += seed[pos:pos+1] + sign[pos:pos+1]
	}
	return code // length 8 chars
}

// Verify the code with sign string whether valid.
//
//	// get ecc sign string whereever came from
//	sign := "mhdWY0hJZmBLO4PnxTSWeUd2yqDUUgHyoAbnMjnOZpjo5IVlayfdrkDLsTquvj7nEkpqlSZCKIWx1OhuDq1ZLg=="
//	code := secure.GetSignCode(sign)
//	rst  := secure.ViaSignCode(sign, code) // rst == true
func (s *SeedSign) ViaCode(sign, code string) bool {
	sl, cl := len(sign), len(code)
	if sl == 0 || cl%2 != 0 || sl < (cl/2) {
		return false
	}

	seed := s.getSignSeed(sign)
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

// Use ECC private cert file to encrypt plaintext as sign string.
//
// # WARNING:
//	- Even use the same plaintext and ECC crtfile to sign, it allows
//	output the different sign strings.
func (s *SeedSign) EccSign(plaintext string, crtfile string) (string, error) {
	prikey, err := LoadEccPemFile(crtfile)
	if err != nil {
		return "", err
	}
	return EccSign(plaintext, prikey)
}

// Use ECC public key pem content to verify plaintext and encrypted sign string.
func (s *SeedSign) EccVerify(plaintext, sign, pubpem string) (bool, error) {
	pubkey, err := EccPubKey(pubpem)
	if err != nil {
		return false, err
	}
	return EccVerify(plaintext, sign, pubkey)
}

// Use RSA private key content to encrypt plaintext as sign string (formated as base64 string).
//
// # WARNING:
//	- Use the same plaintext and RSA private key to sign, it allows
//	output the same sign strings.
func (s *SeedSign) RsaSign(plaintext string, prikey string) (string, error) {
	return RSASignB64(prikey, plaintext)
}

// Use RSA public key content to verify plaintext ena encrypted base64 sign string.
func (s *SeedSign) RsaVerify(plaintext, sign, pubkey string) (bool, error) {
	if signbytes, err := Base64ToByte(sign); err != nil {
		return false, err
	} else {
		err = RSAVerify(pubkey, plaintext, signbytes)
		return err == nil, err
	}
}

// Encode signature plaintexts for next ECC or RSA sign and verfiy.
func SignPlaintext(data string, extras ...string) string {
	texts := []string{}
	if data != "" {
		texts = append(texts, data)
	}

	for _, text := range extras {
		if text != "" {
			texts = append(texts, text)
		}
	}
	return strings.Join(texts, "\n")
}
