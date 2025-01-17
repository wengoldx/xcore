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

const _validateCodes = "10X98765432"

var _validateWeights = []int{7, 9, 10, 5, 8, 4, 2, 1, 6, 3, 7, 9, 10, 5, 8, 4, 2}

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
			sum += digit * _validateWeights[i]
		}
	}

	index := sum % 11
	code := _validateCodes[index : index+1]
	return code == last
}

/* ========== Foreign-Relate passport authentication ========== */

// All nations informations on ISO 3166-1
//	@See more visit https://baike.baidu.com/item/ISO%203166-1/5269555?fr=ge_ala
var _nations = NewSets().AddStrings([]string{
	"AFG", "ALA", "ALB", "DZA", "ASM", "AND", "AGO", "AIA", "ATA", "ATG", "ARG", "ARM", "ABW", "AUS", "AUT",
	"AZE", "BHS", "BHR", "BGD", "BRB", "BLR", "BEL", "BLZ", "BEN", "BMU", "BTN", "BOL", "BIH", "BWA", "BVT",
	"BRA", "IOT", "BRN", "BGR", "BFA", "BDI", "KHM", "CMR", "CAN", "CPV", "CYM", "CAF", "TCD", "CHL", "CHN",
	"CXR", "CCK", "COL", "COM", "COG", "COD", "COK", "CRI", "CIV", "HRV", "CUB", "CYP", "CZE", "DNK", "DJI",
	"DMA", "DOM", "ECU", "EGY", "SLV", "GNQ", "ERI", "EST", "ETH", "FLK", "FRO", "FJI", "FIN", "FRA", "GUF",
	"PYF", "ATF", "GAB", "GMB", "GEO", "DEU", "GHA", "GIB", "GRC", "GRL", "GRD", "GLP", "GUM", "GTM", "GGY",
	"GIN", "GNB", "GUY", "HTI", "HMD", "VAT", "HND", "HKG", "HUN", "ISL", "IND", "IDN", "IRN", "IRQ", "IRL",
	"IMN", "ISR", "ITA", "JAM", "JPN", "JEY", "JOR", "KAZ", "KEN", "KIR", "PRK", "KOR", "KWT", "KGZ", "LAO",
	"LVA", "LBN", "LSO", "LBR", "LBY", "LIE", "LTU", "LUX", "MAC", "MKD", "MDG", "MWI", "MYS", "MDV", "MLI",
	"MLT", "MHL", "MTQ", "MRT", "MUS", "MYT", "MEX", "FSM", "MDA", "MCO", "MNG", "MNE", "MSR", "MAR", "MOZ",
	"MMR", "NAM", "NRU", "NPL", "NLD", "ANT", "NCL", "NZL", "NIC", "NER", "NGA", "NIU", "NFK", "MNP", "NOR",
	"OMN", "PAK", "PLW", "PSE", "PAN", "PNG", "PRY", "PER", "PHL", "PCN", "POL", "PRT", "PRI", "QAT", "REU",
	"ROU", "RUS", "RWA", "SHN", "KNA", "LCA", "SPM", "VCT", "WSM", "SMR", "STP", "SAU", "SEN", "SRB", "SYC",
	"SLE", "SGP", "SVK", "SVN", "SLB", "SOM", "ZAF", "SGS", "ESP", "LKA", "SDN", "SUR", "SJM", "SWZ", "SWE",
	"CHE", "SYR", "TWN", "TJK", "TZA", "THA", "TLS", "TGO", "TKL", "TON", "TTO", "TUN", "TUR", "TKM", "TCA",
	"TUV", "UGA", "UKR", "ARE", "GBR", "USA", "UMI", "URY", "UZB", "VUT", "VEN", "VNM", "VGB", "VIR", "WLF",
	"ESH", "YEM", "YUG", "ZMB", "ZWE",
})

// Verify Nation abbreviation if validate on 3 chars
func IsVaildNation(abbr string) bool {
	return _nations.Contain(strings.ToUpper(abbr))
}
