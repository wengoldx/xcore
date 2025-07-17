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
	"testing"

	"github.com/wengoldx/xcore/utils"
)

func TestQueryBuilder(t *testing.T) {
	q := NewQuery("test_table")

	query, arg := q.Tags("f1", "f2", "f3").Build()
	fmt.Println(query, "-", arg)

	query, arg = q.Wheres(Wheres{"w1=?": 1, "w2<>?": 2}).Build()
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
	i := NewInsert("test_table")
	v1 := KValues{"manager": "123456", "perfer": 2}
	v2 := KValues{"manager": "qwertyu", "perfer": 6}

	query, arg := i.Values(v1).Build()
	fmt.Println(query, "-", arg)

	query, arg = i.Values(v1, v2).Build()
	fmt.Println(query, "-", arg)
}
