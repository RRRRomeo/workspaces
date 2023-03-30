/*
 * @Author: hfouyang hfouyang@quant360.com
 * @Date: 2023-03-20 15:02:38
 * @LastEditors: hfouyang hfouyang@quant360.com
 * @LastEditTime: 2023-03-29 10:34:11
 * @FilePath: /map_chan/cmd/mkt_idx_cmd/mkt_id_cmd.go
 */
package main

import (
	"errors"
	"io"
	"log"
	"map_chan/btree_idx_demo"
	"map_chan/internal/mkt_idx_part"
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
	wg.Done()
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
	ch := make(chan *btree_idx_demo.Idx_node, 1000)
	wg.Add(11)
	i = 10
	go ReadFileToBuildNodeIntoChan("/home/ty/data/221_data/csv_to_bin/20220701/000004.csv.3in1.0", ch)
	go ReadFileToBuildNodeIntoChan("/home/ty/data/221_data/csv_to_bin/20220701/000005.csv.3in1.0", ch)
	go ReadFileToBuildNodeIntoChan("/home/ty/data/221_data/csv_to_bin/20220701/000038.csv.3in1.0", ch)
	go ReadFileToBuildNodeIntoChan("/home/ty/data/221_data/csv_to_bin/20220701/000006.csv.3in1.0", ch)
	go ReadFileToBuildNodeIntoChan("/home/ty/data/221_data/csv_to_bin/20220701/000007.csv.3in1.0", ch)
	go ReadFileToBuildNodeIntoChan("/home/ty/data/221_data/csv_to_bin/20220701/000008.csv.3in1.0", ch)
	go ReadFileToBuildNodeIntoChan("/home/ty/data/221_data/csv_to_bin/20220701/000009.csv.3in1.0", ch)
	go ReadFileToBuildNodeIntoChan("/home/ty/data/221_data/csv_to_bin/20220701/000010.csv.3in1.0", ch)
	go ReadFileToBuildNodeIntoChan("/home/ty/data/221_data/csv_to_bin/20220701/000011.csv.3in1.0", ch)
	go ReadFileToBuildNodeIntoChan("/home/ty/data/221_data/csv_to_bin/20220701/000012.csv.3in1.0", ch)
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
	time.Sleep(10 * time.Millisecond)
	go SortItemWithChan(ch, mkt_idx)
	// mkt_idx.DumpIdx()
	wg.Wait()

}
