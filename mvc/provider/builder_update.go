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

// Build a query string for sql update.
//
//	UPDATE table
//		SET v1=?, v2=?, v3=?...
//		WHERE wherer AND field IN (v1,v2...) AND field2 LIKE '%%filter%%'
//
// See QueryBuilder, InsertBuilder, DeleteBuilder.
type UpdateBuilder struct {
	BaseBuilder

	table  string  // Table name for update
	values KValues // Target fields and values to update.
	wheres Wheres  // Where conditions and args values.
	sep    string  // Where conditions connector, one of 'AND', 'OR', ' ', default ''.
	ins    string  // Where in conditions.
	like   string  // Like conditions string.
}

var _ SQLBuilder = (*UpdateBuilder)(nil)

// Create a UpdateBuilder instance to build a query string.
func NewUpdate(table string) *UpdateBuilder {
	return &UpdateBuilder{table: table}
}

/* ------------------------------------------------------------------- */
/* SQL Action Utils By Using master Provider                           */
/* ------------------------------------------------------------------- */

func (b *UpdateBuilder) Exec() error   { return b.master.Exec(b) }   // Update target record without check.
func (b *UpdateBuilder) Update() error { return b.master.Update(b) } // Update target record and check changes counts.

/* ------------------------------------------------------------------- */
/* SQL Action Builder Methonds                                         */
/* ------------------------------------------------------------------- */

// Specify master provider.
func (b *UpdateBuilder) Master(master *TableProvider) *UpdateBuilder {
	b.master = master
	return b
}

// Specify the target table for query.
func (b *UpdateBuilder) Table(table string) *UpdateBuilder {
	b.table = table
	return b
}

// Specify the values of row to update.
//
//	values := KValues{
//		"":       123456,   // Filter out empty field
//		"Age":    16,
//		"Male":   true,
//		"Name":   "ZhangSan",
//		"Height": 176.8,
//		"Secure": nil,      // Filter out nil value
//	}
//	// => SET Age=?, Male=?, Name=?, Height=?
//	// => values: []any{16, true, "ZhangSan", 176.8}
func (b *UpdateBuilder) Values(row KValues) *UpdateBuilder {
	b.values = row
	return b
}

// Specify the where conditions and args for query.
//
//	where = provider.Wheres{
//		"acc=?":"123", "age>=?":20, "role<>?":"admin",
//	}
//	// => WHERE acc=? AND age>=? AND role<>?
//	// => args ("123", 20, "admin")
func (b *UpdateBuilder) Wheres(where Wheres) *UpdateBuilder {
	b.wheres = where
	return b
}

// Specify the where in condition with field and args for query.
//
//	builder.WhereIn("id", []any{1, 2}) // => WHERE id IN (1, 2)
func (b *UpdateBuilder) WhereIn(field string, args []any) *UpdateBuilder {
	b.ins = b.FormatWhereIn(field, args)
	return b
}

// Specify the where in condition with field and args for query.
func (b *UpdateBuilder) WhereSep(sep string) *UpdateBuilder {
	switch s := strings.ToUpper(sep); s {
	case "AND", "OR", " " /* for none where connector */ :
		b.sep = s
	}
	return b
}

// Specify the like condition for query.
//
//	builder.Like("acc", "zhang") // => acc LIKE '%%zhang%%'
func (b *UpdateBuilder) Like(field, filter string) *UpdateBuilder {
	b.like = b.FormatLike(field, filter)
	return b
}

// Build the update action sql string and args for provider to update datas.
//
//	UPDATE table
//		SET v1=?, v2=?, v3=?...
//		WHERE wherer AND field IN (v1,v2...) AND field2 LIKE '%%filter%%'
func (b *UpdateBuilder) Build() (string, []any) {
	sep := utils.Condition(b.sep == "", "AND", b.sep)

	tags, args := b.FormatSets(b.values)                      // SET v1=?,v2=?...
	where, wvs := b.BuildWheres(b.wheres, b.ins, b.like, sep) // WHERE wheres AND field IN (v1,v2...) AND field2 LIKE '%%filter%%'
	args = append(args, wvs...)

	query := "UPDATE %s SET %s %s"
	query = fmt.Sprintf(query, b.table, tags, where)
	query = strings.TrimSuffix(query, " ")
	return query, args
}

// Reset builder datas for next prepare and build.
func (b *UpdateBuilder) Reset() *UpdateBuilder {
	clear(b.values)
	clear(b.wheres)
	b.sep, b.ins, b.like = "", "", ""
	return b
}
