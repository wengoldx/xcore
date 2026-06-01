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

// A interface implement by QUID builder to build
// a sql string for database access.
type Builder interface {

	// Build sql string and return args.
	//
	//	@return string Builded standard SQL query string.
	//	@return []any  SQL args for builded query string.
	Build(debug ...bool) (string, []any)
}

// A interface for set DBClient client to data provider.
//
// Such as provider.BaseProvider, provider.TableProvider.
type Provider interface {
	SetClient(client DBClient)
}

// A interface implement by provider.TableProvider to export utils.
type ProviderUtils interface {

	/* ------------------------------------------------------------------- */
	/* For Query Utils                                                     */
	/* ------------------------------------------------------------------- */

	Has(b Builder) (bool, error)
	None(b Builder) (bool, error)
	Count(b Builder) (int, error)
	OneScan(b Builder, cb ScanCallback) error
	OneDone(b Builder, done ...DoneCallback) error
	Query(b Builder, cb ScanCallback) error
	Array(b Builder, cr Creator) error
	Column(b Builder, sr Scaner) error

	/* ------------------------------------------------------------------- */
	/* For Insert Utils                                                    */
	/* ------------------------------------------------------------------- */

	Insert(b Builder) (int64, error)
	InsertCheck(b Builder) error
	InsertUncheck(b Builder) error

	/* ------------------------------------------------------------------- */
	/* For Update & Delete Utils                                           */
	/* ------------------------------------------------------------------- */

	Exec(b Builder) error                // For Insert, Update, Delete.
	ExecResult(b Builder) (int64, error) // For Insert, Update, Delete.
	Update(b Builder) error
	Delete(b Builder) error
}
