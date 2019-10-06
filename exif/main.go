package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/mknote"
)

var filelist = make([]string, 0, 1024)

func init() {
	if len(os.Args) < 2 {
		log.Fatalf("syntax error: %s <path where files with EXIF data are located>\n", os.Args[0])
	}
	dir := os.Args[1]

	if info, err := os.Stat(dir); os.IsNotExist(err) || !info.IsDir() {
		log.Fatal("error: `", dir, "` is not a valid path to a directory.")
	}
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		log.Fatal("error: Path `", dir, "` does not exist:", err)
	}

	err := filepath.Walk(dir, walker)
	if err != nil {
		log.Fatal("error: Cannot read list of files:", err)
	}
}

func main() {
	for _, fn := range filelist {
		handleFile(fn)
	}
}

func handleFile(fn string) {
	f, err := os.Open(fn)
	if err != nil {
		log.Fatal("error: Opening file", fn, "failed:", err)
	}
	defer f.Close()

	exif.RegisterParsers(mknote.All...)

	x, err := exif.Decode(f)
	if err != nil {
		log.Fatalf("Could not decode %q: %s", fn, err)
	}
	tm, _ := x.DateTime()
	fmt.Println("Taken: ", tm)

	lat, long, err := x.LatLong()
	if err == nil {
		fmt.Println("lat, long: ", lat, ", ", long)
	}
}

func walker(path string, info os.FileInfo, err error) error {
	if err != nil {
		fmt.Println("error:", err)
		return err
	}
	if !info.IsDir() {
		filelist = append(filelist, path)
	}
	return nil
}
