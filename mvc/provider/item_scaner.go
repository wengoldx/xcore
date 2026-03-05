// Copyright (c) 2018-Now Dunyu All Rights Reserved.
//
// Author      : https://www.wengold.net
// Email       : support@wengold.net
//
// Prismy.No | Date       | Modified by. | Description
// -------------------------------------------------------------------
// 00001       2019/05/22   yangping       New version
// -------------------------------------------------------------------

package pd

import "github.com/wengoldx/xcore/invar"

// A interface implement for query single column all values,
// it only for QueryBuilder.Column() method.
//
// Call NewScaner() to create a values scaner, or with parser
// function to parse each item value on scaning.
type Scaner interface {
	CreateItem() any         // New a item and return out params.
	AppendItem(iv any) error // Parse item values after row scaned.
}

// New a Scaner instance to get single column all values.
//
//	emails := []string{}
//	h.Querier().Tags("uid").Wheres(...).
//	  Column(pd.NewScaner(&emails, func(iv *string) {
//		  *iv = decode(*iv)
//	  }))
func NewScaner[T any](outs *[]T, parser ...ParserFunc[T]) Scaner {
	scaner := &ItemScaner[T]{outs:outs}
	if len(parser) > 0 && parser[0] != nil {
		scaner.parseFunc = parser[0]
	}
	return scaner
}

// Row record value scaner, for read and parse item.
type ItemScaner[T any] struct {
	outs      *[]T
	parseFunc ParserFunc[T]
}

var _ Scaner = (*ItemScaner[any])(nil)

// Return a new row value.
func (ic *ItemScaner[T]) CreateItem() any {
	 var item T;
	 return &item
 }

// Parse item value and append into outs array.
func (ic *ItemScaner[T]) AppendItem(iv any) error {
	if iv != nil {
		if item, ok := iv.(*T); ok {
			if ic.parseFunc != nil {
				ic.parseFunc(item);
			}
			*ic.outs = append(*ic.outs, *item)
			return nil
		}
	}
	return invar.ErrInvalidData
 }

