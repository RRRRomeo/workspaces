/*
 * @Author: hfouyang hfouyang@quant360.com
 * @Date: 2023-03-30 13:26:39
 * @LastEditors: hfouyang hfouyang@quant360.com
 * @LastEditTime: 2023-03-30 13:37:41
 * @FilePath: /githuab.com/RRRRomeo/workspaces/cmd/using_ldflags_cmd/using_ldflags_cmd.go
 */
package main

import (
	"log"

	testldflags "githuab.com/RRRRomeo/workspaces/test/test_ldflags"
)

func main() {
	log.Printf("version:%s\n", testldflags.Get())
}
