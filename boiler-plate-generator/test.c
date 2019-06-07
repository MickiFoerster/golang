#include <errno.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/types.h>
#include <sys/wait.h>
#include <unistd.h>

int main(int argc, char *argv[]) {
  int filedes[2];

  pid_t pid = fork();
	if (pid == -1 ) {
		perror("fork failed");
		exit(EXIT_FAILURE);
	} else if(pid == 0) { // child process
          int rc;
          const char path[] = "/bin/ls";
          char *const args[] = {"ls", "-l", "-h", "/boot", NULL};

          // Reroute stdout to pipe
          while (dup2(filedes[1], STDOUT_FILENO) == -1 && errno == EINTR)
            ;
          close(filedes[0]); // close exit for child
          close(filedes[1]); // close entrance of pipe
          rc = execv(path, args);
          fprintf(stderr, "child: may not happen");
          if (rc < 0) {
            fprintf(stderr, "execv failed: %s\n", strerror(errno));
          }
        } else { // parent
          fprintf(stderr, "parent: %d\n", pid);
          close(filedes[1]); // close entrance of pipe in parent process
          // Now parent process reads from exit of pipe
          char buf[4096];
          for (;;) {
            fprintf(stderr, "call read()\n");
            ssize_t n = read(filedes[0], buf, sizeof(buf));
            fprintf(stderr, "back from read(): %ld bytes read\n", n);
            if (n == -1) {
              if (errno == EINTR || errno == EAGAIN) {
                continue;
              } else {
                perror("read");
                waitpid(pid, NULL, 0);
                close(filedes[0]); // close pipe exit for parent
                exit(EXIT_FAILURE);
              }
            } else if (n == 0) {
              break; // EOF
            } else {
              fprintf(stdout, "catched from child process: %s\n", buf);
            }
          }
          waitpid(pid, NULL, 0);
          close(filedes[0]); // close pipe exit for parent
        }

  return 0;
}

