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

package tasks

import (
	"math/rand/v2"
	"strconv"
	"testing"
	"time"

	"github.com/wengoldx/xcore/logger"
)

/* ------------------------------------------------------------------- */
/* For Queur Task Tests                                                */
/* ------------------------------------------------------------------- */

func TestQueueTask(t *testing.T) {
	handler := TaskHandlerFunc(ExecCallback)
	qtask := NewQueueTask(handler, WithInterrupt(false), WithInterval(25*time.Millisecond))

	logger.I("@@ Prepare cancel ids.")
	cids := []string{}
	for i := 1; i <= 10; i++ {
		cids = append(cids, strconv.Itoa(rand.IntN(45)+5))
	}

	logger.I("@@ Post 50 test tasks...")
	for i := 1; i <= 50; i++ {
		logger.I("Post task:", i)
		qtask.Post(NewTask(strconv.Itoa(i), i))
	}

	// waiting for executing 1 ~ 2 tasks.
	time.Sleep(100 * time.Millisecond)

	logger.I("@@ Request cancels...")
	go qtask.Cancels(cids...)

	logger.I("@@ Continue post 10 test tasks...")
	qtask.SetInterval(0)
	for i := 51; i <= 60; i++ {
		logger.I("Post task:", i)
		qtask.Post(NewTask(strconv.Itoa(i), i))
	}

	logger.I("@@ Waiting 1 seconds")
	time.Sleep(1 * time.Second)

	logger.I("@@ Call exit and wait 1 second.")
	qtask.Exit()
	time.Sleep(1 * time.Second)
	logger.I("@@ Finish Test!")
}

func ExecCallback(data *Task) error {
	time.Sleep(25 * time.Millisecond)
	logger.I(" - Executed task:", data.ID)
	return nil
}

func TestSwitchTask(t *testing.T) {
	handler := TaskHandlerFunc(ExecCallback)
	qtask := NewQueueTask(handler, WithInterrupt(false))

	logger.I("@@ Random from and to ids.")
	from := strconv.Itoa(rand.IntN(24) + 10)
	to := strconv.Itoa(rand.IntN(24) + 36)

	logger.I("@@ Post 60 test tasks...")
	for i := 1; i <= 60; i++ {
		logger.I("Post task:", i)
		qtask.Post(NewTask(strconv.Itoa(i), i))
	}

	logger.I("@@ Swith", from, ">", to)
	if !qtask.Switch(from, to) {
		logger.I("@@ Switch failed!")
	}

	// waiting for executing 1 second.
	time.Sleep(2 * time.Second)
	qtask.Exit()
	logger.I("@@ Finish Test!")
}
