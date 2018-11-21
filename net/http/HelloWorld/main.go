package main

import (
	"fmt"
	"net/http"
)

type myArbitraryType float64

func (t myArbitraryType) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Client connected")
	fmt.Fprintln(w, "Welcome to this great web app")
}

func main() {
	var t myArbitraryType
	http.ListenAndServe(":3000", t)
}
