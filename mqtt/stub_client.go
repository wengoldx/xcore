// Copyright (c) 2018-Now Dunyu All Rights Reserved.
//
// Author      : https://www.wengold.net
// Email       : support@wengold.net
//
// Prismy.No | Date       | Modified by. | Description
// -------------------------------------------------------------------
// 00001       2022/11/08   jidi           New version
// -------------------------------------------------------------------

package mqtt

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"reflect"

	"github.com/astaxie/beego"
	mq "github.com/eclipse/paho.mqtt.golang"
	"github.com/wengoldx/xcore/invar"
	"github.com/wengoldx/xcore/logger"
	"github.com/wengoldx/xcore/secure"
	"github.com/wengoldx/xcore/utils"
)

// MQTT stub to manager MQTT connection.
//
// As usages you can connect remote MQTT broker and get client instance by follow usecases.
//
//
// UseCase 1 : Using nacos MQTT configs and connect without callbacks.
//
//	if err := mqtt.NewClient(data); err != nil {
//		mqxlog.E("Connect client err:", err)
//		return
//	}
//
//
// UseCase 2 : Using nacos MQTT configs and connect with callbacks, it same as call mqtt.SetMQClient(data, 2, ConnectHandler).
//
//	stub := mqtt.SetOptions(2)
//	stub.ConnectHandler = ConnectHandler
//	if err := mqtt.NewClient(data); err != nil {
//		mqxlog.E("Connect client err:", err)
//		return
//	}
//
//
// UseCase 3 : Using singleton stub set custom configs and connect.
//
//	stub := mqtt.Singleton()
//	// Here use stub.Options to set broker configs
//	//      use stub.xxxHandler to set connect, disconnect, message handlers
//	if err := stub.Connect(stub.GetConnOptions()); err != nil {
//		mqxlog.E("Connect client err:", err)
//		return
//	}
type MqttStub struct {
	Options           *Options                 // Broker host and port, login secure datas, client id
	Client            mq.Client                // MQTT client instance
	ConnectHandler    mq.OnConnectHandler      // Connect callback handler
	DisconnectHandler mq.ConnectionLostHandler // Disconnect callback handler
	ReconnectHandler  mq.ReconnectHandler      // Reconnect callback handler
	MessageHandler    mq.MessageHandler        // Default publish message callback handler
	qos               byte                     // The default qos for publish or subscribe
	remain            bool                     // The default remain flag
}

// Singleton mqtt stub instance
var mqttStub *MqttStub

// Object logger with [MQX] mark for MQTT module
var mqxlog = logger.CatLogger("MQX")

// Default connect handler, change it before call GetConnOptions().
var connHandler mq.OnConnectHandler = func(client mq.Client) {
	serve, opt := beego.BConfig.AppName, client.OptionsReader()
	mqxlog.I("Server", serve, "connected mqtt as client:", opt.ClientID())
}

// Default disconnect handler, change it before call GetConnOptions().
var lostHandler mq.ConnectionLostHandler = func(client mq.Client, err error) {
	serve, opt := beego.BConfig.AppName, client.OptionsReader()
	mqxlog.W("Server", serve, "disconnect mqtt client:", opt.ClientID())
}

const (
	protFormatTCP = "tcp://%s:%v" // Mqtt protocol of TCP
	protFormatSSL = "ssl://%s:%v" // Mqtt protocel of SSL
)

// Return Mqtt global singleton
func Singleton() *MqttStub {
	if mqttStub == nil {
		mqttStub = &MqttStub{
			Options: &Options{}, Client: nil,
			ConnectHandler: connHandler, DisconnectHandler: lostHandler,
			qos: byte(0), remain: false,
		}
	}
	return mqttStub
}

// Create a new mqtt client and connect with MQTT broker, the client using
// 'tcp' protocol and fixed id as format 'server@12345678'.
//
//	* The configs input param maybe set as json string from Nacos Configs Server
//	* Or input Options object refrence created at local
func NewClient(configs any, server ...string) error {
	svr := utils.VarString(server, beego.BConfig.AppName)
	stub := Singleton()

	// parse MQTT connect configs from json string or Options object refrence
	switch reflect.ValueOf(configs).Interface().(type) {
	case string:
		if err := stub.parseConfig(configs.(string), svr); err != nil {
			return err
		}
	case *Options:
		stub.Options = configs.(*Options)
		if stub.Options.ClientID == "" {
			stub.Options.ClientID = svr + "." + secure.NewCode()
		}
	default:
		return invar.ErrInvalidConfigs
	}

	opt := stub.GetConnOptions() // using default tcp protocol
	if err := stub.Connect(opt); err != nil {
		mqxlog.E("New", svr, "mqtt client err:", err)
		return err
	}
	return nil
}

// Set default qos and remain flag
func SetOptions(qos byte, remain ...bool) *MqttStub {
	stub := Singleton()
	if len(remain) > 0 {
		stub.remain = remain[0]
	}
	stub.qos = qos
	return stub
}

// Create mqtt config, default connection protocol using tcp, you can
// set mode 'tls' and cert files to using ssl protocol.
func (stub *MqttStub) GetConnOptions(mode ...string) *mq.ClientOptions {
	options, protocol := mq.NewClientOptions(), protFormatTCP
	if len(mode) > 0 && mode[0] == "tls" {
		protocol = protFormatSSL
		if tlscfg := stub.newTLSConfig(); tlscfg != nil {
			options.SetTLSConfig(tlscfg)
		}
	}

	broker := fmt.Sprintf(protocol, stub.Options.Host, stub.Options.Port)
	options.AddBroker(broker)
	options.SetClientID(stub.Options.ClientID)
	options.SetUsername(stub.Options.User.Account)
	options.SetPassword(stub.Options.User.Password)
	options.SetAutoReconnect(true)

	// set callback handlers if exist
	options.SetOnConnectHandler(stub.ConnectHandler)
	options.SetConnectionLostHandler(stub.DisconnectHandler)
	options.SetReconnectingHandler(stub.ReconnectHandler)
	options.SetDefaultPublishHandler(stub.MessageHandler)
	return options
}

// New client from given options and connect with broker
func (stub *MqttStub) Connect(opt *mq.ClientOptions) error {
	stub.Client = mq.NewClient(opt)
	if token := stub.Client.Connect(); token.Wait() && token.Error() != nil {
		stub.Client = nil
		mqxlog.E("Connect mqtt client, err:", token.Error())
		return token.Error()
	}
	return nil
}

// Publish empty message topic, it same use for just notify
func (stub *MqttStub) Notify(topic string, Qos ...byte) error {
	return stub.PublishOptions(topic, nil, stub.remain, Qos...)
}

// Publish indicate topic message, the Qos can be set current call in 0 ~ 2
func (stub *MqttStub) Publish(topic string, data any, Qos ...byte) error {
	return stub.PublishOptions(topic, data, stub.remain, Qos...)
}

// Publish indicate topic message with input remain flag and Qos options,
//
// Notice that the data will encode as json bytes array if value type is Struct,
// Pointer or map, or instead nil data to empty bytes array.
func (stub *MqttStub) PublishOptions(topic string, data any, remain bool, Qos ...byte) error {
	if stub.Client == nil {
		mqxlog.E("Abort publish topic:", topic, "on nil client!!")
		return invar.ErrInvalidClient
	}

	var payload any
	if data != nil {
		switch reflect.ValueOf(data).Kind() {
		case reflect.Struct, reflect.Pointer, reflect.Map:
			if buffer, err := json.Marshal(data); err != nil {
				return err
			} else {
				payload = buffer
			}
		default:
			payload = data
		}
	} else {
		payload = []byte{} // Instead nil to empty bytes
	}

	qosv := stub.qos
	if len(Qos) > 0 && Qos[0] > 0 && Qos[0] <= 2 {
		qosv = Qos[0]
	}

	token := stub.Client.Publish(topic, qosv, remain, payload)
	if token.Wait() && token.Error() != nil {
		mqxlog.E("Publish topic:", topic, "err:", token.Error())
		return token.Error()
	}

	mqxlog.I("Published topic:", topic)
	return nil
}

// Subscribe given topic and set callback
func (stub *MqttStub) Subscribe(topic string, hanlder mq.MessageHandler, Qos ...byte) error {
	if stub.Client == nil {
		mqxlog.E("Abort subscribe topic:", topic, "on nil client!!")
		return invar.ErrInvalidClient
	}

	qosv := stub.qos
	if len(Qos) > 0 && Qos[0] > 0 && Qos[0] <= 2 {
		qosv = Qos[0]
	}

	token := stub.Client.Subscribe(topic, qosv, hanlder)
	if token.Wait() && token.Error() != nil {
		mqxlog.E("Subscribe topic:", topic, "err:", token.Error())
		return token.Error()
	}
	mqxlog.I("Subscribed topic:", topic)
	return nil
}

// Return mqtt broker host, port and login user after mqttStub established
func (stub *MqttStub) GetOptions() *Options {
	return &Options{
		Host: stub.Options.Host,
		Port: stub.Options.Port,
		User: stub.Options.User,
	}
}

// Load and create secure configs for TLS protocol to connect.
func (stub *MqttStub) newTLSConfig() *tls.Config {
	opts := stub.Options
	ca, err := os.ReadFile(opts.CAFile)
	if err != nil {
		mqxlog.E("Read CA file err:", err)
		return nil
	}

	certpool := x509.NewCertPool()
	certpool.AppendCertsFromPEM(ca)
	tlsConfig := &tls.Config{RootCAs: certpool}

	// Import client certificate/key pair
	if opts.CerFile != "" && opts.KeyFile != "" {
		keyPair, err := tls.LoadX509KeyPair(opts.CerFile, opts.KeyFile)
		if err != nil {
			mqxlog.E("Load cert and key err:", err)
			return nil
		}

		tlsConfig.ClientAuth = tls.NoClientCert
		tlsConfig.ClientCAs = nil
		tlsConfig.InsecureSkipVerify = true
		tlsConfig.Certificates = []tls.Certificate{keyPair}
	}
	return tlsConfig
}

// Parse mqtt broker and all user datas from nacos config center
func (stub *MqttStub) parseConfig(data, svr string) error {
	cfgs := &MqttConfigs{}
	if err := json.Unmarshal([]byte(data), &cfgs); err != nil {
		mqxlog.E("Unmarshal mqtt settings, err:", err)
		return err
	}

	// Create client configs and fix the id as 'server@123456789'
	if user, ok := cfgs.Users[svr]; !ok {
		return errors.New("Not found mqtt user: " + svr)
	} else {
		stub.Options.Host = cfgs.Broker.Host
		stub.Options.Port = cfgs.Broker.Port
		stub.Options.User = user
		stub.Options.CAFile = cfgs.CAFile
		stub.Options.CerFile = cfgs.CerFile
		stub.Options.KeyFile = cfgs.KeyFile

		// Random client id if not fixed
		if stub.Options.ClientID == "" {
			stub.Options.ClientID = svr + "." + secure.NewCode()
		}
	}
	return nil
}

// ----------------------------------------
// MQTT Local Client Setup
// ----------------------------------------

// A simple way to setup mqtt client by given connect handler and optional params,
// you can subscribe mqtt topics on return true. The optional params must set as
// sequence [qos:int, remain:bool, svr:string]
func SetupMQClient(data string, handler mq.OnConnectHandler, opts ...any) bool {
	optslen := len(opts)
	qos, remain, svr := 0, false, ""

	// get optional params if exist
	if optslen > 0 {
		qos = opts[0].(int)
	}
	if optslen > 1 {
		remain = opts[1].(bool)
	}
	if optslen > 2 {
		svr = opts[2].(string)
	}

	stub := SetOptions(byte(qos), remain)
	if handler != nil {
		stub.ConnectHandler = handler
	}

	if err := NewClient(data, svr); err != nil {
		mqxlog.E("Create mqtt client, err:", err)
		return false
	}
	return true
}

// A simple way to setup mqtt remote logger from original nacos config datas, it called
// by internal nacos module when local server listen nacos.DID_MQTT_AGENTS config key.
func SetupMQLogger(data string) {
	SetupLogger(GetOptions(data, "public"))
}
