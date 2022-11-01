package main

import (
	"context"
	"fmt"
)

func main() {
	WithTimeout(context.Background())
	fmt.Println("exiting")
}
