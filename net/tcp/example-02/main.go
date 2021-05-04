package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

func main() {
	op := flag.String("type", "", "server (s) or client (c) ?")
	address := flag.String("addr", ":1234", "address? host:port")
	flag.Parse()

	if len(*op) == 0 || len(*address) == 0 {
		fmt.Fprintln(os.Stderr, "syntax error: Use -h for help")
		os.Exit(1)
	}
	fmt.Println(*op)
	fmt.Println(*address)

	switch strings.ToUpper(*op) {
	case "S":
		runServer(*address)
	case "C":
		runClient(*address)
	}
}

func runClient(address string) error {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return err
	}
	defer conn.Close()

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Enter text to send:>")
	for scanner.Scan() {
		fmt.Println("Writing ", scanner.Text())
		conn.Write(append(scanner.Bytes(), '\r'))

		conn.SetReadDeadline(time.Now().Add(5 * time.Second))
		buffer := make([]byte, 4096)

		n, err := conn.Read(buffer)
		if err != nil {
			if err == io.EOF {
				log.Println("Connection was closed")
				return nil
			}
			log.Fatal(err)
		}
		fmt.Println(string(buffer[:n]))
	}

	return scanner.Err()
}

func runServer(address string) error {
	l, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	defer l.Close()
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Printf("error while accepting client connection: %v\n", err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)
	for {
		conn.SetDeadline(time.Now().Add(5 * time.Second))
		line, err := reader.ReadString('\r')
		if err != nil {
			if err == io.EOF {
				log.Println("Connection closed")
				return
			}
		}
		fmt.Printf("Received %s from address %s\n", line[:len(line)-1], conn.RemoteAddr())
		writer.WriteString("Message received ...")
		writer.Flush()
	}
}
