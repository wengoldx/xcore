// Copyright (c) 2018-Now Dunyu All Rights Reserved.
//
// Author      : https://www.wengold.net
// Email       : support@wengold.net
//
// Prismy.No | Date       | Modified by. | Description
// -------------------------------------------------------------------
// 00001       2019/05/22   yangping       New version
// -------------------------------------------------------------------

package utils

import (
	"container/list"
	"fmt"
	"sync"

	"github.com/wengoldx/xcore/invar"
)

// Queue the type of queue with sync lock
//
//				--------- <- Head
//	Quere Top : |   1   | -> Pop
//				+-------+
//				|   2   |
//				+-------+
//				|  ...  |
//				+-------+
//		Push -> |   n   | : Queue Back (or Bottom)
//				+-------+
type Queue struct {
	list  *list.List
	mutex sync.Mutex
}

// Create a new queue instance.
func NewQueue() *Queue {
	return &Queue{list: list.New()}
}

// Push push a data to queue back if the data not nil
func (q *Queue) Push(data any) {
	if data == nil {
		return
	}
	q.mutex.Lock()
	defer q.mutex.Unlock()
	q.list.PushBack(data)
}

// Pop pick and remove the front data of queue,
// it will return invar.ErrEmptyData error if the queue is empty
func (q *Queue) Pop() (any, error) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	if e := q.list.Front(); e != nil {
		q.list.Remove(e)
		return e.Value, nil
	}
	return nil, invar.ErrEmptyData
}

// Pick pick but not remove the front data of queue,
// it will return invar.ErrEmptyData error if the queue is empty
func (q *Queue) Pick() (any, error) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	if e := q.list.Front(); e != nil {
		return e.Value, nil
	}
	return nil, invar.ErrEmptyData
}

// Clear clear the queue all data
func (q *Queue) Clear() {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	for e := q.list.Front(); e != nil; {
		en := e.Next()
		q.list.Remove(e)
		e = en
	}
}

// Len return the length of queue
func (q *Queue) Len() int {
	return q.list.Len()
}

/* --------------------------- */
/* Extention Utils             */
/* --------------------------- */

// Fetching result return by callback.
type Result int

// Fetching result to remove item or interupt.
const (
	KEEP_FETCHING   Result = iota // By Default, not remove and keep fetching.
	REMOVE_CONTINUE               // Remove item and continue fetching.
	REMOVE_INTERUPT               // Remove item and interupt fetching.
	REMAIN_INTERUPT               // Found item, interupt fetching.
)

// Fetch queue nodes, use callback returns to remove node or interupt fetch.
// Notice that DO NOT do heavy performence codes in callback, exist a lock here!
//
// For example caller code like (Events is a Queue object):
//
//	// Try fetch waiting requests to send to idle clusters
//	h.Events.Fetch(func(request any) int {
//		if clusterid := h.getIdleCluster(); clusterid != "" {
//			go h.sendRequest(clusterid, request)
//			return utils.REMOVE_CONTINUE
//		}
//
//		// Remain request node and interupt fetch
//		return utils.REMAIN_INTERUPT
//	})
func (q *Queue) Fetch(callback func(value any) Result) {
	if callback != nil {
		q.mutex.Lock()
		defer q.mutex.Unlock()

		for e := q.list.Front(); e != nil; {
			result := callback(e.Value)

			en := e.Next()
			switch result {
			case REMOVE_CONTINUE:
				q.list.Remove(e)
			case REMOVE_INTERUPT:
				q.list.Remove(e)
				return
			case REMAIN_INTERUPT:
				return
			}
			e = en
		}
	}
}

// Dump print out the queue data.
// this method maybe just use for debug to out put queue items
func (q *Queue) Dump() {
	fmt.Println(">>> Dump queue: (front -> back)")
	for e := q.list.Front(); e != nil; e = e.Next() {
		logs := fmt.Sprintf("    : %v", e.Value)
		fmt.Println(logs)
	}
	fmt.Println("<<< End dump!")
}
