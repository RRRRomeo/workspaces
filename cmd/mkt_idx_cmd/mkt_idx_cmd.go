/*
 * @Author: hfouyang hfouyang@quant360.com
 * @Date: 2023-03-20 15:02:38
 * @LastEditors: hfouyang hfouyang@quant360.com
 * @LastEditTime: 2023-04-03 10:04:12
 * @FilePath: /map_chan/cmd/mkt_idx_cmd/mkt_id_cmd.go
 */
package main

import (
	"errors"
	"io"
	"log"
	"map_chan/btree_idx_demo"
	"map_chan/internal/mkt_idx_part"
	"os"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

var wg sync.WaitGroup
var i int32

// func testSndToChan(r *mkt_idx_part.Reader, ch chan mkt_idx_part.Tlv[any]) {
// 	for {
// 		e := r.ReadToChan(ch)
// 		if e != nil {
// 			close(ch)
// 			break
// 		}
// 	}
// 	wg.Done()
// }

// func testRedFromChan(r *mkt_idx_part.Reader, ch chan mkt_idx_part.Tlv[any]) {
// 	for {
// 		tlv, ok := <-ch
// 		if !ok {
// 			// log.Printf("read from ch fail\n")
// 			// close(ch)
// 			break
// 		}
// 		DumpTlv(tlv)
// 	}
// }

func ReadFileToBuildNodeIntoChan(fn string, ch chan *btree_idx_demo.Idx_node) {
	// log.Printf("start read file to tlv")
	r := mkt_idx_part.NewReader(1, fn)
	for {
		tlv := &mkt_idx_part.Tlv[any]{}
		e := r.ReadTo(tlv)
		if e != nil {
			if e == io.EOF {
				r.ReleaseReader()
				// close(ch)
			}
			log.Printf("read err:%s\n", e)
			break
		}
		node := buildItemWithTlv(r, tlv)
		mkt_idx_part.PushItemIntoChan(ch, node)
	}
	atomic.AddInt32(&i, -1)
	if i == 0 {
		close(ch)
	}
	// wg.Done()
}

func ReadFileToBuildHashNodeIntoChan(fn string, ch chan *btree_idx_demo.Hash_Node) {
	// log.Printf("start read file to tlv")
	r := mkt_idx_part.NewReader(1, fn)
	for {
		tlv := &mkt_idx_part.Tlv[any]{}
		e := r.ReadTo(tlv)
		if e != nil {
			if e == io.EOF {
				r.ReleaseReader()
				// close(ch)
			}
			log.Printf("read err:%s\n", e)
			break
		}
		hashNode := buildHashNodeWithTlv(r, tlv)
		mkt_idx_part.PushHashNodeIntoChan(ch, hashNode)
	}
	atomic.AddInt32(&i, -1)
	if i == 0 {
		close(ch)
	}
	wg.Done()
}

func SortHashNodeWithChan(ch chan *btree_idx_demo.Hash_Node, ht *btree_idx_demo.HashTable) error {
	var err error = nil
	defer wg.Done()

	if ch == nil || ht == nil {
		return errors.New("inner params has nil item")
	}

	for {
		// log.Printf("im listening....\n")
		err = ht.PopHashNodeFromChanAndCompare(ch)
		if err != nil {
			// close(ch)
			break
		}
	}

	ht.HashDump()
	log.Printf("table len:%d\n", ht.GetLen())
	// midx.WriteToFile("./20220701.idx")
	return err
}

func SortItemWithChan(ch chan *btree_idx_demo.Idx_node, midx *mkt_idx_part.Mkt_idx) error {
	var err error = nil
	defer wg.Done()

	if ch == nil || midx == nil {
		return errors.New("inner params has nil item")
	}

	for {
		// log.Printf("im listening....\n")
		err = midx.PopItemFromChanAndCompare(ch)
		if err != nil {
			// close(ch)
			break
		}
	}

	midx.DumpIdx()
	log.Printf("list len:%d\n", midx.GetListLen())
	// midx.WriteToFile("./20220701.idx")
	return err
}

func buildHashNodeWithTlv(r *mkt_idx_part.Reader, tlv *mkt_idx_part.Tlv[any]) *btree_idx_demo.Hash_Node {
	var bidIdx uint32
	var dat int32
	var tim int32
	var instrId uint16
	typ := btree_idx_demo.ETICK
	off := r.Off

	switch tlv.Header.DataTyp {
	case mkt_idx_part.TICK_TYPE:
		tick := tlv.Data.(mkt_idx_part.MdsL2Trade)
		bidIdx = tick.ApplSeqNum
		dat = tick.TradeDate
		instrId = uint16(tick.InstrId)
		off -= mkt_idx_part.TICKLEN
	case mkt_idx_part.ORDER_TYPE:
		ord := tlv.Data.(mkt_idx_part.MdsL2Order)
		bidIdx = ord.ApplSeqNum
		dat = ord.TradeDate
		instrId = uint16(ord.InstrId)
		off -= mkt_idx_part.ORDLEN
		typ = btree_idx_demo.EORD
	case mkt_idx_part.SNAPSHOT_TYPE:
		ss := tlv.Data.(mkt_idx_part.MdsMktSZL2Snapshot)
		tim = ss.Header.UpdateTime
		dat = ss.Header.TradeDate
		instrId = uint16(ss.Header.InstrId)
		off -= mkt_idx_part.SNAPLEN
		typ = btree_idx_demo.ESNAP
	}
	return btree_idx_demo.NewHashNode(bidIdx, off, dat, typ, tim, instrId)

}

func buildItemWithTlv(r *mkt_idx_part.Reader, tlv *mkt_idx_part.Tlv[any]) *btree_idx_demo.Idx_node {
	var bidIdx uint32
	var dat int32
	var tim int32
	var instrId uint16
	typ := btree_idx_demo.ETICK
	off := r.Off

	switch tlv.Header.DataTyp {
	case mkt_idx_part.TICK_TYPE:
		tick := tlv.Data.(mkt_idx_part.MdsL2Trade)
		bidIdx = tick.ApplSeqNum
		dat = tick.TradeDate
		instrId = uint16(tick.InstrId)
		off -= mkt_idx_part.TICKLEN
	case mkt_idx_part.ORDER_TYPE:
		ord := tlv.Data.(mkt_idx_part.MdsL2Order)
		bidIdx = ord.ApplSeqNum
		dat = ord.TradeDate
		instrId = uint16(ord.InstrId)
		off -= mkt_idx_part.ORDLEN
		typ = btree_idx_demo.EORD
	case mkt_idx_part.SNAPSHOT_TYPE:
		ss := tlv.Data.(mkt_idx_part.MdsMktSZL2Snapshot)
		tim = ss.Header.UpdateTime
		dat = ss.Header.TradeDate
		instrId = uint16(ss.Header.InstrId)
		off -= mkt_idx_part.SNAPLEN
		typ = btree_idx_demo.ESNAP
	}
	return btree_idx_demo.NewNode(uint32(bidIdx), off, dat, typ, tim, uint16(instrId), nil)

}

// func DumpTlv(tlv mkt_idx_part.Tlv[any]) {

// 	switch tlv.Header.DataTyp {
// 	case mkt_idx_part.TICK_TYPE:
// 		tick := tlv.Data.(mkt_idx_part.MdsL2Trade)
// 		log.Printf("tick:%v\n", tick)
// 	case mkt_idx_part.ORDER_TYPE:
// 		ord := tlv.Data.(mkt_idx_part.MdsL2Order)
// 		log.Printf("ord:%v\n", ord)
// 	case mkt_idx_part.SNAPSHOT_TYPE:
// 		ss := tlv.Data.(mkt_idx_part.MdsMktSZL2Snapshot)
// 		log.Printf("snapshot:%v\n", ss)
// 	}
// }

func main() {
	runtime.GOMAXPROCS(4)
	mkt_idx := mkt_idx_part.NewMktIdx()
	// ht := btree_idx_demo.NewHashTable(100000000)
	ch := make(chan *btree_idx_demo.Idx_node, 1000)
	// wg.Add(11)
	// i = 10
	// go ReadFileToBuildNodeIntoChan("/home/ty/data/221_data/csv_to_bin/20220701/000004.csv.3in1.0", ch)
	// go ReadFileToBuildNodeIntoChan("/home/ty/data/221_data/csv_to_bin/20220701/000005.csv.3in1.0", ch)
	// go ReadFileToBuildNodeIntoChan("/home/ty/data/221_data/csv_to_bin/20220701/000038.csv.3in1.0", ch)
	// go ReadFileToBuildNodeIntoChan("/home/ty/data/221_data/csv_to_bin/20220701/000006.csv.3in1.0", ch)
	// go ReadFileToBuildNodeIntoChan("/home/ty/data/221_data/csv_to_bin/20220701/000007.csv.3in1.0", ch)
	// go ReadFileToBuildNodeIntoChan("/home/ty/data/221_data/csv_to_bin/20220701/000008.csv.3in1.0", ch)
	// go ReadFileToBuildNodeIntoChan("/home/ty/data/221_data/csv_to_bin/20220701/000009.csv.3in1.0", ch)
	// go ReadFileToBuildNodeIntoChan("/home/ty/data/221_data/csv_to_bin/20220701/000010.csv.3in1.0", ch)
	// go ReadFileToBuildNodeIntoChan("/home/ty/data/221_data/csv_to_bin/20220701/000011.csv.3in1.0", ch)
	// go ReadFileToBuildNodeIntoChan("/home/ty/data/221_data/csv_to_bin/20220701/000012.csv.3in1.0", ch)
	// go ReadFileToBuildNodeIntoChan(r4, ch)
	// fd := "/home/ty/data/221_data/csv_to_bin/20220701/"
	// files, err := os.ReadDir(fd)
	// if err != nil {
	// 	return
	// }
	// len_files := len(files)
	// log.Printf("len:%d\n", len_files)
	// for idx := 0; idx < 10; idx += 5 {
	// 	wg.Add(1)
	// 	i += 1
	// 	log.Printf("%s\n", fd+files[idx].Name())
	// 	go ReadFileToBuildNodeIntoChan(fd+files[idx].Name(), ch)
	// 	go ReadFileToBuildNodeIntoChan(fd+files[idx+1].Name(), ch)
	// 	go ReadFileToBuildNodeIntoChan(fd+files[idx+2].Name(), ch)
	// 	go ReadFileToBuildNodeIntoChan(fd+files[idx+3].Name(), ch)
	// 	go ReadFileToBuildNodeIntoChan(fd+files[idx+4].Name(), ch)
	// 	// log.Printf("file:%s\n", files[idx].Name())
	// }
	RunWithGoRoutine(8, "/home/ty/data/221_data/csv_to_bin/20220701/", ch)
	time.Sleep(10 * time.Millisecond)
	wg.Add(1)
	go SortItemWithChan(ch, mkt_idx)
	// mkt_idx.DumpIdx()
	wg.Wait()

}

func RunWithGoRoutine(routines uint8, fd string, ch chan *btree_idx_demo.Idx_node) error {
	files, err := os.ReadDir(fd)
	if err != nil {
		return err
	}
	names := make([]string, 0)
	for _, f := range files {
		names = append(names, f.Name())
	}

	nameroutine := make([][]string, routines)
	step := len(names) / int(routines)
	// log.Printf("step:%d\n", step)
	for i := 0; i < int(routines); i++ {
		nameroutine[i] = names[i*step : (i+1)*step]
	}
	if step*int(routines) != len(names) {
		nameroutine[7] = append(nameroutine[7], names[(8*step)+1:]...)
	}

	for i := 0; i < int(routines); i++ {
		wg.Add(1)
		go goReadFilesToBuildNodeIntoChan(fd, nameroutine[i], ch)
	}
	// log.Printf("names len:%d, name:%v\n", len(names), names)
	// log.Printf("======================================\n")
	// log.Printf("nameroutins len:%d, nameroutine:%v\n", len(nameroutine), nameroutine)

	return nil
}

func goReadFilesToBuildNodeIntoChan(fd string, files []string, ch chan *btree_idx_demo.Idx_node) {
	for _, f := range files {
		// log.Printf("fd+f:%s\n", fd+f)
		ReadFileToBuildNodeIntoChan(fd+f, ch)
	}
	wg.Done()
}
