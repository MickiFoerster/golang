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
	http.HandleFunc("/shell", shell)
	http.HandleFunc("/checkIfCommandExist", checkIfCommandExist)
	http.Handle("/favicon.ico", http.NotFoundHandler())

	http.ListenAndServe(":8080", nil)
}

func index(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		/*
			fmt.Println(`method is POST`)
			host := req.FormValue("host")
			port := req.FormValue("port")
			fmt.Fprintln(w, `<html>`)
			fmt.Fprintln(w, `<body>`)
			fmt.Fprintf(w, `<pre>%s</pre>`, host)
			fmt.Fprintf(w, `<pre>%s</pre>`, port)
			fmt.Fprintln(w, `</body>`)
			fmt.Fprintln(w, `</html>`)
		*/
		http.Redirect(w, req, "/shell", http.StatusSeeOther)
		return
	}
	fmt.Println(`method is GET`)
	tpl.ExecuteTemplate(w, "index.gohtml", nil)
}

func checkIfCommandExist(w http.ResponseWriter, req *http.Request) {

}

func shell(w http.ResponseWriter, req *http.Request) {
	tpl.ExecuteTemplate(w, "shell.gohtml", nil)
}
