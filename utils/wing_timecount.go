// Copyright (c) 2018-Now Dunyu All Rights Reserved.
//
// Author      : https://www.wengold.net
// Email       : support@wengold.net
//
// Prismy.No | Date       | Modified by. | Description
// -------------------------------------------------------------------
// 00001       2025/05/28   youhei         New version
// -------------------------------------------------------------------

package utils

import (
	"fmt"
	"time"

	"github.com/wengoldx/xcore/logger"
)

// Time counter for calculate used durations.
type timeCounter struct {
	start int64 // Start time in unix nano second.
	tick  int64 // Tick start time for count tick interval.
}

// Create a new time counter for calculate duration and logs.
func NewTimeCounter() timeCounter {
	return timeCounter{start: time.Now().UnixNano()}
}

// Count the used duration time after counter create or called Reset().
func (c *timeCounter) Count() int64 {
	used := time.Now().UnixNano() - c.start
	return used
}

// Count the tick interval after counter create or called Reset().
func (c *timeCounter) Tick() int64 {
	c.tick = time.Now().UnixNano() - c.tick
	return c.tick
}

// Reset start time and clear used duration value.
func (c *timeCounter) Reset() {
	c.start = time.Now().UnixNano()
	c.tick = c.start
}

// Count and logout the used duration on auto calculate unit.
func (c *timeCounter) LogUsed(msg string, islogger ...bool) {
	used := c.Count()
	if Variable(islogger, false) {
		logger.If("%s: %s", msg, formatDuration(used))
		return
	}
	fmt.Printf("%s: %s\n", msg, formatDuration(used))
}

// Count and logout the used duration on auto calculate unit.
func (c *timeCounter) LogTick(msg string, islogger ...bool) {
	interval := c.Tick()
	if Variable(islogger, false) {
		logger.If("%s: %s", msg, formatDuration(interval))
		return
	}
	fmt.Printf("%s: %s\n", msg, formatDuration(interval))
}

const (
	_tc_s  = int64(time.Second)      // 1000 x 1000 x 1000 ns
	_tc_ms = int64(time.Millisecond) // 1000 x 1000 ns
	_tc_us = int64(time.Microsecond) // 1000 ns
)

// Format used time like: 1.234s, 56.789ms, 123.456us, 234ns.
func formatDuration(used int64) string {
	if used > _tc_s {
		return fmt.Sprintf("%v.%v s", used/_tc_s, (used%_tc_s)/_tc_ms)
	} else if used > _tc_ms {
		return fmt.Sprintf("%v.%v ms", used/_tc_ms, (used%_tc_ms)/_tc_us)
	} else if used > _tc_us {
		return fmt.Sprintf("%v.%v us", used/_tc_us, used%_tc_us)
	} else {
		return fmt.Sprintf("%v ns", used)
	}
}
