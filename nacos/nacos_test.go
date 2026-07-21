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
	"testing"

	xt "github.com/wengoldx/xcore/utils/xtest"
)

const _test_config = `{
    "email" : {
        "smtps": {
            "qq": {
                "host":"smtp.exmail.qq.com",
                "port":465
            },
            "ali": {
                "host":"smtp.qiye.aliyun.com",
                "port":465
            }
        },
        "serves": {
            "myserver-1": {
                "acc": {
                    "user":"user1@email.com",
                    "pwd":"123456"
                },
                "web": {
                    "user":"user2@email.com",
                    "pwd":"654321"
                }
            },
            "myserver-2": {
                "chat": {
                    "user":"user3@email.com",
                    "pwd":"123456"
                }
            }
        }
    }
}`

func TestParseConfig(t *testing.T) {
	ac := &AccConfs{}
	if err := json.Unmarshal([]byte(_test_config), ac); err != nil {
		t.Fatal("Error:", err)
	}
	xt.LogI("Parsed out:\n", ac)

	for serve, senders := range ac.Email.Serves {
		for tag, sender := range senders {
			xt.LogI("> server:", serve, "sender:", tag,
				"acc:", sender.User, "pwd:", sender.Pwd)
		}
	}

	senders := ac.Email.Serves["myserver-1"]
	xt.LogI("Senders:", senders)
}

func TestMatchProxyIP(t *testing.T) {
	xt.UseDebugLogger()

	proxy := "192.168.1.100"
	if result, err := matchProxyIP(proxy); err != nil {
		t.Fatal(err)
	} else if proxy != result {
		xt.LogI("Using local ip:", result)
	} else {
		xt.LogI("None local ip!")
	}
}
