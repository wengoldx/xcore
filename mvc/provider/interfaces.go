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
	"database/sql"
)

// A interface for database client to implement.
//
// Such as mysql.MySQL, mssql.MSSQL, sqlite.Sqlite database connection.
type DBClient interface {
	DB() *sql.DB    // Return database connected client.
	Connect() error // Connect client with database server.
	Close() error   // Disconect and close database client
}

// A interface for set DBClient client to data provider.
//
// Such as provider.BaseProvider, provider.TableProviderImpl.
type Provider interface {
	SetClient(client DBClient)
}

// A interface implement by array elems creator to return
// out values of columns.
//
// It only for QueryBuilder.Array().
type SQLCreator interface {

	// Create a new item and return out values.
	NewItem() []any
}

// A interface implement by QUID builder to build
// a sql string for database access.
type SQLBuilder interface {

	// Build sql string and return args.
	//
	//	@return string Builded standard SQL query string.
	//	@return []any  SQL args for builded query string.
	Build(debug ...bool) (string, []any)
}

// A interface implement by builder.BuilderImpl to
// suport SQL string build utils.
type BaseBuilder interface {
	SetProvider(provider TableProvider) // Set TableProvider instance.
	HasProvider() bool                  // Check TableProvider whether inited.

	/* ------------------------------------------------------------------- */
	/* SQL String Build Utils                                              */
	/* ------------------------------------------------------------------- */

	FormatJoins(tables Joins) string
	FormatWheres(wheres Wheres, sep ...string) (string, []any)
	FormatWhereIn(field string, args []any) string
	FormatOrder(field string, desc ...bool) string
	FormatLimit(n int) string
	FormatLike(field, filter string, pattern ...string) string
	FormatInserts(values KValues) (string, string, []any)
	FormatSets(values KValues) (string, []any)
	FormatValues(values KValues) string
	CheckWhere(wheres string) string
	CheckLimit(query string) string
	BuildWheres(wheres Wheres, ins, like string, sep ...string) (string, []any)
	JoinWheres(wheres ...string) string
	JoinAndWheres(wheres ...string) string
	JoinOrWheres(wheres ...string) string
	ParseOut(out any) ([]string, []any)
}

// A interface implement by builder.QuerierImpl to suport SQL query.
type QueryBuilder interface {
	SQLBuilder
	BaseBuilder

	/* ------------------------------------------------------------------- */
	/* Provider Utils For Query                                            */
	/* ------------------------------------------------------------------- */

	Has() (bool, error)
	None() (bool, error)
	Count() (int, error)
	OneScan(cb ScanCallback) error
	OneDone(done ...DoneCallback) error
	Query(cb ScanCallback) error
	Array(cr SQLCreator) error

	/* ------------------------------------------------------------------- */
	/* SQL String Build Utils                                              */
	/* ------------------------------------------------------------------- */

	Joins(tables Joins) QueryBuilder
	Tags(tag ...string) QueryBuilder
	Outs(outs ...any) QueryBuilder
	TagOut(tag string, out any) QueryBuilder
	Parse(out any) QueryBuilder
	Wheres(where Wheres) QueryBuilder
	WhereIn(field string, args []any) QueryBuilder
	WhereSep(sep string) QueryBuilder
	OrderBy(field string, desc ...bool) QueryBuilder
	Like(field, filter string, pattern ...string) QueryBuilder
	Limit(limit int) QueryBuilder
	Reset() QueryBuilder
}

// A interface implement by builder.InserterImpl to suport SQL insert.
type InsertBuilder interface {
	SQLBuilder
	BaseBuilder

	/* ------------------------------------------------------------------- */
	/* Provider Utils For Insert                                           */
	/* ------------------------------------------------------------------- */

	Exec() error
	Insert() (int64, error)
	InsertCheck() error
	InsertUncheck() error

	/* ------------------------------------------------------------------- */
	/* SQL String Build Utils                                              */
	/* ------------------------------------------------------------------- */

	Values(row ...KValues) InsertBuilder
	ValuesSize() int
	Reset() InsertBuilder
}

// A interface implement by builder.UpdaterImpl to suport SQL update.
type UpdateBuilder interface {
	SQLBuilder
	BaseBuilder

	/* ------------------------------------------------------------------- */
	/* Provider Utils For Update                                           */
	/* ------------------------------------------------------------------- */

	Exec() error
	Update() error

	/* ------------------------------------------------------------------- */
	/* SQL String Build Utils                                              */
	/* ------------------------------------------------------------------- */

	Values(row KValues) UpdateBuilder
	Wheres(where Wheres) UpdateBuilder
	WhereIn(field string, args []any) UpdateBuilder
	WhereSep(sep string) UpdateBuilder
	Like(field, filter string, pattern ...string) UpdateBuilder
	Reset() UpdateBuilder
}

// A interface implement by builder.DeleterImpl to suport SQL delete.
type DeleteBuilder interface {
	SQLBuilder
	BaseBuilder

	/* ------------------------------------------------------------------- */
	/* Provider Utils For Delete                                           */
	/* ------------------------------------------------------------------- */

	Exec() error
	Delete() error

	/* ------------------------------------------------------------------- */
	/* SQL String Build Utils                                              */
	/* ------------------------------------------------------------------- */

	Wheres(where Wheres) DeleteBuilder
	WhereIn(field string, args []any) DeleteBuilder
	WhereSep(sep string) DeleteBuilder
	Like(field, filter string, pattern ...string) DeleteBuilder
	Limit(limit int) DeleteBuilder
	Reset() DeleteBuilder
}

// A interface implement by builder.TableProviderImpl to suport table datas access.
type TableProvider interface {
	Provider

	/* ------------------------------------------------------------------- */
	/* Create QueryBuilder, InsertBuilder, UpdateBuilder, DeleteBuilder    */
	/* ------------------------------------------------------------------- */

	Querier(t ...string) QueryBuilder
	Inserter(t ...string) InsertBuilder
	Updater(t ...string) UpdateBuilder
	Deleter(t ...string) DeleteBuilder

	/* ------------------------------------------------------------------- */
	/* Provider Utils For Query, Insert, Update, Delete                    */
	/* ------------------------------------------------------------------- */

	Has(builder QueryBuilder) (bool, error)
	None(builder QueryBuilder) (bool, error)
	Count(builder QueryBuilder) (int, error)
	Exec(builder SQLBuilder) error
	ExecResult(builder SQLBuilder) (int64, error)
	OneScan(builder QueryBuilder, cb ScanCallback) error
	OneDone(builder QueryBuilder, done DoneCallback, outs ...any) error
	OneOuts(builder QueryBuilder, outs ...any) error
	Query(builder QueryBuilder, cb ScanCallback) error
	Array(builder QueryBuilder, creator SQLCreator) error
	Insert(builder InsertBuilder) (int64, error)
	InsertCheck(builder InsertBuilder) error
	InsertUncheck(builder InsertBuilder) error
	Update(builder UpdateBuilder) error
	Delete(builder DeleteBuilder) error
}
