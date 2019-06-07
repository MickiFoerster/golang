package main

import (
	"log"
	"os"
	"text/template"
)

const Ccode = `
#include <unistd.h>

int main(int argc, char *argv[]) {
  pid_t pid = fork();
	if (pid == -1 ) {
		perror("fork failed");
		exit(EXIT_FAILURE);
	} else if(pid == 0) { // child process
	  char *const args[] = {
		{{range $arg := .}} "{{$arg}}", {{end}}
		  NULL
		}
	} else { // parent
	  waitpid(pid, NULL, 0);
	}

  return 0;
}
`

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("syntax error: At least one parameter must be provided. For example:\n%s ls -l -t", os.Args[0])
	}

	t := template.Must(template.New("Ccode").Parse(Ccode))
	t.Execute(os.Stdout, os.Args[1:])
}
