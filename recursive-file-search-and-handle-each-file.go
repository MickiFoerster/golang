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
	done        = make(chan struct{})
)

func main() {
	dir := checkArgs()
	go task1()
	go consumeFilepaths()
	go fileSearch(dir)
	<-done
}

func consumeFilepaths() {
	for fn := range filelist {
		input_task1 <- fn
	}
	close(input_task1)
}

func fileSearch(dir string) {
	err := filepath.Walk(dir, walker)
	if err != nil {
		log.Fatal("error: Cannot read list of files:", err)
	}
	close(filelist)
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

func task1() {
	task1WorkerDone := make(chan struct{}, numOfConcurrentWorker)
	for fn := range input_task1 {
		task1WorkerDone <- struct{}{}
		go func() {
			fmt.Println(fn, "starts ... ")
			time.Sleep(time.Second)
			fmt.Println(fn, "ends")
			<-task1WorkerDone
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
