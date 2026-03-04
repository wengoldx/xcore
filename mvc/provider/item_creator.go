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

// A interface implement by array elems creator to return
// out values of columns.
//
// It only for QueryBuilder.Array().
type Creator interface {

	// New a item and return out values.
	CreateItem() []any
}

// New a Creator instance to generate target module items object.
//
//	datas := []*types.Account{}
//	err := h.Querier().Tags("uid", "email").Wheres(...).
//		Array(pd.NewCreator(func(iv *types.Account) []any {
//			datas = append(datas, iv)
//			return []any{&iv.UID, &iv.Email}
//		}))
func NewCreator[T any](cb func(iv *T) []any) *ItemCreator[T] {
	return &ItemCreator[T]{getItemOuts : cb}
}

// A table data module struct as ORM object.
type ItemCreator[T any] struct {
	getItemOuts func(iv *T) []any
}

var _ Creator = (*ItemCreator[any])(nil)

// New a item and return out values.
func (ic *ItemCreator[T]) CreateItem() []any {
	 var item T; 
	 return ic.getItemOuts(&item)
 }

