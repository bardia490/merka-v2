package utility

import (
	"errors"
	"os"
)

// returns true if the file does exist in the relative [path]
func DoesFileExist(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true // File exists
	}
	if errors.Is(err, os.ErrNotExist) {
		return false // File does not exist
	}
	// Note: Any other error means something else went wrong
	// (e.g., permission denied, broken symlink, etc.)
	return false
}
