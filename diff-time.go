package main

import (
	"fmt"
	"log"
	"time"
)

func main() {
	t1, err := time.Parse(time.RFC3339, "2023-08-17T20:10:38.510979Z")
	if err != nil {
		log.Fatal(err)
	}
	t2, err := time.Parse(time.RFC3339, "2023-08-19T10:31:24.690176Z")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(t2.Sub(t1))
}
