// Copyright (c) 2018-Now Dunyu All Rights Reserved.
//
// Author      : https://www.wengold.net
// Email       : support@wengold.net
//
// Prismy.No | Date       | Modified by. | Description
// -------------------------------------------------------------------
// 00001       2019/05/22   yangping       New version
// -------------------------------------------------------------------

package httpx

import (
	"net"
	"net/url"

	"github.com/wengoldx/xcore/invar"
	"github.com/wengoldx/xcore/logger"
)

// GetIP get just ip not port from controller.Ctx.Request.RemoteAddr of beego
func GetIP(remoteaddr string) string {
	ip, _, _ := net.SplitHostPort(remoteaddr)
	if ip == "::1" {
		ip = "127.0.0.1"
	}
	logger.I("Got ip [", ip, "] from [", remoteaddr, "]")
	return ip
}

// GetLocalIPs get all the loacl IP of current deploy machine
func GetLocalIPs() ([]string, error) {
	netfaces, err := net.Interfaces()
	if err != nil {
		logger.E("Get ip interfaces err:", err)
		return nil, err
	}

	ips := []string{}
	for _, netface := range netfaces {
		addrs, err := netface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			switch v := addr.(type) {
			case *net.IPNet:
				if v.IP.IsGlobalUnicast() {
					ips = append(ips, v.IP.String())
				}
			}
		}
	}

	// Check the result list is empty
	if len(ips) == 0 {
		return nil, invar.ErrNotFound
	}

	return ips, nil
}

// EncodeUrl encode url params
func EncodeUrl(rawurl string) string {
	enurl, err := url.Parse(rawurl)
	if err != nil {
		logger.E("Encode urlm err:", err)
		return rawurl
	}
	enurl.RawQuery = enurl.Query().Encode()
	return enurl.String()
}
