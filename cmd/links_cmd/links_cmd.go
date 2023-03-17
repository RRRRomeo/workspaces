/*
 * @Author: hfouyang hfouyang@quant360.com
 * @Date: 2023-03-17 10:45:46
 * @LastEditors: hfouyang hfouyang@quant360.com
 * @LastEditTime: 2023-03-17 16:07:30
 * @FilePath: /map_chan/cmd/links_cmd/links_cmd.go
 */
package main

import (
	"fmt"
	"map_chan/btree_idx_demo"
	"map_chan/map_go"
)

func main() {
	l := btree_idx_demo.Links()
	fmt.Printf("links:%v\n", l)

	for i := 0; i < 100; i++ {
		nn := btree_idx_demo.NewNode(uint32(i), 0, map_go.Test_RandomString(8), nil)
		btree_idx_demo.LinksPushBack(l, nn)
	}
	fmt.Printf("links:%v\n", l)
	// for i := 0; i < 100; i++ {
	// 	h := btree_idx_demo.LinksTop(l)
	// 	fmt.Printf("links.val:%v\n", *h)
	// }
	btree_idx_demo.LinkkDump(l)
}
