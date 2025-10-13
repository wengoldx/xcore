// Copyright (c) 2018-Now Dunyu All Rights Reserved.
//
// Author      : https://www.wengold.net
// Email       : support@wengold.net
//
// Prismy.No | Date       | Modified by. | Description
// -------------------------------------------------------------------
// 00001       2022/02/09   yangping       New version
// -------------------------------------------------------------------

package clients

import (
	sio "github.com/googollee/go-socket.io"
	"github.com/wengoldx/xcore/invar"
)

// socket connected client
type client struct {
	id     string     // Client id, maybe same as uuid
	option string     // Client optional data, can save as nickname, category, type, struct json string and so on
	socket sio.Socket // Client socket.io connection.
}

// Create a new client with unique id.
func newClient(cid string) *client {
	return &client{id: cid, socket: nil}
}

// Return client id that set when register on client connnecte.
func (c *client) UID() string {
	return c.id
}

// Return client extra optional data.
func (c *client) Option() string {
	return c.option
}

// Set client optional data.
func (c *client) SetOption(opt string) {
	c.option = opt
}

// Send signaling message to client.
func (c *client) Send(evt, msg string) error {
	if !c.registered() {
		return invar.ErrInvalidState
	}
	if err := c.socket.Emit(evt, msg); err != nil {
		siolog.E("Send [", evt, "], err:", err)
		return err
	}
	siolog.I("Send to", c.id, "[", evt, "] >>", msg)
	return nil
}

// Push client join to given room.
func (c *client) Join(room string) error {
	if !c.registered() {
		return invar.ErrInvalidState
	} else if room == "" {
		return invar.ErrInvalidParams
	}

	// check client if already joined given room
	rooms := c.socket.Rooms()
	for _, joinedroom := range rooms {
		if joinedroom == room {
			return nil
		}
	}

	siolog.I("Client:", c.id, "join room:", room)
	return c.socket.Join(room)
}

// Pull client leave from given room.
func (c *client) Leave(room string) error {
	if !c.registered() {
		return invar.ErrInvalidState
	} else if room == "" {
		return invar.ErrInvalidParams
	}

	siolog.I("Client:", c.id, "leave room:", room)
	return c.socket.Leave(room)
}

// Pull client leave all joined rooms.
func (c *client) LeaveRooms() error {
	if !c.registered() {
		return invar.ErrInvalidState
	}

	rooms := c.socket.Rooms()
	for _, room := range rooms {
		if err := c.socket.Leave(room); err != nil {
			return err
		}
	}

	siolog.I("Client:", c.id, "leave all rooms")
	return nil
}

// Return client joined rooms, it maybe nil when not joined.
func (c *client) Rooms() ([]string, error) {
	if !c.registered() {
		return nil, invar.ErrInvalidState
	}
	return c.socket.Rooms(), nil
}

// Broadcast signaling message to rooms that given by input param or all client joined.
//
// # NOTICE:
//
// The input param rooms should already joined by client when you want only
// broadcast event and message to indicate part of joined rooms.
//
//	rooms := client.Rooms()
//	client.Broadcast("evt-string", "message-content", rooms[0])
//	client.Broadcast("evt-string", "message-content", rooms[0], rooms[1])
//	client.Broadcast("evt-string", "message-content", rooms[:2]...)
//
// Or, not set input param rooms, it will broadcast event and message to all
// rooms that joined by client as:
//
//	client.Broadcast("evt-string", "message-content")
func (c *client) Broadcast(evt, msg string, rooms ...string) error {
	if !c.registered() {
		return invar.ErrInvalidState
	}

	// get target rooms from input params or joined
	var tagrooms []string
	if len(rooms) > 0 {
		tagrooms = rooms
	} else {
		tagrooms = c.socket.Rooms()
	}

	// execute broadcast to valid target rooms
	for _, room := range tagrooms {
		if err := c.socket.BroadcastTo(room, evt, msg); err != nil {
			siolog.E("Client", c.id, "broadcast [", evt, "] err:", err)
			return err
		}
		siolog.I("Broadcast to", c.id, "[", evt, "]", room, ">>", msg)
	}
	return nil
}

// --------

// Binds the socket with client.
func (c *client) register(sc sio.Socket, opt string) error {
	if c.registered() {
		cid, sid := c.id, sc.Id()
		siolog.E("Client", cid, "duplicate bind socket", sid)
		return invar.ErrDupRegister
	}
	c.socket, c.option = sc, opt
	return nil
}

// Unbind the socket with client and disconnect.
func (c *client) deregister() {
	if c.registered() {
		sid := c.socket.Id()
		siolog.I("Client", c.id, "unbind socket", sid)
		c.socket.Disconnect()
		c.socket = nil
	}
}

// Check client if bind with valid socket, true is bind.
func (c *client) registered() bool {
	return c.socket != nil
}
