// Copyright (c) 2018-Now Dunyu All Rights Reserved.
//
// Author      : https://www.wengold.net
// Email       : support@wengold.net
//
// Prismy.No | Date       | Modified by. | Description
// -------------------------------------------------------------------
// 00001       2019/05/22   yangping       New version
// -------------------------------------------------------------------

package provider

import (
	"fmt"
	"testing"

	pd "github.com/wengoldx/xcore/mvc/provider"
	"github.com/wengoldx/xcore/mvc/provider/builder"
	"github.com/wengoldx/xcore/utils"
)

func TestQueryBuilder(t *testing.T) {
	q := builder.NewQuery("test_table")

	query, arg := q.Tags("f1", "f2", "f3").Build()
	fmt.Println(query, "-", arg)

	query, arg = q.Wheres(pd.Wheres{"w1=?": 1, "w2<>?": 2}).Build()
	fmt.Println(query, "-", arg)

	query, arg = q.WhereIn("in3", utils.ToAnys([]int64{3, 4, 5})).Build()
	fmt.Println(query, "-", arg)

	query, arg = q.WhereIn("in3", utils.ToAnys([]float64{1.2, 0.34, 5.6789})).Build()
	fmt.Println(query, "-", arg)

	query, arg = q.WhereIn("in3", utils.ToAnys([]bool{true, false})).Build()
	fmt.Println(query, "-", arg)

	query, arg = q.OrderBy("order4", true).Build()
	fmt.Println(query, "-", arg)

	query, arg = q.Like("like5", "keyword").Build()
	fmt.Println(query, "-", arg)

	query, arg = q.Limit(6).Build()
	fmt.Println(query, "-", arg)

	q.Reset()
	query, arg = q.Build()
	fmt.Println(query, "-", arg)
}

func TestInsertBuilder(t *testing.T) {
	i := builder.NewInsert("test_table")
	v1 := pd.KValues{"manager": "123456", "perfer": 2}
	v2 := pd.KValues{"manager": "qwertyu", "perfer": 6, "obj": nil}
	v3 := pd.KValues{"": 123, "manager": "poiuytr", "perfer": 8, "obj": nil}

	query, arg := i.Values(v1).Build()
	fmt.Println(query, "-", arg)

	query, arg = i.Values(v3, v2, v1).Build()
	fmt.Println(query, "-", arg)
}
