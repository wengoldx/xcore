// Copyright (c) 2018-Now Dunyu All Rights Reserved.
//
// Author      : https://www.wengold.net
// Email       : support@wengold.net
//
// Prismy.No | Date       | Modified by. | Description
// -------------------------------------------------------------------
// 00001       2023/04/18   tangxiaoyu     New version
// -------------------------------------------------------------------

package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/astaxie/beego"
	"github.com/wengoldx/xcore/logger"
)

// Swagger.json field keywords for version 2.0.0.
//
// # WARNING:
//
//	DO NOT CHANGE THEME IF YOU KNOWE HOW TO CHANGE IT!
const (
	_swagger_json_file = "./swagger/swagger.json"
	_key_paths         = "paths"
	_key_method_get    = "get"
	_key_method_post   = "post"
	_key_summary       = "summary"
	_key_group_tags    = "tags"
	_key_group_name    = "name"
	_key_group_desc    = "description"
)

// A router informations
type Router struct {
	Router string `json:"router"` // restful router full path start /
	Method string `json:"method"` // restful router http method, such as GET, POST...
	Group  string `json:"group"`  // beego controller keyworld
	EnDesc string `json:"endesc"` // beego controller router summary
}

// A group informations
type Group struct {
	Name   string `json:"group"`  // group name as swagger controller path, like '/v3/acc'
	EnDesc string `json:"endesc"` // beego controller router summary
}

// All routers of one server
type Routers struct {
	Groups  []*Group  `json:"groups"`  // groups as swagger controllers
	Routers []*Router `json:"routers"` // routers parsed from swagger.json file
}

// Parse total routers and update description on chinese for local server routers,
// then marshal to string for next to push to nacos config server.
func UpdateRouters(data string) (string, error) {
	routers, err := loadSwaggerRouters()
	if err != nil {
		logger.E("Load local swagger, err:", err)
		return "", err
	}

	svr := beego.BConfig.AppName
	rsmap := parseNacosRouters(data)
	if rsmap != nil {
		rsmap[svr] = routers
	} else {
		rsmap = make(map[string]*Routers)
		rsmap[svr] = routers
	}

	swagger, err := json.Marshal(rsmap)
	if err != nil {
		logger.E("Marshal routers, err:", err)
		return "", err
	}

	logger.D("Updated routers apis for", svr)
	return string(swagger), nil
}

const _router_ver_file = "./routers/router.go"
const _router_ver_key = "// @APIVersion "

// Update swagger version by search keyword '// @APIVersion '
// in project ./routers/router.go code file.
//
// ---
//
//	                                         end
//	                                         v
//	router.go content : '// @APIVersion 1.2.3'
//	                     ^              ^
//	                     start          verstart
//
// ---
//
// This method is useful to update swagger API version when
// excute 'bee run -gendoc=true' command in project root dir.
//
// # USAGE:
//
// The best timing to call it before database inited in init().
//
//	package content
//	import "github.com/wengoldx/xcore/utils"
//	func init() {
//	    utils.UpdateVersion(types.APP_VERSION)
//	    // ...
//	}
func UpdateVersion(ver string) {
	// clear all black chars out from version string!
	ver = strings.TrimSpace(strings.ReplaceAll(ver, " ", ""))

	// check version not empty and router.go code file exist!
	if ver != "" && IsFile(_router_ver_file) {
		if content := ReadTextFile(_router_ver_file); content != "" {

			// router.go code file exist and not empty!
			if start := strings.Index(content, _router_ver_key); start != -1 {
				// '// @APIVersion 1.2.3' seed to '1' position.
				verstart := start + len(_router_ver_key)

				lineend := strings.Index(content[verstart:], "\n")
				if lineend > 0 && lineend < len("xx.xx.xx") {
					// '// @APIVersion 1.2.3' seed to line end.
					end := verstart + lineend

					// not upgrade if version up to date!
					if old := content[verstart:end]; old == ver {
						logger.I("Swagger version is last!")
						return
					}

					// upgrade swagger version and overwrite to router.go file.
					newver := fmt.Sprintf(_router_ver_key+"%s", ver)
					content = content[:start] + newver + content[end:]
					if err := os.WriteFile(_router_ver_file, []byte(content), 0755); err != nil {
						logger.E("Update swagger version, err:", err)
						return
					}
					logger.I("Upgrade swagger version:", ver)
				}
			}
		}
	}
}

/* ------------------------------------------------------------------- */

// Load local server routers from swagger.json file.
//
//	{
//	    "swagger": "2.0",
//	    "info": { ... },
//	    "basePath": "/myserver",
//	    "paths": {
//	        "/debug/api": {
//	            "get": {
//	                "tags": [ "debug" ],
//	                "summary": "Api Test",
//	                "description": "Test Api request status.",
//	                "responses": {
//	                    "200": { "description": "" },
//	                    "404": { "description": "Server internal error." }
//	                }
//	            }
//	        }, ...
//	    },
//	    "definitions": { ... },
//	    "tags": [
//	        {
//	            "name": "debug",
//	            "description": "Debug controller to easy test on develop mode."
//	        }, ...
//	    ]
//	}
func loadSwaggerRouters() (*Routers, error) {
	buff, err := os.ReadFile(_swagger_json_file)
	if err != nil {
		logger.E("Load swagger routers, err:", err)
		return nil, err
	}

	// parse 'paths' and 'tags' fields.
	routers := make(map[string]any)
	if err := json.Unmarshal(buff, &routers); err != nil {
		logger.E("Unmarshal swagger routers err:", err)
		return nil, err
	}
	logger.I("Loaded swagger json, start parse routers")
	out := &Routers{}

	// fetch routers from 'paths' field values.
	if ps, ok := routers[_key_paths]; ok && ps != nil {
		paths := ps.(map[string]any) // '/debug/api': any.

		// fetch target router infos such as '/debug/api'.
		for path, pathvals := range paths {
			router := &Router{Router: path} // parse router path

			// parse http GET or POST methods.
			var mvs any
			pvs := pathvals.(map[string]any) // ['get'|'post']: any
			if pmg, ok := pvs[_key_method_get]; ok && pmg != nil {
				router.Method, mvs = "GET", pmg
			} else if pmp, ok := pvs[_key_method_post]; ok && pmp != nil {
				router.Method, mvs = "POST", pmp
			} else {
				logger.W("Invalid method of path:", path)
				continue
			}

			// parse beego controller group name.
			method := mvs.(map[string]any)
			if gps, ok := method[_key_group_tags]; ok && gps != nil {
				groups := reflect.ValueOf(gps)
				router.Group = groups.Index(0).Interface().(string)
			}

			// parse router path english description.
			if desc, ok := method[_key_summary]; ok && desc != nil {
				router.EnDesc = desc.(string)
			}

			// append the router into routers array.
			// logger.D("> Parsed ["+router.Method+"]\tpath:", path, "\tdesc:", router.EnDesc)
			out.Routers = append(out.Routers, router)
		}
	}

	// fetch routers from 'tags' field values.
	if gps, ok := routers[_key_group_tags]; ok && gps != nil {
		groups := gps.([]any) // parse all group array.

		for _, group := range groups {
			t := group.(map[string]any)
			gp := &Group{}

			// parse group name value.
			if gpn, ok := t[_key_group_name]; ok && gpn != nil {
				gp.Name = gpn.(string)
			}

			// parse group english description.
			if gpd, ok := t[_key_group_desc]; ok && gpd != nil {
				gp.EnDesc = strings.TrimRight(gpd.(string), "\n")
				if end := strings.Index(gp.EnDesc, "\n"); end != -1 {
					gp.EnDesc = gp.EnDesc[:end]
				}
			}

			// append the group into groups array
			// logger.D("# Parsed group ["+gp.Name+"] \t desc:", gp.EnDesc)
			out.Groups = append(out.Groups, gp)
		}
	}

	logger.I("Finish parsed swagger routers")
	return out, nil
}

// Parse servers routers from nacos config data.
func parseNacosRouters(data string) map[string]*Routers {
	routers := make(map[string]*Routers)
	if data != "" && data != "{}" { // check data if empty
		if err := json.Unmarshal([]byte(data), &routers); err != nil {
			logger.E("Unmarshal nacos routers, err:", err)
			return nil
		}
	}
	return routers
}
