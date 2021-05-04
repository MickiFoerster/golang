package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:12345")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	fmt.Fprintf(conn, "Hello server %v from client %v\r\n", conn.RemoteAddr(), conn.LocalAddr())
}
