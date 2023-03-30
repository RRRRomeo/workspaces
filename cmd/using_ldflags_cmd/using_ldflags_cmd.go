/*
 * @Author: hfouyang hfouyang@quant360.com
 * @Date: 2023-03-30 13:26:39
 * @LastEditors: hfouyang hfouyang@quant360.com
 * @LastEditTime: 2023-03-30 13:37:41
 * @FilePath: /map_chan/cmd/using_ldflags_cmd/using_ldflags_cmd.go
 */
package main

import (
	"log"
	testldflags "map_chan/test_ldflags"
)

func main() {
	log.Printf("version:%s\n", testldflags.Get())
}
