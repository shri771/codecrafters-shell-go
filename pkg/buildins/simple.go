package buildins

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

var ErrExit = errors.New("exit shell")

func exitCMD(args []string, stdin io.Reader, stdout, stderr io.Writer) error {
	return ErrExit
}

func echoCMD(args []string, stdin io.Reader, stdout, stderr io.Writer) error {
	_, err := fmt.Fprintln(stdout, strings.Join(args, " "))
	return err
}

func typeCMD(args []string, stdin io.Reader, stdout, stderr io.Writer) error {
	if len(args) == 0 {
		return nil
	}
	program := args[0]

	if !IsBuiltin(program) {
		path, err := LookUpPath(program)
		if err != nil {
			return err
		}
		if path == "" {
			fmt.Fprintf(stdout, "%s not found\n", program)
		} else {
			fmt.Fprintf(stdout, "%s is %s\n", program, path)
		}
	} else {
		fmt.Fprintf(stdout, "%s is a shell builtin\n", program)
	}

	return nil
}

func pwdCMD(args []string, stdin io.Reader, stdout, stderr io.Writer) error {
	// Get the working dir
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	_, err = fmt.Fprintln(stdout, cwd)
	return err
}
