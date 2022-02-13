package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func main() {
	buf := readBuf()
	if err := parse(buf); err != nil {
		fmt.Println("Parsing failed:", err)
		os.Exit(1)
	}
	fmt.Println("Parsing successfully")
}

func readBuf() []byte {
	reader := bufio.NewReader(os.Stdin)

	buf := make([]byte, 4096)
	n, err := reader.Read(buf)
	if err != nil {
		if err == io.EOF {
			os.Exit(0)
		}
		fmt.Println(err)
		os.Exit(1)
	}
	buf = buf[:n]

	return buf
}
