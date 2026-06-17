package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/codecrafters-io/shell-starter-go/pkg/buildins"
	"github.com/codecrafters-io/shell-starter-go/pkg/utils"
)

func main() {
	reader := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("$ ")
		reader.Scan()

		line := reader.Text()

		// Sanitize Args
		// parts := strings.Fields(strings.ToLower(line))
		parts := strings.Fields(line)

		if len(parts) == 0 {
			continue
		}
		program := parts[0]
		var args []string

		if len(parts) > 1 {
			args = parts[1:]
		}

		cliCmd, ok := buildins.GetCommand(program)
		if ok {
			err := cliCmd.Callback(args)
			if err != nil {
				fmt.Printf("Unable to run program\n: %s", err)
			}
		} else {
			err := utils.RunProgram(program, args)
			if err != nil {
				fmt.Printf("Unable to run program\n: %s", err)
			}
		}

	}
}
