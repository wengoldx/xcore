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
	"container/list"
	"sync"
	"time"

	"github.com/wengoldx/xcore/invar"
	"github.com/wengoldx/xcore/logger"
	"github.com/wengoldx/xcore/utils"
)

/* ------------------------------------------------------------------- */
/* For Task Data                                                       */
/* ------------------------------------------------------------------- */

// Task data with unique id.
type Task struct {
	ID   string // Task unique id.
	Data any    // Task custom datas.
}

// Create a task for Push to queue.
func NewTask(id string, data any) *Task {
	return &Task{id, data}
}

/* ------------------------------------------------------------------- */
/* For QueueTask                                                       */
/* ------------------------------------------------------------------- */

// Task monitor to execute queue tasks in sequence.
type QueueTask struct {
	queue *list.List // Tasks queue.
	mutex sync.Mutex

	postchan chan utils.TNone // Block chan for queue task PIPO.
	exitchan chan utils.TNone // Block chan for exit queue task monitor.
	opts     Options          // Queue task options.
}

// Create a new queue task and start as runtime monitor.
//
// Set custom interval duration and interrupt flag by call utils.WithInterrupt(),
// utils.WithInterval() like follow codes:
//
//	task := utils.NewQueueTask(handler,
//		utils.WithInterrupt(true),                 // interrupt monitor when case error.
//		utils.WithInterval(20 * time.Millisecond), // waiting 20ms between two task.
//	)
//	task.Post(taskdata)
func NewQueueTask(handler TaskHandler, opts ...Option) *QueueTask {
	qt := &QueueTask{
		queue:    list.New(),
		postchan: make(chan utils.TNone),
		exitchan: make(chan utils.TNone),
		opts: Options{
			interrupt: false, // not interrupt by default.
			interval:  0,     // non-waiting delay.
			limits:    0,     // non-limit.
		},
	}

	// init task monitor options.
	for _, optFunc := range opts {
		optFunc(qt)
	}

	// start task monitor to listen task insert.
	go qt.startMonitor(handler)
	it, iv := qt.opts.interrupt, qt.opts.interval
	logger.I("Start QueueTask monitor > interrupt:", it, "interval:", iv)
	return qt
}

// Set task monitor interrupt filter times.
func (t *QueueTask) SetInterrupt(interrupt bool) {
	t.opts.interrupt = interrupt
}

// Set waiting interval between tasks, the value must >= 0.
func (t *QueueTask) SetInterval(interval time.Duration) {
	if interval >= 0 {
		t.opts.interval = interval
	}
}

// Set maximums task items limit, the value must >= 0, default 0 non-limit.
func (t *QueueTask) SetLimits(limits int) {
	if limits >= 0 {
		t.opts.limits = limits
	}
}

// Return quenu item counts.
func (t *QueueTask) Counts() int {
	return t.queue.Len()
}

// Push a new task to monitor at queue backend.
func (t *QueueTask) Post(task *Task) error {
	if task == nil || task.ID == "" || task.Data == nil {
		return invar.ErrInvalidData
	}

	if limits := t.opts.limits; limits > 0 && t.queue.Len() > limits {
		logger.E("QueueTask too heavy on oversize", limits)
		return invar.ErrPoolFull
	}

	/*
	 * FIXME: The task handler will called as PIPO by using chan
	 * 'postchan' to blocking in gorutine, so it no-need to check
	 * whether handler methods executing toggether when multiple
	 * post requests comming.
	 */
	t.push(task)
	go func() { t.postchan <- utils.NONE }()
	return nil
}

// Cancels the target waiting tasks by id, it will fetch all tasks.
func (t *QueueTask) Cancels(tags ...string) []string {
	if len(tags) > 0 {
		return t.removes(tags...)
	}
	return tags
}

// Switch the target two tasks by given ids when both exist.
func (t *QueueTask) Switch(from, to string) bool {
	if from != "" && to != "" && from != to {
		return t.switchs(from, to)
	}
	return false
}

// Clear waiting tasks and exit the QueueTask monitor.
func (t *QueueTask) Exit() {
	t.clear() // clear tasks first.
	go func() { t.exitchan <- utils.NONE }()
}

// Start task monitor to listen tasks pushed into queue, and execute it.
func (t *QueueTask) startMonitor(handler TaskHandler) {
	if handler == nil {
		logger.E("Nil handler, exit monitor!")
		return
	}

	for {
		select {
		case <-t.exitchan: // stop task monitor.
			logger.I("Exist QueueTask monitor!")
			return

		case <-t.postchan: // blocking and waiting task post.
			// popup the topmost task to execte.
			task, err := t.pop()
			if err != nil {
				/*
				 * FIXME: Any tasks maybe removed when user handled cancel
				 * action, but the postchan requests not removed toggether,
				 * so HERE MUST filter out the invalid request chans when
				 * QueueTask is empty!
				 */
				continue // queue maybe empty.
			}

			if err := handler.ExecQueueTask(task); err != nil {
				logger.E("Execute task, err:", err)
				if t.opts.interrupt {
					logger.I("Interrupted QueueTask monitor!")
					return
				}
			}

			// waiting for next if need.
			if t.opts.interval > 0 {
				time.Sleep(t.opts.interval)
			}
		}
	}
}

/* ------------------------------------------------------------------- */
/* For QueueTask Internal Methods                                      */
/* ------------------------------------------------------------------- */

// Push a task to queue back.
func (t *QueueTask) push(task *Task) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	t.queue.PushBack(task)
}

// Pick and remove the front task of queue,
// it will return invar.ErrEmptyData error when queue is empty.
func (t *QueueTask) pop() (*Task, error) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	if e := t.queue.Front(); e != nil {
		t.queue.Remove(e)
		return e.Value.(*Task), nil
	}
	return nil, invar.ErrEmptyData
}

// Clear all queue tasks.
func (t *QueueTask) clear() {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	for e := t.queue.Front(); e != nil; {
		next := e.Next()
		t.queue.Remove(e)
		e = next
	}
}

// Remove tasks by given ids.
func (t *QueueTask) removes(tags ...string) []string {
	ids := utils.NewSets[string]().Add(tags...)
	cnt := ids.Size()
	logger.I("Fetching and cancel tasks:", ids.Array())

	t.mutex.Lock()
	defer t.mutex.Unlock()

	for e := t.queue.Front(); e != nil && cnt > 0; {
		next := e.Next()
		task, ok := e.Value.(*Task)
		if ok && task != nil && ids.Contain(task.ID) {
			t.queue.Remove(e)

			logger.I("> Canceled task:", task.ID)
			ids.Remove(task.ID) // remove target found item id.
			cnt--               // decrease the cancel ids count.
		}
		e = next
	}
	return ids.Array()
}

// Swtich the given element values when found.
func (t *QueueTask) switchs(from, to string) bool {
	var fe *list.Element
	var te *list.Element

	t.mutex.Lock()
	defer t.mutex.Unlock()

	for e := t.queue.Front(); e != nil; e = e.Next() {
		task, ok := e.Value.(*Task)
		if ok && task != nil {
			switch task.ID {
			case from:
				fe = e
			case to:
				te = e
			}
		}

		// switch elements values when found.
		if fe != nil && te != nil {
			fe.Value, te.Value = te.Value, fe.Value
			logger.I("Switched tasks", from, "<->", to)
			return true
		}
	}
	return false
}
