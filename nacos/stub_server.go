// Copyright (c) 2018-Now Dunyu All Rights Reserved.
//
// Author      : https://www.wengold.net
// Email       : support@wengold.net
//
// Prismy.No | Date       | Modified by. | Description
// -------------------------------------------------------------------
// 00001       2022/05/11   yangping     New version
// -------------------------------------------------------------------

package nacos

import (
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/model"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"github.com/wengoldx/xcore/invar"
)

// Nacos naming client stub
type ServerStub struct {
	Namespace   string                      // Namespace id, it defined on console at first
	NacosServer string                      // Nacos server host ip
	Stub        naming_client.INamingClient // Nacos naming client instance
}

// Callback to listen server register status changed
type SubscribeCallback func(services []model.Instance, err error)

// Generate a ServerStub instance
//	@params ns  string server namespace id
//	@params svr string remote nacos server ip address
//	@return - *ServerStub server stub instance
func NewServerStub(ns, svr string) *ServerStub {
	return &ServerStub{Namespace: ns, NacosServer: svr}
}

// Setup server stub instance, it must set stub namespace and
// nacos server values before call this function
func (s *ServerStub) Setup() error {
	if s.Namespace == "" || s.NacosServer == "" {
		return invar.ErrUnperparedState
	}

	ncp := genClientParam(s.Namespace, s.NacosServer)
	client, err := clients.NewNamingClient(ncp)
	if err != nil {
		return err
	}

	s.Stub = client
	return nil
}

// Register business server into nacos remote server
//	@params name string   business server name
//	@params host string   business server deploied ip or domain
//	@params port uint64   business server port, must over 3000
//	@params opts []string 0:group name, 1:cluster name of business server
//	@return - error handle exception
func (s *ServerStub) Register(name, host string, port uint64, opts ...string) error {
	if s.Stub == nil {
		return invar.ErrInvalidClient
	} else if name == "" || host == "" || port < 3000 {
		return invar.ErrInvalidParams
	}

	dip := vo.RegisterInstanceParam{
		Ip: host, Port: port, ServiceName: name,
		Weight: 10, Enable: true, Healthy: true, Ephemeral: true,
	}

	if cnt := len(opts); cnt > 0 {
		if opts[0] != "" /* group name */ {
			dip.GroupName = opts[0]
		}
		if cnt > 1 && opts[1] != "" /* cluster name */ {
			dip.ClusterName = opts[1]
		}
	}

	if rst, err := s.Stub.RegisterInstance(dip); err != nil {
		return err
	} else if !rst {
		return invar.ErrInvalidState
	}
	return nil
}

// Deregister business server out of nacos remote server
//	@params name string   business server name
//	@params host string   business server deploied ip or domain
//	@params port uint64   business server port, must over 3000
//	@params opts []string 0:group name, 1:cluster name of business server
//	@return - error handle exception
func (s *ServerStub) Deregister(name, host string, port uint64, opts ...string) error {
	if s.Stub == nil {
		return invar.ErrInvalidClient
	} else if name == "" || host == "" || port < 3000 {
		return invar.ErrInvalidParams
	}

	dip := vo.DeregisterInstanceParam{
		Ip: host, Port: port, ServiceName: name,
		Ephemeral: true, //it must be true
	}

	if cnt := len(opts); cnt > 0 {
		if opts[0] != "" /* group name */ {
			dip.GroupName = opts[0]
		}
		if cnt > 1 && opts[1] != "" /* cluster name */ {
			dip.Cluster = opts[1]
		}
	}

	if rst, err := s.Stub.DeregisterInstance(dip); err != nil {
		return err
	} else if !rst {
		return invar.ErrInvalidState
	}
	return nil
}

// Get business server registry informations from nacos remote server
//	@params name string   business server name
//	@params opts []string 0:group name, 1~n:clusters name of business server
//	@return - error handle exception
func (s *ServerStub) GetServer(name string, opts ...string) (*model.Service, error) {
	if s.Stub == nil {
		return nil, invar.ErrInvalidClient
	} else if name == "" {
		return nil, invar.ErrInvalidParams
	}

	sp := vo.GetServiceParam{ServiceName: name}
	if cnt := len(opts); cnt > 0 {
		if opts[0] != "" /* group name */ {
			sp.GroupName = opts[0]
		}
		if cnt > 1 && opts[1] != "" /* clusters, just check first one */ {
			sp.Clusters = opts[1:]
		}
	}

	if service, err := s.Stub.GetService(sp); err != nil {
		return nil, err
	} else {
		return &service, nil
	}
}

// Subscribe business server registry changed event
//	@params name string            business server name
//	@params cb   SubscribeCallback server registry changed callback
//	@params opts []string          0:group name, 1~n:clusters name of business server
//	@return - error handle exception
func (s *ServerStub) Subscribe(name string, cb SubscribeCallback, opts ...string) error {
	if s.Stub == nil {
		return invar.ErrInvalidClient
	} else if name == "" || cb == nil {
		return invar.ErrInvalidParams
	}

	sp := s.genSubscribeParam(name, cb, opts...)
	if err := s.Stub.Subscribe(sp); err != nil {
		return err
	}
	return nil
}

// Unsubscribe business server registry changed event
//	@params name string            business server name
//	@params cb   SubscribeCallback server registry changed callback
//	@params opts []string          0:group name, 1~n:clusters name of business server
//	@return - error handle exception
func (s *ServerStub) Unsubscribe(name string, cb SubscribeCallback, opts ...string) error {
	if s.Stub == nil {
		return invar.ErrInvalidClient
	} else if name == "" || cb == nil {
		return invar.ErrInvalidParams
	}

	sp := s.genSubscribeParam(name, cb, opts...)
	if err := s.Stub.Unsubscribe(sp); err != nil {
		return err
	}
	return nil
}

// Generate subscribe param for server registry changed events
func (s *ServerStub) genSubscribeParam(name string, cb SubscribeCallback, opts ...string) *vo.SubscribeParam {
	param := &vo.SubscribeParam{ServiceName: name, SubscribeCallback: cb}
	if cnt := len(opts); cnt > 0 {
		if opts[0] != "" /* group name */ {
			param.GroupName = opts[0]
		}
		if cnt > 1 && opts[1] != "" /* clusters, just check first one */ {
			param.Clusters = opts[1:]
		}
	}
	return param
}
