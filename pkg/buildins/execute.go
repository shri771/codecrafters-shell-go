package buildins

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
)

// redirection holds one parsed redirection operator.
type redirection struct {
	fd         int  // 1 = stdout, 2 = stderr
	appendMode bool // true for >>, false for >
	filename   string
}

type parsedCommand struct {
	program      string
	args         []string
	redirections []redirection
}

func ExecuteLine(line string, stdin io.Reader, stdout, stderr io.Writer) error {
	line = strings.TrimSpace(line)
	if line == "" {
		return nil
	}

	segments := splitPipeline(line)
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

// splitPipeline splits a command line on unquoted '|' characters.
func splitPipeline(line string) []string {
	var segments []string
	var current strings.Builder
	inSingle := false
	inDouble := false
	for i := 0; i < len(line); i++ {
		ch := line[i]
		switch {
		case ch == '\\' && !inSingle:
			// backslash escape: consume next char literally
			current.WriteByte(ch)
			if i+1 < len(line) {
				i++
				current.WriteByte(line[i])
			}
		case ch == '\'' && !inDouble:
			inSingle = !inSingle
			current.WriteByte(ch)
		case ch == '"' && !inSingle:
			inDouble = !inDouble
			current.WriteByte(ch)
		case ch == '|' && !inSingle && !inDouble:
			segments = append(segments, current.String())
			current.Reset()
		default:
			current.WriteByte(ch)
		}
	}
	segments = append(segments, current.String())
	return segments
}

// parseRedirOp recognises a redirection operator token.
// Returns (fd, appendMode, true) on match, or (0, false, false) otherwise.
func parseRedirOp(token string) (fd int, appendMode bool, ok bool) {
	switch token {
	case ">", "1>":
		return 1, false, true
	case ">>", "1>>":
		return 1, true, true
	case "2>":
		return 2, false, true
	case "2>>":
		return 2, true, true
	}
	return 0, false, false
}

// parseCommand tokenizes a single command segment, respecting shell quoting,
// and extracts any redirection operators.
func parseCommand(segment string) (parsedCommand, error) {
	tokens, err := shellTokenize(strings.TrimSpace(segment))
	if err != nil {
		return parsedCommand{}, err
	}
	if len(tokens) == 0 {
		return parsedCommand{}, errors.New("empty command in pipeline")
	}

	var args []string
	var redirections []redirection

	for i := 0; i < len(tokens); i++ {
		fd, appendMode, isRedir := parseRedirOp(tokens[i])
		if isRedir {
			if i+1 >= len(tokens) {
				return parsedCommand{}, fmt.Errorf("syntax error: no filename after %s", tokens[i])
			}
			redirections = append(redirections, redirection{
				fd:         fd,
				appendMode: appendMode,
				filename:   tokens[i+1],
			})
			i++ // skip the filename token
		} else {
			args = append(args, tokens[i])
		}
	}

	if len(args) == 0 {
		return parsedCommand{}, errors.New("empty command in pipeline")
	}

	return parsedCommand{
		program:      args[0],
		args:         args[1:],
		redirections: redirections,
	}, nil
}

// shellTokenize splits a shell command line into tokens, stripping quotes and
// handling backslash escapes.
func shellTokenize(s string) ([]string, error) {
	var tokens []string
	var current strings.Builder
	inToken := false

	for i := 0; i < len(s); i++ {
		ch := s[i]
		switch {
		case ch == '\\' && i+1 < len(s):
			// outside quotes: backslash escapes next character
			i++
			current.WriteByte(s[i])
			inToken = true
		case ch == '\'':
			// single-quoted: everything literal until closing '
			inToken = true
			i++
			for i < len(s) && s[i] != '\'' {
				current.WriteByte(s[i])
				i++
			}
			// i now points at closing ' (or end of string)
		case ch == '"':
			// double-quoted: backslash only escapes \ and "
			inToken = true
			i++
			for i < len(s) && s[i] != '"' {
				if s[i] == '\\' && i+1 < len(s) && (s[i+1] == '\\' || s[i+1] == '"') {
					i++
					current.WriteByte(s[i])
				} else {
					current.WriteByte(s[i])
				}
				i++
			}
			// i now points at closing " (or end of string)
		case ch == ' ' || ch == '\t':
			if inToken {
				tokens = append(tokens, current.String())
				current.Reset()
				inToken = false
			}
		default:
			current.WriteByte(ch)
			inToken = true
		}
	}

	if inToken {
		tokens = append(tokens, current.String())
	}

	return tokens, nil
}

// applyRedirections opens files for all redirections and returns updated
// stdout/stderr writers plus a cleanup function that closes all opened files.
func applyRedirections(
	redirections []redirection,
	stdout, stderr io.Writer,
) (io.Writer, io.Writer, func(), error) {
	var openFiles []*os.File

	cleanup := func() {
		for _, f := range openFiles {
			f.Close()
		}
	}

	for _, redir := range redirections {
		flags := os.O_WRONLY | os.O_CREATE
		if redir.appendMode {
			flags |= os.O_APPEND
		} else {
			flags |= os.O_TRUNC
		}
		f, err := os.OpenFile(redir.filename, flags, 0644)
		if err != nil {
			cleanup()
			return stdout, stderr, func() {}, err
		}
		openFiles = append(openFiles, f)
		switch redir.fd {
		case 1:
			stdout = f
		case 2:
			stderr = f
		}
	}

	return stdout, stderr, cleanup, nil
}

func executeCommand(
	command parsedCommand,
	stdin io.Reader,
	stdout, stderr io.Writer,
	inPipeline bool,
) error {
	// Apply any redirections for this command
	stdout, stderr, cleanup, err := applyRedirections(command.redirections, stdout, stderr)
	if err != nil {
		return err
	}
	defer cleanup()

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
