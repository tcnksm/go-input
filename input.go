/*
go-input is package to interact with user input from command line interface.

http://github.com/tcnksm/go-input
*/
package input

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strconv"
)

var (
	defaultWriter = os.Stdout
	defaultReader = os.Stdin
)

// UI
type UI struct {
	// Writer is where output is written. For example query
	// to the user will be written. By default, it's os.Stdout.
	Writer io.Writer

	// Reader is source of input. By default, it's os.Stdin.
	Reader io.Reader
}

type Options struct {
	// Default is the value when no thing is innputted
	Default string

	Loop bool
}

func (i *UI) Select(query string, list []string, opts *Options) (string, error) {

	wr := i.Writer
	if wr == nil {
		wr = defaultWriter
	}

	rd := i.Reader
	if rd == nil {
		rd = defaultReader
	}

	defaultIndex := -1
	defaultVal := opts.Default
	if defaultVal != "" {
		for i, item := range list {
			if item == defaultVal {
				defaultIndex = i
			}
		}

		// DefaultVal is set but does'nt exist in list
		if defaultIndex == -1 {
			// This error message is not for user
			// Should be found while development
			return "", fmt.Errorf("opt.Default is specied but does not exst in list")
		}
	}

	// Construct the query to the user
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("%s\n\n", query))
	for i, item := range list {
		buf.WriteString(fmt.Sprintf("%d. %s\n", i+1, item))
	}
	buf.WriteString("\n")

	// Prompt the message
	fmt.Fprintf(wr, buf.String())

	resultCh := make(chan string, 1)
	errCh := make(chan error, 1)
	go func() {
		// Loop only when error by invalid user input and opts.Loop is true.
		for {
			// Construct the asking line to input
			var buf bytes.Buffer
			buf.WriteString("Enter a number")

			// Add default val if provided
			if defaultIndex >= 0 {
				buf.WriteString(fmt.Sprintf(" (Default is %d)", defaultIndex+1))
			}

			buf.WriteString(": ")
			fmt.Fprintf(wr, buf.String())

			// Read user input from reader.
			var line string
			if _, err := fmt.Fscanln(rd, &line); err != nil {
				// Handle error if it's not `unexpected newline`
				if err.Error() != "unexpected newline" {
					errCh <- fmt.Errorf("failed to read the input: %s", err)
					break
				}
			}

			// execSelect selects a item from list by user input.
			result, err := execSelect(list, line, defaultIndex)
			if err != nil {

				// Don't loop and just return error if Loop is false
				if !opts.Loop {
					errCh <- err
					break
				}

				// Check error and if it's possible to ask again to user
				// then provide appropriate message and run loop again
				switch err {
				case ErrorEmpty:
					fmt.Fprintf(wr, "Input must not be empty. Answer by a number.\n\n")
					continue
				case ErrorNotNumber:
					fmt.Fprintf(wr,
						"%q is not a valid input. Answer by a number.\n\n", line)
					continue
				case ErrorOutOfRange:
					fmt.Fprintf(wr,
						"%q is not a valid choice. Choose from 1 to %d.\n\n",
						line, len(list))
					continue
				default:
					// If other error is returned, it means asking again is
					// impossible
					errCh <- err
					break
				}
			}

			resultCh <- result
			break
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)
	defer signal.Stop(sigCh)

	select {
	case result := <-resultCh:
		// Insert the new line for next output
		fmt.Fprintf(wr, "\n")
		return result, nil
	case err := <-errCh:
		// Insert the new line for next output
		fmt.Fprintf(wr, "\n")
		return "", err
	case <-sigCh:
		// Insert the new line for next output
		fmt.Fprintf(wr, "\n")
		return "", errors.New("interrupted")
	}
}

// execSelect selects a item from list by user input.
// It checks input meets the condition to choose answer from list and if not
// returns appropriate error. See more about error in `error.go` file.
func execSelect(list []string, input string, defaultIndex int) (string, error) {
	if input == "" {
		if defaultIndex >= 0 {
			return list[defaultIndex], nil
		}
		return "", ErrorEmpty
	}

	// Convert user input string to int val
	n, err := strconv.Atoi(input)
	if err != nil {
		return "", ErrorNotNumber
	}

	// Check answer is in range of list
	if n < 1 || len(list) < n {
		return "", ErrorOutOfRange
	}

	return list[n-1], nil
}
