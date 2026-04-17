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
	Has(b Builder) (bool, error)                   // For Query.
	None(b Builder) (bool, error)                  // For Query.
	Count(b Builder) (int, error)                  // For Query.
	OneScan(b Builder, cb ScanCallback) error      // For Query.
	OneDone(b Builder, done ...DoneCallback) error // For Query.
	Query(b Builder, cb ScanCallback) error        // For Query.
	Array(b Builder, cr Creator) error             // For Query.
	Column(b Builder, sr Scaner) error             // For Query.
	Insert(b Builder) (int64, error)               // For Insert.
	InsertCheck(b Builder) error                   // For Insert.
	InsertUncheck(b Builder) error                 // For Insert.
	Exec(b Builder) error                          // For Insert, Update, Delete.
	Update(b Builder) error                        // For Update.
	Delete(b Builder) error                        // For Delete.
}
