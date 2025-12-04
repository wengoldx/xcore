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

package xtask

import (
	"math/rand/v2"
	"runtime"
	"strconv"
	"testing"
	"time"

	"github.com/wengoldx/xcore/logger"
)

/* ------------------------------------------------------------------- */
/* For Queur Task Tests                                                */
/* ------------------------------------------------------------------- */

func TestQueueTask(t *testing.T) {
	handler := TaskHandlerFunc(execTaskWorks)
	qtask := NewQueueTask(handler, WithInterrupt(false), WithInterval(10*time.Millisecond))
	start := runtime.NumGoroutine()

	// create 10 cancel task ids.
	cids, cnt := []string{}, 150
	for i := 0; i < 10; i++ {
		cids = append(cids, strconv.Itoa(rand.IntN(90)+10))
	}

	logger.I("Post", cnt, "Tasks (", start, ")...")
	for i := 1; i <= cnt; i++ {
		// logger.I("|- Post Task:", i)
		qtask.Post(NewTask(strconv.Itoa(i), i))
	}
	logger.I("- Posted", cnt, "Tasks | blocking rountines:", runtime.NumGoroutine()-start)

	// waiting for executing 1 ~ 10 tasks.
	time.Sleep(250 * time.Millisecond)
	logger.I("Post Cancels...")
	go qtask.Cancels(cids...)

	qtask.SetInterval(0)
	logger.I("Continue Post", cnt, "Tasks...")
	for i, max := cnt+1, 2*cnt+1; i <= max; i++ {
		// logger.I("|- Post Task:", i)
		qtask.Post(NewTask(strconv.Itoa(i), i))
	}

	logger.I("- Posted", cnt, "Tasks | blocking rountines:", runtime.NumGoroutine()-start)
	logger.I("Waiting 5 seconds for exit...")
	time.Sleep(5 * time.Second)
	qtask.Exit()

	// test post error when monitor closed.
	qtask.Post(NewTask(strconv.Itoa(500), 500))
	logger.I("Finished, rountines:", runtime.NumGoroutine())
}

func execTaskWorks(data *Task) error {
	time.Sleep(15 * time.Millisecond)
	logger.I("| - Executed task:", data.ID)
	return nil
}

func TestSwitchTask(t *testing.T) {
	handler := TaskHandlerFunc(execTaskWorks)
	qtask := NewQueueTask(handler, WithInterrupt(false))

	from := strconv.Itoa(rand.IntN(24) + 10)
	to := strconv.Itoa(rand.IntN(24) + 36)
	logger.I("Random from:", from, "and to:", to, "ids")

	logger.I("Post 60 test tasks...")
	for i := 1; i <= 60; i++ {
		// logger.I("| - Post Task:", i)
		qtask.Post(NewTask(strconv.Itoa(i), i))
	}
	time.Sleep(20 * time.Millisecond)

	logger.I("Swith", from, ">", to)
	if !qtask.Switch(from, to) {
		logger.I("? - Switch failed!")
	}

	qtask.Exit()
	time.Sleep(100 * time.Millisecond)
	logger.I("Finished, rountines:", runtime.NumGoroutine())
}
