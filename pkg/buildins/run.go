package buildins

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func RunProgram(program string, args []string) error {
	// Find path
	path, err := LookUpPath(program)
	if err != nil {
		return err
	}
	if path == "" {
		fmt.Printf("%s: command not found\n", program)
		return nil
	}

	_ = strings.Split(fmt.Sprintf("%s  %s", program, strings.Join(args, " ")), "&&")

	// Check if in which mode to run program
	if len(args) > 0 && args[len(args)-1] == "&" {
		args = args[:len(args)-1]
		err = runInBackground(program, args, path)
		if err != nil {
			return err
		}
	} else {
		err = runInForeground(program, args)
		if err != nil {
			return err
		}
	}

	return nil
}

func runInForeground(program string, args []string) error {
	cmd := exec.Command(program, args...)
	// Capture the out and ins
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout

	err := cmd.Run()
	if err != nil {
		return nil
	}
	return nil
}

func runInBackground(program string, args []string, path string) error {
	job := CreateJob()

	cmd := exec.Command(program, args...)

	// Capture the Std
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout

	// Set the Job
	job.SetCmdUsed(fmt.Sprintf("%s %s", program, strings.Join(args, " ")))

	err := cmd.Start()
	if err != nil {
		return nil
	}
	job.SetStatus(Running)
	job.SetPID(cmd.Process.Pid)
	DefaultJobStore.Add(job)

	// Display info
	fmt.Printf("[%d] %d\n", job.GetJobNumber(), job.GetPID())

	return nil
}

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

func SplitMultipleCMD(line string) []string {
	return strings.Split(line, "&&")
}

func ParseCommandLine(line string) (string, []string) {
	parts := strings.Fields(line)

	if len(parts) == 0 {
		return "", nil
	}

	if len(parts) == 1 {
		return parts[0], nil
	}
	return parts[0], parts[1:]
}
