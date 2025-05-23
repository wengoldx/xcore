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
	"testing"
	"time"

	"github.com/wengoldx/xcore/logger"
)

func TestQueueTask(t *testing.T) {
	handler := TaskHandlerFunc(ExecCallback)
	qtask := NewQueueTask(handler, WithInterrupt(false), WithInterval(25*time.Millisecond))
	for i := 0; i < 50; i++ {
		logger.I("Post task:", i)
		qtask.Post(i)
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
