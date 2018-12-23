package main

import (
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"time"
)

var tpl *template.Template
var pseudoRandomGenSeed = rand.NewSource(time.Now().UnixNano())
var pseudoRandomGen = rand.New(pseudoRandomGenSeed)
var taskCounter = 1

func init() {
	tpl = template.Must(template.ParseFiles("tpl.gohtml"))
}

func main() {
	http.HandleFunc("/", serveMainRoute)
	http.Handle("/favicon.ico", http.NotFoundHandler())
	serverAddress := ":1234"
	log.Println("Server is running on ", serverAddress)
	log.Fatal(http.ListenAndServe(serverAddress, nil))
}

func serveMainRoute(w http.ResponseWriter, req *http.Request) {
	log.Printf("Serving URL %q", req.URL)
	err := req.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	if len(req.PostForm) > 0 {
		for _, values := range req.PostForm {
			for _, val := range values {
				log.Println(val)
			}
		}
	}

	r1 := pseudoRandomGen.Intn(10) + 1
	r1s := fmt.Sprint(r1)
	r2 := pseudoRandomGen.Intn(10) + 1
	r2s := fmt.Sprint(r2)
	data := struct {
		Challenge   string
		Answerlabel string
		Counter     int
	}{
		Challenge:   "Was ist " + r1s + " + " + r2s + "?",
		Answerlabel: "Antwort",
		Counter:     taskCounter,
	}
	err = tpl.ExecuteTemplate(w, "tpl.gohtml", data)
	if err != nil {
		log.Fatal(err)
	}
	taskCounter++
}
