/*
 * @Author: hfouyang hfouyang@quant360.com
 * @Date: 2023-03-20 14:42:14
 * @LastEditors: hfouyang hfouyang@quant360.com
 * @LastEditTime: 2023-03-30 18:05:33
 * @FilePath: /map_chan/internal/mkt_idx_part/mkt_idx.go
 */
package mkt_idx_part

import (
	"errors"
	"log"
	"map_chan/btree_idx_demo"
	"sync"
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

	log.Printf("item:%v\n", item)

	i.Lock()
	defer i.Unlock()

	e := CompareIt(i.list, item)
	if e != nil {
		return e
	}
	return nil
}

// func CompareIt2(l *btree_idx_demo.Idx_links, n *btree_idx_demo.Idx_node) error {
// 	if l == nil || n == nil {
// 		return errors.New("inner list or node is nil")
// 	}
// 	typ := n.Typ
// 	cur := l.H

// 	// 排序
// 	sortLinks(cur, typ)

// 	// 二分查找
// 	idx := searchForInsertPos(cur, n, typ)
// 	if idx == -1 {
// 		btree_idx_demo.LinksPushBack(l, n)
// 	} else {
// 		insertNode(cur, idx, n)
// 		l.L++
// 	}

// 	return nil

// }

// func sortLinks(h *btree_idx_demo.Idx_node, typ uint16) {
// 	// 将链表节点按照 typ 所定义的比较规则排序
// 	switch typ {
// 	case btree_idx_demo.ESNAP:
// 		sort.SliceStable(h, func(i, j int) bool {
// 			return h[i].Tim < h[j].Tim
// 		})
// 	default:
// 		sort.SliceStable(h, func(i, j int) bool {
// 			return h[i].BidIdx < h[j].BidIdx
// 		})
// 	}
// }

// func searchForInsertPos(h *btree_idx_demo.Idx_links, n *btree_idx_demo.Idx_node, typ uint16) int {
// 	// 二分查找插入位置
// 	switch typ {
// 	case btree_idx_demo.ESNAP:
// 		for low, high := 0, len(h)-1; low <= high; {
// 			mid := (low + high) / 2
// 			if h[mid].Tim == n.Tim {
// 				return mid
// 			} else if h[mid].Tim < n.Tim {
// 				low = mid + 1
// 			} else {
// 				high = mid - 1
// 			}
// 		}
// 		return low // 找到插入位置，返回索引
// 	default:
// 		for low, high := 0, len(h)-1; low <= high; {
// 			mid := (low + high) / 2
// 			if h[mid].BidIdx == n.BidIdx {
// 				return mid
// 			} else if h[mid].BidIdx < n.BidIdx {
// 				low = mid + 1
// 			} else {
// 				high = mid - 1
// 			}
// 		}
// 		return low // 找到插入位置，返回索引
// 	}
// 	return -1 // 错误情况
// }

// func insertNode(h []*btree_idx_demo.Idx_links, idx int, n *btree_idx_demo.Idx_node) {
// 	// 在 idx 处插入新节点 n
// 	// 注意需要移动后面的节点
// 	for i := len(h) - 1; i >= idx; i-- {
// 		h[i+1] = h[i]
// 	}
// 	h[idx] = &btree_idx_demo.Idx_links{
// 		BidIdx: n.BidIdx,
// 		Tim:    n.Tim,
// 		N:      n.N,
// 	}
// }
