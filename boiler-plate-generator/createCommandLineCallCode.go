package main

import (
	"fmt"
	"log"
	"os"
	"text/template"
)

var tpl *template.Template

type execvCall struct {
	Path string
	Args []string
}

func init() {
	tpl = template.Must(template.ParseFiles("createCommandLineCallCodeC.gotemplate"))
}

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("syntax error: At least one parameter must be provided. For example:\n%s ls -l -t", os.Args[0])
	}

	execv := execvCall{os.Args[1], os.Args[1:]}

	if _, err := os.Stat(execv.Path); os.IsNotExist(err) {
		fmt.Fprintln(os.Stderr, "warning: First argument has to be full path to existent file.")
	}

	err := tpl.ExecuteTemplate(os.Stdout, "createCommandLineCallCodeC.gotemplate", execv)
	if err != nil {
		log.Fatal(err)
	}
}
