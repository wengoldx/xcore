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
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/model"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"github.com/wengoldx/xcore/invar"
	"github.com/wengoldx/xcore/logger"
	"github.com/wengoldx/xcore/utils"
)

// -------- Auto Register Define --------

// Server register informations
type ServerItem struct {
	Name     string         // Server name, same as beego app name
	Callback ServerCallback // Server register datas changed callback
}

// Callback to listen server address and port changes
type ServerCallback func(svr, addr string, gport, hport int)

// Register current server to nacos, you must set configs in app.conf
//	@return - *ServerStub nacos server stub instance
//
//	`NOTICE` : nacos config as follows.
//
// ----
//
//	; Nacos remote server host
//	nacossvr = "10.239.40.24"
//
//	[dev]
//	; Inner net ideal address for dev servers access
//	nacosaddr = "10.239.20.99"
//
//	; Inner net port for grpc access
//	nacosport = 3000
//
//	[prod]
//	; Inner net ideal address for prod servers access
//	nacosaddr = "10.239.40.64"
//
//	; Inner net port for grpc access
//	nacosport = 3000
func RegisterServer() *ServerStub {
	svr := beego.AppConfig.String(configKeySvr)
	if svr == "" {
		panic("Not found nacos server host!")
	}

	// Local server listing ip proxy
	idealip := beego.AppConfig.String(configKeyAddr)
	if idealip == "" {
		panic("Not found idea server ip to register!")
	}

	// Server access port for grpc, it maybe same as httpport config
	// when the local server not support grpc but for http
	port, err := beego.AppConfig.Int(configKeyPort)
	if err != nil || port < 3000 /* remain 0 ~ 3000 */ {
		panic("Not found port number or less 3000!")
	}

	// Namespace id of local server, and parse local ip
	ns := utils.Condition(beego.BConfig.RunMode == "prod", NS_PROD, NS_DEV).(string)
	addr, err := matchProxyIP(idealip)
	if err != nil {
		panic("Find proxy local ip, err:" + err.Error())
	}

	// Generate nacos server stub and setup it
	stub := NewServerStub(ns, svr)
	if err := stub.Setup(); err != nil {
		panic(err)
	}

	// Fixed app name as nacos server name to register,
	// and pick server port from config 'nacosport' not form 'httpport' value,
	// becase it maybe support either grpc or http hanlder to accesse.
	//
	// And here not use cluster name, please keep it empty!
	// And last FIX the group name as 'group.wengold'.
	app, gp := beego.BConfig.AppName, GP_WENGOLD
	if err := stub.Register(app, addr, uint64(port), gp); err != nil {
		panic(err)
	}

	logmsg := fmt.Sprintf("%s@%s:%v", app, addr, port)
	logger.I("Registered server on", logmsg)
	return stub
}

// Listing services address and port changes, it will call the callback
// immediately to return target service host when them allready registerd
// to service central of nacos.
//	@params servers []*ServerItem target server registry informations
func (ss *ServerStub) ListenServers(servers []*ServerItem) {
	for _, s := range servers {
		if err := ss.Subscribe(s.Name, s.OnChanged, GP_WENGOLD); err != nil {
			panic("Subscribe server " + s.Name + " err:" + err.Error())
		}
	}
}

// Subscribe callback called when target service address and port changed
func (si *ServerItem) OnChanged(services []model.Instance, err error) {
	if err != nil {
		logger.E("Received server", si.Name, "change, err:", err)
		return
	}

	if len(services) > 0 {
		addr, port := services[0].Ip, services[0].Port

		// Paser httpport from metadata map if it exist
		meta, httpport := services[0].Metadata, 0
		if meta != nil {
			if hp, ok := meta[configKeyHPort]; ok {
				httpport, _ = strconv.Atoi(hp)
			}
		}

		logmsg := fmt.Sprintf("%s@%s:%v - httpport:%v", si.Name, addr, port, httpport)
		logger.I("Update server to", logmsg)
		si.Callback(si.Name, addr, int(port), httpport)
	}
}

// -------- Config Service Define --------

// Meta config informations
type MetaConfig struct {
	Stub      *ConfigStub                   // Nacos config client instance
	Callbacks map[string]MetaConfigCallback // Meta config changed callback maps, key is dataid
}

// Callback to listen server address and port changes
type MetaConfigCallback func(dataId, data string)

// Generate a meta config client to get or listen configs changes
//	@return - *MetaConfig nacos config client instance
//
//	`NOTICE` : nacos config as follows.
//
// ----
//
//	; Nacos remote server host
//	nacossvr = "10.239.40.24"
func GenMetaConfig() *MetaConfig {
	svr := beego.AppConfig.String(configKeySvr)
	if svr == "" {
		panic("Not found nacos server host!")
	}

	// Namespace id of meta configs
	ns := utils.Condition(beego.BConfig.RunMode == "prod", NS_PROD, NS_DEV).(string)

	// Generate nacos config stub and setup it
	stub := NewConfigStub(ns, svr)
	if err := stub.Setup(); err != nil {
		panic("Gen config stub, err:" + err.Error())
	}

	// Fix the all config group as wengold
	cbs := make(map[string]MetaConfigCallback)
	return &MetaConfig{
		Stub: stub, Callbacks: cbs,
	}
}

// Get and listing the config of indicated dataId
func (mc *MetaConfig) ListenConfig(dataId string, cb MetaConfigCallback) {
	mc.Callbacks[dataId] = cb // cache callback
	gp := GP_WENGOLD

	// get config first before listing
	data, err := mc.Stub.GetString(dataId, gp)
	if err != nil {
		panic("Get config " + dataId + "err: " + err.Error())
	}
	cb(dataId, data)

	// listing config changes
	logger.I("Start listing config dateId:", dataId)
	mc.Stub.Listen(dataId, gp, mc.OnChanged)
}

// Get and listing the configs of indicated dataIds
func (mc *MetaConfig) ListenConfigs(dataIds []string, cb MetaConfigCallback) {
	for _, dataId := range dataIds {
		if dataId == "" || cb == nil {
			continue
		}
		mc.ListenConfig(dataId, cb)
	}
}

// Listing callback called when target configs changed
func (mc *MetaConfig) OnChanged(namespace, group, dataId, data string) {
	if namespace != NS_DEV && namespace != NS_PROD {
		logger.E("Invalid meta config ns:", namespace)
		return
	}

	if callback, ok := mc.Callbacks[dataId]; ok {
		logger.I("Update config dataId", dataId, "to:", data)
		callback(dataId, data)
	}
}

// Get config data from nacos server by given data id
func (mc *MetaConfig) GetConfig(dataId string) (string, error) {
	return mc.Stub.GetString(dataId, GP_WENGOLD)
}

// Push config data to indicated nacos config
func (mc *MetaConfig) PushConfig(dataId, data string) error {
	return mc.Stub.Publish(dataId, GP_WENGOLD, data)
}

// Parse local server swagger and upload routers to nacos
func (mc *MetaConfig) UploadRouters() error {
	nrouters, err := mc.GetConfig(DID_API_ROUTERS)
	if err != nil {
		logger.E("Pull nacos routers, err:", err)
		return err
	}

	// update local server swagger routers
	if routers, err := utils.UpdateRouters(nrouters); err == nil {
		mc.PushConfig(DID_API_ROUTERS, routers)
	}
	return nil
}

// Update routers chinese descriptions and upload to nacos
func (mc *MetaConfig) UpdateChineses(descs []*utils.SvrDesc) error {
	nrouters, err := mc.GetConfig(DID_API_ROUTERS)
	if err != nil {
		logger.E("Pull nacos routers, err:", err)
		return err
	}

	// update routers chineses descriptions
	routers, err := utils.UpdateChineses(nrouters, descs)
	if err != nil {
		logger.E("Update routers chineses, err:", err)
		return err
	}

	mc.PushConfig(DID_API_ROUTERS, routers)
	return nil
}

// ----------------------------------------

// Generate nacos client config, contain nacos remote server and
// current business servers configs, this client keep alive with
// 5s pingpong heartbeat and output logs on warn leven.
//
//	`NOTICE`
//
// - Remote direct nacos server need access on http://{svr}:8848/nacos
//
// - Nginx proxy vip server need access on http://{svr}:3608/nacos
func genClientParam(ns, svr string) vo.NacosClientParam {
	sc := []constant.ServerConfig{
		{Scheme: "http", ContextPath: "/nacos", IpAddr: svr, Port: 3608},
	}

	// logs config
	logcfg := &constant.ClientLogRollingConfig{
		MaxSize: 10, MaxBackups: 10, // max 10 files and each max 10MB
	}

	// client config
	cc := &constant.ClientConfig{
		NamespaceId:         ns,
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              nacosDirLogs,
		CacheDir:            nacosDirCache,
		LogRollingConfig:    logcfg,
		LogLevel:            nacosLogLevel,
		Username:            nacosSysSecure, // secure account
		Password:            nacosSysSecure, // secure passowrd
	}

	return vo.NacosClientParam{
		ClientConfig: cc, ServerConfigs: sc,
	}
}

// Parse and return the local register IP that meets the conditions
func matchProxyIP(proxy string) (string, error) {
	segments := strings.Split(proxy, ".")
	if len(segments) != 4 {
		logger.E("Invalid nocos proxy ip:", proxy)
		return "", invar.ErrInvalidParams
	}

	segment := `.((2((5[0-5])|([0-4]\d)))|([0-1]?\d{1,2}))`
	condition := strings.Join(segments[0:3], ".") + segment
	reg, err := regexp.Compile(condition)
	if err != nil {
		logger.E("Compile regular err:", err)
		return "", err
	}

	matchips, err := utils.GetLocalIPs()
	if err != nil {
		return "", err
	}

	satisfips := []string{}
	for _, v := range matchips {
		// Find ideal ip exists and return it if found
		ip := reg.FindString(v)
		if ip == proxy {
			logger.I("Direct use ideal ip:", proxy)
			return ip, nil
		} else if ip != "" {
			satisfips = append(satisfips, ip)
		}
	}

	if len(satisfips) > 0 {
		logger.I("Use dynamic ip:", satisfips[0])
		return satisfips[0], nil
	}

	logger.E("Not found any local ips, just return proxy:", proxy)
	return proxy, nil
}
