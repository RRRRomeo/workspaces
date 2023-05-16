/*
 * @Author: hfouyang hfouyang@quant360.com
 * @Date: 2023-03-28 14:40:25
 * @LastEditors: hfouyang hfouyang@quant360.com
 * @LastEditTime: 2023-03-28 14:42:16
 * @FilePath: /githuab.com/RRRRomeo/workspaces/pb_lines/pb_lines.go
 */
package pb_lines

import (
	"time"

	"github.com/cheggaaa/pb/v3"
)

func Lines_start() {
	count := 100
	bar := pb.StartNew(count)

	for i := 0; i < count; i++ {
		time.Sleep(time.Millisecond * 50)
		bar.Increment()
	}

	bar.Finish()
}
