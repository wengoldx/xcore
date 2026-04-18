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

	"github.com/wengoldx/xcore/utils"
)

// A callback for scan query result records, it will interrupt
// scanning when callback return error.
type ScanCallback func(rows *sql.Rows) error

// A callback for format insert values as string to insert record.
type InsertCallback func(index int) string

// A callback for handle transaction by call BaseProvider.Trans().
type TransCallback func(tx *sql.Tx) error

// A callback for handle transaction by call TableProvider.Trans().
type TranerCallback func(t *Traner) error

// A callback for single record query finaly notify.
type DoneCallback func()

/* ------------------------------------------------------------------- */
/* Table-Alias typed data for Join-Query                               */
/* ------------------------------------------------------------------- */

// Table name with join alias for multi-table join,
// the datas format as table:alias.
//
//	account:a // table name "account", alias "a".
//
// # WARNING:
//	- None check the duplicate alias error for multiple tables!
//	- The alias will be overwritten when table name same!
type Joins map[string]string

// Append a table-alias into table Joins.
func (t *Joins) Add(table, alias string) *Joins {
	(*t)[table] = alias
	return t
}

// Remove target table out from table Joins.
func (t *Joins) Remove(table string) *Joins {
	delete(*t, table)
	return t
}

/* ------------------------------------------------------------------- */
/* Key-Value typed data for Insert, Update                             */
/* ------------------------------------------------------------------- */

// Fields name and referened values for insert or update,
// the datas format as column:value.
//
//	user_name:'xiaoming' // column 'user_name', value 'xiaoming'.
//
// # WARNING:
//	- None check the duplicate alias error for multiple tables!
//	- The alias will be overwritten when table name same!
type KValues map[string]any

// Append a key-value into KValues.
func (v *KValues) Add(key string, value any) *KValues {
	(*v)[key] = value
	return v
}

// Append mutiple key-value into KValues.
func (v *KValues) Adds(values KValues) *KValues {
	for key, value := range values {
		(*v)[key] = value
	}
	return v
}

// Remove target key-value out from KValues.
func (v *KValues) Remove(keys ...string) *KValues {
	for _, key := range keys {
		delete(*v, key)
	}
	return v
}

/* ------------------------------------------------------------------- */
/* Wheres typed data for Query, Update, Delete                         */
/* ------------------------------------------------------------------- */

// Fields name and referened values for construct where condition string,
// the datas format as condition:value, for example:
//
//	uid=?:'123456'     // condition is "uid=?", the value "123456".
//	uid='123456':nil   // condition is "uid='123456'" without value.
//	t1.uid=t2.user:nil // confition for joined tables where string.
type Wheres map[string]any

// Append a where condition into Wheres.
func (w *Wheres) Add(condition string, arg any) *Wheres {
	(*w)[condition] = arg
	return w
}

// Append multiple where conditions into Wheres.
func (w *Wheres) Adds(conditions Wheres) *Wheres {
	for condition, arg := range conditions {
		(*w)[condition] = arg
	}
	return w
}

// Remove target where condition out from Wheres.
func (w *Wheres) Remove(conditions ...string) *Wheres {
	for _, condition := range conditions {
		delete(*w, condition)
	}
	return w
}

// Target column name and args for build where in condition,
// call pd.NewIn() to create it.
type In struct {
	field string // target column name.
	args  []any  // where in values.
}

// Create a where in condition object.
func NewIn[T any](field string, args []T) *In {
 	return &In{field: field, args: utils.ToAnys(args)}
}

// Return the where in field and args.
func (w *In) Get() (string, []any) {
	return w.field, w.args
}
