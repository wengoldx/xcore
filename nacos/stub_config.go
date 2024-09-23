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
	"encoding/json"

	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"github.com/wengoldx/xcore/invar"
)

// Nacos config client stub
type ConfigStub struct {
	Namespace   string                      // Namespace id, it defined on console at first
	NacosServer string                      // Nacos server host ip
	Stub        config_client.IConfigClient // Nacos configs client instance
}

// Nacos config keyword datas
type ConfigKey struct {
	DataID string // config data id
	Group  string // config group name
}

// Callback to listen config changed
type ListenCallback func(namespace, group, dataId, data string)

// Generate a ConfigStub instance
//	@params ns  string config namespace id
//	@params svr string remote nacos server ip address
//	@return - *ConfigStub config stub instance
func NewConfigStub(ns, svr string) *ConfigStub {
	return &ConfigStub{Namespace: ns, NacosServer: svr}
}

// Setup config stub instance, it must set stub namespace and
// nacos server values before call this function
func (c *ConfigStub) Setup() error {
	if c.Namespace == "" || c.NacosServer == "" {
		return invar.ErrUnperparedState
	}

	ncp := genClientParam(c.Namespace, c.NacosServer)
	client, err := clients.NewConfigClient(ncp)
	if err != nil {
		return err
	}

	c.Stub = client
	return nil
}

// Publish config to nacos server, just support string and json struct datas
//	@params did    string    data id of nacos defined
//	@params group  string    group name of nacos defined
//	@params config interface string or struct type config content
//	@return - error handle exception
func (c *ConfigStub) Publish(did, group string, config any) error {
	if c.Stub == nil {
		return invar.ErrInvalidClient
	} else if did == "" || group == "" || config == nil {
		return invar.ErrInvalidParams
	}

	content := ""
	switch config := config.(type) {
	case string:
		content = config
	default:
		cfg, err := json.Marshal(config)
		if err != nil {
			return err
		}
		content = string(cfg)
	}

	// config key = dataId+group+namespaceId
	cp := vo.ConfigParam{DataId: did, Group: group, Content: content}
	_, err := c.Stub.PublishConfig(cp)
	return err
}

// Listen same one config changed, and notify to update
//	@params did   string         data id of nacos defined
//	@params group string         group name of nacos defined
//	@params cb    ListenCallback config changed callback
//	@return - error handle exception
func (c *ConfigStub) Listen(did, group string, cb ListenCallback) error {
	if c.Stub == nil {
		return invar.ErrInvalidClient
	} else if did == "" || group == "" || cb == nil {
		return invar.ErrInvalidParams
	}

	cp := vo.ConfigParam{DataId: did, Group: group, OnChange: cb}
	return c.Stub.ListenConfig(cp)
}

// Listen multiple configs changed, and using same callback to update
//	@params keys []ConfigKey    configs keywords
//	@params cb   ListenCallback config changed callback
//	@return - error handle exception
func (c *ConfigStub) Listens(keys []ConfigKey, cb ListenCallback) error {
	for _, key := range keys {
		if err := c.Listen(key.DataID, key.Group, cb); err != nil {
			return err
		}
	}
	return nil
}

// Cancel listen configs changed
//	@params keys []ConfigKey configs keywords
//	@return - error handle exception
func (c *ConfigStub) CancelListens(keys []ConfigKey) error {
	for _, key := range keys {
		cp := vo.ConfigParam{DataId: key.DataID, Group: key.Group}
		if err := c.Stub.CancelListenConfig(cp); err != nil {
			return err
		}
	}
	return nil
}

// Get string config content
//	@params did   string data id of nacos defined
//	@params group string group name of nacos defined
//	@return - string config content string
//			- error  handle exception
func (c *ConfigStub) GetString(did, group string) (string, error) {
	if c.Stub == nil {
		return "", invar.ErrInvalidClient
	} else if did == "" || group == "" {
		return "", invar.ErrInvalidParams
	}

	cp := vo.ConfigParam{DataId: did, Group: group}
	return c.Stub.GetConfig(cp)
}

// Get struct config content from josn string
//	@params did   string data id of nacos defined
//	@params group string group name of nacos defined
//	@return - out   config struct data
//			- error handle exception
func (c *ConfigStub) GetStruct(did, group string, out any) error {
	if content, err := c.GetString(did, group); err != nil {
		return err
	} else {
		if err = json.Unmarshal([]byte(content), out); err != nil {
			return err
		}
	}
	return nil
}

// Delete registered config
//	@params did   string data id of nacos defined
//	@params group string group name of nacos defined
//	@return - error handle exception
func (c *ConfigStub) Delete(did, group string) error {
	if c.Stub == nil {
		return invar.ErrInvalidClient
	} else if did == "" || group == "" {
		return invar.ErrInvalidParams
	}

	cp := vo.ConfigParam{DataId: did, Group: group}
	_, err := c.Stub.DeleteConfig(cp)
	return err
}

// Delete registered configs
//	@params keys []ConfigKey configs keywords
//	@return - error handle exception
func (c *ConfigStub) Deletes(keys []ConfigKey) error {
	for _, key := range keys {
		if err := c.Delete(key.DataID, key.Group); err != nil {
			return err
		}
	}
	return nil
}
