// Package ring implements operations on circular lists.
package main

import (
	"container/ring"
	"fmt"
)

func main() {
	// Create a new ring of size 5
	r := ring.New(7)

	// Get the length of the ring
	n := r.Len()

	// Initialize the ring with some integer values
	for i := 0; i < n; i++ {
		r.Value = i
		r = r.Next()
	}

	// Iterate through the ring and print its contents
	r.Do(func(p interface{}) {
		fmt.Println(p.(int))
	})

	fmt.Println("Extend ring to see what happen with contained elements")
	for i := 0; i < n/2; i++ {
		r.Value = 10 + i
		r = r.Next()
	}

	// Iterate through the ring and print its contents
	r.Do(func(p interface{}) {
		fmt.Println(p.(int))
	})
}
