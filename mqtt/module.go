// Copyright (c) 2018-Now Dunyu All Rights Reserved.
//
// Author      : https://www.wengold.net
// Email       : support@wengold.net
//
// Prismy.No | Date       | Modified by. | Description
// -------------------------------------------------------------------
// 00001       2024/03/15   youhei         New version
// -------------------------------------------------------------------

package mqtt

// Optionals mqtt client configs to connect MQTT broker
type Options struct {
	Host     string // remote MQTT broker host
	Port     int    // remote MQTT broker port number
	ClientID string // current client unique id on broker
	User     *User  // user account and password to connect broker
	CAFile   string // CA cert file for TSL
	CerFile  string // certificate/key file for TSL
	KeyFile  string // secure key file for TSL
}

// MQTT configs pasered from nacos configs server
type MqttConfigs struct {
	Broker  Broker           `json:"broker"`
	Users   map[string]*User `json:"users"`
	CAFile  string           `json:"ca"`
	CerFile string           `json:"cert"`
	KeyFile string           `json:"key"`
}

// MQTT broker address
type Broker struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

// MQTT client login user
type User struct {
	Account  string `json:"user"`
	Password string `json:"pwd"`
}
