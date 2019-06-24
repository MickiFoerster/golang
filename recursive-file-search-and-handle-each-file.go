package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

const numOfConcurrentWorker = 8

var (
	filelist    = make(chan string)
	input_task1 = make(chan string)
)

func main() {
	dir := checkArgs()
	done := fileSearch(dir)
	<-done
}

func checkArgs() string {
	if len(os.Args) < 2 {
		log.Fatalf("syntax error: %s <path where recursive search should start>\n", os.Args[0])
	}
	dir := os.Args[1]

	if info, err := os.Stat(dir); os.IsNotExist(err) || !info.IsDir() {
		log.Fatal("error: `", dir, "` is not a valid path to a directory.")
	}
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		log.Fatal("error: Path `", dir, "` does not exist:", err)
	}

	d, err := filepath.EvalSymlinks(dir)
	if err != nil {
		return dir
	}
	return d
}

func consumeFilepaths(done chan struct{}) {
	go task1(done)
	for fn := range filelist {
		input_task1 <- fn
	}
	close(input_task1)
}

func fileSearch(dir string) chan struct{} {
	c := make(chan struct{})

	go func() {
		go consumeFilepaths(c)
		err := filepath.Walk(dir, walker)
		if err != nil {
			log.Fatal("error: Cannot read list of files:", err)
		}
		close(filelist)
	}()

	return c
}

func task1(done chan struct{}) {
	semaphore := make(chan struct{}, numOfConcurrentWorker)
	for fn := range input_task1 {
		semaphore <- struct{}{}
		go func() {
			fmt.Println(fn, "starts ... ")
			time.Sleep(time.Second)
			fmt.Println(fn, "ends")
			<-semaphore
		}()
	}
	done <- struct{}{}
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
