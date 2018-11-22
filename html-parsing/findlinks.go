// findlinks gets a list of URLs as input and prints all links it finds inside
// the HTML documents to the standard output.
package main

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"strings"
	"sync"

	"golang.org/x/net/html"
)

var wg sync.WaitGroup
var links = make(map[string][]string, len(os.Args))

func main() {
	if len(os.Args) == 1 {
		fmt.Fprintf(os.Stderr,
			"syntax error: %s <URL1> [<URL2>] [...]\n", path.Base(os.Args[0]))
		os.Exit(1)
	}
	for _, url := range os.Args[1:] {
		if !strings.HasPrefix(url, "http") && !strings.HasPrefix(url, "https") {
			fmt.Fprintf(os.Stderr, "Argument %q seems not to be a valid URL.", url)
			continue
		}
		wg.Add(1)
		go htmlParse(url)
	}
	wg.Wait()

	for url, linklist := range links {
		fmt.Println(url, ":")
		for _, link := range linklist {
			fmt.Println(" ", link)
		}
	}
}

// htmlParse get a URL as input and downloads the corresponding document and
// starts the search for links afterwards.
func htmlParse(url string) {
	defer wg.Done()
	htmldoc, err := http.Get(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: Could not download %q\n", url)
		return
	}

	htmlrootnode, err := html.Parse(htmldoc.Body)
	htmldoc.Body.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: HTML parser failed: %v\n", err)
		return
	}
	links[url] = visit(nil, htmlrootnode)
}

// visit appends to links each link found in n and returns the result
func visit(links []string, n *html.Node) []string {
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, a := range n.Attr {
			if a.Key == "href" {
				links = append(links, a.Val)
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		links = visit(links, c)
	}
	return links
}
