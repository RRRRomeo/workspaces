package main

import (
	"log"

	golang "gitlab.quant360.com/algo_strategy_group/algo_trading/api/golang"
)

func main() {
	errAlgoInit := golang.InitApi("./pkg/algo_trading/conf", "./pkg/algo_trading/data", "./pkg/algo_trading/logs", 20220701) // 算法初始化
	if errAlgoInit != nil {
		log.Printf("algo init fail:%s\n", errAlgoInit)
		return
		// panic(errAlgoInit)
	}
	golang.CloseApi()

}
