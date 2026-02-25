// Copyright (c) 2018-Now Dunyu All Rights Reserved.
//
// Author      : https://www.wengold.net
// Email       : support@wengold.net
//
// Prismy.No | Date       | Modified by. | Description
// -------------------------------------------------------------------
// 00001       2019/05/22   yangping       New version
// -------------------------------------------------------------------

package utils

import (
	"bytes"
	"strings"

	"github.com/russross/blackfriday/v2"
	"github.com/wengoldx/xcore/logger"
	"golang.org/x/net/html"
)

// Parse markdown content to text content.
func MarkdownToText(src string) (string, error) {
	if src = strings.TrimSpace(src); src == "" {
		return "", nil
	}

	// markdown content to html content.
	htmlstring := string(blackfriday.Run([]byte(src)))

	// html content to pick text strings.
	doc, err := html.Parse(strings.NewReader(htmlstring))
	if err != nil {
		logger.E("Parsing markdown to html, err:", err)
		return "", err
	}

	var text bytes.Buffer
	walkHtmlNode(doc, &text, func(n *html.Node) bool {
		if n.Type == html.TextNode {
			text.WriteString(n.Data)
		}
		return true
	})
	return text.String(), nil
}

// Walk html content nodes to filter out all text typed nodes.
func walkHtmlNode(n *html.Node, out *bytes.Buffer, fetchFuc func(*html.Node) bool) {
	for c := n.FirstChild; c != nil; {
		if !fetchFuc(c) {
			if c.Type == html.ElementNode && c.Data == "br" {
				out.WriteString("\n")
			}
		} else {
			walkHtmlNode(c, out, fetchFuc)
		}
		c = c.NextSibling
	}
}
