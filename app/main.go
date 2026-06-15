package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var validCmd = []string{"echo", "cd"}

type cliCommand struct {
	name     string
	callback func(string) error
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
			cliCmd.callback(strings.Join(args, " "))
		} else {
			fmt.Printf("%s: command not found\n", program)
		}

	}
}

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"cd": {
			name: "cd",
			callback: func(string) error {
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
	}
}

func exit(arg string) error {
	os.Exit(0)
	return nil
}

func echo(arg string) error {
	fmt.Println(arg)
	return nil
}
