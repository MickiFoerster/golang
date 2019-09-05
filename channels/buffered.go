package main

import (
	"log"
	"time"
)

func test(sz int) {
	ch := make(chan int, sz)
	go func() {
		ch <- 42
		log.Println("42 has been sent")
	}()
	go func() {
		ch <- 23
		log.Println("23 has been sent")
	}()
	log.Println("Received: ", <-ch)
	//time.Sleep(5 * time.Second)
	log.Println("Received: ", <-ch)
	time.Sleep(1 * time.Second)
}

func main() {
	log.Println("start of test with buffer size 0")
	test(0)
	log.Println("start of test with buffer size 1")
	test(1)
}
