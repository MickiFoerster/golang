package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path"
	"strings"
	"sync"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

var sessionOutput = make(map[string]string)
var wg sync.WaitGroup
var logFile *os.File

func init() {
	logFile, err := os.Create(fmt.Sprintf("%s.log", path.Base(os.Args[0])))
	if err != nil {
		fmt.Println("Could not create log file:", err)
	}
	log.SetOutput(logFile)
}

func main() {
	if logFile != nil {
		defer logFile.Close()
	}
	if len(os.Args) > 1 {
		for _, host := range os.Args[1:] {
			wg.Add(1)
			log.Println("Connect to host", host)
			go connectToHost(host)
		}
		wg.Wait()
		for k, v := range sessionOutput {
			fmt.Println(k, v)
		}
	} else {
		fmt.Printf("Give the host names or IP addresses as command line option:\n"+
			"%s host1 host2 ...\n", os.Args[0])
	}
}

func connectToHost(host string) {
	conn, err := startSSHConnection(host)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		conn.Close()
		log.Println("connection closed")
		wg.Done()
	}()

	session, err := conn.NewSession()
	if err != nil {
		fmt.Printf("error:%s: Could not create SSH session", host)
		return
	}

	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	if err := session.RequestPty("xterm", 80, 40, modes); err != nil {
		fmt.Printf("error:%s: Could not get PTY\n", host)
		return
	}

	stdout, err := session.StdoutPipe()
	if err != nil {
		fmt.Println("Could not redirect stdout to temporary file")
		return
	}
	go func() {
		buffer := make([]byte, 4096)
		for {
			n, err := stdout.Read(buffer)
			log.Printf("DEBUG: read %d bytes: %X\n", n, string(buffer[:n]))
			if err != nil {
				if err == io.EOF {
					break
				}
				fmt.Println("error while reading remote stdout:", err)
				break
			}
			processSessionOutput(conn, string(buffer[:n]))
		}
		session.Close()
		log.Println("session closed")
	}()
	session.Run("hostname")
}

func processSessionOutput(conn ssh.Conn, output string) {
	curValue := sessionOutput[conn.RemoteAddr().String()]
	s := strings.TrimRight(output, "\r\n")
	sessionOutput[conn.RemoteAddr().String()] = curValue + s
}

func startSSHConnection(host string) (*ssh.Client, error) {
	s := strings.Split(host, "@")
	var hostname string
	var user string
	switch len(s) {
	case 1:
		user = os.Getenv("USER")
		hostname = s[0]
	case 2:
		user = s[0]
		hostname = s[1]
	}

	var port string
	s = strings.Split(hostname, ":")
	switch len(s) {
	case 1:
		port = "22"
		hostname = s[0]
	case 2:
		hostname = s[0]
		port = s[1]
	}

	log.Println("user:", user)
	log.Println("hostname:", hostname)
	log.Println("port:", port)
	sshConfig := &ssh.ClientConfig{
		User:            user,
		Auth:            []ssh.AuthMethod{sshAgent()},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	hostPlusPort := fmt.Sprintf("%s:%s", hostname, port)
	log.Println("Try to connect to ", hostPlusPort)
	conn, err := ssh.Dial("tcp", hostPlusPort, sshConfig)
	if err != nil {
		return nil, fmt.Errorf("Could not connect to %q:", hostPlusPort, err)
	}
	log.Println("Successful connected to ", hostPlusPort)
	return conn, nil
}

func sshAgent() ssh.AuthMethod {
	sshAgent, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK"))
	if err != nil {
		log.Fatal(err)
	}
	return ssh.PublicKeysCallback(agent.NewClient(sshAgent).Signers)
}
