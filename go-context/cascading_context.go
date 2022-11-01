package main

import (
	"context"
	"fmt"
	"time"
)

func CascadingContext(ctx context.Context) {
	c := make(chan string)
	go func() {
		time.Sleep(30 * time.Second)
		c <- "one"
	}()

	ctx1 := context.Context(ctx)
	ctx2, cancel2 := context.WithTimeout(ctx1, 20*time.Second)
	ctx3, cancel3 := context.WithTimeout(ctx2, 10*time.Second) // derives from ctx2
	ctx4, cancel4 := context.WithTimeout(ctx3, 30*time.Second) // derives from ctx3
	ctx5, cancel5 := context.WithTimeout(ctx4, 5*time.Second)  // derives from ctx4

	defer cancel2()
	defer cancel3()
	defer cancel4()
	defer cancel5()

	select {
	case <-ctx2.Done():
		fmt.Println("ctx2 closed! error: ", ctx2.Err())
	case <-ctx3.Done():
		fmt.Println("ctx3 closed! error: ", ctx3.Err())
	case <-ctx4.Done():
		fmt.Println("ctx4 closed! error: ", ctx4.Err())
	case <-ctx5.Done():
		fmt.Println("ctx5 closed! error: ", ctx5.Err())
	case msg := <-c:
		fmt.Println("received", msg)
	}
}
