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
//		WHERE wherers AND field IN (v1,v2...) AND field2 LIKE '%%filter%%'
type UpdateBuilder struct {
	BaseBuilder
	values pd.KValues // Target fields and values to update.
	wheres pd.Wheres  // Where conditions and args values.
	sep    string     // Where conditions connector, one of 'AND', 'OR', ' ', default ''.
	ins    string     // Where in conditions.
	like   string     // Like conditions string.
}

var _ pd.Builder = (*UpdateBuilder)(nil)

// Create a UpdateBuilder instance to build a query string.
func NewUpdate(table string, provider ...pd.ProviderUtils) *UpdateBuilder {
	return &UpdateBuilder{BaseBuilder: *NewBuilder(table, provider...)}
}

/* ------------------------------------------------------------------- */
/* For Provider Update Utils                                           */
/* ------------------------------------------------------------------- */

// Update record without check.
//
//	h.Updater().Values(pd.KValues{"role": admin}).Wheres(pd.Wheres{"uid=?": uid}).Exec()
//	// not check updated row count.
func (b *UpdateBuilder) Exec() error { return b.provider.Exec(b) }

// Update record and check changed counts.
//
//	h.Updater().Values(pd.KValues{"role": admin}).Wheres(pd.Wheres{"uid=?": uid}).Exec()
//	// check updated row count.
func (b *UpdateBuilder) Update() error { return b.provider.Update(b) }

/* ------------------------------------------------------------------- */
/* For SQL String Build Utils                                          */
/* ------------------------------------------------------------------- */

// Specify the values of row to update.
//
//	values := KValues{
//		"":       123456,   // Filter out empty field
//		"Age":    16,
//		"Male":   true,
//		"Name":   "ZhangSan",
//		"Height": 176.8,
//		"Secure": nil,      // Set value as NULL
//	}
//	// => SET Age=?, Male=?, Name=?, Height=?, Secure=NULL
//	// => values: []any{16, true, "ZhangSan", 176.8, nil}
func (b *UpdateBuilder) Values(row pd.KValues) *UpdateBuilder {
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
func (b *UpdateBuilder) Wheres(where pd.Wheres) *UpdateBuilder {
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

// Specify the where in condition with field and args for query.
//
//	builder.In(pd.NewIn("id", []int{1,2})) // => WHERE id IN (1, 2)
func (b *UpdateBuilder) In(in *pd.In) *UpdateBuilder {
	b.ins = b.FormatWhereIn(in.Get())
	return b
}

// Specify the like condition for query.
//
//	builder.Like("acc", "zhang")           // => acc LIKE '%%zhang%%'
//	builder.Like("acc", "zhang", "perfix") // => acc LIKE 'zhang%%'
//	builder.Like("acc", "zhang", "suffix") // => acc LIKE '%%zhang'
func (b *UpdateBuilder) Like(field, filter string, pattern ...string) *UpdateBuilder {
	b.like = b.FormatLike(field, filter, pattern...)
	return b
}

// Reset builder datas for next prepare and build.
func (b *UpdateBuilder) Reset() *UpdateBuilder {
	clear(b.values)
	clear(b.wheres)
	b.sep, b.ins, b.like = "", "", ""
	return b
}

/* ------------------------------------------------------------------- */
/* For SQL Builder interface                                           */
/* ------------------------------------------------------------------- */

// Build the update action sql string and args for provider to update datas.
//
//	UPDATE table
//		SET v1=?, v2=?, v3=?...
//		WHERE wherer AND field IN (v1,v2...) AND field2 LIKE '%%filter%%'
func (b *UpdateBuilder) Build(debug ...bool) (string, []any) {
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
