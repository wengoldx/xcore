// Copyright (c) 2018-Now Dunyu All Rights Reserved.
//
// Author      : https://www.wengold.net
// Email       : support@wengold.net
//
// Prismy.No | Date       | Modified by. | Description
// -------------------------------------------------------------------
// 00001       2019/05/22   yangping       New version
// -------------------------------------------------------------------

// The follow use mysql to describe database init follows:
//
// Step 1: Open and connect mysql databse.
//      -  mysql.Open("utf8mb4") ---> client ------------------+
//      -                            +------+                  |
//      -                            | conn | <-- *sql.DB      |
//      -                            +------+                  |
//                                                             |
// Step 2: Define table and generate instance.                 |
//      -  type MyTable struct {                               |
//      -      provider.TableProvider                          |
//      -  }                                                   |
//      -                                                      |
//      -  var mytable = MyTable{  --------------+             |
//      -      *mysql.NewTable("mytable", log)   |             |
//      -  }                                     |             |
//                                    mytable <--+             |
//                                   +--------+                |
//             DBClient (mysql) -->  | client | <--------------+
//                    "mytable" -->  | table  |        ^
//                                   +--------+        |
//                                                     |
// Step 3: Bind client with tables.                    |
//      -  mysql.BindTables([]pd.Provider{mytable}) ---+

// Package content provider for simple code to connect and
// access mysql or mssql database datas.
package pd
