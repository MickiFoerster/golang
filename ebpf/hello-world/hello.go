package main

import (
	"C"

	bpf "github.com/aquasecurity/tracee/libbpfgo"
)
import (
	"fmt"
	"log"
	"os"
	"os/signal"
)

func main() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	b, err := bpf.NewModuleFromFile("hello.bpf.o")
	if err != nil {
		log.Fatal("error while loading BPF from file:", err)
	}
	defer b.Close()

	// Load program into kernel
	if err := b.BPFLoadObject().Error; err != nil {
		log.Fatal("error while loading program into kernel:", err)
	}

	// Get function hello and attach it to kprobe
	p, err := b.GetProgram("hello")
	if err != nil {
		log.Fatal("error while GetProgram: ", err)
	}

	_, err = p.AttachKprobe("__x64_sys_execve")
	if err != nil {
		log.Fatal("error while AttachKprobe(): ", err)
	}

	//bpf.TracePrint()
	<-sig
	fmt.Println("main is done")
}
