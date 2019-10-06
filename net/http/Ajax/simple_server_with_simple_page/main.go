package main

import (
	"fmt"
	"html/template"
	"net/http"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("templates/*"))
}

func main() {
	GetCPUUsage()
	http.HandleFunc("/", index)
	http.HandleFunc("/getCPUusage", getCPUusage)
	http.ListenAndServe(":8080", nil)
}

func index(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "index.gohtml", nil)
}

var cpuusage float64

func getCPUusage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, fmt.Sprintf("%0.2f%%", cpuusage))
}
