/*
 * @Author: hfouyang hfouyang@quant360.com
 * @Date: 2023-03-16 10:51:18
 * @LastEditors: hfouyang hfouyang@quant360.com
 * @LastEditTime: 2023-03-31 13:55:56
 * @FilePath: /map_chan/btree_idx_demo/mkt_data.go
 * @Description: using btree func to save the mktdata idx;
 */
package btree_idx_demo

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"os"
	"sync"
)

const (
	ETICK uint16 = iota
	EORD
	ESNAP
)

type Idx_node struct {
	Write_Node
	BidIdx uint32
	Tim    int32
	N      *Idx_node
}

type Idx_links struct {
	H *Idx_node
	T *Idx_node
	L uint64
}

type Write_Node struct {
	Typ       uint16
	SZInStrId uint16
	Dat       int32
	Off       uint32
}

// var derr = errors.New("default err")

func NewNode(idx uint32, off uint32, date int32, typ uint16, tim int32, inStrId uint16, n *Idx_node) *Idx_node {
	return &Idx_node{
		Write_Node: Write_Node{
			Off:       off,
			Typ:       typ,
			Dat:       date,
			SZInStrId: inStrId,
		},
		BidIdx: idx,
		Tim:    tim,
		N:      n,
	}
}

func Links() *Idx_links {
	return &Idx_links{
		H: nil,
		T: nil,
		L: 0,
	}

}

func LinksPushBack(lk *Idx_links, ele *Idx_node) {
	e := *ele
	if lk == nil {
		panic(errors.New("lk is nil"))
	}

	// ... init first
	if lk.H == nil && lk.T == nil {
		// ...
		lk.H = &e
		lk.T = lk.H
		lk.L++
		return
	}

	// ... normal
	if lk.H != nil && lk.T != nil {
		// ...
		t := lk.T
		t.N = &e
		lk.T = t.N

		lk.L++
	}
}

func LinksPop(lk *Idx_links) *Idx_node {
	if lk == nil {
		panic(errors.New("lk is nil"))
	}

	if lk.L == 0 {
		panic(errors.New("lk len == 0"))
	}

	if lk.H == nil {
		panic(errors.New("lk.h is nil"))
	}
	v := lk.H

	// ...
	lk.H = lk.H.N
	lk.L--
	return v
}

func LinksDump(lk *Idx_links) {
	if lk == nil {
		panic(errors.New("lk is nil"))
	}

	if lk.L == 0 {
		panic(errors.New("lk len == 0"))
	}

	if lk.H == nil {
		panic(errors.New("lk.h is nil"))
	}

	fmt.Printf("lk len:%d\n", lk.L)
	for it := lk.H; it != nil; {
		log.Printf("it p:%p, v:%v\n", it, it)
		it = it.N
	}
}

func CompareId(p1 *Idx_node, p2 *Idx_node) *Idx_node {
	if p1 == nil || p2 == nil {
		return nil
	}
	if p1.BidIdx >= p2.BidIdx {
		return p2
	}
	return p1
}

func WriteListIntoFile(lk *Idx_links, fp string) error {
	if lk == nil {
		panic(errors.New("inner links is nil"))
	}

	fd, err := os.OpenFile(fp, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0755)
	if err != nil {
		return err
	}
	defer fd.Close()

	writer := bufio.NewWriter(fd)
	for it := lk.H; it != nil; {
		buf := &bytes.Buffer{}
		wt := Write_Node{
			Typ:       it.Typ,
			SZInStrId: it.SZInStrId,
			Dat:       it.Dat,
			Off:       it.Off,
		}

		err := binary.Write(buf, binary.LittleEndian, wt)
		if err != nil {
			log.Printf("bin write err:%s\n", err)
			break
		}

		_, err = writer.Write(buf.Bytes())

		if err != nil {
			log.Printf("文件创建并写入失败！错误：%s\n", err)
			break
		}
		writer.Flush()
		it = it.N
	}
	return nil
}

func (lk *Idx_links) GetLen() uint64 {
	return lk.L
}

type Hash_Node struct {
	// ... key ==> BidIdx
	Write_Node
	BidIdx uint32
	Tim    int32
	N      *Hash_Node
}

type HashTable struct {
	sync.Mutex
	l  uint32
	tb map[uint32]*Hash_Node
}

func NewHashNode(idx uint32, off uint32, date int32, typ uint16, tim int32, inStrId uint16) *Hash_Node {
	return &Hash_Node{
		Write_Node: Write_Node{
			Typ:       typ,
			SZInStrId: inStrId,
			Dat:       date,
			Off:       off,
		},
		BidIdx: idx,
		Tim:    tim,
		N:      nil,
	}
}

func NewHashTable(size uint32) *HashTable {
	return &HashTable{
		tb: make(map[uint32]*Hash_Node, size),
		l:  0,
	}
}

func (ht *HashTable) DelHashNode(k uint32) {
	ht.Lock()
	defer ht.Unlock()
	delete(ht.tb, k)
}

func (ht *HashTable) Insert(k uint32, v *Hash_Node) error {
	if v == nil {
		return errors.New("val is nil")
	}
	ht.Lock()
	defer ht.Unlock()

	if k == 0 { // ... k == 0 ==> this value is snapshot
		// binary search
		idx, err := binSearch(ht, v.Tim)
		if err != nil {
			return err
		}
		ht.tb[idx].N = v
		ht.l++
		return nil
	}
	ht.tb[k] = v
	ht.l++

	return nil
}

func (ht *HashTable) Get(k uint32) (*Hash_Node, error) {
	ht.Lock()
	defer ht.Unlock()
	v, ok := ht.tb[k]
	if !ok {
		return nil, errors.New("dont get meet val")
	}

	return v, nil
}

func (ht *HashTable) HashDump() {
	var i uint32 = 0
	for ; i < ht.l; i++ {
		v, ok := ht.tb[i]
		if !ok {
			continue
		}
		log.Printf("v:%v\n", v)
		for v.N != nil {
			log.Printf("v.N:%v\n", v.N)
			v.N = v.N.N
		}
	}
}

func binSearch(beSearch *HashTable, k int32) (uint32, error) {
	if beSearch == nil {
		return 0, errors.New("the be searcher is nil")
	}
	// var left, mid, right uint32 = 0, 1, beSearch.l

	// for mid > 0 && mid < beSearch.l {
	// 	mid = (left + right) / 2
	// 	log.Printf("left:%d, right:%d, mid:%d\n", left, right, mid)
	// 	if beSearch.tb[mid].Tim <= k && beSearch.tb[mid+1].Tim > k {
	// 		return mid, nil
	// 	}
	// 	if beSearch.tb[mid].Tim < k {
	// 		left = mid
	// 		continue
	// 	}
	// 	if beSearch.tb[mid].Tim > k {
	// 		right = mid
	// 		continue
	// 	}
	// }
	karr := make([]uint32, 0)
	for k, _ := range beSearch.tb {
		karr = append(karr, k)
	}

	for i := len(karr) - 1; i >= 0; i-- {
		if beSearch.tb[karr[i]].Tim > k {
			continue
		}

		if beSearch.tb[karr[i]].Tim <= k {
			return karr[i], nil
		}
	}

	return 0, errors.New("dont search correct k")
}

func (ht *HashTable) PopHashNodeFromChanAndCompare(ch chan *Hash_Node) error {
	if ch == nil {
		return errors.New("chan is nil")
	}
	item, ok := <-ch
	if !ok {
		return errors.New("get item from chan is fail")
	}

	log.Printf("item:%v\n", item)

	var k uint32 = 0
	if item.Typ != ESNAP {
		k = item.BidIdx
	}

	err := ht.Insert(k, item)
	if err != nil {
		return err
	}
	return nil
}

func (ht *HashTable) GetLen() uint32 {
	return ht.l
}
