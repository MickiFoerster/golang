package main

import "fmt"

func main() {
	v := []uint8{2, 5, 1, 6, 8, 4, 4, 10, 234, 123, 88, 33, 54, 42}
	//insertion_sort(v)
	insertion_search_sort_with_binary_search(v)
	fmt.Println(v)
}

func insertion_sort(v []uint8) {
	for i := 1; i < len(v); i++ {
		t := v[i]
		j := i
		for ; j > 0 && v[j-1] > t; j-- {
			v[j] = v[j-1]
		}
		v[j] = t
	}
}

// We exploit the fact that the array left from index i is sorted
// Benchmark shows that this is about 20% faster (N=1024*256)
func insertion_search_sort_with_binary_search(v []uint8) {
	for i := 1; i < len(v); i++ {
		t := v[i]
		l := 0
		r := i
		for l < r {
			m := (l + r) / 2
			if v[m] < t { // insertion location is on right half
				l = m + 1
			} else { // otherwise continue with left half
				r = m
			}
		}
		// Here, index r points to the location where t is to be inserted

		// Move all elements up to r one position to the right
		for j := i; j > r; j-- {
			v[j] = v[j-1]
		}
		v[r] = t
	}
}
