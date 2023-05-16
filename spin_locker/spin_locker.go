/*
 * @Author: hfouyang hfouyang@quant360.com
 * @Date: 2023-03-14 10:01:34
 * @LastEditors: hfouyang hfouyang@quant360.com
 * @LastEditTime: 2023-03-20 14:54:35
 * @FilePath: /githuab.com/RRRRomeo/workspaces/spin_locker/spin_locker.go
 */
package spin_locker

import (
	"runtime"
	"sync"
	"sync/atomic"
)

type spinLocker uint32
type Spin spinLocker

var maxBackoff = 16

func (s *spinLocker) Lock() {
	backoff := 1
	for !atomic.CompareAndSwapUint32((*uint32)(s), 0, 1) {
		for i := 0; i < backoff; i++ {
			runtime.Gosched()
		}

		if backoff < maxBackoff {
			backoff <<= 1
		}
	}
}

func (s *spinLocker) Unlock() {
	atomic.StoreUint32((*uint32)(s), 0)
}

func NewSpinLocker() sync.Locker {
	return new(spinLocker)
}
