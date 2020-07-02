package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

var wg sync.WaitGroup

func executor(ctx context.Context) <-chan struct{} {
	overall := make(chan struct{})
	wg.Add(2)
	task(ctx, "A", 100*time.Millisecond)
	task(ctx, "B", 500*time.Millisecond)
	go func() {
		wg.Wait()
		log.Println("all tasks are done")
		overall <- struct{}{}
	}()
	return overall
}

func task(ctx context.Context, name string, waittime time.Duration) {
	go func() {
		defer wg.Done()
		counter := 0
		for {
			select {
			case <-ctx.Done():
				fmt.Println("DONE")
				return
			default:
				fmt.Println("task ", name, ": ", counter)
				time.Sleep(waittime)
			}
			counter += 1
		}
	}()
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	done := executor(ctx)

	go func() {
		reader := bufio.NewReader(os.Stdin)
		fmt.Println("Press <RETURN> to end")
		reader.ReadString('\n')
		log.Println("main cancels sub tasks")
		cancel()
	}()

	<-done
	err := ctx.Err()
	if err != nil {
		fmt.Printf("cancelation reason: %v\n", err)
	}
}
