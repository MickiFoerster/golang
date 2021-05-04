package main

import (
	"fmt"
	"time"
)

func main() {
	//go counter()
	go ticker()
	time.Sleep(15 * time.Second)
}

func ticker() {
	d := time.Duration(250) * time.Millisecond
	ticker := time.NewTicker(d)
	i := 0
	for t := range ticker.C {
		i++
		fmt.Println("Count", i, "at", t)
	}
}

func counter() {
	i := 0
	d := time.Duration(250) * time.Millisecond
	for {
		t := time.NewTimer(d)
		<-t.C
		i += 250
		fmt.Printf("waited %v, overall %v\n", d, time.Duration(i)*time.Millisecond)
	}
}
