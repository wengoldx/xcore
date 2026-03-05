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
	"errors"
	"strings"
)

// WingErr constom error with code, you can use it as standry error object.
//
//	var err error
//	err = invar.ErrNotFound
type WingErr struct{ error }

var (
	ErrNotFound           = WingErr{errors.New("notfound")}                                // Error: notfound.
	ErrNotChanged         = WingErr{errors.New("not changed")}                             // Error: not changed.
	ErrNotInserted        = WingErr{errors.New("not inserted")}                            // Error: not inserted.
	ErrInvalidNum         = WingErr{errors.New("invalid number")}                          // Error: invalid number.
	ErrInvalidAccount     = WingErr{errors.New("invalid account")}                         // Error: invalid account.
	ErrInvalidToken       = WingErr{errors.New("invalid token")}                           // Error: invalid token.
	ErrInvalidRole        = WingErr{errors.New("invalid role")}                            // Error: invalid role.
	ErrInvalidClient      = WingErr{errors.New("invalid client")}                          // Error: invalid client.
	ErrInvalidDevice      = WingErr{errors.New("invalid device")}                          // Error: invalid device.
	ErrInvalidParams      = WingErr{errors.New("invalid params")}                          // Error: invalid params.
	ErrInvalidData        = WingErr{errors.New("invalid data")}                            // Error: invalid data.
	ErrInvalidState       = WingErr{errors.New("invalid state")}                           // Error: invalid state.
	ErrInvalidPhone       = WingErr{errors.New("invalid phone")}                           // Error: invalid phone.
	ErrInvalidEmail       = WingErr{errors.New("invalid email")}                           // Error: invalid email.
	ErrInvalidOptions     = WingErr{errors.New("invalid options")}                         // Error: invalid options.
	ErrInvalidConfigs     = WingErr{errors.New("invalid configs")}                         // Error: invalid configs.
	ErrInvaildTime        = WingErr{errors.New("invaild time")}                            // Error: invaild time.
	ErrInvalidName        = WingErr{errors.New("invaild name")}                            // Error: invaild name.
	ErrInvalidFile        = WingErr{errors.New("invalid file")}                            // Error: invalid file.
	ErrBadPubKey          = WingErr{errors.New("invalid public key")}                      // Error: invalid public key.
	ErrBadPriKey          = WingErr{errors.New("invalid private key")}                     // Error: invalid private key.
	ErrDupRegister        = WingErr{errors.New("duplicated registration")}                 // Error: duplicated registration.
	ErrDupLogin           = WingErr{errors.New("duplicated login")}                        // Error: duplicated login.
	ErrDupData            = WingErr{errors.New("duplicated data")}                         // Error: duplicated data.
	ErrDupAccount         = WingErr{errors.New("duplicated account")}                      // Error: duplicated account.
	ErrDupName            = WingErr{errors.New("duplicate name")}                          // Error: duplicate name.
	ErrDupKey             = WingErr{errors.New("duplicate key")}                           // Error: duplicate key.
	ErrDupEntry           = WingErr{errors.New("duplicate entry")}                         // Error: duplicate entry.
	ErrBadDBConnect       = WingErr{errors.New("database not connnect")}                   // Error: database not connnect.
	ErrBadSQLBuilder      = WingErr{errors.New("bad sql string builder")}                  // Error: bad sql string builder.
	ErrBadModelCreator    = WingErr{errors.New("bad model creator")}                       // Error: bad model creator.
	ErrTagOffline         = WingErr{errors.New("target offline")}                          // Error: target offline.
	ErrClientOffline      = WingErr{errors.New("client offline")}                          // Error: client offline.
	ErrTokenExpired       = WingErr{errors.New("token expired")}                           // Error: token expired.
	ErrUnkownCharType     = WingErr{errors.New("unkown chars type")}                       // Error: unkown chars type.
	ErrUnperparedState    = WingErr{errors.New("unperpared state")}                        // Error: unperpared state.
	ErrOrmNotUsing        = WingErr{errors.New("orm not using")}                           // Error: orm not using.
	ErrNoneRowFound       = WingErr{errors.New("none row found")}                          // Error: none row found.
	ErrSendFailed         = WingErr{errors.New("failed to send")}                          // Error: failed to send.
	ErrAuthDenied         = WingErr{errors.New("permission denied")}                       // Error: permission denied.
	ErrKeyLenSixteen      = WingErr{errors.New("require sixteen-length secret key")}       // Error: require sixteen-length secret key.
	ErrOverTimes          = WingErr{errors.New("over retry times")}                        // Error: over retry times.
	ErrSetFrameNil        = WingErr{errors.New("failed clear frame meta")}                 // Error: failed clear frame meta.
	ErrNotSupport         = WingErr{errors.New("operation not support")}                   // Error: operation not support.
	ErrFailSendHead       = WingErr{errors.New("failed send head bytes")}                  // Error: failed send head bytes.
	ErrFailSendBody       = WingErr{errors.New("failed send body bytes")}                  // Error: failed send body bytes.
	ErrReadBytes          = WingErr{errors.New("error read bytes")}                        // Error: error read bytes.
	ErrInternalServer     = WingErr{errors.New("internal server error")}                   // Error: internal server error.
	ErrCreateByte         = WingErr{errors.New("failed create bytes: system protection")}  // Error: failed create bytes: system protection.
	ErrFileNotFound       = WingErr{errors.New("file not found")}                          // Error: file not found.
	ErrDownloadFile       = WingErr{errors.New("failed download file")}                    // Error: failed download file.
	ErrOpenSourceFile     = WingErr{errors.New("failed open source file")}                 // Error: failed open source file.
	ErrAlreadyConn        = WingErr{errors.New("already connected")}                       // Error: already connected.
	ErrEmptyReponse       = WingErr{errors.New("received empty response")}                 // Error: received empty response.
	ErrReadConf           = WingErr{errors.New("failed load config file")}                 // Error: failed load config file.
	ErrUnexpectedDir      = WingErr{errors.New("expect file path not directory")}          // Error: expect file path not directory.
	ErrWriteMD5           = WingErr{errors.New("failed write to md5")}                     // Error: failed write to md5.
	ErrWriteOut           = WingErr{errors.New("failed write out")}                        // Error: failed write out.
	ErrHandleDownload     = WingErr{errors.New("failed handle download file")}             // Error: failed handle download file.
	ErrFullConnPool       = WingErr{errors.New("connection pool is full")}                 // Error: connection pool is full.
	ErrPoolSize           = WingErr{errors.New("thread pool size value must be positive")} // Error: thread pool size value must be positive.
	ErrPoolFull           = WingErr{errors.New("pool is full, can not take any more")}     // Error: pool is full, can not take any more.
	ErrCheckDB            = WingErr{errors.New("check database: failed retry many times")} // Error: check database: failed retry many times.
	ErrFetchDB            = WingErr{errors.New("fetch database connection timeout")}       // Error: fetch database connection timeout.
	ErrReadFileBody       = WingErr{errors.New("failed read file content")}                // Error: failed read file content.
	ErrNilFrame           = WingErr{errors.New("frame is null")}                           // Error: frame is null.
	ErrNoStorage          = WingErr{errors.New("no storage server available")}             // Error: no storage server available.
	ErrUnmatchLen         = WingErr{errors.New("unmatch download file length")}            // Error: unmatch download file length.
	ErrCopyFile           = WingErr{errors.New("failed copy file")}                        // Error: failed copy file.
	ErrEmptyData          = WingErr{errors.New("empty data")}                              // Error: empty data.
	ErrImgOverSize        = WingErr{errors.New("image file size over")}                    // Error: image file size over.
	ErrAudioOverSize      = WingErr{errors.New("audio file size over")}                    // Error: audio file size over.
	ErrVideoOverSize      = WingErr{errors.New("video file size over")}                    // Error: video file size over.
	ErrNoAssociatedExpire = WingErr{errors.New("no associated expire")}                    // Error: no associated expire.
	ErrUnsupportFormat    = WingErr{errors.New("unsupported format data")}                 // Error: unsupported format data.
	ErrUnsupportedFile    = WingErr{errors.New("unsupported file format")}                 // Error: unsupported file format.
	ErrUnexistKey         = WingErr{errors.New("unexist key")}                             // Error: unexist key.
	ErrUnexistRedisKey    = WingErr{errors.New("unexist redis key")}                       // Error: unexist redis key.
	ErrUnexistLifecycle   = WingErr{errors.New("unexist lifecycle configs")}               // Error: unexist lifecycle configs.
	ErrSetLifecycleTag    = WingErr{errors.New("failed set file lifecycle tag")}           // Error: failed set file lifecycle tag.
	ErrInactiveAccount    = WingErr{errors.New("inactive status account")}                 // Error: inactive status account.
	ErrCaseException      = WingErr{errors.New("case exception")}                          // Error: case exception.
)

// Create a WingErr from given message.
func NewError(message string) *WingErr {
	return &WingErr{errors.New(message)}
}

// Return WingErr object replic with additions message.
//
//	err := invar.ErrNotFound.Replic("column xxx is missing")
//	// err message is: notfound - column xxx is missing
func (w *WingErr) Replic(additions ...string) *WingErr {
	if len(additions) > 0 {
		return NewError(w.Error() + " - " + strings.Join(additions, " "))
	}
	return NewError(w.Error())
}

// Return true if tow error messages matched.
func (w *WingErr) Equal(o error) bool {
	return w.Error() == o.Error()
}

// Return true if error message contain the WingErr string.
func (w *WingErr) SubOf(e error) bool {
	return IsError(e, w.Error())
}

// Check if error message contain given string
func IsError(e error, m string) bool {
	es, ms := strings.ToLower(e.Error()), strings.ToLower(m)
	return strings.Contains(es, ms)
}

// Check given error if contain 'duplicate' string for as deplicate error.
func IsDupError(e error) bool {
	return IsError(e, "duplicate")
}
