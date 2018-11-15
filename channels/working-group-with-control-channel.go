package main

import "fmt"

func main() {
	c := make(chan int)
	controlChannel := make(chan bool)

	go func() {
		for i := 0; i < 10; i++ {
			c <- i
		}
		controlChannel <- true
	}()

	go func() {
		for i := 10; i < 20; i++ {
			c <- i
		}
		controlChannel <- true
	}()

	go func() {
		<-controlChannel
		<-controlChannel
		close(c)
		close(controlChannel)
	}()

	for n := range c {
		fmt.Println(n)
	}
}
