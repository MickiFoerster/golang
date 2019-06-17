package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup

	broadcast_ch := make(chan struct{})
	broadcast_ch2 := make(chan struct{})

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			select {
			case <-broadcast_ch:
				fmt.Println("broadcast received by routine", i)
			case <-broadcast_ch2:
				fmt.Println("broadcast2 received by routine", i)
			}
		}(i)
	}
	fmt.Println("sleep")
	time.Sleep(2 * time.Second)
	fmt.Println("broadcast by closing the channel")
	close(broadcast_ch)
	close(broadcast_ch2)
	wg.Wait()
}
