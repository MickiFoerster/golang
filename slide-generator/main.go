package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseFiles("tpl.gohtml"))
}

func main() {
	http.HandleFunc("/", serveMainRoute)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func serveMainRoute(w http.ResponseWriter, req *http.Request) {
	data := struct {
		Title    string
		Headline string
		Items    []string
		Author   string
		Date     string
	}{
		Title:    "Free Lunch is Over",
		Headline: "Overview",
		Items: []string{
			"Introduction",
			"Approach",
			"Results",
			"Summary",
		},
		Author: "John Doo",
		Date:   fmt.Sprint(time.Now().Format(time.RFC850)),
	}
	err := tpl.ExecuteTemplate(w, "tpl.gohtml", data)
	if err != nil {
		log.Fatal(err)
	}
}
