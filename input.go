/*
go-input is package to interact with user input from command line interface.

http://github.com/tcnksm/go-input
*/
package input

import (
	"bytes"
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
	// Default is the value when no thing is innputted.
	Default string

	// Loop continues to asking user to input until getting valid input.
	Loop bool

	// Required returns error when input is empty.
	Required bool

	// Hide hides user input is prompting console.
	Hide bool
}

func (i *UI) Ask(query string, opts *Options) (string, error) {

	// Set the default writer & reader if not provided
	wr, rd := i.Writer, i.Reader
	if wr == nil {
		wr = defaultWriter
	}
	if rd == nil {
		rd = defaultReader
	}

	// Construct the query to the user
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("%s\n", query))
	fmt.Fprintf(wr, buf.String())

	resultCh := make(chan string, 1)
	errCh := make(chan error, 1)
	go func() {
		// Loop only when error by invalid user input and opts.Loop is true.
		for {
			// Construct the asking line to input
			var buf bytes.Buffer
			buf.WriteString("Enter a value")

			// Add default val if provided
			if opts.Default != "" {
				buf.WriteString(fmt.Sprintf(" (Default is %s)", opts.Default))
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

			// Use default value if provided
			if line == "" && opts.Default != "" {
				resultCh <- opts.Default
				return
			}

			resultCh <- line
			return
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
		return "", ErrInterrupted
	}
}

func (i *UI) Select(query string, list []string, opts *Options) (string, error) {

	// Set the default writer & reader if not provided
	wr, rd := i.Writer, i.Reader
	if wr == nil {
		wr = defaultWriter
	}
	if rd == nil {
		rd = defaultReader
	}

	// Find default index which opts.Default indicates
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

	// Prompt the query
	buf.WriteString("\n")
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
					return
				}

				// Check error and if it's possible to ask again to user
				// then provide appropriate message and run loop again
				switch err {
				case ErrEmpty:
					fmt.Fprintf(wr, "Input must not be empty. Answer by a number.\n\n")
					continue
				case ErrNotNumber:
					fmt.Fprintf(wr,
						"%q is not a valid input. Answer by a number.\n\n", line)
					continue
				case ErrOutOfRange:
					fmt.Fprintf(wr,
						"%q is not a valid choice. Choose from 1 to %d.\n\n",
						line, len(list))
					continue
				default:
					// If other error is returned, it means asking again is
					// impossible
					errCh <- err
					return
				}
			}

			resultCh <- result
			return
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
		return "", ErrInterrupted
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
		return "", ErrEmpty
	}

	// Convert user input string to int val
	n, err := strconv.Atoi(input)
	if err != nil {
		return "", ErrNotNumber
	}

	// Check answer is in range of list
	if n < 1 || len(list) < n {
		return "", ErrOutOfRange
	}

	return list[n-1], nil
}
