package buildins

import (
	"errors"
	"fmt"
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
	if len(segments) != 2 {
		return fmt.Errorf("only two-command pipelines are supported")
	}

	left, err := parseCommand(segments[0])
	if err != nil {
		return err
	}
	right, err := parseCommand(segments[1])
	if err != nil {
		return err
	}

	return executeTwoCommandPipeline(left, right, stdin, stdout, stderr)
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

func executeTwoCommandPipeline(
	left, right parsedCommand,
	stdin io.Reader,
	stdout, stderr io.Writer,
) error {
	reader, writer := io.Pipe()
	var waitGroup sync.WaitGroup
	errorsChannel := make(chan error, 2)

	waitGroup.Add(2)
	go func() {
		defer waitGroup.Done()
		err := executeCommand(left, stdin, writer, stderr, true)
		if errors.Is(err, ErrExit) || errors.Is(err, io.ErrClosedPipe) {
			err = nil
		}
		writer.CloseWithError(err)
		errorsChannel <- err
	}()
	go func() {
		defer waitGroup.Done()
		err := executeCommand(right, reader, stdout, stderr, true)
		if errors.Is(err, ErrExit) {
			err = nil
		}
		reader.CloseWithError(err)
		errorsChannel <- err
	}()

	waitGroup.Wait()
	close(errorsChannel)

	for err := range errorsChannel {
		if err != nil {
			return err
		}
	}
	return nil
}
