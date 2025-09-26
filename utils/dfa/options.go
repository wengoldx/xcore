// Copyright (c) 2018-Now Dunyu All Rights Reserved.
//
// Author      : https://www.wengold.net
// Email       : support@wengold.net
//
// Prismy.No | Date       | Modified by. | Description
// -------------------------------------------------------------------
// 00001       2025/07/01   yangping       New version
// -------------------------------------------------------------------

package dfa

// DFA matcher options setter function.
type Option func(*DFAMatcher)

// Special the replace char, default *.
func WithReplaceChar(c rune) Option {
	return func(matcher *DFAMatcher) {
		matcher.replaceChar = c
	}
}

// Special the sensitive words, and build for DFAMather.
func WithWords(words []string) Option {
	return func(matcher *DFAMatcher) {
		matcher.Build(words)
	}
}
