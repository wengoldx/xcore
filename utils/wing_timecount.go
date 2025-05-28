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
	used  int64 // Used durations of unix nano second.
}

// Create a new time counter for calculate duration and logs.
func NewTimeCounter() timeCounter {
	return timeCounter{
		start: time.Now().UnixNano(),
	}
}

// Count the used duration time after counter create or called Reset().
func (c *timeCounter) Count() int64 {
	c.used = time.Now().UnixNano() - c.start
	return c.used
}

// Reset start time and clear used duration value.
func (c *timeCounter) Reset() {
	c.start, c.used = time.Now().UnixNano(), 0
}

// Count and logout the used duration on auto calculate unit.
func (c *timeCounter) LogUsed(msg string, custom ...bool) {
	used := c.Count()
	if len(custom) > 0 && custom[0] {
		logger.If("%s: %s", msg, formatDuration(used))
	} else {
		fmt.Printf("%s: %s\n", msg, formatDuration(used))
	}
}

// Format used time like: 1.23.345678s, 90.123456ms, 78.901us, 234ns.
func formatDuration(used int64) string {
	if used > int64(time.Second) {
		return fmt.Sprintf("%v.%v.%vs", used/int64(time.Second),
			(used%int64(time.Second))/int64(time.Millisecond),
			used%int64(time.Millisecond))
	} else if used > time.Hour.Milliseconds() {
		return fmt.Sprintf("%v.%vms", used/int64(time.Millisecond),
			used%int64(time.Millisecond))
	} else if used > time.Hour.Microseconds() {
		return fmt.Sprintf("%v.%vus", used/int64(time.Microsecond),
			used%int64(time.Microsecond))
	} else {
		return fmt.Sprintf("%vns", used)
	}
}
