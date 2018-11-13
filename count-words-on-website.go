package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	website, err := http.Get("https://www.sueddeutsche.de/")
	if err != nil {
		log.Fatal(err)
	}

	/*
		body, err := ioutil.ReadAll(website.Body)
		if err != nil {
			log.Fatal(err)
		}
		defer website.Body.Close()
		bodyString := string(body)
		fmt.Println(bodyString)
	*/

	words := make(map[string]int)

	scanner := bufio.NewScanner(website.Body)
	scanner.Split(bufio.ScanWords)

	for scanner.Scan() {
		words[scanner.Text()]++
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "error while reading website body: ", err)
		os.Exit(1)
	}

	for word := range words {
		fmt.Printf("%s : %d\n", word, words[word])
	}
}
