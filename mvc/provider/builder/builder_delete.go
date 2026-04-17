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

// Build a query string for sql delete.
//
//	DELETE FROM table
//		WHERE wheres AND field IN (v1,v2...) AND field2 LIKE '%%filter%%'
//		LIMIT limit.
type DeleteBuilder struct {
	BaseBuilder
	wheres pd.Wheres // Where conditions and args values.
	sep    string    // Where conditions connector, one of 'AND', 'OR', ' ', default ''.
	ins    string    // Where in conditions.
	like   string    // Like conditions string.
	limit  int       // Limit number.
}

var _ pd.Builder = (*DeleteBuilder)(nil)

// Create a DeleteBuilder instance to build a query string.
func NewDelete(table string, provider ...pd.ProviderUtils) *DeleteBuilder {
	return &DeleteBuilder{BaseBuilder: *NewBuilder(table, provider...)}
}

/* ------------------------------------------------------------------- */
/* For Provider Delete Utils                                           */
/* ------------------------------------------------------------------- */

// Delete record without check.
//
//	h.Deleter().Wheres(pd.Wheres{"uid=?": uid}).Exec()
//	// not check delete row count.
func (b *DeleteBuilder) Exec() error { return b.provider.Exec(b) }

// Delete record and check deleted counts.
//
//	h.Deleter().Wheres(pd.Wheres{"uid=?": uid}).Delete()
//	// check delete row count.
func (b *DeleteBuilder) Delete() error { return b.provider.Delete(b) }

/* ------------------------------------------------------------------- */
/* For SQL String Build Utils                                          */
/* ------------------------------------------------------------------- */

// Specify the where conditions and args for query.
//
//	where = pd.Wheres{
//		"acc=?":"123", "age>=?":20, "role<>?":"admin",
//	}
//	// => WHERE acc=? AND age>=? AND role<>?
//	// => args ("123", 20, "admin")
func (b *DeleteBuilder) Wheres(where pd.Wheres) *DeleteBuilder {
	b.wheres = where
	return b
}

// Specify the where in condition with field and args for query.
//
//	builder.WhereIn("id", []any{1, 2}) // => WHERE id IN (1, 2)
func (b *DeleteBuilder) WhereIn(field string, args []any) *DeleteBuilder {
	b.ins = b.FormatWhereIn(field, args)
	return b
}

// Specify the where in condition with field and args for query.
func (b *DeleteBuilder) WhereSep(sep string) *DeleteBuilder {
	switch s := strings.ToUpper(sep); s {
	case "AND", "OR", " " /* for none where connector */ :
		b.sep = s
	}
	return b
}

// Specify the where in condition with field and args for query.
//
//	builder.In(pd.NewIn("id", []int{1,2})) // => WHERE id IN (1, 2)
func (b *DeleteBuilder) In(in *pd.In) *DeleteBuilder {
	b.ins = b.FormatWhereIn(in.Get())
	return b
}

// Specify the like condition for query.
//
//	builder.Like("acc", "zhang")           // => acc LIKE '%%zhang%%'
//	builder.Like("acc", "zhang", "perfix") // => acc LIKE 'zhang%%'
//	builder.Like("acc", "zhang", "suffix") // => acc LIKE '%%zhang'
func (b *DeleteBuilder) Like(field, filter string, pattern ...string) *DeleteBuilder {
	b.like = b.FormatLike(field, filter, pattern...)
	return b
}

// Specify the limit result for query.
//
//	builder.Limit(20) // => LIMIT 20
func (b *DeleteBuilder) Limit(limit int) *DeleteBuilder {
	b.limit = limit
	return b
}

// Reset builder datas for next prepare and build.
func (b *DeleteBuilder) Reset() *DeleteBuilder {
	clear(b.wheres)
	b.sep, b.ins, b.like = "", "", ""
	b.limit = 0
	return b
}

/* ------------------------------------------------------------------- */
/* For SQL Builder interface                                           */
/* ------------------------------------------------------------------- */

// Build the delete action sql string and args for provider to delete datas.
//
//	DELETE FROM table
//		WHERE wherer AND field IN (v1,v2...) AND field2 LIKE '%%filter%%'
//		LIMIT limit.
func (b *DeleteBuilder) Build(debug ...bool) (string, []any) {
	sep := utils.Condition(b.sep == "", "AND", b.sep)
	where, args := b.BuildWheres(b.wheres, b.ins, b.like, sep) // WHERE wheres AND field IN (v1,v2...) AND field2 LIKE '%%filter%%'
	limit := b.FormatLimit(b.limit)                            // LIMIT n

	query := "DELETE FROM %s %s %s"
	query = fmt.Sprintf(query, b.table, where, limit)
	query = strings.TrimRight(query, " ")

	if utils.Variable(debug, false) {
		logger.D("[DELETE] SQL:", query, "|", args)
	}
	return query, args
}
