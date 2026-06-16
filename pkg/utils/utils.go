package utils

import (
	"errors"
	"fmt"
	"os"
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

func RunProgram(program string, args []string) error {
	// status := 0
	// Find path
	path, err := LookUpPath(program)
	if err != nil {
		return err
	}
	if path == "" {
		fmt.Printf("%s: command not found\n", program)
		return nil
	}

	// Run the program
	cmd := exec.Command(program, args...)
	// Capture the out and ins
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	err = cmd.Run()
	if err != nil {
		return nil
	}
	// if err != nil {
	// 	if exitErr, ok := err.(*exec.ExitError); ok {
	// 		status = exitErr.ExitCode()
	// 	} else {
	// 		// Program failed to start, not just exited non-zero.
	// 		status = 1
	// 	}
	// }

	return nil
}
