package buildins

import (
	"fmt"
	"os"
	"strings"

	"github.com/codecrafters-io/shell-starter-go/pkg/utils"
)

func exitCMD(args []string) error {
	os.Exit(0)
	return nil
}

func echoCMD(args []string) error {
	fmt.Println(strings.Join(args, " "))
	return nil
}

func typeCMD(args []string) error {
	var program string
	if args != nil {
		program = args[0]
	}

	if !IsBuiltin(program) {
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

func pwdCMD(args []string) error {
	// Get the working dir
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	fmt.Println(cwd)
	return nil
}
