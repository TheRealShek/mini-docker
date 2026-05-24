package runtime

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"syscall"
)

var (
	// Test seam: keep the rootfs path configurable so unit tests can point it at a temp directory.
	alpineRootfsPath = "/tmp/alpine-rootfs"
	// Test seam: let tests observe the built command without launching a child process.
	runCommandFn = func(cmd *exec.Cmd) error {
		return cmd.Run()
	}
)

// Parent Side
func Run() error {
	ok, err := checkDirectoryExists(alpineRootfsPath)
	if ok {
		log.Println("alpine-rootfs exists")
	} else {
		log.Fatal(err)
	}

	fmt.Println("parent: spawning container...")
	return runCommandFn(newParentCommand())
}

func newParentCommand() *exec.Cmd {
	// Separated from Run so tests can verify the command wiring without executing it.
	// This creates a NEW process running the SAME binary again.
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
	return cmd
}

// Child Side
/*
lower/  -> (alpine-rootfs) internal base image layer
upper/  -> internal writable container layer
merged/ -> filesystem exposed to container
*/
func ContainerInit() error {
	// MS_PRIVATE -> "Stop sharing mount events."
	// MS_REC -> "Apply recursively to all submounts."
	err := syscall.Mount(
		"",
		"/",
		"",
		syscall.MS_PRIVATE|syscall.MS_REC,
		"",
	)
	if err != nil {
		return err
	}
	fmt.Printf("child: inside init, pid=%d \n", os.Getpid())

	err = makeMultipleDirectories(
		"/tmp/container/upper",
		"/tmp/container/work",
		"/tmp/container/merged",
	)
	if err != nil {
		return err
	}

	// Combine lower + upper layers to merged filesystem
	err = syscall.Mount(
		"overlay",
		"/tmp/container/merged",
		"overlay",
		0,
		"lowerdir=/tmp/alpine-rootfs,upperdir=/tmp/container/upper,workdir=/tmp/container/work",
	)
	if err != nil {
		return err
	}

	pivotRoot("/tmp/container/merged")

	// Mount the Linux proc filesystem inside the container.
	// /proc is not a normal directory on disk.
	// The Linux kernel generates its contents dynamically.
	err = syscall.Mount(
		"proc",
		"/proc",
		"proc",
		0,
		"",
	)
	if err != nil {
		return err
	}

	// syscall.Exec() replaces the CURRENT child runtime process.
	// Three args: path, argv (first element is program name), env.
	return syscall.Exec(
		"/bin/sh",
		[]string{"/bin/sh"},
		os.Environ(),
	)
}
