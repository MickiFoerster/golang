package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
)

var verbose = flag.Bool("v", false, "show verbose progress messages")

func main() {
	// Determine the initial directories from command line
	flag.Parse()
	roots := flag.Args()
	if len(roots) == 0 {
		roots = []string{"."}
	}

	// Traverse each root of the file tree in parallel
	var wg sync.WaitGroup
	for _, root := range roots {
		wg.Add(1)
		go walkDir(root, &wg)
	}
	wg.Wait()
}

func walkDir(dir string, wg *sync.WaitGroup) {
	defer wg.Done()
	for _, entry := range dirents(dir) {
		if entry.IsDir() {
			if entry.Name() == ".git" {
				var stdout, stderr bytes.Buffer
				cmd := exec.Command("git", "pull")
				cmd.Dir = dir
				cmd.Stdout = &stdout
				cmd.Stderr = &stderr
				err := cmd.Run()
				if err != nil {
					fmt.Printf("'git pull' in %s failed:\nstdout: %s\nstderr: %s\n", dir, stdout, stderr)
					continue
				}
				fmt.Printf("executed 'git pull' in %s successfully\n", dir)

				// Test if submodules exist, if so then update submodules
				gitmodules := filepath.Join(dir, ".gitmodules")
				if _, err := os.Stat(gitmodules); !os.IsNotExist(err) {
					cmd = exec.Command("git", "submodule", "update", "--init", "--recursive")
					cmd.Dir = dir
					cmd.Stdout = &stdout
					cmd.Stderr = &stderr
					err = cmd.Run()
					if err != nil {
						fmt.Printf("'git pull' in %s failed:\nstdout: %s\nstderr: %s\n", dir, stdout, stderr)
						continue
					}
					fmt.Printf("executed 'git submodule update' in %s successfully\n", dir)
				}
			} else {
				subdir := filepath.Join(dir, entry.Name())
				wg.Add(1)
				go walkDir(subdir, wg)
			}
		}
	}
}

// sema is a counting semaphore for limiting concurrency
var sema = make(chan struct{}, 20)

func dirents(dir string) []os.FileInfo {
	sema <- struct{}{}        // acquire token
	defer func() { <-sema }() // release token
	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: error: %v\n", os.Args[0], err)
		return nil
	}
	return entries
}
