// htmlOutline prints the outline of a HTML document behind the URL given
// as command line option.
package main

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"strings"

	"golang.org/x/net/html"
)

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
	outline(nil, htmlRootNode)
}

func outline(stack []string, node *html.Node) {
	//fmt.Println("Encountered HTML node type", node.Type, "Data:", node.Data)
	if node.Type == html.ElementNode {
		stack = append(stack, node.Data)
		fmt.Println(stack)
	}
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		outline(stack, c)
	}
}
