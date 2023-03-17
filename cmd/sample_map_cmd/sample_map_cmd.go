/*
 * @Author: hfouyang hfouyang@quant360.com
 * @Date: 2023-03-13 17:47:30
 * @LastEditors: hfouyang hfouyang@quant360.com
 * @LastEditTime: 2023-03-14 09:32:56
 * @FilePath: /map_chan/cmd/sample_map_cmd.go
 */
package main

import (
	"log"
	"map_chan/map_go"
)

func main() {
	wg := map_go.GetWg()
	log.Printf("main....\n")
	m, err := map_go.SampleMap_Init(100)
	log.Printf("tttt\n")
	if err != nil {
		return
	}
	log.Printf("tttt2:%v\n", m.Dat)
	for i := 0; i < 100; i++ {

		wg.Add(1)
		go map_go.Test_GoRunTime_Write(m, i)
		wg.Add(1)
		go map_go.Test_GoRunTime_Read(m, i)
	}
	wg.Wait()
	log.Printf("main....end\n")

}
