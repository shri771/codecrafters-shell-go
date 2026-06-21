package buildins

import (
	"errors"
	"io"
	"strings"
	"sync"
)

type parsedCommand struct {
	program string
	args    []string
}

func ExecuteLine(line string, stdin io.Reader, stdout, stderr io.Writer) error {
	line = strings.TrimSpace(line)
	if line == "" {
		return nil
	}

	segments := strings.Split(line, "|")
	if len(segments) == 1 {
		command, err := parseCommand(segments[0])
		if err != nil {
			return err
		}
		return executeCommand(command, stdin, stdout, stderr, false)
	}
	commands := make([]parsedCommand, 0, len(segments))
	for _, segment := range segments {
		command, err := parseCommand(segment)
		if err != nil {
			return err
		}
		commands = append(commands, command)
	}

	return executePipeline(commands, stdin, stdout, stderr)
}

func parseCommand(segment string) (parsedCommand, error) {
	parts := strings.Fields(strings.TrimSpace(segment))
	if len(parts) == 0 {
		return parsedCommand{}, errors.New("empty command in pipeline")
	}

	return parsedCommand{
		program: parts[0],
		args:    parts[1:],
	}, nil
}

func executeCommand(
	command parsedCommand,
	stdin io.Reader,
	stdout, stderr io.Writer,
	inPipeline bool,
) error {
	if builtin, ok := GetCommand(command.program); ok {
		if inPipeline && (command.program == "cd" || command.program == "exit") {
			return nil
		}
		return builtin.Callback(command.args, stdin, stdout, stderr)
	}

	return runProgramWithIO(
		command.program,
		command.args,
		stdin,
		stdout,
		stderr,
		!inPipeline,
	)
}

func executePipeline(
	commands []parsedCommand,
	stdin io.Reader,
	stdout, stderr io.Writer,
) error {
	readers := make([]*io.PipeReader, len(commands)-1)
	writers := make([]*io.PipeWriter, len(commands)-1)
	for index := range readers {
		readers[index], writers[index] = io.Pipe()
	}

	var waitGroup sync.WaitGroup
	errorsChannel := make(chan error, len(commands))

	for index, command := range commands {
		commandInput := stdin
		if index > 0 {
			commandInput = readers[index-1]
		}

		commandOutput := stdout
		if index < len(commands)-1 {
			commandOutput = writers[index]
		}

		waitGroup.Add(1)
		go func(
			index int,
			command parsedCommand,
			commandInput io.Reader,
			commandOutput io.Writer,
		) {
			defer waitGroup.Done()

			err := executeCommand(
				command,
				commandInput,
				commandOutput,
				stderr,
				true,
			)
			if errors.Is(err, ErrExit) || errors.Is(err, io.ErrClosedPipe) {
				err = nil
			}

			if index < len(writers) {
				writers[index].CloseWithError(err)
			}
			if index > 0 {
				readers[index-1].CloseWithError(err)
			}
			errorsChannel <- err
		}(index, command, commandInput, commandOutput)
	}

	waitGroup.Wait()
	close(errorsChannel)

	for err := range errorsChannel {
		if err != nil {
			return err
		}
	}
	return nil
}
