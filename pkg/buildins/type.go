package buildins

import "io"

type CliCommand struct {
	Name     string
	Callback func([]string, io.Reader, io.Writer, io.Writer) error
}

func GetCommands() map[string]CliCommand {
	return map[string]CliCommand{
		"echo": {
			Name:     "echo",
			Callback: echoCMD,
		},
		"exit": {
			Name:     "exit",
			Callback: exitCMD,
		},
		"type": {
			Name:     "type",
			Callback: typeCMD,
		},
		"pwd": {
			Name:     "pwd",
			Callback: pwdCMD,
		},
		"cd": {
			Name:     "cd",
			Callback: cdCMD,
		},
		"jobs": {
			Name:     "jobs",
			Callback: jobsCMD,
		},
	}
}

func GetCommand(name string) (CliCommand, bool) {
	cmd, ok := GetCommands()[name]
	return cmd, ok
}

func IsBuiltin(name string) bool {
	_, ok := GetCommand(name)
	return ok
}
