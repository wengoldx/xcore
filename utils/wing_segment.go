// Copyright (c) 2018-Now Dunyu All Rights Reserved.
//
// Author      : https://www.wengold.net
// Email       : support@wengold.net
//
// Prismy.No | Date       | Modified by. | Description
// -------------------------------------------------------------------
// 00001       2022/06/07   jidi           New version
// -------------------------------------------------------------------

package utils

import (
	"unicode"

	"github.com/go-ego/gse"
)

type WordSegment struct {
	segment gse.Segmenter
}

var GseSegmenter *WordSegment

func SegmentHelper() *WordSegment {
	GseSegmenter = &WordSegment{}
	upperDir := GetUpperFileDir()
	file := upperDir + "/source/s_1.txt," + upperDir + "/source/t_1.txt"
	GseSegmenter.segment.LoadDict(file)
	return GseSegmenter
}

// CutWord split clothes title, save keywords and set keyword weight
func (c *WordSegment) CutWord(params string) []string {
	key_words := c.filterSpaceSymbols(params)
	words := c.segment.CutSearch(key_words, true)
	return words
}

// filterSpaceSymbols filter the spaces and symbols in the string array
func (c *WordSegment) filterSpaceSymbols(key_words string) string {
	words := []rune{}
	for _, v := range key_words {
		if unicode.IsSpace(v) || unicode.IsPunct(v) {
			continue
		}
		words = append(words, v)
	}
	return string(words)
}
