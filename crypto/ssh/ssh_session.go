package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"sync"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

var n sync.WaitGroup

func main() {
	for _, sshhost := range os.Args[1:] {
		n.Add(1)
		go startSSHConnection(sshhost)
	}
	n.Wait()
}

func startSSHConnection(host string) {
	sshConfig := &ssh.ClientConfig{
		User:            "root",
		Auth:            []ssh.AuthMethod{sshAgent()},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	hostPlusPort := fmt.Sprintf("%s:%d", host, 22)
	fmt.Println("Try to connect to ", hostPlusPort)
	connection, err := ssh.Dial("tcp", hostPlusPort, sshConfig)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Successful connected to ", hostPlusPort)

	cmds := []string{
		"hostname",
		"ls -l",
	}
	for _, cmd := range cmds {
		session, err := connection.NewSession()
		if err != nil {
			log.Fatal(err)
		}

		modes := ssh.TerminalModes{
			ssh.ECHO:          0,     // disable echoing
			ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
			ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
		}

		if err := session.RequestPty("xterm", 80, 40, modes); err != nil {
			session.Close()
			log.Fatal("request for pseudo terminal failed: %s", err)
		}

		stdout, err := session.StdoutPipe()
		if err != nil {
			log.Fatal("Unable to setup stdout for session: %v", err)
		}
		go io.Copy(os.Stdout, stdout)
		stderr, err := session.StderrPipe()
		if err != nil {
			log.Fatal("Unable to setup stderr for session: %v", err)
		}
		go io.Copy(os.Stdout, stderr)

		session.Setenv("PS1", "")
		session.Run(cmd)
		session.Close()
	}

	n.Done()
}

func sshAgent() ssh.AuthMethod {
	sshAgent, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK"))
	if err != nil {
		log.Fatal(err)
	}
	return ssh.PublicKeysCallback(agent.NewClient(sshAgent).Signers)
}

