package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"
)

type Effector func(context.Context) (string, error)

func Retry(effector Effector, retries int, delay time.Duration) Effector {
	return func(ctx context.Context) (string, error) {
		for r := 0; r < retries; r++ {
			count := ctx.Value("try").(int)
			ctx := context.WithValue(ctx, "try", count+1)
			res, err := effector(ctx)
			if err == nil {
				return res, err
			}
			log.Printf("Function call failed, retrying in %v", delay)

			select {
			case <-time.After(delay):
			case <-ctx.Done():
				return "", err
			}
		}
		return "", errors.New("ahmed")
	}
}

func GetPdfUrl(ctx context.Context) (string, error) {
	count := ctx.Value("try").(int)
	if count <= 6 {
		return "", errors.New("boom")
	} else {
		return "https://linktopdf.com", nil
	}
}

func main() {
	r := Retry(GetPdfUrl, 5, 2*time.Second)
	ctx := context.WithValue(context.Background(), "try", 0)
	res, err := r(ctx)
	fmt.Println(res, err)
}
