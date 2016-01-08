package input

import (
	"fmt"
	"io"
	"log"
	"os"
)

// Raw
type Raw struct {
	Writer io.Writer
	Mask   bool
}

// Readline reads the given file one by one
func (rw *Raw) readline(f *os.File) (string, error) {
	var resultBuf []byte
	for {
		var buf [1]byte
		n, err := f.Read(buf[:])
		if err != nil && err != io.EOF {
			log.Fatal(err)
		}

		if n == 0 || buf[0] == '\n' || buf[0] == '\r' {
			break
		}

		if rw.Mask {
			fmt.Fprintf(rw.Writer, "*")
		}
		resultBuf = append(resultBuf, buf[0])
	}

	return string(resultBuf), nil
}
