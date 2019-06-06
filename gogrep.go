package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
)

const numOfConcurrentWorker = 8

var (
	filelist = make(chan string)
	done     = make(chan struct{})
	wg       sync.WaitGroup
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("syntax error: Give pattern to search for as first argument")
	}

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
	defer wg.Done()

	pattern := os.Args[1]
	patternlen := len(os.Args[1])
	//regExpr := regexp.MustCompile(pattern)

	f, err := os.Open(fn)
	if err != nil {
		return err
	}
	defer f.Close()

	linecounter := 1
	linestart := 0

	buf := make([]byte, 32)
	n := 0
	patternidx := 0
	for {
		var err error
		linematches := false
		linestart = 0
		n, err = f.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		for i := 0; i < n; i++ {
			b := buf[i]
			if b == pattern[patternidx] {
				if patternidx+1 == patternlen {
					linematches = true
					patternidx = 0
				} else {
					patternidx++
				}
			} else {
				patternidx = 0
			}

			if b == '\n' {
				if linematches {
					fmt.Printf("%s:%d: %s\n", fn, linecounter, string(buf[linestart:i+1]))
					linematches = false
				}

				linecounter++
				if i+1 < n {
					linestart = i + 1
				}
			}
		}
		if linematches {
			fmt.Printf("%s:%d: %s\n", fn, linecounter, string(buf[linestart:n]))
		}
	}

	return nil
}
