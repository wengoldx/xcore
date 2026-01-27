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
//		WHERE wherer AND field IN (v1,v2...) AND field2 LIKE '%%filter%%'
//		LIMIT limit.
//
// See QuerierImpl, InserterImpl, UpdaterImpl.
type DeleterImpl struct {
	BuilderImpl

	wheres pd.Wheres // Where conditions and args values.
	sep    string    // Where conditions connector, one of 'AND', 'OR', ' ', default ''.
	ins    string    // Where in conditions.
	like   string    // Like conditions string.
	limit  int       // Limit number.
}

var _ pd.SQLBuilder = (*DeleterImpl)(nil)
var _ pd.DeleteBuilder = (*DeleterImpl)(nil)

// Create a DeleteBuilder instance to build a query string.
func NewDelete(table string) pd.DeleteBuilder {
	return &DeleterImpl{BuilderImpl: NewBuilder(table)}
}

/* ------------------------------------------------------------------- */
/* SQL Action Utils By Using master Provider                           */
/* ------------------------------------------------------------------- */

func (b *DeleterImpl) Exec() error   { return b.provider.Exec(b) }   // Delete record without check.
func (b *DeleterImpl) Delete() error { return b.provider.Delete(b) } // Delete record and check deleted counts.

/* ------------------------------------------------------------------- */
/* For DeleteBuilder interface                                         */
/* ------------------------------------------------------------------- */

// Specify the where conditions and args for query.
//
//	where = pd.Wheres{
//		"acc=?":"123", "age>=?":20, "role<>?":"admin",
//	}
//	// => WHERE acc=? AND age>=? AND role<>?
//	// => args ("123", 20, "admin")
func (b *DeleterImpl) Wheres(where pd.Wheres) pd.DeleteBuilder {
	b.wheres = where
	return b
}

// Specify the where in condition with field and args for query.
//
//	builder.WhereIn("id", []any{1, 2}) // => WHERE id IN (1, 2)
func (b *DeleterImpl) WhereIn(field string, args []any) pd.DeleteBuilder {
	b.ins = b.FormatWhereIn(field, args)
	return b
}

// Specify the where in condition with field and args for query.
func (b *DeleterImpl) WhereSep(sep string) pd.DeleteBuilder {
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
func (b *DeleterImpl) Like(field, filter string, pattern ...string) pd.DeleteBuilder {
	b.like = b.FormatLike(field, filter, pattern...)
	return b
}

// Specify the limit result for query.
//
//	builder.Limit(20) // => LIMIT 20
func (b *DeleterImpl) Limit(limit int) pd.DeleteBuilder {
	b.limit = limit
	return b
}

// Reset builder datas for next prepare and build.
func (b *DeleterImpl) Reset() pd.DeleteBuilder {
	clear(b.wheres)
	b.sep, b.ins, b.like = "", "", ""
	b.limit = 0
	return b
}

/* ------------------------------------------------------------------- */
/* For SQLBuilder interface                                            */
/* ------------------------------------------------------------------- */

// Build the delete action sql string and args for provider to delete datas.
//
//	DELETE FROM table
//		WHERE wherer AND field IN (v1,v2...) AND field2 LIKE '%%filter%%'
//		LIMIT limit.
func (b *DeleterImpl) Build(debug ...bool) (string, []any) {
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
