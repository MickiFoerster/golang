// findlinks gets a list of URLs as input and prints all links it finds inside
// the HTML documents to the standard output.
package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"sync"

	"golang.org/x/net/html"
)

const websiteJSONFileName = "websites.json"

type website struct {
	URL      string
	BodyHash string
	Links    []string
}

var (
	wg       sync.WaitGroup
	websites = make(map[string]website, len(os.Args))
)

func init() {
	jsondata, err := ioutil.ReadFile(websiteJSONFileName)
	if err != nil {
		fmt.Fprintln(os.Stderr, "warning: Could not read from JSON file:", err)
		return
	}
	var w []website
	if err := json.Unmarshal(jsondata, &w); err != nil {
		log.Fatal(err)
	}
	for _, ws := range w {
		websites[ws.URL] = ws
	}
	fmt.Printf("Initialize websites data structure from JSON file %q successfully.\n", websiteJSONFileName)
}

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

	var w []website
	for _, ws := range websites {
		w = append(w, ws)
		/*
			fmt.Println(url, ":", ws.BodyHash)
			for _, link := range ws.Links {
				fmt.Println("  ", link)
			}
		*/
	}

	jsondata, err := json.MarshalIndent(w, "", "    ")
	if err != nil {
		log.Fatal(err)
	}

	fp, err := os.Create(websiteJSONFileName)
	if err != nil {
		log.Fatal(err)
	}
	defer fp.Close()
	if _, err := fp.Write(jsondata); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Slice of websites were written to file successfully.\n")
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

	hash := sha256.Sum256(data)
	hashstr := fmt.Sprintf("%x", hash)

	if hashstr == websites[url].BodyHash {
		fmt.Printf("Body of %q same as before, so skip parsing.\n", url)
		return
	}

	htmlrootnode, err := html.Parse(bytes.NewReader(data))
	if err != nil {
		fmt.Fprintln(os.Stderr, "error: HTML parser failed:", err)
		return
	}

	websites[url] = website{
		URL:      url,
		BodyHash: hashstr,
		Links:    visit(nil, htmlrootnode, resp),
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
