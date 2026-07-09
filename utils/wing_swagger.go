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
	"github.com/wengoldx/xcore/invar"
	"github.com/wengoldx/xcore/logger"
)

// Swagger.json field keywords for version 2.0.0.
//
// # WARNING:
//
//	DO NOT CHANGE THEME IF YOU KNOWE HOW TO CHANGE IT!
const (
	_swagger_json_file = "./swagger/swagger.json"
	_key_paths         = "paths"       // ROUTER KEY.
	_key_method_get    = "get"         // ROUTER KEY.
	_key_method_post   = "post"        // ROUTER KEY.
	_key_summary       = "summary"     // ROUTER KEY.
	_key_tags          = "tags"        // ROUTER & GROUP KEY.
	_key_name          = "name"        // GROUP KEY.
	_key_desc          = "description" // GROUP KEY.
)

// Controller api router path.
type Router struct {
	Router string `json:"router"` // api router full path, like '/debug/token'.
	Method string `json:"method"` // api router http method, 'GET' or 'POST'.
	Group  string `json:"group"`  // api router group's tag.
}

// All routers of one server
type Routers struct {
	Tags  []string  `json:"tags"`  // server controllers, from 'tags' array of swagger.json.
	Paths []*Router `json:"paths"` // controller api routers, from 'paths.{router}'
}

// Role router policys of server to auto set api access permissions.
type RPolicy struct {
	App     string   // server app name.
	Role    string   // role key, such as 'admin', 'user', 'comp', 'mach', 'part', 'rb'.
	Policy  string   // api path as role policy, like 'server/v4/utils/admin/*'.
	Methods []string // api methods, one or all 'GET' and 'POST'.
}

// Parse total routers and update description on chinese for local server routers,
// then marshal to string for next to push to nacos config server.
func UpdateRouters(data string) (string, error) {
	routers, err := loadSwaggerRouters(_swagger_json_file)
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

// Parse nocos swagger apis datas and return all roles routers
// for next auto append restful apis access permissions.
//
// # WARNIG:
//   - Call this method after nacos config server connected.
func ParseRouters(data string) []*RPolicy {
	rsmap := parseNacosRouters(data)
	if rsmap != nil {
		policys := []*RPolicy{}
		for app, routers := range rsmap {
			rps := parseRoleRouters(app, routers.Paths)
			for _, rp := range rps {
				policys = append(policys, rp)
			}
		}
		return policys
	}
	return nil
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
//	    "tags": [
//	        {
//	            "name": "debug",
//	            "description": "Debug controller to easy test on develop mode."
//	        }, ...
//	    ]
//	}
func loadSwaggerRouters(sf string) (*Routers, error) {
	buff, err := os.ReadFile(sf) // swagger.json file path.
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
			if tags, ok := method[_key_tags]; ok && tags != nil {
				groups := reflect.ValueOf(tags)
				router.Group = groups.Index(0).Interface().(string)
			}

			// append the router into routers array.
			// logger.D("> Parsed ["+router.Method+"]\tpath:", path, "\tmethod:", router.Group)
			out.Paths = append(out.Paths, router)
		}
	}

	// fetch routers from 'tags' field values.
	if tags, ok := routers[_key_tags]; ok && tags != nil {
		groups := tags.([]any) // parse all group array.

		for _, group := range groups {
			tag := group.(map[string]any)

			// parse group name value.
			if tn, ok := tag[_key_name]; ok && tn != nil {
				name := tn.(string)

				// append the tag into tags array.
				// logger.D("# Parsed group:", name)
				out.Tags = append(out.Tags, name)
			}
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

// Parse role routers from nacos router datas.
func parseRoleRouters(app string, routers []*Router) map[string]*RPolicy {
	outs := make(map[string]*RPolicy)
	for _, r := range routers {
		// '/v4/mach/code/auth'    : 'v4/mach'  -> path = 'code/auth'
		// '/v4/utils/admin/confs' : 'v4/utils' -> path = 'admin/confs'
		path := strings.TrimPrefix(r.Router, "/"+r.Group+"/")
		if path != "" {
			if seg := strings.Split(path, "/"); len(seg) > 0 {
				if key := seg[0]; invar.IsRoleKey(key) {
					// api = 'server/v4/utils/admin/*'
					api := fmt.Sprintf("/%s/%s/%s/*", app, r.Group, key)
					if role, ok := outs[api]; ok {
						if !Contain(role.Methods, r.Method) {
							role.Methods = append(role.Methods, r.Method)
						}
						continue
					}
					outs[api] = &RPolicy{App: app, Role: key,
						Policy: api, Methods: []string{r.Method}}
				}
			}
		}
	}
	return outs
}
