package runtime

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

// Parent Side
func Run() error {
	fmt.Println("parent: spawning container...")
	// relaunch same binary inside new namespaces
	cmd := exec.Command("/proc/self/exe")
	cmd.Env = append(os.Environ(), "CONTAINER_INIT=1")

	// Connect child process stdio to current terminal.
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	// passing root settings to the child processes using a linux kernel call
	/*
		namespace is assigned by kernel during clone()/fork().
		you can't join a new namespace after process already started (without unshare() syscall).
	*/

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWPID | // container gets own pid table
			syscall.CLONE_NEWUTS | // container gets own hostname
			syscall.CLONE_NEWNS, // container gets own mount table
	}
	return cmd.Run()
}

// Child Side
func ContainerInit() {
	fmt.Printf("child: inside init, pid=%d \n", os.Getpid())

	// prepare a child thread
	// Three args: path, argv (first element is program name), env.
	syscall.Exec("/bin/sh", []string{"/bin/sh"}, os.Environ())

}
