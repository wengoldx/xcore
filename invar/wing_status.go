// Copyright (c) 2018-Now Dunyu All Rights Reserved.
//
// Author      : https://www.wengold.net
// Email       : support@wengold.net
//
// Prismy.No | Date       | Modified by. | Description
// -------------------------------------------------------------------
// 00001       2019/05/22   yangping       New version
// -------------------------------------------------------------------
package invar

import (
	"net/http"
)

const (
	StatusOK             = http.StatusOK          // success status
	StatusExError        = http.StatusAccepted    // response extend error on 202 code
	StatusBadFile        = http.StatusNoContent   // response file unsaved on 204 code
	E400ParseParams      = http.StatusBadRequest
	E401Unauthorized     = http.StatusUnauthorized
	E403PermissionDenied = http.StatusForbidden
	E404Exception        = http.StatusNotFound
	E405FuncDisabled     = http.StatusMethodNotAllowed
	E406InputParams      = http.StatusNotAcceptable
	E408Timeout          = http.StatusRequestTimeout
	E409Duplicate        = http.StatusConflict
	E410Gone             = http.StatusGone
	E412InvalidState     = http.StatusPreconditionFailed
	E423Locked           = http.StatusLocked
	E426UpgradeRequired  = http.StatusUpgradeRequired
)

var statusText = map[int]string{
	StatusOK:             "OK",
	StatusExError:        "Response Extend Error",
	StatusBadFile:        "Response File Unsaved Error",
	E400ParseParams:      "Parse Input Params Error",
	E401Unauthorized:     "Unauthorized",
	E403PermissionDenied: "Permission Denied",
	E404Exception:        "Case Exception",
	E405FuncDisabled:     "Function Disabled",
	E406InputParams:      "Invalid Input Params",
	E408Timeout:          "Request Timeout",
	E409Duplicate:        "Duplicate Request",
	E410Gone:             "Gone",
	E412InvalidState:     "Invalid State",
	E423Locked:           "Resource Locked",
	E426UpgradeRequired:  "Upgrade Header Required",
}

// StatusText returns a text for the HTTP status code,
// It returns the empty string if the code is unknown.
func StatusText(code int) string {
	codetext := statusText[code]
	if codetext == "" {
		return http.StatusText(code)
	}
	return codetext
}
