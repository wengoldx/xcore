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

// Node data for construct chain list for DFA algorithm.
type Node struct {
	End  bool           // The flag indicate the node whther end node.
	Next map[rune]*Node // Point the next node of DFA nodes chain.
}

// Add the given word to nodes chain start current node.
func (n *Node) AddWord(word string) {
	node, chars := n, []rune(word)
	for index := range chars {
		node = node.AddChild(chars[index])
	}
	node.End = true
}

// Add the next unexist node and return it, or return the exit one.
func (n *Node) AddChild(c rune) *Node {
	if n.Next == nil {
		n.Next = make(map[rune]*Node)
	}

	if next, ok := n.Next[c]; ok {
		return next
	} else {
		next = &Node{End: false, Next: nil}
		n.Next[c] = next
		return next
	}
}

// Find child node from current node, return nil if unexist.
func (n *Node) FindChild(c rune) *Node {
	if n.Next != nil {
		if next, ok := n.Next[c]; ok {
			return next
		}
	}
	return nil
}
