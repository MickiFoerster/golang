package main

import (
	"log"
	"net/http"
	"path"
)

func main() {
	http.HandleFunc("/", kenny)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func kenny(w http.ResponseWriter, req *http.Request) {
	p := path.Join("/tmp", "a.mp3")

	/* First method
	f, err := os.Open(p)
	if err != nil {
		http.Error(w, "file not found", 404)
		return
	}
	defer f.Close()
	io.Copy(w, f)
	*/

	// Second Method
	http.ServeFile(w, req, p)
}
