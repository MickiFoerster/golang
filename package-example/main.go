package main

import (
	"apackage"
	"fmt"
)

func main() {
	fmt.Println("Using apackage")
	apackage.A()
	apackage.B()
	fmt.Println(apackage.MyConstant)
}
