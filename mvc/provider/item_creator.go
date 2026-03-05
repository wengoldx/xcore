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

// A interface implement for query array rows with struct object type,
// and only for QueryBuilder.Array() method.
//
// Call NewCreator() to create a item factory to create items, or with
// parser function to parse each item on scaning.
type Creator interface {
	CreateItem() (any, []any) // New a item and return out params.
	AppendItem(iv any) error  // Parse item values after row scaned.
}

type GetterFunc[T any] func(iv *T) []any
type ParserFunc[T any] func(iv *T) error

// New a Creator instance to generate target module items object.
//
//	datas := []*types.Account{
//		UID, Email string	
//	}
//	
//	// UseAage 1: only for query array records.
//	h.Querier().Tags("uid", "email").Wheres(...).
//	  Array(pd.NewCreator(&datas, func(iv *types.Account) []any {
//		  return []any{&iv.UID, &iv.Email}
//	  }))
//
//	// UseAage 2: query array records and parse item value on scaning.
//	h.Querier().Tags("uid", "email").Wheres(...).
//	  Array(pd.NewCreator(&datas, func(iv *types.Account) []any {
//		  return []any{&iv.UID, &iv.Email}
//	  }, func(iv *types.Account) {
//		  iv.Email = decode(iv.Email)
//	  }))
func NewCreator[T any](outs *[]*T, getter GetterFunc[T], parser ...ParserFunc[T]) *ItemCreator[T] {
	creator := &ItemCreator[T]{outs:outs, getFunc : getter}
	if len(parser) > 0 && parser[0] != nil {
		creator.parseFunc = parser[0]
	}
	return creator
}

// Row record data struct creator, for create and parse item.
type ItemCreator[T any] struct {
	outs      *[]*T
	getFunc   GetterFunc[T]
	parseFunc ParserFunc[T]
}

var _ Creator = (*ItemCreator[any])(nil)

// New a item and return out values.
func (ic *ItemCreator[T]) CreateItem() (any, []any) {
	 var item T;
	 return &item, ic.getFunc(&item)
 }

// Parse item scaned values if parser exist.
func (ic *ItemCreator[T]) AppendItem(iv any) error {
	if iv != nil {
		if item, ok := iv.(*T); ok {
			*ic.outs = append(*ic.outs, item)
			if ic.parseFunc != nil {
				return ic.parseFunc(item)
			}
			return nil
		}
	}
	return invar.ErrInvalidData
 }
