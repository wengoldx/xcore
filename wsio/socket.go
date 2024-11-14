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
	"net/http"
	"strings"
	"sync"
	"time"
	"unsafe"

	"github.com/astaxie/beego"
	sio "github.com/googollee/go-socket.io"
	"github.com/wengoldx/xcore/invar"
	"github.com/wengoldx/xcore/logger"
	"github.com/wengoldx/xcore/wsio/clients"
)

// Client datas for temp cache
type clientOpt struct {
	CID string // Client unique id (maybe same as account uuid)
	Opt string // Client optional data
}

// Socket.io connecte information, it will generate a socket server
// and register as http Handler to listen socket signalings.
//
// `NOTICE` :
//
// - The request client MUST set header have 'Author' : WENGOLD-V1.2 and
// no-empty auth Token as 'Token' : 'plaintext uuid or base64 formated string'.
//
// - DO NOT CHANGE THE go-socket.io MODULE VERSION, FIX IT IN 1.0.1 for
// all UE4/UE5, Nodejs, python3 clients. the go-socket.io version 1.6.2
// not matche or well support UE4/UE5 (No disconnected notify on server).
//
// see more [UE Socket.IO Plugin Usage](http://10.239.20.244:8090/pages/viewpage.action?pageId=7110992).
//
// ----
//
// `USAGE` :
//
//	// routers.go : register socket events
//	import "github.com/wengoldx/xcore/wsio"
//
//	init() {
//		handlers := Handlers{
//			AuthHandler: authHandlerFunc, ConnHandler: connHandlerFunc,
//			WillHandler: willHandlerFunc, DiscHandler: discHandlerFunc,
//		}
//		events := map[string]wsio.SignalingEvent{
//			"evt_msg"   : (&SocketController{Evt: "evt_msg"}).askMessage,
//			"evt_create": (&SocketController{Evt: "evt_create"}).createRoom,
//		}
//		wsio.SetupServer(handlers, events)
//	}
//
// You may config socket ping interval, timeout and using optinal data
// check in app.conf file as follow:
//
//	[wsio]
//	; Heartbeat ping interval, default 30 seconds
//	interval = 30
//
//	; Max heartbeat ping timeout, default 60 seconds
//	timeout = 60
//
//	; Using client optional data check, default false
//	optinal = false
//
//	; Using idles client function, default false
//	idels = false
type wingSIO struct {
	// mutex sync lock, protect client connecting
	lock sync.Mutex

	// socket server
	server *sio.Server

	// socket golbal handlers.
	handlers Handlers

	// socket events map, registry on server start.
	events map[string]SignalingEvent

	// http request pointer to client, cache datas temporary
	// only for client authenticate-connect process.
	options map[uintptr]*clientOpt

	// `WARNING` :
	//
	// the go-socket.io will call onConnect duplicate times when socket.io
	// client version not matched, so handle the valid first and abort the
	// invalid next time.
	//
	//	see more : go-socket.io@v1.0.1/parse.go   > Decode() > NextReader()
	//			 : go-socket.io@v1.0.1/socket.go  > loop() > for { onPacket }
	//			 : go-socket.io@v1.0.1/handler.go > onPacket()
	onceBunds map[uintptr]string // http request url to empty string (not used)
}

// Socket.IO handlers for auth, connect, disconnect.
type Handlers struct {
	AuthHandler AuthHandler       // socket callback to authenticate client.
	ConnHandler ConnectHandler    // socket callback to perpare after client connected.
	WillHandler WillDisconHandler // socket callback to do will before client disconnect.
	DiscHandler DisconnectHandler // socket callback to release after client disconnected.
}

// Socket connection server
var wsc *wingSIO

// Object logger with [SIO] mark for socket.io module
var siolog = logger.CatLogger("SIO")

var (
	serverPingInterval = 30 * time.Second
	serverPingTimeout  = 60 * time.Second
	maxConnectCount    = 200000

	// Check client option if empty when connnection is established,
	// if optinal data is empty the connect will not establish and disconnect.
	viaOption = false
)

// Setup socket.io server by manual called.
func SetupServer(cbs Handlers, evts map[string]SignalingEvent) {
	loadWsioConfigs()
	wsc = &wingSIO{
		options:   make(map[uintptr]*clientOpt),
		onceBunds: make(map[uintptr]string),
		handlers:  cbs, events: evts,
	}

	// set http handler for socke.io
	handler, err := wsc.createHandler()
	if err != nil {
		panic(err)
	}

	// set socket.io routers
	beego.Handler("/"+beego.BConfig.AppName+"/socket.io", handler)
	siolog.I("Initialized socket.io routers...")
}

// Read wsio configs from file
func loadWsioConfigs() {
	interval := beego.AppConfig.DefaultInt64("wsio::interval", 30)
	serverPingInterval = time.Duration(interval) * time.Second

	timeout := beego.AppConfig.DefaultInt64("wsio::timeout", 60)
	serverPingTimeout = time.Duration(timeout) * time.Second

	viaOption := beego.AppConfig.DefaultBool("wsio::optinal", false)
	siolog.I("Server configs, interval:", interval, "timeout:", timeout,
		"optional:", viaOption)
}

// Create http handler for socket.io
func (cc *wingSIO) createHandler() (http.Handler, error) {
	server, err := sio.NewServer(nil)
	if err != nil {
		return nil, err
	}
	cc.server = server

	// set socket.io ping interval and timeout
	siolog.I("Set ping-pong and timeout")
	server.SetPingInterval(serverPingInterval)
	server.SetPingTimeout(serverPingTimeout)

	// set max connection count
	server.SetMaxConnection(maxConnectCount)

	// set auth middleware for socket.io connection
	server.SetAllowRequest(func(req *http.Request) error {
		if err = cc.onAuthentication(req); err != nil {
			siolog.E("Authenticate err:", err)
			return err
		}
		return nil
	})

	// set connection event
	server.On("connection", func(sc sio.Socket) {
		cc.onConnect(sc)
	})

	// set disconnection event
	server.On("disconnection", func(sc sio.Socket) {
		cc.onDisconnected(sc)
	})

	cc.registryEvents()
	siolog.I("Created socket.io handler")
	return server, nil
}

// Set adapter to register socket signaling events.
func (cc *wingSIO) registryEvents() {
	if len(cc.events) == 0 {
		siolog.W("Not registry any socket events !!")
		return
	}

	// register socket signaling events
	for evt, callback := range cc.events {
		if evt != "" && callback != nil {
			controller := &WsioController{Evt: evt, hander: callback}
			if err := wsc.server.On(evt, controller.hander); err != nil {
				siolog.E("Bind socket event:", evt, "err:", err)
				continue
			}
		}
	}
	siolog.I("Bund socket events")
}

// Internal event of authentication, to get auth datas from request header
// and then call outside registered authentication handler to vertify token.
func (cc *wingSIO) onAuthentication(req *http.Request) error {
	author, token := req.Header.Get("Author"), ""
	if author == "" {
		author = req.Header.Get("Authoration")
	}

	/*
	 * `NOTICE` :
	 *
	 * (1). Try get Author from http header for Python3 and Unreal client;
	 *      but React frontend and Wechat App not support auth header as well,
	 *      so tail token string after request connect url for them.
	 *
	 * (2). Auth header 'WENGOLD' may upgraded to upper versions, so need just
	 *      check the perfix better than total match.
	 */
	if author != "" {
		// use auth header function for Python3 and Unreal client
		token = req.Header.Get("Token")
		if !strings.HasPrefix(author, "WENGOLD") || token == "" {
			siolog.E("Invalid authoration:", author, "token:", token)
			return invar.ErrAuthDenied
		}
	} else {
		// use URL + token string for React frontend and Wechat app client
		if err := req.ParseForm(); err != nil {
			siolog.E("Parse request form, err:", err)
			return err
		} else if token = req.Form.Get("token"); token == "" {
			siolog.E("Failed get token from request url!")
			return invar.ErrAuthDenied
		}
	}

	// auth client token by handler if set Authenticate function
	// handler, or just use token as uuid when not set.
	uuid, option := token, ""
	if cc.handlers.AuthHandler != nil {
		cid, opt, err := cc.handlers.AuthHandler(req.Form, token)
		if err != nil || cid == "" {
			siolog.E("Invalid cid:", cid, "or case err:", err)
			return invar.ErrAuthDenied
		} else if viaOption && opt == "" {
			siolog.E("Empty client", cid, "option data!")
			return invar.ErrAuthDenied
		}

		siolog.I("Decoded client token, cid:", cid, "opt:", opt)
		uuid, option = cid, opt
	}

	// bind http.Request -> cid
	h := uintptr(unsafe.Pointer(req))
	cc.bindHTTP2UUIDLocked(h, uuid, option)
	return nil
}

// onConnect event of connect
func (cc *wingSIO) onConnect(sc sio.Socket) {
	// found client id and unbind -> http.Request
	h := uintptr(unsafe.Pointer(sc.Request()))
	if _, ok := cc.onceBunds[h]; ok /* already bund */ {
		siolog.W("Duplicate onConnect, abort for", h)
		return
	}
	cc.onceBunds[h] = "" // cache the first time
	defer cc.clearBundCache(h)

	co := cc.unbindUUIDFromHTTPLocked(h)
	if co == nil || co.CID == "" {
		siolog.E("Invalid socket request bind!")
		sc.Disconnect()
		return
	}

	clientPool := clients.Clients()
	if err := clientPool.Register(sc, co.CID, co.Opt); err != nil {
		siolog.E("Failed register socket client:", co.CID)
		sc.Disconnect()
		return
	}

	// handle connect callback for socket with client id
	if cc.handlers.ConnHandler != nil {
		if err := cc.handlers.ConnHandler(sc, co.CID, co.Opt); err != nil {
			siolog.E("Client:", co.CID, "connect socket err:", err)
			sc.Disconnect()
		}
	}
	siolog.I("Connected socket client:", co.CID)
}

// onDisconnected event of disconnect
func (cc *wingSIO) onDisconnected(sc sio.Socket) {
	clientPool := clients.Clients()

	// hangle will callback before diconnect when
	// socket client object still valid.
	cid := clientPool.ClientID(sc.Id())
	if cc.handlers.WillHandler != nil && cid != "" {
		cc.handlers.WillHandler(sc, cid)
	}

	// deregister socket and disconnect it, then
	// call disconnected handler to release sources.
	opt := clientPool.Deregister(sc)
	if cc.handlers.DiscHandler != nil {
		cc.handlers.DiscHandler(cid, opt)
	}
}

// Bind http request pointer -> cid on locked status
func (cc *wingSIO) bindHTTP2UUIDLocked(h uintptr, cid, opt string) {
	cc.lock.Lock()
	defer cc.lock.Unlock()

	cc.options[h] = &clientOpt{CID: cid, Opt: opt}
}

// Unbind cid -> http request pointer on locked status
func (cc *wingSIO) unbindUUIDFromHTTPLocked(h uintptr) *clientOpt {
	cc.lock.Lock()
	defer cc.lock.Unlock()

	if data, ok := cc.options[h]; ok {
		co := &clientOpt{CID: data.CID, Opt: data.Opt}
		delete(cc.options, h)
		return co
	}
	return nil
}

// Clear the bind cache after 10ms
func (cc *wingSIO) clearBundCache(h uintptr) {
	go func(c *wingSIO, h uintptr) {
		time.Sleep(50 * time.Millisecond)
		delete(c.onceBunds, h)
	}(cc, h)
}
