/*
 * @Author: hfouyang hfouyang@quant360.com
 * @Date: 2023-03-13 17:47:17
 * @LastEditors: hfouyang hfouyang@quant360.com
 * @LastEditTime: 2023-03-14 09:22:40
 * @FilePath: /githuab.com/RRRRomeo/workspaces/cmd/new_map_cmd.go
 */
package main

import (
	"githuab.com/RRRRomeo/workspaces/map_go"
)

func main() {
	wg := map_go.GetWg()
	mp := map_go.MallocNewMap(100)
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go map_go.Test_Write(mp, i)
		wg.Add(1)
		go map_go.Test_Read(mp, i)
		// wg.Done()
	}
	wg.Wait()
}
