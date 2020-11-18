package main

import (
	"fmt"
	"log"

	"gopkg.in/yaml.v2"
)

type StructB struct {
	B       string  `yaml:"b"`
	D       float64 `yaml:"d"`
	StructA struct {
		A string `yaml:"a"`
		C int64  `yaml:"c"`
	} `yaml:"subobject"`
}

var data = `
b: a string from struct B
d: 42.3
subobject:
    a: a string from struct A
    c: 23
`

func main() {
	var b StructB

	err := yaml.Unmarshal([]byte(data), &b)
	if err != nil {
		log.Fatalf("cannot unmarshal data: %v", err)
	}
	fmt.Printf("%v\n", b.StructA)
	fmt.Printf("%v\n", b)
}
