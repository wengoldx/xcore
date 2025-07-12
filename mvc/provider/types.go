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
)

// A callback for scan query result records, it will interrupt
// scanning when callback return error.
type ScanCallback func(rows *sql.Rows) error

// A callback for format insert values as string to insert record.
type InsertCallback func(index int) string

// A callback for handle transaction by call provider.Trans().
type TransCallback func(tx *sql.Tx) error

// A callback for single record query finaly notify.
type DoneCallback func()

/* ------------------------------------------------------------------- */
/* Key-Value                                                           */
/* ------------------------------------------------------------------- */

// Fields name and referened values for insert or update.
type KValues map[string]any

// Append a key-value into KValues.
func (v *KValues) Add(key string, value any) *KValues {
	(*v)[key] = value
	return v
}

// Remove target key-value out from KValues.
func (v *KValues) Remove(key string) *KValues {
	delete(*v, key)
	return v
}

/* ------------------------------------------------------------------- */
/* Wheres                                                              */
/* ------------------------------------------------------------------- */

// Fields name and referened values for construct where condition string.
type Wheres map[string]any

// Append a where condition into Wheres.
func (w *Wheres) Add(condition string, arg any) *Wheres {
	(*w)[condition] = arg
	return w
}

// Remove target where condition out from Wheres.
func (w *Wheres) Remove(condition string) *Wheres {
	delete(*w, condition)
	return w
}
