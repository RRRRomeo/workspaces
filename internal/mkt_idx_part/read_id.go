/*
 * @Author: hfouyang hfouyang@quant360.com
 * @Date: 2023-03-20 15:33:13
 * @LastEditors: hfouyang hfouyang@quant360.com
 * @LastEditTime: 2023-03-29 10:57:40
 * @FilePath: /map_chan/internal/mkt_idx_part/read_id.go
 */
package mkt_idx_part

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"log"
	"os"
	"unsafe"
)

const (
	TICKLEN uint32 = uint32(unsafe.Sizeof(MdsL2Trade{})) + 8
	ORDLEN  uint32 = uint32(unsafe.Sizeof(MdsL2Order{})) + 8
	SNAPLEN uint32 = uint32(unsafe.Sizeof(MdsMktSZL2Snapshot{})) + 8
)

type idx_reader interface {
	ReadIdx(*os.File)
}

type Reader struct {
	r_id uint32
	Off  uint32
	f    *os.File
}

func NewReader(id uint32, ffp string) *Reader {
	out, e := os.Open(ffp)
	if e != nil {
		return nil
	}

	return &Reader{
		r_id: id,
		Off:  0,
		f:    out,
	}
}

func (r *Reader) GetId() uint32 {
	return r.r_id
}

func (r *Reader) ReleaseReader() {
	r.f.Close()
	r.f = nil
}

func (r *Reader) ReadTo(tlv *Tlv[any]) error {
	var err error
	var typ uint32
	var siz uint32

	if tlv == nil {
		return errors.New("inner uploader is nil")
	}

	if r.f == nil {
		log.Fatalf("inner reader fd nil\n")
		return err
	}

	typBuf := make([]byte, 4)
	sizBuf := make([]byte, 4)

	n, e := r.f.Read(typBuf)
	if e != nil {
		// ... if EOF close file
		if e == io.EOF {
			r.f.Close()
		}
		// log.Fatalf("read typ fail:%s\n", e)
		return e
	}
	e = binary.Read(bytes.NewReader(typBuf), binary.LittleEndian, &typ)
	if e != nil {
		log.Fatalf("tran typ fail:%s\n", e)
		return e
	}
	r.Off += uint32(n)

	n, e = r.f.Read(sizBuf)
	if e != nil {
		log.Fatalf("read size fail:%s\n", e)
		return e
	}
	e = binary.Read(bytes.NewReader(sizBuf), binary.LittleEndian, &siz)
	if e != nil {
		log.Fatalf("read size fail:%s\n", e)
		return e
	}
	r.Off += uint32(n)

	dataBuf := make([]byte, siz)
	n, e = r.f.Read(dataBuf)
	if e != nil {
		log.Fatalf("read size fail:%s\n", e)
		return e
	}
	if uint32(n) != siz {

		log.Fatalf("read data err:buf len:%d != siz:%d\n", n, siz)
		return e
	}
	r.Off += uint32(n)

	tlv.Header.DataTyp = typ
	tlv.Header.BufLen = siz

	switch typ {
	case TICK_TYPE:
		l2Trade := MdsL2Trade{}
		e := binary.Read(bytes.NewReader(dataBuf), binary.LittleEndian, &l2Trade)
		if e != nil {
			log.Fatalf("read l2Trade fail\n")
			return e
		}
		tlv.Data = l2Trade

	case ORDER_TYPE:
		l2Order := MdsL2Order{}
		e := binary.Read(bytes.NewReader(dataBuf), binary.LittleEndian, &l2Order)
		if e != nil {
			log.Fatalf("read l2Trade fail\n")
			return e
		}
		tlv.Data = l2Order

	case SNAPSHOT_TYPE:
		l2Snapshot := MdsMktSZL2Snapshot{}
		e := binary.Read(bytes.NewReader(dataBuf), binary.LittleEndian, &l2Snapshot)
		if e != nil {
			log.Fatalf("read l2Snapshot fail\n")
			return e
		}
		tlv.Data = l2Snapshot
	}
	return nil
}

func (r *Reader) ReadToChan(ch chan Tlv[any]) error {
	var err error
	tlv := Tlv[any]{}

	if r.f == nil {
		log.Fatalf("inner reader fd nil\n")
		return err
	}
	e := r.ReadTo(&tlv)
	if e != nil {
		return e
	}

	ch <- tlv
	return nil
}

// func getTypLen(typ uint16) uint32 {
// 	var l uint32 = 0
// 	switch typ {
// 	case btree_idx_demo.ETICK:
// 		l = TICKLEN
// 	case btree_idx_demo.EORD:
// 		l = ORDLEN
// 	case btree_idx_demo.ESNAP:
// 		l = SNAPLEN
// 	}

// 	return l
// }

func (r *Reader) ReadFrom(off uint32, tlv *Tlv[any]) error {
	// l := getTypLen(typ)
	if off != 0 {
		r.f.Seek(int64(off), 0)
	}
	e := r.ReadTo(tlv)
	if e != nil {
		return e
	}

	// ... offset reset
	return nil

}
