// htmlElementCounter creates a histogram of HTML elements of a document
// provided as URL on command line.
package main

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"strings"

	"golang.org/x/net/html"
)

var histogram = make(map[string]uint32)

func main() {
	if len(os.Args) == 1 {
		fmt.Fprintf(os.Stderr,
			"syntax error: %s <URL>\n", path.Base(os.Args[0]))
		os.Exit(1)
	}
	url := os.Args[1]
	if !strings.HasPrefix(url, "http") && !strings.HasPrefix(url, "https") {
		fmt.Fprintf(os.Stderr, "Argument %q seems not to be a valid URL.", url)
		os.Exit(1)
	}
	htmlParse(url)

	for k, v := range histogram {
		fmt.Printf("%-15s: %8d\n", k, v)
	}
}

// htmlParse get a URL as input and downloads the corresponding document and
// starts the search for links afterwards.
func htmlParse(url string) {
	htmldoc, err := http.Get(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: Could not download %q\n", url)
		return
	}

	htmlRootNode, err := html.Parse(htmldoc.Body)
	htmldoc.Body.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: HTML parser failed: %v\n", err)
		return
	}
	countElements(htmlRootNode)
}

func countElements(node *html.Node) {
	//fmt.Println("Encountered HTML node type", node.Type, "Data:", node.Data)
	if node.Type == html.ElementNode {
		histogram[node.Data]++
	}
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		countElements(c)
	}
}
