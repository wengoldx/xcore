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

// Region regions data
type Region struct {
	string          // Simple use Region value as string type
	Phone    string // Phone region number
	TimeDiff string // Region timezoom diff
	CnName   string // Contry chinese name
}

const (
	Angola              = "Angola"
	Afghanistan         = "Afghanistan"
	Albania             = "Albania"
	Algeria             = "Algeria"
	Andorra             = "Andorra"
	Anguilla            = "Anguilla"
	AntiguaBarbuda      = "Antigua and Barbuda"
	Argentina           = "Argentina"
	Armenia             = "Armenia"
	Ascension           = "Ascension"
	Australia           = "Australia"
	Austria             = "Austria"
	Azerbaijan          = "Azerbaijan"
	Bahamas             = "Bahamas"
	Bahrain             = "Bahrain"
	Bangladesh          = "Bangladesh"
	Barbados            = "Barbados"
	Belarus             = "Belarus"
	Belgium             = "Belgium"
	Belize              = "Belize"
	Benin               = "Benin"
	BermudaIs           = "Bermuda Is."
	Bolivia             = "Bolivia "
	Botswana            = "Botswana"
	Brazil              = "Brazil"
	Brunei              = "Brunei"
	Bulgaria            = "Bulgaria"
	BurkinaFaso         = "Burkina-faso"
	Burma               = "Burma"
	Burundi             = "Burundi"
	Cameroon            = "Cameroon"
	Canada              = "Canada"
	CaymanIs            = "Cayman Is."
	CentralAfricanRep   = "Central African Republic"
	Chad                = "Chad"
	Chile               = "Chile"
	China               = "China"
	Colombia            = "Colombia"
	Congo               = "Congo"
	CookIs              = "Cook Is."
	CostaRica           = "Costa Rica"
	Cuba                = "Cuba"
	Cyprus              = "Cyprus"
	CzechRep            = "Czech Republic"
	Denmark             = "Denmark"
	Djibouti            = "Djibouti"
	DominicaRep         = "Dominica Rep."
	Ecuador             = "Ecuador"
	Egypt               = "Egypt"
	EISalvador          = "EI Salvador"
	Estonia             = "Estonia"
	Ethiopia            = "Ethiopia"
	Fiji                = "Fiji"
	Finland             = "Finland"
	France              = "France"
	FrenchGuiana        = "French Guiana"
	Gabon               = "Gabon"
	Gambia              = "Gambia"
	Georgia             = "Georgia"
	Germany             = "Germany"
	Ghana               = "Ghana"
	Gibraltar           = "Gibraltar"
	Greece              = "Greece"
	Grenada             = "Grenada"
	Guam                = "Guam"
	Guatemala           = "Guatemala"
	Guinea              = "Guinea"
	Guyana              = "Guyana"
	Haiti               = "Haiti"
	Honduras            = "Honduras"
	Hongkong            = "Hongkong"
	Hungary             = "Hungary"
	Iceland             = "Iceland"
	India               = "India"
	Indonesia           = "Indonesia"
	Iran                = "Iran"
	Iraq                = "Iraq"
	Ireland             = "Ireland"
	Israel              = "Israel"
	Italy               = "Italy"
	IvoryCoast          = "Ivory Coast"
	Jamaica             = "Jamaica"
	Japan               = "Japan"
	Jordan              = "Jordan"
	Kampuchea           = "Kampuchea (Cambodia)"
	Kazakstan           = "Kazakstan"
	Kenya               = "Kenya"
	Korea               = "Korea"
	Kuwait              = "Kuwait"
	Kyrgyzstan          = "Kyrgyzstan"
	Laos                = "Laos"
	Latvia              = "Latvia"
	Lebanon             = "Lebanon"
	Lesotho             = "Lesotho"
	Liberia             = "Liberia"
	Libya               = "Libya"
	Liechtenstein       = "Liechtenstein"
	Lithuania           = "Lithuania"
	Luxembourg          = "Luxembourg"
	Macao               = "Macao"
	Madagascar          = "Madagascar"
	Malawi              = "Malawi"
	Malaysia            = "Malaysia"
	Maldives            = "Maldives"
	Mali                = "Mali"
	Malta               = "Malta"
	MarianaIs           = "Mariana Is"
	Martinique          = "Martinique"
	Mauritius           = "Mauritius"
	Mexico              = "Mexico"
	MoldovaRep          = "Republic of Moldova"
	Monaco              = "Monaco"
	Mongolia            = "Mongolia"
	MontserratIs        = "Montserrat Is"
	Morocco             = "Morocco"
	Mozambique          = "Mozambique"
	Namibia             = "Namibia"
	Nauru               = "Nauru"
	Nepal               = "Nepal"
	NetheriandsAntilles = "Netheriands Antilles"
	Netherlands         = "Netherlands"
	NewZealand          = "New Zealand"
	Nicaragua           = "Nicaragua"
	Niger               = "Niger"
	Nigeria             = "Nigeria"
	NorthKorea          = "North Korea"
	Norway              = "Norway"
	Oman                = "Oman"
	Pakistan            = "Pakistan"
	Panama              = "Panama"
	PapuaNewCuinea      = "Papua New Cuinea"
	Paraguay            = "Paraguay"
	Peru                = "Peru"
	Philippines         = "Philippines"
	Poland              = "Poland"
	FrenchPolynesia     = "French Polynesia"
	Portugal            = "Portugal"
	PuertoRico          = "Puerto Rico"
	Qatar               = "Qatar"
	Reunion             = "Reunion"
	Romania             = "Romania"
	Russia              = "Russia"
	SaintLueia          = "Saint Lueia"
	SaintVincent        = "Saint Vincent"
	SamoaEastern        = "Samoa Eastern"
	SamoaWestern        = "Samoa Western"
	SanMarino           = "San Marino"
	SaoTomePrincipe     = "Sao Tome and Principe"
	SaudiArabia         = "Saudi Arabia"
	Senegal             = "Senegal"
	Seychelles          = "Seychelles"
	SierraLeone         = "Sierra Leone"
	Singapore           = "Singapore"
	Slovakia            = "Slovakia"
	Slovenia            = "Slovenia"
	SolomonIs           = "Solomon Is"
	Somali              = "Somali"
	SouthAfrica         = "South Africa"
	Spain               = "Spain"
	SriLanka            = "Sri Lanka"
	StLucia             = "St.Lucia"
	StVincent           = "St.Vincent"
	Sudan               = "Sudan"
	Suriname            = "Suriname"
	Swaziland           = "Swaziland"
	Sweden              = "Sweden"
	Switzerland         = "Switzerland"
	Syria               = "Syria"
	Taiwan              = "Taiwan"
	Tajikstan           = "Tajikstan"
	Tanzania            = "Tanzania"
	Thailand            = "Thailand"
	Togo                = "Togo"
	Tonga               = "Tonga"
	TrinidadTobago      = "Trinidad and Tobago"
	Tunisia             = "Tunisia"
	Turkey              = "Turkey"
	Turkmenistan        = "Turkmenistan"
	Uganda              = "Uganda"
	Ukraine             = "Ukraine"
	UnitedArabEmirates  = "United Arab Emirates"
	UnitedKiongdom      = "United Kiongdom "
	USA                 = "United States of America"
	Uruguay             = "Uruguay"
	Uzbekistan          = "Uzbekistan"
	Venezuela           = "Venezuela"
	Vietnam             = "Vietnam"
	Yemen               = "Yemen"
	Yugoslavia          = "Yugoslavia"
	Zimbabwe            = "Zimbabwe"
	Zaire               = "Zaire"
	Zambia              = "Zambia"
)

// regionsCache regions information cache
var regionsCache = map[string]*Region{
	/*                    CODE        PHONE         TIME DIFF      COUNTRY */
	Angola:              {"AO" /* */, "244" /*  */, "-7" /*    */, "安哥拉"},
	Afghanistan:         {"AF" /* */, "93" /*   */, "0" /*     */, "阿富汗"},
	Albania:             {"AL" /* */, "355" /*  */, "-7" /*    */, "阿尔巴尼亚"},
	Algeria:             {"DZ" /* */, "213" /*  */, "-8" /*    */, "阿尔及利亚"},
	Andorra:             {"AD" /* */, "376" /*  */, "-8" /*    */, "安道尔共和国"},
	Anguilla:            {"AI" /* */, "1264" /* */, "-12" /*   */, "安圭拉岛"},
	AntiguaBarbuda:      {"AG" /* */, "1268" /* */, "-12" /*   */, "安提瓜和巴布达"},
	Argentina:           {"AR" /* */, "54" /*   */, "-11" /*   */, "阿根廷"},
	Armenia:             {"AM" /* */, "374" /*  */, "-6" /*    */, "亚美尼亚"},
	Ascension:           {"" /*   */, "247" /*  */, "-8" /*    */, "阿森松"},
	Australia:           {"AU" /* */, "61" /*   */, "2" /*     */, "澳大利亚"},
	Austria:             {"AT" /* */, "43" /*   */, "-7" /*    */, "奥地利"},
	Azerbaijan:          {"AZ" /* */, "994" /*  */, "-5" /*    */, "阿塞拜疆"},
	Bahamas:             {"BS" /* */, "1242" /* */, "-13" /*   */, "巴哈马"},
	Bahrain:             {"BH" /* */, "973" /*  */, "-5" /*    */, "巴林"},
	Bangladesh:          {"BD" /* */, "880" /*  */, "-2" /*    */, "孟加拉国"},
	Barbados:            {"BB" /* */, "1246" /* */, "-12" /*   */, "巴巴多斯"},
	Belarus:             {"BY" /* */, "375" /*  */, "-6" /*    */, "白俄罗斯"},
	Belgium:             {"BE" /* */, "32" /*   */, "-7" /*    */, "比利时"},
	Belize:              {"BZ" /* */, "501" /*  */, "-14" /*   */, "伯利兹"},
	Benin:               {"BJ" /* */, "229" /*  */, "-7" /*    */, "贝宁"},
	BermudaIs:           {"BM" /* */, "1441" /* */, "-12" /*   */, "百慕大群岛"},
	Bolivia:             {"BO" /* */, "591" /*  */, "-12" /*   */, "玻利维亚"},
	Botswana:            {"BW" /* */, "267" /*  */, "-6" /*    */, "博茨瓦纳"},
	Brazil:              {"BR" /* */, "55" /*   */, "-11" /*   */, "巴西"},
	Brunei:              {"BN" /* */, "673" /*  */, "0" /*     */, "文莱"},
	Bulgaria:            {"BG" /* */, "359" /*  */, "-6" /*    */, "保加利亚"},
	BurkinaFaso:         {"BF" /* */, "226" /*  */, "-8" /*    */, "布基纳法索"},
	Burma:               {"MM" /* */, "95" /*   */, "-1.3" /*  */, "缅甸"},
	Burundi:             {"BI" /* */, "257" /*  */, "-6" /*    */, "布隆迪"},
	Cameroon:            {"CM" /* */, "237" /*  */, "-7" /*    */, "喀麦隆"},
	Canada:              {"CA" /* */, "1" /*    */, "-13" /*   */, "加拿大"},
	CaymanIs:            {"" /*   */, "1345" /* */, "-13" /*   */, "开曼群岛"},
	CentralAfricanRep:   {"CF" /* */, "236" /*  */, "-7" /*    */, "中非共和国"},
	Chad:                {"TD" /* */, "235" /*  */, "-7" /*    */, "乍得"},
	Chile:               {"CL" /* */, "56" /*   */, "-13" /*   */, "智利"},
	China:               {"CN" /* */, "86" /*   */, "0" /*     */, "中国"},
	Colombia:            {"CO" /* */, "57" /*   */, "0" /*     */, "哥伦比亚"},
	Congo:               {"CG" /* */, "242" /*  */, "-7" /*    */, "刚果"},
	CookIs:              {"CK" /* */, "682" /*  */, "-18.3" /* */, "库克群岛"},
	CostaRica:           {"CR" /* */, "506" /*  */, "-14" /*   */, "哥斯达黎加"},
	Cuba:                {"CU" /* */, "53" /*   */, "-13" /*   */, "古巴"},
	Cyprus:              {"CY" /* */, "357" /*  */, "-6" /*    */, "塞浦路斯"},
	CzechRep:            {"CZ" /* */, "420" /*  */, "-7" /*    */, "捷克"},
	Denmark:             {"DK" /* */, "45" /*   */, "-7" /*    */, "丹麦"},
	Djibouti:            {"DJ" /* */, "253" /*  */, "-5" /*    */, "吉布提"},
	DominicaRep:         {"DO" /* */, "1890" /* */, "-13" /*   */, "多米尼加共和国"},
	Ecuador:             {"EC" /* */, "593" /*  */, "-13" /*   */, "厄瓜多尔"},
	Egypt:               {"EG" /* */, "20" /*   */, "-6" /*    */, "埃及"},
	EISalvador:          {"SV" /* */, "503" /*  */, "-14" /*   */, "萨尔瓦多"},
	Estonia:             {"EE" /* */, "372" /*  */, "-5" /*    */, "爱沙尼亚"},
	Ethiopia:            {"ET" /* */, "251" /*  */, "-5" /*    */, "埃塞俄比亚"},
	Fiji:                {"FJ" /* */, "679" /*  */, "4" /*     */, "斐济"},
	Finland:             {"FI" /* */, "358" /*  */, "-6" /*    */, "芬兰"},
	France:              {"FR" /* */, "33" /*   */, "-8" /*    */, "法国"},
	FrenchGuiana:        {"GF" /* */, "594" /*  */, "-12" /*   */, "法属圭亚那"},
	Gabon:               {"GA" /* */, "241" /*  */, "-7" /*    */, "加蓬"},
	Gambia:              {"GM" /* */, "220" /*  */, "-8" /*    */, "冈比亚"},
	Georgia:             {"GE" /* */, "995" /*  */, "0" /*     */, "格鲁吉亚"},
	Germany:             {"DE" /* */, "49" /*   */, "-7" /*    */, "德国"},
	Ghana:               {"GH" /* */, "233" /*  */, "-8" /*    */, "加纳"},
	Gibraltar:           {"GI" /* */, "350" /*  */, "-8" /*    */, "直布罗陀"},
	Greece:              {"GR" /* */, "30" /*   */, "-6" /*    */, "希腊"},
	Grenada:             {"GD" /* */, "1809" /* */, "-14" /*   */, "格林纳达"},
	Guam:                {"GU" /* */, "1671" /* */, "2" /*     */, "关岛"},
	Guatemala:           {"GT" /* */, "502" /*  */, "-14" /*   */, "危地马拉"},
	Guinea:              {"GN" /* */, "224" /*  */, "-8" /*    */, "几内亚"},
	Guyana:              {"GY" /* */, "592" /*  */, "-11" /*   */, "圭亚那"},
	Haiti:               {"HT" /* */, "509" /*  */, "-13" /*   */, "海地"},
	Honduras:            {"HN" /* */, "504" /*  */, "-14" /*   */, "洪都拉斯"},
	Hongkong:            {"HK" /* */, "852" /*  */, "0" /*     */, "香港"},
	Hungary:             {"HU" /* */, "36" /*   */, "-7" /*    */, "匈牙利"},
	Iceland:             {"IS" /* */, "354" /*  */, "-9" /*    */, "冰岛"},
	India:               {"IN" /* */, "91" /*   */, "-2.3" /*  */, "印度"},
	Indonesia:           {"ID" /* */, "62" /*   */, "-0.3" /*  */, "印度尼西亚"},
	Iran:                {"IR" /* */, "98" /*   */, "-4.3" /*  */, "伊朗"},
	Iraq:                {"IQ" /* */, "964" /*  */, "-5" /*    */, "伊拉克"},
	Ireland:             {"IE" /* */, "353" /*  */, "-4.3" /*  */, "爱尔兰"},
	Israel:              {"IL" /* */, "972" /*  */, "-6" /*    */, "以色列"},
	Italy:               {"IT" /* */, "39" /*   */, "-7" /*    */, "意大利"},
	IvoryCoast:          {"" /*   */, "225" /*  */, "-6" /*    */, "科特迪瓦"},
	Jamaica:             {"JM" /* */, "1876" /* */, "-12" /*   */, "牙买加"},
	Japan:               {"JP" /* */, "81" /*   */, "1" /*     */, "日本"},
	Jordan:              {"JO" /* */, "962" /*  */, "-6" /*    */, "约旦"},
	Kampuchea:           {"KH" /* */, "855" /*  */, "-1" /*    */, "柬埔寨"},
	Kazakstan:           {"KZ" /* */, "327" /*  */, "-5" /*    */, "哈萨克斯坦"},
	Kenya:               {"KE" /* */, "254" /*  */, "-5" /*    */, "肯尼亚"},
	Korea:               {"KR" /* */, "82" /*   */, "1" /*     */, "韩国"},
	Kuwait:              {"KW" /* */, "965" /*  */, "-5" /*    */, "科威特"},
	Kyrgyzstan:          {"KG" /* */, "331" /*  */, "-5" /*    */, "吉尔吉斯坦"},
	Laos:                {"LA" /* */, "856" /*  */, "-1" /*    */, "老挝"},
	Latvia:              {"LV" /* */, "371" /*  */, "-5" /*    */, "拉脱维亚"},
	Lebanon:             {"LB" /* */, "961" /*  */, "-6" /*    */, "黎巴嫩"},
	Lesotho:             {"LS" /* */, "266" /*  */, "-6" /*    */, "莱索托"},
	Liberia:             {"LR" /* */, "231" /*  */, "-8" /*    */, "利比里亚"},
	Libya:               {"LY" /* */, "218" /*  */, "-6" /*    */, "利比亚"},
	Liechtenstein:       {"LI" /* */, "423" /*  */, "-7" /*    */, "列支敦士登"},
	Lithuania:           {"LT" /* */, "370" /*  */, "-5" /*    */, "立陶宛"},
	Luxembourg:          {"LU" /* */, "352" /*  */, "-7" /*    */, "卢森堡"},
	Macao:               {"MO" /* */, "853" /*  */, "0" /*     */, "澳门"},
	Madagascar:          {"MG" /* */, "261" /*  */, "-5" /*    */, "马达加斯加"},
	Malawi:              {"MW" /* */, "265" /*  */, "-6" /*    */, "马拉维"},
	Malaysia:            {"MY" /* */, "60" /*   */, "-0.5" /*  */, "马来西亚"},
	Maldives:            {"MV" /* */, "960" /*  */, "-7" /*    */, "马尔代夫"},
	Mali:                {"ML" /* */, "223" /*  */, "-8" /*    */, "马里"},
	Malta:               {"MT" /* */, "356" /*  */, "-7" /*    */, "马耳他"},
	MarianaIs:           {"" /*   */, "1670" /* */, "1" /*     */, "马里亚那群岛"},
	Martinique:          {"" /*   */, "596" /*  */, "-12" /*   */, "马提尼克"},
	Mauritius:           {"MU" /* */, "230" /*  */, "-4" /*    */, "毛里求斯"},
	Mexico:              {"MX" /* */, "52" /*   */, "-15" /*   */, "墨西哥"},
	MoldovaRep:          {"MD" /* */, "373" /*  */, "-5" /*    */, "摩尔多瓦"},
	Monaco:              {"MC" /* */, "377" /*  */, "-7" /*    */, "摩纳哥"},
	Mongolia:            {"MN" /* */, "976" /*  */, "0" /*     */, "蒙古"},
	MontserratIs:        {"MS" /* */, "1664" /* */, "-12" /*   */, "蒙特塞拉特岛"},
	Morocco:             {"MA" /* */, "212" /*  */, "-6" /*    */, "摩洛哥"},
	Mozambique:          {"MZ" /* */, "258" /*  */, "-6" /*    */, "莫桑比克"},
	Namibia:             {"NA" /* */, "264" /*  */, "-7" /*    */, "纳米比亚"},
	Nauru:               {"NR" /* */, "674" /*  */, "4" /*     */, "瑙鲁"},
	Nepal:               {"NP" /* */, "977" /*  */, "-2.3" /*  */, "尼泊尔"},
	NetheriandsAntilles: {"" /*   */, "599" /*  */, "-12" /*   */, "荷属安的列斯"},
	Netherlands:         {"NL" /* */, "31" /*   */, "-7" /*    */, "荷兰"},
	NewZealand:          {"NZ" /* */, "64" /*   */, "4" /*     */, "新西兰"},
	Nicaragua:           {"NI" /* */, "505" /*  */, "-14" /*   */, "尼加拉瓜"},
	Niger:               {"NE" /* */, "227" /*  */, "-8" /*    */, "尼日尔"},
	Nigeria:             {"NG" /* */, "234" /*  */, "-7" /*    */, "尼日利亚"},
	NorthKorea:          {"KP" /* */, "850" /*  */, "1" /*     */, "朝鲜"},
	Norway:              {"NO" /* */, "47" /*   */, "-7" /*    */, "挪威"},
	Oman:                {"OM" /* */, "968" /*  */, "-4" /*    */, "阿曼"},
	Pakistan:            {"PK" /* */, "92" /*   */, "-2.3" /*  */, "巴基斯坦"},
	Panama:              {"PA" /* */, "507" /*  */, "-13" /*   */, "巴拿马"},
	PapuaNewCuinea:      {"PG" /* */, "675" /*  */, "2" /*     */, "巴布亚新几内亚"},
	Paraguay:            {"PY" /* */, "595" /*  */, "-12" /*   */, "巴拉圭"},
	Peru:                {"PE" /* */, "51" /*   */, "-13" /*   */, "秘鲁"},
	Philippines:         {"PH" /* */, "63" /*   */, "0" /*     */, "菲律宾"},
	Poland:              {"PL" /* */, "48" /*   */, "-7" /*    */, "波兰"},
	FrenchPolynesia:     {"PF" /* */, "689" /*  */, "3" /*     */, "法属玻利尼西亚"},
	Portugal:            {"PT" /* */, "351" /*  */, "-8" /*    */, "葡萄牙"},
	PuertoRico:          {"PR" /* */, "1787" /* */, "-12" /*   */, "波多黎各"},
	Qatar:               {"QA" /* */, "974" /*  */, "-5" /*    */, "卡塔尔"},
	Reunion:             {"" /*   */, "262" /*  */, "-4" /*    */, "留尼旺"},
	Romania:             {"RO" /* */, "40" /*   */, "-6" /*    */, "罗马尼亚"},
	Russia:              {"RU" /* */, "7" /*    */, "-5" /*    */, "俄罗斯"},
	SaintLueia:          {"LC" /* */, "1758" /* */, "-12" /*   */, "圣卢西亚"},
	SaintVincent:        {"VC" /* */, "1784" /* */, "-12" /*   */, "圣文森特岛"},
	SamoaEastern:        {"" /*   */, "684" /*  */, "-19" /*   */, "东萨摩亚(美)"},
	SamoaWestern:        {"" /*   */, "685" /*  */, "-19" /*   */, "西萨摩亚"},
	SanMarino:           {"SM" /* */, "378" /*  */, "-7" /*    */, "圣马力诺"},
	SaoTomePrincipe:     {"ST" /* */, "239" /*  */, "-8" /*    */, "圣多美和普林西比"},
	SaudiArabia:         {"SA" /* */, "966" /*  */, "-5" /*    */, "沙特阿拉伯"},
	Senegal:             {"SN" /* */, "221" /*  */, "-8" /*    */, "塞内加尔"},
	Seychelles:          {"SC" /* */, "248" /*  */, "-4" /*    */, "塞舌尔"},
	SierraLeone:         {"SL" /* */, "232" /*  */, "-8" /*    */, "塞拉利昂"},
	Singapore:           {"SG" /* */, "65" /*   */, "0.3" /*   */, "新加坡"},
	Slovakia:            {"SK" /* */, "421" /*  */, "-7" /*    */, "斯洛伐克"},
	Slovenia:            {"SI" /* */, "386" /*  */, "-7" /*    */, "斯洛文尼亚"},
	SolomonIs:           {"SB" /* */, "677" /*  */, "3" /*     */, "所罗门群岛"},
	Somali:              {"SO" /* */, "252" /*  */, "-5" /*    */, "索马里"},
	SouthAfrica:         {"ZA" /* */, "27" /*   */, "-6" /*    */, "南非"},
	Spain:               {"ES" /* */, "34" /*   */, "-8" /*    */, "西班牙"},
	SriLanka:            {"LK" /* */, "94" /*   */, "0" /*     */, "斯里兰卡"},
	StLucia:             {"LC" /* */, "1758" /* */, "-12" /*   */, "圣卢西亚"},
	StVincent:           {"VC" /* */, "1784" /* */, "-12" /*   */, "圣文森特"},
	Sudan:               {"SD" /* */, "249" /*  */, "-6" /*    */, "苏丹"},
	Suriname:            {"SR" /* */, "597" /*  */, "-11.3" /* */, "苏里南"},
	Swaziland:           {"SZ" /* */, "268" /*  */, "-6" /*    */, "斯威士兰"},
	Sweden:              {"SE" /* */, "46" /*   */, "-7" /*    */, "瑞典"},
	Switzerland:         {"CH" /* */, "41" /*   */, "-7" /*    */, "瑞士"},
	Syria:               {"SY" /* */, "963" /*  */, "-6" /*    */, "叙利亚"},
	Taiwan:              {"TW" /* */, "886" /*  */, "0" /*     */, "台湾省"},
	Tajikstan:           {"TJ" /* */, "992" /*  */, "-5" /*    */, "塔吉克斯坦"},
	Tanzania:            {"TZ" /* */, "255" /*  */, "-5" /*    */, "坦桑尼亚"},
	Thailand:            {"TH" /* */, "66" /*   */, "-1" /*    */, "泰国"},
	Togo:                {"TG" /* */, "228" /*  */, "-8" /*    */, "多哥"},
	Tonga:               {"TO" /* */, "676" /*  */, "4" /*     */, "汤加"},
	TrinidadTobago:      {"TT" /* */, "1809" /* */, "-12" /*   */, "特立尼达和多巴哥"},
	Tunisia:             {"TN" /* */, "216" /*  */, "-7" /*    */, "突尼斯"},
	Turkey:              {"TR" /* */, "90" /*   */, "-6" /*    */, "土耳其"},
	Turkmenistan:        {"TM" /* */, "993" /*  */, "-5" /*    */, "土库曼斯坦"},
	Uganda:              {"UG" /* */, "256" /*  */, "-5" /*    */, "乌干达"},
	Ukraine:             {"UA" /* */, "380" /*  */, "-5" /*    */, "乌克兰"},
	UnitedArabEmirates:  {"AE" /* */, "971" /*  */, "-4" /*    */, "阿拉伯联合酋长国"},
	UnitedKiongdom:      {"GB" /* */, "44" /*   */, "-8" /*    */, "英国"},
	USA:                 {"US" /* */, "1" /*    */, "-13" /*   */, "美国"},
	Uruguay:             {"UY" /* */, "598" /*  */, "-10.3" /* */, "乌拉圭"},
	Uzbekistan:          {"UZ" /* */, "233" /*  */, "-5" /*    */, "乌兹别克斯坦"},
	Venezuela:           {"VE" /* */, "58" /*   */, "-12.3" /* */, "委内瑞拉"},
	Vietnam:             {"VN" /* */, "84" /*   */, "-1" /*    */, "越南"},
	Yemen:               {"YE" /* */, "967" /*  */, "-5" /*    */, "也门"},
	Yugoslavia:          {"YU" /* */, "381" /*  */, "-7" /*    */, "南斯拉夫"},
	Zimbabwe:            {"ZW" /* */, "263" /*  */, "-6" /*    */, "津巴布韦"},
	Zaire:               {"ZR" /* */, "243" /*  */, "-7" /*    */, "扎伊尔"},
	Zambia:              {"ZM" /* */, "260" /*  */, "-6" /*    */, "赞比亚"},
}

// GetRegion get region information by country
func GetRegion(country string) *Region {
	if region, ok := regionsCache[country]; ok {
		return &Region{country, region.Phone, region.TimeDiff, region.CnName}
	}
	return nil
}

// GetRegionByCode get country and region information by code and phone
func GetRegionByCode(code string, phone ...string) (string, *Region) {
	regions, lastcountry := make(map[string]*Region), ""
	for country, region := range regionsCache {
		if region.string == code {
			regions[country] = region
			lastcountry = country
		}
	}

	if len(phone) > 0 && phone[0] != "" {
		for country, region := range regions {
			if region.Phone == phone[0] {
				return country, &Region{
					region.string, region.Phone, region.TimeDiff, region.CnName,
				}
			}
		}
	} else if lastcountry != "" {
		region := regions[lastcountry]
		return lastcountry, &Region{
			region.string, region.Phone, region.TimeDiff, region.CnName,
		}
	}
	return "", nil
}
