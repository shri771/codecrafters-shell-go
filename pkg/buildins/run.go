package buildins

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

func RunProgram(program string, args []string) error {
	return runProgramWithIO(program, args, os.Stdin, os.Stdout, os.Stderr, true)
}

func runProgramWithIO(
	program string,
	args []string,
	stdin io.Reader,
	stdout, stderr io.Writer,
	allowBackground bool,
) error {
	// Find path
	path, err := LookUpPath(program)
	if err != nil {
		return err
	}
	if path == "" {
		fmt.Fprintf(stderr, "%s: command not found\n", program)
		return nil
	}

	_ = strings.Split(fmt.Sprintf("%s  %s", program, strings.Join(args, " ")), "&&")

	// Check if in which mode to run program
	if allowBackground && len(args) > 0 && args[len(args)-1] == "&" {
		args = args[:len(args)-1]
		err = runInBackground(program, args, stdin, stdout, stderr)
		if err != nil {
			return err
		}
	} else {
		err = runInForeground(program, args, stdin, stdout, stderr)
		if err != nil {
			return err
		}
	}

	return nil
}

func runInForeground(
	program string,
	args []string,
	stdin io.Reader,
	stdout, stderr io.Writer,
) error {
	cmd := exec.Command(program, args...)
	// Capture the out and ins
	cmd.Stderr = stderr
	cmd.Stdin = stdin
	cmd.Stdout = stdout

	err := cmd.Run()
	if err != nil {
		return nil
	}
	return nil
}

func runInBackground(
	program string,
	args []string,
	stdin io.Reader,
	stdout, stderr io.Writer,
) error {
	job := CreateJob()

	cmd := exec.Command(program, args...)

	// Capture the Std
	cmd.Stderr = stderr
	cmd.Stdin = stdin
	cmd.Stdout = stdout

	// Set the Job
	job.SetCmdUsed(fmt.Sprintf("%s %s", program, strings.Join(args, " ")))
	job.SetCmd(cmd)

	err := cmd.Start()
	if err != nil {
		return nil
	}
	job.SetStatus(Running)
	job.SetPID(cmd.Process.Pid)
	DefaultJobStore.Add(job)

	// Display info
	fmt.Fprintf(stdout, "[%d] %d\n", job.GetJobNumber(), job.GetPID())

	// Wait for the process to finish in background
	go func() {
		waitErr := cmd.Wait()
		if waitErr == nil {
			job.SetStatus(Done)
		} else {
			job.SetStatus(Failed)
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
