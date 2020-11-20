package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"debug/elf"

	"github.com/u-root/u-root/pkg/ldd"
)

func main() {
	const cwd_string = "current working directory"

	startdir := flag.String("startdir", cwd_string,
		"Starting directory of recursive search.\nThe default is the current working directory.\n")
	flag.Parse()

	// Read current working directory in case no dir is given per option
	if *startdir == cwd_string {
		c, err := os.Getwd()
		if err != nil {
			log.Fatalf("error: could not read current working directory: %v", err)
		}
		*startdir = c
	}

	fi, err := os.Lstat(*startdir)
	if err != nil {
		log.Fatalf("error: cannot read %q: %v", *startdir, err)
	}

	if fm := fi.Mode(); !fm.IsDir() {
		log.Fatalf("error: %q is not an directory.\n", *startdir)
	}

	ch := getElfFiles(*startdir)
	for fn := range ch {
		lst, err := ldd.List([]string{fn})
		if err != nil {
			log.Fatalf("error: could not read library dependency information of %q: %v", fn, err)
		}
		fmt.Println(fn, ":")
		for _, lib := range lst {
			fmt.Println("\t", lib)
		}
	}
	log.Println("Consumer has been finished")
}

func getElfFiles(dir string) chan string {
	ch := make(chan string)

	go func() {
		err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				log.Printf("error: %v\n", err)
				return err
			}
			if info.IsDir() {
				return nil // NOOP
			} else if info.Mode()&os.ModeSymlink != 0 {
				rp, err := realpath(path)
				if err != nil {
					log.Printf("error: could not get real path of %q: %v\n", path, err)
					return nil
				}
				path = rp
			}

			f, err := elf.Open(path)
			if err != nil {
				return nil
			}
			f.Close()
			ch <- path

			return nil
		})
		if err != nil {
			log.Fatalf("error: cannot start walk through file system: %v", err)
		}
		log.Println("Producer has been finished")
		close(ch)
	}()

	return ch
}

func realpath(symlink string) (string, error) {
	rp := ""
	p, err := os.Readlink(symlink)
	if err != nil {
		return "", fmt.Errorf("could not read symbolic link: %v", err)
	}
	if !filepath.IsAbs(p) {
		s := fmt.Sprintf("%s/%s", filepath.Dir(symlink), p)
		abs, err := filepath.Abs(s)
		if err != nil {
			return "", fmt.Errorf("cannot read absolute symlink %q: %v", s, err)
		}
		rp = abs
	} else {
		rp = p
	}

	fi, err := os.Lstat(rp)
	if err != nil {
		return "", fmt.Errorf("could not stat file %q: %v", rp, err)
	}

	if fi.Mode()&os.ModeSymlink != 0 {
		rp, err := realpath(rp)
		if err != nil {
			return "", fmt.Errorf("recursive call to realpath failed: %v", err)
		}
		return rp, nil
	}
	return rp, nil
}
