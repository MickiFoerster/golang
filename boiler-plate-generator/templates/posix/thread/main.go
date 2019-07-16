package main

import (
	"html/template"
	"log"
	"os"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("*.gotemplate"))
}

func main() {
	type arg struct {
		Type string
		Name string
	}
	type Thread struct {
		WorkerFunctionReturnType string
		WorkerFunctionDefintions string
		WorkerFunctionName       string
		WorkerFunctionArgs       []arg
		WorkerFunctionReturnStmt string
	}

	thread := Thread{
		WorkerFunctionReturnType: "void*",
		WorkerFunctionName:       "task",
		WorkerFunctionArgs:       []arg{{"char*", "buffer"}, {"int", "n"}},
		WorkerFunctionDefintions: "",
		WorkerFunctionReturnStmt: "return NULL;",
	}

	err := tpl.ExecuteTemplate(os.Stdout, "posix-thread.gotemplate", thread)
	if err != nil {
		log.Fatalln(err)
	}
}
