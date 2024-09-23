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

// WingErr constom error with code
type WingErr struct {
	error     // Simple use WingErr value as error type
	Code  int // Extend error code, start 0x1000
}

// WExErr extend error with code and error message
type WExErr struct {
	Code    int    `json:"code"    description:"Extend error code"`
	Message string `json:"message" description:"Extend error message"`
}

var (
	ErrNotFound            = WingErr{errors.New("not fount") /*                                    */, 0x1000}
	ErrInvalidNum          = WingErr{errors.New("invalid number") /*                               */, 0x1001}
	ErrInvalidAccount      = WingErr{errors.New("invalid account") /*                              */, 0x1002}
	ErrInvalidToken        = WingErr{errors.New("invalid token") /*                                */, 0x1003}
	ErrInvalidRole         = WingErr{errors.New("invalid role") /*                                 */, 0x1004}
	ErrInvalidClient       = WingErr{errors.New("invalid client") /*                               */, 0x1005}
	ErrInvalidDevice       = WingErr{errors.New("invalid device") /*                               */, 0x1006}
	ErrInvalidParams       = WingErr{errors.New("invalid params") /*                               */, 0x1007}
	ErrInvalidData         = WingErr{errors.New("invalid data") /*                                 */, 0x1008}
	ErrInvalidState        = WingErr{errors.New("invalid state") /*                                */, 0x1009}
	ErrInvalidPhone        = WingErr{errors.New("invalid phone") /*                                */, 0x100A}
	ErrInvalidEmail        = WingErr{errors.New("invalid email") /*                                */, 0x100B}
	ErrInvalidOptions      = WingErr{errors.New("invalid options") /*                              */, 0x100C}
	ErrInvalidRedisOptions = WingErr{errors.New("invalid redis options") /*                        */, 0x100D}
	ErrInvalidConfigs      = WingErr{errors.New("invalid config datas") /*                         */, 0x100E}
	ErrInvaildExecTime     = WingErr{errors.New("invaild execute time") /*                         */, 0x100F}
	ErrInvalidRealname     = WingErr{errors.New("invaild realname") /*                             */, 0x1010}
	ErrTagOffline          = WingErr{errors.New("target offline") /*                               */, 0x1011}
	ErrClientOffline       = WingErr{errors.New("client offline") /*                               */, 0x1012}
	ErrDupRegister         = WingErr{errors.New("duplicated registration") /*                      */, 0x1013}
	ErrDupLogin            = WingErr{errors.New("duplicated admin login") /*                       */, 0x1014}
	ErrDupData             = WingErr{errors.New("duplicated data") /*                              */, 0x1015}
	ErrDupAccount          = WingErr{errors.New("duplicated account") /*                           */, 0x1016}
	ErrDupName             = WingErr{errors.New("duplicate name") /*                               */, 0x1017}
	ErrDupKey              = WingErr{errors.New("duplicate key") /*                                */, 0x1018}
	ErrTokenExpired        = WingErr{errors.New("token expired") /*                                */, 0x1019}
	ErrBadPublicKey        = WingErr{errors.New("invalid public key") /*                           */, 0x101A}
	ErrBadPrivateKey       = WingErr{errors.New("invalid private key") /*                          */, 0x101B}
	ErrUnkownCharType      = WingErr{errors.New("unkown chars type") /*                            */, 0x101C}
	ErrUnperparedState     = WingErr{errors.New("unperpared state") /*                             */, 0x101D}
	ErrOrmNotUsing         = WingErr{errors.New("orm not using") /*                                */, 0x101E}
	ErrNoneRowFound        = WingErr{errors.New("none row found") /*                               */, 0x101F}
	ErrNotChanged          = WingErr{errors.New("not changed") /*                                  */, 0x1020}
	ErrNotInserted         = WingErr{errors.New("not inserted") /*                                 */, 0x1021}
	ErrSendFailed          = WingErr{errors.New("failed to send") /*                               */, 0x1022}
	ErrAuthDenied          = WingErr{errors.New("permission denied") /*                            */, 0x1023}
	ErrKeyLenSixteen       = WingErr{errors.New("require sixteen-length secret key") /*            */, 0x1024}
	ErrOverTimes           = WingErr{errors.New("over retry times") /*                             */, 0x1025}
	ErrSetFrameNil         = WingErr{errors.New("failed clear frame meta") /*                      */, 0x1026}
	ErrOperationNotSupport = WingErr{errors.New("operation not support") /*                        */, 0x1027}
	ErrSendHeadBytes       = WingErr{errors.New("failed send head bytes") /*                       */, 0x1028}
	ErrSendBodyBytes       = WingErr{errors.New("failed send body bytes") /*                       */, 0x1029}
	ErrReadBytes           = WingErr{errors.New("error read bytes") /*                             */, 0x102A}
	ErrInternalServer      = WingErr{errors.New("internal server error") /*                        */, 0x102B}
	ErrCreateByte          = WingErr{errors.New("failed create bytes: system protection") /*       */, 0x102C}
	ErrFileNotFound        = WingErr{errors.New("file not found") /*                               */, 0x102D}
	ErrDownloadFile        = WingErr{errors.New("failed download file") /*                         */, 0x102E}
	ErrOpenSourceFile      = WingErr{errors.New("failed open source file") /*                      */, 0x102F}
	ErrAlreadyConn         = WingErr{errors.New("already connected") /*                            */, 0x1030}
	ErrEmptyReponse        = WingErr{errors.New("received empty response") /*                      */, 0x1031}
	ErrReadConf            = WingErr{errors.New("failed load config file") /*                      */, 0x1032}
	ErrUnexpectedDir       = WingErr{errors.New("expect file path not directory") /*               */, 0x1033}
	ErrWriteMD5            = WingErr{errors.New("failed write to md5") /*                          */, 0x1034}
	ErrWriteOut            = WingErr{errors.New("failed write out") /*                             */, 0x1035}
	ErrHandleDownload      = WingErr{errors.New("failed handle download file") /*                  */, 0x1036}
	ErrFullConnPool        = WingErr{errors.New("connection pool is full") /*                      */, 0x1037}
	ErrPoolSize            = WingErr{errors.New("thread pool size value must be positive") /*      */, 0x1038}
	ErrPoolFull            = WingErr{errors.New("pool is full, can not take any more") /*          */, 0x1039}
	ErrCheckDB             = WingErr{errors.New("check database: failed retry many times") /*      */, 0x103A}
	ErrFetchDB             = WingErr{errors.New("fetch database connection timeout") /*            */, 0x103B}
	ErrReadFileBody        = WingErr{errors.New("failed read file content") /*                     */, 0x103C}
	ErrNilFrame            = WingErr{errors.New("frame is null") /*                                */, 0x103D}
	ErrNoStorage           = WingErr{errors.New("no storage server available") /*                  */, 0x103E}
	ErrUnmatchLen          = WingErr{errors.New("unmatch download file length") /*                 */, 0x103F}
	ErrCopyFile            = WingErr{errors.New("failed copy file") /*                             */, 0x1040}
	ErrEmptyData           = WingErr{errors.New("empty data") /*                                   */, 0x1041}
	ErrImgOverSize         = WingErr{errors.New("image file size over") /*                         */, 0x1042}
	ErrAudioOverSize       = WingErr{errors.New("audio file size over") /*                         */, 0x1043}
	ErrVideoOverSize       = WingErr{errors.New("video file size over") /*                         */, 0x1044}
	ErrNoAssociatedExpire  = WingErr{errors.New("no associated expire") /*                         */, 0x1045}
	ErrUnsupportFormat     = WingErr{errors.New("unsupported format data") /*                      */, 0x1046}
	ErrUnsupportedFile     = WingErr{errors.New("unsupported file format") /*                      */, 0x1047}
	ErrUnexistKey          = WingErr{errors.New("unexist key") /*                                  */, 0x1048}
	ErrUnexistRedisKey     = WingErr{errors.New("unexist redis key") /*                            */, 0x1049}
	ErrUnexistLifecycle    = WingErr{errors.New("unexist lifecycle configs") /*                    */, 0x104A}
	ErrSetLifecycleTag     = WingErr{errors.New("failed set file lifecycle tag") /*                */, 0x104B}
	ErrInactiveAccount     = WingErr{errors.New("inactive status account") /*                      */, 0x104C}
	ErrCaseException       = WingErr{errors.New("case exception") /*                               */, 0x104D}
)

// Create a WExErr from given code and message
func ExErr(code int, message string) *WExErr {
	return &WExErr{Code: code, Message: message}
}

// Create a WExErr from given code and message
func StErr(code int, message string) (int, *WExErr) {
	return StatusExError, &WExErr{Code: code, Message: message}
}

// Create a WingErr from given code and message
func WErr(code int, message string) *WingErr {
	return &WingErr{errors.New(message), code}
}

// Return standry error from WingErr object.
//
//	// Simple to get standry error object
//	var err error
//	err = invar.ErrNotFoune.Err()
func (w *WingErr) Err() error {
	return w
}

// Return WExErr extend error from WingErr object.
//
//	// Simple to get WExErr extend error
//	var exerr *invar.WExErr
//	exerr = invar.ErrNotFoune.ExErr()
//
//	// Directly using WingErr as error value
//	var err error
//	err := invar.ErrNotFoune
func (w *WingErr) ExErr() *WExErr {
	return &WExErr{Code: w.Code, Message: w.Error()}
}

// Return HTTP response code and WExErr extend error.
//
//	// Using for Restful API to response custom status and message.
//	http_resp_code, err := invar.ErrNotFoune.StErr()
func (w *WingErr) StErr() (int, *WExErr) {
	return StatusExError, &WExErr{Code: w.Code, Message: w.Error()}
}

/////////////////////////////////////

// Equal tow error if message same on char case
func EqualError(a, b error) bool {
	return a.Error() == b.Error()
}

// Equal tow error if message same ignoral char case
func EqualErrorFold(a, b error) bool {
	return strings.EqualFold(a.Error(), b.Error())
}

// Check if error message contain given error string
func ErrorContain(s, sub error) bool {
	return strings.Contains(s.Error(), sub.Error())
}

// Check if error message start given perfix
func ErrorStart(s, sub error) bool {
	return strings.HasPrefix(s.Error(), sub.Error())
}

// Check if error message start given perfix
func ErrorEnd(s, sub error) bool {
	return strings.HasSuffix(s.Error(), sub.Error())
}

// Check if error message contain given string
func IsError(e error, s string) bool {
	esu, su := strings.ToLower(e.Error()), strings.ToLower(s)
	return strings.Contains(esu, su)
}

// Check given error if duplicated errors
func IsDupError(e error) bool {
	return ErrorContain(e, ErrDupRegister) || IsError(e, "Duplicate entry") ||
		ErrorContain(e, ErrDupData) || ErrorContain(e, ErrDupAccount) ||
		ErrorContain(e, ErrDupName) || ErrorContain(e, ErrDupKey) ||
		ErrorContain(e, ErrDupLogin)
}
