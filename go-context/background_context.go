package main

import (
	"context"
	"fmt"
)

// background or todo context doesn't contain deadline / done channel
func BackgroundContext(ctx context.Context) {
	_, ok := ctx.Deadline()
	if !ok {
		fmt.Println("no deadline is set")
	}
	done := ctx.Done()
	if done == nil {
		fmt.Println("channel is nil")
	}
}
