// Copyright (c) 2018-Now Dunyu All Rights Reserved.
//
// Author      : https://www.wengold.net
// Email       : support@wengold.net
//
// Prismy.No | Date       | Modified by. | Description
// -------------------------------------------------------------------
// 00001       2019/05/22   yangping       New version
// -------------------------------------------------------------------

package builder

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	pd "github.com/wengoldx/xcore/mvc/provider"
	wt "github.com/wengoldx/xcore/utils/xtest"
)

/* ------------------------------------------------------------------- */
/* For BaseBuilder Tests                                               */
/* ------------------------------------------------------------------- */

func TestFormatJoins(t *testing.T) {
	cases := []*wt.TestCase{
		wt.NewCase("Check normal datas", "table_a AS a, table_b AS b", pd.Joins{"table_a": "a", "table_b": "b"}),
		wt.NewCase("Check empty table ", "table_a AS a", pd.Joins{"table_a": "a", "": "b"}),
		wt.NewCase("Check empty alias ", "table_b AS b", pd.Joins{"table_a": "", "table_b": "b"}),
		wt.NewCase("Check all emptys  ", "", pd.Joins{"": ""}),
	}

	builder := NewBuilder("")
	for _, c := range cases {
		rst := builder.FormatJoins(c.Params.(pd.Joins))
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

	builder := NewBuilder("")
	for _, c := range cases {
		param := c.Params.(LikeData)
		rst := builder.FormatLike(param.Field, param.Filter, param.Pattern...)
		if want := c.Want.(string); rst != want {
			t.Fatal("BaseBuilder.FormatLike error > want:", want, "but result is", rst)
		}
	}
}

type MyTestDoc struct {
	Doc string `column:"doc"`
}

type MyTestData struct {
	Name   string     `column:"name"`
	Age    int64      `column:"age"`
	Labels []string   `column:"labels"`
	Doc    MyTestDoc  `column:"doc"`
	DocPtr *MyTestDoc `column:"docptr"`
	NilPtr *MyTestDoc `column:"nilptr"`
}

func (m *MyTestDoc) String() string {
	return m.Doc
}

func TestParseOut(t *testing.T) {
	data := &MyTestData{}
	fmt.Println("MyTestData json tag:")

	builder := NewBuilder("")
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
	fmt.Println("---")
	fmt.Println("MyTestData result:")
	fmt.Println(" data.Name   =", data.Name)
	fmt.Println("     .Age    =", data.Age)
	fmt.Println("     .Labels =", data.Labels)
	fmt.Println("     .Doc    =", data.Doc)
	fmt.Println("     .DocPtr =", data.DocPtr)
	fmt.Println("     .NilPtr =", data.NilPtr)
}

func TestItemCreator(t *testing.T) {
	datas := []*MyTestData{}
	creator := pd.NewCreator(&datas, func(iv *MyTestData) []any {
		return []any{&iv.Name, &iv.Age, &iv.Labels, &iv.Doc, &iv.DocPtr, &iv.NilPtr}
	}, func(iv *MyTestData) error {
		if iv.Age%2 == 0 {
			iv.NilPtr = iv.DocPtr
		}
		return nil
	})

	for i := 0; i < 10; i++ {
		item, fields := creator.CreateItem()
		*(fields[0].(*string)) = "zhangsan" + strconv.Itoa(i)
		*(fields[1].(*int64)) = 19 + int64(i)
		*(fields[2].(*[]string)) = []string{"label-" + strconv.Itoa(i)}
		*(fields[3].(*MyTestDoc)) = MyTestDoc{Doc: "test doc"}
		*(fields[4].(**MyTestDoc)) = &MyTestDoc{Doc: "test doc ptr!"}
		*(fields[5].(**MyTestDoc)) = nil
		/* Here for rows scaning... */
		creator.AppendItem(item)
	}

	for _, data := range datas {
		fmt.Println("Data:", data)
		if data.Age%2 == 0 && data.NilPtr == nil {
			t.Fatal("Parsed item value, not changed!")
		}
	}
}

func TestItemScaner(t *testing.T) {
	datas := []string{}
	scaner := pd.NewScaner(&datas, func(iv *string) error {
		if strings.HasSuffix(*iv, "2") {
			*iv = "changed!"
		}
		return nil
	})

	for i := 0; i < 5; i++ {
		item := scaner.CreateItem()
		*(item.(*string)) = "scaning " + strconv.Itoa(i)
		/* Here for rows scaning... */
		scaner.AppendItem(item)
	}

	for index, data := range datas {
		fmt.Println("Data:", data)
		if index == 2 && data != "changed!" {
			t.Fatal("Parsed item value, not changed!")
		}
	}
}

// TODO
// ...
