package main

import (
	"fmt"
	"runtime"
)

func main() {
	printMemStats()
	runtime.GC()
	printMemStats()
}

func printMemStats() {
	var memstats runtime.MemStats
	runtime.ReadMemStats(&memstats)
	fmt.Println("mem.Alloc:", memstats.Alloc)
	fmt.Println("mem.TotalAlloc:", memstats.TotalAlloc)
	fmt.Println("mem.HeapAlloc:", memstats.HeapAlloc)
	fmt.Println("mem.NumGC:", memstats.NumGC)
}
