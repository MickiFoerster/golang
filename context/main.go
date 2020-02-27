package main

import (
	"context"
	"fmt"
	"log"
	"time"
)

func main() {
	//ctx, cancel := context.WithCancel(context.Background())
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// time.AfterFunc(time.Second, cancel)

	go func() {
		time.Sleep(100 * time.Millisecond)
		cancel()
	}()

	sleepAndTalk(ctx, 5*time.Second, "hello")
}

func sleepAndTalk(ctx context.Context, d time.Duration, msg string) {
	select {
	case <-time.After(d):
		fmt.Println(msg)
	case <-ctx.Done():
		log.Println(ctx.Err())
	}
}
