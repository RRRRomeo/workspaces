/*
 * @Author: hfouyang hfouyang@quant360.com
 * @Date: 2023-03-21 10:00:41
 * @LastEditors: hfouyang hfouyang@quant360.com
 * @LastEditTime: 2023-03-31 17:22:08
 * @FilePath: /map_chan/test/func_test.go
 */
package test

import (
	"log"
	"map_chan/internal/mkt_idx_part"
	"os"
	"sync"
	"testing"
)

var wg sync.WaitGroup

func tetSndToChan(r *mkt_idx_part.Reader, ch chan mkt_idx_part.Tlv[any]) {
	for {
		e := r.ReadToChan(ch)
		if e != nil {
			close(ch)
			break
		}
	}
	wg.Done()
}

func tetRedFromChan(r *mkt_idx_part.Reader, ch chan mkt_idx_part.Tlv[any]) {
	for {
		_, ok := <-ch
		if !ok {
			// log.Printf("read from ch fail\n")
			// close(ch)
			break
		}
		// DumpTlv(tlv)
	}
}

func DumpTlv(tlv mkt_idx_part.Tlv[any]) {

	switch tlv.Header.DataTyp {
	case mkt_idx_part.TICK_TYPE:
		tick := tlv.Data.(mkt_idx_part.MdsL2Trade)
		log.Printf("tick:%v\n", tick)
	case mkt_idx_part.ORDER_TYPE:
		ord := tlv.Data.(mkt_idx_part.MdsL2Order)
		log.Printf("ord:%v\n", ord)
	case mkt_idx_part.SNAPSHOT_TYPE:
		ss := tlv.Data.(mkt_idx_part.MdsMktSZL2Snapshot)
		log.Printf("snapshot:%v\n", ss)
	}
}

// func TestRead(t *testing.T) {
// 	// idxs := mkt_idx_part.NewMktIdx()
// 	// wg = *map_go.GetWg()
// 	// for i := 0; i < 100; i++ {
// 	// 	nn := btree_idx_demo.NewNode(uint32(i), 0, map_go.Test_RandomString(8), nil)
// 	// 	idxs.PushBack(nn)
// 	// }
// 	// for i := 0; i < 100; i++ {
// 	// h := idxs.PopIdx()
// 	// fmt.Printf("links.val:%v\n", *h)
// 	// }
// 	// idxs.DumpIdx()

// 	r := mkt_idx_part.NewReader(1, "/home/ty/data/221_data/20220701/000001.csv.3in1.0")
// 	ch := make(chan mkt_idx_part.Tlv[any])
// 	wg.Add(1)
// 	go tetSndToChan(r, ch)
// 	tetRedFromChan(r, ch)
// 	wg.Wait()
// }

func RunWithGoRoutine(routines uint8, fd string) error {
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
	log.Printf("step:%d\n", step)
	for i := 0; i < int(routines); i++ {
		nameroutine[i] = names[i*step : (i+1)*step]
	}
	if step*int(routines) != len(names) {
		nameroutine[7] = append(nameroutine[7], names[(8*step)+1:]...)
	}

	log.Printf("names len:%d, name:%v\n", len(names), names)
	log.Printf("======================================\n")
	log.Printf("nameroutins len:%d, nameroutine:%v\n", len(nameroutine), nameroutine)

	return nil
}

func TestWithGoRoutine(t *testing.T) {
	fd := "/home/ty/data/221_data/csv_to_bin/20220701/"
	RunWithGoRoutine(8, fd)
}
