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

// client datas for temp cache
type clientOpt struct {
	UID string // Client unique id
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
// see more [UE Socket.IO Pluge Usage](http://10.239.20.244:8090/pages/viewpage.action?pageId=7110992).
//
// ----
//
// `USAGE` :
//
//	// routers.go : register socket events
//	import "github.com/wengoldx/xcore/wsio"
//
//	init() {
//		// set socket io handler and signalings
//		wsio.SetHandlers(ctrl.Authenticate, nil, nil)
//		adaptor := &ctrl.DefSioAdaptor{}
//		if err := wsio.SetAdapter(adaptor); err != nil {
//			panic(err)
//		}
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
type wingSIO struct {
	// Mutex sync lock, protect client connecting
	lock sync.Mutex

	// socket server
	server *sio.Server

	// socket golbal handler to execute clients authenticate action.
	authHandler AuthHandler

	// socket golbal handler to execute clients connect action.
	connHandler ConnectHandler

	// socket golbal handler to execute clients disconnect actions
	discHandler DisconnectHandler

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

// Socket connection server
var wsc *wingSIO

// Object logger with [SIO] perfix for socket.io module
var siolog = logger.NewLogger("SIO")

var (
	serverPingInterval = 30 * time.Second
	serverPingTimeout  = 60 * time.Second
	maxConnectCount    = 200000

	// Check client option if empty when connnection is established,
	// if optinal data is empty the connect will not establish and disconnect.
	usingOption = false
)

func init() {
	setupWsioConfigs()
	wsc = &wingSIO{
		options:   make(map[uintptr]*clientOpt),
		onceBunds: make(map[uintptr]string),
	}

	// set http handler for socke.io
	handler, err := wsc.createHandler()
	if err != nil {
		panic(err)
	}

	// set socket.io routers
	beego.Handler("/"+beego.BConfig.AppName+"/socket.io", handler)
	siolog.I("Initialized routers...")
}

// read wsio configs from file
func setupWsioConfigs() {
	interval := beego.AppConfig.DefaultInt64("wsio::interval", 30)
	serverPingInterval = time.Duration(interval) * time.Second

	timeout := beego.AppConfig.DefaultInt64("wsio::timeout", 60)
	serverPingTimeout = time.Duration(timeout) * time.Second

	using := beego.AppConfig.DefaultBool("wsio::optinal", false)
	usingOption = using

	// logout the configs value
	siolog.I("Server configs interval:", interval,
		"timeout:", timeout, "optional:", using)
}

// createHandler create http handler for socket.io
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

	siolog.I("Created handler")
	return server, nil
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
		// Use auth header function for Python3 and Unreal client
		token = req.Header.Get("Token")
		if !strings.HasPrefix(author, "WENGOLD") || token == "" {
			siolog.E("Invalid authoration:", author, "token:", token)
			return invar.ErrAuthDenied
		}
	} else {
		// Use URL + token string for React frontend and Wechat app client
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
	if cc.authHandler != nil {
		uid, opt, err := cc.authHandler(token)
		if err != nil || uid == "" {
			siolog.E("Invalid uid:", uid, "or case err:", err)
			return invar.ErrAuthDenied
		} else if usingOption && opt == "" {
			siolog.E("Empty client", uid, "option data!")
			return invar.ErrAuthDenied
		}

		siolog.I("Decoded client token, uuid:", uid, "opt:", opt)
		uuid, option = uid, opt
	}

	// bind http.Request -> uuid
	h := uintptr(unsafe.Pointer(req))
	cc.bindHTTP2UUIDLocked(h, uuid, option)
	return nil
}

// onConnect event of connect
func (cc *wingSIO) onConnect(sc sio.Socket) {
	// found client uuid and unbind -> http.Request
	h := uintptr(unsafe.Pointer(sc.Request()))
	if _, ok := cc.onceBunds[h]; ok /* already bund */ {
		siolog.W("Duplicate onConnect, abort for", h)
		return
	}
	cc.onceBunds[h] = "" // cache the first time

	co := cc.unbindUUIDFromHTTPLocked(h)
	if co == nil || co.UID == "" {
		siolog.E("Invalid socket request bind!")
		sc.Disconnect()
		return
	}

	clientPool := clients.Clients()
	if err := clientPool.Register(co.UID, sc, co.Opt); err != nil {
		siolog.E("Failed register socket client:", co.UID)
		sc.Disconnect()
		return
	}

	// handle connect callback for socket with uuid
	if cc.connHandler != nil {
		if err := cc.connHandler(co.UID, co.Opt); err != nil {
			siolog.E("Client:", co.UID, "connect socket err:", err)
			sc.Disconnect()
		}
	}
	siolog.I("Connected client:", co.UID)
	go cc.clearBundCache(h)
}

// onDisconnected event of disconnect
func (cc *wingSIO) onDisconnected(sc sio.Socket) {
	uuid, opt := clients.Clients().Deregister(sc)
	if cc.discHandler != nil {
		cc.discHandler(uuid, opt)
	}
}

// bindHTTP2UUIDLocked bind http request pointer -> uuid on locked status
func (cc *wingSIO) bindHTTP2UUIDLocked(h uintptr, uuid, opt string) {
	cc.lock.Lock()
	defer cc.lock.Unlock()

	cc.options[h] = &clientOpt{UID: uuid, Opt: opt}
}

// unbindUUIDFromHTTPLocked unbind uuid -> http request pointer on locked status
func (cc *wingSIO) unbindUUIDFromHTTPLocked(h uintptr) *clientOpt {
	cc.lock.Lock()
	defer cc.lock.Unlock()

	if data, ok := cc.options[h]; ok {
		co := &clientOpt{UID: data.UID, Opt: data.Opt}
		delete(cc.options, h)
		return co
	}
	return nil
}

// Clear the bind cache after 10ms
func (cc *wingSIO) clearBundCache(h uintptr) {
	time.Sleep(50 * time.Millisecond)
	delete(cc.onceBunds, h)
}
