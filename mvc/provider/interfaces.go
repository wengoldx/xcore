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
type SQLBuilder interface {

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
	Has(b SQLBuilder) (bool, error)                   // For Query.
	None(b SQLBuilder) (bool, error)                  // For Query.
	Count(b SQLBuilder) (int, error)                  // For Query.
	OneScan(b SQLBuilder, cb ScanCallback) error      // For Query.
	OneDone(b SQLBuilder, done ...DoneCallback) error // For Query.
	Query(b SQLBuilder, cb ScanCallback) error        // For Query.
	Array(b SQLBuilder, cr Creator) error             // For Query.
	Insert(b SQLBuilder) (int64, error)               // For Insert.
	InsertCheck(b SQLBuilder) error                   // For Insert.
	InsertUncheck(b SQLBuilder) error                 // For Insert.
	Exec(b SQLBuilder) error                          // For Insert, Update, Delete.
	Update(b SQLBuilder) error                        // For Update.
	Delete(b SQLBuilder) error                        // For Delete.
}
