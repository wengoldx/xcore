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
	"github.com/wengoldx/xcore/invar"
)

// Simple provider for using builder to build query string
// and args for database datas access.
type SimpleProvider struct {
	BaseProvider
	table string // Table name
}

// var _ DataProvider = (*SimpleProvider)(nil)

// Create a SimpleProvider with given database client.
func NewSimpler(client DBClient, opts ...Option) *SimpleProvider {
	sp := &SimpleProvider{
		BaseProvider: BaseProvider{client, &BaseBuilder{}},
	}

	for _, optFunc := range opts {
		optFunc(sp)
	}
	return sp
}

// The setter for set SimpleProvider fields.
type Option func(provider *SimpleProvider)

// Specify the table name.
func WithTable(table string) Option {
	return func(provider *SimpleProvider) {
		provider.table = table
	}
}

/* Create and return builder instance */

func (p *SimpleProvider) Querier() *QueryBuilder   { return NewQuery(p.table) }
func (p *SimpleProvider) Inserter() *InsertBuilder { return NewInsert(p.table) }
func (p *SimpleProvider) Updater() *UpdateBuilder  { return NewUpdate(p.table) }
func (p *SimpleProvider) Deleter() *DeleteBuilder  { return NewDelete(p.table) }

/* ------------------------------------------------------------------- */
/* Using Builder To Construct Query String For Database Access         */
/* ------------------------------------------------------------------- */

// Check the target record whether exist by the given QueryBuilder to
// build query string.
//
// Use None() method to check whether unexist.
func (p *SimpleProvider) Has(builder *QueryBuilder) (bool, error) {
	query, args := builder.Build()
	return p.BaseProvider.Has(query, args...)
}

// Check the target record whether unexist by the given QueryBuilder to
// build query string.
//
// Use Has() method to check has result.
func (p *SimpleProvider) None(builder *QueryBuilder) (bool, error) {
	has, err := p.Has(builder)
	return !has, err
}

// Count records by the given builder to build a query string, it will
// return 0 when notfound anyone.
//
// Use BaseProvider.Count() method to direct execute query string.
func (p *SimpleProvider) Count(builder *QueryBuilder) (int, error) {
	query, args := builder.Build()
	return p.BaseProvider.Count(query, args...)
}

// Execute the query string builded from given QueryBuilder, InsertBuilder,
// UpdateBuilder or DeleteBuilder, it not check the affected row counts.
//
// Use BaseProvider.Exec() method to direct execute query string.
func (p *SimpleProvider) Exec(builder SQLBuilder) error {
	query, args := builder.Build()
	return p.BaseProvider.Exec(query, args...)
}

// Execute the query string builded from given QueryBuilder, InsertBuilder,
// UpdateBuilder or DeleteBuilder, and check the affected row counts.
//
// Use BaseProvider.Exec() method to direct execute query string.
func (p *SimpleProvider) ExecResult(builder SQLBuilder) (int64, error) {
	query, args := builder.Build()
	return p.BaseProvider.ExecResult(query, args...)
}

// Query one record by given builder builded query string, and read datas
// from scan callback.
//
// Use BaseProvider.One() method to direct execute query string.
func (p *SimpleProvider) One(builder *QueryBuilder, cb ScanCallback) error {
	query, args := builder.Build()
	return p.BaseProvider.One(query, cb, args...)
}

// Query records by given builder builded query string, and read datas
// from scan callback.
//
// Use BaseProvider.Query() method to direct execute query string.
func (p *SimpleProvider) Query(builder *QueryBuilder, cb ScanCallback) error {
	query, args := builder.Build()
	return p.BaseProvider.Query(query, cb, args...)
}

// Insert the given values as a record into target table, and return
// the inserted id of the 'auto increment' primary key field.
//
// Use BaseProvider.Insert() method to direct execute query string.
func (p *SimpleProvider) Insert(builder *InsertBuilder) (int64, error) {
	if cnt := len(builder.rows); cnt != 1 {
		return -1, invar.ErrInvalidData
	}
	query, args := builder.Build()
	return p.BaseProvider.Insert(query, args...)
}

// Insert the given rows into target table without check insert counts.
//
// Use BaseProvider.Inserts() method to direct execute query string.
func (p *SimpleProvider) Inserts(builder *InsertBuilder) (int64, error) {
	if cnt := len(builder.rows); cnt < 1 {
		return -1, invar.ErrInvalidData
	}
	query, _ := builder.Build()
	return p.BaseProvider.ExecResult(query)
}

// Update target record by given builder to build a query string, it will
// return invar.ErrNotChanged error when none updated.
//
// Use BaseProvider.Update() method to direct execute query string.
func (p *SimpleProvider) Update(builder *UpdateBuilder) error {
	query, args := builder.Build()
	return p.BaseProvider.Update(query, args...)
}

// Delete records by the given builder to build a query string, it will
// return invar.ErrNotChanged error when none deleted.
//
// Use BaseProvider.Delete() method to direct execute query string.
func (p *SimpleProvider) Delete(builder *DeleteBuilder) error {
	query, args := builder.Build()
	return p.BaseProvider.Delete(query, args...)
}
