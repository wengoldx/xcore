// Copyright (c) 2018-Now Dunyu All Rights Reserved.
//
// Author      : https://www.wengold.net
// Email       : support@wengold.net
//
// Prismy.No | Date       | Modified by. | Description
// -------------------------------------------------------------------
// 00001       2019/05/22   yangping       New version
// -------------------------------------------------------------------
package invar

import (
	"regexp"
)

/* The follow gegex strings copied from regex.go of Validator */

const (
	alphaRegexString                 = "^[a-zA-Z]+$"
	alphaNumericRegexString          = "^[a-zA-Z0-9]+$"
	alphaUnicodeRegexString          = "^[\\p{L}]+$"
	alphaUnicodeNumericRegexString   = "^[\\p{L}\\p{N}]+$"
	numericRegexString               = "^[-+]?[0-9]+(?:\\.[0-9]+)?$"
	numberRegexString                = "^[0-9]+$"
	hexadecimalRegexString           = "^(0[xX])?[0-9a-fA-F]+$"
	hexColorRegexString              = "^#(?:[0-9a-fA-F]{3}|[0-9a-fA-F]{4}|[0-9a-fA-F]{6}|[0-9a-fA-F]{8})$"
	rgbRegexString                   = "^rgb\\(\\s*(?:(?:0|[1-9]\\d?|1\\d\\d?|2[0-4]\\d|25[0-5])\\s*,\\s*(?:0|[1-9]\\d?|1\\d\\d?|2[0-4]\\d|25[0-5])\\s*,\\s*(?:0|[1-9]\\d?|1\\d\\d?|2[0-4]\\d|25[0-5])|(?:0|[1-9]\\d?|1\\d\\d?|2[0-4]\\d|25[0-5])%\\s*,\\s*(?:0|[1-9]\\d?|1\\d\\d?|2[0-4]\\d|25[0-5])%\\s*,\\s*(?:0|[1-9]\\d?|1\\d\\d?|2[0-4]\\d|25[0-5])%)\\s*\\)$"
	rgbaRegexString                  = "^rgba\\(\\s*(?:(?:0|[1-9]\\d?|1\\d\\d?|2[0-4]\\d|25[0-5])\\s*,\\s*(?:0|[1-9]\\d?|1\\d\\d?|2[0-4]\\d|25[0-5])\\s*,\\s*(?:0|[1-9]\\d?|1\\d\\d?|2[0-4]\\d|25[0-5])|(?:0|[1-9]\\d?|1\\d\\d?|2[0-4]\\d|25[0-5])%\\s*,\\s*(?:0|[1-9]\\d?|1\\d\\d?|2[0-4]\\d|25[0-5])%\\s*,\\s*(?:0|[1-9]\\d?|1\\d\\d?|2[0-4]\\d|25[0-5])%)\\s*,\\s*(?:(?:0.[1-9]*)|[01])\\s*\\)$"
	hslRegexString                   = "^hsl\\(\\s*(?:0|[1-9]\\d?|[12]\\d\\d|3[0-5]\\d|360)\\s*,\\s*(?:(?:0|[1-9]\\d?|100)%)\\s*,\\s*(?:(?:0|[1-9]\\d?|100)%)\\s*\\)$"
	hslaRegexString                  = "^hsla\\(\\s*(?:0|[1-9]\\d?|[12]\\d\\d|3[0-5]\\d|360)\\s*,\\s*(?:(?:0|[1-9]\\d?|100)%)\\s*,\\s*(?:(?:0|[1-9]\\d?|100)%)\\s*,\\s*(?:(?:0.[1-9]*)|[01])\\s*\\)$"
	emailRegexString                 = "^(?:(?:(?:(?:[a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+(?:\\.([a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+)*)|(?:(?:\\x22)(?:(?:(?:(?:\\x20|\\x09)*(?:\\x0d\\x0a))?(?:\\x20|\\x09)+)?(?:(?:[\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x7f]|\\x21|[\\x23-\\x5b]|[\\x5d-\\x7e]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(?:(?:[\\x01-\\x09\\x0b\\x0c\\x0d-\\x7f]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}]))))*(?:(?:(?:\\x20|\\x09)*(?:\\x0d\\x0a))?(\\x20|\\x09)+)?(?:\\x22))))@(?:(?:(?:[a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(?:(?:[a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])(?:[a-zA-Z]|\\d|-|\\.|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*(?:[a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.)+(?:(?:[a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(?:(?:[a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])(?:[a-zA-Z]|\\d|-|\\.|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*(?:[a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.?$"
	e164RegexString                  = "^\\+[1-9]?[0-9]{7,14}$"
	base64RegexString                = "^(?:[A-Za-z0-9+\\/]{4})*(?:[A-Za-z0-9+\\/]{2}==|[A-Za-z0-9+\\/]{3}=|[A-Za-z0-9+\\/]{4})$"
	base64URLRegexString             = "^(?:[A-Za-z0-9-_]{4})*(?:[A-Za-z0-9-_]{2}==|[A-Za-z0-9-_]{3}=|[A-Za-z0-9-_]{4})$"
	iSBN10RegexString                = "^(?:[0-9]{9}X|[0-9]{10})$"
	iSBN13RegexString                = "^(?:(?:97(?:8|9))[0-9]{10})$"
	uUID3RegexString                 = "^[0-9a-f]{8}-[0-9a-f]{4}-3[0-9a-f]{3}-[0-9a-f]{4}-[0-9a-f]{12}$"
	uUID4RegexString                 = "^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$"
	uUID5RegexString                 = "^[0-9a-f]{8}-[0-9a-f]{4}-5[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$"
	uUIDRegexString                  = "^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"
	uUID3RFC4122RegexString          = "^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-3[0-9a-fA-F]{3}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$"
	uUID4RFC4122RegexString          = "^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-4[0-9a-fA-F]{3}-[89abAB][0-9a-fA-F]{3}-[0-9a-fA-F]{12}$"
	uUID5RFC4122RegexString          = "^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-5[0-9a-fA-F]{3}-[89abAB][0-9a-fA-F]{3}-[0-9a-fA-F]{12}$"
	uUIDRFC4122RegexString           = "^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$"
	uLIDRegexString                  = "^[A-HJKMNP-TV-Z0-9]{26}$"
	md4RegexString                   = "^[0-9a-f]{32}$"
	md5RegexString                   = "^[0-9a-f]{32}$"
	sha256RegexString                = "^[0-9a-f]{64}$"
	sha384RegexString                = "^[0-9a-f]{96}$"
	sha512RegexString                = "^[0-9a-f]{128}$"
	ripemd128RegexString             = "^[0-9a-f]{32}$"
	ripemd160RegexString             = "^[0-9a-f]{40}$"
	tiger128RegexString              = "^[0-9a-f]{32}$"
	tiger160RegexString              = "^[0-9a-f]{40}$"
	tiger192RegexString              = "^[0-9a-f]{48}$"
	aSCIIRegexString                 = "^[\x00-\x7F]*$"
	printableASCIIRegexString        = "^[\x20-\x7E]*$"
	multibyteRegexString             = "[^\x00-\x7F]"
	dataURIRegexString               = `^data:((?:\w+\/(?:([^;]|;[^;]).)+)?)`
	latitudeRegexString              = "^[-+]?([1-8]?\\d(\\.\\d+)?|90(\\.0+)?)$"
	longitudeRegexString             = "^[-+]?(180(\\.0+)?|((1[0-7]\\d)|([1-9]?\\d))(\\.\\d+)?)$"
	sSNRegexString                   = `^[0-9]{3}[ -]?(0[1-9]|[1-9][0-9])[ -]?([1-9][0-9]{3}|[0-9][1-9][0-9]{2}|[0-9]{2}[1-9][0-9]|[0-9]{3}[1-9])$`
	hostnameRegexStringRFC952        = `^[a-zA-Z]([a-zA-Z0-9\-]+[\.]?)*[a-zA-Z0-9]$`                                                                   // https://tools.ietf.org/html/rfc952
	hostnameRegexStringRFC1123       = `^([a-zA-Z0-9]{1}[a-zA-Z0-9-]{0,62}){1}(\.[a-zA-Z0-9]{1}[a-zA-Z0-9-]{0,62})*?$`                                 // accepts hostname starting with a digit https://tools.ietf.org/html/rfc1123
	fqdnRegexStringRFC1123           = `^([a-zA-Z0-9]{1}[a-zA-Z0-9-]{0,62})(\.[a-zA-Z0-9]{1}[a-zA-Z0-9-]{0,62})*?(\.[a-zA-Z]{1}[a-zA-Z0-9]{0,62})\.?$` // same as hostnameRegexStringRFC1123 but must contain a non numerical TLD (possibly ending with '.')
	btcAddressRegexString            = `^[13][a-km-zA-HJ-NP-Z1-9]{25,34}$`                                                                             // bitcoin address
	btcAddressUpperRegexStringBech32 = `^BC1[02-9AC-HJ-NP-Z]{7,76}$`                                                                                   // bitcoin bech32 address https://en.bitcoin.it/wiki/Bech32
	btcAddressLowerRegexStringBech32 = `^bc1[02-9ac-hj-np-z]{7,76}$`                                                                                   // bitcoin bech32 address https://en.bitcoin.it/wiki/Bech32
	ethAddressRegexString            = `^0x[0-9a-fA-F]{40}$`
	ethAddressUpperRegexString       = `^0x[0-9A-F]{40}$`
	ethAddressLowerRegexString       = `^0x[0-9a-f]{40}$`
	uRLEncodedRegexString            = `^(?:[^%]|%[0-9A-Fa-f]{2})*$`
	hTMLEncodedRegexString           = `&#[x]?([0-9a-fA-F]{2})|(&gt)|(&lt)|(&quot)|(&amp)+[;]?`
	hTMLRegexString                  = `<[/]?([a-zA-Z]+).*?>`
	jWTRegexString                   = "^[A-Za-z0-9-_]+\\.[A-Za-z0-9-_]+\\.[A-Za-z0-9-_]*$"
	splitParamsRegexString           = `'[^']*'|\S+`
	bicRegexString                   = `^[A-Za-z]{6}[A-Za-z0-9]{2}([A-Za-z0-9]{3})?$`
	semverRegexString                = `^(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)(?:-((?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+([0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?$` // numbered capture groups https://semver.org/
	dnsRegexStringRFC1035Label       = "^[a-z]([-a-z0-9]*[a-z0-9]){0,62}$"
)

var (
	alphaRegex                 = regexp.MustCompile(alphaRegexString)
	alphaNumericRegex          = regexp.MustCompile(alphaNumericRegexString)
	alphaUnicodeRegex          = regexp.MustCompile(alphaUnicodeRegexString)
	alphaUnicodeNumericRegex   = regexp.MustCompile(alphaUnicodeNumericRegexString)
	numericRegex               = regexp.MustCompile(numericRegexString)
	numberRegex                = regexp.MustCompile(numberRegexString)
	hexadecimalRegex           = regexp.MustCompile(hexadecimalRegexString)
	hexColorRegex              = regexp.MustCompile(hexColorRegexString)
	rgbRegex                   = regexp.MustCompile(rgbRegexString)
	rgbaRegex                  = regexp.MustCompile(rgbaRegexString)
	hslRegex                   = regexp.MustCompile(hslRegexString)
	hslaRegex                  = regexp.MustCompile(hslaRegexString)
	e164Regex                  = regexp.MustCompile(e164RegexString)
	emailRegex                 = regexp.MustCompile(emailRegexString)
	base64Regex                = regexp.MustCompile(base64RegexString)
	base64URLRegex             = regexp.MustCompile(base64URLRegexString)
	iSBN10Regex                = regexp.MustCompile(iSBN10RegexString)
	iSBN13Regex                = regexp.MustCompile(iSBN13RegexString)
	uUID3Regex                 = regexp.MustCompile(uUID3RegexString)
	uUID4Regex                 = regexp.MustCompile(uUID4RegexString)
	uUID5Regex                 = regexp.MustCompile(uUID5RegexString)
	uUIDRegex                  = regexp.MustCompile(uUIDRegexString)
	uUID3RFC4122Regex          = regexp.MustCompile(uUID3RFC4122RegexString)
	uUID4RFC4122Regex          = regexp.MustCompile(uUID4RFC4122RegexString)
	uUID5RFC4122Regex          = regexp.MustCompile(uUID5RFC4122RegexString)
	uUIDRFC4122Regex           = regexp.MustCompile(uUIDRFC4122RegexString)
	uLIDRegex                  = regexp.MustCompile(uLIDRegexString)
	md4Regex                   = regexp.MustCompile(md4RegexString)
	md5Regex                   = regexp.MustCompile(md5RegexString)
	sha256Regex                = regexp.MustCompile(sha256RegexString)
	sha384Regex                = regexp.MustCompile(sha384RegexString)
	sha512Regex                = regexp.MustCompile(sha512RegexString)
	ripemd128Regex             = regexp.MustCompile(ripemd128RegexString)
	ripemd160Regex             = regexp.MustCompile(ripemd160RegexString)
	tiger128Regex              = regexp.MustCompile(tiger128RegexString)
	tiger160Regex              = regexp.MustCompile(tiger160RegexString)
	tiger192Regex              = regexp.MustCompile(tiger192RegexString)
	aSCIIRegex                 = regexp.MustCompile(aSCIIRegexString)
	printableASCIIRegex        = regexp.MustCompile(printableASCIIRegexString)
	multibyteRegex             = regexp.MustCompile(multibyteRegexString)
	dataURIRegex               = regexp.MustCompile(dataURIRegexString)
	latitudeRegex              = regexp.MustCompile(latitudeRegexString)
	longitudeRegex             = regexp.MustCompile(longitudeRegexString)
	sSNRegex                   = regexp.MustCompile(sSNRegexString)
	hostnameRegexRFC952        = regexp.MustCompile(hostnameRegexStringRFC952)
	hostnameRegexRFC1123       = regexp.MustCompile(hostnameRegexStringRFC1123)
	fqdnRegexRFC1123           = regexp.MustCompile(fqdnRegexStringRFC1123)
	btcAddressRegex            = regexp.MustCompile(btcAddressRegexString)
	btcUpperAddressRegexBech32 = regexp.MustCompile(btcAddressUpperRegexStringBech32)
	btcLowerAddressRegexBech32 = regexp.MustCompile(btcAddressLowerRegexStringBech32)
	ethAddressRegex            = regexp.MustCompile(ethAddressRegexString)
	ethAddressRegexUpper       = regexp.MustCompile(ethAddressUpperRegexString)
	ethAddressRegexLower       = regexp.MustCompile(ethAddressLowerRegexString)
	uRLEncodedRegex            = regexp.MustCompile(uRLEncodedRegexString)
	hTMLEncodedRegex           = regexp.MustCompile(hTMLEncodedRegexString)
	hTMLRegex                  = regexp.MustCompile(hTMLRegexString)
	jWTRegex                   = regexp.MustCompile(jWTRegexString)
	splitParamsRegex           = regexp.MustCompile(splitParamsRegexString)
	bicRegex                   = regexp.MustCompile(bicRegexString)
	semverRegex                = regexp.MustCompile(semverRegexString)
	dnsRegexRFC1035Label       = regexp.MustCompile(dnsRegexStringRFC1035Label)
)

func RegexAlpha(src string) bool                 { return alphaRegex.MatchString(src) }
func RegexAlphaNumeric(src string) bool          { return alphaNumericRegex.MatchString(src) }
func RegexAlphaUnicode(src string) bool          { return alphaUnicodeRegex.MatchString(src) }
func RegexAlphaUnicodeNumeric(src string) bool   { return alphaUnicodeNumericRegex.MatchString(src) }
func RegexNumeric(src string) bool               { return numericRegex.MatchString(src) }
func RegexNumber(src string) bool                { return numberRegex.MatchString(src) }
func RegexHexadecimal(src string) bool           { return hexadecimalRegex.MatchString(src) }
func RegexHexColor(src string) bool              { return hexColorRegex.MatchString(src) }
func RegexRgb(src string) bool                   { return rgbRegex.MatchString(src) }
func RegexRgba(src string) bool                  { return rgbaRegex.MatchString(src) }
func RegexHsl(src string) bool                   { return hslRegex.MatchString(src) }
func RegexHsla(src string) bool                  { return hslaRegex.MatchString(src) }
func RegexE164(src string) bool                  { return e164Regex.MatchString(src) }
func RegexEmail(src string) bool                 { return emailRegex.MatchString(src) }
func RegexBase64(src string) bool                { return base64Regex.MatchString(src) }
func RegexBase64URL(src string) bool             { return base64URLRegex.MatchString(src) }
func RegexISBN10(src string) bool                { return iSBN10Regex.MatchString(src) }
func RegexISBN13(src string) bool                { return iSBN13Regex.MatchString(src) }
func RegexUUID3(src string) bool                 { return uUID3Regex.MatchString(src) }
func RegexUUID4(src string) bool                 { return uUID4Regex.MatchString(src) }
func RegexUUID5(src string) bool                 { return uUID5Regex.MatchString(src) }
func RegexUUID(src string) bool                  { return uUIDRegex.MatchString(src) }
func RegexUUID3RFC4122(src string) bool          { return uUID3RFC4122Regex.MatchString(src) }
func RegexUUID4RFC4122(src string) bool          { return uUID4RFC4122Regex.MatchString(src) }
func RegexUUID5RFC4122(src string) bool          { return uUID5RFC4122Regex.MatchString(src) }
func RegexUUIDRFC4122(src string) bool           { return uUIDRFC4122Regex.MatchString(src) }
func RegexULID(src string) bool                  { return uLIDRegex.MatchString(src) }
func RegexMd4(src string) bool                   { return md4Regex.MatchString(src) }
func RegexMd5(src string) bool                   { return md5Regex.MatchString(src) }
func RegexSha256(src string) bool                { return sha256Regex.MatchString(src) }
func RegexSha384(src string) bool                { return sha384Regex.MatchString(src) }
func RegexSha512(src string) bool                { return sha512Regex.MatchString(src) }
func RegexRipemd128(src string) bool             { return ripemd128Regex.MatchString(src) }
func RegexRipemd160(src string) bool             { return ripemd160Regex.MatchString(src) }
func RegexTiger128(src string) bool              { return tiger128Regex.MatchString(src) }
func RegexTiger160(src string) bool              { return tiger160Regex.MatchString(src) }
func RegexTiger192(src string) bool              { return tiger192Regex.MatchString(src) }
func RegexASCII(src string) bool                 { return aSCIIRegex.MatchString(src) }
func RegexPrintableASCII(src string) bool        { return printableASCIIRegex.MatchString(src) }
func RegexMultibyte(src string) bool             { return multibyteRegex.MatchString(src) }
func RegexDataURI(src string) bool               { return dataURIRegex.MatchString(src) }
func RegexLatitude(src string) bool              { return latitudeRegex.MatchString(src) }
func RegexLongitude(src string) bool             { return longitudeRegex.MatchString(src) }
func RegexSSN(src string) bool                   { return sSNRegex.MatchString(src) }
func RegexHostRFC952(src string) bool            { return hostnameRegexRFC952.MatchString(src) }
func RegexHostRFC1123(src string) bool           { return hostnameRegexRFC1123.MatchString(src) }
func RegexFqdnRFC1123(src string) bool           { return fqdnRegexRFC1123.MatchString(src) }
func RegexBtcAddress(src string) bool            { return btcAddressRegex.MatchString(src) }
func RegexBtcUpperAddressBech32(src string) bool { return btcUpperAddressRegexBech32.MatchString(src) }
func RegexBtcLowerAddressBech32(src string) bool { return btcLowerAddressRegexBech32.MatchString(src) }
func RegexEthAddress(src string) bool            { return ethAddressRegex.MatchString(src) }
func RegexEthAddressUpper(src string) bool       { return ethAddressRegexUpper.MatchString(src) }
func RegexEthAddressLower(src string) bool       { return ethAddressRegexLower.MatchString(src) }
func RegexURLEncoded(src string) bool            { return uRLEncodedRegex.MatchString(src) }
func RegexHTMLEncoded(src string) bool           { return hTMLEncodedRegex.MatchString(src) }
func RegexHTML(src string) bool                  { return hTMLRegex.MatchString(src) }
func RegexJWT(src string) bool                   { return jWTRegex.MatchString(src) }
func RegexSplitParams(src string) bool           { return splitParamsRegex.MatchString(src) }
func RegexBic(src string) bool                   { return bicRegex.MatchString(src) }
func RegexSemver(src string) bool                { return semverRegex.MatchString(src) }
func RegexDnsRFC1035Label(src string) bool       { return dnsRegexRFC1035Label.MatchString(src) }
