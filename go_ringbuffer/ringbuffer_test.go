package gingbuffer

import (
	"testing"

	"gitlab.quant360.com/hfouyang/qlog/log"
)

func Test_gingbuff_Put(t *testing.T) {
	// gbuf, err := NewBuffer(1024 * 4)
	// if err != nil {
	// 	log.Debugf("NewGingBuffer fail:%s\n", err)
	// 	return
	// }
	Init()

	testData := "test put the data 0123456789\t\n"
	putBuf := []byte(testData)
	puter := NewPuter()
	puter.Put(putBuf)
	puter.Put(putBuf)
	puter.Put(putBuf)
	puter.Put(putBuf)

	geter := NewGeter()
	geter2 := NewGeter()
	reads, err := geter.Get()
	if err != nil {
		log.Errf("geter get data fail:%s\n", err)
		return
	}
	str := string(reads)
	log.Debugf("read bytes:%s\n", str)
	reads2, err := geter.Get()
	if err != nil {
		log.Errf("geter get data fail:%s\n", err)
		return
	}
	str2 := string(reads2)
	log.Debugf("read bytes:%s\n", str2)

	log.Debugf("the puter id:%d, the geter id:%d, the geter2id:%d\n", puter.GetId(), geter.GetId(), geter2.GetId())

}
