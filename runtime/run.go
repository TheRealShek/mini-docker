package runtime

import (
	"os"
	"os/exec"
	"syscall"
)

func Run() error {
	// prepare a child thread
	cmd := exec.Command("/bin/sh")

	// passing root settings to the child processes using a linux kernel call
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWPID | // container gets own pid table
			syscall.CLONE_NEWUTS | // container gets own hostname
			syscall.CLONE_NEWNS, // container gets own mount table
	}

	// Connect child process stdio to current terminal.
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
