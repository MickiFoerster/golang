package main

import (
	"bytes"
	"log"
	"os"
	"os/exec"
)

func main() {
	buf := `This 
is a Test 
over multiple 
lines.

Escape characters are \t also contained\n
\n

`
	cmd := exec.Command("cat")
	cmd.Stdin = bytes.NewBuffer([]byte(buf))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	if err := cmd.Wait(); err != nil {
		log.Fatal(err)
	}
}
