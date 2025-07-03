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
	"fmt"
	"strings"

	"github.com/wengoldx/xcore/utils"
)

type BaseBuilder struct {
}

// Format where conditions to string with args, by default join conditions
// with AND connector, but can change to OR or empty connector by set 'ornnoe'
// param.
//
//	- ornone not set  : use AND connector.
//	- ornone set true : use OR  connector.
//	- ornone set false: not use any connector, by inset with condition as 'condition AND', 'condition OR'.
func (b *BaseBuilder) FormatWheres(wheres Wheres, ornone ...bool) (string, []any) {
	where, args := "", []any{}
	if len(wheres) > 0 {
		conditions := []string{}
		for condition, arg := range wheres {
			conditions = append(conditions, condition)
			args = append(args, arg)
		}

		// join conditions as:
		//
		// - WHERE  condition1 AND   condition2 AND condition3.
		// - WHERE  condition1 OR    condition2 OR  condition3.
		// - WHERE 'condition1 AND' 'condition2 OR' condition3.
		sep := " AND "
		if len(ornone) > 0 {
			sep = utils.Condition(ornone[0], " OR ", " ").(string)
		}
		where = "WHERE " + strings.Join(conditions, sep)
	}
	return where, args
}

// Format order by condition to string.
func (b *BaseBuilder) FormatOrder(field string, desc bool) string {
	if field != "" {
		order := utils.Condition(desc, "DESC", "ASC").(string)
		return fmt.Sprintf("ORDER BY %s %s", field, order)
	}
	return ""
}

// Format limit condition to string.
func (b *BaseBuilder) FormatLimit(limit int) string {
	if limit > 0 {
		return fmt.Sprintf("LIMIT %d", limit)
	}
	return ""
}
