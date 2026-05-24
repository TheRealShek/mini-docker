package runtime

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCheckDirectoryExists(t *testing.T) {
	// t.TempDir() auto-cleaned after test ends
	tmp := t.TempDir()

	// create a file inside tmp to test the "is file" case
	tmpFile := filepath.Join(tmp, "testfile.txt")
	os.WriteFile(tmpFile, []byte("data"), 0644)

	tests := []struct {
		name    string
		path    string
		wantOk  bool
		wantErr bool
	}{
		{
			name:    "valid directory",
			path:    tmp,
			wantOk:  true,
			wantErr: false,
		},
		{
			name:    "path is a file not dir",
			path:    tmpFile,
			wantOk:  false,
			wantErr: false,
		},
		{
			name:    "non existent path",
			path:    "/this/does/not/exist",
			wantOk:  false,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := checkDirectoryExists(tt.path)

			if (err != nil) != tt.wantErr {
				t.Fatalf("error = %v, wantErr = %v", err, tt.wantErr)
			}
			if got != tt.wantOk {
				t.Errorf("got = %v, want = %v", got, tt.wantOk)
			}
		})
	}
}

func TestMakeMultipleDirectories(t *testing.T) {
	tmp := t.TempDir()

	tests := []struct {
		name    string
		paths   []string
		wantErr bool
	}{
		{
			name:    "single directory",
			paths:   []string{filepath.Join(tmp, "one")},
			wantErr: false,
		},
		{
			name: "multiple directories",
			paths: []string{
				filepath.Join(tmp, "two"),
				filepath.Join(tmp, "three"),
			},
			wantErr: false,
		},
		{
			name:    "nested directory",
			paths:   []string{filepath.Join(tmp, "a", "b", "c")},
			wantErr: false,
		},
		{
			name:    "empty input",
			paths:   []string{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := makeMultipleDirectories(tt.paths...)

			if (err != nil) != tt.wantErr {
				t.Fatalf("error = %v, wantErr = %v", err, tt.wantErr)
			}

			// verify dirs actually exist on disk
			for _, p := range tt.paths {
				ok, err := checkDirectoryExists(p)
				if err != nil || !ok {
					t.Errorf("directory not created: %s", p)
				}
			}
		})
	}
}
