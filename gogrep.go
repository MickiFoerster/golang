package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
)

const numOfConcurrentWorker = 8

var (
	filelist = make(chan string)
	done     = make(chan struct{})
)

func main() {
	go func() {
		concurrentSemaphore := make(chan struct{}, 32)
		for fn := range filelist {
			concurrentSemaphore <- struct{}{} // down
			go func(path string) {
				err := grep(path)
				if err != nil {
					log.Fatal(err)
				}
				<-concurrentSemaphore // up
			}(fn)
		}
		done <- struct{}{}
	}()
	go func() {
		err := filepath.Walk(".", walker)
		if err != nil {
			log.Fatal("error: Cannot read list of files:", err)
		}
		close(filelist)
	}()
	<-done
}

func walker(path string, info os.FileInfo, err error) error {
	if err != nil {
		fmt.Println("error:", err)
		return err
	}
	if !info.IsDir() {
		filelist <- path
	}
	return nil
}

func grep(fn string) error {
	if len(os.Args) != 2 {
		return fmt.Errorf("syntax error: Give pattern to search for as first argument")
	}

	f, err := os.Open(fn)
	if err != nil {
		return err
	}
	pattern := fmt.Sprintf(`.*%s.*`, os.Args[1])
	regExpr := regexp.MustCompile(pattern)
	scanner := bufio.NewScanner(f)
	counter := 1
	for scanner.Scan() {
		if matched := regExpr.Match([]byte(scanner.Text())); matched {
			fmt.Printf("%s:%d: %s\n", fn, counter, scanner.Text())
		}
		counter++
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error: read operation failed:", err)
	}

	return nil
}
