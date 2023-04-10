package main

import (
	"fmt"
	"map_chan/log/api"
	"sync"
)

func main() {
	wg := sync.WaitGroup{}
	// loggger := log.NewLoggerWithOutter(log.LOG_LEVEL_DBG, 1, "./test.log")
	loggger := api.LoggerInit(0, 1, "./test.log")
	loggger.Infof("test logger using: %d\n", 1)
	loggger.Warnf("test logger using: %d\n", 1)
	loggger.Errf("test logger using: %d\n", 1)
	// loggger.Fatalf("test logger using: %d\n", 1)
	// log.DelLogger(loggger)
	loggger.Debugf("test logger using: %d\n", 1)

	fmt.Printf("=============================================\n")
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			api.Debugf("test logger using debug\n")
			wg.Done()
		}()

		wg.Add(1)
		go func() {
			api.Infof("test logger using info\n")
			wg.Done()
		}()

		wg.Add(1)
		go func() {
			api.Warnf("test logger using Warn\n")
			wg.Done()
		}()

		wg.Add(1)
		go func() {
			api.Errf("test logger using Err\n")
			wg.Done()
		}()
	}
	wg.Wait()
	api.Infof("test logger using info: end\n")
}
