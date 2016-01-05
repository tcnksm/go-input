package input

import "errors"

var (
	ErrInterrupted = errors.New("interrupted")
	ErrEmpty       = errors.New("default value is not provided but input is empty")
	ErrNotNumber   = errors.New("input must be number")
	ErrOutOfRange  = errors.New("input is out of range")
)
