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

// ----------------------------------------
// Variable
// ----------------------------------------

// Return the first element string, or default if empty.
func VarString(src []string, def string) string {
	if len(src) > 0 && src[0] != "" {
		return src[0]
	}
	return def
}

// Return the first element int (> 0), or default if empty.
func VarInt(src []int, def int) int {
	if len(src) > 0 && src[0] > 0 {
		return src[0]
	}
	return def
}

// Return the first element int64 (> 0), or default if empty.
func VarInt64(src []int64, def int64) int64 {
	if len(src) > 0 && src[0] > 0 {
		return src[0]
	}
	return def
}

// Return the first element bool, or default if empty.
func VarBool(src []bool, def bool) bool {
	if len(src) > 0 && src[0] {
		return src[0]
	}
	return def
}

// ----------------------------------------
// Contain
// ----------------------------------------

// Contain check the given string list if contains item.
//
// You should call Distinct() to filter out the repeat items.
func Contain(list []string, item string) bool {
	for _, v := range list {
		if v == item {
			return true
		}
	}
	return false
}

// Contain check the given int list if contains item.
func ContainInt(list []int, item int) bool {
	for _, v := range list {
		if v == item {
			return true
		}
	}
	return false
}

// Contain check the given int64 list if contains item.
func ContainInt64(list []int64, item int64) bool {
	for _, v := range list {
		if v == item {
			return true
		}
	}
	return false
}

// Contains check the given strings items if contains in totals.
func Contains(totals []string, items []string) bool {
	return NewSets().AddStrings(totals).ContainStrings(items)
}

// ContainInts check the given int items if contains in totals.
func ContainInts(totals []int, items []int) bool {
	return NewSets().AddInts(totals).ContainInts(items)
}

// ContainInt64s check the given int64 items if contains in totals.
func ContainInt64s(totals []int64, items []int64) bool {
	return NewSets().AddInt64s(totals).ContainInt64s(items)
}

// ----------------------------------------
// Fileter
// ----------------------------------------

// Filters remove the strings from items which not exist in totals.
func Filters(totals []string, items []string) []string {
	if len(totals) == 0 || len(items) == 0 {
		return []string{}
	}
	return NewSets().AddStrings(totals).FilterStrings(items)
}

// FilterInts remove the int values from items which not exist in totals.
func FilterInts(totals []int, items []int) []int {
	if len(totals) == 0 || len(items) == 0 {
		return []int{}
	}
	return NewSets().AddInts(totals).FilterInts(items)
}

// FilterInt64s remove the int64 values from items which not exist in totals.
func FilterInt64s(totals []int64, items []int64) []int64 {
	if len(totals) == 0 || len(items) == 0 {
		return []int64{}
	}
	return NewSets().AddInt64s(totals).FilterInt64s(items)
}

// ----------------------------------------
// Exist
// ----------------------------------------

// ExistInts check the given int items if any exist in totals.
func ExistInts(totals []int, items []int) bool {
	return NewSets().AddInts(totals).ExistInts(items)
}

// ExistInt64s check the given int64 items if any exist in totals.
func ExistInt64s(totals []int64, items []int64) bool {
	return NewSets().AddInt64s(totals).ExistInt64s(items)
}

// ExistStrings check the given string items if any exist in totals.
func ExistStrings(totals []string, items []string) bool {
	return NewSets().AddStrings(totals).ExistStrings(items)
}

// ----------------------------------------
// Distinct
// ----------------------------------------

// Distinct remove duplicate string from given array.
//
// You should call Contain() only for check if exist sub string.
func Distinct(src []string) []string {
	return NewSets().AddStrings(src).ArrayString()
}

// Distinct remove duplicate int from given array.
func DistInts(src []int) []int {
	return NewSets().AddInts(src).ArrayInt()
}

// Distinct remove duplicate int64 from given array.
func DistInt64s(src []int64) []int64 {
	return NewSets().AddInt64s(src).ArrayInt64()
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
// Join
// ----------------------------------------

// JoinInts join int numbers as string '1,2,3' with ',' default separator,
// or custom separator '-' like '1-2-3'.
func JoinInts(nums []int, sep ...string) string {
	vs := []string{}
	if len(nums) > 0 {
		for _, num := range nums {
			if v := strconv.Itoa(num); v != "" {
				vs = append(vs, v)
			}
		}
	}
	return strings.Join(vs, VarString(sep, ","))
}

// JoinInt64s join int64 numbers as string '1,2,3' with ',' default separator,
// or custom separator '-' like '1-2-3'.
func JoinInt64s(nums []int64, sep ...string) string {
	vs := []string{}
	if len(nums) > 0 {
		for _, num := range nums {
			if v := strconv.FormatInt(num, 10); v != "" {
				vs = append(vs, v)
			}
		}
	}
	return strings.Join(vs, VarString(sep, ","))
}

// JoinFloats join float64 numbers as string '0.1,2,3.45' with ',' default separator,
// or custom separator '-' like '0.1-2-3.45'.
func JoinFloats(nums []float64, sep ...string) string {
	vs := []string{}
	if len(nums) > 0 {
		for _, num := range nums {
			if v := strconv.FormatFloat(num, 'f', 2, 64); v != "" {
				vs = append(vs, v)
			}
		}
	}
	return strings.Join(vs, VarString(sep, ","))
}

// JoinLines combine strings into multiple lines
func JoinLines(inputs ...string) string {
	packet := ""
	for _, line := range inputs {
		packet += line + "\n"
	}
	return packet
}

// ----------------------------------------
// Reversal
// ----------------------------------------

// Reversal ints string to int array with default separator , char or custom separator.
func ReverInts(src string, sep ...string) []int {
	vs := strings.Split(src, VarString(sep, ","))

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
	vs := strings.Split(src, VarString(sep, ","))

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
