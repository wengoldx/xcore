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
//	s := &SampleProvider{*mysql.GetTabler(
//		provider.WithTable("sample"),  //set table name.
//		provider.WithDriver("sqlite"), //set builder driver, default 'sql'.
//	)}
//
// Use mysql.GetTabler(), mysql.GetTabler() sqlite.GetTabler() of mvc inner packages
// to create TableProvider with connected mysql, mssql, sqlite database client.
type TableProvider struct {
	BaseProvider
	table string // Table name
	debug bool   // Debug mode, default false.
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

// Specify the debug mode.
func WithDebug(debug bool) Option {
	return func(provider *TableProvider) {
		provider.debug = debug
	}
}

// Specify the builder sql driver, one of 'sql', 'sqlite'.
func WithDriver(driver string) Option {
	return func(provider *TableProvider) {
		if provider.Builder != nil {
			provider.Builder.driver = driver
		}
	}
}

/* ------------------------------------------------------------------- */
/* Create and Return Builder Instance FOR QIUD Actions                 */
/* ------------------------------------------------------------------- */

func (p *TableProvider) Q() *QueryBuilder  { return NewQuery(p.table, p.Builder.driver).Master(p) }
func (p *TableProvider) I() *InsertBuilder { return NewInsert(p.table, p.Builder.driver).Master(p) }
func (p *TableProvider) U() *UpdateBuilder { return NewUpdate(p.table, p.Builder.driver).Master(p) }
func (p *TableProvider) D() *DeleteBuilder { return NewDelete(p.table, p.Builder.driver).Master(p) }

/* ------------------------------------------------------------------- */
/* Using Builder To Construct Query String For Database Access         */
/* ------------------------------------------------------------------- */

// Setup TableProvider with database client and options.
func (p *TableProvider) Setup(client DBClient, opts ...Option) {
	p.BaseProvider = *NewProvider(client)
	for _, optFunc := range opts {
		optFunc(p)
	}
}

// Check the target record whether exist by the given QueryBuilder to
// build query string.
//
// Use None() method to check whether unexist.
func (p *TableProvider) Has(builder *QueryBuilder) (bool, error) {
	query, args := builder.Build(p.debug)
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
	query, args := builder.Build(p.debug)
	return p.BaseProvider.Count(query, args...)
}

// Execute the query string builded from given QueryBuilder, InsertBuilder,
// UpdateBuilder or DeleteBuilder, it not check the affected row counts.
//
// Use BaseProvider.Exec() method to direct execute query string.
func (p *TableProvider) Exec(builder SQLBuilder) error {
	query, args := builder.Build(p.debug)
	return p.BaseProvider.Exec(query, args...)
}

// Execute the query string builded from given QueryBuilder, InsertBuilder,
// UpdateBuilder or DeleteBuilder, and check the affected row counts.
//
// Use BaseProvider.Exec() method to direct execute query string.
func (p *TableProvider) ExecResult(builder SQLBuilder) (int64, error) {
	query, args := builder.Build(p.debug)
	return p.BaseProvider.ExecResult(query, args...)
}

// Query one record by given builder builded query string, and read datas
// from scan callback.
//
// Use BaseProvider.One() method to direct execute query string.
func (p *TableProvider) One(builder *QueryBuilder, cb ScanCallback) error {
	query, args := builder.Build(p.debug)
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
	query, args := builder.Build(p.debug)
	return p.BaseProvider.OneDone(query, outs, done, args...)
}

// Query records by given builder builded query string, and read datas
// from scan callback.
//
// Use BaseProvider.Query() method to direct execute query string.
func (p *TableProvider) Query(builder *QueryBuilder, cb ScanCallback) error {
	query, args := builder.Build(p.debug)
	return p.BaseProvider.Query(query, cb, args...)
}

// Insert the given rows into target table and return inserted row id of
// single value, or inserted rows count of multiple values.
//
// Use BaseProvider.Insert() method to direct execute query string.
func (p *TableProvider) Insert(builder *InsertBuilder) (int64, error) {
	query, args := builder.Build(p.debug)
	if cnt := len(builder.rows); cnt <= 0 {
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
func (p *TableProvider) InsertCheck(builder *InsertBuilder) error {
	_, err := p.Insert(builder)
	return err
}

// Insert the given rows into target table without check insert counts.
//
// Use BaseProvider.Insert() method to direct execute query string.
func (p *TableProvider) InsertUncheck(builder *InsertBuilder) error {
	if cnt := len(builder.rows); cnt <= 0 {
		return invar.ErrInvalidData
	}
	return p.Exec(builder)
}

// Update target record by given builder to build a query string, it will
// return invar.ErrNotChanged error when none updated.
//
// Use BaseProvider.Update() method to direct execute query string.
func (p *TableProvider) Update(builder *UpdateBuilder) error {
	query, args := builder.Build(p.debug)
	return p.BaseProvider.Update(query, args...)
}

// Delete records by the given builder to build a query string, it will
// return invar.ErrNotChanged error when none deleted.
//
// Use BaseProvider.Delete() method to direct execute query string.
func (p *TableProvider) Delete(builder *DeleteBuilder) error {
	query, args := builder.Build(p.debug)
	return p.BaseProvider.Delete(query, args...)
}
