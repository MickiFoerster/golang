package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

type Options []struct {
	Option struct {
		Name         string `yaml:"name"`
		Abbreviation string `yaml:"abbreviation"`
		HasArg       struct {
			Type             string `yaml:"type"`
			OptionalArgument struct {
				Type   string   `yaml:"type"`
				Values []string `yaml:"values"`
			} `yaml:"optional_argument"`
		} `yaml:"has_arg"`
	} `yaml:"option"`
}

func main() {
	const fn = "getoptions.yaml"
	data, err := ioutil.ReadFile(fn)
	if err != nil {
		log.Fatalf("error: could not read file: %v", err)
	}

	var opts Options
	err = yaml.Unmarshal(data, &opts)
	if err != nil {
		log.Fatalf("cannot unmarshal data: %v", err)
	}
	for _, opt := range opts {
		fmt.Printf("%v\n", opt)
	}
}
