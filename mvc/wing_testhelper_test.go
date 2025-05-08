// Copyright (c) 2018-Now Dunyu All Rights Reserved.
//
// Author      : https://www.wengold.net
// Email       : support@wengold.net
//
// Prismy.No | Date       | Modified by. | Description
// -------------------------------------------------------------------
// 00001       2019/05/22   yangping       New version
// -------------------------------------------------------------------

package mvc

import (
	"math/rand/v2"
	"strconv"
	"strings"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql" // mysql
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
		Order   string
	}{
		{"Insert then query by id", time.Now().Unix(), time.Now().Format("2006-01-02 15:04:05.000"), "id"},
		{"Insert then query by num", time.Now().Unix(), time.Now().Format("2006-01-02 15:04:05"), "number"},
		{"Insert then query by time", time.Now().Unix(), time.Now().Format("2006-01-02 15:04:05"), "creates"},
	}

	// connect test database first.
	if err := setupTestDatabase(); err != nil {
		t.Fatal("Connect test database, err:", err)
	}
	defer CloseMySQL("mysql")

	ut, table := UTest(), "UnitTest"
	for _, c := range cases {
		t.Run(c.Case, func(t *testing.T) {
			if last, err := insertTestValues(ut, c.Case, c.Creates, c.Number); err != nil {
				t.Fatal("Insert a test record, err:", err)
			} else if id, err := ut.LastID(table, "number"); err != nil {
				t.Fatal("Query last id, err:", err)
			} else if last != id {
				t.Fatal("Verify last id failed!")
			} else if text, err := ut.LastField(table, "text", c.Order); err != nil {
				t.Fatal("Query last field, err:", err)
			} else if text != c.Case {
				t.Fatal("Verify last field failed!")
			} else {
				t.Log("Inserted id:", id, "- Order by", c.Order+",", "text:", text)
				ut.DelOne("UnitTest", "id", last)
			}
		})
	}
}

// Test Clear, LastIDs, DelMults, Datas.
func TestRecordMults(t *testing.T) {
	cases := []struct {
		Case    string
		Number  int64
		Creates string
	}{
		{"Insert and query mults", time.Now().Unix(), time.Now().Format("2006-01-02 15:04:05.000")},
	}

	// connect test database first.
	if err := setupTestDatabase(); err != nil {
		t.Fatal("Connect test database, err:", err)
	}
	defer CloseMySQL("mysql")

	ut, table := UTest(), "UnitTest"
	for _, c := range cases {
		t.Run(c.Case, func(t *testing.T) {
			ut.Clear(table)

			// insert 3 test record for next query and delete.
			if id1, err := insertTestValues(ut, c.Case, c.Creates, c.Number+rand.Int64()); err != nil {
				t.Fatal("Insert record 1, err:", err)
			} else if id2, err := insertTestValues(ut, c.Case, c.Creates, c.Number+rand.Int64()); err != nil {
				t.Fatal("Insert record 2, err:", err)
			} else if id3, err := insertTestValues(ut, c.Case, c.Creates, c.Number+rand.Int64()); err != nil {
				t.Fatal("Insert record 3, err:", err)
			} else {

				// test query last ids.
				lasts, err := ut.LastIDs(table, "id", 0)
				if err != nil {
					t.Fatal("Query remained ids, err:", err)
				} else if len(lasts) < 3 {
					t.Fatal("Not enough inserted ids!")
				}
				t.Log("Queried last ids:", lasts[0], lasts[1], lasts[2])

				strids := []string{}
				// test query datas.
				for _, id := range lasts {
					strids = append(strids, strconv.FormatInt(id, 10))
				}
				values, err := ut.Datas(table, "number", "id", strings.Join(strids, ","))
				if err != nil {
					t.Fatal("Query inserted numbers, err:", err)
				}
				t.Log("Queried numbers:", strings.Join(values, ","))

				// test delete multi record by ids.
				strids = strids[:len(strids)-1]
				ut.DelMults(table, "id", strings.Join(strids, ","))

				// test last id verify.
				if lasts[2] != id1 && lasts[2] != id2 && lasts[2] != id3 {
					t.Fatal("Queried ids exist no-inserted id!")
				}

				ut.Clear(table)
			}
		})
	}
}

// -------------------------------------------------------------------
// Private methods define.
// -------------------------------------------------------------------

// Connect the fixed test database with secure account.
func setupTestDatabase() error {
	/* instead the database valid configs before excute unit test! */
	confs := &MyConfs{Host: "192.168.1.100:3306", User: "user", Pwd: "123456", Name: "testdb"}
	if err := OpenMySQL2("utf8mb4", confs); err != nil {
		return err
	}
	return nil
}

// Insert test record data to test database, and return inserted id.
func insertTestValues(t *utestHelper, text, creates string, number int64) (int64, error) {
	return t.Insert("INSERT UnitTest (text, number, creates) VALUE (?, ?, ?)", text, number, creates)
}
