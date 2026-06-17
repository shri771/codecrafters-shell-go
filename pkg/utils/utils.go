package utils

import (
	"strings"
)

// Clean args
func CleanArgs(line string) string {
	line = strings.ToLower(line)
	line = strings.TrimSpace(line)

	return line
}
