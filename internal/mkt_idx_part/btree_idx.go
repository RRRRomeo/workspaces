/*
 * @Author: hfouyang hfouyang@quant360.com
 * @Date: 2023-03-21 14:46:00
 * @LastEditors: hfouyang hfouyang@quant360.com
 * @LastEditTime: 2023-03-23 10:24:44
 * @FilePath: /map_chan/internal/mkt_idx_part/btree_idx.go
 */
package mkt_idx_part

import (
	"log"
)

type tree_node struct {
	isSnap  int8
	applNum uint64
	bcTime  int64
	off     int64
	ffp     string
}

type Idx_tree struct {
	val   tree_node
	left  *Idx_tree
	right *Idx_tree
}

func NewTreeNode(s int8, applId uint64, offset int64, ffn string) *tree_node {
	return &tree_node{
		isSnap:  s,
		applNum: applId,
		off:     offset,
		ffp:     ffn,
	}
}

func NewIdxTreeWith(node *tree_node) *Idx_tree {
	return &Idx_tree{
		val:   *node,
		left:  nil,
		right: nil,
	}
}

func (t *Idx_tree) IdxTree_PushData(d *tree_node) error {
	if d == nil {
		log.Fatalf("inner data idx_node is nil\n")
	}

	cur := t
	for cur != nil {
		if d.isSnap == 1 {
			if d.bcTime >= cur.val.bcTime {
				cur = cur.right
				continue
			}

			cur = cur.left
			continue
		}

		if d.applNum >= cur.val.applNum {
			cur = cur.right
			continue
		}
		cur = cur.left
		continue
	}

	tmp := NewIdxTreeWith(d)
	cur = tmp
	return nil
}
