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

	// Check if in which mode to run program
	if strings.HasSuffix(fmt.Sprintf("%s %s", program, strings.Join(args, " ")), "&") {
		err = runInBackground(program, args, path)
		if err != nil {
			return err
		}
	} else {
		err = runInForeground(path, args)
		if err != nil {
			return err
		}
	}

	return nil
}

func runInForeground(path string, args []string) error {
	cmd := exec.Command(path, args...)
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
	// Set the Job
	job.SetCmdUsed(fmt.Sprintf("%s  %s", program, strings.Join(args, " ")))

	err := cmd.Start()
	if err != nil {
		return nil
	}
	job.SetStatus(Running)
	job.SetPID(cmd.Process.Pid)
	DefaultJobStore.Add(job)
	// job.SetJobNumber(DefaultJobStore.RunningCount())
	job.SetJobNumber(len(DefaultJobStore.jobs))

	// Display info
	fmt.Printf("[%d] %d\n", job.GetJobNumber(), job.GetPID())

	// Wait for the program to finish running
	go func() {
		err = cmd.Wait()
		if err != nil {
			job.SetStatus(Failed)
		} else {
			job.SetStatus(Done)
		}
	}()

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

// Run the program
// if err != nil {
// 	if exitErr, ok := err.(*exec.ExitError); ok {
// 		status = exitErr.ExitCode()
// 	} else {
// 		// Program failed to start, not just exited non-zero.
// 		status = 1
// 	}
// }

// status := 0
