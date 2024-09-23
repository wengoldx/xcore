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
	"strings"
)

const (
	lang_MIN  int = iota - 1
	Lang_arIL     // Arabic (Israel)
	Lang_arEG     // Arabic (Egypt)
	Lang_zhCN     // Chinese Simplified
	Lang_zhTW     // Chinese Tradition
	Lang_zhHK     // Chinese Hongkong
	Lang_nlNL     // Dutch(Netherlands)
	Lang_nlBE     // Dutch(Netherlands)
	Lang_enUS     // English(United States)
	Lang_enAU     // English(Australia)
	Lang_enCA     // English(Canada)
	Lang_enIN     // English(India)
	Lang_enIE     // English(Ireland)
	Lang_enNZ     // English(New Zealand)
	Lang_enSG     // English(Singapore)
	Lang_enZA     // English(South Africa)
	Lang_enGB     // English(United Kingdom)
	Lang_frFR     // French
	Lang_frBE     // French
	Lang_frCA     // French
	Lang_frCH     // French
	Lang_deDE     // German
	Lang_deLI     // German
	Lang_deAT     // German
	Lang_deCH     // German
	Lang_itIT     // Italian
	Lang_itCH     // Italian
	Lang_ptBR     // Portuguese
	Lang_ptPT     // Portuguese
	Lang_esES     // Spanish
	Lang_esUS     // Spanish
	Lang_bnBD     // Bengali
	Lang_bnIN     // Bengali
	Lang_hrHR     // Croatian
	Lang_csCZ     // Czech
	Lang_daDK     // Danish
	Lang_elGR     // Greek
	Lang_heIL     // Hebrew
	Lang_iwIL     // Hebrew
	Lang_hiIN     // Hindi
	Lang_huHU     // Hungarian
	Lang_inID     // Indonesian
	Lang_jaJP     // Japanese
	Lang_koKR     // Korean
	Lang_msMY     // Malay
	Lang_faIR     // Perisan
	Lang_plPL     // Polish
	Lang_roRO     // Romanian
	Lang_ruRU     // Russian
	Lang_srRS     // Serbian
	Lang_svSE     // Swedish
	Lang_thTH     // Thai
	Lang_trTR     // Turkey
	Lang_urPK     // Urdu
	Lang_viVN     // Vietnamese
	Lang_caES     // Catalan
	Lang_lvLV     // Latviesu
	Lang_ltLT     // Lithuanian
	Lang_nbNO     // Norwegian
	Lang_skSK     // slovencina
	Lang_slSI     // Slovenian
	Lang_bgBG     // bulgarian
	Lang_ukUA     // Ukrainian
	Lang_tlPH     // Filipino
	Lang_fiFI     // Finnish
	Lang_afZA     // Afrikaans
	Lang_rmCH     // Romansh
	Lang_myZG     // Burmese
	Lang_myMM     // Burmese
	Lang_kmKH     // Khmer
	Lang_amET     // Amharic
	Lang_beBY     // Belarusian
	Lang_etEE     // Estonian
	Lang_swTZ     // Swahili
	Lang_zuZA     // Zulu
	Lang_azAZ     // Azerbaijani
	Lang_hyAM     // Armenian
	Lang_kaGE     // Georgian
	Lang_loLA     // Laotian
	Lang_mnMN     // Mongolian
	Lang_neNP     // Nepali
	Lang_kkKZ     // Kazakh
	Lang_siLK     // Sinhala
	lang_MAX
)

// Language language information
type Lang struct {
	int           // Simple use Lang value as int type
	Key    string // Language unique key
	EnName string // Language english name
	CnName string // Language chinese name
}

const (
	// InvalidLangCode invalid language code
	InvalidLangCode int = -1

	// LangsSeparator multi-langguages separator
	LangsSeparator = ","
)

// langsCache languages information cache
var langsCache = map[int]Lang{
	Lang_arIL: {Lang_arIL, "ar_IL", "Arabic(Israel)" /*          */, "阿拉伯语(以色列)"},
	Lang_arEG: {Lang_arEG, "ar_EG", "Arabic(Egypt)" /*           */, "阿拉伯语(埃及)"},
	Lang_zhCN: {Lang_zhCN, "zh_CN", "Chinese Simplified" /*      */, "中文简体"},
	Lang_zhTW: {Lang_zhTW, "zh_TW", "Chinese Tradition" /*       */, "中文繁体"},
	Lang_zhHK: {Lang_zhHK, "zh_HK", "Chinese Hongkong" /*        */, "中文(香港)"},
	Lang_nlNL: {Lang_nlNL, "nl_NL", "Dutch(Netherlands)" /*      */, "荷兰语"},
	Lang_nlBE: {Lang_nlBE, "nl_BE", "Dutch(Netherlands)" /*      */, "荷兰语(比利时)"},
	Lang_enUS: {Lang_enUS, "en_US", "English(United States)" /*  */, "英语(美国)"},
	Lang_enAU: {Lang_enAU, "en_AU", "English(Australia)" /*      */, "英语(澳大利亚)"},
	Lang_enCA: {Lang_enCA, "en_CA", "English(Canada)" /*         */, "英语(加拿大)"},
	Lang_enIN: {Lang_enIN, "en_IN", "English(India)" /*          */, "英语(印度)"},
	Lang_enIE: {Lang_enIE, "en_IE", "English(Ireland)" /*        */, "英语(爱尔兰)"},
	Lang_enNZ: {Lang_enNZ, "en_NZ", "English(New Zealand)" /*    */, "英语(新西兰)"},
	Lang_enSG: {Lang_enSG, "en_SG", "English(Singapore)" /*      */, "英语(新加波)"},
	Lang_enZA: {Lang_enZA, "en_ZA", "English(South Africa)" /*   */, "英语(南非)"},
	Lang_enGB: {Lang_enGB, "en_GB", "English(United Kingdom)" /* */, "英语(英国)"},
	Lang_frFR: {Lang_frFR, "fr_FR", "French" /*                  */, "法语"},
	Lang_frBE: {Lang_frBE, "fr_BE", "French" /*                  */, "法语(比利时)"},
	Lang_frCA: {Lang_frCA, "fr_CA", "French" /*                  */, "法语(加拿大)"},
	Lang_frCH: {Lang_frCH, "fr_CH", "French" /*                  */, "法语(瑞士)"},
	Lang_deDE: {Lang_deDE, "de_DE", "German" /*                  */, "德语"},
	Lang_deLI: {Lang_deLI, "de_LI", "German" /*                  */, "德语(列支敦斯登)"},
	Lang_deAT: {Lang_deAT, "de_AT", "German" /*                  */, "德语(奥地利)"},
	Lang_deCH: {Lang_deCH, "de_CH", "German" /*                  */, "德语(瑞士)"},
	Lang_itIT: {Lang_itIT, "it_IT", "Italian" /*                 */, "意大利语"},
	Lang_itCH: {Lang_itCH, "it_CH", "Italian" /*                 */, "意大利语(瑞士)"},
	Lang_ptBR: {Lang_ptBR, "pt_BR", "Portuguese" /*              */, "葡萄牙语（巴西）"},
	Lang_ptPT: {Lang_ptPT, "pt_PT", "Portuguese" /*              */, "葡萄牙语"},
	Lang_esES: {Lang_esES, "es_ES", "Spanish" /*                 */, "西班牙语"},
	Lang_esUS: {Lang_esUS, "es_US", "Spanish" /*                 */, "西班牙语(美国)"},
	Lang_bnBD: {Lang_bnBD, "bn_BD", "Bengali" /*                 */, "孟加拉语"},
	Lang_bnIN: {Lang_bnIN, "bn_IN", "Bengali" /*                 */, "孟加拉语(印度)"},
	Lang_hrHR: {Lang_hrHR, "hr_HR", "Croatian" /*                */, "克罗地亚语"},
	Lang_csCZ: {Lang_csCZ, "cs_CZ", "Czech" /*                   */, "捷克语"},
	Lang_daDK: {Lang_daDK, "da_DK", "Danish" /*                  */, "丹麦语"},
	Lang_elGR: {Lang_elGR, "el_GR", "Greek" /*                   */, "希腊语"},
	Lang_heIL: {Lang_heIL, "he_IL", "Hebrew" /*                  */, "希伯来语(以色列)"},
	Lang_iwIL: {Lang_iwIL, "iw_IL", "Hebrew" /*                  */, "希伯来语(以色列)"},
	Lang_hiIN: {Lang_hiIN, "hi_IN", "Hindi" /*                   */, "印度语"},
	Lang_huHU: {Lang_huHU, "hu_HU", "Hungarian" /*               */, "匈牙利语"},
	Lang_inID: {Lang_inID, "in_ID", "Indonesian" /*              */, "印度尼西亚语"},
	Lang_jaJP: {Lang_jaJP, "ja_JP", "Japanese" /*                */, "日语"},
	Lang_koKR: {Lang_koKR, "ko_KR", "Korean" /*                  */, "韩语（朝鲜语）"},
	Lang_msMY: {Lang_msMY, "ms_MY", "Malay" /*                   */, "马来语"},
	Lang_faIR: {Lang_faIR, "fa_IR", "Perisan" /*                 */, "波斯语"},
	Lang_plPL: {Lang_plPL, "pl_PL", "Polish" /*                  */, "波兰语"},
	Lang_roRO: {Lang_roRO, "ro_RO", "Romanian" /*                */, "罗马尼亚语"},
	Lang_ruRU: {Lang_ruRU, "ru_RU", "Russian" /*                 */, "俄罗斯语"},
	Lang_srRS: {Lang_srRS, "sr_RS", "Serbian" /*                 */, "塞尔维亚语"},
	Lang_svSE: {Lang_svSE, "sv_SE", "Swedish" /*                 */, "瑞典语"},
	Lang_thTH: {Lang_thTH, "th_TH", "Thai" /*                    */, "泰语"},
	Lang_trTR: {Lang_trTR, "tr_TR", "Turkey" /*                  */, "土耳其语"},
	Lang_urPK: {Lang_urPK, "ur_PK", "Urdu" /*                    */, "乌尔都语"},
	Lang_viVN: {Lang_viVN, "vi_VN", "Vietnamese" /*              */, "越南语"},
	Lang_caES: {Lang_caES, "ca_ES", "Catalan" /*                 */, "加泰隆语(西班牙)"},
	Lang_lvLV: {Lang_lvLV, "lv_LV", "Latviesu" /*                */, "拉脱维亚语"},
	Lang_ltLT: {Lang_ltLT, "lt_LT", "Lithuanian" /*              */, "立陶宛语"},
	Lang_nbNO: {Lang_nbNO, "nb_NO", "Norwegian" /*               */, "挪威语"},
	Lang_skSK: {Lang_skSK, "sk_SK", "slovencina" /*              */, "斯洛伐克语"},
	Lang_slSI: {Lang_slSI, "sl_SI", "Slovenian" /*               */, "斯洛文尼亚语"},
	Lang_bgBG: {Lang_bgBG, "bg_BG", "bulgarian" /*               */, "保加利亚语"},
	Lang_ukUA: {Lang_ukUA, "uk_UA", "Ukrainian" /*               */, "乌克兰语"},
	Lang_tlPH: {Lang_tlPH, "tl_PH", "Filipino" /*                */, "菲律宾语"},
	Lang_fiFI: {Lang_fiFI, "fi_FI", "Finnish" /*                 */, "芬兰语"},
	Lang_afZA: {Lang_afZA, "af_ZA", "Afrikaans" /*               */, "南非语"},
	Lang_rmCH: {Lang_rmCH, "rm_CH", "Romansh" /*                 */, "罗曼什语(瑞士)"},
	Lang_myZG: {Lang_myZG, "my_ZG", "Burmese(Zawgyi)" /*         */, "缅甸语"},
	Lang_myMM: {Lang_myMM, "my_MM", "Burmese" /*                 */, "缅甸语"},
	Lang_kmKH: {Lang_kmKH, "km_KH", "Khmer" /*                   */, "柬埔寨语"},
	Lang_amET: {Lang_amET, "am_ET", "Amharic" /*                 */, "阿姆哈拉语(埃塞俄比亚)"},
	Lang_beBY: {Lang_beBY, "be_BY", "Belarusian" /*              */, "白俄罗斯语"},
	Lang_etEE: {Lang_etEE, "et_EE", "Estonian" /*                */, "爱沙尼亚语"},
	Lang_swTZ: {Lang_swTZ, "sw_TZ", "Swahili" /*                 */, "斯瓦希里语(坦桑尼亚)"},
	Lang_zuZA: {Lang_zuZA, "zu_ZA", "Zulu" /*                    */, "祖鲁语(南非)"},
	Lang_azAZ: {Lang_azAZ, "az_AZ", "Azerbaijani" /*             */, "阿塞拜疆语"},
	Lang_hyAM: {Lang_hyAM, "hy_AM", "Armenian" /*                */, "亚美尼亚语(亚美尼亚)"},
	Lang_kaGE: {Lang_kaGE, "ka_GE", "Georgian" /*                */, "格鲁吉亚语(格鲁吉亚)"},
	Lang_loLA: {Lang_loLA, "lo_LA", "Laotian" /*                 */, "老挝语(老挝)"},
	Lang_mnMN: {Lang_mnMN, "mn_MN", "Mongolian" /*               */, "蒙古语"},
	Lang_neNP: {Lang_neNP, "ne_NP", "Nepali" /*                  */, "尼泊尔语"},
	Lang_kkKZ: {Lang_kkKZ, "kk_KZ", "Kazakh" /*                  */, "哈萨克语"},
	Lang_siLK: {Lang_siLK, "si_LK", "Sinhala" /*                 */, "僧加罗语(斯里兰卡)"},
}

// GetLanguage get language information by code
func GetLanguage(code int) *Lang {
	if lang, ok := langsCache[code]; ok {
		return &Lang{code, lang.Key, lang.EnName, lang.CnName}
	}
	return &Lang{InvalidLangCode, "", "", ""}
}

// GetLangCode get language code by key
func GetLangCode(key string) int {
	for _, lang := range langsCache {
		if lang.Key == key {
			return lang.int
		}
	}
	return InvalidLangCode
}

// IsValidLang check the given language code if valid
func IsValidLang(code int) bool {
	return code > lang_MIN && code < lang_MAX
}

// AppendLangs append a language key to multi-language string
func AppendLangs(langs string, code int) string {
	langkey := GetLanguage(code).Key
	if langkey != "" && !strings.Contains(langs, langkey) {
		langs = strings.Trim(langs, " ")            // Trim leading and trailing empty chars
		langs = strings.Trim(langs, LangsSeparator) // Trim leading and trailing separators
		return strings.Join([]string{langs, langkey}, LangsSeparator)
	}
	return langs
}

// RemoveLangs remove a language key outof multi-language string
func RemoveLangs(langs string, code int) string {
	langkey := GetLanguage(code).Key
	if langkey != "" && strings.Contains(langs, langkey) {
		langsarr := strings.Split(langs, LangsSeparator)
		for i, existlang := range langsarr {
			existlang = strings.Trim(existlang, " ") // Trim leading and trailing empty chars
			if existlang == langkey {
				last := append(langsarr[:i], langsarr[i+1:]...)
				return strings.Join(last, LangsSeparator)
			}
		}
	}
	return langs
}

// IsContain check the language if exist in multi-language string
func IsContain(langs string, code int) bool {
	langkey := GetLanguage(code).Key
	return (langkey != "" && strings.Contains(langs, langkey))
}
