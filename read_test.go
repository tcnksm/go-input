package input

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"testing"
)

func TestRead(t *testing.T) {
	var stringWithSpace = fmt.Sprintf("taichi nakashima%s", LineSep)
	var expectedWithSpace = "taichi nakashima"
	var stringNoSpace = fmt.Sprintf("passw0rd%s", LineSep)
	var expectedNoSpace = "passw0rd"

	cases := []struct {
		opts      *readOptions
		userInput io.Reader
		expect    string
	}{
		{
			opts: &readOptions{
				mask:    false,
				maskVal: "",
			},
			userInput: bytes.NewBufferString(stringNoSpace),
			expect:    expectedNoSpace,
		},

		{
			opts: &readOptions{
				mask:    false,
				maskVal: "",
			},
			userInput: bytes.NewBufferString(stringWithSpace),
			expect:    expectedWithSpace,
		},

		// No good way to test masking...
	}

	for i, tc := range cases {
		ui := &UI{
			Writer: ioutil.Discard,
			Reader: tc.userInput,
		}

		out, err := ui.read(tc.opts)
		if err != nil {
			t.Fatalf("#%d expect not to be error: %s", i, err)
		}

		if out != tc.expect {
			t.Fatalf("#%d expect %q to be eq %q", i, out, tc.expect)
		}
	}
}
