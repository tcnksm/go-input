# Example of go-input

This directory contains examples to use `go-input`.

To run simplest asking,

```bash
$ go run ask.go
What is your name?
Enter a value: tcnksm
```

To ask password (which you don't want to prompt user input),

```bash
$ go run password.go
What is your password?
Enter a value: *******
```

To ask selecting from list, 

```bash
$ go run select.go
Which language do you prefer to use?

1. go
2. Go
3. golang

Enter a number (Default is 2):
```
