package main

import (
	"fmt"
	"regexp"
)

func main() {
	patterns := []string{
		"a",
		"b",
		"test",
		"(ab)*X[cd]+",
	}

	re := "^("
	for i := 0; i < len(patterns); i++ {
		re += patterns[i]
		if i+1 < len(patterns) {
			re += "|"
		}
	}
	regexpString := re + ")$"
	fmt.Println("Regular expression is ", regexpString)
	regexpr := regexp.MustCompile(regexpString)

	for _, s := range []string{"asdf", "a", "b", "atest", "Xcd", "ababXcd", "abXcdcdcdcd", "abX"} {
		test(regexpr, s)
	}
}

func test(r *regexp.Regexp, s string) {
	if matched := r.MatchString(s); matched {
		fmt.Println(s, " matches")
	} else {
		fmt.Println(s, " is no match")
	}
}
