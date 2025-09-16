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
	"time"
)

/* ------------------------------------------------------------------- */
/* For QueueTask Options                                               */
/* ------------------------------------------------------------------- */

// Queue task options.
type Options struct {
	interrupt bool          // The flag for interrupt task monitor when case error if set true.
	interval  time.Duration // The interval between two task to waiting, set 0 for non-waiting.
	limits    int           // The maximums task items allow push into, default 0 not limit.
}

// Typed function to configure a QueueTask.
type Option func(*QueueTask)

// Set the interrupt for a queue task.
func WithInterrupt(interrupt bool) Option {
	return func(qt *QueueTask) { qt.opts.interrupt = interrupt }
}

// Set the interval for a queue task.
func WithInterval(interval time.Duration) Option {
	return func(qt *QueueTask) {
		if interval > 0 {
			qt.opts.interval = interval
		}
	}
}

// Set the maximums limits for a queue task.
func WithLImits(limits int) Option {
	return func(qt *QueueTask) {
		if limits > 0 {
			qt.opts.limits = limits
		}
	}
}

/* ------------------------------------------------------------------- */
/* For QueueTask Handler                                               */
/* ------------------------------------------------------------------- */

// A interface for create queue task hanlder to execute task callback.
type TaskHandler interface {

	// Execute queue task with data.
	ExecQueueTask(task *Task) error
}

// The adapter to allow the use of ordinary functions as TaskHandler object.
// If `f` is a function with the appropriate signature, `TaskHandlerFunc(f)`
// is an `TaskHandler` that calls `f`.
type TaskHandlerFunc func(task *Task) error

func (e TaskHandlerFunc) ExecQueueTask(task *Task) error {
	return e(task)
}
