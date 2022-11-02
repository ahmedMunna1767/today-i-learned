package main

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

func ContextTestServer() {
	http.ListenAndServe(":8000", http.HandlerFunc(handler))
}

func handler(w http.ResponseWriter, r *http.Request) {
	s := time.Now()
	ctx := r.Context()
	go whatever(ctx)
	select {
	case <-time.After(2 * time.Second):
		w.Write([]byte(s.String() + "<-- -->" + time.Now().String()))
	case <-ctx.Done():
		fmt.Println("request cancelled")
		return
	}
	fmt.Println("exiting")
}

func whatever(ctx context.Context) {
	<-ctx.Done()
	fmt.Println("waited 5 seconds")
}
