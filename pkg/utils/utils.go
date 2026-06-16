package utils

import (
	"errors"
	"os/exec"
)

// Find exec Path
func LookUpPath(program string) (string, error) {
	path, err := exec.LookPath(program)
	if err != nil {
		if errors.Is(err, exec.ErrNotFound) {
			return "", nil
		} else {
			return "", err
		}
	}
	return path, nil
}
