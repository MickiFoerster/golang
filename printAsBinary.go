package main

import ( 
	"fmt"
	"os"
	"strconv"
)

func main() {
	for _, arg := range os.Args[1:] {
		number, err := strconv.Atoi(arg)
		if err != nil {
			fmt.Printf("Cannot convert '%s' into a number.\n", arg)
			continue
		}
		fmt.Printf("%08v = %032b\n", arg, number)
	}
}
