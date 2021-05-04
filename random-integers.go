package main

import (
	"fmt"
	"math/rand"
	"time"
)

var histogram = map[int]int{}

const numSamples = 4000

func main() {
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < numSamples; i++ {
		histogram[rand.Intn(32)]++
	}

	fmt.Println("Number | Occurrences")
	for i := 0; i < numSamples; i++ {
		if v, ok := histogram[i]; ok {
			fmt.Printf("   %03d | ", i)
			for j := 0; j < v; j++ {
				fmt.Print("*")
			}
			fmt.Println()
		}
	}
}
