package main

import (
	"fmt"
	"html/template"
	"net/http"
	"time"
)

var (
	tpl      *template.Template
	cpuusage float64
)

const updateTime = 1000 * time.Millisecond

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
	tpl.ExecuteTemplate(w, "index.gohtml", updateTime.Milliseconds())
}

func getCPUusage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, fmt.Sprintf(`{"CPUUsage":"%6.2f%%"}`, cpuusage))
}
