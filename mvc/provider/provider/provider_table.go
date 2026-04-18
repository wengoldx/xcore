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
	"database/sql"

	"github.com/astaxie/beego"
	"github.com/wengoldx/xcore/invar"
	pd "github.com/wengoldx/xcore/mvc/provider"
	"github.com/wengoldx/xcore/mvc/provider/builder"
	"github.com/wengoldx/xcore/utils"
)

// Table provider for using builder to build query string and args for
// database datas access.
//
// Usage: Define the custom provider as follow code.
//
//	// define the custom provider.
//	type SampleTable struct { *provider.TableProvider }
//	s := &SampleTable{*mysql.NewTable("sample")}
//	mysql.SetClient(s)
//
//	// or create directory if client exist.
//	s := provider.NewTableProvider(client,
//		provider.WithTable("sample"), provider.WithDebug(true))
//
// Use mysql.NewTable(), mysql.NewTable() sqlite.NewTable() of mvc inner packages
// to create TableProvider with connected mysql, mssql, sqlite database client.
type TableProvider struct {
	BaseProvider
	table string // Table name.
	debug bool   // Debug flag for print SQL actions, default false.
}

var _ pd.Provider = (*TableProvider)(nil)
var _ pd.ProviderUtils = (*TableProvider)(nil)

// Create a TableProvider with given database client.
//
// # WARNING:
//
// This method call by 'sqlite3', 'mysql', 'mssql' module called to create
// target table and bind with connected database client instance.
func NewTableProvider(client pd.DBClient, opts ...Option) *TableProvider {
	logsql := beego.AppConfig.String("logger::logsql") == "on"
	tp := &TableProvider{debug: logsql}
	tp.BaseProvider = *NewBaseProvider(client)
	for _, optFunc := range opts {
		optFunc(tp)
	}
	return tp
}

// The setter for set TableProvider options.
type Option func(provider *TableProvider)

// Specify the table name.
func WithTable(table string) Option {
	return func(provider *TableProvider) { provider.table = table }
}

/* ------------------------------------------------------------------- */
/* Create and Return Builder Instance FOR QUID Actions                 */
/* ------------------------------------------------------------------- */

// Set current table provider debug flag, true is on, false not.
func (p *TableProvider) Debug(onoff bool) *TableProvider {
	p.debug = onoff
	return p
}

// Create a query builder to query table records.
//
//	SELECT tags FROM table
//		WHERE wheres AND field IN (v1,v2...) AND field2 LIKE '%%filter%%'
//		ORDER BY order DESC
//		LIMIT limit.
func (p *TableProvider) Querier(t ...string) *builder.QueryBuilder {
	return builder.NewQuery(utils.Variable(t, p.table), p)
}

// Create a insert builder to insert records to table.
//
//	`MySQL & MSSQL`: INSERT table (tags) VALUES (?, ?, ?)...
//	`SQLITE`       : INSERT INTO table (tags) VALUES (?, ?, ?)...
func (p *TableProvider) Inserter(t ...string) *builder.InsertBuilder {
	return builder.NewInsert(utils.Variable(t, p.table), p)
}

// Create a update builder to update table records.
//
//	UPDATE table
//		SET v1=?, v2=?, v3=?...
//		WHERE wherers AND field IN (v1,v2...) AND field2 LIKE '%%filter%%'
func (p *TableProvider) Updater(t ...string) *builder.UpdateBuilder {
	return builder.NewUpdate(utils.Variable(t, p.table), p)
}

// Create a delete builder to delete table records.
//
//	DELETE FROM table
//		WHERE wheres AND field IN (v1,v2...) AND field2 LIKE '%%filter%%'
//		LIMIT limit.
func (p *TableProvider) Deleter(t ...string) *builder.DeleteBuilder {
	return builder.NewDelete(utils.Variable(t, p.table), p)
}

/* ------------------------------------------------------------------- */
/* Using Builder To Construct Query String For Database Access         */
/* ------------------------------------------------------------------- */

// Check the target record whether exist by the given QueryBuilder to
// build query string, it no-need set any tags.
//
// # USAGE
//
//	h.Querier().Wheres(pd.Wheres{"account=?": acc}).Has()
//
// Use None() method to check whether unexist.
func (p *TableProvider) Has(b pd.Builder) (bool, error) {
	if qb, ok := b.(*builder.QueryBuilder); ok {
		query, args := qb.Tags("*").Build(p.debug)
		return p.BaseProvider.Has(query, args...)
	}
	return false, invar.ErrBadSQLBuilder
}

// Check the target record whether unexist by the given QueryBuilder to
// build query string, it no-need set any tags.
//
// # USAGE
//
//	h.Querier().Wheres(pd.Wheres{"account=?": acc}).None()
//
// Use Has() method to check has result.
func (p *TableProvider) None(builder pd.Builder) (bool, error) {
	has, err := p.Has(builder)
	return !has, err
}

// Count records by the given builder to build a query string, it will
// return 0 when notfound anyone, it no-need set any tags.
//
// # USAGE
//
//	h.Querier().Wheres(pd.Wheres{"role=?": "admin"}).Count()
//
// Use BaseProvider.Count() method to direct execute query string.
func (p *TableProvider) Count(b pd.Builder) (int, error) {
	if qb, ok := b.(*builder.QueryBuilder); ok {
		query, args := qb.Tags("COUNT(*)").Build(p.debug)
		return p.BaseProvider.Count(query, args...)
	}
	return 0, invar.ErrBadSQLBuilder
}

// Query one record by given builder builded query string, and read datas
// from scan callback.
//
// # NOTICE:
//	- Use BaseProvider.One() method to direct execute query string.
func (p *TableProvider) OneScan(b pd.Builder, cb pd.ScanCallback) error {
	if qb, ok := b.(*builder.QueryBuilder); ok {
		query, args := qb.Build(p.debug)
		return p.BaseProvider.One(query, cb, args...)
	}
	return invar.ErrBadSQLBuilder
}

// Query the top one record and return the results without scaner
// callback, it canbe set the finally done callback called when
// result success read.
func (p *TableProvider) OneDone(b pd.Builder, done ...pd.DoneCallback) error {
	if qb, ok := b.(*builder.QueryBuilder); ok {
		query, args := qb.Build(p.debug)
		if cb := utils.Variable(done, nil); cb != nil {
			return p.BaseProvider.OneDone(query, qb.GetOuts(), cb, args...)
		}
		return p.BaseProvider.OneDone(query, qb.GetOuts(), nil, args...)
	}
	return invar.ErrBadSQLBuilder
}

// Query records by given builder builded query string, and read datas
// from scan callback.
//
// Use BaseProvider.Query() method to direct execute query string.
func (p *TableProvider) Query(b pd.Builder, cb pd.ScanCallback) error {
	if qb, ok := b.(*builder.QueryBuilder); ok {
		query, args := qb.Build(p.debug)
		return p.BaseProvider.Query(query, cb, args...)
	}
	return invar.ErrBadSQLBuilder
}

// Query records by given builder builded query string, and read datas
// from ItemCreator instance.
//
//	type MyAcc struct { Name string }
//
//	accs := []*MyAcc{}
//	creator := pd.NewCreator(func(iv *MyAcc) []any {
//		return []any{&iv.Name}
//	}, /* func(iv *MyAcc) {} */) // or append parser function.
//	h.Querier().Tags("name").Wheres(pd.Wheres{"role=?": "admin"}).Array(creator)
func (p *TableProvider) Array(b pd.Builder, creator pd.Creator) error {
	return p.Query(b, func(rows *sql.Rows) error {
		item, outs := creator.CreateItem() // item is *T type, outs all & pointers!
		if err := rows.Scan(outs...); err != nil {
			return err
		}
		return creator.AppendItem(item)
	})
}

// Query single column values by given builder builded query string,
// and read datas from ItemScaner instance.
//
//	names := []string{}
//	scaner := pd.NewScaner(&names/* , func(iv *string) {} */)
//	h.Querier().Tags("name").Wheres(pd.Wheres{"role=?": "admin"}).Column(scaner)
func (p *TableProvider) Column(b pd.Builder, scaner pd.Scaner) error {
	return p.Query(b, func(rows *sql.Rows) error {
		out := scaner.CreateItem() // out is *T type!
		if err := rows.Scan(out); err != nil {
			return err
		}
		return scaner.AppendItem(out)
	})
}

// Execute the query string builded from given QueryBuilder, InsertBuilder,
// UpdateBuilder or DeleteBuilder, it not check the affected row counts.
//
// Use BaseProvider.Exec() method to direct execute query string.
func (p *TableProvider) Exec(b pd.Builder) error {
	query, args := b.Build(p.debug)
	return p.BaseProvider.Exec(query, args...)
}

// Execute the query string builded from given QueryBuilder, InsertBuilder,
// UpdateBuilder or DeleteBuilder, and check the affected row counts.
//
// Use BaseProvider.Exec() method to direct execute query string.
func (p *TableProvider) ExecResult(b pd.Builder) (int64, error) {
	query, args := b.Build(p.debug)
	return p.BaseProvider.ExecResult(query, args...)
}

// Insert the given rows into target table and return inserted row id of
// single value, or inserted rows count of multiple values.
//
// Use BaseProvider.Insert() method to direct execute query string.
func (p *TableProvider) Insert(b pd.Builder) (int64, error) {
	if ib, ok := b.(*builder.InsertBuilder); ok {
		query, args := b.Build(p.debug)
		if cnt := ib.ValRows(); cnt <= 0 {
			return -1, invar.ErrInvalidData
		} else if cnt == 1 {
			return p.BaseProvider.Insert(query, args...)
		}
		return p.BaseProvider.ExecResult(query)
	}
	return 0, invar.ErrBadSQLBuilder
}

// Insert the given rows into target table, and check inserted result
// but not return insert id or counts.
//
// Use BaseProvider.Insert() method to direct execute query string.
func (p *TableProvider) InsertCheck(b pd.Builder) error {
	_, err := p.Insert(b)
	return err
}

// Insert the given rows into target table without check insert counts.
//
// Use BaseProvider.Insert() method to direct execute query string.
func (p *TableProvider) InsertUncheck(b pd.Builder) error {
	if ib, ok := b.(*builder.InsertBuilder); ok {
		if ib.ValRows() <= 0 {
			return invar.ErrInvalidData
		}
		return p.Exec(b)
	}
	return invar.ErrBadSQLBuilder
}

// Update target record by given builder to build a query string, it will
// return invar.ErrNotChanged error when none updated.
//
// Use BaseProvider.Update() method to direct execute query string.
func (p *TableProvider) Update(b pd.Builder) error {
	if ub, ok := b.(*builder.UpdateBuilder); ok {
		query, args := ub.Build(p.debug)
		return p.BaseProvider.Update(query, args...)
	}
	return invar.ErrBadSQLBuilder
}

// Delete records by the given builder to build a query string, it will
// return invar.ErrNotChanged error when none deleted.
//
// Use BaseProvider.Delete() method to direct execute query string.
func (p *TableProvider) Delete(b pd.Builder) error {
	if rb, ok := b.(*builder.DeleteBuilder); ok {
		query, args := rb.Build(p.debug)
		return p.BaseProvider.Delete(query, args...)
	}
	return invar.ErrBadSQLBuilder
}

// Excute multiple transactions, it will rollback when cased one error.
//
//	// Excute 4 transactions in callback with different query1 ~ 4
//	err := provider.Trans(
//		func(t *pd.Traner) error { return tr.Query(query1, func(rows *sql.Rows) error {
//				// Fetch all rows to get result datas...
//			}, args...) },
//		func(t *pd.Traner) error { return tr.Insert(query2, nil, args...) },
//		func(t *pd.Traner) error { return tr.Inserts(query3, pd.NewTxInserter(rows, func(iv *MyStruct) string {
//			return fmt.Sprintf("(%v, '%v')", iv.D1, iv.D2)
//		})),
//		func(t *pd.Traner) error { return tr.Exec(query4, args...) })
func (p *TableProvider) Trans(cbs ...pd.TranerCallback) error {
	if !p.prepared() || len(cbs) == 0 {
		return invar.ErrBadDBConnect
	}

	tx, err := p.client.DB().Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()
	traner := (*pd.Traner)(tx)
	for _, cb := range cbs {
		if err := cb(traner); err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}
