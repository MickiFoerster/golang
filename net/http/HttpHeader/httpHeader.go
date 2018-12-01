package main

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"strings"
	"sync"
)

var wg sync.WaitGroup

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

// htmlParse asks the header information from the given URL
func htmlParse(url string) {
	defer wg.Done()
	resp, err := http.Head(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: Could not download HEAD from %q\n", url)
		return
	}

	for k, v := range resp.Header {
		fmt.Printf("%s : %s\n", k, v)
	}
}
