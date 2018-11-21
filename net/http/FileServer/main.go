package main

import (
	"log"
	"net/http"
)

func main() {
	http.Handle("/", http.FileServer(http.Dir("/home")))
	http.Handle("/tempfiles", http.StripPrefix("/tempfiles", http.FileServer(http.Dir("/tmp"))))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
