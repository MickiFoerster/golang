package main

import (
	"log"
	"time"
)

func main() {
	ch := make(chan bool)
	for i := 0; i < 3; i++ {
		go func(i int) {
			<-ch
			log.Println("After recv", i)
		}(i)
	}
	log.Println("Sleeping ...")
	time.Sleep(5 * time.Second)
	log.Println("Continue ...")
	close(ch)
	time.Sleep(3 * time.Second)
	log.Println("Done")
}
