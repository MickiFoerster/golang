package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/fatih/color"
)

var tplclient *template.Template
var tplserver *template.Template

type unixDomainSocket struct {
	UnixDomainSocketPath string
	BufferSize           uint64
}

func init() {
	tplclient = template.Must(template.ParseFiles("createClientServerCommunicatingWithUnixSockets-client.gotemplate"))
	tplserver = template.Must(template.ParseFiles("createClientServerCommunicatingWithUnixSockets-server.gotemplate"))
}

func main() {
	uds := unixDomainSocket{
		UnixDomainSocketPath: "/tmp/ASDF",
		BufferSize:           4096,
	}

	const clientSourceFilename = "uds-client.c"
	const serverSourceFilename = "uds-server.c"
	for _, fn := range []string{clientSourceFilename, serverSourceFilename} {
		f, err := os.Create(fn)
		if err != nil {
			log.Fatal(err)
		}

		switch fn {
		case clientSourceFilename:
			err = tplclient.Execute(f, uds)
		case serverSourceFilename:
			err = tplserver.Execute(f, uds)
		default:
			log.Fatal("unexpected file name appeared:", fn)
		}
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

	// Copy tmpfile content to original file
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
	args := []string{
		"-std=c11",
		"-Wall",
		"-Werror",
		"-pthread",
		"-o",
		strings.TrimSuffix(inputfile, filepath.Ext(inputfile)) + "-" + compiler,
		inputfile,
	}
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
