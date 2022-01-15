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
		fmt.Fprintln(os.Stderr, "Please give hostname:port as parameter")
		return
	}

	host_and_port := os.Args[1]
	sock, err := net.ResolveUDPAddr("udp", host_and_port)
	if err != nil {
		log.Fatalln(err)
	}

	conn, err := net.DialUDP("udp", nil, sock)
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	fmt.Printf("Connected via UDP to %s\n", conn.RemoteAddr())

	var counter uint64 = 0
	for {
		b := make([]byte, 8)
		binary.LittleEndian.PutUint64(b, counter)
		n, err := conn.Write(b)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error while sending %v: %v", counter, err)
			continue
		}
		fmt.Printf("%v bytes and counter %v has been sent\n", n, counter)

		counter++
		//time.Sleep(time.Millisecond)
	}
}
