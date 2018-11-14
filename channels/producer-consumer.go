package main

import (
	"fmt"
)

func main() {
	c := make(chan int)
	done := make(chan struct{})

	go producer(c)
	go consumer(c, done)
	go consumer(c, done)
	go consumer(c, done)
	go consumer(c, done)

	<-done
	<-done
	<-done
	<-done
}

func producer(c chan<- int) {
	for i := 0; i < 100; i++ {
		c <- i
	}
	close(c)
}

func consumer(c <-chan int, done chan<- struct{}) {
	for i := range c {
		fmt.Print("    ", i)
	}
	done <- struct{}{}
}
