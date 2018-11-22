package main

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/fsnotify/fsnotify"
)

func main() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "syntax error: %q <Go file to execute>\n", path.Base(os.Args[0]))
		os.Exit(1)
	}

	GoFile := os.Args[1]

	done := make(chan struct{})
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				log.Println("event:", event)
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Println("modified file:", event.Name)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add(GoFile)
	if err != nil {
		log.Fatal(err)
	}
	<-done
}
