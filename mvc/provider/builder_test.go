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
)

func TestQueryBuilder(t *testing.T) {
	q := NewQuery("test_table")

	query, arg := q.Outs("f1", "f2", "f3").Build()
	fmt.Println(query, "-", arg)

	query, arg = q.Wheres(Wheres{"w1=?": 1, "w2<>?": 2}).Build()
	fmt.Println(query, "-", arg)

	query, arg = q.WhereIn("in3", q.Int64Anys([]int64{3, 4, 5})).Build()
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
