/*
 * @Author: hfouyang hfouyang@quant360.com
 * @Date: 2023-03-15 17:21:32
 * @LastEditors: hfouyang hfouyang@quant360.com
 * @LastEditTime: 2023-03-16 10:45:10
 * @FilePath: /githuab.com/RRRRomeo/workspaces/cmd/spinlocker_cmd/spinlocker_cmd.go
 */
package main

import (
	"log"
	"sync"

	"githuab.com/RRRRomeo/workspaces/map_go"
	"githuab.com/RRRRomeo/workspaces/spin_locker"
)

type TestStruct struct {
	spnlck sync.Locker
	mp     map[int]string
}

var wg sync.WaitGroup

func NewTestStruct() *TestStruct {
	return &TestStruct{
		spnlck: spin_locker.NewSpinLocker(),
		mp:     make(map[int]string),
	}
}

func Write(t *TestStruct, k int, v string) bool {
	t.spnlck.Lock()
	defer t.spnlck.Unlock()
	defer wg.Done()

	t.mp[k] = v
	log.Printf("w k:%d v:%s\n", k, v)
	return true
}

func Read(t *TestStruct, k int) string {
	t.spnlck.Lock()
	defer t.spnlck.Unlock()
	defer wg.Done()

	v, ok := t.mp[k]
	if ok {
		log.Printf("r k:%d v:%s\n", k, v)
		return v
	}

	return ""
}

func main() {
	wg = sync.WaitGroup{}
	t := NewTestStruct()
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go Write(t, i, map_go.Test_RandomString(12))
		wg.Add(1)
		go Read(t, i)
	}
	wg.Wait()

}
