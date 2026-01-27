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

import (
	"database/sql"

	"github.com/wengoldx/xcore/invar"
)

// A interface implement by array elems creator to return
// out values of columns.
//
// It only for QueryBuilder.Array().
type ModuleCreator interface {

	// Create a new item and return out values.
	Generate() []any
}

// Return module fields pointers for scan query results.
type GetFields[T any] func(iv *T) []any

// A table data module struct as ORM object.
type Module[T any] struct {
	OutsFunc GetFields[T]
}

var _ ModuleCreator = (*Module[any])(nil)

// Create a ModuleCreator instance to generate target module items object.
func NewCreator[T any](cb GetFields[T]) *Module[T] {
	return &Module[T]{OutsFunc : cb}
}

// Create a new item and return out values.
func (ic *Module[T]) Generate() []any {
	 var item T; 
	 return ic.OutsFunc(&item)
 }

 /* ------------------------------------------------------------------- */
/* Util Methods For package callable                                   */
/* ------------------------------------------------------------------- */

// Query the target column values and return array by callback.
//
// # USAGE:
//
//	// case 1: new a builder and set exist provider.
//	files := []string{}
//	builder := pd.NewQuery("mytable").Master(myprovider)
//	err := pd.QueryColumn(builder.Tags("file").Wheres(pd.Wheres{"uid": uid}), &files)
//
//	// case 2: or, use exist provider get builder.
//	builder = myprovider.Querier()
func QueryColumn[T any](builder QueryBuilder, outs *[]T) error {
	if outs == nil || builder == nil || !builder.HasProvider() {
		return invar.ErrBadDBConnect
	}

	return builder.Query(func(rows *sql.Rows) error {
		var v T
		if err := rows.Scan(&v); err != nil {
			return err
		}
		*outs = append(*outs, v)
		return nil
	})
}