// Copyright (c) 2018-Now Dunyu All Rights Reserved.
//
// Author      : https://www.wengold.net
// Email       : support@wengold.net
//
// Prismy.No | Date       | Modified by. | Description
// -------------------------------------------------------------------
// 00001       2019/05/22   yangping       New version
// -------------------------------------------------------------------

package tester

import (
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql" // mysql
	pd "github.com/wengoldx/xcore/mvc/provider"
	"github.com/wengoldx/xcore/mvc/provider/mysqlc"
)

// -------------------------------------------------------------------
// USAGE: Enter ~/xcore/secure, and excute command to test.
//
//	go test -v -cover
// -------------------------------------------------------------------

// Test OpenMySQL2, LastID, DelOne.
func TestRecordOne(t *testing.T) {
	cases := []struct {
		Case    string // as text
		Number  int64
		Creates string
		Tag     string
	}{
		{"Insert & Query Num ", time.Now().Unix(), time.Now().Format("2006-01-02 15:04:05"), "number"},
	}

	// connect test database first.
	if err := setupTestDatabase(); err != nil {
		t.Fatal("Connect test database, err:", err)
	}
	defer mysqlc.Close("mysql")

	ut, text := NewHelper(), ""
	for _, c := range cases {
		t.Run(c.Case, func(t *testing.T) {
			if last, err := insertTestValues(ut, c.Case, c.Creates, c.Number); err != nil {
				t.Fatal("Insert a test record, err:", err)
			} else if id, err := ut.LastID(_ut_table); err != nil {
				t.Fatal("Query last id, err:", err)
			} else if last != id {
				t.Fatal("Verify last id failed!")
			} else if err := ut.Target(_ut_table, "text", pd.Wheres{c.Tag: c.Number}, &text); err != nil {
				t.Fatal("Query last field, err:", err)
			} else if text != c.Case {
				t.Fatal("Verify last field failed!")
			} else {
				t.Log("Inserted id:", id, "- Query by", c.Tag+",", "text:", text)
				ut.DeleteBy(_ut_table, pd.Wheres{"id": last})
			}
		})
	}
}

// Test Clear, LastIDs, DelMults, Datas.
func TestRecordMults(t *testing.T) {
	// cases := []struct {
	// 	Case    string
	// 	Number  int64
	// 	Creates string
	// }{
	// 	{"Insert and query mults", time.Now().Unix(), time.Now().Format("2006-01-02 15:04:05.000")},
	// }

	// // connect test database first.
	// if err := setupTestDatabase(); err != nil {
	// 	t.Fatal("Connect test database, err:", err)
	// }
	// defer CloseMySQL("mysql")

	// ut := UTest()
	// for _, c := range cases {
	// 	t.Run(c.Case, func(t *testing.T) {
	// 		ut.Clear(_ut_table)

	// 		// insert 3 test record for next query and delete.
	// 		if id1, err := insertTestValues(ut, c.Case, c.Creates, c.Number+rand.Int64()); err != nil {
	// 			t.Fatal("Insert record 1, err:", err)
	// 		} else if id2, err := insertTestValues(ut, c.Case, c.Creates, c.Number+rand.Int64()); err != nil {
	// 			t.Fatal("Insert record 2, err:", err)
	// 		} else if id3, err := insertTestValues(ut, c.Case, c.Creates, c.Number+rand.Int64()); err != nil {
	// 			t.Fatal("Insert record 3, err:", err)
	// 		} else {

	// 			// test query last ids.
	// 			lasts, err := ut.LastIDs(_ut_table, "id", 0)
	// 			if err != nil {
	// 				t.Fatal("Query remained ids, err:", err)
	// 			} else if len(lasts) < 3 {
	// 				t.Fatal("Not enough inserted ids!")
	// 			}
	// 			t.Log("Queried last ids:", lasts[0], lasts[1], lasts[2])

	// 			strids := []string{}
	// 			// test query datas.
	// 			for _, id := range lasts {
	// 				strids = append(strids, strconv.FormatInt(id, 10))
	// 			}
	// 			values, err := ut.Datas(_ut_table, "number", "id", strings.Join(strids, ","))
	// 			if err != nil {
	// 				t.Fatal("Query inserted numbers, err:", err)
	// 			}
	// 			t.Log("Queried numbers:", strings.Join(values, ","))

	// 			// test delete multi record by ids.
	// 			strids = strids[:len(strids)-1]
	// 			ut.DelMults(_ut_table, "id", strings.Join(strids, ","))

	// 			// test last id verify.
	// 			if lasts[2] != id1 && lasts[2] != id2 && lasts[2] != id3 {
	// 				t.Fatal("Queried ids exist no-inserted id!")
	// 			}

	// 			ut.Clear(_ut_table)
	// 		}
	// 	})
	// }
}

// -------------------------------------------------------------------
// Private methods define.
// -------------------------------------------------------------------

const _ut_table = "UnitTest"

// Connect the fixed test database with secure account.
//
//	NOTICE: The testdb must create as as follows script, and change the connect
//	configs to valid before excute go test of wing_texthepler_test.go:
//
//	`
//	CREATE DATABASE IF NOT EXISTS testdb CHARACTER SET utf8mb4 COLLATE utf8mb4_bin;
//
//	USE testdb;
//
//	CREATE TABLE IF NOT EXISTS UnitTest (
//	    id               int                NOT NULL AUTO_INCREMENT,
//	    text             varchar (256)      DEFAULT '',
//	    number           int                DEFAULT 0,
//	    creates          datetime           DEFAULT CURRENT_TIMESTAMP,
//	    PRIMARY KEY (id)
//	) COMMENT='UnitTest table';
//	`
func setupTestDatabase() error {
	opts := mysqlc.DefaultOptions("mysql")
	opts.Host = "192.168.1.100:3306"
	opts.User = "user"
	opts.Password = "123456"
	opts.Database = "testdb"

	if err := mysqlc.OpenWithOptions(opts, "utf8mb4"); err != nil {
		return err
	}
	return nil
}

// Insert test record data to test database, and return inserted id.
func insertTestValues(t *helper, text, creates string, number int64) (int64, error) {
	return t.AddWithID(_ut_table, map[string]any{
		"text": text, "number": number, "creates": creates,
	})
}
