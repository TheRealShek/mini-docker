package runtime

import (
	"os"
	"os/exec"
	"syscall"
	"testing"
)

func TestRunBuildsParentCommand(t *testing.T) {
	oldRootfsPath := alpineRootfsPath
	oldRunCommandFn := runCommandFn
	t.Cleanup(func() {
		alpineRootfsPath = oldRootfsPath
		runCommandFn = oldRunCommandFn
	})

	alpineRootfsPath = t.TempDir()

	var captured *exec.Cmd
	runCommandFn = func(cmd *exec.Cmd) error {
		captured = cmd
		return nil
	}

	if err := Run(); err != nil {
		t.Fatalf("Run() error = %v, want nil", err)
	}

	if captured == nil {
		t.Fatal("Run() did not build a command")
	}

	if got, want := captured.Path, "/proc/self/exe"; got != want {
		t.Fatalf("command path = %q, want %q", got, want)
	}

	if len(captured.Args) == 0 || captured.Args[0] != "/proc/self/exe" {
		t.Fatalf("command args = %v, want first arg %q", captured.Args, "/proc/self/exe")
	}

	if got := captured.Env[len(captured.Env)-1]; got != "CONTAINER_INIT=1" {
		t.Fatalf("command env tail = %q, want %q", got, "CONTAINER_INIT=1")
	}

	if captured.Stdin != os.Stdin {
		t.Fatal("command stdin was not wired to os.Stdin")
	}
	if captured.Stdout != os.Stdout {
		t.Fatal("command stdout was not wired to os.Stdout")
	}
	if captured.Stderr != os.Stderr {
		t.Fatal("command stderr was not wired to os.Stderr")
	}

	if captured.SysProcAttr == nil {
		t.Fatal("command SysProcAttr was not configured")
	}

	wantFlags := uintptr(syscall.CLONE_NEWPID | syscall.CLONE_NEWUTS | syscall.CLONE_NEWNS)
	if captured.SysProcAttr.Cloneflags != wantFlags {
		t.Fatalf("command clone flags = %v, want %v", captured.SysProcAttr.Cloneflags, wantFlags)
	}
}
