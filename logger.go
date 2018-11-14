package main

import (
	"fmt"
	"log"
	"os"
)

func init() {
	fd, err := os.Create("log.txt")
	if err != nil {
		fmt.Println(err)
	}
	log.SetOutput(fd)
}

func main() {
	_, err := os.Open("not-existing-file.txt")
	if err != nil {
		log.Fatalln("error happened", err)
	}
}
