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

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/toolbox"
	"github.com/wengoldx/xcore/invar"
	"github.com/wengoldx/xcore/logger"
)

/* --------------------------- */
/* Timer Task Base On toolbox  */
/* --------------------------- */

// Task datas for multipe generate
type WTask struct {
	Name    string           // monitor task name
	Func    toolbox.TaskFunc // monitor task execute function
	Spec    string           // monitor task interval
	ForProd bool             // indicate the task only for prod mode, default no limit
}

// Add a single monitor task to list
func AddTask(tname, spec string, f toolbox.TaskFunc) {
	monitor := toolbox.NewTask(tname, spec, f)
	monitor.ErrLimit = 0

	logger.I("Create task:", tname, "and add to list")
	toolbox.AddTask(tname, monitor)
}

// Generate tasks and start them as monitors.
func StartTasks(monitors []*WTask) {
	for _, m := range monitors {
		if m.ForProd && beego.BConfig.RunMode != "prod" {
			logger.W("Filter out task:", m.Name, "on dev mode")
			continue
		}
		AddTask(m.Name, m.Spec, m.Func)
	}

	toolbox.StartTask()
	logger.I("Started all monitors")
}

// Return task if exist, or nil when unexist
func GetTask(tname string) *toolbox.Task {
	if tasker, ok := toolbox.AdminTaskList[tname]; ok {
		return tasker.(*toolbox.Task)
	}
	return nil
}

/* --------------------------- */
/* Custom Task On Queue        */
/* --------------------------- */

// Task monitor to execute queue tasks in sequence
type QTask struct {
	queue     *Queue
	interrupt bool
	interval  time.Duration
	executing bool
}

// Block chan for TTack queue PIPO
var ttaskchan = make(chan string)

// TaskCallback task callback function
type TaskCallback func(data any) error

// Generat a new task monitor instance.
//
// Custom interval duration and interrupt flag by input params as follow:
//
//	interrupt := 1  // interrupt to execut the remain tasks when case error
//	interval := 500 // sleep interval between tasks in microseconds
//	task := comm.GenTask(callback, interrupt, interval)
//	task.Post(taskdata)
func GenQTask(callback TaskCallback, configs ...int) *QTask {
	task := &QTask{
		queue: GenQueue(), interrupt: false, interval: 0, executing: false,
	}

	// set task configs from given data
	if configs != nil {
		task.interrupt = len(configs) > 0 && configs[0] > 0
		if len(configs) > 1 && configs[1] > 0 {
			task.interval = time.Duration(configs[1] * 1000)
		}
	}

	// start task monitor to listen task insert
	go task.startTaskMonitor(callback)
	logger.I("Excuting task monitor:{interrupt:", task.interrupt, ", interval:", task.interval, "}")
	return task
}

// Set task monitor interrupt filter times
func (t *QTask) SetInterrupt(interrupt bool) {
	t.interrupt = interrupt
}

// Set waiting interval between tasks in microseconds, and it must > 0.
func (t *QTask) SetInterval(interval int) {
	if interval > 0 {
		t.interval = time.Duration(interval * 1000)
	}
}

// Push a new task to monitor queue back
func (t *QTask) Post(taskdata any, maxlimits ...int) error {
	if taskdata == nil {
		logger.E("Invalid data, abort push to queue!")
		return invar.ErrInvalidData
	}

	if len(maxlimits) > 0 && maxlimits[0] > 0 && t.queue.Len() > maxlimits[0] {
		logger.E("Task queue too heavy on oversize", maxlimits[0])
		return invar.ErrPoolFull
	}

	t.queue.Push(taskdata)
	t.asyncPostNext("Post")
	return nil
}

// Start runtime to post action
func (t *QTask) asyncPostNext(action string) {
	logger.D("Start runtime for [" + action + "] action")
	go func() { ttaskchan <- action }()
}

// Start task monitor to listen tasks pushed into queue, and execute it
func (t *QTask) startTaskMonitor(callback TaskCallback) {
	for {
		logger.I("Blocking for task require select...")

		select {
		case action := <-ttaskchan:
			logger.I("Received request from:", action)
			if callback == nil {
				logger.E("Nil task callback, abort request")
				break
			}

			// check current if executing status
			if t.executing {
				logger.W("Bussying now, try the next time...")
				break
			}

			// flag on executing and popup the topmost task to execte
			t.executing = true
			taskdata, err := t.queue.Pop()
			if err != nil {
				t.executing = false
				logger.I("Executed all tasks")
				break
			}

			if err := callback(taskdata); err != nil {
				logger.E("Execute task callback err:", err)
				if t.interrupt {
					logger.I("Interrupted tasks when case error")
					t.executing = false
					break
				}
			}
			if t.interval > 0 {
				logger.I("Waiting to next task after:", t.interval)
				time.Sleep(t.interval)
			}
			t.executing = false
		}
	}
}
