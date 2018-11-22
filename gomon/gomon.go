package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
)

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c)
	go func() {
		for {
			sign := <-c
			fmt.Println("got signal:", sign)
		}
	}()

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "syntax error: %q <Go file to execute>\n", path.Base(os.Args[0]))
		os.Exit(1)
	}

	if !strings.HasSuffix(path.Ext(os.Args[1]), ".go") {
		fmt.Fprintf(os.Stderr, "The extension of %q is not '.go'\n", os.Args[1])
	}

	build(os.Args[1])

	done := make(chan struct{})
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				log.Println("event:", event)
				err = watcher.Add(os.Args[1])
				if err != nil {
					log.Fatal(err)
				}
				fmt.Println("Restart Go file", executeableName)
				go startCmd()
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	go startCmd()
	<-done
}

var cmd *exec.Cmd
var executeableName string

func build(goFile string) {
	buildcmd := exec.Command("go", "build", goFile)
	err := buildcmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	executeableName = path.Join(pwd, goFile[:len(goFile)-len(".go")])
	fmt.Println(executeableName)
}

func startCmd() {
	if cmd != nil {
		err := cmd.Process.Kill()
		if err != nil {
			log.Fatal(err)
		}
		time.Sleep(time.Second)
	}
	cmd = exec.Command(executeableName)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	go io.Copy(os.Stdout, stdout)

	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Command %q started with PID %d\n", executeableName, cmd.Process.Pid)

	if err := cmd.Wait(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Command finished")
}
