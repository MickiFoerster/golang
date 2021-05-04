package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", ":12345")
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()
	fmt.Println("Servers listens on", listener.Addr())

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
		}
		go handle(conn)

	}
}

func handle(conn net.Conn) {
	conndetails := fmt.Sprintf("Connection details:  local=%v <-> remote=%v", conn.LocalAddr(), conn.RemoteAddr())
	fmt.Println(conndetails, "\n")

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		ln := scanner.Text()
		fmt.Println("client:", ln)
		fmt.Fprintln(conn, fmt.Sprintln("server:", ln))
	}

	err := conn.Close()
	if err != nil {
		log.Printf("Could not close connection (%v)", conndetails)
	}
	fmt.Printf("Connection closed (%v)\n", conndetails)
}
