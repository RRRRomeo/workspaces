/*
 * @Author: hfouyang hfouyang@quant360.com
 * @Date: 2023-03-13 17:00:58
 * @LastEditors: hfouyang hfouyang@quant360.com
 * @LastEditTime: 2023-03-14 09:36:16
 * @FilePath: /githuab.com/RRRRomeo/workspaces/cmd/test.go
 */
package main

import (
	"fmt"
	"time"
)

func Test_Read(m *map[int]int, k int) {
	// for i := 0; i < 1000; i++ {
	// 	go func() {
	_, ok := (*m)[k]
	if ok {
		// ...
		fmt.Printf("its ok\n")
	}
	// 	}()
	// }
}

// err1:
//
//	return
func Test_Write(m *map[int]int, k int, v int) {
	// for i := 0; i < 1000; i++ {
	// 	go func() {
	(*m)[k] = v
	// 	}()
	// }
}

func main() {
	m := make(map[int]int)

	for i := 0; i < 1000; i++ {
		go Test_Write(&m, i, i+1)
		go Test_Read(&m, i)
	}

	time.Sleep(100 * time.Second)
}
