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

import "github.com/wengoldx/xcore/utils"

// A interface for implement DFA matcher.
type Matcher interface {
	Build(words []string)                        // Build DFA nodes tree from given words.
	MatchReplace(text string) ([]string, string) // Filter out the sensitive words and return the replaced text.
	MatchWords(text string) []string             // Filter out the sensitive words.
}

var _ = (*DFAMatcher)(nil)

// DFA matcher for filter sensitive words.
type DFAMatcher struct {
	root        *Node // The root node of sensitive nodes chain.
	replaceChar rune  // The replace char, default *.
}

// Create a new DFAMatcher instance for search sensitive words.
func NewMatcher(options ...Option) Matcher {
	matcher := &DFAMatcher{root: &Node{End: false}, replaceChar: '*'}
	for _, option := range options {
		option(matcher)
	}
	return matcher
}

// Build DFA nodes tree from given words as follow struct:
//
// words = []string{"生产", "生日", "敏感词", ...}
//
//	               +-- [产]
//	              /
//	       +-- [生] -- [日]
//	       |
//	[root] --- [敏] -- [感] -- [词]
//	       |
//	       ...
func (d *DFAMatcher) Build(words []string) {
	for _, word := range words {
		d.root.AddWord(word)
	}
}

// Filter out the sensitive words and return the replaced text.
func (d *DFAMatcher) MatchReplace(text string) ([]string, string) {
	return d.match(text, true)
}

// Filter out the sensitive words.
func (d *DFAMatcher) MatchWords(text string) []string {
	sensitives, _ := d.match(text, false)
	return sensitives
}

// Match the given text to filter out the sensitive words, then return the replaced result text.
func (d *DFAMatcher) match(text string, replace bool) ([]string, string) {
	if d.root == nil {
		return nil, text
	}

	// prepare chars and replace copy.
	chars, replaces := []rune(text), []rune{}
	length := len(chars)
	if replace {
		replaces = make([]rune, length)
		copy(replaces, chars)
	}

	// fetch sensitives words and replace them if request.
	sets := utils.NewSets[string]()
	for i, c := range chars {
		if node := d.root.FindChild(c); node != nil {
			for j := i + 1; j < length && node != nil; j++ {
				if node.End {
					sets.Add(string(chars[i:j]))
					if replace {
						d.replaceRune(replaces, i, j)
					}
				}
				node = node.FindChild(chars[j])
			}

			if node != nil && node.End {
				sets.Add(string(chars[i:length]))
				if replace {
					d.replaceRune(replaces, i, length)
				}
			}
		}
	}
	return sets.Array(), string(replaces)
}

// Using the given char to replace the found sensitive words.
func (d *DFAMatcher) replaceRune(chars []rune, start, end int) {
	for i := start; i < end; i++ {
		chars[i] = d.replaceChar
	}
}
