/*
 * @Author: hfouyang hfouyang@quant360.com
 * @Date: 2023-03-20 15:02:38
 * @LastEditors: hfouyang hfouyang@quant360.com
 * @LastEditTime: 2023-04-03 10:04:12
 * @FilePath: /githuab.com/RRRRomeo/workspaces/cmd/mkt_idx_cmd/mkt_id_cmd.go
 */
package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"githuab.com/RRRRomeo/workspaces/btree_idx_demo"
	"githuab.com/RRRRomeo/workspaces/internal/mkt_idx_part"
	"githuab.com/RRRRomeo/workspaces/qsorter"
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

func ReadFileToBuildHashNodeIntoChan(fn string, ch chan *qsorter.Qsorter_node) {
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

func SortHashNodeWithChan(ch chan *qsorter.Qsorter_node, q *qsorter.Qsorter) error {
	var err error = nil
	defer wg.Done()

	if ch == nil {
		return errors.New("inner params has nil item")
	}

	for {
		node, ok := <-ch
		if !ok {
			break
		}

		// log.Printf("node:%v\n", node)
		q.Store(node)
		time.Sleep(2 * time.Microsecond)
	}

	wg.Add(1)
	go popAndSaveToFile("./20220701-2460-858.idx", q)
	// midx.WriteToFile("./20220518-456789-10-11-12-38.idx")
	wg.Wait()
	return err
}

func popAndSaveToFile(fp string, q *qsorter.Qsorter) {
	defer wg.Done()
	q.Pop(fp)
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

	// midx.DumpIdx()
	log.Printf("list len:%d\n", midx.GetListLen())
	midx.WriteToFile("./20220518-456789-10-11-12-38.idx")
	return err
}

func buildHashNodeWithTlv(r *mkt_idx_part.Reader, tlv *mkt_idx_part.Tlv[any]) *qsorter.Qsorter_node {
	var bidIdx uint32 = 0
	// var dat int32
	var tim int32 = 0
	var instrId int32
	// typ := btree_idx_demo.ETICK
	off := r.Off

	switch tlv.Header.DataTyp {
	case mkt_idx_part.TICK_TYPE:
		tick := tlv.Data.(mkt_idx_part.MdsL2Trade)
		bidIdx = tick.ApplSeqNum
		tim = tick.TransactTime
		instrId = tick.InstrId
		off -= mkt_idx_part.TICKLEN

	case mkt_idx_part.ORDER_TYPE:
		ord := tlv.Data.(mkt_idx_part.MdsL2Order)
		bidIdx = ord.ApplSeqNum
		tim = ord.TransactTime
		instrId = ord.InstrId
		off -= mkt_idx_part.ORDLEN

	case mkt_idx_part.SNAPSHOT_TYPE:
		ss := tlv.Data.(mkt_idx_part.MdsMktSZL2Snapshot)
		tim = ss.Header.UpdateTime
		instrId = ss.Header.InstrId
		off -= mkt_idx_part.SNAPLEN

	}
	return &qsorter.Qsorter_node{
		Qsorter_write_node: qsorter.Qsorter_write_node{
			SZInStrId: instrId,
			Off:       off,
		},
		BidIdx: bidIdx,
		Tim:    tim,
		N:      nil,
	}

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

func main() {
	runtime.GOMAXPROCS(4)
	// mkt_idx := mkt_idx_part.NewMktIdx()
	sorter := qsorter.NewQSorter()
	ch := make(chan *qsorter.Qsorter_node, 1000)
	wg.Add(6)
	i = 5

	// =================================================================================================/
	// ./20220518-456789-10-11-12-38.idx
	// go ReadFileToBuildHashNodeIntoChan("/home/ty/data/221_data/csv_to_bin/20220518/000004.csv.3in1.0", ch)
	// go ReadFileToBuildHashNodeIntoChan("/home/ty/data/221_data/csv_to_bin/20220518/000005.csv.3in1.0", ch)
	// go ReadFileToBuildHashNodeIntoChan("/home/ty/data/221_data/csv_to_bin/20220518/000006.csv.3in1.0", ch)
	// go ReadFileToBuildHashNodeIntoChan("/home/ty/data/221_data/csv_to_bin/20220518/000007.csv.3in1.0", ch)
	// go ReadFileToBuildHashNodeIntoChan("/home/ty/data/221_data/csv_to_bin/20220518/000008.csv.3in1.0", ch)
	// go ReadFileToBuildHashNodeIntoChan("/home/ty/data/221_data/csv_to_bin/20220518/000009.csv.3in1.0", ch)
	// go ReadFileToBuildHashNodeIntoChan("/home/ty/data/221_data/csv_to_bin/20220518/000010.csv.3in1.0", ch)
	// go ReadFileToBuildHashNodeIntoChan("/home/ty/data/221_data/csv_to_bin/20220518/000011.csv.3in1.0", ch)
	// go ReadFileToBuildHashNodeIntoChan("/home/ty/data/221_data/csv_to_bin/20220518/000012.csv.3in1.0", ch)
	// go ReadFileToBuildHashNodeIntoChan("/home/ty/data/221_data/csv_to_bin/20220518/000038.csv.3in1.0", ch)

	// =================================================================================================/

	go ReadFileToBuildHashNodeIntoChan("/home/ty/data/221_data/csv_to_bin/20220701/002594.csv.3in1.0", ch)
	go ReadFileToBuildHashNodeIntoChan("/home/ty/data/221_data/csv_to_bin/20220701/300750.csv.3in1.0", ch)
	go ReadFileToBuildHashNodeIntoChan("/home/ty/data/221_data/csv_to_bin/20220701/002460.csv.3in1.0", ch)
	go ReadFileToBuildHashNodeIntoChan("/home/ty/data/221_data/csv_to_bin/20220701/000858.csv.3in1.0", ch)
	go ReadFileToBuildHashNodeIntoChan("/home/ty/data/221_data/csv_to_bin/20220701/300390.csv.3in1.0", ch)

	time.Sleep(10 * time.Millisecond)
	go SortHashNodeWithChan(ch, sorter)
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

func genStkFd(day int32, stks []string) []string {
	tmp := make([]string, 0)
	for _, stk := range stks {
		fp := fmt.Sprintf("/home/ty/data/221_data/csv_to_bin/%d/%s.csv.3in1.0", day, stk)
		tmp = append(tmp, fp)
	}
	return tmp
}

func GoGenIdxWithDayAndStks(day int32, stks []string, status qsorter.Status) {
	ch := make(chan *btree_idx_demo.Idx_node, 1000)
	if len(stks) == 1 {
		RunWithGoRoutine(4, stks[0], ch)
		status("ok...\n")
		return
	}

	fps := genStkFd(day, stks)
	goReadFilesToBuildNodeIntoChan(fd, stks, ch)

}
