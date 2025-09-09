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
	"time"

	"github.com/wengoldx/xcore/invar"
	"github.com/wengoldx/xcore/logger"
)

/* --------------------------- */
/* Custom Task On Queue        */
/* --------------------------- */

// Task monitor to execute queue tasks in sequence.
type QueueTask struct {
	queue       *Queue           // Task object queue.
	postchan    chan EmptyStruct // Block chan for queue task PIPO.
	interrupt   bool             // The flag for interrupt task monitor when case error if set true.
	interval    time.Duration    // The interval between two task to waiting, set 0 for non-waiting.
	authHandler AuthHandler      // The handler for return task unique id for any fetch actions.
}

// Typed function to return target item unique id for fetch actions.
type AuthHandler func(taskdata any) string

// Typed function to configure a QueueTask.
type Option func(*QueueTask)

// Set the interrupt for a queue task.
func WithInterrupt(interrupt bool) Option {
	return func(qt *QueueTask) { qt.interrupt = interrupt }
}

// Set the interval for a queue task.
func WithInterval(interval time.Duration) Option {
	return func(qt *QueueTask) {
		if interval > 0 {
			qt.interval = interval
		}
	}
}

// Set the auth handler function for return target task unique id when fetching.
func WithAuthHandler(handler AuthHandler) Option {
	return func(qt *QueueTask) { qt.authHandler = handler }
}

// A interface for create queue task hanlder to execute task callback.
type TaskHandler interface {

	// Execute queue task with data.
	ExecQueueTask(data any) error
}

// The adapter to allow the use of ordinary functions as TaskHandler object.
// If `f` is a function with the appropriate signature, `TaskHandlerFunc(f)`
// is an `TaskHandler` that calls `f`.
type TaskHandlerFunc func(data any) error

func (e TaskHandlerFunc) ExecQueueTask(data any) error {
	return e(data)
}

// Create a new queue task and start as runtime monitor.
//
// Set custom interval duration and interrupt flag by call utils.WithInterrupt(),
// utils.WithInterval() like follow codes:
//
//	task := utils.NewQueueTask(handler,
//		utils.WithInterrupt(true),                 // interrupt monitor when case error.
//		utils.WithInterval(20 * time.Millisecond), // waiting 500ms between two task.
//	)
//	task.Post(taskdata)
func NewQueueTask(handler TaskHandler, opts ...Option) *QueueTask {
	task := &QueueTask{
		queue:     NewQueue(),
		postchan:  make(chan EmptyStruct),
		interrupt: false, // not interrupt by default.
		interval:  0,     // non-waiting delay.
	}
	for _, optFunc := range opts {
		optFunc(task)
	}

	// start task monitor to listen task insert.
	go task.startTaskMonitor(handler)
	logger.I("Start QueueTask monitor > interrupt:",
		task.interrupt, "interval:", task.interval)
	return task
}

// Set task monitor interrupt filter times.
func (t *QueueTask) SetInterrupt(interrupt bool) {
	t.interrupt = interrupt
}

// Set waiting interval between tasks, the value must >= 0.
func (t *QueueTask) SetInterval(interval time.Duration) {
	if interval >= 0 {
		t.interval = interval
	}
}

// Push a new task into monitor at queue back.
func (t *QueueTask) Post(taskdata any, maxlimits ...int) error {
	if taskdata == nil {
		return invar.ErrInvalidData
	}

	if ml := Variable(maxlimits, 0); ml > 0 && t.queue.Len() > ml {
		logger.E("Task queue too heavy on oversize", ml)
		return invar.ErrPoolFull
	}

	t.queue.Push(taskdata)
	//  NOTICE: The task handler will called as PIPO by using
	//  the 'postchan' to blocking in gorutine, so it no-need
	//  to check whether the handlers method execute toggether
	//  when multiple post requst comming.
	go func() { t.postchan <- E_ }()
	return nil
}

// Cancels the waiting tasks, use AuthHandler callback return task
// unique id to find target cancal task on fetching.
//
//	queuetask.Cancels(func(taskdata any) string {
//		if item, ok := taskdata.(*QueueItem); ok {
//			return item.TaskID
//		}
//		return ""
//	}, cancelid)
//
// Or, set the AuthHandler as option by utils.WithAuthHandler() when
// create QueueTask instance, then call CancelsHandler(cancelids...)
// to cancel target tasks.
func (t *QueueTask) Cancels(handler AuthHandler, tags ...string) {
	if len(tags) > 0 {
		ids := NewSets[string]().Add(tags...)
		cnt := ids.Size()

		t.queue.Fetch(func(d any) Result {
			if id := handler(d); id != "" && ids.Contain(id) {
				rst := Condition(cnt == 1, REMOVE_INTERUPT, REMOVE_CONTINUE)
				logger.I("Canceled task:", id) // "- result", rst, "ids:", ids.Array(), "cnt:", cnt)

				ids.Remove(id) // remove target found item id.
				cnt--          // decrease the cancel ids count.
				return rst
			}
			return KEEP_FETCHING
		})
	}
}

// Cancels the waiting tasks by the inited AuthHandler.
func (t *QueueTask) CancelsHandler(tags ...string) {
	if t.authHandler != nil {
		t.Cancels(t.authHandler, tags...)
	}
}

// Start task monitor to listen tasks pushed into queue, and execute it.
func (t *QueueTask) startTaskMonitor(handler TaskHandler) {
	if handler == nil {
		logger.E("Nil handler, exit monitor!")
		return
	}

	for {
		<-t.postchan // blocking and waiting task post.

		// popup the topmost task to execte.
		taskdata, err := t.queue.Pop()
		if err != nil {
			break // queue maybe empty.
		}

		if err := handler.ExecQueueTask(taskdata); err != nil {
			logger.E("Execute queue task, err:", err)
			if t.interrupt {
				logger.I("Interrupted QueueTask monitor!")
				return
			}
		}

		// waiting for next if need.
		if t.interval > 0 {
			time.Sleep(t.interval)
		}
	}
}
