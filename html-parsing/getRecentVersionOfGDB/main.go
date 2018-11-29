package main

import (
	"fmt"
	"log"
	"net/http"

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

func traverseDOM(node *html.Node) {
	fmt.Println(node.Data)
	if node.Data == "The most recent release" {
		fmt.Println(node.Type)
	}
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		traverseDOM(c)
	}
}
