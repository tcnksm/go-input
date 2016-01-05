package input

import "errors"

var (
	ErrorEmpty      = errors.New("default value is not provided but input is empty")
	ErrorNotNumber  = errors.New("input must be number")
	ErrorOutOfRange = errors.New("input is out of range")
)
