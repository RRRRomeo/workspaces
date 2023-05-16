package getter

import (
	"fmt"
	"sync"
	"testing"

	"gitlab.quant360.com/hfouyang/qlog/log"
	"gitlab.quant360.com/hsw/sz_match/types"
)

var wg sync.WaitGroup

func TestDataLoader_Get(t *testing.T) {
	wg = sync.WaitGroup{}
	type fields struct {
		fp   string
		fn   string
		date string
	}
	type args struct {
		// ch chan types.MdData
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
		{"001", fields{"/home/ty/data/221_data/csv_to_bin/", "000001", "20220701"}, args{}},
		{"002", fields{"/home/ty/data/221_data/csv_to_bin/", "300081", "20220701"}, args{}},
		{"003", fields{"/home/ty/data/221_data/csv_to_bin/", "300774", "20220701"}, args{}},
		// {"004", fields{"/home/ty/data/221_data/cdv_to_bin/", "300081", "20220701"}, args{}},
		// {"005", fields{"/home/ty/data/221_data/cdv_to_bin/", "300772", "20220701"}, args{}},
		// {"006", fields{"/home/ty/data/221_data/cdv_to_bin/", "000005", "20220701"}, args{}},
	}
	for i := 0; i < len(tests)-2; i++ {
		wg.Add(3)
		go t.Run(tests[i].name, func(t *testing.T) {
			// log.Debugf("test start\n")
			d := NewDataLoader(tests[i].fields.fp, tests[i].fields.fn, tests[i].fields.date)
			fp := fmt.Sprintf("./%s.log", tests[i].fields.fn)
			logger := log.LoggerInit(0, 1, fp)
			defer log.LoggerDeinit(logger)

			ch := make(chan types.MdData, 1000)
			go d.Get(ch)

			// time.Sleep(200 * time.Millisecond)
			for {
				// log.Debugf("<-ch getter...\n")
				data, ok := <-ch
				if !ok {
					// eontinue
					break
				}

				logger.Infof("%v\n", data)
			}
			wg.Done()

		})
		go t.Run(tests[i].name, func(t *testing.T) {
			// log.Debugf("test start\n")
			d := NewDataLoader(tests[i+1].fields.fp, tests[i+1].fields.fn, tests[i+1].fields.date)
			fp := fmt.Sprintf("./%s.log", tests[i+1].fields.fn)
			logger := log.LoggerInit(0, 1, fp)
			defer log.LoggerDeinit(logger)

			ch := make(chan types.MdData, 1000)
			go d.Get(ch)

			// time.Sleep(200 * time.Millisecond)
			for {
				// log.Debugf("<-ch getter...\n")
				data, ok := <-ch
				if !ok {
					// eontinue
					break
				}

				logger.Infof("%v\n", data)
			}
			wg.Done()

		})
		go t.Run(tests[i].name, func(t *testing.T) {
			// log.Debugf("test start\n")
			d := NewDataLoader(tests[i+2].fields.fp, tests[i+2].fields.fn, tests[i+2].fields.date)
			fp := fmt.Sprintf("./%s.log", tests[i+2].fields.fn)
			logger := log.LoggerInit(0, 1, fp)
			defer log.LoggerDeinit(logger)

			ch := make(chan types.MdData, 1000)
			go d.Get(ch)

			// time.Sleep(200 * time.Millisecond)
			for {
				// log.Debugf("<-ch getter...\n")
				data, ok := <-ch
				if !ok {
					// eontinue
					break
				}

				logger.Infof("%v\n", data)
			}
			wg.Done()

		})
		wg.Wait()
	}
}
