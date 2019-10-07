package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"
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
		http.Redirect(w, req, "/shell", http.StatusSeeOther)
		return
	}
	fmt.Println(`method is GET`)
	tpl.ExecuteTemplate(w, "index.gohtml", nil)
}

func checkIfCommandExist(w http.ResponseWriter, req *http.Request) {
	cmd, err := ioutil.ReadAll(req.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not read body: %v\n", err)
	}

	which, err := exec.Command("which", string(cmd)).Output()
	if err != nil {
		fmt.Fprintf(w, "")
	}

	path := strings.TrimRight(string(which), " \r\n")
	fmt.Fprintf(w, path)
}

func shell(w http.ResponseWriter, req *http.Request) {
	tpl.ExecuteTemplate(w, "shell.gohtml", nil)
}
