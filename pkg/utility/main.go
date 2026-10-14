package utility

import (
	"bytes"
	"errors"
	"os"
	"strconv"
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
	// Any other error means something else went wrong
	// (e.g., permission denied, broken symlink, etc.)
	return false
}

// AppendPrettyFloat formats a float and inserts 'sep' every 3 digits
// before the decimal point without any heap allocations.
func AppendPrettyFloat(dst []byte, val float64, prec int, sep byte) []byte {
	// Stack buffer to format the raw float first
	var scratch [64]byte
	raw := strconv.AppendFloat(scratch[:0], val, 'f', prec, 64)

	// Locate sign and decimal point
	start := 0
	if len(raw) > 0 && (raw[0] == '-' || raw[0] == '+') {
		start = 1
	}

	dotIdx := bytes.IndexByte(raw, '.')
	if dotIdx == -1 {
		dotIdx = len(raw)
	}

	intDigits := dotIdx - start
	if intDigits <= 0 {
		return append(dst, raw...)
	}

	// Calculate how many separators are needed: (n - 1) / 3
	//numSeps := (intDigits - 1) / 3

	// Write prefix (optional +/- sign)
	//if start == 1 {
	//	dst = append(dst, raw[0])
	//}

	// Write integer digits with separators
	for i := range intDigits {
		// Insert separator if remaining digits are a multiple of 3
		rem := intDigits - i
		if i > 0 && rem%3 == 0 {
			dst = append(dst, sep)
		}
		dst = append(dst, raw[start+i])
	}

	// Append decimal portion (e.g., ".25")
	if dotIdx < len(raw) {
		dst = append(dst, raw[dotIdx:]...)
	}

	return dst
}
