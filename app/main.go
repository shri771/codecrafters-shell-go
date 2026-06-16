package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/codecrafters-io/shell-starter-go/pkg/utils"
)

var validCmd = []string{"echo", "cd"}

type cliCommand struct {
	name     string
	callback func([]string) error
}

func main() {
	reader := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("$ ")
		reader.Scan()

		line := reader.Text()

		// Sanitize Args
		parts := strings.Fields(strings.ToLower(line))

		if len(parts) == 0 {
			continue
		}
		program := parts[0]
		var args []string

		if len(parts) != 0 {
			args = parts[1:]
		}

		availableCmd := getCommands()
		cliCmd, ok := availableCmd[program]
		if ok {
			err := cliCmd.callback(args)
			if err != nil {
				fmt.Printf("Unable to run program\n: %s", err)
			}
		} else {
			err := runProgram(program, args)
			if err != nil {
				fmt.Printf("Unable to run program\n: %s", err)
			}
		}

	}
}

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"cd": {
			name: "cd",
			callback: func([]string) error {
				return nil
			},
		},
		"echo": {
			name:     "echo",
			callback: echoCMD,
		},
		"exit": {
			name:     "exit",
			callback: exitCMD,
		},
		"type": {
			name:     "type",
			callback: typeCMD,
		},
	}
}

func exitCMD(args []string) error {
	os.Exit(0)
	return nil
}

func echoCMD(args []string) error {
	fmt.Println(strings.Join(args, " "))
	return nil
}

func typeCMD(args []string) error {
	availableCmd := getCommands()
	program := args[0]

	_, ok := availableCmd[program]
	if !ok {
		path, err := utils.LookUpPath(program)
		if err != nil {
			return err
		}
		if path == "" {
			fmt.Printf("%s not found\n", program)
		} else {
			fmt.Printf("%s is %s\n", program, path)
		}
	} else {
		fmt.Printf("%s is a shell builtin\n", program)
	}

	return nil
}

func runProgram(program string, args []string) error {
	// status := 0
	// Find path
	path, err := utils.LookUpPath(program)
	if err != nil {
		return err
	}
	if path == "" {
		fmt.Printf("%s not found\n", program)
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
