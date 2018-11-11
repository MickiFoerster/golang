package main

import (
	"fmt"
	"net/http"
)

type mytype int

func (m mytype) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "This is the handler for the response")
}

func main() {
	var asdf mytype
	http.ListenAndServe(":8080", asdf)
}

