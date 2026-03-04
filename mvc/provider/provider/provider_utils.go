// Copyright (c) 2018-Now Dunyu All Rights Reserved.
//
// Author      : https://www.wengold.net
// Email       : support@wengold.net
//
// Prismy.No | Date       | Modified by. | Description
// -------------------------------------------------------------------
// 00001       2019/05/22   yangping       New version
// -------------------------------------------------------------------

package provider

import (
	"database/sql"

	"github.com/wengoldx/xcore/invar"
	"github.com/wengoldx/xcore/mvc/provider/builder"
)

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
func QueryColumn[T any](b *builder.QueryBuilder, outs *[]T) error {
	if outs == nil || b == nil || !b.HasProvider(){
		return invar.ErrBadDBConnect
	}

	return b.Query(func(rows *sql.Rows) error {
		var v T
		if err := rows.Scan(&v); err != nil {
			return err
		}
		*outs = append(*outs, v)
		return nil
	})
}