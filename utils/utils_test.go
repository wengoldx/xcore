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

package utils

import (
	"fmt"
	"testing"

	wt "github.com/wengoldx/xcore/utils/testx"
)

/* ------------------------------------------------------------------- */
/* For File Utils Tests                                                */
/* ------------------------------------------------------------------- */

func TestNormalizePath(t *testing.T) {
	// FIXME : for windows system want string.
	cases := []*wt.TestCase{
		wt.NewCase("Check 1", "1\\2\\4\\5\\6", "  /  1//2\\3/..///4/./5/6\\\\"),
		wt.NewCase("Check 2", "1\\2\\3", "    1/2//3/     "),
		wt.NewCase("Check 3", "1 \\2\\3", "/  1 /2\\3\\    "),
		wt.NewCase("Check 4", ".", ""),
	}

	wt.TestMults(t, cases, func(param any) any {
		return NormalizePath(param.(string))
	})
}

func TestSplitSuffix(t *testing.T) {
	cases := []*wt.TestCase{
		wt.NewCase("Check 1", "123", "/1/2/   123  .pdf"),
		wt.NewCase("Check 1", "123", "123.pdf"),
		wt.NewCase("Check 2", "123", "123"),
		wt.NewCase("Check 3", "", ".pdf"),
		wt.NewCase("Check 4", "", ""),
	}

	wt.TestMults(t, cases, func(param any) any {
		rst, suffix := SplitSuffix(param.(string))
		fmt.Println("suffix:", suffix)
		return rst
	})
}
