package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"

	"github.com/codecrafters-io/shell-starter-go/pkg/buildins"
)

func main() {
	reader := bufio.NewScanner(os.Stdin)
	// Running jobs
	for {
		buildins.ReapCompletedJobs()
		fmt.Print("$ ")
		if !reader.Scan() {
			break
		}

		err := buildins.ExecuteLine(reader.Text(), os.Stdin, os.Stdout, os.Stderr)
		if errors.Is(err, buildins.ErrExit) {
			return
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to run program: %s\n", err)
		}
	}
}
