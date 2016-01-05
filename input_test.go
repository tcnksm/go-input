package input

import (
	"bytes"
	"io"
	"io/ioutil"
	"testing"
)

func TestAsk(t *testing.T) {
	cases := []struct {
		opts      *Options
		userInput io.Reader
		expect    string
	}{
		{
			opts:      &Options{},
			userInput: bytes.NewBufferString("Taichi\n"),
			expect:    "Taichi",
		},

		{
			opts: &Options{
				Default: "Nakashima",
			},
			userInput: bytes.NewBufferString("\n"),
			expect:    "Nakashima",
		},
	}

	for i, c := range cases {
		ui := &UI{
			Writer: ioutil.Discard,
			Reader: c.userInput,
		}

		ans, err := ui.Ask("", c.opts)
		if err != nil {
			t.Fatalf("#%d expect not to occurr error: %s", i, err)
		}

		if ans != c.expect {
			t.Fatalf("#%d expect %q to be eq %q", i, ans, c.expect)
		}
	}
}

func TestSelect(t *testing.T) {
	cases := []struct {
		list      []string
		opts      *Options
		userInput io.Reader
		expect    string
	}{
		{
			list:      []string{"A", "B", "C"},
			opts:      &Options{},
			userInput: bytes.NewBufferString("1\n"),
			expect:    "A",
		},

		{
			list: []string{"A", "B", "C"},
			opts: &Options{
				Default: "A",
			},
			userInput: bytes.NewBufferString("\n"),
			expect:    "A",
		},

		{
			list: []string{"A", "B", "C"},
			opts: &Options{
				Default: "A",
			},
			userInput: bytes.NewBufferString("3\n"),
			expect:    "C",
		},

		// Loop
		{
			list: []string{"A", "B", "C"},
			opts: &Options{
				Loop: true,
			},
			userInput: bytes.NewBufferString("\n3\n"),
			expect:    "C",
		},

		// Loop
		{
			list: []string{"A", "B", "C"},
			opts: &Options{
				Loop: true,
			},
			userInput: bytes.NewBufferString("\n\n\n\n\n2\n"),
			expect:    "B",
		},

		// Loop
		{
			list: []string{"A", "B", "C"},
			opts: &Options{
				Loop: true,
			},
			userInput: bytes.NewBufferString("4\n3\n"),
			expect:    "C",
		},

		// Loop
		{
			list: []string{"A", "B", "C"},
			opts: &Options{
				Loop: true,
			},
			userInput: bytes.NewBufferString("A\n3\n"),
			expect:    "C",
		},
	}

	for i, c := range cases {
		ui := &UI{
			Writer: ioutil.Discard,
			Reader: c.userInput,
		}

		ans, err := ui.Select("", c.list, c.opts)
		if err != nil {
			t.Fatalf("#%d expect not to occurr error: %s", i, err)
		}

		if ans != c.expect {
			t.Fatalf("#%d expect %q to be eq %q", i, ans, c.expect)
		}
	}
}

func TestExecSelect(t *testing.T) {
	cases := []struct {
		list         []string
		input        string
		defaultIndex int
		expect       string
		err          error
	}{
		{
			list:         []string{"A", "B", "C"},
			input:        "1",
			defaultIndex: -1,
			expect:       "A",
			err:          nil,
		},

		{
			list:         []string{"A", "B", "C"},
			input:        "",
			defaultIndex: 0,
			expect:       "A",
			err:          nil,
		},

		{
			list:         []string{"A", "B", "C"},
			input:        "",
			defaultIndex: -1,
			expect:       "",
			err:          ErrEmpty,
		},

		{
			list:         []string{"A", "B", "C"},
			input:        "A",
			defaultIndex: -1,
			expect:       "",
			err:          ErrNotNumber,
		},

		{
			list:         []string{"A", "B", "C"},
			input:        "80",
			defaultIndex: -1,
			expect:       "",
			err:          ErrOutOfRange,
		},
	}

	for i, c := range cases {
		out, err := execSelect(c.list, c.input, c.defaultIndex)
		if err != c.err {
			t.Fatalf("#%d expect %q to be eq %q", i, err, c.err)
		}

		if c.err != nil {
			continue
		}

		if out != c.expect {
			t.Fatalf("#%d expect %q to be eq %q", i, out, c.expect)
		}
	}

}

func TestSelect_invalidDefault(t *testing.T) {
	ui := &UI{
		Writer: ioutil.Discard,
	}
	_, err := ui.Select("Which?", []string{"A", "B", "C"}, &Options{
		// "D" is not in select target list
		Default: "D",
	})

	if err == nil {
		t.Fatal("expect err to be occurr")
	}
}
