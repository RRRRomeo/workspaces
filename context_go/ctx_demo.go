/*
 * @Author: hfouyang hfouyang@quant360.com
 * @Date: 2023-03-13 15:22:20
 * @LastEditors: hfouyang hfouyang@quant360.com
 * @LastEditTime: 2023-03-14 09:56:21
 * @FilePath: /map_chan/context_go/ctx_demo.go
 * @Description: learn the ctx point within src code;
 */
package context_go

import (
	"context"
	"fmt"
	"time"
)

/**
 * @description: context module init
 * @return {void}
 */
func Context_TestInit() {
	ctx, close := context.WithTimeout(context.Background(), 1*time.Second)
	defer close()

	for {
		select {
		case <-ctx.Done():
			fmt.Printf("the ctx done\n")
			goto brk
		default:
			fmt.Printf("default case\n")
		}
	}

brk:
	fmt.Printf("brk out;\n")
}
