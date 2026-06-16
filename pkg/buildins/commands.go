package buildins

type CliCommand struct {
	Name     string
	Callback func([]string) error
}

func GetCommands() map[string]CliCommand {
	return map[string]CliCommand{
		"cd": {
			Name: "cd",
			Callback: func([]string) error {
				return nil
			},
		},
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
	}
}
