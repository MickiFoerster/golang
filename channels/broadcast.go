package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup

	broadcast_ch := make(chan int)
	broadcast_ch2 := make(chan int)

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			select {
			case b := <-broadcast_ch:
				fmt.Printf("go routine %d: broadcast received by routine: %v\n", i, b)
			case b := <-broadcast_ch2:
				fmt.Printf("go routine %d: broadcast2 received by routine: %v\n", i, b)
			}
		}(i)
	}
	fmt.Println("sleep")
	time.Sleep(2 * time.Second)
	fmt.Println("broadcast by closing the channel")
	close(broadcast_ch)
	time.Sleep(1 * time.Second)
	fmt.Println("broadcast again on channel 2 by closing the channel")
	close(broadcast_ch2)
	wg.Wait()
}
