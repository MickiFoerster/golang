package main

import (
	"fmt"
	"log"
	"os"
	"text/template"
)

var tpl *template.Template

type awkCall struct {
	SectionName string
	ValueName   string
	NewValue    string
}

func init() {
	tpl = template.Must(template.ParseFiles("awk.gotemplate"))
}

func main() {
	if len(os.Args) < 4 {
		log.Fatalf("syntax error: Please give section name, value to change, and new value as parameter\n")
	}

	awk := awkCall{
		SectionName: os.Args[1],
		ValueName:   os.Args[2],
		NewValue:    os.Args[3],
	}

	const fn = "awk-ChangeIniFile.sh"
	f, err := os.Create(fn)
	if err != nil {
		log.Fatal(err)
	}
	_ = f.Chmod(0700)

	err = tpl.Execute(f, awk)
	if err != nil {
		log.Fatal(err)
	}

	f.Close()
	fmt.Printf("AWK script has been written to %q.\n", fn)
}
