package main

import (
	"fmt"
	"regexp"
)

var httpSchemeRE = regexp.MustCompile(`^https?:`) // "http:" or "https:"

func main() {
	s := "http://www.example.com"
	test(s)
	s = "htts://www.example.com"
	test(s)
	s = "https://www.example.com"
	test(s)
}

func test(s string) {
	if matched := httpSchemeRE.Match([]byte(s)); matched {
		fmt.Println(s, " matches")
	} else {
		fmt.Println(s, " does not match")
	}
}
