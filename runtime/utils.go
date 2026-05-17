package runtime

import "os"

// returns (isDir, nil) if Stat succeeds; returns (false, err) if Stat fails
func checkDirectoryExists(path string) (bool, error) {
	info, err := os.Stat(path)
	if err == nil {
		return info.IsDir(), nil
	}
	return false, err
}

// creates multiple directories with the specified paths
func makeMultipleDirectories(paths ...string) error {
	for _, path := range paths {
		err := os.MkdirAll(path, 0755)
		if err != nil {
			return err
		}
	}
	return nil
}
