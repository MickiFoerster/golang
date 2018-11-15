package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

var wg sync.WaitGroup
var sessionOutput map[string]chan []byte

func main() {
	for _, host := range os.Args[1:] {
		fmt.Println("Connect to host", host)
		wg.Add(1)
		go connectToHost(host)
	}
	wg.Wait()
}

func connectToHost(host string) {
	defer wg.Done()
	conn, err := startSSHConnection(host)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	tmpfilename, err := createWindowForOutput()
	if err != nil {
		fmt.Println("Could not create window")
		return
	}
	//defer os.Remove(tmpfilename)

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

	tmpfile, err := os.Create(tmpfilename)
	if err != nil {
		fmt.Println("Could not open temporary file")
		return
	}
	stdout, err := session.StdoutPipe()
	if err != nil {
		fmt.Println("Could not redirect stdout to temporary file")
		return
	}
	go func() {
		n, err := io.Copy(tmpfile, stdout)
		if err != nil {
			fmt.Println("error: io.Copy failed:", err)
		}
		fmt.Printf("%d bytes copied\n", n)
		session.Close()
		tmpfile.Close()
		fmt.Println("copier finished")
	}()

	fmt.Println("Execute command")
	session.Run("ls -l")
	fmt.Println("Finished")
	time.Sleep(100 * time.Millisecond)
}

func createWindowForOutput() (string, error) {
	tmpfile, err := ioutil.TempFile("", "sshoutput")
	if err != nil {
		return "", err
	}
	tmpfile.Close()
	cmd := exec.Command("/usr/bin/xterm", "-hold", "-e", "tail", "-f", tmpfile.Name())
	err = cmd.Run()
	fmt.Println("xterm:", err)
	return tmpfile.Name(), nil
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

	fmt.Println(user)
	fmt.Println(hostname)
	fmt.Println(port)
	sshConfig := &ssh.ClientConfig{
		User:            user,
		Auth:            []ssh.AuthMethod{sshAgent()},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	hostPlusPort := fmt.Sprintf("%s:%s", hostname, port)
	fmt.Println("Try to connect to ", hostPlusPort)
	conn, err := ssh.Dial("tcp", hostPlusPort, sshConfig)
	if err != nil {
		return nil, fmt.Errorf("Could not connect to %q:", hostPlusPort, err)
	}
	fmt.Println("Successful connected to ", hostPlusPort)
	return conn, nil
}

func sshAgent() ssh.AuthMethod {
	sshAgent, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK"))
	if err != nil {
		log.Fatal(err)
	}
	return ssh.PublicKeysCallback(agent.NewClient(sshAgent).Signers)
}
