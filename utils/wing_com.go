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
	"cmp"
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

// Standardy build-in types of golang.
type BuildIn interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 |
	float32 | float64 | bool | string | any
}

// Return the postive value when condition is true, or return negative value.
//
//	See utils.BuildIn for more types defined informations.
func Condition[T BuildIn](condition bool, postive T, negative T) T {
	if condition {
		return postive
	}
	return negative
}

// Return the first element value is exist, or return default, it not check
// the got value whether valid.
//
//	See utils.BuildIn for more types defined informations.
func Variable[T BuildIn](src []T, def T) T {
	if len(src) > 0 {
		return src[0]
	}
	return def
}

// Translate strict build-in types values to any type array.
//
//	See utils.BuildIn for more types defined informations.
func ToAnys[T BuildIn](values []T) []any {
	args := []any{}
	for _, value := range values {
		args = append(args, value)
	}
	return args
}

// Translate strict build-in types values to strings array, it will
// auto append '' when values is string type like '12345'.
//
//	See utils.BuildIn for more types defined informations.
func ToStrings[T BuildIn](values []T) []string {
	vs := []string{}
	for _, value := range values {
		switch vt := any(value).(type) {
		case string:
			vs = append(vs, fmt.Sprintf("'%v'", vt))
		default:
			vs = append(vs, fmt.Sprintf("%v", values))
		}
	}
	return vs
}

// Join the strict build-in types values as string like "1,2,3",
// "true,false", "4.5,6.78", "'abc','bce'" and so on, then append 
// the joined string into the query which set by caller.
//
// 1. Usage for number values as:
//
//	- values: []int64{1, 2, 3}
//	- query : "SELECT * FROM tablename WHERE id IN (%s)"
//	- -> The result is "SELECT * FROM tablename WHERE id IN (1,2,3)".
//
// 2. Usage for string values as:
//
//	- values: []string{"123", "abc", "hello"}
//	- query : "SELECT * FROM tablename WHERE tag IN (%s)"
//	- -> The result is "SELECT * FROM tablename WHERE tag IN ('123','abc','hello')".
//
// 3. Usage for none query as:
//
//	- values: []string{"123", "abc", "hello"}
//	- -> The result is "'123','abc','hello'".
//
// Other more usages support by change the values type to defferent one.
//
//	See utils.BuildIn for more types defined informations.
func Joins[T BuildIn] (values []T, query ...string) string {
	vs := ToStrings(values)
	if len(vs) > 0 {
		// Append values into none-empty query string
		if q := Variable(query, ""); q != "" {
			return fmt.Sprintf(q, strings.Join(vs, ","))
		}
		return strings.Join(vs, ",")
	}
	return ""
}

// Join the strict build-in types vlues with sep char to string.
func JoinsSep[T BuildIn](values []T, sep string) string {
	return strings.Join(ToStrings(values), sep)
}

// Join the input string as multiple lines.
func JoinLines(values ...string) string {
	packet := ""
	for _, line := range values {
		packet += line + "\n"
	}
	return packet
}

// Clamp the number in given range.
func Clamp[T cmp.Ordered](num T, minmum T, maximum T) T {
	return min(max(num, minmum), maximum)
}

// Check the given list whether contain the item value.
func Contain[T comparable](list []T, item T) bool {
	for _, v := range list {
		if v == item {
			return true
		}
	}
	return false
}

// Check the given items whether all contains in totals.
func Contains[T any](totals []T, items []T) bool {
	return NewSets[T]().Add(totals...).Contain(items...)
}

// Remove the invalid values from items which not exist in totals.
func Filters[T any](totals []T, items []T) []T {
	if len(totals) == 0 || len(items) == 0 {
		return []T{}
	}
	return NewSets[T]().Add(totals...).Filters(items...)
}

// Check the given items whether exist anyone in totals.
func Exists[T any](totals []T, items []T) bool {
	return NewSets[T]().Add(totals...).Exist(items...)
}

// Remove duplicate items from the given array.
func Distinct[T any](src []T) []T {
	return NewSets[T]().Add(src...).Array()
}

// TrimEmpty remove empty string, it maybe return empty result array.
func TrimEmpty(src []string) []string {
	dst := []string{}
	for _, str := range src {
		if str != "" {
			dst = append(dst, str)
		}
	}
	return dst
}

// ----------------------------------------
// Translate
// ----------------------------------------

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

// ----------------------------------------
// Reversal
// ----------------------------------------

// Reversal ints string to int array with default separator , char or custom separator.
func ReverInts(src string, sep ...string) []int {
	vs := strings.Split(src, Variable(sep, ","))

	out := []int{}
	for _, v := range vs {
		if vi, err := strconv.Atoi(v); err == nil {
			out = append(out, vi)
		}
	}
	return out
}

// Reversal int64s string to int64 array with default separator , char or custom separator.
func ReverInt64s(src string, sep ...string) []int64 {
	vs := strings.Split(src, Variable(sep, ","))

	out := []int64{}
	for _, v := range vs {
		if vi, err := strconv.ParseInt(v, 10, 64); err == nil {
			out = append(out, vi)
		}
	}
	return out
}

// ----------------------------------------
// Split
// ----------------------------------------

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

// ----------------------------------------
// Others
// ----------------------------------------

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
