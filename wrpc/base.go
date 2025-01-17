// Copyright (c) 2018-Now Dunyu All Rights Reserved.
//
// Author      : https://www.wengold.net
// Email       : support@wengold.net
//
// Prismy.No | Date       | Modified by. | Description
// -------------------------------------------------------------------
// 00001       2019/05/22   yangping       New version
// -------------------------------------------------------------------

package wrpc

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/xml"
	"fmt"
	"net"

	"github.com/astaxie/beego"
	"github.com/wengoldx/xcore/logger"
	"github.com/wengoldx/xcore/nacos"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	acc "github.com/wengoldx/xcore/wrpc/accservice/proto"
	mea "github.com/wengoldx/xcore/wrpc/measure/proto"
	wss "github.com/wengoldx/xcore/wrpc/webss/proto"
	chat "github.com/wengoldx/xcore/wrpc/wgchat/proto"
	pay "github.com/wengoldx/xcore/wrpc/wgpay/proto"
)

const (
	SVR_ACC  = "accservice" // server name of AccService backend
	SVR_MEA  = "measure"    // server name of Measure    backend
	SVR_WSS  = "webss"      // server name of WebSS      backend
	SVR_CHAT = "wgchat"     // server name of WgChat     backend
	SVR_PAY  = "wgpay"      // server name of WgPay      backend
)

// GrpcHandlerFunc grpc server handler for register
type GrpcHandlerFunc func(svr *grpc.Server)

type GrpcStub struct {
	// Grpc handler certs
	Certs map[string]*nacos.GrpcCert

	// Current grpc server if registried
	isRegistried bool

	// Global handler function to return grpc server handler
	SvrHandlerFunc GrpcHandlerFunc

	// GRPC agent clients, when server register and listen to nacos
	Acc  acc.AccClient     // SVR_ACC  : Acc        GRPC client, maybe null
	Mea  mea.MeaClient     // SVR_MEA  : Measure    GRPC client, maybe null
	Wss  wss.WebssClient   // SVR_WSS  : Web static GRPC client, maybe null
	Chat chat.WgchatClient // SVR_CHAT : Chat       GRPC client, maybe null
	Pay  pay.WgpayClient   // SVR_PAY  : Pay        GRPC client, maybe null
}

// Singleton grpc stub instance
var grpcStub *GrpcStub

// Object logger with [RPC] mark for GRPC module
var rpclog = logger.CatLogger("RPC")

// Return Grpc global singleton
func Singleton() *GrpcStub {
	if grpcStub == nil {
		grpcStub = &GrpcStub{
			isRegistried: false, Certs: make(map[string]*nacos.GrpcCert),
		}
	}
	return grpcStub
}

// Start and excute grpc server, you numst setup global grpc
// register handler first as follow.
//
// `USAGE`
//
//	// set grpc server register handler
//	stub := wrpc.Singleton()
//	stub.SvrHandlerFunc = func(svr *grpc.Server) {
//		proto.RegisterAccServer(svr, &(handler.Acc{}))
//	}
//
//	// parse grps certs before register
//	stub.ParseCerts(data)
//
//	// register local server as grpc server
//	go stub.StartGrpcServer()
func (stub *GrpcStub) StartGrpcServer() {
	if stub.SvrHandlerFunc == nil {
		rpclog.E("Not setup grpc handler!")
		return
	} else if stub.isRegistried {
		return // drop the duplicate registry
	}

	svrname := beego.BConfig.AppName
	rpclog.I("Register grpc server:", svrname)

	secure, ok := stub.Certs[svrname]
	if !ok || secure.Key == "" || secure.Pem == "" {
		rpclog.E("Not found", svrname, "grpc cert, abort register!")
		return
	}

	// load grpc grpc server local port to listen
	port := beego.AppConfig.String("nacosport")
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		rpclog.E("Listen grpc server, err:", err)
		return
	}

	// generate TLS cert from pem datas
	cert, err := tls.X509KeyPair([]byte(secure.Pem), []byte(secure.Key))
	if err != nil {
		rpclog.E("Create grpc cert, err:", err)
		return
	}

	// generate grpc server handler with TLS secure
	cred := credentials.NewServerTLSFromCert(&cert)
	svr := grpc.NewServer(grpc.Creds(cred))
	stub.SvrHandlerFunc(svr)
	rpclog.I("Running grpc server:", svrname, "on port", port)

	stub.isRegistried = true
	defer func(stub *GrpcStub) { stub.isRegistried = false }(stub)
	if err := svr.Serve(lis); err != nil {
		rpclog.E("Start grpc server, err:", err)
	}
}

// Parse grpc certs from nacos configs and register local server
// as grpc server handler, then start and listen.
func (stub *GrpcStub) ParseAndStart(data string) error {
	if err := stub.ParseCerts(data); err != nil {
		return err
	}

	go stub.StartGrpcServer()
	return nil
}

// Deprecated: use wrpc.NewClient instead it.
func (stub *GrpcStub) GenClient(svrkey, addr string, port int) { stub.NewClient(svrkey, addr, port) }

// Create a new grpc client as stub handler.
func (stub *GrpcStub) NewClient(svrkey, addr string, port int) {
	if svrkey != SVR_ACC && svrkey != SVR_MEA && svrkey != SVR_WSS &&
		svrkey != SVR_CHAT && svrkey != SVR_PAY {
		rpclog.E("Invalid target grpc server:", svrkey)
		return
	}

	secure, ok := stub.Certs[svrkey]
	if !ok || secure.Key == "" || secure.Pem == "" {
		rpclog.E("Not found target grpc cert of", svrkey)
		return
	}

	// generate TLS cert from pem datas
	cp := x509.NewCertPool()
	if !cp.AppendCertsFromPEM([]byte(secure.Pem)) {
		rpclog.E("Failed generate grpc cert!")
		return
	}

	// generate grpc client handler with TLS secure
	grpcsvr := fmt.Sprintf("%s:%d", addr, port)
	cred := credentials.NewClientTLSFromCert(cp, svrkey)
	conn, err := grpc.Dial(grpcsvr, grpc.WithTransportCredentials(cred))
	if err != nil {
		rpclog.E("Dial grpc address", grpcsvr, " fialed", err)
		return
	}

	// content grpc client by server name
	switch svrkey {
	case SVR_ACC:
		stub.Acc = acc.NewAccClient(conn)
	case SVR_MEA:
		stub.Mea = mea.NewMeaClient(conn)
	case SVR_WSS:
		stub.Wss = wss.NewWebssClient(conn)
	case SVR_CHAT:
		stub.Chat = chat.NewWgchatClient(conn)
	case SVR_PAY:
		stub.Pay = pay.NewWgpayClient(conn)
	default:
		rpclog.E("Invalid grpc server:", svrkey+"@"+grpcsvr)
		return
	}
	rpclog.I("Connected with", svrkey+"@"+grpcsvr)
}

// Parse all grpc certs from nacos config data, and cache to certs map
func (stub *GrpcStub) ParseCerts(data string) error {
	certs := nacos.GrpcCerts{}
	if err := xml.Unmarshal([]byte(data), &certs); err != nil {
		rpclog.E("Parse grpc certs, err:", err)
		return err
	}

	for _, cert := range certs.Certs {
		rpclog.D("Update", cert.Svr, "grpc cert")
		stub.Certs[cert.Svr] = &nacos.GrpcCert{
			Svr: cert.Svr, Key: cert.Key, Pem: cert.Pem,
		}
	}
	return nil
}

// ----------------------------------------
// Account Authentications Request
// ----------------------------------------

// Auth header token and return account uuid.
//
// `NOTICE` that the returned password fixed as empty string for services
// except backend account server.
func (stub *GrpcStub) AuthHeaderToken(token string) (string, string) {
	if stub.Acc == nil || token == "" {
		rpclog.E("Acc grpc uninit, or invalie token!")
		return "", ""
	}

	param := &acc.Token{Token: token}
	resp, err := stub.Acc.ViaToken(context.Background(), param)
	if err != nil {
		rpclog.E("Auth grpc token, err:", err)
		return "", ""
	}
	return resp.Acc /* uuid */, "" /* password */
}

// Auth account role from http header
func (stub *GrpcStub) AuthHeaderRole(uuid, url, method string) bool {
	if stub.Acc == nil || uuid == "" || url == "" || method == "" {
		rpclog.E("Acc grpc uninit, or invalid role!")
		return false
	}

	param := &acc.Role{Uuid: uuid, Router: url, Method: method}
	resp, err := stub.Acc.ViaRole(context.Background(), param)
	if err != nil {
		rpclog.E("Auth", uuid, "role, err:", err)
		return false
	}
	return resp.Pass
}

// ----------------------------------------
// GRPC Local Service Setup
// ----------------------------------------

// Target services listing callback for create grpc client after registred.
func listingCallback(svr string, addr string, port int) {
	Singleton().NewClient(svr, addr, port)
}

// Register server to nacos and listen tags for grpc
func registryAndUploadRouters(mc *nacos.MetaConfig, servers ...string) {
	svrstub := nacos.RegisterServer()
	if len(servers) > 0 {
		svrs := []*nacos.ServerItem{}
		for _, sn := range servers {
			svr := &nacos.ServerItem{Name: sn, Callback: listingCallback}
			svrs = append(svrs, svr)
		}
		svrstub.ListenServers(svrs)
	}
	mc.UploadRouters()
}

// Parse certs to running as grpc server and update swagger routers,
// it will listen target services if you input them infos.
//
// ---
//
//	stub := wrpc.Singleton()
//	stub.ParseAndStart(data) // stub.ParseCerts(data);
//							 // go stub.StartGrpcServer()
//
//	mvc.registryAndUploadRouters(mc, servers...)
//	// svrstub := nacos.RegisterServer()
//	// svrstub.ListenServers(servers...)
//	// mc.UploadRouters()
//
// ---
func SetupAsServer(mc *nacos.MetaConfig, data string, servers ...string) *GrpcStub {
	// Parse grpc certs and start as grpc server handler
	Singleton().ParseAndStart(data)

	// Register server to nacos and listen tags for grpc
	registryAndUploadRouters(mc, servers...)
	return grpcStub
}

// Parse certs to running as grpc client and update swagger routers,
// it will listen target services if you input them infos.
//
// ---
//
//	wrpc.Singleton().ParseCerts(data);
//
//	mvc.registryAndUploadRouters(mc, servers...)
//	// svrstub := nacos.RegisterServer()
//	// svrstub.ListenServers(servers...)
//	// mc.UploadRouters()
//
// ---
func SetupAsClient(mc *nacos.MetaConfig, data string, servers ...string) *GrpcStub {
	// Parse grpc certs and start as grpc server handler
	Singleton().ParseCerts(data)

	// Register server to nacos and listen tags for grpc
	registryAndUploadRouters(mc, servers...)
	return grpcStub
}
