// Copyright (c) 2018-Now Dunyu All Rights Reserved.
//
// Author      : https://www.quantkernel.com
// Email       : ping.yang@wengold.net
//
// Prismy.No | Date       | Modified by. | Description
// -------------------------------------------------------------------
// 00001       2025/05/06   youhei         New version
// -------------------------------------------------------------------

package cmdx

import (
	"encoding/base64"
	"os"
	"os/exec"

	"github.com/wengoldx/xcore/logger"
)

/* ------------------------------------------------------------------- */
/* Command Utils                                                       */
/* ------------------------------------------------------------------- */

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
func Shutdown() bool {
	cmd := exec.Command("shutdown", "-h", "now")
	if err := cmd.Run(); err != nil {
		logger.E("Shutdown, err:", err.Error())
		return false
	}
	return true
}

/* ------------------------------------------------------------------- */
/* Encode & Decode Script Utils                                        */
/* ------------------------------------------------------------------- */

// Encode script file content to base64 formated string,
// the encoded string easy define as a const string in project code,
// to hide and recreate script file whenever on need.
//
// # USAGE:
//
//	script_file := "./sample_script.sh"
//	enstr := cmde.EncodeScript(script_file)
//	fmt.Println(enstr) // print encoded base64 string.
//
//	// define const script var what you want.
//	const SAMPLE_SCRIPT = "the-print-encoded-base64-string"
//
//	See cmde.DecodeScript() to decode and save to script file.
func EncodeScript(script string) (string, error) {
	buf, err := os.ReadFile(script)
	if err != nil {
		return "", err
	}
	en := base64.StdEncoding
	return en.EncodeToString(buf), nil
}

// Decode script content and save to target executeable scrupt file.
//
// # USAGE:
//
//	// SAMPLE_SCRIPT : encoded script file content,
//	// script_file   : output script shell file.
//	cmde.DecodeScript(SAMPLE_SCRIPT, script_file)
//
//	See cmde.EncodeScript() to encode script content.
func DecodeScript(src, script string) error {
	en := base64.StdEncoding
	buf, err := en.DecodeString(src)
	if err != nil {
		return err
	}

	return os.WriteFile(script, buf, 0755)
}
