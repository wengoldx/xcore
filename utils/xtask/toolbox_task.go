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
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/toolbox"
	"github.com/wengoldx/xcore/logger"
)

/* ------------------------------------------------------------------- */
/* For Toolbox Time Task                                               */
/* ------------------------------------------------------------------- */

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
	monitor.ErrLimit = 0 // not interupt when case error.
	toolbox.AddTask(tname, monitor)
}

// Create tasks and start them as monitors.
func StartTasks(monitors []*WTask) {
	for _, m := range monitors {
		if m.ForProd && beego.BConfig.RunMode != "prod" {
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
