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
	"math/rand/v2"
	"strconv"
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

	// prepare cancel ids.
	logger.I("@@ Prepare cancel ids.")
	cids := []string{}
	for i := 1; i <= 15; i++ {
		cids = append(cids, strconv.Itoa(rand.IntN(30)+21))
	}

	// post 50 task
	logger.I("@@ Post 50 test tasks...")
	for i := 1; i <= 50; i++ {
		logger.I("Post task:", i)
		qtask.Post(strconv.Itoa(i))
	}

	// waiting for executing 1 ~ 20 tasks.
	time.Sleep(1 * time.Second)

	logger.I("@@ Request cancels...")
	go qtask.Cancels(func(taskdata any) string {
		if cid, ok := taskdata.(string); ok {
			// logger.I("Get task item:", cid)
			return cid
		}

		logger.E("!! INVALID TASK DATA !!")
		return ""
	}, cids...)

	logger.I("@@ Post 10 test tasks...")
	qtask.SetInterval(0)
	for i := 51; i <= 60; i++ {
		logger.I("Post task:", i)
		qtask.Post(strconv.Itoa(i))
	}

	logger.I("@@ Waiting 3 seconds")
	time.Sleep(3 * time.Second)

	logger.I("@@ Call exit and wait 1 second.")
	qtask.Exit()
	time.Sleep(1 * time.Second)
	logger.I("@@ Finish Test!")
}

func ExecCallback(data any) error {
	index := data.(string)

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
		NewCase("Check 1", "1\\2\\4\\5\\6", "  /  1//2\\3/..///4/./5/6\\\\"),
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

func TestSplitSuffix(t *testing.T) {
	cases := []*TestCase{
		NewCase("Check 1", "123", "/1/2/   123  .pdf"),
		NewCase("Check 1", "123", "123.pdf"),
		NewCase("Check 2", "123", "123"),
		NewCase("Check 3", "", ".pdf"),
		NewCase("Check 4", "", ""),
	}

	for _, c := range cases {
		rst, suffix := SplitSuffix(c.Params.(string))
		if want := c.Want.(string); rst != want {
			t.Fatal("FileBaseName error > want:", want, "but result is", rst)
		}
		fmt.Println("suffix:", suffix)
	}
}
