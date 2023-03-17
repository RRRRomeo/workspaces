/*
 * @Author: hfouyang hfouyang@quant360.com
 * @Date: 2023-03-16 10:51:18
 * @LastEditors: hfouyang hfouyang@quant360.com
 * @LastEditTime: 2023-03-17 15:04:22
 * @FilePath: /map_chan/btree_idx_demo/mkt_data.go
 * @Description: using btree func to save the mktdata idx;
 */
package btree_idx_demo

import (
	"errors"
	"fmt"
)

type idx_node struct {
	idx uint32
	off int64
	ffp string
	n   *idx_node
}

type idx_links struct {
	h *idx_node
	t *idx_node
	l uint32
}

// var derr = errors.New("default err")

func NewNode(idx uint32, off int64, ffp string, n *idx_node) *idx_node {
	return &idx_node{
		idx: idx,
		off: off,
		ffp: ffp,
		n:   n,
	}
}

func Links() *idx_links {
	return &idx_links{
		h: nil,
		t: nil,
		l: 0,
	}

}

func LinksPushBack(lk *idx_links, ele *idx_node) {
	e := *ele
	if lk == nil {
		panic(errors.New("lk is nil"))
	}

	// ... init first
	if lk.h == nil && lk.t == nil {
		// ...
		lk.h = &e
		lk.t = lk.h
		lk.l++
		return
	}

	// ... normal
	if lk.h != nil && lk.t != nil {
		// ...
		t := lk.t
		t.n = &e
		lk.t = t.n

		lk.l++
	}
}

func LinksPop(lk *idx_links) *idx_node {
	if lk == nil {
		panic(errors.New("lk is nil"))
	}

	if lk.l == 0 {
		panic(errors.New("lk len == 0"))
	}

	if lk.h == nil {
		panic(errors.New("lk.h is nil"))
	}
	v := lk.h

	// ...
	lk.h = lk.h.n
	lk.l--
	return v
}

func LinkkDump(lk *idx_links) {
	if lk == nil {
		panic(errors.New("lk is nil"))
	}

	if lk.l == 0 {
		panic(errors.New("lk len == 0"))
	}

	if lk.h == nil {
		panic(errors.New("lk.h is nil"))
	}

	fmt.Printf("lk len:%d\n", lk.l)
	for it := lk.h; it != nil; {
		fmt.Printf("it p:%p, v:%v\n", it, it)
		it = it.n
	}
}
