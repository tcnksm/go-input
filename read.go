package input

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
)

var (
	ErrInterrupted = errors.New("interrupted")
)

type maskOptions struct {
	mask    bool
	maskVal string
}

// read reads input from UI.Reader
func (i *UI) read(opts *maskOptions) (string, error) {
	// sigCh is channel which is watch Interruptted signal (SIGINT)
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)
	defer signal.Stop(sigCh)

	var resultStr string
	var resultErr error
	doneCh := make(chan struct{})

	go func() {
		defer close(doneCh)

		if opts.Mask {
			f, ok := i.Reader.(*os.File)
			if !ok {
				resultErr = fmt.Errorf("reader must be a file")
				return
			}

			i.maskWriter = i.Writer
			i.mask, i.maskVal = opts.mask, opts.maskVal

			resultStr, resultErr = i.rawRead(f)
		} else {
			_, err := fmt.Fscanln(i.Reader, &resultStr)
			if err != nil && err.Error() != "unexpected newline" {
				resultErr = fmt.Errorf("failed to read the input: %s", err)
			}
		}
	}()

	select {
	case <-sigCh:
		return "", ErrInterrupted
	case <-doneCh:
		return resultStr, resultErr
	}
}

// rawReadline tries to return a single line, not including the end-of-line
// bytes with raw Mode (without prompting nothing). Or if provided show some
// value instead of actual value.
func (i *UI) rawReadline(f *os.File) (string, error) {
	var resultBuf []byte
	for {
		var buf [1]byte
		n, err := f.Read(buf[:])
		if err != nil && err != io.EOF {
			return "", err
		}

		if n == 0 || buf[0] == '\n' || buf[0] == '\r' {
			break
		}

		if buf[0] == 3 {
			return "", ErrInterrupted
		}

		if i.mask {
			fmt.Fprintf(i.maskWriter, i.maskVal)
		}

		resultBuf = append(resultBuf, buf[0])
	}

	return string(resultBuf), nil
}
