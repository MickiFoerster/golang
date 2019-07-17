package main

import (
	"html/template"
	"log"
	"os"
)

var funcMap = template.FuncMap{
	"plus1": func(x int) int { return x + 1 },
}

var tpl *template.Template

func init() {
	tpl = template.Must(template.New("c++thread.gotemplate").Funcs(funcMap).ParseFiles("c++thread.gotemplate"))
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
		WorkerFunctionReturnType: "void",
		WorkerFunctionName:       "task",
		WorkerFunctionArgs:       []arg{{"char*", "buffer"}, {"int", "n"}},
		WorkerFunctionDefintions: "",
		WorkerFunctionReturnStmt: "",
	}

	f, err := os.Create("main.cpp")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	err = tpl.Execute(f, thread)
	if err != nil {
		log.Fatalln(err)
	}
}
