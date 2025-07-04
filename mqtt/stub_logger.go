// Copyright (c) 2018-Now Dunyu All Rights Reserved.
//
// Author      : https://www.wengold.net
// Email       : support@wengold.net
//
// Prismy.No | Date       | Modified by. | Description
// -------------------------------------------------------------------
// 00001       2019/05/22   yangping       New version
// -------------------------------------------------------------------

package mqtt

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	mq "github.com/eclipse/paho.mqtt.golang"
	"github.com/wengoldx/xcore/utils"
)

const (
	adapterMqtt   = "mqtt"                    // Adapter name of logger ouput by mqtt
	logTopicPre   = "wengold/logs/"           // Mqtt logger publish topic prefix
	logTimeLayout = "2006/01/02 15:04:05.000" // Format string for time output
	appPrefixLen  = 12                        // Fixed app prefix length
)

// Custom logger to output logs by mqtt on single connected client chanel.
//
// NOTICE:
//	- It will not print out error event logger failed generated!
//	- The logger just output logs over level of INFO.
//
// UseCase 1 : Using the custom broker options to setup logger
//
//	mqtt.SetupLogger(&mqtt.Options {
//		Host: "192.168.1.10", Port: 1883,
//		User: &mqtt.User {
//			Account: "username", Password: "password"
//		}
//	})
//
//
// UseCase 2 : Get broker options from nacos configs buffer
//
//	// the data is mqtt configs received from nacos server
//	mqtt.SetupLogger(mqtt.GetOptions(data))        // use current app name as user key
//	mqtt.SetupLogger(mqtt.GetOptions(data, "svr")) // input indicated user key, e.g. 'svr'
//
//
// UseCase 3 : Get broker options from exist mqttStub singleton
//
//	stub := mqtt.Singleton()
//	// Do samething to setup mqttStub
//	mqtt.SetupLogger(stub.GetOptions())
//
// All of the usecase must caller call mqtt.SetupLogger() to set options,
// connect with remote broker and then register logger as beego logger.
type mqttLogger struct {
	Options   *Options  // mqtt broker configs
	Stub      mq.Client // mqtt client instanse
	Topic     string    // publish topic
	AppPerfix string    // app name as log message prefix in fixed len string
}

// Register mqtt logger as a beego logs, it will create
// single mqtt client to output logs
func SetupLogger(opts *Options) {
	if opts != nil {
		getMqttLogger := func() logs.Logger {
			return &mqttLogger{
				AppPerfix: getAppPrefix(),
				Topic:     logTopicPre + beego.BConfig.AppName,
				Options:   opts,
			}
		}
		logs.Register(adapterMqtt, getMqttLogger)
		logs.SetLogger(adapterMqtt, "mqtt-logger")
	}
}

// Parse and return mqtt broker configs, it maybe nil returned
func GetOptions(data string, svr ...string) *Options {
	cfgs := &MqttConfigs{}
	if err := json.Unmarshal([]byte(data), &cfgs); err != nil {
		mqxlog.E("Unmarshal mqtt options err:", err)
		return nil
	}

	userkey := utils.Variable(svr, beego.BConfig.AppName)
	if user, ok := cfgs.Users[userkey]; ok {
		mqxlog.I("Got mqtt options of user:", userkey)
		return &Options{
			Host: cfgs.Broker.Host,
			Port: cfgs.Broker.Port,
			User: user,
		}
	}
	return nil
}

// Return fixed length app name string as log message prefix at format of:
// App-Name    | 2024/05/23 11:15:03:246 [E] [ctrl_exercise.go:189] QRCode() xxxxxxxx
func getAppPrefix() string {
	prefix := beego.BConfig.AppName
	if pl := len(prefix); pl < appPrefixLen {
		for i := appPrefixLen - pl; i > 0; i-- {
			prefix += " "
		}
	} else if pl >= appPrefixLen {
		prefix = prefix[:appPrefixLen]
	}
	return prefix + " | "
}

// Init mqtt logger topic and connect client with id of 'logger.appname'
func (w *mqttLogger) Init(config string) error {
	options := mq.NewClientOptions()
	broker := fmt.Sprintf(protFormatTCP, w.Options.Host, w.Options.Port)
	options.AddBroker(broker)
	options.SetClientID("logger." + beego.BConfig.AppName)
	options.SetUsername(w.Options.User.Account)
	options.SetPassword(w.Options.User.Password)
	options.SetAutoReconnect(true)

	w.Stub = mq.NewClient(options)
	if token := w.Stub.Connect(); token.Wait() && token.Error() != nil {
		w.Stub = nil

		// Delete mqtt logger from beego logs when connect failed
		logs.GetBeeLogger().DelLogger(adapterMqtt)
		mqxlog.E("Setup mqtt logger err:", token.Error())
		return nil // not return error
	}

	mqxlog.I("Connected mqtt logger!")
	return nil
}

// Publish logs above info level after mqtt client connected
func (w *mqttLogger) WriteMsg(when time.Time, msg string, level int) error {
	if w.Stub != nil && level <= logs.LevelInfo && msg != "" {
		msg = w.AppPerfix + when.Format(logTimeLayout) + " " + msg
		w.Stub.Publish(w.Topic, 0, false, msg)
	}
	return nil
}

// Disconnect mqtt client if living
func (w *mqttLogger) Destroy() {
	if w.Stub != nil && w.Stub.IsConnected() {
		w.Stub.Disconnect(0)
		w.Stub = nil
	}
}

// Do nothing here, none cache to output
func (w *mqttLogger) Flush() {}
