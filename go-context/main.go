package main

import (
	"context"
	"fmt"
)

func main() {
	BackgroundContext(context.Background())
	fmt.Println("exiting")
}
