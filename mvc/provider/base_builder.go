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

	"github.com/wengoldx/xcore/logger"
	"github.com/wengoldx/xcore/utils"
)

// A interface implement by CUDA builder to build
// a sql string for database access.
type SQLBuilder interface {
	Build() (string, []any) // Build sql string and return args.
}

// The base builder to support util methods to simple build a
// sql string for database CUDA actions.
type BaseBuilder struct {
	master *TableProvider // Table provider for execute sql actions.
}

var _ SQLBuilder = (*BaseBuilder)(nil)

// Empty build method, override by CUDA builders.
func (b *BaseBuilder) Build() (string, []any) {
	return "", []any{} // not implement Build method.
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
			logger.D("- Where:", condition, "-", "arg:", arg)
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
			switch c := strings.TrimSpace(strings.ToUpper(sep[0])); c {
			case "AND", "OR":
				connector = " " + c + " "
			case "":
				connector = " "
			}
		}
		where = "WHERE " + strings.Join(conditions, connector)
	}
	logger.D("Formated where:", where, "with arg:", args)
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
		values := strings.Join(utils.ToStrings(args), ",")
		return fmt.Sprintf("%s IN (%s)", field, values)
	}
	return ""
}

// Format order by condition to string.
//
//	- desc = true : ORDER BY field DESC
//	- desc = false: ORDER BY field ASC
func (b *BaseBuilder) FormatOrder(field string, desc ...bool) string {
	if field != "" {
		isdesc := utils.Variable(desc, true) // default for DESC.
		order := utils.Condition(isdesc, "DESC", "ASC")
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

// Fetch the KValues items and return the joined fields, ? holders, and args.
//
//	values := KValues{
//		"":       123456,   // Filter out empty field
//		"Age":    16,
//		"Male":   true,
//		"Name":   "ZhangSan",
//		"Height": 176.8,
//		"Secure": nil,      // Filter out nil value
//	}
//	// => Age=?, Male=?, Name=?, Height=?
//	// => ?,?,?,?
//	// => []any{16, true, "ZhangSan", 176.8}
func (b *BaseBuilder) FormatInserts(values KValues) (string, string, []any) {
	fields, holders, args := "", "", []any{}
	if cnt := len(values); cnt > 0 {
		tags := []string{}
		for key, arg := range values {
			tags = append(tags, key)
			args = append(args, arg)
		}

		fields = strings.Join(tags, ", ")
		holders = strings.Repeat("?,", cnt)
		holders = strings.TrimSuffix(holders, ",")
	}
	return fields, holders, args
}

// Fetch the KValues items and return the joined fields and args.
//
//	values := KValues{
//		"":       123456,   // Filter out empty field
//		"Age":    16,
//		"Male":   true,
//		"Name":   "ZhangSan",
//		"Height": 176.8,
//		"Secure": nil,      // Filter out nil value
//	}
//	// => Age=?, Male=?, Name=?, Height=?
//	// => []any{16, true, "ZhangSan", 176.8}
func (b *BaseBuilder) FormatSets(values KValues) (string, []any) {
	fields, args := "", []any{}
	if cnt := len(values); cnt > 0 {
		sets := []string{}
		for key, arg := range values {

			sets = append(sets, key+"=?")
			args = append(args, arg)
		}

		fields = strings.Join(sets, ", ")
	}
	return fields, args
}

// Fetch the KValues items and return the formated sets string.
//
//	values := KValues{
//		"":       123456,   // Filter out empty field
//		"Age":    16,
//		"Male":   true,
//		"Name":   "ZhangSan",
//		"Height": 176.8,
//		"Secure": nil,      // Filter out nil value
//	}
//	// => Age=16, Male=true, Name='ZhangSan', Height=176.8
func (p *BaseBuilder) FormatValues(values KValues) string {
	sets := []string{}
	for key, value := range values {
		if key != "" && value != nil {
			switch v := value.(type) {
			case string:
				sets = append(sets, key+"='"+v+"'")
			default:
				sets = append(sets, fmt.Sprintf(key+"=%v", v))
			}
		}
	}
	return strings.Join(sets, ",")
}

// Ensure where condition prefixed 'WHERE' keyword when not empty.
func (b *BaseBuilder) CheckWhere(wheres string) string {
	wheres = strings.TrimSpace(wheres)
	if wheres != "" && !strings.HasPrefix(wheres, "WHERE ") {
		wheres = "WHERE " + wheres
	}
	return wheres
}

// Ensure query string must tail 'LIMIT 1' for query the top one record.
func (b *BaseBuilder) CheckLimit(query string) string {
	query = strings.TrimSpace(query)
	if query != "" && !strings.HasSuffix(query, "LIMIT 1") &&
		!strings.HasSuffix(query, "limit 1") {
		query += " LIMIT 1"
	}
	return query
}

// Build where conditions, append where ins, like conditions if exist.
//
//	- WHERE wheres
//	- WHERE wheres AND field IN (v1,v2...)
//	- WHERE wheres AND field IN (v1,v2...) AND field2 LIKE '%%filter%%'
//	- WHERE field IN (v1,v2...) AND field2 LIKE '%%filter%%'
//	- WHERE field LIKE '%%filter%%'
//
// Use FormatWheres(), FormatWhereIn() to format Wheres data or where in condition.
func (b *BaseBuilder) BuildWheres(wheres Wheres, ins, like string, sep ...string) (string, []any) {
	where, args := b.FormatWheres(wheres, sep...) // WHERE wheres
	if where != "" {
		// WHERE wheres AND field IN (v1,v2...)
		if ins != "" {
			where += " AND " + ins
		}

		// WHERE wheres AND field IN (v1,v2...) AND field2 LIKE '%%filter%%'
		if like != "" {
			where += " AND " + like
		}
	} else {
		if ins != "" {
			// WHERE field IN (v1,v2...) AND field2 LIKE '%%filter%%'
			where = "WHERE " + ins
			if like != "" {
				where += " AND " + like
			}
		} else if like != "" {
			// WHERE field LIKE '%%filter%%'
			where = "WHERE " + like
		}
	}
	return where, args
}

// Join the given where conditions without input AND and OR connectors.
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
