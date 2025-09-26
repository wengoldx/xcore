// Copyright (c) 2018-Now Dunyu All Rights Reserved.
//
// Author      : https://www.wengold.net
// Email       : support@wengold.net
//
// Prismy.No | Date       | Modified by. | Description
// -------------------------------------------------------------------
// 00001       2025/07/01   yangping       New version
// -------------------------------------------------------------------

package dfa

import (
	"fmt"
	"testing"
)

// -------------------------------------------------------
// USAGE: Enter current folder and excute command to test.
//
//	go test -v -cover
// -------------------------------------------------------

func TestDFAMatcher(t *testing.T) {
	words := []string{"测试", "太阳"}
	text := "我正在执行测试，结果一定要准，不然看不到明天的太阳了"
	matcher := NewMatcher(WithWords(words))

	fmt.Println("Match Words:")
	sensitives := matcher.MatchWords(text)
	for _, s := range sensitives {
		fmt.Println("- Sensitive:", s)
	}

	fmt.Println("Match Replace:")
	_, replace := matcher.MatchReplace(text)
	fmt.Println("- Replace:", replace)
}
