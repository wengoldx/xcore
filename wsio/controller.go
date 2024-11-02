// Copyright (c) 2018-Now Dunyu All Rights Reserved.
//
// Author      : https://www.wengold.net
// Email       : support@wengold.net
//
// Prismy.No | Date       | Modified by. | Description
// -------------------------------------------------------------------
// 00001       2022/02/09   yangping       New version
// -------------------------------------------------------------------

package wsio

import (
	"encoding/json"
	"net/url"

	sio "github.com/googollee/go-socket.io"
)

// Auth client outset, it will disconnect when return no-nil error
//	@param forms form values parse from url
//	@param token client login jwt-token contain uuid or optional data in claims key string
//	@return - string client uuid
//			- any client optional data parsed from token
//			- error Exception message
type AuthHandler func(forms url.Values, token string) (string, string, error)

// Client connected callback, it will disconnect when return no-nil error
//	@param sc current socket client
//	@param uuid client unique id
//	@param option client login optional data, maybe nil
//	@return - error Exception message
type ConnectHandler func(sc sio.Socket, uuid, option string) error

// Client will disconnected handler function, it called before socket client disconnect.
//	@param sc current socket client
//	@param uuid client unique id
type WillDisconHandler func(sc sio.Socket, uuid string)

// Client disconnected handler function
//	@param uuid client unique id
//	@param option client login optional data, maybe nil
//
// `NOTICE` :
//
// The client of uuid already released when call this event function.
type DisconnectHandler func(uuid, option string)

// Socket signlaing event function
type SignalingEvent func(sc sio.Socket, uuid, params string) string

// Socket event ack
type EventAck struct {
	State   int    `json:"state"`
	Message string `json:"message"`
}

const (
	// StSuccess success status
	StSuccess = iota + 1

	// StError error status
	StError
)

// Response normal ack to socket client
func AckResp(msg string) string {
	resp, _ := json.Marshal(&EventAck{
		State: StSuccess, Message: msg,
	})
	return string(resp)
}

// Response error ack to socket client
func AckError(msg string) string {
	resp, _ := json.Marshal(&EventAck{
		State: StError, Message: msg,
	})
	siolog.E("Response err >>", msg)
	return string(resp)
}
