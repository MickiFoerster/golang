package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
)

func main() {
	cert, err := tls.LoadX509KeyPair("cert.pem", "cert-key.pem")
	if err != nil {
		log.Fatal(err)
	}
	cfg := &tls.Config{Certificates: []tls.Certificate{cert}}
	listener, err := tls.Listen("tcp", ":2000", cfg)
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()
	remhost := conn.RemoteAddr().String()
	fmt.Fprintln(conn, "Hello ", remhost)
	var buffer []byte
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Read %d bytes\n", n)
	if n > 0 {
		fmt.Fprintln(conn, buffer[:n])
	}
}
