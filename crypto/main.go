package main

import (
	"fmt"
	"os"
	"crypto/sha256"
	"crypto/sha512"
)

func main() {
	for _, arg := range os.Args[1:] {
		f, err := os.Open(arg)
		if err != nil {
			fmt.Printf("error: Could not open file %s: %s\n", arg, err)
			continue
		}
		defer f.Close()
		buf := make([]byte, 4096)
		buffer := make([]byte, 10000)
		// Need to increase buffer if necessary:
		// See https://blog.golang.org/go-slices-usage-and-internals
		h256 := sha256.New()
		h512 := sha512.New()
		index := 0
		for {
			count, err := f.Read(buf)
			if err != nil {
				fmt.Printf("error: Could not read from file %s: %s\n", arg, err)
				continue
			}
			copied := copy(buffer[index:], buf[:count])
		  fmt.Printf("copied %d bytes\n", copied)
			index += copied
			if count < 4096 {
				break
			}
		}
		//fmt.Printf("%x\n", buffer)
		fmt.Printf("%x  %s (sha256)\n", h256.Sum(buffer), arg)
		fmt.Printf("%x  %s (sha512)\n", h512.Sum(buffer), arg)
	}
}
