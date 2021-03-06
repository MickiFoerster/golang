package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"text/template"

	"github.com/fatih/color"
)

var tpl *template.Template

type execvCall struct {
	Path string
	Args []string
}

func init() {
	tpl = template.Must(template.ParseFiles("createCommandLineCallCodeC.gotemplate"))
}

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("syntax error: At least one parameter must be provided. For example:\n%s ls -l -t", os.Args[0])
	}

	execv := execvCall{os.Args[1], os.Args[1:]}

	if _, err := os.Stat(execv.Path); os.IsNotExist(err) {
		log.Fatal(fmt.Sprintf("error: First argument has to be full path to existent file, %q does not exist.", execv.Path))
	}

	const fn = "main.c"
	f, err := os.Create(fn)
	if err != nil {
		log.Fatal(err)
	}

	err = tpl.Execute(f, execv)
	if err != nil {
		log.Fatal(err)
	}

	f.Close()
	fmt.Printf("C code has been written to %q.\n", fn)

	// Postprocess with clang-format
	applyClangFormat(fn)

	fmt.Println("Now, we test this code by using:")
	testOutputWithCompiler("gcc", fn)
	testOutputWithCompiler("clang", fn)
}

func applyClangFormat(fn string) {
	clangformat := exec.Command("clang-format", fn)
	stdout, err := clangformat.StdoutPipe()
	if err != nil {
		fmt.Println("Could not redirect stdout of clang-format", err)
		return
	}
	reader := bufio.NewReader(stdout)
	if err = clangformat.Start(); err != nil {
		fmt.Println("Could not start clang-format", err)
		return
	}

	tmpfile, err := ioutil.TempFile("", "clangformat")
	if err != nil {
		fmt.Println("could create temporary file for applying clang-format", err)
		return
	}
	defer os.Remove(tmpfile.Name())

	for {
		buf := make([]byte, 4096)
		n, err := reader.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("Could not read from stdout of clang-format", err)
			return
		}
		n, err = tmpfile.Write(buf[:n])
		if err != nil {
			fmt.Println("could not write to temporary file", err)
			return
		}
	}
	tmpfile.Close()

	if clangformat.Wait(); err != nil {
		fmt.Println("Wait() failed for clang-format", err)
		return
	}

	// Copy tmpfile content to main.c
	src, err := os.Open(tmpfile.Name())
	if err != nil {
		fmt.Println("could not open temporary file", err)
		return
	}
	defer src.Close()

	dst, err := os.Create(fn)
	if err != nil {
		fmt.Println("could not open target file", fn, err)
		return
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	if err != nil {
		fmt.Println("could not copy temporary file to", fn, err)
		return
	}
}

func testOutputWithCompiler(compiler string, inputfile string) {
	args := []string{"-std=c11", "-ggdb3", "-Wall", "-Werror", inputfile}
	s := fmt.Sprint(compiler)
	for _, arg := range args {
		s += fmt.Sprintf(" %s", arg)
	}
	s += fmt.Sprint(": ")

	cmd := exec.Command(compiler, args...)
	err := cmd.Run()
	if err != nil {
		fmt.Printf("%-40s", s)
		color.Red(" %10s ", "[failed]")
		fmt.Println(err)
	} else {
		fmt.Printf("%-40s", s)
		color.Green(" %10s ", "[OK]")
	}
}
