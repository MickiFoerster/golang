package main

import (
	"fmt"
	"sync"
)

func main() {
	num_of_workers := 1000

	errCh := make(chan error, num_of_workers)
	var wg sync.WaitGroup
	wg.Add(num_of_workers)

	for i := 0; i < num_of_workers; i++ {
		go func(i int) {
			defer wg.Done()
			if i == 42 {
				errCh <- fmt.Errorf("Thread %v is returning an error", i)
				return
			}
		}(i)
	}
	wg.Wait()

	select {
	case err := <-errCh:
		fmt.Printf("error: %v\n", err)
	default:
		fmt.Println("all successful")
	}
}
