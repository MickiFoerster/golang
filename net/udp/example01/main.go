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
		runUDPServer(*address)
	case "C":
		runUDPClient(*address)
	}
}

func runUDPClient(address string) error {
	conn, err := net.Dial("udp", address)
	if err != nil {
		return err
	}
	defer conn.Close()

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Enter text to send:>")
	for scanner.Scan() {
		fmt.Println("Writing ", scanner.Text())
		conn.Write(scanner.Bytes())

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

func runUDPServer(address string) error {
	packetconn, err := net.ListenPacket("udp", address)
	if err != nil {
		return err
	}
	defer packetconn.Close()

	buffer := make([]byte, 4096)
	for {
		n, addr, err := packetconn.ReadFrom(buffer)
		if n > 0 {
			fmt.Printf("Received %q from address %s\n", string(buffer[:n]), addr)
		}
		if err != nil {
			fmt.Printf("error while reading packet: %v\n", err)
		}

		// Signal back to sender
		n, err = packetconn.WriteTo(
			[]byte(
				fmt.Sprintf("ACK for message %s %s",
					addr,
					string(buffer[:n]),
				),
			),
			addr,
		)
		if err != nil {
			fmt.Printf("error while sending response: %v\n", err)
		}
	}
}
