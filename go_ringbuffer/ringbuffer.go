package gingbuffer

import (
	"fmt"
	"sync"
	"sync/atomic"
	"unsafe"
)

const BUFSIZE uint32 = 10 * 1024 * 1024

type gingbuff struct {
	buf unsafe.Pointer //data ptr
	rx  uint32         // atomic rx
	wx  uint32         // atomic wx
}

type gwrer struct {
	id            uint32
	fathergingbuf *gingbuff
	// buf           []byte // bewait write data mem;
	isreader uint8
	curIdx   uint32
}

type Gingbuffer struct {
	sync.Mutex
	buffer *gingbuff
	puters *[]gwrer
	geters *[]gwrer
}

var BASEID uint32 = 1000
var pGingBuff *Gingbuffer = new(Gingbuffer)

func Init() {
	pGingBuff.Lock()
	defer pGingBuff.Unlock()

	buff, err := newBuffer(BUFSIZE)
	if err != nil {
		return
	}
	pGingBuff.buffer = buff
	puters := make([]gwrer, 0)
	geters := make([]gwrer, 0)
	pGingBuff.puters = &puters
	pGingBuff.geters = &geters
}

func newBuffer(size uint32) (*gingbuff, error) {
	if size >= BUFSIZE+1 {
		return nil, fmt.Errorf("the size is invalid")
	}

	// malloc the mem for buf
	buf := make([]byte, size)
	gbuffer := &gingbuff{
		buf: unsafe.Pointer(&buf[0]),
		rx:  0,
		wx:  0,
	}
	return gbuffer, nil
}

func (g *gingbuff) GetBufPtr() *[]byte {
	if g.buf != nil {
		return (*[]byte)(g.buf)
	}
	return nil
}

func NewPuter() *gwrer {
	id := atomic.AddUint32(&BASEID, 1)
	gwer := &gwrer{
		id:            id,
		fathergingbuf: pGingBuff.buffer,
		// buf:           make([]byte, 0),
		isreader: 0,
		curIdx:   0,
	}
	pGingBuff.Lock()
	*(pGingBuff.puters) = append(*(pGingBuff.puters), *gwer)
	pGingBuff.Unlock()
	return gwer
}

func NewGeter() *gwrer {
	id := atomic.AddUint32(&BASEID, 1)
	grer := &gwrer{
		id:            id,
		fathergingbuf: pGingBuff.buffer,
		// buf:           make([]byte, 0),
		isreader: 1,
		curIdx:   0,
	}
	pGingBuff.Lock()
	*(pGingBuff.geters) = append(*(pGingBuff.geters), *grer)
	pGingBuff.Unlock()
	return grer
}

func (wr *gwrer) Put(data []byte) (uint32, error) {
	var status int8 = 0
	wr.fathergingbuf.put(data, &status)
	if status != 0 {
		return 0, fmt.Errorf("%d", status)
	}
	wr.curIdx += uint32(len(data))
	return uint32(len(data)), nil
}

func (wr *gwrer) Get() ([]byte, error) {
	bufptr := wr.fathergingbuf.GetBufPtr()
	if bufptr == nil {
		return nil, fmt.Errorf("get buf ptr fail")
	}
	wr.curIdx += wr.fathergingbuf.wx - wr.fathergingbuf.rx
	return (*bufptr)[:wr.fathergingbuf.wx], nil
}

func (wr *gwrer) GetId() uint32 {
	return wr.id
}

func (g *gingbuff) put(data []byte, status *int8) {
	wx, rx := g.wx, g.rx
	if wx+uint32(len(data)) > BUFSIZE { //atomic
		*status = 0x1f
		return
	}

	if (rx != 0) && (wx+uint32(len(data)) >= rx) { // atomic
		*status = 0x1e
		return
	}

	gbufs := (*[]byte)(g.buf)

	*gbufs = append(*gbufs, data...)
	atomic.AddUint32(&g.wx, uint32(len(data))) // atomic
}

// func (g *gingbuff) Read()
