/*
 * @Author: hfouyang hfouyang@quant360.com
 * @Date: 2023-03-16 10:51:18
 * @LastEditors: hfouyang hfouyang@quant360.com
 * @LastEditTime: 2023-03-27 14:33:50
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
