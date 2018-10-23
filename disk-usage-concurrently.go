package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

var verbose = flag.Bool("v", false, "show verbose progress messages")

func main() {
	// Determine the initial directories from command line
	flag.Parse()
	roots := flag.Args()
	if len(roots) == 0 {
		roots = []string{"."}
	}

	// Traverse the file tree in background task
	fileSizes := make(chan int64)
	go func() {
		for _, root := range roots {
			walkDir(root, fileSizes)
		}
		close(fileSizes)
	}()

	// Print results periodically
	var nfiles, nbytes int64
	for size := range fileSizes {
		nfiles++
		nbytes += size
	}
	printDiskUsage(nfiles, nbytes)

}

func printDiskUsage(nfiles, nbytes int64) {
	if nbytes < 10000 {
		fmt.Printf("%d files %.1f Byte\n", nfiles, float64(nbytes))
		return
	}
	nbytes >>= 10 
	if nbytes < 10000 {
		fmt.Printf("%d files %.1f KB\n", nfiles, float64(nbytes))
		return
	}
	nbytes >>= 10 
	if nbytes < 10000 {
		fmt.Printf("%d files %.1f MB\n", nfiles, float64(nbytes))
		return
	}
	nbytes >>= 10 
	fmt.Printf("%d files %.1f GB\n", nfiles, float64(nbytes)/1e9)
}

func walkDir(dir string, fileSizes chan<- int64) {
	for _, entry := range dirents(dir) {
		if entry.IsDir() {
			subdir := filepath.Join(dir, entry.Name())
			walkDir(subdir, fileSizes)
		} else {
			fileSizes <- entry.Size()
		}
	}
}

func dirents(dir string) []os.FileInfo {
	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: error: %v\n", os.Args[0], err)
		return nil
	}
	return entries
}


