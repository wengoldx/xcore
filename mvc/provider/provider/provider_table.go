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
//	s := &SampleTable{*mysql.NewTable("sample", _logsql)}
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
	debug bool   // Debug mode for show builded query string, default false.
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
	tp := &TableProvider{}
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

// Specify the debug mode.
func WithDebug(debug bool) Option {
	return func(provider *TableProvider) { provider.debug = debug }
}

/* ------------------------------------------------------------------- */
/* Create and Return Builder Instance FOR QUID Actions                 */
/* ------------------------------------------------------------------- */

// Create a query builder to query table records.
func (p *TableProvider) Querier(t ...string) *builder.QueryBuilder {
	return builder.NewQuery(utils.Variable(t, p.table), p)
}

// Create a insert builder to insert records to table.
func (p *TableProvider) Inserter(t ...string) *builder.InsertBuilder {
	return builder.NewInsert(utils.Variable(t, p.table), p)
}

// Create a update builder to update table records.
func (p *TableProvider) Updater(t ...string) *builder.UpdateBuilder {
	return builder.NewUpdate(utils.Variable(t, p.table), p)
}

// Create a delete builder to delete table records.
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
func (p *TableProvider) Has(b pd.SQLBuilder) (bool, error) {
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
func (p *TableProvider) None(builder pd.SQLBuilder) (bool, error) {
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
func (p *TableProvider) Count(b pd.SQLBuilder) (int, error) {
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
func (p *TableProvider) OneScan(b pd.SQLBuilder, cb pd.ScanCallback) error {
	if qb, ok := b.(*builder.QueryBuilder); ok {
		query, args := qb.Build(p.debug)
		return p.BaseProvider.One(query, cb, args...)
	}
	return invar.ErrBadSQLBuilder
}

// Query the top one record and return the results without scaner
// callback, it canbe set the finally done callback called when
// result success read.
func (p *TableProvider) OneDone(b pd.SQLBuilder, done ...pd.DoneCallback) error {
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
func (p *TableProvider) Query(b pd.SQLBuilder, cb pd.ScanCallback) error {
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
//		datas= append(datas, iv)
//		return []any{&iv.Name}
//	}, /* func(iv *MyAcc) {} */) // or append parser function.
//	h.Querier().Outs("name").Wheres(pd.Wheres{"role=?": "admin"}).Array(creator)
func (p *TableProvider) Array(b pd.SQLBuilder, creator pd.Creator) error {
	return p.Query(b, func(rows *sql.Rows) error {
		item, outs := creator.CreateItem()
		if err := rows.Scan(outs...); err != nil {
			return err
		}
		return creator.ParseItem(item)
	})
}

// Query single column values by given builder builded query string,
// and read datas from ItemScaner instance.
//
//	names := []string{}
//	scaner := pd.NewScaner(&names/* , func(iv *string) {} */)
//	h.Querier().Outs("name").Wheres(pd.Wheres{"role=?": "admin"}).Column(scaner)
func (p *TableProvider) Column(b pd.SQLBuilder, scaner pd.Scaner) error {
	return p.Query(b, func(rows *sql.Rows) error {
		out := scaner.CreateItem()
		if err := rows.Scan(&out); err != nil {
			return err
		}
		return scaner.AppendItem(out)
	})
}

// Execute the query string builded from given QueryBuilder, InsertBuilder,
// UpdateBuilder or DeleteBuilder, it not check the affected row counts.
//
// Use BaseProvider.Exec() method to direct execute query string.
func (p *TableProvider) Exec(b pd.SQLBuilder) error {
	query, args := b.Build(p.debug)
	return p.BaseProvider.Exec(query, args...)
}

// Execute the query string builded from given QueryBuilder, InsertBuilder,
// UpdateBuilder or DeleteBuilder, and check the affected row counts.
//
// Use BaseProvider.Exec() method to direct execute query string.
func (p *TableProvider) ExecResult(b pd.SQLBuilder) (int64, error) {
	query, args := b.Build(p.debug)
	return p.BaseProvider.ExecResult(query, args...)
}

// Insert the given rows into target table and return inserted row id of
// single value, or inserted rows count of multiple values.
//
// Use BaseProvider.Insert() method to direct execute query string.
func (p *TableProvider) Insert(b pd.SQLBuilder) (int64, error) {
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
func (p *TableProvider) InsertCheck(b pd.SQLBuilder) error {
	_, err := p.Insert(b)
	return err
}

// Insert the given rows into target table without check insert counts.
//
// Use BaseProvider.Insert() method to direct execute query string.
func (p *TableProvider) InsertUncheck(b pd.SQLBuilder) error {
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
func (p *TableProvider) Update(b pd.SQLBuilder) error {
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
func (p *TableProvider) Delete(b pd.SQLBuilder) error {
	if rb, ok := b.(*builder.DeleteBuilder); ok {
		query, args := rb.Build(p.debug)
		return p.BaseProvider.Delete(query, args...)
	}
	return invar.ErrBadSQLBuilder
}
