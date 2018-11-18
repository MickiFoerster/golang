package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/exec"
	"strings"
	"sync"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

var wg sync.WaitGroup
var sessionOutput map[string]chan []byte

func main() {
	if len(os.Args) > 1 {
		ch2b := make(chan string)
		var clients []chan string
		go func() {
			for cmd := range ch2b {
				fmt.Println("Entered command", cmd)
				fmt.Printf("Send command to %d clients\n", len(clients))
				for _, cch := range clients {
					fmt.Println("Broadcast", cmd)
					cch <- cmd
				}
			}
		}()
		for _, host := range os.Args[1:] {
			ch := make(chan string)
			clients = append(clients, ch)
			fmt.Println("DEBUG: number of clients:", len(clients))
			fmt.Println("Connect to host", host)
			wg.Add(1)
			go connectToHost(host, ch)
		}
		wg.Wait()
		reader := bufio.NewReader(os.Stdin)
		for {
			fmt.Print("ssh Broadcast>")
			cmd, err := reader.ReadString('\n')
			if err != nil {
				break
			}
			fmt.Println("Put into channel", cmd)
			ch2b <- cmd
		}
	} else {
		fmt.Printf("Give the host names or IP addresses as command line option:\n"+
			"%s host1 host2 ...\n", os.Args[0])
	}
}

func connectToHost(host string, ch chan string) {
	conn, err := startSSHConnection(host)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	tmpfile, err := ioutil.TempFile("", "sshoutput")
	if err != nil {
		fmt.Printf("error: Could create temporary file:", err)
		return
	}
	fmt.Println("Create temp file")
	go createWindowForOutput(tmpfile.Name())
	fmt.Println("temp file created")

	wg.Done() // send main ready signal

	for cmd := range ch {
		fmt.Println("Received command:", cmd)
		session, err := conn.NewSession()
		if err != nil {
			fmt.Printf("error:%s: Could not create SSH session", host)
			continue
		}

		modes := ssh.TerminalModes{
			ssh.ECHO:          0,     // disable echoing
			ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
			ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
		}

		if err := session.RequestPty("xterm", 80, 40, modes); err != nil {
			fmt.Printf("error:%s: Could not get PTY\n", host)
			continue
		}

		stdout, err := session.StdoutPipe()
		if err != nil {
			fmt.Println("Could not redirect stdout to temporary file")
			return
		}
		go func() {
			_, err := io.Copy(tmpfile, stdout)
			if err != nil {
				fmt.Println("error: io.Copy failed:", err)
			}
			session.Close()
		}()
		session.Run(cmd)
	}
	tmpfile.Close()
	os.Remove(tmpfile.Name())
}

func createWindowForOutput(tmpFilename string) {
	cmd := exec.Command("/usr/bin/xterm", "-hold", "-e", "tail", "-f", tmpFilename)
	err := cmd.Run()
	fmt.Println("xterm:", err)
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
