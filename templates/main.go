package main

import (
	"log"
	"os"
	"text/template"
)

func main() {
	tpl, err := template.ParseFiles("tpl.gohtml")
	if err != nil {
		log.Fatal(err)
	}
	fd, err := os.Create("index.html")
	if err != nil {
		log.Fatal(err)
	}
	defer fd.Close()

	err = tpl.Execute(fd, nil)
	if err != nil {
		log.Fatal(err)
	}
}
