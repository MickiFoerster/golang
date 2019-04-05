// printLinks prints all links found in a website given by URL on command line
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"sync"

	"golang.org/x/net/html"
)

type WDRKiRaKa struct {
	Context           string `json:"@context"`
	Type              string `json:"@type"`
	Name              string `json:"name"`
	URL               string `json:"url"`
	Image             string `json:"image"`
	ProductionCompany struct {
		Type string `json:"@type"`
		Name string `json:"name"`
	} `json:"productionCompany"`
	IsFamilyFriendly bool   `json:"isFamilyFriendly"`
	AccessMode       string `json:"accessMode"`
	ThumbnailURL     string `json:"thumbnailUrl"`
	Episodes         []struct {
		Type       string `json:"@type"`
		Name       string `json:"name"`
		URL        string `json:"url"`
		AccessMode string `json:"accessMode"`
		Author     struct {
			Type string `json:"@type"`
			Name string `json:"name"`
		} `json:"author"`
		CopyrightHolder struct {
			Type string `json:"@type"`
			Name string `json:"name"`
		} `json:"copyrightHolder"`
		Creator struct {
			Type string `json:"@type"`
			Name string `json:"name"`
		} `json:"creator"`
		IsFamilyFriendly bool   `json:"isFamilyFriendly"`
		ThumbnailURL     string `json:"thumbnailUrl"`
	} `json:"episodes"`
}

var (
	jsonData WDRKiRaKa
	wg       sync.WaitGroup
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
		if strings.HasSuffix(link, ".mp3") {
			wg.Add(1)
			downloadMP3(link)
		}
	}
}

// visit appends to links each link found in n and returns the result
func visit(links []string, n *html.Node, resp *http.Response) []string {
	if n.Type == html.ElementNode && n.Data == "script" {
		for _, attr := range n.Attr {
			if attr.Key == "type" && attr.Val == "application/ld+json" {
				for c := n.FirstChild; c != nil; c = c.NextSibling {
					if c.Type == html.TextNode {
						//fmt.Printf("Found: %v\n", c.Data)
						err := json.Unmarshal([]byte(c.Data), &jsonData)
						if err != nil {
							log.Fatal("error json unmarshalling:", err)
						}
						//fmt.Println(jsonData)
						for _, episode := range jsonData.Episodes {
							links = append(links, episode.URL)
						}
					}
				}
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		links = visit(links, c, resp)
	}
	return links
}

func downloadMP3(mp3URL string) error {
	response, err := http.Get(mp3URL)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	filename := path.Base(response.Request.URL.Path)
	outputFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	_, err = io.Copy(outputFile, response.Body)
	if err != nil {
		return err
	}

	defer func() {
		fmt.Printf("File %q successfully downloaded", filename)
	}()

	return nil
}
