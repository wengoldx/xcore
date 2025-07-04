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

// Simple provider for using builder to build query string
// and args for database datas access.
type SimpleProvider struct {
	BaseProvider
}

// var _ DataProvider = (*SimpleProvider)(nil)

// Create a SimpleProvider with given database client.
func NewSimpler(client DBClient) *SimpleProvider {
	return &SimpleProvider{BaseProvider{client, &BaseBuilder{}}}
}

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

// Insert the given values as a record into target table, and return
// the inserted id of the 'auto increment' primary key field.
//
// Use BaseProvider.Insert() method to direct execute query string.
func (p *SimpleProvider) Insert(builder InsertBuilder) (int64, error) {
	query, args := builder.Build()
	return p.BaseProvider.Insert(query, args...)
}

// Delete records by the given builder to build a query string, it will
// return invar.ErrNotChanged error when none deleted.
//
// Use BaseProvider.Delete() method to direct execute query string.
func (p *SimpleProvider) Delete(builder *DeleteBuilder) error {
	query, args := builder.Build()
	return p.BaseProvider.Delete(query, args...)
}
