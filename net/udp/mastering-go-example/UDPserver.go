package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	if len(os.Args) == 1 {
		fmt.Fprintln(os.Stderr, "Please give port as parameter")
		return
	}

	port := ":" + os.Args[1]
	sock, err := net.ResolveUDPAddr("udp", port)
	if err != nil {
		log.Fatalln(err)
	}

	conn, err := net.ListenUDP("udp", sock)
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	var counter uint64 = 0
	b := make([]byte, 8)
	for {
		n, _, err := conn.ReadFrom(b)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error while reading: %v\n", err)
			continue
		}
		if n != 8 {
			log.Fatalf("Received number of bytes is %v but 8 were expected\n", n)
		}

		received := binary.LittleEndian.Uint64(b)
		if received != counter {
			fmt.Printf("received %v != %v expected\n", received, counter)
		} else {
			fmt.Printf("received %v is as expected\n", received)
		}
		counter++
	}
}
