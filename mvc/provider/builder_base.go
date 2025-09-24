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
	"fmt"
	"strings"

	"github.com/wengoldx/xcore/utils"
)

// The base builder to support util methods to simple build a
// sql string for database CUDA actions.
type BaseBuilder struct {
	master *TableProvider // Table provider for execute sql actions.
}

var _ SQLBuilder = (*BaseBuilder)(nil)

// Empty build method, override by CUDA builders.
func (b *BaseBuilder) Build(debug ...bool) (string, []any) {
	return "", []any{} // not implement Build method.
}

// Format table joins to string for multi-table query, it will filter out the
// empty table or alias join datas.
//
//	tables := provider.Joins{
//		"account":"a", "profile":"b", "other":"", // the 'other' table will filter out!
//	}
//	joins := builder.FormatJoins(tables)
//	fmt.Println(joins) // => account AS a, profile AS b
func (b *BaseBuilder) FormatJoins(tables Joins) string {
	ts := []string{}
	for table, alias := range tables {
		table, alias = strings.TrimSpace(table), strings.TrimSpace(alias)
		if table != "" && alias != "" {
			ts = append(ts, fmt.Sprintf("%s AS %s", table, alias))
		}
	}
	return strings.Join(ts, ", ")
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
			conditions = append(conditions, condition) // append conditions whatever arg is nil.
			if arg != nil {                            // filter out the nil args, it useful for where joins like 'a.acc=b.user'.
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

// Format like condition to string, set pattern one of 'perfix', 'suffix', 'center'
// to make diffrent filter string as follow, by default use 'center' pattern.
//
//	- Perfix pattern: field LIKE 'filter%%'
//	- Center pattern: field LIKE '%%filter%%'
//	- Suffix pattern: field LIKE '%%filter'
func (b *BaseBuilder) FormatLike(field, filter string, pattern ...string) string {
	if field != "" && filter != "" {
		lower := strings.ToLower(utils.Variable(pattern, "center"))
		switch lower {
		case "perfix":
			return field + " LIKE '" + filter + "%%'"
		case "suffix":
			return field + " LIKE '%%" + filter + "'"
		default:
			return field + " LIKE '%%" + filter + "%%'"
		}
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
//	// => Age=?, Male=?, Name=?, Height=?, Secure=?
//	// => ?,?,?,?
//	// => []any{16, true, "ZhangSan", 176.8}
//
// WARNING: This method not well support insert nil value by arg, the nil
// value will inserted like '<nil>' string, not NULL value;
func (b *BaseBuilder) FormatInserts(values KValues) (string, string, []any) {
	fields, holders, args := "", "", []any{}
	if cnt := len(values); cnt > 0 {
		tags := []string{}
		for key, arg := range values {
			if key == "" { // filter out the empty field key.
				continue
			}

			// FIXME: The nil arg will be insert like '<nil>' string by
			// arg, so DO NOT insert nil values if you can if possible!
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
//		"Secure": nil,      // Set value as NULL
//	}
//	// => Age=?, Male=?, Name=?, Height=?, Secure=NULL
//	// => []any{16, true, "ZhangSan", 176.8}
//
// WARNING: This method support update nil arg as NULL value.
func (b *BaseBuilder) FormatSets(values KValues) (string, []any) {
	fields, args := "", []any{}
	if cnt := len(values); cnt > 0 {
		sets := []string{}
		for key, arg := range values {
			if key == "" { // filter out the empty field key.
				continue
			}

			// FIXME: The nil arg will be translate to NULL value
			// for single row or multiple rows insert.
			if arg == nil {
				sets = append(sets, key+"=NULL")
				continue
			}
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
