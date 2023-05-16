/*
 * @Author: hfouyang hfouyang@quant360.com
 * @Date: 2023-03-17 10:45:46
 * @LastEditors: hfouyang hfouyang@quant360.com
 * @LastEditTime: 2023-03-30 13:26:03
 * @FilePath: /githuab.com/RRRRomeo/workspaces/cmd/links_cmd/links_cmd.go
 */
package main

import (
	"fmt"

	"githuab.com/RRRRomeo/workspaces/btree_idx_demo"
)

func main() {
	l := btree_idx_demo.Links()
	fmt.Printf("links:%v\n", l)

	for i := 0; i < 100; i++ {
		nn := btree_idx_demo.NewNode(uint32(i), 0, 0, 1, 0, 0, nil)
		btree_idx_demo.LinksPushBack(l, nn)
	}

	fmt.Printf("links:%v\n", l)
	// for i := 0; i < 100; i++ {
	// 	h := btree_idx_demo.LinksTop(l)
	// 	fmt.Printf("links.val:%v\n", *h)
	// }
	btree_idx_demo.LinksDump(l)
}
