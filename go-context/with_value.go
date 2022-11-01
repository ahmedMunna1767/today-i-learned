package main

import (
	"context"
	"fmt"
)

type contextKey string

func WithValue(ctx context.Context) {
	var authToken contextKey = "auth_token"
	innerCtx := context.WithValue(ctx, authToken, "XYZ_123")
	fmt.Println(innerCtx.Value(authToken))
}
