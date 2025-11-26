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

import "database/sql"

// A interface for database client implement.
type DBClient interface {
	DB() *sql.DB    // Return database connected client.
	Connect() error // Connect client with database server.
	Close() error   // Disconect and close database client
}

// A interface implement by CUDA builder to build
// a sql string for database access.
type SQLBuilder interface {

	// Build sql string and return args.
	//
	//	@return string Builded standard SQL query string.
	//	@return []any  SQL args for builded query string.
	Build(debug ...bool) (string, []any)
}

// A interface implement by SQL model struct to return
// columns names and bind values pointers.
type SQLModel interface {

	// Return model columns and values pointers.
	//	@return KValues Target columns and bind values pointers.
	//
	// # Example
	//
	//	type MyModel struct { Name string }
	//	func (m *MyModel) MapValues() pd.KValues {
	//		return pd.KValues{"name": &m.Name}  // out pointer!
	//	}
	MapOuts() KValues

	// Return model columns and values.
	//	@return KValues Target columns and bind values.
	//
	// # Example
	//
	//	type MyModel struct { Name string }
	//	func (m *MyModel) MapValues() pd.KValues {
	//		return pd.KValues{"name": m.Name}   // value data!
	//	}
	MapValues() KValues
}

// A interface implement by SQL model struct to return
// query columns and bind values of created array items.
type SQLItemCreator interface {

	// Return query target columns.
	//
	// # Example
	//
	//	type MyModel struct { Name string }
	//	func (m *MyModel) GetTags() []string {
	//		return []string{"name"}
	//	}
	GetTags() []string

	// Create a new item and return out values.
	//
	// # Example
	//
	//	type MyModel struct { Name string }
	//	func (m *MyModel) GetOuts() (any, []any) {
	//		item := &MyModel{}
	//		return &item, []any{&item.Name}
	//	}
	GetOuts() (any, []any)
}
