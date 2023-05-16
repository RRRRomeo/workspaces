/*
 * @Author: hfouyang hfouyang@quant360.com
 * @Date: 2023-03-27 17:37:23
 * @LastEditors: hfouyang hfouyang@quant360.com
 * @LastEditTime: 2023-03-27 17:59:13
 * @FilePath: /githuab.com/RRRRomeo/workspaces/test_chan/test_chan.go
 */
package main

import (
	"log"
	"runtime"
	"strconv"
	"strings"
	"sync"

	"gitlab.quant360.com/algo_strategy_group/mds_data_playback/src"
	"gitlab.quant360.com/hsw/sz_match/types"
)

var wg sync.WaitGroup

func Start(fn string, date string) {
	ch := make(chan types.MdData, 1000)
	defer wg.Done()
	go src.Get3In1DataWithChan("/home/ty/data/221_data/csv_to_bin/", fn, date, ch)

	for {
		d, ok := <-ch
		if !ok {
			// close()
			break
		}
		log.Printf("pid:%d start recv data from %s data %v\n", getGoroutineID(), fn, d)
	}

}

func getGoroutineID() int64 {
	buf := make([]byte, 64)
	n := runtime.Stack(buf, false)
	idStr := strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine "))[0]
	id, _ := strconv.ParseInt(idStr, 10, 64)
	return id
}

func main() {
	log.Printf("start\n")
	wg = sync.WaitGroup{}
	log.Printf("Main Goroutine ID:%d\n", getGoroutineID())
	// for i := 1; i < 11; i++ {
	wg.Add(3)
	// fn := fmt.Sprintf("%d", i)
	go Start("000005", "20220701")
	go Start("000038", "20220701")
	go Start("000004", "20220701")
	// }
	wg.Wait()
}
