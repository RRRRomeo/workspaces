/*
 * @Author: hfouyang hfouyang@quant360.com
 * @Date: 2023-03-29 10:23:34
 * @LastEditors: hfouyang hfouyang@quant360.com
 * @LastEditTime: 2023-03-29 11:25:06
 * @FilePath: /githuab.com/RRRRomeo/workspaces/cmd/idx_read_cmd/idx_read_cmd.go
 */
package main

import (
	"fmt"
	"log"

	mktidxcore "githuab.com/RRRRomeo/workspaces/internal/mkt_idx_core"
	"githuab.com/RRRRomeo/workspaces/internal/mkt_idx_part"
)

func CallIdxRead() {
	idxreader, err := mktidxcore.NewIdxReader(20220701)
	if err != nil {
		log.Printf("new Idx reader fail:%s\n", err)
		return
	}
	for {
		node, err := idxreader.Read()
		if err != nil {
			log.Printf("reader read node fail:%s\n", err)
			break
		}

		// log.Printf("CallIdxRead ==> node:%v\n", node)
		sid := mktidxcore.ChangeStrId(node.SZInStrId)
		ffp := fmt.Sprintf("/home/ty/data/221_data/csv_to_bin/%d/%s.csv.3in1.0", 20220701, sid)

		log.Printf("ffp:%s\n", ffp)
		csvreader := mkt_idx_part.NewReader(0, ffp)
		if csvreader == nil {
			continue
		}
		tlv := &mkt_idx_part.Tlv[any]{}
		csvreader.ReadFrom(node.Off, tlv)
		csvreader.ReleaseReader()
		log.Printf("tlv:%v\n", tlv)
	}
}

func main() {
	// log.Printf("start call idxRead\n")
	CallIdxRead()
	// log.Printf("call idxRead end\n")
}
