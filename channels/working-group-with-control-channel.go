package main

import "fmt" 

func main() {
	c := make(chan int)
	control_channel := make(chan bool)

	go func() {
		for i := 0; i< 10; i++ {
			c <- i
		}
		control_channel <- true
	}()

	go func() {
		for i := 10; i< 20; i++ {
			c <- i
		}
		control_channel <- true
	}()

	go func() {
		<-control_channel 
		<-control_channel 
		close(c)
		close(control_channel)
	}()

	for n := range c {
		fmt.Println(n)
	}
}
