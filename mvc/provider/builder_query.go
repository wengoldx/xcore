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

	"github.com/wengoldx/xcore/logger"
	"github.com/wengoldx/xcore/utils"
)

// Build a query string for sql query.
//
//	SELECT tags FROM table
//		WHERE wherer AND field IN (v1,v2...) AND field2 LIKE '%%filter%%'
//		ORDER BY order DESC
//		LIMIT limit.
//
// See InsertBuilder, UpdateBuilder, DeleteBuilder.
type QueryBuilder struct {
	BaseBuilder

	table  string   // Table name for query
	joins  Joins    // Table-Alias for multi-table joins.
	tags   []string // Target fields for output values.
	outs   []any    // The params output query results, only for single query.
	wheres Wheres   // Where conditions and args values.
	sep    string   // Where conditions connector, one of 'AND', 'OR', ' ', default ''.
	ins    string   // Where in conditions.
	like   string   // Like conditions string.
	order  string   // Keyword for order by condition.
	limit  int      // Limit number.

	// Model creater for records query to create a new item object,
	// and get target query columns bind out values, then append the
	// query result to given array.
	ItemCreator SQLItemCreator
}

var _ SQLBuilder = (*QueryBuilder)(nil)

// Create a QueryBuilder instance to build a query string.
func NewQuery(table string) *QueryBuilder {
	return &QueryBuilder{table: table}
}

/* ------------------------------------------------------------------- */
/* SQL Action Utils By Using master Provider                           */
/* ------------------------------------------------------------------- */

func (b *QueryBuilder) Has() (bool, error)          { return b.master.Has(b) }        // Check whether has the target record.
func (b *QueryBuilder) None() (bool, error)         { return b.master.None(b) }       // Check whether unexist the target record.
func (b *QueryBuilder) Count() (int, error)         { return b.master.Count(b) }      // Count the mathed query condition records.
func (b *QueryBuilder) One(cb ScanCallback) error   { return b.master.One(b, cb) }    // Query the top one record.
func (b *QueryBuilder) Query(cb ScanCallback) error { return b.master.Query(b, cb) }  // Query the all matched condition records.
func (b *QueryBuilder) Querys(cb AddCallback) error { return b.master.Querys(b, cb) } // Query the all records with the SQLItemCreator utils.

// Query the top one record and return the results without scaner
// callback, it canbe set the finally done callback called when
// result success read.
//
// # NOTICE
//
// This method also used to query one record of SQLModel data, as follow:
//
//	// Define MyAcc and implement pd.SQLModel interface.
//	type MyAcc struct { Name string }
//	func (c *MyAcc) MapValues() { return map[string]any{"name": &c.Name} }
//
//	acc := MyAcc{}
//	h.Querier().Model(acc).Wheres(pd.Wheres{"role=?": "admin"}).OneOuts()
//	// the query result filled into acc object.
//
//	See `SQLModel` interface for model data define.
func (b *QueryBuilder) OneDone(done ...DoneCallback) error {
	if len(done) > 0 && done[0] != nil {
		return b.master.OneDone(b, done[0], b.outs...)
	} else {
		return b.master.OneOuts(b, b.outs...)
	}
}

/* ------------------------------------------------------------------- */
/* SQL Action Builder Methonds                                         */
/* ------------------------------------------------------------------- */

// Specify master table provider.
func (b *QueryBuilder) Master(master *TableProvider) *QueryBuilder {
	b.master = master
	return b
}

// Specify the target table for query.
func (b *QueryBuilder) Table(table string) *QueryBuilder {
	b.table = table
	return b
}

// Specify the table-alias joins for query.
func (b *QueryBuilder) Joins(tables Joins) *QueryBuilder {
	b.joins = tables
	return b
}

// Specify the target output fields name for query.
func (b *QueryBuilder) Tags(tag ...string) *QueryBuilder {
	b.tags = tag
	return b
}

// Specify the target output params for single query, the
// outs length must same as Tags length.
func (b *QueryBuilder) Outs(outs ...any) *QueryBuilder {
	b.outs = outs
	return b
}

// Specify the target column and output param for single query.
func (b *QueryBuilder) TagOut(tag string, out any) *QueryBuilder {
	b.tags, b.outs = []string{tag}, []any{out}
	return b
}

// Specify the target output fields and params for single query.
//
// # WARNING
//	- The model columns must not empty and out values must pointer like &myvalue.
//	- No-need call Tags() and Outs() again when called this method.
//
// See SQLModel interface for out model data define.
func (b *QueryBuilder) Model(model SQLModel) *QueryBuilder {
	tags, outs := []string{}, []any{}

	tagouts := model.MapValues()
	for tag, out := range tagouts {
		if tag != "" && out != nil {
			tags = append(tags, tag)
			outs = append(outs, out)
		}
	}
	b.tags, b.outs = tags, outs
	return b
}

// Specify the target output fields and params for array query.
//
// # WARNING
//	- The creator.GetTags() returned columns names must not empty.
//	- The creator.GetOuts() returned outs must a point value like &myvalue.
//	- The returned tags and outs array lenght must same.
//	- No-need call Tags() and Outs() again when called this method.
//
// See SQLItemCreator interface for out model data define.
func (b *QueryBuilder) Models(creater SQLItemCreator) *QueryBuilder {
	b.tags = creater.GetTags()
	b.ItemCreator = creater
	return b
}

// Specify the where conditions and args for query.
//
//	where = provider.Wheres{
//		"acc=?":"123", "age>=?":20, "role<>?":"admin",
//	}
//	// => WHERE acc=? AND age>=? AND role<>?
//	// => args ("123", 20, "admin")
func (b *QueryBuilder) Wheres(where Wheres) *QueryBuilder {
	b.wheres = where
	return b
}

// Specify the where in condition with field and args for query.
//
//	builder.WhereIn("id", []any{1, 2}) // => WHERE id IN (1, 2)
func (b *QueryBuilder) WhereIn(field string, args []any) *QueryBuilder {
	b.ins = b.FormatWhereIn(field, args)
	return b
}

// Specify the where in condition with field and args for query.
func (b *QueryBuilder) WhereSep(sep string) *QueryBuilder {
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
func (b *QueryBuilder) OrderBy(field string, desc ...bool) *QueryBuilder {
	b.order = b.FormatOrder(field, desc...)
	return b
}

// Specify the like condition for query.
//
//	builder.Like("acc", "zhang")           // => acc LIKE '%%zhang%%'
//	builder.Like("acc", "zhang", "perfix") // => acc LIKE 'zhang%%'
//	builder.Like("acc", "zhang", "suffix") // => acc LIKE '%%zhang'
func (b *QueryBuilder) Like(field, filter string, pattern ...string) *QueryBuilder {
	b.like = b.FormatLike(field, filter, pattern...)
	return b
}

// Specify the limit result for query.
//
//	builder.Limit(20) // => LIMIT 20
func (b *QueryBuilder) Limit(limit int) *QueryBuilder {
	b.limit = limit
	return b
}

// Build the query action sql string and args for provider to query datas.
//
//	SELECT tags FROM [table | table1 AS a, table2 AS b, ...]
//		WHERE wherer AND field IN (v1,v2...) AND field2 LIKE '%%filter%%'
//		ORDER BY order DESC
//		LIMIT limit.
func (b *QueryBuilder) Build(debug ...bool) (string, []any) {
	sep := utils.Condition(b.sep == "", "AND", b.sep)

	tags := strings.Join(b.tags, ",")                          // out1,out2,out3...
	where, args := b.BuildWheres(b.wheres, b.ins, b.like, sep) // WHERE wheres AND field IN (v1,v2...) AND field2 LIKE '%%filter%%'
	limit := b.FormatLimit(b.limit)                            // LIMIT n

	joins := b.FormatJoins(b.joins)                       // table1 AS a, table2 AS b
	table := utils.Condition(joins != "", joins, b.table) // priority use of joined tables, or use b.table

	query := "SELECT %s FROM %s %s %s %s"
	query = fmt.Sprintf(query, tags, table, where, b.order, limit)
	query = strings.TrimSuffix(query, " ")
	if utils.Variable(debug, false) {
		logger.D("[QUERY] SQL:", query, "|", args)
	}
	return query, args
}

// Reset builder datas for next prepare and build.
func (b *QueryBuilder) Reset() *QueryBuilder {
	clear(b.tags)
	clear(b.wheres)
	clear(b.outs)
	b.sep, b.ins, b.like, b.order = "", "", "", ""
	b.limit = 0
	return b
}
