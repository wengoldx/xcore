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

/* -------------------------- */
/* Internal constants defines */
/* -------------------------- */

const (
	nacosSysSecure = "accessor"      // CICD system secure authentications
	nacosLogLevel  = "warn"          // default level to print nacos logs on warn, it not same nacos-sdk-go:info
	nacosDirLogs   = "./nacos/logs"  // default nacos logs dir
	nacosDirCache  = "./nacos/cache" // default nacos caches dir

	configKeySvr  = "nacossvr"  // Nacos remote server IP address
	configKeyAddr = "nacosaddr" // Local server access IP address
	configKeyPort = "nacosport" // Local server access port for grpc connect
)

/* -------------------------- */
/* Export constants defines   */
/* -------------------------- */

// Nacos namespace string for xcore/nacos
const (
	NS_PROD = "dunyu-server-prod" // PROD namespace id
	NS_DEV  = "dunyu-server-dev"  // DEV  namespace id
)

// Fixed all registered servers and configs named 'wengold' group
const GP_WENGOLD = "group.wengold"

// Nacos data id for xcore/nacos
const (
	DID_ACC_CONFIGS  = "dunyu.acc.configs"  // Fixed group, data id of accservice cofnigs
	DID_API_ROUTERS  = "dunyu.api.routers"  // Fixed group, data id of swagger restful routers
	DID_DTALK_NTFERS = "dunyu.dtalk.ntfers" // Fixed group, data id of dingtalk notifiers
	DID_ES_AGENTS    = "dunyu.es.agents"    // Fixed group, data id of elastic search agents
	DID_GRPC_CERTS   = "dunyu.grpc.certs"   // Fixed group, data id of grpc certs that datas format as xml
	DID_MIO_PATHS    = "dunyu.mio.paths"    // Fixed group, data id of minio source paths
	DID_MIO_USERS    = "dunyu.mio.users"    // Fixed group, data id of minio account key
	DID_MQTT_AGENTS  = "dunyu.mqtt.agents"  // Fixed group, data id of mqtt agents
	DID_OTA_BUILDS   = "dunyu.ota.builds"   // Fixed group, data id of all projects OTA infos, get data from mc.OTA maps
	DID_WX_AGENTS    = "dunyu.wx.agents"    // Fixed group, data id of wechat agents
)

/* -------------------------- */
/* Export Configs defines     */
/* -------------------------- */

// Nacos config for data id DID_ACC_CONFIGS
type AccConfs struct {

	// Email sender service
	Email struct {
		Host     string `json:"host"`
		Port     int    `json:"port"`
		Identity string `json:"identity"`
	} `json:"email"`

	// SMS sender service
	Sms struct {
		Secret    string `json:"secret"`
		KeyID     string `json:"keyid"`
		URLFormat string `json:"urlformat"`
	} `json:"sms"`

	// Account secure settings
	Secures struct {
		SecureSalt   string `json:"secureSalt"`   // Secure salt key to decode account login token
		ApiTaxCode   string `json:"apiTaxCode"`   // Auth code to access API of check company tax code
		ApiIDViaCode string `json:"apiIDViaCode"` // Auth code to access API of identification check
		PageLimits   int    `json:"pageLimits"`   // One times to get list item counts on a page
	} `json:"secure"`

	// Mall account settings
	MallAccs map[string]*MallAcc `json:"mallaccs"`
}

// Nacos config for mall account settings
type MallAcc struct {
	User string `json:"user"`
	Pwd  string `json:"pwd"`
}

// Nacos config for OTA upgrade by using DID_OTA_BUILDS data id
type OTAInfo struct {
	BuildVersion string  `json:"BuildVersion" description:"Build version string"`
	BuildNumber  int     `json:"BuildNumber"  description:"Build number, pase form BuildVersion string as version = major*10000 + middle*100 + minor"`
	DownloadUrl  string  `json:"DownloadUrl"  description:"Bin file download url"`
	UpdateDate   string  `json:"UpdateDate"   description:"Bin file update date"`
	HashSums     string  `json:"HashSums"     description:"Bin file hash sums"`
	BinSizes     float64 `json:"BinSizes"     description:"Bin file sizes in MB"`
}

// Nacos config for DingTalk notify sender
type DTalkSender struct {
	WebHook   string   `json:"webhook"`   // DingTalk group chat session webhook
	Secure    string   `json:"secure"`    // DingTalk group chat senssion secure key
	Receivers []string `json:"receivers"` // The target @ users
}

// Nacos config for GRPC cert content
type GrpcCert struct {
	Svr string `xml:"Server"` // GRPC cert server name
	Key string `xml:"Key"`    // GRPC cert key data
	Pem string `xml:"Pem"`    // GRPC cert pem data
}

// Nacos config for GRPC certs
type GrpcCerts struct {
	Certs []GrpcCert `xml:"Cert"` // GRPC certs
}

// Bucket path bund resource number to export MinIO bucket paths
type ResPath struct {
	Res  string `json:"res"`  // Resource number as unique id used by outside to bind real bucket path
	Path string `json:"path"` // Real bucket path of MinIO service
}
