// +build linux darwin freebsd
package input

import (
	"os"

	"golang.org/x/crypto/ssh/terminal"
)

// RawRead reads file with raw mode (without prompting to terminal).
func (rw *Raw) Read(f *os.File) (string, error) {
	// MakeRaw put the terminal connected to the given file descriptor
	// into raw mode
	fd := int(f.Fd())
	oldState, err := terminal.MakeRaw(fd)
	if err != nil {
		return "", err
	}
	defer terminal.Restore(fd, oldState)

	return rw.readline(f)
}
