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
	"strings"

	"github.com/wengoldx/xcore/invar"
	"github.com/wengoldx/xcore/logger"
)

// Encode url params for http GET method request.
func EncodeUrl(rawurl string) string {
	enurl, err := url.Parse(rawurl)
	if err != nil {
		logger.E("Encode urlm err:", err)
		return rawurl
	}
	enurl.RawQuery = enurl.Query().Encode()
	return enurl.String()
}

// Return frontend request ip address from controller.Ctx.Request.RemoteAddr
// of beego without port number, such as '192.168.1.100'.
func GetIP(remoteaddr string) string {
	ip, _, _ := net.SplitHostPort(remoteaddr)
	if ip == "::1" {
		ip = "127.0.0.1"
	}
	logger.I("Got ip [", ip, "] from [", remoteaddr, "]")
	return ip
}

// Return all local ip address of current device, set filter to
// true to filter out docker ips that perfixed '172.', and it will
// return invar.ErrNotFound error when not found.
func GetLocalIPs(filter ...bool) ([]string, error) {
	if netfaces, err := net.Interfaces(); err != nil {
		logger.E("Get net interfaces, err:", err)
		return nil, err
	} else {
		ips := []string{}
		isfilter := len(filter) > 0 && filter[0]
		for _, netface := range netfaces {
			addrs, err := netface.Addrs()
			if err != nil {
				continue
			}

			for _, addr := range addrs {
				switch v := addr.(type) {
				case *net.IPNet:
					if v.IP.IsGlobalUnicast() {
						ip := v.IP.String()
						if isfilter && strings.HasPrefix(ip, "172.") {
							continue
						}
						ips = append(ips, ip)
					}
				}
			}
		}

		// check the result list is empty
		if len(ips) == 0 {
			return nil, invar.ErrNotFound
		}
		return ips, nil
	}
}

// Return all hardware mac and ip address of , it filter out
// the docker ips that perfixed '172.', then output mac and ip
// to string as 'xx:xx:xx:xx:xx:xx / 192.168.1.100'.
func GetMacIPs() ([]string, error) {
	if netfaces, err := net.Interfaces(); err != nil {
		logger.E("Get net interfaces, err:", err)
		return nil, err
	} else {
		outs := []string{}
		for _, netface := range netfaces {
			addrs, err := netface.Addrs()
			if err != nil {
				continue
			}

			mac := netface.HardwareAddr.String()
			for _, addr := range addrs {
				switch v := addr.(type) {
				case *net.IPNet:
					if v.IP.IsGlobalUnicast() {
						if ip := v.IP.String(); !strings.HasPrefix(ip, "172.") {
							address := mac + " / " + ip
							outs = append(outs, address)
						}
					}
				}
			}
		}
		return outs, nil
	}
}
