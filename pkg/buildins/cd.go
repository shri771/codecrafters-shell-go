package buildins

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

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
