// Copyright (c) 2018-Now Dunyu All Rights Reserved.
//
// Author      : https://www.wengold.net
// Email       : support@wengold.net
//
// Prismy.No | Date       | Modified by. | Description
// -------------------------------------------------------------------
// 00001       2019/05/22   yangping       New version
// -------------------------------------------------------------------

package tester

import (
	pd "github.com/wengoldx/xcore/mvc/provider"
	"github.com/wengoldx/xcore/utils"
)

// Unit test helper options setter.
type Option func(t *helper)

// Apply all unit test helper options settngs.
func applyOptions(t *helper, opts ...Option) {
	for _, optFunc := range opts {
		optFunc(t)
	}
}

// Special target field name as input param or output field,
// such as uid.
func WithTag(tag string) Option {
	return func(u *helper) { u.tag = tag }
}

// Special where condition fields and values.
func WithWhere(where pd.Wheres) Option {
	return func(u *helper) { u.where = where }
}

// Special limit number for query result.
func WithLimit(limit int) Option {
	return func(u *helper) { u.limit = limit }
}

// Special order by conditions.
func WithOrder(order string, desc ...string) Option {
	return func(u *helper) {
		u.desc = utils.VarString(desc, "DESC")
		u.order = order
	}
}
