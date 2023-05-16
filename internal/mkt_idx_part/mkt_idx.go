/*
 * @Author: hfouyang hfouyang@quant360.com
 * @Date: 2023-03-20 14:42:14
 * @LastEditors: hfouyang hfouyang@quant360.com
 * @LastEditTime: 2023-03-31 13:57:53
 * @FilePath: /githuab.com/RRRRomeo/workspaces/internal/mkt_idx_part/mkt_idx.go
 */
package mkt_idx_part

import (
	"errors"
	"log"
	"sync"

	"githuab.com/RRRRomeo/workspaces/btree_idx_demo"
	"githuab.com/RRRRomeo/workspaces/qsorter"
)

type idx_manager interface {
	PopIdx() *btree_idx_demo.Idx_node
	PushBack(*btree_idx_demo.Idx_node)
}

type Mkt_idx struct {
	sync.Mutex
	list *btree_idx_demo.Idx_links
}

func NewMktIdx() *Mkt_idx {
	return &Mkt_idx{
		list: btree_idx_demo.Links(),
	}
}

func (i *Mkt_idx) PushBack(n *btree_idx_demo.Idx_node) {
	i.Lock()
	btree_idx_demo.LinksPushBack(i.list, n)
	i.Unlock()
}

func (i *Mkt_idx) PopIdx() *btree_idx_demo.Idx_node {
	i.Lock()
	defer i.Unlock()
	return btree_idx_demo.LinksPop(i.list)
}

func (i *Mkt_idx) DumpIdx() {
	i.Lock()
	defer i.Unlock()

	btree_idx_demo.LinksDump(i.list)
}

func (i *Mkt_idx) WriteToFile(fp string) {
	i.Lock()
	defer i.Unlock()

	err := btree_idx_demo.WriteListIntoFile(i.list, fp)
	if err != nil {
		log.Panicf("err:%s\n", err)
	}
}

func (i *Mkt_idx) GetListLen() uint64 {
	return i.list.L
}
func PushHashNodeIntoChan(ch chan *qsorter.Qsorter_node, item *qsorter.Qsorter_node) {
	if item == nil {
		return
	}

	ch <- item
}

func PushItemIntoChan(ch chan *btree_idx_demo.Idx_node, item *btree_idx_demo.Idx_node) {
	if item == nil {
		return
	}

	ch <- item
}

func CompareIt(l *btree_idx_demo.Idx_links, n *btree_idx_demo.Idx_node) error {
	if l == nil || n == nil {
		return errors.New("inner list or node is nil")
	}
	typ := n.Typ
	cur := l.H

	if typ == btree_idx_demo.ESNAP {
		for cur != nil {
			// cur.Time < n.Time ===> cur ++
			if cur.Tim < n.Tim {
				if cur.N != nil && cur.N.Tim >= n.Tim {
					// insert cur station
					n.N = cur.N
					cur.N = n
					l.L++
					return nil
				}

				cur = cur.N
				continue
			}

			n.N = cur
			cur = n
			l.L++
			return nil
		}
		btree_idx_demo.LinksPushBack(l, n)
		return nil
	}

	// ... typ == ENOTSNAP
	for cur != nil {
		if cur.BidIdx < n.BidIdx {
			if cur.N != nil && cur.N.BidIdx >= n.BidIdx {
				n.N = cur.N
				cur.N = n
				l.L++
				return nil
			}
			cur = cur.N
			continue
		}
		n.N = cur
		cur = n
		l.L++
		return nil
	}

	btree_idx_demo.LinksPushBack(l, n)
	return nil
}

func (i *Mkt_idx) PopItemFromChanAndCompare(ch chan *btree_idx_demo.Idx_node) error {
	if ch == nil {
		return errors.New("chan is nil")
	}
	item, ok := <-ch
	if !ok {
		return errors.New("get item from chan is fail")
	}

	// log.Printf("item:%v\n", item)

	i.Lock()
	defer i.Unlock()

	e := CompareIt(i.list, item)
	if e != nil {
		return e
	}
	return nil
}
