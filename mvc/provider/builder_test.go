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
	"fmt"
	"strconv"
	"testing"

	wt "github.com/wengoldx/xcore/utils/testx"
)

/* ------------------------------------------------------------------- */
/* For BaseBuilder Tests                                               */
/* ------------------------------------------------------------------- */

func TestFormatJoins(t *testing.T) {
	cases := []*wt.TestCase{
		wt.NewCase("Check normal datas", "table_a AS a, table_b AS b", Joins{"table_a": "a", "table_b": "b"}),
		wt.NewCase("Check empty table ", "table_a AS a", Joins{"table_a": "a", "": "b"}),
		wt.NewCase("Check empty alias ", "table_b AS b", Joins{"table_a": "", "table_b": "b"}),
		wt.NewCase("Check all emptys  ", "", Joins{"": ""}),
	}

	builder := &BaseBuilder{}
	for _, c := range cases {
		rst := builder.FormatJoins(c.Params.(Joins))
		if want := c.Want.(string); rst != want {
			t.Fatal("BaseBuilder.FormatJoins error > want:", want, "but result is", rst)
		}
	}
}

func TestFormatWheres(t *testing.T)  { /* TODO */ }
func TestFormatWhereIn(t *testing.T) { /* TODO */ }
func TestFormatOrder(t *testing.T)   { /* TODO */ }
func TestFormatLimit(t *testing.T)   { /* TODO */ }

func TestFormatLike(t *testing.T) {
	type LikeData struct {
		Field   string
		Filter  string
		Pattern []string
	}
	cases := []*wt.TestCase{
		wt.NewCase("Check normal", "filed LIKE '%%filter%%'", LikeData{"filed", "filter", []string{}}),
		wt.NewCase("Check perfix", "filed LIKE 'filter%%'", LikeData{"filed", "filter", []string{"perfix"}}),
		wt.NewCase("Check suffix", "filed LIKE '%%filter'", LikeData{"filed", "filter", []string{"suffix"}}),
		wt.NewCase("Check emptys", "", LikeData{"", "", []string{}}),
	}

	builder := &BaseBuilder{}
	for _, c := range cases {
		param := c.Params.(LikeData)
		rst := builder.FormatLike(param.Field, param.Filter, param.Pattern...)
		if want := c.Want.(string); rst != want {
			t.Fatal("BaseBuilder.FormatLike error > want:", want, "but result is", rst)
		}
	}
}

type MyTestDoc struct {
	Doc string `json:"doc"`
}

type MyTestData struct {
	Name   string     `json:"name"`
	Age    int64      `json:"age"`
	Labels []string   `json:"labels"`
	Doc    MyTestDoc  `json:"doc"`
	DocPtr *MyTestDoc `json:"docptr"`
	NilPtr *MyTestDoc `json:"nilptr"`
}

func TestParseOut(t *testing.T) {
	data := &MyTestData{}
	fmt.Println("MyTestData json tag:")

	builder := &BaseBuilder{}
	tags, outs := builder.ParseOut(data)
	for index, tag := range tags {
		fmt.Println(fmt.Sprintf("% 12s", tag), "- out:", outs[index])
		switch tag {
		case "name":
			nameptr := outs[index].(*string)
			*nameptr = "zhangsan-" + strconv.Itoa(index)
		case "age":
			ageptr := outs[index].(*int64)
			*ageptr = 20 + int64(index)
		case "labels":
			labelsptr := outs[index].(*[]string)
			*labelsptr = []string{"label1", "label2", "label3"}
		case "doc":
			docptr := outs[index].(*MyTestDoc)
			*docptr = MyTestDoc{Doc: "my-doc"}
		case "docptr":
			docpptr := outs[index].(**MyTestDoc)
			*docpptr = &MyTestDoc{Doc: "my-doc-ptr"}
		case "nilptr":
			docpptr := outs[index].(**MyTestDoc)
			*docpptr = nil
		}
	}
	fmt.Println("")
	fmt.Println("MyTestData result:")
	fmt.Println(" data.Name   =", data.Name)
	fmt.Println("     .Age    =", data.Age)
	fmt.Println("     .Labels =", data.Labels)
	fmt.Println("     .Doc    =", data.Doc)
	fmt.Println("     .DocPtr =", data.DocPtr)
	fmt.Println("     .NilPtr =", data.NilPtr)
}

// TODO
// ...
