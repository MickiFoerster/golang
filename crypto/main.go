package main

import (
	"fmt"
	"os"
	"io"
	"bytes"
	"crypto/sha256"
	"crypto/sha512"
)

func main() {
	counter := 0
	results256 := make(map[string][]byte)
	results512 := make(map[string][]byte)
	for _, arg := range os.Args[1:] {
		counter++
		f, err := os.Open(arg)
		if err != nil {
			fmt.Printf("error: Could not open file %s: %s\n", arg, err)
			continue
		}
		defer f.Close()
		buf := make([]byte, 4096)
		// Need to increase buffer if necessary:
		// See https://blog.golang.org/go-slices-usage-and-internals
		h256 := sha256.New()
		h512 := sha512.New()
		for {
			count, err := f.Read(buf)
			if err != nil {
				if err == io.EOF {
					break
				}
				fmt.Printf("error: Could not read from file %s: %s\n", arg, err)
				continue
			}
		  //fmt.Printf("%d bytes successfully read\n", count)
			bytereader := bytes.NewReader(buf[:count])
			_, err = io.Copy(h256, bytereader)
			if  err != nil {
				fmt.Printf("error: Could not copy bytes from file into sha256 hash function: %s\n", err)
				continue
			}
		  //fmt.Printf("%d bytes successfully copied to sha256 hash function\n", copied)
			bytereader.Seek(0, io.SeekStart)
			_, err = io.Copy(h512, bytereader);
			if err != nil {
				fmt.Printf("error: Could not copy bytes from file into sha512 hash function: %s\n", err)
				continue
			}
		  //fmt.Printf("%d bytes successfully copied to sha512 hash function\n", copied)
			if count < 4096 {
				break
			}
		  //fmt.Printf("Continue to read from file\n")
		}
		//fmt.Printf("%s = %x\n", arg, h256.Sum(nil))
		results256[arg] = h256.Sum(nil)
		results512[arg] = h512.Sum(nil)
	}

	if counter > 0 {
		for key, value := range results256 {
			fmt.Printf("sha256: %128x  %s\n", value, key)
			fmt.Printf("sha512: %128x  %s\n", results512[key], key)
		}
	} else {
		fmt.Printf("syntax error: Give files as argument.")
	}
}
