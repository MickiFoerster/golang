package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

var verbose = flag.Bool("v", false, "show verbose progress messages")
var concurrency = make(chan struct{}, 32)

func main() {
	// Determine the initial directories from command line
	flag.Parse()
	roots := flag.Args()
	if len(roots) == 0 {
		roots = []string{"."}
	}

	// Traverse each root of the file tree in parallel
	for _, root := range roots {
		walkDir(root)
	}
}

func walkDir(dir string) {
    concurrency <- struct{}{}
	for _, entry := range dirents(dir) {
		if entry.IsDir() {
			if entry.Name() == ".git" {
				cmd := exec.Command("git", "pull")
				cmd.Dir = dir
				stdout, err := cmd.StdoutPipe()
				if err != nil {
					log.Println("Could not get stdout from command:", err)
					return
				}
				go io.Copy(os.Stdout, stdout)

				stderr, err := cmd.StderrPipe()
				if err != nil {
					log.Println("Could not get stderr from command:", err)
					return
				}
				go io.Copy(os.Stderr, stderr)

				err = cmd.Start()
				if err != nil {
					log.Println("Could not start command:", err)
					return
				}

				if err := cmd.Wait(); err != nil {
					fmt.Printf("'git pull' in %s failed: %s\n", dir, err)
					continue
				}
				fmt.Printf("executed 'git pull' in %s successfully\n", dir)

				// Test if submodules exist, if so then update submodules
				gitmodules := filepath.Join(dir, ".gitmodules")
				if _, err := os.Stat(gitmodules); !os.IsNotExist(err) {
					cmd = exec.Command("git", "submodule", "update", "--init", "--recursive")
					cmd.Dir = dir
					stdout, err := cmd.StdoutPipe()
					if err != nil {
						log.Println("Could not get stdout from command:", err)
						return
					}
					go io.Copy(os.Stdout, stdout)

					stderr, err := cmd.StderrPipe()
					if err != nil {
						log.Println("Could not get stderr from command:", err)
						return
					}
					go io.Copy(os.Stderr, stderr)

					err = cmd.Start()
					if err != nil {
						log.Println("Could not start command:", err)
						return
					}

					if err := cmd.Wait(); err != nil {
						log.Printf("'git submodule update' in %s failed:", err)
						continue
					}
					fmt.Printf("executed 'git submodule update' in %s successfully\n", dir)
				}
			} else {
				subdir := filepath.Join(dir, entry.Name())
				walkDir(subdir)
			}
		}
	}
    <-concurrency
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
