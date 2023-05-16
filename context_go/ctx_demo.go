package context_go

import (
	"context"
	"fmt"
	"time"
)

func Context_TestInit() {
	ctx, close := context.WithTimeout(context.Background(), 1*time.Second)
	defer close()

	for {
		select {
		case <-ctx.Done():
			fmt.Printf("the ctx done\n")
			goto brk
		default:
			fmt.Printf("default case\n")
		}
	}

brk:
	fmt.Printf("brk out;\n")
}
