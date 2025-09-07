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
	"math/rand/v2"
	"testing"
	"time"

	"github.com/wengoldx/xcore/logger"
)

// Test case datas for multiple testing.
type TestCase struct {
	Case   string
	Params any
	Want   any
}

// Return test case object for easy multipe testing.
func NewCase(label string, want any, param any) *TestCase {
	return &TestCase{Case: label, Want: want, Params: param}
}

/* ------------------------------------------------------------------- */
/* For Queur Task Tests                                                */
/* ------------------------------------------------------------------- */

func TestQueueTask(t *testing.T) {
	handler := TaskHandlerFunc(ExecCallback)
	qtask := NewQueueTask(handler, WithInterrupt(false), WithInterval(25*time.Millisecond))
	for i := 0; i < 50; i++ {
		logger.I("Post task:", i)
		qtask.Post(i)
	}

	for i := 0; i < 10; i++ {
		cid := rand.IntN(40) + 10
		logger.I("Request cancel:", cid)
		qtask.Cancels(func(taskdata any) (bool, bool) {
			if cid == taskdata.(int) {
				logger.I("- Canceled task:", cid)
				return true, true
			}
			return false, false
		})
	}

	time.Sleep(1 * time.Second)

	qtask.SetInterval(0)
	for i := 50; i < 70; i++ {
		logger.I("Post task:", i)
		qtask.Post(i)
	}
	time.Sleep(5 * time.Second)
}

func ExecCallback(data any) error {
	index := data.(int)

	start := time.Now().UnixNano()
	time.Sleep(25 * time.Millisecond)
	used := (time.Now().UnixNano() - start) / int64(time.Millisecond)
	logger.I(" - Executed task:", index, "used time:", used)
	return nil
}

/* ------------------------------------------------------------------- */
/* For File Utils Tests                                                */
/* ------------------------------------------------------------------- */

func TestNormalizePath(t *testing.T) {
	// FIXME : for windows system want string.
	cases := []*TestCase{
		NewCase("Check 1", "1\\2\\4\\5", "1/2//3/../4/./5/"),
		NewCase("Check 2", "1\\2\\3", "    1/2//3/     "),
		NewCase("Check 3", "1 \\2\\3", "/  1 /2\\3\\    "),
		NewCase("Check 4", ".", ""),
	}

	for _, c := range cases {
		rst := NormalizePath(c.Params.(string))
		if want := c.Want.(string); rst != want {
			t.Fatal("NormalizePath error > want:", want, "but result is", rst)
		}
	}
}

func TestFileBaseName(t *testing.T) {
	cases := []*TestCase{
		NewCase("Check 1", "123", "   123  .pdf"),
		NewCase("Check 1", "123", "123.pdf"),
		NewCase("Check 2", "123", "123"),
		NewCase("Check 3", "", ".pdf"),
		NewCase("Check 4", "", ""),
	}

	for _, c := range cases {
		rst := FileBaseName(c.Params.(string))
		if want := c.Want.(string); rst != want {
			t.Fatal("FileBaseName error > want:", want, "but result is", rst)
		}
	}
}
