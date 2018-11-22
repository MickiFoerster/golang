package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", foo)
	http.Handle("/favicon.ico", http.NotFoundHandler())
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func foo(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	fmt.Println(r.URL.Path)
	fmt.Fprintln(w, "Example for sending Internal Server Error")
}
