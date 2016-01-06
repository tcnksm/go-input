package input

import (
	"bytes"
	"fmt"
	"os"
	"os/signal"
)

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
