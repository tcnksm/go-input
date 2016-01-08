package input

import (
	"bytes"
	"fmt"
	"os"
	"os/signal"
)

// Ask asks user to input an answer about query. It shows query to user
// and ask input. It returns answer as string. If it catches the SIGINT
// stops reading user input and returns error.
func (i *UI) Ask(query string, opts *Options) (string, error) {

	// Set the default writer & reader if not provided
	wr, rd := i.Writer, i.Reader
	if wr == nil {
		wr = defaultWriter
	}
	if rd == nil {
		rd = defaultReader
	}

	// Construct the query to the user and show it.
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("%s\n", query))
	fmt.Fprintf(wr, buf.String())

	// resultCh is channel receives result string from user input.
	resultCh := make(chan string, 1)

	// errCh is channel receives error while reading user input.
	errCh := make(chan error, 1)

	// sigCh is channel which is watch Interruptted signal (SIGINT)
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)
	defer signal.Stop(sigCh)

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
			if opts.Hide || opts.Mask {

				raw := &Raw{
					Writer: wr,
					Mask:   opts.Mask,
				}

				f, ok := rd.(*os.File)
				if !ok {
					errCh <- fmt.Errorf("reader must be a file")
					return
				}

				var err error
				line, err = raw.Read(f)
				if err != nil {
					errCh <- err
					return
				}
			} else {
				if _, err := fmt.Fscanln(rd, &line); err != nil {
					// Handle error if it's not `unexpected newline`
					if err.Error() != "unexpected newline" {
						errCh <- fmt.Errorf("failed to read the input: %s", err)
						return
					}
				}

			}

			// Use default value if provided
			if line == "" && opts.Default != "" {
				resultCh <- opts.Default
				return
			}

			if line == "" && opts.Required {
				if !opts.Loop {
					errCh <- ErrEmpty
					return
				}

				fmt.Fprintf(wr, "Input must not be empty.\n\n")
				continue
			}

			resultCh <- line
			return
		}
	}()

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
