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
	http.HandleFunc("/", index)
	http.ListenAndServe(":8080", nil)
}

func index(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		host := req.FormValue("host")
		port := req.FormValue("port")
		fmt.Println(w, `<html>`)
		fmt.Println(w, `<body>`)
		fmt.Fprintf(w, `<pre>%s</pre>`, host)
		fmt.Fprintf(w, `<pre>%s</pre>`, port)
		fmt.Println(w, `</body>`)
		fmt.Println(w, `</html>`)
		return
	}
	tpl.ExecuteTemplate(w, "login.gohtml", nil)
}
