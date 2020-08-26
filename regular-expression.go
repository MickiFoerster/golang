package main

import (
	"fmt"
	"regexp"
)

var (
	httpSchemeRE = regexp.MustCompile(`^https?:`) // "http:" or "https:"
	re           = regexp.MustCompile(`([a-z]+)://([^:]+):?([0-9]*)`)
)

func main() {
	s := "http://www.example.com"
	test(s)
	s = "htts://www.example.com"
	test(s)
	s = "https://www.example.com"
	test(s)
	s = "ssh://www.example.com"
	test(s)
	s = "ssh://www.example.com:22"
	test(s)
}

func test(s string) {
	if matched := httpSchemeRE.Match([]byte(s)); matched {
		fmt.Println(s, " matches")
	} else {
		fmt.Println(s, " does not match")
	}

	fmt.Printf("%q\n", re.FindAllStringSubmatch(s, -1))
}
