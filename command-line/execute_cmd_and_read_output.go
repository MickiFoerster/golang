package main

import (
	"fmt"
	"log"
	"os/exec"
)

func main() {
	gituser, err := exec.Command("git", "config", "--global", "user.name").Output()
	if err != nil {
		log.Fatal(err)
	}
	if gituser[len(gituser)-1] == '\n' {
		gituser = gituser[0 : len(gituser)-1]
	}
	gitemail, err := exec.Command("git", "config", "--global", "user.email").Output()
	if err != nil {
		log.Fatal(err)
	}
	if gitemail[len(gitemail)-1] == '\n' {
		gitemail = gitemail[0 : len(gitemail)-1]
	}
	fmt.Printf("Git configured user is %s (%s)\n", gituser, gitemail)
}
