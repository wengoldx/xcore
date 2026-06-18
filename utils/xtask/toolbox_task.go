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

package xtask

import (
	"github.com/astaxie/beego/toolbox"
	"github.com/wengoldx/xcore/logger"
)

/* ------------------------------------------------------------------- */
/* For Toolbox Time Task                                               */
/* ------------------------------------------------------------------- */

// Task datas for multipe generate.
//
// The spec six columns mean ( * * * * * *)：
//
//	| second | minture | hour | day  | month | week                 |
//	| 0-59   | 0-59    | 1-23 | 1-31 | 1-12  | 0-6 (0 means Sunday) |
//
// Cron signals：
//
//	`*`  : any time
//	`,`  : separate signal
//	`-`  : duration
//	`/n` : do as n times of time duration
//
// Cron spec samples:
//
// ---
//
//	0/30 * * * * *                  every 30s
//	0 43 21 * * *                   21:43
//	0 15 05 * * *                   05:15
//	0 0 17 * * *                    17:00
//	0 0 17 * * 1                    17:00 in every Monday
//	0 0,10 17 * * 0,2,3             17:00 and 17:10 in every Sunday, Tuesday and Wednesday
//	0 0-10 17 1 * *                 17:00 to 17:10 in 1 min duration each time on the first day of month
//	0 0 0 1,15 * 1                  0:00 on the 1st day and 15th day of month
//	0 42 4 1 * *                    4:42 on the 1st day of month
//	0 0 21 * * 1-6                  21:00 from Monday to Saturday
//	0 0,10,20,30,40,50 * * * *      every 10 min duration
//	0 */10 * * * *                  every 10 min duration
//	0 * 1 * * *                     1:00 to 1:59 in 1 min duration each time
//	0 0 1 * * *                     1:00
//	0 0 */1 * * *                   0 min of hour in 1 hour duration
//	0 0 * * * *                     0 min of hour in 1 hour duration
//	0 2 8-20/3 * * *                8:02, 11:02, 14:02, 17:02, 20:02
//	0 30 5 1,15 * *                 5:30 on the 1st day and 15th day of month
//
// ---
type WTask struct {
	Name    string           // monitor task name
	Func    toolbox.TaskFunc // monitor task execute function
	Spec    string           // monitor task interval
	ForProd bool             // indicate the task only for prod mode, default no limit
}

// Add a single monitor task to list
func AddTask(tname, spec string, f toolbox.TaskFunc) {
	monitor := toolbox.NewTask(tname, spec, f)
	monitor.ErrLimit = 0 // not interupt when case error.
	toolbox.AddTask(tname, monitor)
}

// Create tasks and start them as monitors.
func StartTasks(monitors []*WTask) {
	for _, m := range monitors {
		if m.ForProd && !logger.IsRunProd() {
			logger.W("Filter out task:", m.Name, "on dev mode")
			continue
		}
		AddTask(m.Name, m.Spec, m.Func)
	}
	toolbox.StartTask()
}

// Return task if exist, or nil when unexist
func GetTask(tname string) *toolbox.Task {
	if tasker, ok := toolbox.AdminTaskList[tname]; ok {
		return tasker.(*toolbox.Task)
	}
	return nil
}
