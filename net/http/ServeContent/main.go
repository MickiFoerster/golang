package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/", kenny)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func kenny(w http.ResponseWriter, req *http.Request) {
	f, err := os.Open("Kenny.jpg")
	if err != nil {
		http.Error(w, "file not found", 404)
		return
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		http.Error(w, "file not found", 404)
		return
	}

	http.ServeContent(w, req, f.Name(), fi.ModTime(), f)
}
