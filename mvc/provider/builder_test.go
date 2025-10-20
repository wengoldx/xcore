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

	builder := newBuilder()
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

	builder := newBuilder()
	for _, c := range cases {
		param := c.Params.(LikeData)
		rst := builder.FormatLike(param.Field, param.Filter, param.Pattern...)
		if want := c.Want.(string); rst != want {
			t.Fatal("BaseBuilder.FormatLike error > want:", want, "but result is", rst)
		}
	}
}

// TODO
// ...
