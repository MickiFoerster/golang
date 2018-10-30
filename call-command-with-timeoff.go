package main

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"time"
)

var ready chan struct{}

func main() {
	ready = make(chan struct{})
	fmt.Println("Start 'ls -l'")
	go lscmd()
	<- ready
	fmt.Println("Start 'sleep 5'")
	go sleep()
	select {
	case <-time.After(2 * time.Second):
		fmt.Println("timeout")
	case <-ready:
		fmt.Println("finished")
	}
}

func lscmd() {
	var stdout, stderr bytes.Buffer
  cmd := exec.Command("ls", "-l")
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		log.Fatal("error: command execution failed: %s", err)
	}
	if stdout.Len() > 0 { 
		fmt.Printf("stdout: %s\n", stdout.String())
	}
	if stderr.Len() > 0 { 
		fmt.Printf("stderr: '%s'\n", stderr.String())
	}
	ready <- struct{}{}
}

func sleep() {
	cmd := exec.Command("sleep", "5")
	err := cmd.Run()
	if err != nil {
		log.Fatal("error: %s", err)
	}
	ready <- struct{}{}
}
