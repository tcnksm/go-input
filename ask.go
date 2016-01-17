package input

import (
	"bytes"
	"fmt"
)

// Ask asks the user for input using the given query. The response is
// returned as string. Error is returned based on the given option.
// If the user sends SIGINT (Ctrl+C) while reading input, it catches
// it and return it as a error.
func (i *UI) Ask(query string, opts *Options) (string, error) {
	i.once.Do(i.setDefault)

	// Display the query to the user.
	fmt.Fprintf(i.Writer, "%s\n", query)

	// resultStr and resultErr are return val of this function
	var resultStr string
	var resultErr error
	for {

		// Construct the instruction to user.
		var buf bytes.Buffer
		buf.WriteString("Enter a value")
		if opts.Default != "" {
			buf.WriteString(fmt.Sprintf(" (Default is %s)", opts.Default))
		}

		// Display the instruction to user and ask to input.
		buf.WriteString(": ")
		fmt.Fprintf(i.Writer, buf.String())

		// Read user input from UI.Reader.
		line, err := i.read(opts.readOpts())
		if err != nil {
			resultErr = err
			break
		}

		// line is empty but default is provided returns it
		if line == "" && opts.Default != "" {
			resultStr = opts.Default
			break
		}

		if line == "" && opts.Required {
			if !opts.Loop {
				resultErr = ErrEmpty
				break
			}

			fmt.Fprintf(i.Writer, "Input must not be empty.\n\n")
			continue
		}

		// validate input by custom fuction
		validate := opts.validateFunc()
		if err := validate(line); err != nil {
			if !opts.Loop {
				resultErr = err
				break
			}

			fmt.Fprintf(i.Writer, "Failed to validate input string: %s\n\n", err)
			continue
		}

		// Reach here means it gets ideal input.
		resultStr = line
		break
	}

	// Insert the new line for next output
	fmt.Fprintf(i.Writer, "\n")

	return resultStr, resultErr
}
