package main

import (
	"fmt"
	"io"
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
	defer wg.Done()

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

	linecounter := 1
	linestart := 0

	// Replace this by Bayer Moore Search
	buf := make([]byte, 32)
	n := 0
	for {
		var err error
		linestart = 0
		if n == 0 {
			n, err = f.Read(buf)
		} else {
			for i := 0; i < patternlen; i++ {
				buf[i] = buf[n-patternlen+i]
			}
			n, err = f.Read(buf[patternlen:])
			n += patternlen
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		for i := 0; i < n; i++ {
			b := buf[i]
			if b == '\n' {
				if patternlen <= i-linestart+1 {
					log.Printf("look into line %v '%s'\n", linecounter, string(buf[linestart:i]))
					if matched := regExpr.Match(buf[linestart:i]); matched {
						fmt.Printf("%s:%d: %s\n", fn, linecounter, string(buf[linestart:i]))
					}
				}
				linecounter++
				linestart = i + 1
			}
		}
	}

	return nil
}
