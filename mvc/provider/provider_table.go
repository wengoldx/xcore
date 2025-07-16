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
	"github.com/wengoldx/xcore/invar"
)

// Table provider for using builder to build query string and args for
// database datas access.
//
// Usage: Define the custom provider as follow code.
//
//	// define the custom provider.
//	type SampleProvider struct {
//		provider.TableProvider
//	}
//	s := &SampleProvider{*mysqlc.GetTabler(
//		provider.WithTable("sample"), //set table name.
//	)}
//
// Use mysqlc.GetTabler(), mysqlc.GetTabler() to create TableProvider with
// connected mysql or mssql database client.
type TableProvider struct {
	BaseProvider
	table string // Table name
}

var _ TableSetup = (*TableProvider)(nil)

// Create a TableProvider with given database client.
func NewTabler(client DBClient, opts ...Option) *TableProvider {
	tp := &TableProvider{}
	tp.Setup(client, opts...)
	return tp
}

// The setter for set TableProvider fields.
type Option func(provider *TableProvider)

// Specify the table name.
func WithTable(table string) Option {
	return func(provider *TableProvider) {
		provider.table = table
	}
}

/* ------------------------------------------------------------------- */
/* Create and Return Builder Instance                                  */
/* ------------------------------------------------------------------- */

func (p *TableProvider) Querier() *QueryBuilder   { return NewQuery(p.table).Master(p) }
func (p *TableProvider) Inserter() *InsertBuilder { return NewInsert(p.table).Master(p) }
func (p *TableProvider) Updater() *UpdateBuilder  { return NewUpdate(p.table).Master(p) }
func (p *TableProvider) Deleter() *DeleteBuilder  { return NewDelete(p.table).Master(p) }

/* ------------------------------------------------------------------- */
/* Using Builder To Construct Query String For Database Access         */
/* ------------------------------------------------------------------- */

// Setup TableProvider with database client and options.
func (p *TableProvider) Setup(client DBClient, opts ...Option) {
	p.BaseProvider = BaseProvider{client, &BaseBuilder{}}
	for _, optFunc := range opts {
		optFunc(p)
	}
}

// Check the target record whether exist by the given QueryBuilder to
// build query string.
//
// Use None() method to check whether unexist.
func (p *TableProvider) Has(builder *QueryBuilder) (bool, error) {
	query, args := builder.Build()
	return p.BaseProvider.Has(query, args...)
}

// Check the target record whether unexist by the given QueryBuilder to
// build query string.
//
// Use Has() method to check has result.
func (p *TableProvider) None(builder *QueryBuilder) (bool, error) {
	has, err := p.Has(builder)
	return !has, err
}

// Count records by the given builder to build a query string, it will
// return 0 when notfound anyone.
//
// Use BaseProvider.Count() method to direct execute query string.
func (p *TableProvider) Count(builder *QueryBuilder) (int, error) {
	query, args := builder.Build()
	return p.BaseProvider.Count(query, args...)
}

// Execute the query string builded from given QueryBuilder, InsertBuilder,
// UpdateBuilder or DeleteBuilder, it not check the affected row counts.
//
// Use BaseProvider.Exec() method to direct execute query string.
func (p *TableProvider) Exec(builder SQLBuilder) error {
	query, args := builder.Build()
	return p.BaseProvider.Exec(query, args...)
}

// Execute the query string builded from given QueryBuilder, InsertBuilder,
// UpdateBuilder or DeleteBuilder, and check the affected row counts.
//
// Use BaseProvider.Exec() method to direct execute query string.
func (p *TableProvider) ExecResult(builder SQLBuilder) (int64, error) {
	query, args := builder.Build()
	return p.BaseProvider.ExecResult(query, args...)
}

// Query one record by given builder builded query string, and read datas
// from scan callback.
//
// Use BaseProvider.One() method to direct execute query string.
func (p *TableProvider) One(builder *QueryBuilder, cb ScanCallback) error {
	query, args := builder.Build()
	return p.BaseProvider.One(query, cb, args...)
}

// Query one record by given builder builded query string, and return the
// result datas by given outs params.
//
// Use BaseProvider.OneDone() method to direct execute query string.
func (p *TableProvider) OneOuts(builder *QueryBuilder, outs ...any) error {
	return p.OneDone(builder, nil, outs...)
}

// Query one record by given builder builded query string, and return the
// result datas by given outs params, finally call done callback to translate
// the outs datas before provider method returned.
//
// Use BaseProvider.OneDone() method to direct execute query string.
func (p *TableProvider) OneDone(builder *QueryBuilder, done DoneCallback, outs ...any) error {
	query, args := builder.Build()
	return p.BaseProvider.OneDone(query, outs, done, args...)
}

// Query records by given builder builded query string, and read datas
// from scan callback.
//
// Use BaseProvider.Query() method to direct execute query string.
func (p *TableProvider) Query(builder *QueryBuilder, cb ScanCallback) error {
	query, args := builder.Build()
	return p.BaseProvider.Query(query, cb, args...)
}

// Insert the given rows into target table and return inserted row id of
// single value, or inserted rows count of multiple values.
//
// Use BaseProvider.Insert() method to direct execute query string.
func (p *TableProvider) Insert(builder *InsertBuilder) (int64, error) {
	query, args := builder.Build()
	if cnt := len(builder.rows); cnt <= 0 {
		return -1, invar.ErrInvalidData
	} else if cnt == 1 {
		return p.BaseProvider.Insert(query, args...)
	} else {
		return p.BaseProvider.ExecResult(query)
	}
}

// Insert the given rows into target table without check insert counts.
//
// Use BaseProvider.Insert() method to direct execute query string.
func (p *TableProvider) InsertUncheck(builder *InsertBuilder) error {
	if cnt := len(builder.rows); cnt <= 0 {
		return invar.ErrInvalidData
	}
	query, args := builder.Build()
	return p.BaseProvider.Exec(query, args)
}

// Update target record by given builder to build a query string, it will
// return invar.ErrNotChanged error when none updated.
//
// Use BaseProvider.Update() method to direct execute query string.
func (p *TableProvider) Update(builder *UpdateBuilder) error {
	query, args := builder.Build()
	return p.BaseProvider.Update(query, args...)
}

// Delete records by the given builder to build a query string, it will
// return invar.ErrNotChanged error when none deleted.
//
// Use BaseProvider.Delete() method to direct execute query string.
func (p *TableProvider) Delete(builder *DeleteBuilder) error {
	query, args := builder.Build()
	return p.BaseProvider.Delete(query, args...)
}
