/*
 * @Author: hfouyang hfouyang@quant360.com
 * @Date: 2023-03-29 09:43:05
 * @LastEditors: hfouyang hfouyang@quant360.com
 * @LastEditTime: 2023-03-29 13:05:57
 * @FilePath: /githuab.com/RRRRomeo/workspaces/internal/mkt_idx_core/idx_reader.go
 */
package mktidxcore

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"

	"githuab.com/RRRRomeo/workspaces/qsorter"
)

const (
	TIMEMIN uint32 = 19000101
	TIMEMAX uint32 = 20991231
)

type Reader interface {
	Read()
}

type IdxReader struct {
	date uint32
	fd   *os.File
	mp   sync.Map // ... fds
}

func NewIdxReader(date uint32) (*IdxReader, error) {
	if date <= TIMEMIN || date >= TIMEMAX {
		return nil, errors.New("the date invalid")
	}

	ffp := fmt.Sprintf("/home/ty/data/221_data/csv_to_bin/idx/%d.idx", date)
	tmp, e := os.Open(ffp)
	if e != nil {
		return nil, e
	}

	return &IdxReader{
		date: date,
		fd:   tmp,
		mp:   sync.Map{},
	}, nil
}

func makeBufAndRead(f *os.File, buflen uint32, dataloader any) error {
	if buflen <= 0 {
		return errors.New("buflen is valiad")
	}

	buf := make([]byte, buflen)
	n, e := f.Read(buf)

	if n != int(buflen) {
		return errors.New("read len isnot equal to buf len")
	}

	if e != nil {
		// ... if EOF close file
		if e == io.EOF {
			f.Close()
		}
		// log.Fatalf("read typ fail:%s\n", e)
		return e
	}

	e = binary.Read(bytes.NewReader(buf), binary.LittleEndian, dataloader)
	if e != nil {
		log.Fatalf("tran typ fail:%s\n", e)
		return e
	}

	return nil
}

func (r *IdxReader) Read() (*qsorter.Qsorter_write_node, error) {
	// ... return the typ off and instrId
	node := &qsorter.Qsorter_write_node{}
	var err error

	if r.fd == nil {
		log.Fatalf("inner reader fd nil\n")
		return nil, err
	}

	e := makeBufAndRead(r.fd, 4, &node.SZInStrId)
	if e != nil {
		return nil, e
	}

	e = makeBufAndRead(r.fd, 4, &node.Off)
	if e != nil {
		return nil, e
	}

	// log.Printf("node:%v\n", node)
	return node, nil
}

func ChangeStrId(id int32) string {
	ids := strconv.Itoa(int(id))
	l := len(ids)
	subs := strings.Repeat("0", 6-l)
	return fmt.Sprintf("%s%d", subs, id)
}
