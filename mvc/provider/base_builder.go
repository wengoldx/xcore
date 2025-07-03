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
	"strconv"
	"strings"

	"github.com/wengoldx/xcore/utils"
)

type BaseBuilder struct {
}

// Fetch the KValues items and return the joined fields, ? holders, and args.
func (b *BaseBuilder) FormatInserts(values KValues) (string, string, []any) {
	fields, holders, args := "", "", []any{}
	if cnt := len(values); cnt > 0 {
		outs := []string{}
		for key, arg := range values {
			outs = append(outs, key)
			args = append(args, arg)
		}

		fields = strings.Join(outs, ", ")
		holders = strings.TrimSuffix(strings.Repeat("?,", cnt), ",")
	}
	return fields, holders, args
}

// Format where conditions to string with args, by default join conditions with
// AND connector, but can change to OR or empty connector by set 'connector' param.
//
//	- not set or set AND : use AND connector.
//	- set OR             : use OR  connector.
//	- set empty string   : tail connector inside where condition like 'condition AND', 'condition OR'.
//
// WARNING: Here will filter out the nil values in wheres.
func (b *BaseBuilder) FormatWheres(wheres Wheres, sep ...string) (string, []any) {
	where, args := "", []any{}
	if len(wheres) > 0 {
		conditions := []string{}
		for condition, arg := range wheres {
			if arg != nil {
				conditions = append(conditions, condition)
				args = append(args, arg)
			}
		}

		// join conditions as:
		//
		// - WHERE condition1 AND condition2 AND condition3
		// - WHERE condition1 OR  condition2 OR  condition3
		// - WHERE condition1 AND condition2 OR  condition3
		connector := " AND "
		if len(sep) > 0 {
			switch c := strings.ToUpper(sep[0]); c {
			case "AND", "OR":
				connector = " " + c + " "
			case "":
				connector = " "
			}
		}
		where = "WHERE " + strings.Join(conditions, connector)
	}
	return where, args
}

// Format where in condition to string without perfix 'WHERE' keyword.
//
//	- int number args  : field IN (1,2,3)
//	- float number args: field IN (1.2,2.3,3.45)
//	- string args      : field IN ('1','2','3')
//
// WARNING: Here will filter out the nil values in args.
func (b *BaseBuilder) FormatWhereIn(field string, args []any) string {
	if len(args) > 0 {
		values := strings.Join(b.ToStrings(args), ",")
		return fmt.Sprintf("%s IN (%s)", field, values)
	}
	return ""
}

// Format order by condition to string.
//
//	- desc = true : ORDER BY field DESC
//	- desc = false: ORDER BY field ASC
func (b *BaseBuilder) FormatOrder(field string, desc bool) string {
	if field != "" {
		order := utils.Condition(desc, "DESC", "ASC").(string)
		return fmt.Sprintf("ORDER BY %s %s", field, order)
	}
	return ""
}

// Format limit condition to string.
//
//	- output string: LIMIT n
func (b *BaseBuilder) FormatLimit(n int) string {
	if n > 0 {
		return fmt.Sprintf("LIMIT %d", n)
	}
	return ""
}

// Format like condition to string.
//
//	- output string: field LIKE '%%filter%%'
func (b *BaseBuilder) FormatLike(field, filter string) string {
	if field != "" && filter != "" {
		return field + " LIKE '%%" + filter + "%%'"
	}
	return ""
}

// Ensure where condition prefixed 'WHERE' keyword when not empty.
func (b *BaseBuilder) CheckWhere(wheres string) string {
	wheres = strings.TrimSpace(wheres)
	if wheres != "" && !strings.HasPrefix(wheres, "WHERE") {
		wheres = "WHERE " + wheres
	}
	return wheres
}

// Translate build-in types values to strings, it only support the types as follow,
// or return empty string array when contain any unsupport types value.
//
//	- int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64
//	- float32, float64
//	- bool
//	- string
//
// WARNING: Here will filter out the nil values in wheres.
func (b *BaseBuilder) ToStrings(values []any) []string {
	vs := []string{}
	for _, value := range values {
		if value == nil {
			continue
		}

		switch v := value.(type) {
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
			vs = append(vs, fmt.Sprintf("%d", v))
		case float32:
			vs = append(vs, strconv.FormatFloat(float64(v), 'f', -1, 64))
		case float64:
			vs = append(vs, strconv.FormatFloat(v, 'f', -1, 64))
		case bool:
			vs = append(vs, utils.Condition(v, "true", "false").(string))
		case string:
			vs = append(vs, "'"+v+"'") // 'value'
		default:
			return []string{}
		}
	}
	return vs
}

// Translate string number array to any type array.
func (b *BaseBuilder) ToAnys(values []string) []any {
	args := []any{}
	for _, value := range values {
		args = append(args, value)
	}
	return args
}

// Translate int number array to any type array.
func (b *BaseBuilder) IntAnys(values []int) []any {
	args := []any{}
	for _, value := range values {
		args = append(args, value)
	}
	return args
}

// Translate int64 number array to any type array.
func (b *BaseBuilder) Int64Anys(values []int64) []any {
	args := []any{}
	for _, value := range values {
		args = append(args, value)
	}
	return args
}

// Translate float64 number array to any type array.
func (b *BaseBuilder) FloatAnys(values []float64) []any {
	args := []any{}
	for _, value := range values {
		args = append(args, value)
	}
	return args
}

// Join values as string like "1,2.3,'456',true", or append the values
// string into query strings, the input params as formart:
//
//	- values: []any{1, 2.3, "456", true}
//	- query : "SELECT * FROM tablename WHERE id IN (%s)"
//
// The result is "SELECT * FROM tablename WHERE id IN (1,2.3,'456',true)".
//
//	WARNING: The values only support int, int64, float64, bool, string types!
func (b *BaseBuilder) Joins(values []any, query ...string) string {
	vs := b.ToStrings(values)
	if len(vs) > 0 {
		// Append values into none-empty query string
		if q := utils.VarString(query, ""); q != "" {
			return fmt.Sprintf(q, strings.Join(vs, ","))
		}
		return strings.Join(vs, ",")
	}
	return ""
}

// Join int64 values as string like "1,2,3".
//
// See Joins() method for link different types values.
func (b *BaseBuilder) JoinInts(values []int64, query ...string) string {
	return b.Joins(b.Int64Anys(values), query...)
}

// Join string values as string like "'1','2','3'".
//
// See Joins() method for link different types values.
func (b *BaseBuilder) JoinStrings(values []string, query ...string) string {
	return b.Joins(b.ToAnys(values), query...)
}

// Join the given where conditions without input AND and OR connectors.
//
// Set FormatWheres() method to known more where connectors.
func (b *BaseBuilder) JoinWheres(wheres ...string) string {
	return strings.Join(wheres, " ")
}

// Join the given where conditions with input AND connectors.
func (b *BaseBuilder) JoinAndWheres(wheres ...string) string {
	return strings.Join(wheres, " AND ")
}

// Join the given where conditions with input OR connectors.
func (b *BaseBuilder) JoinOrWheres(wheres ...string) string {
	return strings.Join(wheres, " OR ")
}
