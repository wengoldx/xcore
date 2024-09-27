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
	acc "github.com/wengoldx/xcore/wrpc/accservice/proto"
	mea "github.com/wengoldx/xcore/wrpc/measure/proto"
	webss "github.com/wengoldx/xcore/wrpc/webss/proto"
	chat "github.com/wengoldx/xcore/wrpc/wgchat/proto"
	pay "github.com/wengoldx/xcore/wrpc/wgpay/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	SvrAcc  = "accservice" // server name of AccService backend
	SvrMea  = "measure"    // server name of Measure    backend
	SvrWss  = "webss"      // server name of WebSS      backend
	SvrChat = "wgchat"     // server name of WgChat     backend
	SvrPay  = "wgpay"      // server name of WgPay      backend
)

// GrpcHandlerFunc grpc server handler for register
type GrpcHandlerFunc func(svr *grpc.Server)

type GrpcStub struct {
	Certs   map[string]*nacos.GrpcCert // Grpc handler certs
	Clients map[string]any             // Grpc client handlers

	// Current grpc server if registried
	isRegistried bool

	// Global handler function to return grpc server handler
	SvrHandlerFunc GrpcHandlerFunc
}

// Singleton grpc stub instance
var grpcStub *GrpcStub

// Return Grpc global singleton
func Singleton() *GrpcStub {
	if grpcStub == nil {
		grpcStub = &GrpcStub{
			isRegistried: false,
			Certs:        make(map[string]*nacos.GrpcCert),
			Clients:      make(map[string]any),
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
		logger.E("Not setup global grpc handler!")
		return
	} else if stub.isRegistried {
		return // drop the duplicate registry
	}

	svrname := beego.BConfig.AppName
	logger.I("Register Grpc server:", svrname)

	secure, ok := stub.Certs[svrname]
	if !ok || secure.Key == "" || secure.Pem == "" {
		logger.E("Not found", svrname, "grpc cert, abort register!")
		return
	}

	// load grpc grpc server local port to listen
	port := beego.AppConfig.String("nacosport")
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		logger.E("Listen grpc server, err:", err)
		return
	}

	// generate TLS cert from pem datas
	cert, err := tls.X509KeyPair([]byte(secure.Pem), []byte(secure.Key))
	if err != nil {
		logger.E("Create grpc cert, err:", err)
		return
	}

	// generate grpc server handler with TLS secure
	cred := credentials.NewServerTLSFromCert(&cert)
	svr := grpc.NewServer(grpc.Creds(cred))
	stub.SvrHandlerFunc(svr)
	logger.I("Running Grpc server:", svrname, "on port", port)

	stub.isRegistried = true
	defer func(stub *GrpcStub) { stub.isRegistried = false }(stub)
	if err := svr.Serve(lis); err != nil {
		logger.E("Start grpc server, err:", err)
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

// Generate grpc client handler
func (stub *GrpcStub) GenClient(svrkey, addr string, port int) {
	if svrkey != SvrAcc && svrkey != SvrMea && svrkey != SvrWss &&
		svrkey != SvrChat && svrkey != SvrPay {
		logger.E("Invaoid target grpc server:", svrkey)
		return
	}

	secure, ok := stub.Certs[svrkey]
	if !ok || secure.Key == "" || secure.Pem == "" {
		logger.E("Not found target grpc cert of", svrkey)
		return
	}

	// generate TLS cert from pem datas
	cp := x509.NewCertPool()
	if !cp.AppendCertsFromPEM([]byte(secure.Pem)) {
		logger.E("Failed generate grpc cert!")
		return
	}

	// generate grpc client handler with TLS secure
	grpcsvr := fmt.Sprintf("%s:%d", addr, port)
	cred := credentials.NewClientTLSFromCert(cp, svrkey)
	conn, err := grpc.Dial(grpcsvr, grpc.WithTransportCredentials(cred))
	if err != nil {
		logger.E("dial grpc address", grpcsvr, " fialed", err)
		return
	}

	// content grpc client by server name
	logger.I("Grpc client:", svrkey, "connect", grpcsvr)
	switch svrkey {
	case SvrAcc:
		stub.Clients[svrkey] = acc.NewAccClient(conn)
	case SvrMea:
		stub.Clients[svrkey] = mea.NewMeaClient(conn)
	case SvrWss:
		stub.Clients[svrkey] = webss.NewWebssClient(conn)
	case SvrChat:
		stub.Clients[svrkey] = chat.NewWgchatClient(conn)
	case SvrPay:
		stub.Clients[svrkey] = pay.NewWgpayClient(conn)
	}
}

// Parse all grpc certs from nacos config data, and cache to certs map
func (stub *GrpcStub) ParseCerts(data string) error {
	certs := nacos.GrpcCerts{}
	if err := xml.Unmarshal([]byte(data), &certs); err != nil {
		logger.E("Parse grpc certs, err:", err)
		return err
	}

	for _, cert := range certs.Certs {
		logger.D("Update", cert.Svr, "grpc cert")
		stub.Certs[cert.Svr] = &nacos.GrpcCert{
			Svr: cert.Svr, Key: cert.Key, Pem: cert.Pem,
		}
	}
	return nil
}

// ----------------------------------------
// Get Singleton Instanced GRPC Clients
// ----------------------------------------

// Return AccService grpc client, maybe null if not generate first
func (stub *GrpcStub) Acc() acc.AccClient {
	if client, ok := stub.Clients[SvrAcc]; ok {
		return client.(acc.AccClient)
	}
	return nil
}

// Return Measure grpc client, maybe null if not generate first
func (stub *GrpcStub) Mea() mea.MeaClient {
	if client, ok := stub.Clients[SvrMea]; ok {
		return client.(mea.MeaClient)
	}
	return nil
}

// Return Webss grpc client, maybe null if not generate first
func (stub *GrpcStub) Webss() webss.WebssClient {
	if client, ok := stub.Clients[SvrWss]; ok {
		return client.(webss.WebssClient)
	}
	return nil
}

// Return Wgchat grpc client, maybe null if not generate first
func (stub *GrpcStub) Chat() chat.WgchatClient {
	if client, ok := stub.Clients[SvrChat]; ok {
		return client.(chat.WgchatClient)
	}
	return nil
}

// Return Wgpay grpc client, maybe null if not generate first
func (stub *GrpcStub) Pay() pay.WgpayClient {
	if client, ok := stub.Clients[SvrPay]; ok {
		return client.(pay.WgpayClient)
	}
	return nil
}

// ----------------------------------------
// Account Authentications Request
// ----------------------------------------

// Auth header token and return account uuid and password
func (stub *GrpcStub) AuthHeaderToken(token string) (string, string) {
	if accGRPC := stub.Acc(); accGRPC == nil {
		logger.E("Acc RPC instance not inited!")
		return "", ""
	} else {
		param := &acc.Token{Token: token}
		resp, err := accGRPC.ViaToken(context.Background(), param)
		if err != nil {
			logger.E("RPC auth token, err:", err)
			return "", ""
		}

		return resp.Acc /* uuid */, ""
	}
}

// Auth account role from http header
func (stub *GrpcStub) AuthHeaderRole(uuid, url, method string) bool {
	if accGRPC := stub.Acc(); accGRPC == nil {
		logger.E("Acc RPC instance not inited!")
		return false
	} else {
		param := &acc.Role{Uuid: uuid, Router: url, Method: method}
		resp, err := accGRPC.ViaRole(context.Background(), param)
		if err != nil {
			logger.E("RPC auth", uuid, "role, err:", err)
			return false
		}

		return resp.Pass
	}
}

// ----------------------------------------
// GRPC Local Service Setup
// ----------------------------------------

// Parse certs to running as grpc server and update swagger routers,
// it will listen target services if you input them infos.
func SetupAsServer(mc *nacos.MetaConfig, data string, servers ...*nacos.ServerItem) {
	// Parse grpc certs and start as grpc server handler
	Singleton().ParseAndStart(data)

	// Register server to nacos and listen tags for grpc
	if len(servers) > 0 {
		nacos.RegisterServer().ListenServers(servers)
	}
	mc.UploadRouters()
}

// Parse certs to running as grpc client and update swagger routers,
// it will listen target services if you input them infos.
func SetupAsClient(mc *nacos.MetaConfig, data string, servers ...*nacos.ServerItem) {
	// Parse grpc certs and start as grpc server handler
	Singleton().ParseCerts(data)

	// Register server to nacos and listen tags for grpc
	if len(servers) > 0 {
		nacos.RegisterServer().ListenServers(servers)
	}
	mc.UploadRouters()
}
