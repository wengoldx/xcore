// Copyright (c) 2018-Now Dunyu All Rights Reserved.
//
// Author      : https://www.wengold.net
// Email       : support@wengold.net
//
// Prismy.No | Date       | Modified by. | Description
// -------------------------------------------------------------------
// 00001       2019/05/22   yangping       New version
// 00002       2022/03/26   yangping       Using toolbox.Task
// -------------------------------------------------------------------

package wt

import "testing"

// Test case datas for multiple testing.
type TestCase struct {
	Case   string
	Params any
	Want   any
}

// Test handler for function test.
type TestHandler func(param any) any

// Return test case object for easy multipe testing.
func NewCase(label string, want any, param any) *TestCase {
	return &TestCase{Case: label, Want: want, Params: param}
}

// Test helper to execute multiple cases and check the result.
//
// Example:
//
//	cases := []*wt.TestCase{
//		wt.NewCase("Check 1", "1 \\2\\3", "/  1 /2\\3\\    "),
//		wt.NewCase("Check 2", ".", ""),
//	}
//
//	wt.TestMults(t, cases, func(param any) any {
//		return NormalizePath(param.(string))
//	})
func TestMults(t *testing.T, cases []*TestCase, callback TestHandler) {
	t.Helper()
	for _, c := range cases {
		if want := callback(c.Params); want != c.Want {
			t.Fatal("Failed, want:", c.Want, "but result", want)
		}
	}
}
