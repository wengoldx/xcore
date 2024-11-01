// Copyright (c) 2018-Now Dunyu All Rights Reserved.
//
// Author      : https://www.wengold.net
// Email       : support@wengold.net
//
// Prismy.No | Date       | Modified by. | Description
// -------------------------------------------------------------------
// 00001       2022/02/09   yangping       New version
// -------------------------------------------------------------------

package wsio

import (
	"sync"
)

// Rooms
type WRooms struct {
	lock     sync.Mutex
	Staffers map[string][]string // Staffer client array group by exhibitor aid
	Targets  map[string][]string // client@tagid id to room id
}

var wrm *WRooms

func init() {
	wrm = &WRooms{
		Staffers: make(map[string][]string),
		Targets:  make(map[string][]string),
	}
	siolog.I("Init socket rooms!")
}

// Return WRooms singleton
func GetRooms() *WRooms {
	return wrm
}
