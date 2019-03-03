package main

import (
	"flag"
	"fmt"
)

var ub = flag.Int64("u", 10000000, "Upper bound until prime numbers are computed")

func main() {
	flag.Parse()
	for i := int64(2); i < *ub; i++ {
		foundPrime := true
		if i%2 == 0 {
			foundPrime = false
		} else {
			for j := int64(3); j*j <= i; j += 2 {
				if i%j == 0 {
					foundPrime = false
					break
				}
			}
		}
		if foundPrime {
			fmt.Println(i)
		}
	}
}
