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
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/wengoldx/xcore/invar"
)

// Claims jwt claims data
type Claims struct {
	Keyword string `json:"keyword"`
	jwt.RegisteredClaims
}

// Deprecated: use utils.NewJwtToken instead it.
func GenJwtToken(k, s string, d time.Duration) (string, error) { return NewJwtToken(k, s, d) }

// Create a jwt token with keyword and salt string, the token will expired after the given duration.
func NewJwtToken(keyword, salt string, dur time.Duration) (string, error) {
	expireAt := time.Now().Add(dur)
	claims := Claims{
		keyword,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireAt),
			Issuer:    keyword,
		},
	}

	// create the token using your claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// signs the token with a salt.
	signedToken, err := token.SignedString([]byte(salt))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

// Verify the encoded jwt token with salt string
func ViaJwtToken(signedToken, salt string) (string, error) {
	token, err := jwt.ParseWithClaims(signedToken, &Claims{}, func(token *jwt.Token) (any, error) {
		return []byte(salt), nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims.Keyword, nil
	}
	return "", err
}

// Encode account uuid and optianl datas as claims content of jwt token
func EncClaims(uuid string, params ...string) string {
	sets := []string{uuid}
	if len(params) > 0 {
		sets = append(sets, params...)
	}
	orikey := strings.Join(sets, ";")
	return EncodeBase64(orikey)
}

// Decode claims of jwt token and return datas as string array
func DecClaims(keyword string, count ...int) ([]string, error) {
	orikeys, err := DecodeBase64(keyword)
	if err != nil {
		return nil, err
	}

	sets := strings.Split(orikeys, ";")

	// check claims content fields if give the verify count param
	if cl := len(count); cl > 0 && count[0] > 0 && count[0] != len(sets) {
		return nil, invar.ErrInvalidNum
	}
	return sets, nil
}

// Encode account uuid, password and subject string
//
// Deprecated: Use secure.EncClaims() instead it.
func EncJwtKeyword(uuid, pwd string, subject string) string {
	sets := []string{uuid, pwd, subject}
	orikey := strings.Join(sets, ";")
	return EncodeBase64(orikey)
}

// Decode account uuid, password and subject from jwt keyword string
//
// Deprecated: Use secure.DecClaims() instead it.
func DecJwtKeyword(keyword string) (string, string, string) {
	orikeys, err := DecodeBase64(keyword)
	if err != nil {
		return "", "", ""
	}

	sets := strings.Split(orikeys, ";")
	for i := len(sets); i < 3; i++ {
		sets = append(sets, "")
	}
	return sets[0], sets[1], sets[2]
}
