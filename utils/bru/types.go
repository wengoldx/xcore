// Copyright (c) 2018-Now Dunyu All Rights Reserved.
//
// Author      : https://www.quantkernel.com
// Email       : ping.yang@wengold.net
//
// Prismy.No | Date       | Modified by. | Description
// -------------------------------------------------------------------
// 00001       2025/11/02   youhei         New version
// -------------------------------------------------------------------

package bru

import "github.com/wengoldx/xcore/utils"

// The handler to execute next action.
//
// # WARNING:
//	- Called for save chunk to indexed file by SaveChunk() method.
//	- Called for merge chunks to out single file by Complete() method.
type NextHandler func(se *Session) int

// The fixed chunk size to slice larget file, set as 5MB.
const CHUNK_SIZE = 5 * 1024 * 1024

// A session for one job to upload a single file.
type Session struct {
	/* Init params when session create. -------------------------- */

	SID    string // Session unique id.
	File   string // File name with suffix, without any paths.
	OutDir string // Relative dir to save out file, without file name.
	Hash   string // MD5 hash code of the out file to validate integrity.
	Total  int64  // File total size in bytes.

	/* Extend params for next job. ------------------------------- */

	UID  string // Who create this upload session.
	Opts any    // Optional payload, for next job after uploaded completed.

	/* Runtime status and values. -------------------------------- */

	Upload int64               // current already uploaded size in bytes.
	Counts int                 // total slice chunks counts of the session.
	Chunks map[int]utils.TNone // chunks upload flags, [index : _holder_], true for uploaded.
}

// Session uploading status.
type Status struct {
	SID     string // session uniqu id.
	Upload  int64  // current already uploaded bytes.
	Total   int64  // file content total size in bytes.
	Missing []int  // chunks index list which not uploaded.
}
