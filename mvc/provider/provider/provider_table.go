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
)

// Table provider for using builder to build query string and args for
// database datas access.
//
// Usage: Define the custom provider as follow code.
//
//	// define the custom provider.
//	type SampleTable struct { provider.ProviderImpl }
//	s := &SampleTable{*mysql.NewTable("sample", _logsql)}
//	mysql.SetClient(s)
//
//	// or create directory if client exist.
//	s := provider.NewTableProvider(client,
//		provider.WithTable("sample"), provider.WithDebug(true))
//
// Use mysql.NewTable(), mysql.NewTable() sqlite.NewTable() of mvc inner packages
// to create TableProvider with connected mysql, mssql, sqlite database client.
type ProviderImpl struct {
	BaseProvider
	table string // Table name.
	debug bool   // Debug mode for show builded query string, default false.
}

var _ pd.TableProvider = (*ProviderImpl)(nil)

// Create a TableProvider with given database client.
//
// # WARNING:
//
// This method call by 'sqlite3', 'mysql', 'mssql' module called to create
// target table and bind with connected database client instance.
func NewTableProvider(client pd.DBClient, opts ...Option) pd.TableProvider {
	tp := &ProviderImpl{}
	tp.BaseProvider = *NewBaseProvider(client)
	for _, optFunc := range opts {
		optFunc(tp)
	}
	return tp
}

// The setter for set TableProvider options.
type Option func(provider *ProviderImpl)

// Specify the table name.
func WithTable(table string) Option {
	return func(provider *ProviderImpl) { provider.table = table }
}

// Specify the debug mode.
func WithDebug(debug bool) Option {
	return func(provider *ProviderImpl) { provider.debug = debug }
}

/* ------------------------------------------------------------------- */
/* Create and Return Builder Instance FOR QUID Actions                 */
/* ------------------------------------------------------------------- */

// Create a query builder to query table records.
func (p *ProviderImpl) Querier(t ...string) pd.QueryBuilder {
	query := builder.NewQuery(p.getTable(t...))
	query.SetProvider(p)
	return query
}

// Create a insert builder to insert records to table.
func (p *ProviderImpl) Inserter(t ...string) pd.InsertBuilder {
	insert := builder.NewInsert(p.getTable(t...))
	insert.SetProvider(p)
	return insert
}

// Create a update builder to update table records.
func (p *ProviderImpl) Updater(t ...string) pd.UpdateBuilder {
	update := builder.NewUpdate(p.getTable(t...))
	update.SetProvider(p)
	return update
}

// Create a delete builder to delete table records.
func (p *ProviderImpl) Deleter(t ...string) pd.DeleteBuilder {
	delete := builder.NewDelete(p.getTable(t...))
	delete.SetProvider(p)
	return delete
}

// Return target table name or current provider table name.
func (p *ProviderImpl) getTable(t ...string) string {
	if len(t) > 0 && t[0] != "" {
		return t[0]
	}
	return p.table
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
func (p *ProviderImpl) Has(builder pd.QueryBuilder) (bool, error) {
	query, args := builder.Tags("*").Build(p.debug)
	return p.BaseProvider.Has(query, args...)
}

// Check the target record whether unexist by the given QueryBuilder to
// build query string, it no-need set any tags.
//
// # USAGE
//
//	h.Querier().Wheres(pd.Wheres{"account=?": acc}).None()
//
// Use Has() method to check has result.
func (p *ProviderImpl) None(builder pd.QueryBuilder) (bool, error) {
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
func (p *ProviderImpl) Count(builder pd.QueryBuilder) (int, error) {
	query, args := builder.Tags("COUNT(*)").Build(p.debug)
	return p.BaseProvider.Count(query, args...)
}

// Execute the query string builded from given QueryBuilder, InsertBuilder,
// UpdateBuilder or DeleteBuilder, it not check the affected row counts.
//
// Use BaseProvider.Exec() method to direct execute query string.
func (p *ProviderImpl) Exec(builder pd.SQLBuilder) error {
	query, args := builder.Build(p.debug)
	return p.BaseProvider.Exec(query, args...)
}

// Execute the query string builded from given QueryBuilder, InsertBuilder,
// UpdateBuilder or DeleteBuilder, and check the affected row counts.
//
// Use BaseProvider.Exec() method to direct execute query string.
func (p *ProviderImpl) ExecResult(builder pd.SQLBuilder) (int64, error) {
	query, args := builder.Build(p.debug)
	return p.BaseProvider.ExecResult(query, args...)
}

// Query one record by given builder builded query string, and read datas
// from scan callback.
//
// # NOTICE:
//	- Use BaseProvider.One() method to direct execute query string.
func (p *ProviderImpl) OneScan(builder pd.QueryBuilder, cb pd.ScanCallback) error {
	query, args := builder.Build(p.debug)
	return p.BaseProvider.One(query, cb, args...)
}

// Query one record by given builder builded query string, and return the
// result datas by given outs params, finally call done callback to translate
// the outs datas before provider method returned.
//
// # NOTICE:
//	- Use BaseProvider.OneDone() method to direct execute query string.
//	- Use QueryBuilder.OneDone() method to query result by orm model.
func (p *ProviderImpl) OneDone(builder pd.QueryBuilder, done pd.DoneCallback, outs ...any) error {
	query, args := builder.Build(p.debug)
	return p.BaseProvider.OneDone(query, outs, done, args...)
}

// Query one record by given builder builded query string, and return the
// result datas by given outs params.
//
// # NOTICE:
//	- Use BaseProvider.OneDone() method to direct execute query string.
//	- Use QueryBuilder.OneDone() method to query result by orm model.
func (p *ProviderImpl) OneOuts(builder pd.QueryBuilder, outs ...any) error {
	return p.OneDone(builder, nil, outs...)
}

// Query records by given builder builded query string, and read datas
// from scan callback.
//
// Use BaseProvider.Query() method to direct execute query string.
func (p *ProviderImpl) Query(builder pd.QueryBuilder, cb pd.ScanCallback) error {
	query, args := builder.Build(p.debug)
	return p.BaseProvider.Query(query, cb, args...)
}

// Query records by given builder builded query string, and read datas
// from ElemCreator instance.
//
//	type MyAcc struct { Name string }
//
//	accs := []*MyAcc{}
//	creator := pd.NewCreator(func(iv *MyAcc) []any {
//		datas= append(datas, iv)
//		return []any{&iv.Name}
//	})
//	h.Querier().Wheres(pd.Wheres{"role=?": "admin"}).Array(creator)
func (p *ProviderImpl) Array(builder pd.QueryBuilder, creator pd.SQLCreator) error {
	query, args := builder.Build(p.debug)
	return p.BaseProvider.Query(query, func(rows *sql.Rows) error {
		outs := creator.NewItem()
		if err := rows.Scan(outs...); err != nil {
			return err
		}
		return nil
	}, args...)
}

// Insert the given rows into target table and return inserted row id of
// single value, or inserted rows count of multiple values.
//
// Use BaseProvider.Insert() method to direct execute query string.
func (p *ProviderImpl) Insert(builder pd.InsertBuilder) (int64, error) {
	query, args := builder.Build(p.debug)
	if cnt := builder.ValuesSize(); cnt <= 0 {
		return -1, invar.ErrInvalidData
	} else if cnt == 1 {
		return p.BaseProvider.Insert(query, args...)
	} else {
		return p.BaseProvider.ExecResult(query)
	}
}

// Insert the given rows into target table, and check inserted result
// but not return insert id or counts.
//
// Use BaseProvider.Insert() method to direct execute query string.
func (p *ProviderImpl) InsertCheck(builder pd.InsertBuilder) error {
	_, err := p.Insert(builder)
	return err
}

// Insert the given rows into target table without check insert counts.
//
// Use BaseProvider.Insert() method to direct execute query string.
func (p *ProviderImpl) InsertUncheck(builder pd.InsertBuilder) error {
	if builder.ValuesSize() <= 0 {
		return invar.ErrInvalidData
	}
	return p.Exec(builder)
}

// Update target record by given builder to build a query string, it will
// return invar.ErrNotChanged error when none updated.
//
// Use BaseProvider.Update() method to direct execute query string.
func (p *ProviderImpl) Update(builder pd.UpdateBuilder) error {
	query, args := builder.Build(p.debug)
	return p.BaseProvider.Update(query, args...)
}

// Delete records by the given builder to build a query string, it will
// return invar.ErrNotChanged error when none deleted.
//
// Use BaseProvider.Delete() method to direct execute query string.
func (p *ProviderImpl) Delete(builder pd.DeleteBuilder) error {
	query, args := builder.Build(p.debug)
	return p.BaseProvider.Delete(query, args...)
}
