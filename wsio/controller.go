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

	sio "github.com/googollee/go-socket.io"
	"github.com/wengoldx/xcore/logger"
)

// Auth client outset, it will disconnect when return no-nil error
//	@param token client login jwt-token contain uuid or optional data in claims key string
//	@return - string client uuid
//			- any client optional data parsed from token
//			- error Exception message
type AuthHandler func(token string) (string, string, error)

// Client connected callback, it will disconnect when return no-nil error
//	@param uuid client unique id
//	@param option client login optional data, maybe nil
//	@return - error Exception message
type ConnectHandler func(uuid, option string) error

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

// Socket signaling adaptor to register events
type SignalingAdaptor interface {

	// Retruen socket signaling events
	Signalings() []string

	// Dispath socket signaling callback by event
	Dispatch(evt string) SignalingEvent
}

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
	logger.E("[SIO] Response err >>", msg)
	return string(resp)
}

// Set handler to execute clients authenticate, connect and disconnect.
func SetHandlers(auth AuthHandler, conn ConnectHandler, disc DisconnectHandler) {
	wsc.authHandler, wsc.connHandler, wsc.discHandler = auth, conn, disc
	logger.I("[SIO] Set wsio handlers...")
}

// Set adapter to register socket signaling events.
func SetAdapter(adaptor SignalingAdaptor) error {
	if adaptor == nil {
		logger.W("[SIO] Invalid socket event adaptor!")
		return nil
	}

	evts := adaptor.Signalings()
	if len(evts) == 0 {
		logger.W("[SIO] No signaling event keys!")
		return nil
	}

	// register socket signaling events
	for _, evt := range evts {
		if evt != "" {
			callback := adaptor.Dispatch(evt)
			if callback != nil {
				if err := wsc.server.On(evt, callback); err != nil {
					return err
				}
				logger.I("[SIO] Bind signaling event:", evt)
			}
		}
	}
	return nil
}
