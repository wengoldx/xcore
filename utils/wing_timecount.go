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

const (
	_second      = int64(time.Second)      // 1000 x 1000 x 1000 ns
	_millisecond = int64(time.Millisecond) // 1000 x 1000 ns
	_microsecond = int64(time.Microsecond) // 1000 ns
)

// Format used time like: 1.234s, 56.789ms, 123.456us, 234ns.
func formatDuration(used int64) string {
	if used > _second {
		return fmt.Sprintf("%v.%v s", used/_second, (used%_second)/_millisecond)
	} else if used > _millisecond {
		return fmt.Sprintf("%v.%v ms", used/_millisecond, (used%_millisecond)/_microsecond)
	} else if used > _microsecond {
		return fmt.Sprintf("%v.%v us", used/_microsecond, used%_microsecond)
	} else {
		return fmt.Sprintf("%v ns", used)
	}
}
