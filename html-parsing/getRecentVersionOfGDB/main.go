package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

func main() {
	resp, err := http.Get("https://www.gnu.org/software/gdb/download/")
	handleErr(err)

	node, err := html.Parse(resp.Body)
	resp.Body.Close()
	handleErr(err)
	traverseDOM(node)
}

func handleErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

var recentReleaseLocationFound = false
var recentReleaseLinkFound = false

func traverseDOM(node *html.Node) {
	if recentReleaseLinkFound && recentReleaseLocationFound {
		fmt.Println("The most recent GDB version is:", node.Data)
		recentReleaseLinkFound = false
		recentReleaseLocationFound = false
	}
	if node.Type == html.ElementNode && node.Data == "a" && recentReleaseLocationFound {
		recentReleaseLinkFound = true
	}

	if !recentReleaseLocationFound && strings.Index(node.Data, "The most recent release") != -1 {
		recentReleaseLocationFound = true
	}

	for c := node.FirstChild; c != nil; c = c.NextSibling {
		traverseDOM(c)
	}
}
