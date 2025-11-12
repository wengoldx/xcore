// Copyright (c) 2018-Now Dunyu All Rights Reserved.
//
// Author      : https://www.quantkernel.com
// Email       : ping.yang@wengold.net
//
// Prismy.No | Date       | Modified by. | Description
// -------------------------------------------------------------------
// 00001       2025/05/06   youhei         New version
// -------------------------------------------------------------------

package cmde

import (
	"os/exec"

	"github.com/wengoldx/xcore/logger"
)

// Request reboot current device.
//
// # WARNING:
//	- This method only for linux.
func Reboot() bool {
	cmd := exec.Command("shutdown", "-r", "now")
	if err := cmd.Run(); err != nil {
		logger.E("Reboot, err:", err)
		return false
	}
	return true
}

// Request reboot current device.
//
// # WARNING:
//	- This method only for linux.
func Shutdown(err ConsoleHandler) bool {
	cmd := exec.Command("shutdown", "-h", "now")
	if err := cmd.Run(); err != nil {
		logger.E("Shutdown, err:", err.Error())
		return false
	}
	return true
}
