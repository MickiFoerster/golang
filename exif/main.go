package main

import (
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("syntax error: %s <path to JPG file>\n", os.Args[0])
	}

	f, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal("error: Opening file", os.Args[1], "failed")
	}
	f.Close()
}
