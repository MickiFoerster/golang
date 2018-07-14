package main

import (
	"fmt"
	"os"
	"bufio"
)

func main() {
	files := os.Args[1:]
	if len(files) == 0 {
		fmt.Printf("syntax error: %s <file1> [<file2>, ...]\n", os.Args[0])
	} else {
		for _, arg := range files {
			f, err := os.Open(arg)
			if err != nil {
				fmt.Printf("error: Cannot open file '%s'.\n", arg)
				continue
			}
			fmt.Printf("Read from file '%s'\n", arg)
			input := bufio.NewScanner(f)
			for input.Scan() {
				fmt.Println(input.Text())
			}
			f.Close()
		}
	}
}
