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
	"math/big"
	"math/rand"
	"strings"

	"github.com/wengoldx/xcore/invar"
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
//	like 'o0r522w4'                    |
//	8 chars string           SeedSign.ViaCode()
//
// The SeedSign better for both side all own the same plaintext datas.
type SeedSign struct {
	// Seeds map for generate sign code and verify it,
	// call createSeeds() to setup it before use.
	seeds map[int]string

	// Seeds map radix for get seed string in valid range,
	// it inited when call createSeeds().
	radix int
}

// The default global singleton, call AsDefault() to init it.
var _defsign = &SeedSign{}

// Create a SeedSign object by distinct chars string, it will filter
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

// # FIXME: The follow utils enable called directly to sign plaintext whatever SeedSign seeds inited.
// Utils Methods Start >>

func SeedESign(crt string, ts ...string) (string, error)    { return _defsign.EccSign(crt, ts...) }
func SeedEVerify(s, pub string, ts ...string) (bool, error) { return _defsign.EccVerify(s, pub, ts...) }
func SeedRSign(pri string, ts ...string) (string, error)    { return _defsign.RsaSign(pri, ts...) }
func SeedRVerify(s, pub string, ts ...string) (bool, error) { return _defsign.RsaVerify(s, pub, ts...) }

// Utils Methods End <<

// Check the SeedSign object whether inited.
func (s *SeedSign) isPrepared() bool {
	return s.radix > 0
}

// Create seeds values and cached as a map, search by index.
//
// # WARNING:
//
// The src string MUST all distinct chars and over 2 chars lenght:
//
//	- OK : '1234567890' // good src, all chars distinct.
//	- NG : '1234467890' // have double 4 chars
//	- NG : '1'          // too short src string.
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

	// init seeds radix.
	s.radix = len(s.seeds)
}

// Filter out all duplicate chars, remain the first one and remove
// others, then return the distinct chars string.
//
//	'1234567890' -> '1234567890' OK
//	'1234467890' -> '123467890'  Distinct the duplicate '4'.
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

// Get verify seed by sign encrypted string sum number.
func (s *SeedSign) signSeed(ensign string) string {
	sum := 0
	for _, char := range ensign {
		sum += int(char)
	}
	return s.seeds[sum%s.radix]
}

// Convert sign string to number string.
//
//	input sign: 'ghdWBIEJFuiKgKtL89dfNBfNX7hXKAQj85hP40UcbgC+rPIujfCcac1w6fz/wcdzr1dTAvR2zXfn1yegPnsYDCA='
//	md5 sign  : '124833b4bc9944cd38cf26424e447a5f'
//	sign num  : '1862935100180924569857379440623587583846348235555'
func (s *SeedSign) signSeedNum(sign string) (string, string) {
	sign = MD5Lower(sign)             // encode sign to md5 lower string.
	weight := len(radixCodeCharLoNum) // Fixed using 36 radix!
	if num, ok := new(big.Int).SetString(sign, weight); ok {
		return s.signSeed(sign), num.String()
	}
	return "", ""
}

// Return a random code from sign string of ECC or RSA.
//
// 1. Generate ECC signature.
//
//	pri, pub, _ := secure.NewEccKeys()
//	plaintext := secure.SignPlaintext(data, data1, ...)
//	sign, _ := secure.EccSign(plaintext, pri)
//	// verify by secure.EccVerify(plaintext, sign, pubkey)
//	// Or, call secure.SeedESign() and secure.SeedEVerify().
//
// 2. Generate RSA signature.
//
//	pri, pub, _ := NewRSAKeys(2048)
//	plaintext := secure.SignPlaintext(data, data1, ...)
//	sign, err := secure.RSASignB64(pri, plaintext) // sign base64 string.
//	// signbytes, _ := secure.Base64ToByte(sign)
//	// verify by secure.RSAVerify(pubkey, plaintext, signbytes)
//	// Or, call secure.SeedRSign() and secure.SeedRVerify().
//
// Then call secure.DefSeedSign().ViaCode() to verify sign and code whether matched.
//
// # WARING:
//	- This method need call secure.NewSeedSign() to init seeds first!
func (s *SeedSign) SignCode(sign string) string {
	sign = strings.TrimSpace(sign)
	if !s.isPrepared() || sign == "" {
		return ""
	}

	seed, num := s.signSeedNum(sign)
	weight := len(seed)              // set seed weight as default.
	if nw := len(num); weight > nw { // check sign number lenght.
		weight = nw // use the minimum weight.
	}

	code := "" // random 4 group segements
	for seg := 0; seg < 4; seg++ {
		pos := rand.Intn(weight)
		code += seed[pos:pos+1] + num[pos:pos+1]
	}
	return code // length 8 chars
}

// Verify the code with sign string whether valid.
//
//	ss := secure.DefSeedSign()
//	// get sign string whereever came from ecc or rsa.
//	sign := "mhdWY0hJZmBLO4PnxTSWeUd2yqDUUgHyoAbnMjnOZpjo5IVlayfdrkDLsTquvj7nEkpqlSZCKIWx1OhuDq1ZLg=="
//	code := ss.SignCode(sign)
//	rst  := ss.ViaCode(sign, code) // rst == true
//
// # WARING:
//	- This method need call secure.NewSeedSign() to init seeds first!
func (s *SeedSign) ViaCode(sign, code string) bool {
	sl, cl := len(sign), len(code)
	invalids := (sl == 0 || cl%2 != 0 || sl < (cl/2))
	if !s.isPrepared() || invalids {
		return false
	}

	seed, num := s.signSeedNum(sign)
	for i, weight := 0, len(seed); i < cl-1; i += 2 {
		poschar, digital := code[i], code[i+1]
		pos := strings.Index(seed, string(poschar))
		if pos < 0 || pos >= weight {
			return false
		}

		if num[pos] != digital {
			return false
		}
	}
	return true
}

// Encode signature plaintexts as line by line.
//
//	ss := secure.DefSeedSign()
//	plaintext := ss.SignPlaintext("test_text1", "text2", "t3")
//	// test_text1
//	// text2
//	// t3
func (s *SeedSign) SignPlaintext(texts ...string) string {
	valids := []string{}
	for _, t := range texts {
		if t = strings.TrimSpace(t); t != "" {
			valids = append(valids, t)
		}
	}
	return strings.Join(valids, "\n")
}

// Use ECC private cert file to encrypt plaintext as sign string.
//
// # WARNING:
//	- Even use the same plaintext and ECC crtfile to sign, it always
//	output the different sign strings.
func (s *SeedSign) EccSign(crtfile string, texts ...string) (string, error) {
	if plaintext := s.SignPlaintext(texts...); plaintext != "" {
		prikey, err := LoadEccPemFile(crtfile)
		if err != nil {
			return "", err
		}
		return EccSign(plaintext, prikey)
	}
	return "", invar.ErrEmptyData
}

// Use ECC public key pem content to verify plaintext and encrypted sign string.
func (s *SeedSign) EccVerify(sign, pubpem string, texts ...string) (bool, error) {
	if plaintext := s.SignPlaintext(texts...); plaintext != "" {
		pubkey, err := EccPubKey(pubpem)
		if err != nil {
			return false, err
		}
		return EccVerify(plaintext, sign, pubkey)
	}
	return false, invar.ErrEmptyData
}

// Use RSA private key content to encrypt plaintext as sign string (formated as base64 string).
//
// # WARNING:
//	- Use the same plaintext and RSA private key to sign, it always
//	output the same sign strings.
func (s *SeedSign) RsaSign(prikey string, texts ...string) (string, error) {
	if plaintext := s.SignPlaintext(texts...); plaintext != "" {
		return RSASignB64(prikey, plaintext)
	}
	return "", invar.ErrEmptyData
}

// Use RSA public key content to verify plaintext ena encrypted base64 sign string.
func (s *SeedSign) RsaVerify(sign, pubkey string, texts ...string) (bool, error) {
	if plaintext := s.SignPlaintext(texts...); plaintext != "" {
		if signbytes, err := Base64ToByte(sign); err != nil {
			return false, err
		} else {
			err = RSAVerify(pubkey, plaintext, signbytes)
			return err == nil, err
		}
	}
	return false, invar.ErrEmptyData
}
