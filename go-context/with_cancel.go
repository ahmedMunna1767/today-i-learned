package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

// here we are getting data from randomCharGenerator
// as soon as we get "o", we are breaking the loop
// as there is nothing else to do, main function exits
// calling cancel() function
// as soon as cancel is called , Done channel is closed
// exiting  goroutine in randomCharGenerator
// this way goroutine is exited properly without leaving it unhandled
func WithCancel(ctx context.Context) {
	rand.Seed(time.Now().UnixNano())

	ctx, cancel := context.WithCancel(ctx)
	defer cancel() // cancel when generator is closed and program exits

	for char := range randomCharGenerator(ctx) {
		generatedChar := string(byte(char))
		fmt.Printf("%v\n", generatedChar)

		if generatedChar == "o" {
			break
		}
	}
}

// this function starts a goroutine that creates random characters
// this is a Generator pattern
func randomCharGenerator(ctx context.Context) <-chan int {
	char := make(chan int)

	seedChar := int('a')

	go func() {
		for {
			select {
			case <-ctx.Done():
				// ** this will not print as main function will exit immediately
				fmt.Printf("Yay! we found: %v", seedChar)
				return // returning not to leak the goroutine
			case char <- seedChar:
				seedChar = 'a' + rand.Intn(26)
			}
		}
	}()

	return char
}
