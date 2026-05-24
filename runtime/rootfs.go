package runtime

import (
	"os"
	"path/filepath"
	"syscall"
)

var (
	// Test seams: these wrappers let unit tests validate the pivot-root sequence without invoking real syscalls.
	mkdirAllFn   = os.MkdirAll
	pivotRootFn  = syscall.PivotRoot
	chdirFn      = os.Chdir
	unmountFn    = syscall.Unmount
	removePathFn = os.Remove
)

/*
Make `merged` a valid mount point for `pivot_root`.
Switch root to `merged`, moving old root to `/.pivot_old`.
Unmount old root to isolate the container filesystem.

pivot_root job is simple: make merged become /, and cut the real host filesystem off.
*/
func pivotRoot(newRoot string) error {
	// kernal needs a directory to temporary store old /
	prevRoot := filepath.Join(newRoot, ".pivot_old")
	err := mkdirAllFn(prevRoot, 0700)
	if err != nil {
		return err
	}

	// arg1 = "/tmp/container/merged"  → becomes new /
	// arg2 = "/tmp/container/merged/.pivot_old"  → old host / lands here
	err = pivotRootFn(newRoot, prevRoot)
	if err != nil {
		return err
	}

	// fix working directory as it still points to old root internally
	err = chdirFn("/")
	if err != nil {
		return err
	}

	// detach old host root from the container
	// MNT_DETACH ->detach a filesystem from the directory even if it is currently busy.
	err = unmountFn(
		"/.pivot_old",
		syscall.MNT_DETACH,
	)
	if err != nil {
		return err
	}

	return removePathFn("./.pivot_old")
}
