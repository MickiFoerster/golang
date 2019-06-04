package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sync"
)

const numOfConcurrentWorker = 8

var (
	filelist = make(chan string)
	done     = make(chan struct{})
	wg       sync.WaitGroup
)

func main() {
	go func() {
		concurrentSemaphore := make(chan struct{}, 32)
		for fn := range filelist {
			concurrentSemaphore <- struct{}{} // down
			wg.Add(1)
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
	wg.Wait()
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
	pattern := fmt.Sprintf(`.*%s.*`, os.Args[1])
	patternlen := len(os.Args[1])
	regExpr := regexp.MustCompile(pattern)

	f, err := os.Open(fn)
	if err != nil {
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanBytes)

	linecounter := 1
	counter := 0
	buf := make([]byte, 128)

	for scanner.Scan() {
		bytes := scanner.Bytes()
		for _, b := range bytes {
			if b == '\n' {
				if matched := regExpr.Match(buf[:counter]); matched {
					fmt.Printf("%s:%d: %s\n", fn, linecounter, string(buf))
				}
				linecounter++
				counter = 0
			}
			if counter+1 == len(buf) {
				if matched := regExpr.Match(buf); matched {
					fmt.Printf("%s:%d: %s\n", fn, linecounter, string(buf))
				}

				for i := 0; i < patternlen; i++ {
					buf[i] = buf[len(buf)-patternlen+i]
				}
				counter = patternlen
			}
			buf[counter] = b
			counter++
		}
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error: read operation failed: %s", err)
	}

	if counter > 0 {
		if matched := regExpr.Match(buf[:counter]); matched {
			fmt.Printf("%s:%d: %s\n", fn, linecounter, string(buf))
		}
	}
	wg.Done()

	return nil
}
