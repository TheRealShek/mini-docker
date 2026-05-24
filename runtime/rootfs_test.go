package runtime

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"syscall"
	"testing"
)

var (
	errMkdirAll = errors.New("mkdir failed")
	errPivot    = errors.New("pivot failed")
	errChdir    = errors.New("chdir failed")
	errUnmount  = errors.New("unmount failed")
	errRemove   = errors.New("remove failed")
)

func stubRootfsOps(
	t *testing.T,
	mkdirAllStub func(string, os.FileMode) error,
	pivotRootStub func(string, string) error,
	chdirStub func(string) error,
	unmountStub func(string, int) error,
	removeStub func(string) error,
) {
	t.Helper()

	oldMkdirAll := mkdirAllFn
	oldPivotRoot := pivotRootFn
	oldChdir := chdirFn
	oldUnmount := unmountFn
	oldRemove := removePathFn

	t.Cleanup(func() {
		mkdirAllFn = oldMkdirAll
		pivotRootFn = oldPivotRoot
		chdirFn = oldChdir
		unmountFn = oldUnmount
		removePathFn = oldRemove
	})

	mkdirAllFn = mkdirAllStub
	pivotRootFn = pivotRootStub
	chdirFn = chdirStub
	unmountFn = unmountStub
	removePathFn = removeStub
}

func TestPivotRoot(t *testing.T) {
	tests := []struct {
		name       string
		mkdirErr   error
		pivotErr   error
		chdirErr   error
		unmountErr error
		removeErr  error
		wantErr    error
		wantCalls  []string
	}{
		{
			name:      "mkdir fails",
			mkdirErr:  errMkdirAll,
			wantErr:   errMkdirAll,
			wantCalls: []string{"mkdir"},
		},
		{
			name:      "pivot fails",
			pivotErr:  errPivot,
			wantErr:   errPivot,
			wantCalls: []string{"mkdir", "pivot"},
		},
		{
			name:      "chdir fails",
			chdirErr:  errChdir,
			wantErr:   errChdir,
			wantCalls: []string{"mkdir", "pivot", "chdir"},
		},
		{
			name:       "unmount fails",
			unmountErr: errUnmount,
			wantErr:    errUnmount,
			wantCalls:  []string{"mkdir", "pivot", "chdir", "unmount"},
		},
		{
			name:      "remove fails",
			removeErr: errRemove,
			wantErr:   errRemove,
			wantCalls: []string{"mkdir", "pivot", "chdir", "unmount", "remove"},
		},
		{
			name:      "success",
			wantCalls: []string{"mkdir", "pivot", "chdir", "unmount", "remove"},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			newRoot := filepath.Join(t.TempDir(), "merged")
			expectedPrevRoot := filepath.Join(newRoot, ".pivot_old")
			calls := make([]string, 0, 5)

			stubRootfsOps(
				t,
				func(path string, perm os.FileMode) error {
					calls = append(calls, "mkdir")
					if path != expectedPrevRoot {
						t.Fatalf("mkdirAll path = %q, want %q", path, expectedPrevRoot)
					}
					if perm != 0700 {
						t.Fatalf("mkdirAll perm = %v, want %v", perm, 0700)
					}
					return tt.mkdirErr
				},
				func(newRootArg, prevRootArg string) error {
					calls = append(calls, "pivot")
					if newRootArg != newRoot {
						t.Fatalf("pivotRoot newRoot = %q, want %q", newRootArg, newRoot)
					}
					if prevRootArg != expectedPrevRoot {
						t.Fatalf("pivotRoot prevRoot = %q, want %q", prevRootArg, expectedPrevRoot)
					}
					return tt.pivotErr
				},
				func(path string) error {
					calls = append(calls, "chdir")
					if path != "/" {
						t.Fatalf("chdir path = %q, want %q", path, "/")
					}
					return tt.chdirErr
				},
				func(target string, flags int) error {
					calls = append(calls, "unmount")
					if target != "/.pivot_old" {
						t.Fatalf("unmount target = %q, want %q", target, "/.pivot_old")
					}
					if flags != syscall.MNT_DETACH {
						t.Fatalf("unmount flags = %v, want %v", flags, syscall.MNT_DETACH)
					}
					return tt.unmountErr
				},
				func(path string) error {
					calls = append(calls, "remove")
					if path != "./.pivot_old" {
						t.Fatalf("remove path = %q, want %q", path, "./.pivot_old")
					}
					return tt.removeErr
				},
			)

			err := pivotRoot(newRoot)
			if err != tt.wantErr {
				t.Fatalf("pivotRoot() error = %v, want %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(calls, tt.wantCalls) {
				t.Fatalf("call sequence = %v, want %v", calls, tt.wantCalls)
			}
		})
	}
}
