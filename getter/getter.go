package getter

import (
	"gitlab.quant360.com/algo_strategy_group/mds_data_playback/src"
	"gitlab.quant360.com/hsw/sz_match/types"
)

type getter interface {
	get()
}

type DataLoader struct {
	fp   string
	fn   string
	date string
}

func NewDataLoader(fp string, fn string, date string) *DataLoader {
	return &DataLoader{
		fp:   fp,
		fn:   fn,
		date: date,
	}
}

func (d *DataLoader) get(ch chan types.MdData) {
	// log.Debugf("get...\n")
	for {
		err := src.Get3In1DataWithChan(d.fp, d.fn, d.date, ch)
		if err != nil {
			close(ch)
			// fmt.Println("Dispatch: ", err)
			break
		}
	}
}

func (d *DataLoader) Get(ch chan types.MdData) {
	// log.Debugf("Get...\n")
	if ch == nil {
		return
	}
	d.get(ch)
}
