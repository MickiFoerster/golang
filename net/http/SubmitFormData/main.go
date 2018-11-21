package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

type myArbitraryType float64

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseFiles("form.gohtml"))
}

func (t myArbitraryType) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Client connected")
	//fmt.Fprintln(w, "Welcome to this great web app")

	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}

	tpl.ExecuteTemplate(w, "form.gohtml", r.Form)
}

func main() {
	var t myArbitraryType
	http.ListenAndServe(":3000", t)
}
