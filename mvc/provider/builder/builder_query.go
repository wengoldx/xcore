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

// Build a query string for sql query.
//
//	SELECT tags FROM table
//		WHERE wherer AND field IN (v1,v2...) AND field2 LIKE '%%filter%%'
//		ORDER BY order DESC
//		LIMIT limit.
//
// See InserterImpl, UpdaterImpl, DeleterImpl.
type QuerierImpl struct {
	BuilderImpl

	joins  pd.Joins  // Table-Alias for multi-table joins.
	tags   []string  // Target fields for output values.
	outs   []any     // The params output query results, only for single query.
	wheres pd.Wheres // Where conditions and args values.
	sep    string    // Where conditions connector, one of 'AND', 'OR', ' ', default ''.
	ins    string    // Where in conditions.
	like   string    // Like conditions string.
	order  string    // Keyword for order by condition.
	limit  int       // Limit number.
}

var _ pd.SQLBuilder = (*QuerierImpl)(nil)
var _ pd.QueryBuilder = (*QuerierImpl)(nil)

// Create a QueryBuilder instance to build a query string.
func NewQuery(table string) pd.QueryBuilder {
	return &QuerierImpl{BuilderImpl: NewBuilder(table)}
}

/* ------------------------------------------------------------------- */
/* SQL Action Utils By Using master Provider                           */
/* ------------------------------------------------------------------- */

func (b *QuerierImpl) Has() (bool, error)               { return b.provider.Has(b) }         // Check whether has the target record.
func (b *QuerierImpl) None() (bool, error)              { return b.provider.None(b) }        // Check whether unexist the target record.
func (b *QuerierImpl) Count() (int, error)              { return b.provider.Count(b) }       // Count the mathed query condition records.
func (b *QuerierImpl) OneScan(cb pd.ScanCallback) error { return b.provider.OneScan(b, cb) } // Query the top one record with scan callback.
func (b *QuerierImpl) Query(cb pd.ScanCallback) error   { return b.provider.Query(b, cb) }   // Query the all matched condition records.
func (b *QuerierImpl) Array(cr pd.SQLCreator) error     { return b.provider.Array(b, cr) }   // Query the all records with the SQLCreator utils.

// Query the top one record and return the results without scaner
// callback, it canbe set the finally done callback called when
// result success read.
func (b *QuerierImpl) OneDone(done ...pd.DoneCallback) error {
	if len(done) > 0 && done[0] != nil {
		return b.provider.OneDone(b, done[0], b.outs...)
	} else {
		return b.provider.OneOuts(b, b.outs...)
	}
}

/* ------------------------------------------------------------------- */
/* For QueryBuilder interface                                          */
/* ------------------------------------------------------------------- */

// Specify the table-alias joins for query.
func (b *QuerierImpl) Joins(tables pd.Joins) pd.QueryBuilder {
	b.joins = tables
	return b
}

// Specify the target output fields name for query.
func (b *QuerierImpl) Tags(tag ...string) pd.QueryBuilder {
	b.tags = tag
	return b
}

// Specify the target output params for single query, the
// outs length must same as Tags length.
func (b *QuerierImpl) Outs(outs ...any) pd.QueryBuilder {
	b.outs = outs
	return b
}

// Specify the target column and output param for single query.
func (b *QuerierImpl) TagOut(tag string, out any) pd.QueryBuilder {
	return b.Tags(tag).Outs(out)
}

// Specify the target columns and struct fields for single query.
//
//	Set BuilderImpl.ParseOut() get more info.
func (b *QuerierImpl) Parse(out any) pd.QueryBuilder {
	b.tags, b.outs = b.ParseOut(out)
	return b
}

// Specify the where conditions and args for query.
//
//	where = pd.Wheres{
//		"acc=?":"123", "age>=?":20, "role<>?":"admin",
//	}
//	// => WHERE acc=? AND age>=? AND role<>?
//	// => args ("123", 20, "admin")
func (b *QuerierImpl) Wheres(where pd.Wheres) pd.QueryBuilder {
	b.wheres = where
	return b
}

// Specify the where in condition with field and args for query.
//
//	builder.WhereIn("id", []any{1, 2}) // => WHERE id IN (1, 2)
func (b *QuerierImpl) WhereIn(field string, args []any) pd.QueryBuilder {
	b.ins = b.FormatWhereIn(field, args)
	return b
}

// Specify the where in condition with field and args for query.
func (b *QuerierImpl) WhereSep(sep string) pd.QueryBuilder {
	switch s := strings.ToUpper(sep); s {
	case "AND", "OR", " " /* for none where connector */ :
		b.sep = s
	}
	return b
}

// Specify the order by condition for query.
//
//	builder.OrderBy("id")          // => ORDER BY id DESC
//	builder.OrderBy("slug", false) // => ORDER BY slug ASC
func (b *QuerierImpl) OrderBy(field string, desc ...bool) pd.QueryBuilder {
	b.order = b.FormatOrder(field, desc...)
	return b
}

// Specify the like condition for query.
//
//	builder.Like("acc", "zhang")           // => acc LIKE '%%zhang%%'
//	builder.Like("acc", "zhang", "perfix") // => acc LIKE 'zhang%%'
//	builder.Like("acc", "zhang", "suffix") // => acc LIKE '%%zhang'
func (b *QuerierImpl) Like(field, filter string, pattern ...string) pd.QueryBuilder {
	b.like = b.FormatLike(field, filter, pattern...)
	return b
}

// Specify the limit result for query.
//
//	builder.Limit(20) // => LIMIT 20
func (b *QuerierImpl) Limit(limit int) pd.QueryBuilder {
	b.limit = limit
	return b
}

// Reset builder datas for next prepare and build.
func (b *QuerierImpl) Reset() pd.QueryBuilder {
	clear(b.tags)
	clear(b.wheres)
	clear(b.outs)
	b.sep, b.ins, b.like, b.order = "", "", "", ""
	b.limit = 0
	return b
}

/* ------------------------------------------------------------------- */
/* For SQLBuilder interface                                            */
/* ------------------------------------------------------------------- */

// Build the query action sql string and args for provider to query datas.
//
//	SELECT tags FROM [table | table1 AS a, table2 AS b, ...]
//		WHERE wherer AND field IN (v1,v2...) AND field2 LIKE '%%filter%%'
//		ORDER BY order DESC
//		LIMIT limit.
func (b *QuerierImpl) Build(debug ...bool) (string, []any) {
	sep := utils.Condition(b.sep == "", "AND", b.sep)

	tags := strings.Join(b.tags, ",")                          // out1,out2,out3...
	where, args := b.BuildWheres(b.wheres, b.ins, b.like, sep) // WHERE wheres AND field IN (v1,v2...) AND field2 LIKE '%%filter%%'
	limit := b.FormatLimit(b.limit)                            // LIMIT n

	joins := b.FormatJoins(b.joins)                       // table1 AS a, table2 AS b
	table := utils.Condition(joins != "", joins, b.table) // priority use of joined tables, or use b.table

	query := "SELECT %s FROM %s %s %s %s"
	query = fmt.Sprintf(query, tags, table, where, b.order, limit)
	query = strings.TrimRight(query, " ")
	if utils.Variable(debug, false) {
		logger.D("[QUERY] SQL:", query, "|", args)
	}
	return query, args
}
