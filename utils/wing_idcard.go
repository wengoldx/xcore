// Copyright (c) 2018-Now Dunyu All Rights Reserved.
//
// Author      : https://www.wengold.net
// Email       : support@wengold.net
//
// Prismy.No | Date       | Modified by. | Description
// -------------------------------------------------------------------
// 00001       2019/05/22   yangping       New version
// -------------------------------------------------------------------

package utils

import (
	"strconv"
	"strings"
	"time"

	"github.com/wengoldx/xcore/invar"
	"github.com/wengoldx/xcore/logger"
)

/*
 * ID Card Descriptions :
 *
 * card string length fixed at 18 chars, the first 17 chars must be number.
 * e.g. [42 11 26 19810328 93 5 2]
 *
 * (1).  1 ~  2 digits represent the code of province.
 * (2).  3 ~  4 digits represent the code of city.
 * (3).  5 ~  6 digits represent the code of district.
 * (4).  7 ~ 14 digits represent year, month, and day of birth.
 * (5). 15 ~ 16 digits represent the code of the local police station.
 * (6).      17 diget represent gender, odd numbers for male, other for female.
 * (7).      18 diget is the verification code, must be 0 ~9 or X char.
 */

const validateCodes = "10X98765432"

var validateWeights = []int{7, 9, 10, 5, 8, 4, 2, 1, 6, 3, 7, 9, 10, 5, 8, 4, 2}

// Verify ID Card internal, just simple validate card number only
func IsVaildIDCard(card string) bool {
	card = strings.ToUpper(card)
	if cardlen := len(card); cardlen != 18 {
		logger.E("Invalid ID Card:", card, "lenght:", cardlen)
		return false
	}

	num, last := card[:17], card[17:]
	if !invar.RegexNumber(num) {
		logger.E("Not digites of ID Card:", num)
		return false
	}
	return validateCardNumbers(num, last)
}

// Return birthday as time from given ID Card string
func CardBirthday(card string) (*time.Time, error) {
	if len(card) != 18 {
		return nil, invar.ErrInvalidParams
	}

	birthday := card[6:14]
	bt, err := ParseTime(DateNoneHyphen, birthday)
	if err != nil {
		logger.E("Parse card birthday:", birthday, "err:", err)
		return nil, err
	}
	return &bt, nil
}

// Return gender from given ID Card string, true is male, false is female
func CardGender(card string) (bool, error) {
	if len(card) != 18 {
		return false, invar.ErrInvalidParams
	}

	genderMask, err := strconv.Atoi(string(card[16]))
	if err != nil {
		logger.E("Parse card:", card, "gender, err:", err)
		return false, err
	}
	return genderMask%2 == 1 /* male: 1, female: 2 */, nil
}

// Validate card number self if valide by last code char
// see more http://www.360doc.com/content/22/0112/12/74433059_1012930821.shtml
func validateCardNumbers(num, last string) bool {
	sum := 0
	for i := 0; i < 17; i++ {
		if digit, err := strconv.Atoi(num[i : i+1]); err != nil {
			logger.E("Invalid digit number at:", i, "err:", err)
			return false
		} else {
			sum += digit * validateWeights[i]
		}
	}

	index := sum % 11
	code := validateCodes[index : index+1]
	return code == last
}

/* ========== Foreign-Relate passport authentication ========== */

// All nations informations on ISO 3166-1
//	@See more visit https://baike.baidu.com/item/ISO%203166-1/5269555?fr=ge_ala
var Nations = map[string]byte{
	"AFG": 0, "ALA": 0, "ALB": 0, "DZA": 0, "ASM": 0, "AND": 0, "AGO": 0, "AIA": 0, "ATA": 0, "ATG": 0, "ARG": 0,
	"ARM": 0, "ABW": 0, "AUS": 0, "AUT": 0, "AZE": 0, "BHS": 0, "BHR": 0, "BGD": 0, "BRB": 0, "BLR": 0, "BEL": 0,
	"BLZ": 0, "BEN": 0, "BMU": 0, "BTN": 0, "BOL": 0, "BIH": 0, "BWA": 0, "BVT": 0, "BRA": 0, "IOT": 0, "BRN": 0,
	"BGR": 0, "BFA": 0, "BDI": 0, "KHM": 0, "CMR": 0, "CAN": 0, "CPV": 0, "CYM": 0, "CAF": 0, "TCD": 0, "CHL": 0,
	"CHN": 0, "CXR": 0, "CCK": 0, "COL": 0, "COM": 0, "COG": 0, "COD": 0, "COK": 0, "CRI": 0, "CIV": 0, "HRV": 0,
	"CUB": 0, "CYP": 0, "CZE": 0, "DNK": 0, "DJI": 0, "DMA": 0, "DOM": 0, "ECU": 0, "EGY": 0, "SLV": 0, "GNQ": 0,
	"ERI": 0, "EST": 0, "ETH": 0, "FLK": 0, "FRO": 0, "FJI": 0, "FIN": 0, "FRA": 0, "GUF": 0, "PYF": 0, "ATF": 0,
	"GAB": 0, "GMB": 0, "GEO": 0, "DEU": 0, "GHA": 0, "GIB": 0, "GRC": 0, "GRL": 0, "GRD": 0, "GLP": 0, "GUM": 0,
	"GTM": 0, "GGY": 0, "GIN": 0, "GNB": 0, "GUY": 0, "HTI": 0, "HMD": 0, "VAT": 0, "HND": 0, "HKG": 0, "HUN": 0,
	"ISL": 0, "IND": 0, "IDN": 0, "IRN": 0, "IRQ": 0, "IRL": 0, "IMN": 0, "ISR": 0, "ITA": 0, "JAM": 0, "JPN": 0,
	"JEY": 0, "JOR": 0, "KAZ": 0, "KEN": 0, "KIR": 0, "PRK": 0, "KOR": 0, "KWT": 0, "KGZ": 0, "LAO": 0, "LVA": 0,
	"LBN": 0, "LSO": 0, "LBR": 0, "LBY": 0, "LIE": 0, "LTU": 0, "LUX": 0, "MAC": 0, "MKD": 0, "MDG": 0, "MWI": 0,
	"MYS": 0, "MDV": 0, "MLI": 0, "MLT": 0, "MHL": 0, "MTQ": 0, "MRT": 0, "MUS": 0, "MYT": 0, "MEX": 0, "FSM": 0,
	"MDA": 0, "MCO": 0, "MNG": 0, "MNE": 0, "MSR": 0, "MAR": 0, "MOZ": 0, "MMR": 0, "NAM": 0, "NRU": 0, "NPL": 0,
	"NLD": 0, "ANT": 0, "NCL": 0, "NZL": 0, "NIC": 0, "NER": 0, "NGA": 0, "NIU": 0, "NFK": 0, "MNP": 0, "NOR": 0,
	"OMN": 0, "PAK": 0, "PLW": 0, "PSE": 0, "PAN": 0, "PNG": 0, "PRY": 0, "PER": 0, "PHL": 0, "PCN": 0, "POL": 0,
	"PRT": 0, "PRI": 0, "QAT": 0, "REU": 0, "ROU": 0, "RUS": 0, "RWA": 0, "SHN": 0, "KNA": 0, "LCA": 0, "SPM": 0,
	"VCT": 0, "WSM": 0, "SMR": 0, "STP": 0, "SAU": 0, "SEN": 0, "SRB": 0, "SYC": 0, "SLE": 0, "SGP": 0, "SVK": 0,
	"SVN": 0, "SLB": 0, "SOM": 0, "ZAF": 0, "SGS": 0, "ESP": 0, "LKA": 0, "SDN": 0, "SUR": 0, "SJM": 0, "SWZ": 0,
	"SWE": 0, "CHE": 0, "SYR": 0, "TWN": 0, "TJK": 0, "TZA": 0, "THA": 0, "TLS": 0, "TGO": 0, "TKL": 0, "TON": 0,
	"TTO": 0, "TUN": 0, "TUR": 0, "TKM": 0, "TCA": 0, "TUV": 0, "UGA": 0, "UKR": 0, "ARE": 0, "GBR": 0, "USA": 0,
	"UMI": 0, "URY": 0, "UZB": 0, "VUT": 0, "VEN": 0, "VNM": 0, "VGB": 0, "VIR": 0, "WLF": 0, "ESH": 0, "YEM": 0,
	"YUG": 0, "ZMB": 0, "ZWE": 0,
}

// Verify Nation abbreviation if validate on 3 chars
func IsVaildNation(abbr string) bool {
	abbr = strings.ToUpper(abbr)
	_, ok := Nations[abbr]
	return ok
}
