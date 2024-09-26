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
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"github.com/mozillazg/go-pinyin"
	"github.com/wengoldx/xcore/logger"
)

// Try try-catch-finaly method
func Try(do func(), catcher func(error), finaly ...func()) {
	defer func() {
		if i := recover(); i != nil {
			execption := errors.New(fmt.Sprint(i))
			logger.E("Catched exception:", i)
			catcher(execption)
			if len(finaly) > 0 {
				finaly[0]()
			}
		}
	}()
	do()
}

// Condition return the trueData when pass the condition, or return falseData
//
// `USAGE` :
//
//	// use as follow to return diffrent type value, but the input
//	// true and false params MUST BE no-nil datas.
//	a := Condition(condition, trueString, falseString)	// return any
//	b := Condition(condition, trueInt, falseInt).(int)
//	c := Condition(condition, trueInt64, falseInt64).(int64)
//	d := Condition(condition, trueFloat, falseFloat).(float64)
//	e := Condition(condition, trueDur, falseDur).(time.Duration)
//	f := Condition(condition, trueString, falseString).(string)
func Condition(condition bool, trueData any, falseData any) any {
	if condition {
		return trueData
	}
	return falseData
}

// Check the variable params, return the first value if exist any,
// or return default string
//
// `USAGE` :
//
//	func sample(params ...string) {
//		a := GetVariable(params, "def-value").(string)
//		// Do saming here
//	}
func GetVariable(params any, defvalue any) any {
	pv := reflect.ValueOf(params)
	if pv.Len() > 0 {
		return pv.Index(0).Interface()
	}
	return defvalue
}

// Contain check the given string list if contains item.
//
// You should call Distinct() to filter out the repeat items.
func Contain(list *[]string, item string) bool {
	for _, v := range *list {
		if v == item {
			return true
		}
	}
	return false
}

// Distinct remove duplicate string from given array.
//
// You should call Contain() only for check if exist sub string.
func Distinct(src *[]string) []string {
	dest := make(map[string]byte)
	for _, str := range *src {
		if _, ok := dest[str]; !ok {
			dest[str] = byte(0)
		}
	}

	st := []string{}
	for str := range dest {
		st = append(st, str)
	}
	return st
}

// TrimEmpty remove empty string, it maybe return empty result array.
func TrimEmpty(src *[]string) []string {
	dst := []string{}
	for _, str := range *src {
		if str != "" {
			dst = append(dst, str)
		}
	}
	return dst
}

// To2Digits fill zero if input digit not enough 2
func To2Digits(input any) string {
	return fmt.Sprintf("%02d", input)
}

// To2Digits fill zero if input digit not enough 3
func To3Digits(input any) string {
	return fmt.Sprintf("%03d", input)
}

// ToNDigits fill zero if input digit not enough N
func ToNDigits(input any, n int) string {
	return fmt.Sprintf("%0"+strconv.Itoa(n)+"d", input)
}

// ToMap transform given struct data to map data, the transform struct
// feilds must using json tag to mark the map key.
//
// ---
//
//	type struct Sample {
//		Name string `json:"name"`
//	}
//	d := Sample{ Name : "name_value" }
//	md, _ := comm.ToMap(d)
//	// md data format is {
//	//     "name" : "name_value"
//	// }
func ToMap(input any) (map[string]any, error) {
	out := make(map[string]any)
	buf, err := json.Marshal(input)
	if err != nil {
		logger.E("Marshal input struct err:", err)
		return nil, err
	}

	// json buffer decode to map
	d := json.NewDecoder(bytes.NewReader(buf))
	d.UseNumber()
	if err = d.Decode(&out); err != nil {
		logger.E("Decode json data to map err:", err)
		return nil, err
	}

	return out, nil
}

// ToXMLString transform given struct data to xml string
func ToXMLString(input any) (string, error) {
	buf, err := xml.Marshal(input)
	if err != nil {
		logger.E("Marshal input to XML err:", err)
		return "", err
	}
	return string(buf), nil
}

// ToXMLReplace transform given struct data to xml string, ant then
// replace indicated fileds or values, to form param must not empty,
// but the to param allow set empty when use to remove all form keyworlds.
func ToXMLReplace(input any, from, to string) (string, error) {
	xmlout, err := ToXMLString(input)
	if err != nil {
		return "", err
	}

	trimsrc := strings.TrimSpace(from)
	if trimsrc != "" {
		logger.I("Replace xml string from:", trimsrc, "to:", to)
		xmlout = strings.Replace(xmlout, trimsrc, to, -1)
	}
	return xmlout, nil
}

// JoinLines combine strings into multiple lines
func JoinLines(inputs ...string) string {
	packet := ""
	for _, line := range inputs {
		packet += line + "\n"
	}
	return packet
}

// SplitTrim extend strings.Split to trim space and sub strings before split.
//
//	// Use Case 1: head or end maybe exist space chars or sub strings
//	str := "01/23/34"         // or " 01/23/34", "01/23/34 ", " 01/23/34 "
//	                          // or "/01/23/34", "01/23/34/", "/01/23/34/"
//	                          // or " /01/23/34/ "
//	utils.SplitTrim(str, "/") // output [01 23 34]
//
//	// Use Case 2: source string only have space chars or split strings
//	str := "/"                // or " /", "/ ", " / ", " // ", " //// "
//	utils.SplitTrim(str, "/") // output []
func SplitTrim(src, sub string) []string {
	src = strings.Trim(strings.TrimSpace(src), sub)
	st := strings.Split(src, sub)
	if len(st) == 1 && st[0] == "" {
		return []string{}
	}
	return st
}

// SplitAfterTrim extend strings.SplitAfter to trim all space chars and sub strings.
//
//	// Use Case 1: remove all space chars and filter out sub string items
//	str := "01/23/34"              // or "/01/23/34", " /01/23/34 / ", "/0 1/23 /34"
//	utils.SplitAfterTrim(str, "/") // output [01/ 23/ 34/]
//
//	// Use Case 2: source string only have space chars or split strings
//	str := "/"                     // or " / ", " // ", " / / ", " /// "
//	utils.SplitAfterTrim(str, "/") // output []
//
// `Warning` : DO NOT contain space chars in the sub string!
func SplitAfterTrim(src, sub string) []string {
	src = strings.ReplaceAll(src, " ", "")
	if !strings.HasSuffix(src, sub) {
		src += sub
	}

	st, rst := strings.SplitAfter(src, sub), []string{}
	for _, str := range st {
		if str != "" && str != sub {
			rst = append(rst, str)
		}
	}
	return rst
}

// GetSortKey get first letter of Chinese Pinyin
func GetSortKey(str string) string {
	if str == "" { // check the input param
		return "*"
	}

	// get the first char and verify if it is a~Z char
	firstChar, sortKey := []rune(str)[0], ""
	isAZchar, err := regexp.Match("[a-zA-Z]", []byte(str))
	if err != nil {
		logger.E("Regexp match err:", err)
		return "*"
	}

	if isAZchar {
		sortKey = string(unicode.ToUpper(firstChar))
	} else {
		if unicode.Is(unicode.Han, firstChar) { // chinese
			str1 := pinyin.LazyConvert(string(firstChar), nil)
			s := []rune(str1[0])
			sortKey = string(unicode.ToUpper(s[0]))
		} else if unicode.IsNumber(firstChar) { // number
			sortKey = string(firstChar)
		} else { // other language
			sortKey = "#"
		}
	}
	return sortKey
}

// Indicate given object if the target type.
//
//	type Car struct {
//		Bland string
//	}
//	car := &Car{Bland : "H5"}
//	if comm.Instanceof(car, reflect.TypeOf(&Car{})) {
//		// intvalue is int, but target type is string
//	}
func InstanceOf(object any, tagtype reflect.Type) bool {
	return reflect.TypeOf(object) == tagtype
}
