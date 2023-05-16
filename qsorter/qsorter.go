package qsorter

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"os"
	"sync"

	"gitlab.quant360.com/hfouyang/qlog/log"
)

const MAXSIZE = 10 * 1024 * 1024

type Status func(string)

type Qsorter_write_node struct {
	SZInStrId int32
	Off       uint32
}

type Qsorter_node struct {
	Qsorter_write_node
	BidIdx uint32
	Tim    int32
	N      *Qsorter_node
}

type Qsorter_hash_node struct {
	Qsorter_write_node
	N *Qsorter_hash_node
}

type Qsorter_interface interface {
	Sort() *[]Qsorter_write_node
}

type Qsorter struct {
	mu sync.Mutex
	// be_sort_buf    [MAXSIZE]byte
	// curpos         int32
	sorted_pointer map[int32]*Qsorter_hash_node
}

func NewQSorter() *Qsorter {
	return &Qsorter{
		// curpos:         0,
		sorted_pointer: make(map[int32]*Qsorter_hash_node, 10000000),
	}
}

func (q *Qsorter) Store(node *Qsorter_node) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	tnode := &Qsorter_hash_node{
		Qsorter_write_node: Qsorter_write_node{
			SZInStrId: node.SZInStrId,
			Off:       node.Off,
		},
		N: nil,
	}

	v, ok := q.sorted_pointer[node.Tim]
	if !ok {
		q.sorted_pointer[node.Tim] = tnode
		return nil
	}

	for v != nil {
		// log.Debugf("find next\n")
		if v.N == nil {
			v.N = tnode
			break
		}
		v = v.N
	}

	return nil
}

func (q *Qsorter) Pop(fp string) {
	// for k, v := range q.sorted_pointer {
	fd, err := os.OpenFile(fp, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0755)
	if err != nil {
		return
	}
	defer fd.Close()

	for k := 91500000; ; k++ {

		if k > 151000000 {
			break
		}

		wfd := bufio.NewWriter(fd)
		q.mu.Lock()
		v, ok := q.sorted_pointer[int32(k)]
		if ok {
			// log.Debugf("k:%d\n", k)
			i := 0
			for v != nil {
				log.Infof("%d ==>v:%v\n", i, v)
				n := &Qsorter_write_node{
					SZInStrId: v.SZInStrId,
					Off:       v.Off,
				}

				err := WriteToFile(wfd, n)
				if err != nil {
					log.Errf("write into file fail:%s\n")
				}
				v = v.N
				i++
			}
		}
		delete(q.sorted_pointer, int32(k))
		// atomic.AddInt32(&q.curpos, -1)
		q.mu.Unlock()
		// log.Debugf("map del 1\n")

	}
}

func WriteToFile(wfd *bufio.Writer, node *Qsorter_write_node) error {
	buf := &bytes.Buffer{}

	err := binary.Write(buf, binary.LittleEndian, node)
	if err != nil {
		log.Errf("bin write err:%s\n", err)
		return err
	}

	_, err = wfd.Write(buf.Bytes())

	if err != nil {
		log.Errf("文件创建并写入失败！错误：%s\n", err)
		return err
	}
	wfd.Flush()

	return nil
}
