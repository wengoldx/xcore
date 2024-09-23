# Golang backend server toolkit

The toolkit support to develop golang backend server as fast as quickly based on beego module, it easy to use Elastic, Nacos, MySQL, Redis, MQTT, GRPC services, and support wechat pay interfaces on version 3, the functions as follow.

* `elastic` : Elastic service utils.
* `invar`   : Constant various and HTTP status definitions.
* `logger`  : Custom looger base on beego log, enable remote send by MQTT.
* `mqtt`    : MQTT service and client utils
* `mvc`     : Backend server utils base on beego, export rest4 interfaces.
* `nacos`   : Nacos registry center utils.
* `secure`  : Encript and decript utils of ASE、RSA、MD5、Base64、Hash and so on.
* `utils`   : Common utils of DingTack notify, file access, stack, queue, SMS, email, task and so on.
* `wechat`  : Wechat pay utils base on version 3.
* `wrpc`    : GRPC utils.
* `wsio`    : Socket.IO utils.

### Usage

Construct a HTTP backend service as follow code in golang project.

    package main

    import (
        _ "qks/routers"
        "github.com/wengoldx/wcore/utils"
     )

    func main() {
        utils.HttpServer()
    }

Or get by the follow commands.

> go get git@github.com:wengoldx/xcore.git
