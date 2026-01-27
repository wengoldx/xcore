// Copyright (c) 2018-Now Dunyu All Rights Reserved.
//
// Author      : https://www.wengold.net
// Email       : support@wengold.net
//
// Prismy.No | Date       | Modified by. | Description
// -------------------------------------------------------------------
// 00001       2019/05/22   yangping       New version
// -------------------------------------------------------------------

package builder

import (
	"fmt"
	"strings"

	"github.com/wengoldx/xcore/logger"
	pd "github.com/wengoldx/xcore/mvc/provider"
	"github.com/wengoldx/xcore/utils"
)

// Build a query string for sql update.
//
//	UPDATE table
//		SET v1=?, v2=?, v3=?...
//		WHERE wherer AND field IN (v1,v2...) AND field2 LIKE '%%filter%%'
//
// See QuerierImpl, InserterImpl, DeleterImpl.
type UpdaterImpl struct {
	BuilderImpl

	values pd.KValues // Target fields and values to update.
	wheres pd.Wheres  // Where conditions and args values.
	sep    string     // Where conditions connector, one of 'AND', 'OR', ' ', default ''.
	ins    string     // Where in conditions.
	like   string     // Like conditions string.
}

var _ pd.SQLBuilder = (*UpdaterImpl)(nil)
var _ pd.UpdateBuilder = (*UpdaterImpl)(nil)

// Create a UpdateBuilder instance to build a query string.
func NewUpdate(table string) pd.UpdateBuilder {
	return &UpdaterImpl{BuilderImpl: NewBuilder(table)}
}

/* ------------------------------------------------------------------- */
/* SQL Action Utils By Using master Provider                           */
/* ------------------------------------------------------------------- */

func (b *UpdaterImpl) Exec() error   { return b.provider.Exec(b) }   // Update target record without check.
func (b *UpdaterImpl) Update() error { return b.provider.Update(b) } // Update target record and check changes counts.

/* ------------------------------------------------------------------- */
/* For UpdateBuilder interface                                         */
/* ------------------------------------------------------------------- */

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
func (b *UpdaterImpl) Values(row pd.KValues) pd.UpdateBuilder {
	b.values = row
	return b
}

// Specify the where conditions and args for query.
//
//	where = pd.Wheres{
//		"acc=?":"123", "age>=?":20, "role<>?":"admin",
//	}
//	// => WHERE acc=? AND age>=? AND role<>?
//	// => args ("123", 20, "admin")
func (b *UpdaterImpl) Wheres(where pd.Wheres) pd.UpdateBuilder {
	b.wheres = where
	return b
}

// Specify the where in condition with field and args for query.
//
//	builder.WhereIn("id", []any{1, 2}) // => WHERE id IN (1, 2)
func (b *UpdaterImpl) WhereIn(field string, args []any) pd.UpdateBuilder {
	b.ins = b.FormatWhereIn(field, args)
	return b
}

// Specify the where in condition with field and args for query.
func (b *UpdaterImpl) WhereSep(sep string) pd.UpdateBuilder {
	switch s := strings.ToUpper(sep); s {
	case "AND", "OR", " " /* for none where connector */ :
		b.sep = s
	}
	return b
}

// Specify the like condition for query.
//
//	builder.Like("acc", "zhang")           // => acc LIKE '%%zhang%%'
//	builder.Like("acc", "zhang", "perfix") // => acc LIKE 'zhang%%'
//	builder.Like("acc", "zhang", "suffix") // => acc LIKE '%%zhang'
func (b *UpdaterImpl) Like(field, filter string, pattern ...string) pd.UpdateBuilder {
	b.like = b.FormatLike(field, filter, pattern...)
	return b
}

// Reset builder datas for next prepare and build.
func (b *UpdaterImpl) Reset() pd.UpdateBuilder {
	clear(b.values)
	clear(b.wheres)
	b.sep, b.ins, b.like = "", "", ""
	return b
}

/* ------------------------------------------------------------------- */
/* For SQLBuilder interface                                            */
/* ------------------------------------------------------------------- */

// Build the update action sql string and args for provider to update datas.
//
//	UPDATE table
//		SET v1=?, v2=?, v3=?...
//		WHERE wherer AND field IN (v1,v2...) AND field2 LIKE '%%filter%%'
func (b *UpdaterImpl) Build(debug ...bool) (string, []any) {
	sep := utils.Condition(b.sep == "", "AND", b.sep)

	tags, args := b.FormatSets(b.values)                      // SET v1=?,v2=?...
	where, wvs := b.BuildWheres(b.wheres, b.ins, b.like, sep) // WHERE wheres AND field IN (v1,v2...) AND field2 LIKE '%%filter%%'
	args = append(args, wvs...)

	query := "UPDATE %s SET %s %s"
	query = fmt.Sprintf(query, b.table, tags, where)
	query = strings.TrimRight(query, " ")
	if utils.Variable(debug, false) {
		logger.D("[UPDATE] SQL:", query, "|", args)
	}
	return query, args
}
