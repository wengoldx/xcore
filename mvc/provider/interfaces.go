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
// query tags and out values.
type SQLModel interface {

	// Return query target columns and out values.
	//
	//	@return KValues Target columns and bind out values.
	//
	// # USAGE
	//
	//	type MyModel struct { Name string }
	//	func (m *MyModel) OutTags() pd.KValues {
	//		return pd.KValues{"name": &m.Name}
	//	}
	GetTagOuts() KValues
}

// A interface implement by SQL model struct to return
// query tags and out values of array items.
type SQLItemCreator interface {

	// Return query target columns.
	//
	// # USAGE
	//
	//	type MyModel struct { Name string }
	//	func (m *MyModel) GetTags() []string {
	//		return []string{"name"}
	//	}
	GetTags() []string

	// Create a new item and return out values.
	//
	// # USAGE
	//
	//	type MyModel struct { Name string }
	//	func (m *MyModel) GetOuts() (any, []any) {
	//		item := &MyModel{}
	//		return &item, &item.Name
	//	}
	GetOuts() (any, []any)
}
