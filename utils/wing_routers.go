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
	"os"
	"reflect"
	"strings"

	"github.com/astaxie/beego"
	"github.com/wengoldx/xcore/invar"
	"github.com/wengoldx/xcore/logger"
)

// Swagger.json field keywords for version 2.0.0.
//
// `NOTICE`
//
//	DO NOT CHANGE THEME IF YOU KNOWE HOW TO CHANGE IT!
const (
	swaggerFile  = "./swagger/swagger.json"
	sfPathName   = "paths"
	sfMethodGet  = "get"
	sfMethodPost = "post"
	sfEnDescName = "description"
	sfGroupTags  = "tags"
	sfGroupName  = "name"
)

// A router informations
type Router struct {
	Router string `json:"router"` // restful router full path start /
	Method string `json:"method"` // restful router http method, such as GET, POST...
	Group  string `json:"group"`  // beego controller keyworld
	EnDesc string `json:"endesc"` // english description of router from swagger
	CnDesc string `json:"cndesc"` // chinese description of router manual update by user
}

// A group informations
type Group struct {
	Name   string `json:"group"`  // group name as swagger controller path, like '/v3/acc'
	EnDesc string `json:"endesc"` // english description of group from swagger
	CnDesc string `json:"cndesc"` // chinese description of group manual update by user
}

// All routers of one server
type Routers struct {
	CnName  string    `json:"cnname"`  // backend server chinese name
	Groups  []*Group  `json:"groups"`  // groups as swagger controllers
	Routers []*Router `json:"routers"` // routers parsed from swagger.json file
}

type SvrDesc struct {
	Server  string            `json:"server"`  // backend server english name
	CnName  string            `json:"cnname"`  // backend server chinese name
	Groups  map[string]string `json:"groups"`  // groups  chinese description
	Routers map[string]string `json:"routers"` // routers chinese description
}

// Parse total routers and update description on chinese for local server routers,
// then marshal to string for next to push to nacos config server.
func UpdateRouters(data string) (string, error) {
	routers, err := loadSwaggerRouters()
	if err != nil {
		logger.E("Load local swagger, err:", err)
		return "", err
	}

	rsmap, nrs := parseNacosRouters(data)
	if nrs != nil {
		fetchChineseFields(nrs, routers)
	}

	svr := beego.BConfig.AppName
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

// Parse total routers and update description on chinese for local server routers,
// then marshal to string for next to push to nacos config server.
func UpdateChineses(data string, descs []*SvrDesc) (string, error) {
	rsmap := make(map[string]*Routers)
	if err := json.Unmarshal([]byte(data), &rsmap); err != nil {
		return "", err
	}

	// check total routers map and input chinese values
	if len(rsmap) == 0 || len(descs) == 0 {
		return "", invar.ErrEmptyData
	}

	// fetch all routers and update chinese
	changed := false
	for _, svr := range descs {
		if routers, ok := rsmap[svr.Server]; ok {
			if routers.CnName != svr.CnName {
				routers.CnName, changed = svr.CnName, true
			}

			// update groups chinese descriptions
			for _, group := range routers.Groups {
				if cnname, ok := svr.Groups[group.Name]; ok {
					if group.CnDesc != cnname {
						group.CnDesc, changed = cnname, true
					}
				}
			}

			// update routers chinese descriptions
			for _, router := range routers.Routers {
				if cnname, ok := svr.Routers[router.Router]; ok {
					router.CnDesc, changed = cnname, true
				}
			}
		}
	}

	// check if exist chinese updated
	if !changed {
		return "", invar.ErrNotChanged
	}

	swagger, err := json.Marshal(rsmap)
	if err != nil {
		return "", err
	}

	logger.D("Updated routers chineses")
	return string(swagger), nil
}

// ----------------------------------------

// Load local server routers from swagger.json file.
func loadSwaggerRouters() (*Routers, error) {
	buff, err := os.ReadFile(swaggerFile)
	if err != nil {
		logger.E("Load swagger routers, err:", err)
		return nil, err
	}

	routers := make(map[string]any)
	if err := json.Unmarshal(buff, &routers); err != nil {
		logger.E("Unmarshal swagger routers err:", err)
		return nil, err
	}
	logger.I("Loaded swagger json, start parse routers")
	out := &Routers{}

	// parse routers by path keyword
	if ps, ok := routers[sfPathName]; ok && ps != nil {
		paths := ps.(map[string]any)

		for path, pathvals := range paths {
			router := &Router{Router: path} // parse router path

			// parse http method, HERE only support GET or POST methods
			var mvs any
			pvs := pathvals.(map[string]any)
			if pmg, ok := pvs[sfMethodGet]; ok && pmg != nil {
				router.Method, mvs = "GET", pmg
			} else if pmp, ok := pvs[sfMethodPost]; ok && pmp != nil {
				router.Method, mvs = "POST", pmp
			} else {
				logger.W("Invalid method of path:", path)
				continue
			}

			// parse beego controller group name
			method := mvs.(map[string]any)
			if gps, ok := method[sfGroupTags]; ok && gps != nil {
				groups := reflect.ValueOf(gps)
				router.Group = groups.Index(0).Interface().(string)
			}

			// parse router path english description
			if desc, ok := method[sfEnDescName]; ok && desc != nil {
				router.EnDesc = desc.(string)
			}

			// append the router into routers array
			// logger.D("> Parsed ["+router.Method+"]\tpath:", path, "\tdesc:", router.EnDesc)
			out.Routers = append(out.Routers, router)
		}
	}

	// parse groups by tags keyword
	if gps, ok := routers[sfGroupTags]; ok && gps != nil {
		groups := gps.([]any) // parse all group array

		for _, group := range groups {
			t := group.(map[string]any)
			gp := &Group{}

			// parse group name value
			if gpn, ok := t[sfGroupName]; ok && gpn != nil {
				gp.Name = gpn.(string)
			}

			// parse group english description
			if gpd, ok := t[sfEnDescName]; ok && gpd != nil {
				gp.EnDesc = strings.TrimRight(gpd.(string), "\n")
			}

			// append the group into groups array
			// logger.D("# Parsed group ["+gp.Name+"] \t desc:", gp.EnDesc)
			out.Groups = append(out.Groups, gp)
		}
	}

	logger.I("Finish parsed swagger routers")
	return out, nil
}

// Parse servers routers from nacos config data, then return all backend services
// routers map and local server swagger routers
func parseNacosRouters(data string) (map[string]*Routers, *Routers) {
	routers := make(map[string]*Routers)
	if data != "" && data != "{}" { // check data if empty
		if err := json.Unmarshal([]byte(data), &routers); err != nil {
			logger.E("Unmarshal nacos routers, err:", err)
			return nil, nil
		}
	}

	if len(routers) > 0 {
		svr := beego.BConfig.AppName
		if rs, ok := routers[svr]; ok {
			logger.D("Parsed nacos routers, found", svr)
			return routers, rs
		}
		logger.D("Parsed nacos routers, unexist", svr)
		return routers, nil
	}

	logger.D("Empty nacos routers, data:", data)
	return routers, nil
}

// Fetch the given routers and groups from src param and set chinese description to dest fileds.
func fetchChineseFields(src *Routers, dest *Routers) {
	dest.CnName = Condition(src.CnName != "", src.CnName, dest.CnName)

	/* -------------------------------- */
	/* cache router chinese description */
	/* -------------------------------- */
	routers := make(map[string]string)
	if len(dest.Routers) > 0 {
		for _, router := range src.Routers {
			if router.CnDesc != "" {
				// logger.D("- Cached router ["+router.Router+"]\tchinese desc:", router.CnDesc)
				routers[router.Router] = router.CnDesc
			}
		}
	}

	// set chinese description to dest routers
	if len(routers) > 0 {
		for _, router := range dest.Routers {
			router.CnDesc = routers[router.Router]
		}
	}

	/* -------------------------------- */
	/* cache groups chinese description */
	/* -------------------------------- */
	groups := make(map[string]string)
	if len(dest.Groups) > 0 {
		for _, group := range src.Groups {
			if group.CnDesc != "" {
				// logger.D("= Cached group ["+group.Name+"]\tchinese desc:", group.CnDesc)
				groups[group.Name] = group.CnDesc
			}
		}
	}

	// set chinese description to dest groups
	if len(groups) > 0 {
		for _, group := range dest.Groups {
			group.CnDesc = groups[group.Name]
		}
	}
}
