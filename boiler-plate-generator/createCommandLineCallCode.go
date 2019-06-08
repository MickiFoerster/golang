package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"text/template"

	"github.com/fatih/color"
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

	const fn = "main.c"
	f, err := os.Create(fn)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	err = tpl.Execute(f, execv)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("C code has been written to %q.\n", fn)

	fmt.Println("Consider to test this code by using:")
	fmt.Printf("gcc -std=c11 -Wall -Werror %s\n", fn)
	fmt.Printf("clang -std=c17 -Wall -Werror %s\n", fn)

	fmt.Println("Will that do for you:")

	args := []string{"-std=c11", "-Wall", "-Werror", fn}
	gcc := exec.Command("gcc", args...)
	err = gcc.Run()
	showResult(err, "gcc")

	clang := exec.Command("clang", args...)
	err = clang.Run()
	showResult(err, "clang")
}

func showResult(err error, msg string) {
	fmt.Print(msg)
	if err != nil {
		color.Red(" [failed] ")
		fmt.Println(err)
	} else {
		color.Green(" [OK] ")
	}
}
