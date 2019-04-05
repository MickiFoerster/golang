// printLinks prints all links found in a website given by URL on command line
package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"
	"sync"

	"golang.org/x/net/html"
)

var (
	wg sync.WaitGroup
)

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
}

// htmlParse get a URL as input and downloads the corresponding document and
// starts the search for links afterwards.
func htmlParse(url string) {
	defer wg.Done()
	resp, err := http.Get(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: Could not download %q\n", url)
		return
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error: Cannot read HTML body:", err)
		return
	}
	resp.Body.Close()

	htmlrootnode, err := html.Parse(bytes.NewReader(data))
	if err != nil {
		fmt.Fprintln(os.Stderr, "error: HTML parser failed:", err)
		return
	}
	for _, link := range visit(nil, htmlrootnode, resp) {
		fmt.Printf("Found link %q\n", link)
	}
}

// visit appends to links each link found in n and returns the result
func visit(links []string, n *html.Node, resp *http.Response) []string {
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, a := range n.Attr {
			if a.Key == "href" {
				link, err := resp.Request.URL.Parse(a.Val)
				if err != nil {
					continue
				}
				links = append(links, link.String())
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		links = visit(links, c, resp)
	}
	return links
}
