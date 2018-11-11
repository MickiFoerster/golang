package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"io"
)

func main() {
	go server()
}

func server() {
	// Start to listen on port 8080
	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatalln(err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		// If connection is established start thread to handle it
		go handle(conn)
	}
}

func handle(conn net.Conn) {
	p := make([]byte, 1024)
	for {
		var n int
		n, err := conn.Read(p)
		fmt.Printf("read %d bytes\n", n)
		if err != nil { 
			conn.Close()
			if err == io.EOF {
				log.Println("Found EOF")
				os.Exit(0)
			}
			log.Fatalln(err) 
		}

		if n > 0 {
			for i := 0; i < n; i++ {
				fmt.Printf("%c", p[i])
				fmt.Fprintf(conn, "%c", p[i])
			}
		}
	}
}
