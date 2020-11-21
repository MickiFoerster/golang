package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

type Option struct {
	Name         string `yaml:"name"`
	Abbreviation string `yaml:"abbreviation"`
	HasArg       string `yaml:"has_arg"`
}

type Options struct {
	Options Option `yaml:"option"`
}

func main() {
	//const fn = "getoptions.yaml"
	const fn = "a.yaml"
	data, err := ioutil.ReadFile(fn)
	if err != nil {
		log.Fatalf("error: could not read file: %v", err)
	}

	var opt Options
	err = yaml.Unmarshal(data, &opt)
	if err != nil {
		log.Fatalf("cannot unmarshal data: %v", err)
	}
	fmt.Printf("%v\n", opt)
}
