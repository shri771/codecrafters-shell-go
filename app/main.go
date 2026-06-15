package main

import (
	"bufio"
	"fmt"
	"os"
)

var validCmd = []string{"echo", "cd"}

func main() {
	reader := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("$ ")
		reader.Scan()

		cmd := reader.Text()

		// Check for a valid cmd
		if cmd == "exit" {
			os.Exit(0)
		}
		var isValid bool
		for _, c := range validCmd {
			if c == cmd {
				isValid = true
			}
		}
		if !isValid {
			fmt.Printf("%s: command not found\n", cmd)
		}

	}
}
