package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
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

		parts := strings.Fields(strings.ToLower(line))

		program := parts[0]

		var args []string

		if len(parts) != 0 {
			args = parts[1:]
		}

		availableCmd := getCommands()

		cliCmd, ok := availableCmd[program]
		if ok {
			cliCmd.callback(args)
		} else {
			fmt.Printf("%s: command not found\n", program)
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
			callback: echo,
		},
		"exit": {
			name:     "exit",
			callback: exit,
		},
		"type": {
			name:     "type",
			callback: typeCMD,
		},
	}
}

func exit(args []string) error {
	os.Exit(0)
	return nil
}

func echo(args []string) error {
	fmt.Println(strings.Join(args, " "))
	return nil
}

func typeCMD(args []string) error {
	availableCmd := getCommands()
	arg := args[0]

	_, ok := availableCmd[arg]
	if !ok {
		path, err := exec.LookPath(arg)
		if err != nil {
			if errors.Is(err, exec.ErrNotFound) {
				fmt.Printf("%s not found\n", arg)
				return nil
			} else {
				return err
			}
		}
		fmt.Printf("%s is %s\n", arg, path)

	} else {
		fmt.Printf("%s is a shell builtin\n", arg)
	}

	return nil
}
