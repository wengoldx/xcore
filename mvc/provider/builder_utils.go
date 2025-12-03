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

type FieldsFunc[T any] func(iv *T) []any

type ItemCreator[T any] struct {
	OutsFunc FieldsFunc[T]
}

var _ SQLCreator = (*ItemCreator[any])(nil)

func NewCreator[T any](cb FieldsFunc[T]) *ItemCreator[T] {
	return &ItemCreator[T]{OutsFunc : cb}
}

func (ic *ItemCreator[T]) NewItem() []any {
	 var item T; 
	 return ic.OutsFunc(&item)
 }