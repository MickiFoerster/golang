package main

import "fmt"

// units for kilobyte, megabyte, etc.
const (
	_  = iota             // 0
	KB = 1 << (iota * 10) // 1 << (1 * 10)
	MB = 1 << (iota * 10) // 1 << (1 * 10)
	GB = 1 << (iota * 10) // 1 << (1 * 10)
	TB = 1 << (iota * 10) // 1 << (1 * 10)
)

func main() {
	fmt.Printf("%b\t%d\n", KB, KB)
	fmt.Printf("%b\t%d\n", MB, MB)
	fmt.Printf("%b\t%d\n", GB, GB)
	fmt.Printf("%b\t%d\n", TB, TB)
}
