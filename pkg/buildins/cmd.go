package buildins

import (
	"fmt"
	"os"
	"path/filepath"
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
	availableCmd := GetCommands()
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

func pwdCMD(args []string) error {
	// Get the working dir
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	fmt.Println(cwd)
	return nil
}

func cdCMD(args []string) error {
	var path string
	if len(args) != 1 {
		fmt.Println("The argument should be exactly one")
		return nil
	}

	argPath := args[0]

	if argPath == "~" || strings.HasPrefix(argPath, "~/") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return err
		}

		if argPath == "~" {
			path = homeDir
		} else {
			path = filepath.Join(homeDir, strings.TrimPrefix(argPath, "~/"))
		}
	} else if strings.HasPrefix(argPath, "/") {
		path = argPath
	} else {
		cleanedPath := strings.TrimPrefix(argPath, "./")

		cwd, err := os.Getwd()
		if err != nil {
			return err
		}

		path = filepath.Join(cwd, cleanedPath)
	}

	// Check if the path exist's
	fileInfo, err := os.Stat(path)
	if err != nil {
		fmt.Printf("cd: %s: No such file or directory\n", path)
		return nil
	}

	if fileInfo.IsDir() {
		err := os.Chdir(path)
		if err != nil {
			return err
		}
	} else {
		fmt.Printf("cd: %s: No such file or directory\n", path)
		return nil
	}
	return nil
}
