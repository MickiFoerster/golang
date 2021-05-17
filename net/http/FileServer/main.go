package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	if len(os.Args) != 3 {
		log.Fatalf("syntax error: %s <local folder to serve via http> <port>", os.Args[0])
	}

	dir := os.Args[1]
	port := fmt.Sprintf(":%s", os.Args[2])
	if stat, err := os.Stat(dir); os.IsNotExist(err) {
		log.Fatalf("directory %s does not exist", dir)
	} else if !stat.IsDir() {
		log.Fatalf("%s is not a directory", dir)
	}
	http.Handle("/", http.FileServer(http.Dir(dir)))
	log.Fatal(http.ListenAndServe(port, nil))
}
