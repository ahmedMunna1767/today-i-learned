package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"time"
)

func WithTimeout(ctx context.Context) {
	// context with deadline after 2 millisecond
	ctx, cancel := context.WithTimeout(ctx, 2*time.Millisecond)
	defer cancel()

	lineRead := make(chan string)
	exitChan := make(chan string)

	var fileName = "sample-file.txt"
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("failed opening file: %s", err)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	// goroutine to read file line by line and passing to channel to print
	go func() {
		for scanner.Scan() {
			lineRead <- scanner.Text()
		}

		close(lineRead)
		file.Close()
		exitChan <- "reached end of the file"
	}()

outer:
	for {
		// printing file line by line until deadline is reached
		select {
		case cause := <-exitChan:
			fmt.Println("process stopped. reason: ", cause)
			break outer
		case <-ctx.Done():
			fmt.Println("process stopped. reason: ", ctx.Err())
			break outer
		case line := <-lineRead:
			fmt.Println(line)
		}
	}
}
