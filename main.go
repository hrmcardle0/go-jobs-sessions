package main

import (
	"bytes"
	"fmt"
	"github.com/ryboe/q"
	"golang.org/x/sys/unix"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
)

// hold info about our struct
type Proc struct {
	pid  int // process id
	ppid int // parent id
	pgid int // process group id
	psid int //session id
	tgid int
	tty  int // controlling terminal, converted to int based on bits
}

func main() {

	// get process id
	pid := os.Getpid()

	// get parent process id
	ppid := os.Getppid()

	// get process group (job) id
	pgid, err := syscall.Getpgid(pid)
	if err != nil {
		log.Fatalf("Error: %s", err)
	}

	// get session id
	psid, err := unix.Getsid(pid)
	if err != nil {
		log.Fatalf("Error: %s", err)
	}

	// get thread group Id, the foreground process group ID

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	runCmd := exec.Command("bash", "-c", fmt.Sprintf("cat /proc/%d/stat | cut -d ' ' -f 8", pid))

	runCmd.Stdout = &stdout
	runCmd.Stderr = &stderr
	err = runCmd.Run()
	if err != nil {
		q.Q(runCmd)
		log.Fatalf("Error: %s", err)
	}

	tgid, err := strconv.Atoi(strings.TrimRight(stdout.String(), "\n"))

	if err != nil {
		q.Q(runCmd)
		log.Fatalf("Error: %s", err)
	}
	q.Q(stdout)

	myProc := Proc{
		pid:  pid,
		ppid: ppid,
		pgid: pgid,
		psid: psid,
		tgid: tgid,
	}

	q.Q(myProc)

}
