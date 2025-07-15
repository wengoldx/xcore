// Copyright (c) 2018-Now Dunyu All Rights Reserved.
//
// Author      : https://www.wengold.net
// Email       : support@wengold.net
//
// Prismy.No | Date       | Modified by. | Description
// -------------------------------------------------------------------
// 00001       2019/05/22   yangping       New version
// -------------------------------------------------------------------

package ut

// Unit test helper options setter.
type Option func(t *tester)

// Specify author keyword, one of WENGOLD-V1.1, WENGOLD-V1.2, WENGOLD-V2.0.
func WithAuthor(author string) Option {
	return func(t *tester) { t.author = author }
}

// Specify api url to get token from server.
func WithTokenApi(api string) Option {
	return func(t *tester) { t.tokenApi = api }
}
