package main

import (
	"io"
	"net/http"
	"time"

	"github.com/MickiFoerster/GoExamples/context/log"
)

func main() {
	http.HandleFunc("/", log.Decorate(handler))
	panic(http.ListenAndServe("127.0.0.1:8080", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log.Println(ctx, "handler started\n")
	defer log.Println(ctx, "handler ended\n")

	select {
	case <-time.After(5 * time.Second):
		io.WriteString(w, "hello\n")
	case <-ctx.Done():
		err := ctx.Err()
		log.Println(ctx, err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
