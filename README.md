go-input
====

[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)][license]
[![Go Documentation](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)][godocs]

[license]: https://github.com/tcnksm/go-input/blob/master/LICENSE
[godocs]: http://godoc.org/github.com/tcnksm/go-input

`go-input` is a Go package for reading user input in console. 

## Install

Use `go get` to install this package:

```bash
$ go get github.com/tcnksm/go-input
```

## Usage

The following is the simple example,

```golang
ui := &input.UI{
    Writer: os.Stdout,
    Reader: os.Stdin,
}

query := "What is your name?"
name, err := ui.Ask(query, &input.Options{
    Default: "tcnksm",
    Required: true,
	Loop:     true,
})
```

You can check other examples in [here](/example).

## Contribution

1. Fork ([https://github.com/tcnksm/go-input/fork](https://github.com/tcnksm/go-input/fork))
1. Create a feature branch
1. Commit your changes
1. Rebase your local changes against the master branch
1. Run test suite with the `go test ./...` command and confirm that it passes
1. Run `gofmt -s`
1. Create new Pull Request

## Author

[Taichi Nakashima](https://github.com/tcnksm)
