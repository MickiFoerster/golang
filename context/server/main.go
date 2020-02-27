package main

import (
	"io"
	"log"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe("127.0.0.1:8080", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	log.Printf("handler started\n")
	defer log.Printf("handler ended\n")

	time.Sleep(5 * time.Second)
	io.WriteString(w, "hello\n")
}
