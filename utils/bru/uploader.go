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

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/wengoldx/xcore/invar"
	"github.com/wengoldx/xcore/logger"
	"github.com/wengoldx/xcore/secure"
	"github.com/wengoldx/xcore/utils"
)

// File slicing uploader.
//
// Slice the large file to the specified size chunks and upload them
// one by one, this uploader support breakpoint resume upload function.
//
// # File slice chunks upload proccess
//
//	Frontend Side   | Backend Server Side
//	----------------|----------------------------------------------
//	large.file      -> ~/cache/session-id     +-> ~/target/out.file
//	 |- [chunk 1  ] ->         |- chunk_1     |
//	 |- [chunk 2  ] ->         |- chunk_2     |
//	 |- [chunk ...] ->         |- chunk_...  /   [merge all chunks]
//
// # NOTICE:
//	- When any chunk fails to upload, it can be re-upload again before call Complete().
type BRUploader struct {
	cache    string              // Session chunks saved dir.
	sessions map[string]*Session // Async upload session map.
	mutex    sync.RWMutex        // Sync IO lock.
}

// Create a new BRUploader instance to upload file.
//
// # USAGE
//
// One large file upload as a session async job, and auto delete the
// session datas when file complete upload.
//
//	_bru := bru.NewUploader("./temp")
//	sid := _bru.NewSession(uid, file, outdir, hash, total, opts)
//
//	// the 'chunks' uploaded from frontent client.
//	for index := 0; index < len(chunks); index++ {
//		_bru.SaveChunk(sid, index, chunk, func(se *bru.Session) int {
//			// To check status or datas before save chunk to indexed file.
//			// ...
//			return invar.StatusOK
//		})
//	}
//
//	// merge all chunks to out single file what you want.
//	_bru.Complete(sid, func(se *bru.Session) int {
//		dstfile := filepath.Join(se.OutDir, se.File) // or other file path.
//		return _bru.MergeChunks(se, dstfile)
//	})
//
// # NOTICE:
//	- Try to make 'cache' dir when it empty and unexit.
//	- The session chunks saved current folder if 'cache' empty.
//	- Call Status() to get current uploading status for uncompelete session job.
//	- Call Cancel() to stop and clear target uploading session job.
func NewUploader(cache string) *BRUploader {
	cache = strings.TrimSpace(cache)
	if cache != "" && !utils.IsDir(cache) {
		if err := utils.MakeDirs(cache); err != nil {
			logger.E("Make cache dir:", cache, "err:", err)
		}
	}

	ss := make(map[string]*Session)
	return &BRUploader{cache: cache, sessions: ss}
}

// Create a new session to upload file with custom options.
func (t *BRUploader) NewSession(uid, file, outdir, hash string, total int64, opts ...any) (int, string) {
	sid := secure.NewLowNum() // breakpoint upload session id.
	sedir := filepath.Join(t.cache, sid)
	if err := utils.MakeDirs(sedir); err != nil {
		logger.E("Create session:", sid, "dir, err:", err)
		return invar.E404Exception, ""
	}
	counts := t.countChunks(total)

	t.mutex.Lock()
	t.sessions[sid] = &Session{
		SID: sid, File: file, OutDir: outdir, Hash: hash, Total: total,
		Chunks: make(map[int]utils.TNone), Counts: counts,
		UID: uid, Opts: utils.Variable(opts, nil),
	}
	t.mutex.Unlock()

	logger.I("Created session:", sid, "to upload:", file)
	return invar.StatusOK, sid
}

// Check target chunk uploaded status, and save the chunk datas to cache
// indexed file when it not upload or failed uploaded.
func (t *BRUploader) SaveChunk(sid string, index int, chunk []byte, check ...NextHandler) int {
	se := t.getSession(sid)

	// 1. check chunk data and upload status.
	chunksize := int64(len(chunk))
	if se == nil || chunksize <= 0 {
		logger.E("Unexist session:", sid, "or empty chunk!")
		return invar.E404Exception
	} else if _, ok := se.Chunks[index]; ok {
		logger.W("Chunk index:", index, "already uploaded!")
		return invar.StatusOK
	} else if chunksize != CHUNK_SIZE {
		logger.W("Unmatched chunk size:", chunksize)
	}

	// 2. check permission before save chunk to file.
	if len(check) > 0 && check[0] != nil {
		if state := check[0](se); state != invar.StatusOK {
			logger.E("Denied save session:", sid, "- chunk:", index)
			return state
		}
	}

	// 3. save chunk datas to indexed cache file.
	cf := filepath.Join(t.cache, sid, t.chunkFile(index))
	if err := os.WriteFile(cf, chunk, 0644); err != nil {
		logger.E("Save session:", sid, "chunk:", index, "err:", err)
		return invar.E404Exception
	}

	// 4. update session status.
	t.mutex.Lock()
	se.Chunks[index] = utils.NONE // mark chunk uploaded status.
	se.Upload += chunksize        // count uploaded total size.
	t.mutex.Unlock()
	return invar.StatusOK
}

// Check all chunks upload status and merge to target file.
//
// # WARNING:
//
// You must call MergeChunk() in 'merge' callback, to merge all
// uploaded chunks to out single file, it will salfly to delete
// the cache folder of current session after success merged.
func (t *BRUploader) Complete(sid string, merge NextHandler) int {
	// 1. check upload status.
	se := t.getSession(sid)
	if se == nil || merge == nil {
		logger.E("Unexist session:", sid, "or nil merge hander!")
		return invar.E404Exception
	}

	// 2. check if all chunks are uploaded.
	for i := 0; i < se.Counts; i++ {
		if _, ok := se.Chunks[i]; !ok {
			logger.E("Chunk index:", i, "missing...")
			return invar.E412InvalidState
		}
	}

	// 3. merge chunks to target file.
	if state := merge(se); state != invar.StatusOK {
		logger.E("Failed merge session:", sid, "file!")
		return state
	}

	// 4. update the session status to completed.
	t.removeSession(sid)
	return invar.StatusOK
}

// Merge all chunks to out single file and validate it hash code if
// matched, then delete the session cache folder when merge success.
func (t *BRUploader) MergeChunks(se *Session, outfile string) int {
	sid, hash, totals := se.SID, se.Hash, se.Counts

	// 1. try open or create target file.
	tf, err := os.Create(outfile)
	if err != nil {
		logger.E("Create file:", outfile, "err:", err)
		return invar.E404Exception
	}
	defer tf.Close()

	// 2. Merge all chunks to target file.
	sedir := filepath.Join(t.cache, sid)
	for i := 0; i < totals; i++ {
		chunkfile := filepath.Join(sedir, t.chunkFile(i))
		chunkdata, err := os.ReadFile(chunkfile)
		if err != nil {
			logger.E("Read chunk file, err:", err)
			return invar.E404Exception
		}

		if _, err := tf.Write(chunkdata); err != nil {
			logger.E("Write chunk to file, err:", err)
			return invar.E404Exception
		}
	}

	// 3. check merged file md5 hash code.
	if code, err := utils.FileAbstract(outfile); err != nil {
		logger.E("Validate file:", outfile, "err:", err)
		return invar.E404Exception
	} else if code != hash {
		logger.E("Not matched md5 hash -", outfile)
		return invar.E404Exception
	}

	// 4. clear session cache folder.
	go utils.DeleteFolder(sedir)
	return invar.StatusOK
}

// Return current uploaded bytes and all missing chunks index,
// so the frontend may upload them again.
func (t *BRUploader) Status(sid string) (int, *Status) {
	// 1. check upload status.
	se := t.getSession(sid)
	if se == nil {
		logger.E("Unexist session:", sid)
		return invar.E404Exception, nil
	}

	// 2. retrieve missing blocks.
	idxs := []int{}
	for i := 0; i < se.Counts; i++ {
		if _, ok := se.Chunks[i]; !ok {
			idxs = append(idxs, i)
		}
	}

	// 3. return uploading status.
	return invar.StatusOK, &Status{
		Upload: se.Upload, Total: se.Total, Missing: idxs,
	}
}

// Cancel target uploading session and remove cache folder.
//
// # NOTICE
//	- Set force as true to clear session datas without check user.
//	- By default, the Cancel() will check target request user.
func (t *BRUploader) Cancel(uid, sid string, force ...bool) int {
	// 1. check session status and validate user.
	if !utils.Variable(force, false) {
		se := t.getSession(sid)
		if se == nil || se.UID != uid {
			logger.E("Invalid session:", sid)
			return invar.E404Exception
		}
	}

	// 2. clear uploaded cache folder.
	sedir := filepath.Join(t.cache, sid)
	go utils.DeleteFolder(sedir)

	// 3. delete session.
	t.removeSession(sid)
	return invar.StatusOK
}

/* ------------------------------------------------------------------- */
/* Utils Methods                                                       */
/* ------------------------------------------------------------------- */

// Return session object, find by id on mutex locked.
func (t *BRUploader) getSession(sid string) *Session {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	return t.sessions[sid]
}

// Delete session under mutex locked state, find by id.
func (t *BRUploader) removeSession(sid string) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	delete(t.sessions, sid)
}

// Return total chunks count by file content total size.
func (t *BRUploader) countChunks(total int64) int {
	return int((total + CHUNK_SIZE - 1) / CHUNK_SIZE)
}

// Return chunk file name at target index.
func (t *BRUploader) chunkFile(index int) string {
	return fmt.Sprintf("chunk_%d", index)
}
