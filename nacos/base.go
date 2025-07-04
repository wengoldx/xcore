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
	"fmt"
	"regexp"
	"strings"

	"github.com/astaxie/beego"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/model"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"github.com/wengoldx/xcore/elastic"
	"github.com/wengoldx/xcore/invar"
	"github.com/wengoldx/xcore/logger"
	"github.com/wengoldx/xcore/mqtt"
	"github.com/wengoldx/xcore/utils"
	"github.com/wengoldx/xcore/wechat"
)

// Object logger with [NACOS] mark for nacos module
var naclog = logger.CatLogger("NACOS")

// ----------------------------------------
// Auto Register Define
// ----------------------------------------

// Server register informations
type ServerItem struct {
	Name     string         // Server name, same as beego app name
	Callback ServerCallback // Server register datas changed callback
}

// Callback to listen server address and port changes
type ServerCallback func(svr, addr string, port int)

// Register current server to nacos, you must set configs in app.conf
//	@return - *ServerStub nacos server stub instance
//
//	`NOTICE` : nacos config as follows.
//
// ----
//
//	; Nacos remote server host
//	nacossvr = "xx.xxx.40.218"
//
//	[dev]
//	; Inner net ideal address for dev servers access
//	nacosaddr = "xx.xxx.20.239"
//
//	; Inner net port for grpc access
//	nacosport = 3000
//
//	[prod]
//	; Inner net ideal address for prod servers access
//	nacosaddr = "xx.xxx.40.199"
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
	ns := utils.Condition(beego.BConfig.RunMode == "prod", NS_PROD, NS_DEV)
	addr, err := matchProxyIP(idealip)
	if err != nil {
		panic("Find proxy local ip, err:" + err.Error())
	}

	// Create nacos server stub and setup it
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
	naclog.I("Registered server on", logmsg)
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
		naclog.E("Received server", si.Name, "change, err:", err)
		return
	}

	if len(services) > 0 {
		addr, port := services[0].Ip, services[0].Port
		naclog.I("[GRPC] Server", si.Name, "address changed:", fmt.Sprintf("%s:%v", addr, port))
		si.Callback(si.Name, addr, int(port))
	}
}

// ----------------------------------------
// Config Service Define
// ----------------------------------------

// Meta config informations
type MetaConfig struct {
	Stub      *ConfigStub                   // Nacos config client instance
	Callbacks map[string]MetaConfigCallback // Meta config changed callback maps, key is dataid

	// Custom datas when project register data ids, see parseXxxx inner functions.
	Conf    AccConfs                     // DID_ACC_CONFIGS  : Account configs
	OTA     map[string]*OTAInfo          // DID_OTA_BUILDS   : Projects OTA infos
	Senders map[string]*DTalkSender      // DID_DTALK_NTFERS : DingTalk senders
	Agents  map[string]*wechat.WxIFAgent // DID_WX_AGENTS    : Wechat agents
	Paths   map[string][]*ResPath        // DID_MIO_PATHS    : MinIO export resource paths
	Users   map[string]string            // DID_MIO_USERS    : MinIO access service users
	Words   EWords                       // DID_QK_WORDS     : Excel rule words
}

// Callback to listen server address and port changes
type MetaConfigCallback func(dataId, data string)

// Create meta config client to get or listen configs changes
//	@return - *MetaConfig nacos config client instance
//
//	`NOTICE` : nacos config as follows.
//
// ----
//
//	; Nacos remote server host
//	nacossvr = "xx.xxx.40.218"
func NewMetaConfig() *MetaConfig {
	svr := beego.AppConfig.String(configKeySvr)
	if svr == "" {
		panic("Not found nacos server host!")
	}

	// Namespace id of meta configs
	ns := utils.Condition(beego.BConfig.RunMode == "prod", NS_PROD, NS_DEV)

	// Create nacos config stub and setup it
	stub := NewConfigStub(ns, svr)
	if err := stub.Setup(); err != nil {
		panic("New config stub, err:" + err.Error())
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

	mc.dispathParsers(dataId, data)
	cb(dataId, data)

	// listing config changes
	naclog.I("Start listing config dateId:", dataId)
	mc.Stub.Listen(dataId, gp, mc.OnChanged)
}

// Get and listing the configs of indicated dataIds
func (mc *MetaConfig) ListenConfigs(cb MetaConfigCallback, dataIds ...string) {
	for _, dataId := range dataIds {
		if dataId != "" && cb != nil {
			mc.ListenConfig(dataId, cb)
		}
	}
}

// Listing callback called when target configs changed
func (mc *MetaConfig) OnChanged(namespace, group, dataId, data string) {
	if namespace != NS_DEV && namespace != NS_PROD {
		naclog.E("Invalid meta config ns:", namespace)
		return
	}

	if callback, ok := mc.Callbacks[dataId]; ok {
		mc.dispathParsers(dataId, data)
		// Call callback for next handle
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
		naclog.E("Pull nacos routers, err:", err)
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
		naclog.E("Pull nacos routers, err:", err)
		return err
	}

	// update routers chineses descriptions
	routers, err := utils.UpdateChineses(nrouters, descs)
	if err != nil {
		naclog.E("Update routers chineses, err:", err)
		return err
	}

	mc.PushConfig(DID_API_ROUTERS, routers)
	return nil
}

// ----------------------------------------

// Create nacos client config, contain nacos remote server and
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
		naclog.E("Invalid nocos proxy ip:", proxy)
		return "", invar.ErrInvalidParams
	}

	segment := `.((2((5[0-5])|([0-4]\d)))|([0-1]?\d{1,2}))`
	condition := strings.Join(segments[0:3], ".") + segment
	reg, err := regexp.Compile(condition)
	if err != nil {
		naclog.E("Compile regular err:", err)
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
			naclog.I("Direct use ideal ip:", proxy)
			return ip, nil
		} else if ip != "" {
			satisfips = append(satisfips, ip)
		}
	}

	if len(satisfips) > 0 {
		naclog.I("Use dynamic ip:", satisfips[0])
		return satisfips[0], nil
	}

	naclog.E("Not found any local ips, just return proxy:", proxy)
	return proxy, nil
}

// Dispath custom data parsers to parse data when changed
func (mc *MetaConfig) dispathParsers(dataId, data string) {
	switch dataId {
	case DID_ACC_CONFIGS:
		mc.parseConfigs(data)
	case DID_OTA_BUILDS:
		mc.parseOTAInfo(data)
	case DID_DTALK_NTFERS:
		mc.parseSenders(data)
	case DID_WX_AGENTS:
		mc.parseAgents(data)
	case DID_MIO_USERS:
		mc.parseUsers(data)
	case DID_MIO_PATHS:
		mc.parsePaths(data)
	case DID_QK_WORDS:
		mc.parseWords(data)
	case DID_ES_AGENTS:
		mc.setEsAgent(data)
	case DID_MQTT_AGENTS:
		mqtt.SetupMQLogger(data)
	default: // DID_API_ROUTERS, DID_GRPC_CERTS
		naclog.I("Received configs of", dataId)
	}
}

// Parse acc configs when project register DID_ACC_CONFIGS change event
func (mc *MetaConfig) parseConfigs(data string) {
	conf := AccConfs{}
	if err := json.Unmarshal([]byte(data), &conf); err != nil {
		naclog.E("Unmarshal acc configs, err:", err)
		return
	}
	mc.Conf = conf
	naclog.D("Updated Ass Configs!")
}

// Parse OTA infos when project register DID_OTA_BUILDS change event
func (mc *MetaConfig) parseOTAInfo(data string) {
	ota := make(map[string]*OTAInfo)
	if err := json.Unmarshal([]byte(data), &ota); err != nil {
		naclog.E("Unmarshal OTA infos, err:", err)
		return
	}
	mc.OTA = ota
	naclog.D("Updated OTA infos!")
}

// Parse DingTalk senders when project register DID_DTALK_NTFERS change event
func (mc *MetaConfig) parseSenders(data string) {
	senders := make(map[string]*DTalkSender)
	if err := json.Unmarshal([]byte(data), &senders); err != nil {
		naclog.E("Unmarshal DTalk senders, err:", err)
		return
	}
	mc.Senders = senders
	naclog.D("Updated DTalk senders!")
}

// Parse wechat agents when project register DID_WX_AGENTS change event
func (mc *MetaConfig) parseAgents(data string) {
	agents := make(map[string]*wechat.WxIFAgent)
	if err := json.Unmarshal([]byte(data), &agents); err != nil {
		naclog.E("Unmarshal wechat agents, err:", err)
		return
	}
	mc.Agents = agents
	naclog.D("Updated wechat agents!")
}

// Parse minio service users when project register DID_MIO_USERS change event
func (mc *MetaConfig) parseUsers(data string) {
	users := make(map[string]string)
	if err := json.Unmarshal([]byte(data), &users); err != nil {
		naclog.E("Unmarshal minio users, err:", err)
		return
	}
	mc.Users = users
	naclog.I("Updated minio users!")
}

// Parse minio resource paths when project register DID_MIO_PATHS change event
func (mc *MetaConfig) parsePaths(data string) {
	paths := make(map[string][]*ResPath)
	if err := json.Unmarshal([]byte(data), &paths); err != nil {
		naclog.E("Unmarshal minio paths, err:", err)
		return
	}
	mc.Paths = paths
	naclog.I("Updated minio paths!")
}

// Parse excel rule words when project register DID_QK_WORDS change event
func (mc *MetaConfig) parseWords(data string) {
	words := EWords{}
	if err := json.Unmarshal([]byte(data), &words); err != nil {
		naclog.E("Unmarshal excel rule words, err:", err)
		return
	}
	mc.Words = words
	naclog.I("Updated excel rule words!")
}

// Parse elastic server config and create es client
func (mc *MetaConfig) setEsAgent(data string) {
	c := &ESConfig{}
	if err := json.Unmarshal([]byte(data), &c); err != nil {
		naclog.E("Unmarshal elastic configs, err:", err)
		return
	}

	if err := elastic.NewEsClient(c.Address, c.User, c.Pwd, c.CFP); err != nil {
		naclog.E("Create elastic client, err:", err)
		return
	}
	naclog.I("Update elastic client!")
}
