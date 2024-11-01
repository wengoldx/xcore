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
	"sort"
	"sync"
	"time"

	sio "github.com/googollee/go-socket.io"
	"github.com/wengoldx/xcore/invar"
	"github.com/wengoldx/xcore/logger"
)

// ClientPool client pool
type ClientPool struct {
	lock    sync.Mutex         // Mutex sync lock
	clients map[string]*client // Client map as { client-id : client }
	s2c     map[string]string  // Socket id to client id as { socket-id : client-id }
	idles   map[string]int64   // Idle client weights as { client-id : idle-start-nanosecond }
}

// clientPool singleton instance
var clientPool *ClientPool

// idleWeight idle client weight
type idleWeight struct {
	uuid   string // Client unique id
	weight int64  // Client weight as unix nanosecond start idle
}

func init() {
	clientPool = &ClientPool{
		clients: make(map[string]*client),
		s2c:     make(map[string]string),
		idles:   make(map[string]int64),
	}
}

// Return ClientPool singleton
func Clients() *ClientPool {
	return clientPool
}

// Return client, it maybe nil if unexist.
func (cp *ClientPool) Client(cid string) *client {
	return cp.clients[cid]
}

// Register client and bind socket.
func (cp *ClientPool) Register(cid string, sc sio.Socket, opt string) error {
	cp.lock.Lock()
	defer cp.lock.Unlock()

	if err := cp.registerLocked(cid, sc, opt); err != nil {
		logger.E("[SIO] Regisger client, err:", err.Error())
		return err
	}

	cp.idleLocked(cid)
	return nil
}

// Deregister client and unbind socket.
func (cp *ClientPool) Deregister(sc sio.Socket) (string, string) {
	cp.lock.Lock()
	defer cp.lock.Unlock()

	return cp.deregisterLocked(sc)
}

// Check the client if exist.
func (cp *ClientPool) IsExist(cid string) bool {
	_, ok := cp.clients[cid]
	return ok
}

// Return client id from socket id.
func (cp *ClientPool) ClientID(sid string) string {
	cp.lock.Lock()
	defer cp.lock.Unlock()

	return cp.s2c[sid]
}

// Cache or refresh unix nanosecond time as weight.
func (cp *ClientPool) Idle(cid string) {
	cp.lock.Lock()
	defer cp.lock.Unlock()

	cp.idleLocked(cid)
}

// Remove client out of idles map whatever weight value over zero or not.
func (cp *ClientPool) LeaveIdle(cid string) {
	cp.lock.Lock()
	defer cp.lock.Unlock()

	cp.leaveIdleLocked(cid)
}

// Return idle clients ids
func (cp *ClientPool) Idles() []string {
	var idles []string
	for k := range cp.idles {
		idles = append(idles, k)
	}
	return idles
}

// Move client out of idle state without acquiring the lock.
func (cp *ClientPool) SortIdels() []string {
	if len(cp.idles) == 0 {
		return nil
	}

	var weights []idleWeight
	for k, v := range cp.idles {
		weights = append(weights, idleWeight{uuid: k, weight: v})
	}

	sort.Slice(weights, func(i, j int) bool {
		return weights[i].weight < weights[j].weight
	})

	// add each client's uuid to string array
	uuids := []string{}
	for _, v := range weights {
		uuids = append(uuids, v.uuid)
	}
	return uuids
}

// Count clients by option key, it not count the empty option clients.
func (cp *ClientPool) Counts() map[string]int {
	cnts := make(map[string]int)
	for _, client := range cp.clients {
		if client == nil || client.option == "" {
			continue
		}

		if cnt, ok := cnts[client.option]; ok {
			cnts[client.option] = cnt + 1
		} else {
			cnts[client.option] = 1
		}
	}
	return cnts
}

// -------- quick handle functions for indicate client

// Return client optinal data, it maybe nil.
func (cp *ClientPool) Option(cid string) string {
	if c, ok := cp.clients[cid]; ok {
		return c.option
	}
	return ""
}

// Set the client optinal data, maybe return error if not exist client.
func (cp *ClientPool) SetOption(cid, opt string) error {
	if c, ok := cp.clients[cid]; ok {
		c.option = opt
		return nil
	}
	return invar.ErrNotFound
}

// Send signaling with message to indicate client.
func (cp *ClientPool) Signaling(cid, evt, data string) error {
	if c, ok := cp.clients[cid]; ok {
		return c.Send(evt, data)
	}
	return invar.ErrTagOffline
}

// --------

// Register the client without acquiring the lock.
func (cp *ClientPool) registerLocked(cid string, sc sio.Socket, opt string) error {
	var newOne *client
	sid := sc.Id()

	if oldOne, ok := cp.clients[cid]; ok {
		oldOneID := oldOne.socket.Id()
		if oldOneID == sid {
			logger.W("[SIO] Client", cid, "already bind socket", sid)
			return nil
		}

		logger.W("[SIO] Drop bund socket", oldOneID)
		delete(cp.s2c, oldOneID)
		oldOne.deregister() // reset and disconnet the old socket
		newOne = oldOne
	} else {
		newOne = newClient(cid)
	}

	// bind client with socket
	if err := newOne.register(sc, opt); err != nil {
		return err
	}

	logger.I("[SIO] Client", cid, "bind socket", sid)
	cp.clients[cid] = newOne
	cp.s2c[sid] = cid // same as uuid
	return nil
}

// Deregister the client without acquiring the lock.
func (cp *ClientPool) deregisterLocked(sc sio.Socket) (string, string) {
	sid := sc.Id()
	if cid := cp.s2c[sid]; cid != "" {
		delete(cp.s2c, sid)

		cp.leaveIdleLocked(cid)
		if c := cp.clients[cid]; c != nil {
			delete(cp.clients, cid)
			c.deregister()
			return cid, c.option
		}
	}

	logger.I("[SIO] Disconnect socket", sid)
	sc.Disconnect()
	return "", ""
}

// Increate idle weight for client without acquiring the lock.
func (cp *ClientPool) idleLocked(cid string) {
	cp.idles[cid] = time.Now().UnixNano()
	logger.I("[SIO] Client", cid, "idle...")
}

// Move client out of idle state without acquiring the lock.
func (cp *ClientPool) leaveIdleLocked(cid string) {
	if _, ok := cp.idles[cid]; ok {
		logger.I("[SIO] Client", cid, "leave idle")
		delete(cp.idles, cid)
	}
}
