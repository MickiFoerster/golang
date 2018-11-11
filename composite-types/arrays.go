package main

import "fmt"

func main() {
	var a [3]int = [3]int{1,2,3}
	b := [...]int{1,2,3}
	fmt.Println(a[0])
	fmt.Println(a[len(a)-1])
	fmt.Println(b[len(b)-1])
}
