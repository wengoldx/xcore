// Copyright (c) 2018-Now Dunyu All Rights Reserved.
//
// Author      : https://www.wengold.net
// Email       : support@wengold.net
//
// Prismy.No | Date       | Modified by. | Description
// -------------------------------------------------------------------
// 00001       2022/11/25   jidi           New version
// -------------------------------------------------------------------

package utils

import (
	"sync"
	"time"

	"github.com/wengoldx/xcore/invar"
	"github.com/wengoldx/xcore/logger"
)

type ExecFunc func(data any)

// delayTask delay task node
type task struct {
	cycle_num int64 // number of cycles
	task_info any   // order no or refund no
	execfunc  ExecFunc
}

// TimeWheel delay task sets
type TimeWheel struct {
	sync.Mutex
	current_index int64         // current location ID
	slots         [3600][]*task // one hour is a lap
	start_time    time.Time     // delay queue start time
	is_start      bool
	timer_stop    chan struct{}
}

// Create a new time wheel.
func NewTimeWheel() *TimeWheel {
	wheel := &TimeWheel{
		current_index: 0,
		start_time:    time.Now(),
		is_start:      false,
		timer_stop:    make(chan struct{}),
	}
	for i := 0; i < 3600; i++ {
		wheel.slots[i] = make([]*task, 0)
	}
	return wheel
}

// AddDelayTask insert one new delay task to delay queue
func (d *TimeWheel) AddDelayTask(exectime time.Duration, data any, f ExecFunc) error {
	d.Lock()
	defer d.Unlock()
	if exectime == 0 {
		logger.E("Invaild Exec Time")
		return invar.ErrInvaildExecTime
	}

	cycleNum := int64(exectime/(3600*time.Second)) - 1
	taskIndex := int64(d.current_index+int64(exectime/time.Second)) % 3600
	if taskIndex == d.current_index && cycleNum != 0 {
		cycleNum--
	}
	d.slots[taskIndex] = append(d.slots[taskIndex], &task{
		cycle_num: cycleNum,
		task_info: data,
		execfunc:  f,
	})
	return nil
}

// timeLoop move the scanner once per second
func (d *TimeWheel) timeLoop() {
	t := time.NewTicker(1 * time.Second)
	d.taskExec()
	for {
		select {
		case <-d.timer_stop: // stop time ticker
			t.Stop()
			return
		case <-t.C:
			if d.current_index == 3599 { // used a circular queue for runnable task
				d.current_index = 0
			} else {
				d.current_index++
			}
			go d.taskExec()
		}
	}
}

// taskExec execute the task with the current index and cycle number is 0
func (d *TimeWheel) taskExec() {
	d.Lock()
	tasks := d.slots[d.current_index]
	for i := len(tasks) - 1; i >= 0; i-- {
		v := tasks[i]
		if v.cycle_num == 0 {
			logger.W("exec current task ", d.current_index)
			go v.execfunc(v.task_info)
			d.slots[d.current_index] = append(d.slots[d.current_index][0:i], d.slots[d.current_index][i+1:]...)
		} else {
			v.cycle_num--
		}
	}
	d.Unlock()
}

// Start
func (d *TimeWheel) Start() {
	d.Lock()
	defer d.Unlock()
	if d.is_start {
		return
	}
	d.is_start = true
	go d.timeLoop()
}

// Stop
func (d *TimeWheel) Stop() {
	d.Lock()
	defer d.Unlock()
	var signal struct{}
	d.timer_stop <- signal
}
