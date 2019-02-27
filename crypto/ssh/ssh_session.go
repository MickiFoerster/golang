package main

import (
	"fmt"
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
		User:            "pi",
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

	session, err := connection.NewSession()
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	if err := session.RequestPty("xterm", 80, 40, modes); err != nil {
		session.Close()
		log.Fatal("request for pseudo terminal failed: %s", err)
	}

	stdin, err := session.StdinPipe()
	if err != nil {
		log.Fatal("Unable to setup stdin for session: %v", err)
	}

	session.Stdout = os.Stdout
	session.Stderr = os.Stderr

	// Start remote shell
	err = session.Shell()
	if err != nil {
		log.Fatal(err)
	}

	// Send command by command
	cmds := []string{
		"hostname",
		"ls -l",
		"netstat -tulpn",
		"asdf",
		"ps -ef",
		"exit",
	}
	for _, cmd := range cmds {
		fmt.Println("############ Start of command execution:", cmd, "###############")
		_, err = fmt.Fprintf(stdin, "%s\n", cmd)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("############ End of command execution:", cmd, "###############")
	}

	// Wait for session to finish
	err = session.Wait()
	if err != nil {
		log.Fatal(err)
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
